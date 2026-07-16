package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/oauth2"

	"anidesk/backend/config"
	"anidesk/backend/models"
)

const (
	authURL  = "https://myanimelist.net/v1/oauth2/authorize"
	tokenURL = "https://myanimelist.net/v1/oauth2/token"
)

type Authenticator struct {
	config    *config.Config
	oauthCfg  *oauth2.Config
	verifier  string
}

func New(cfg *config.Config) *Authenticator {
	return &Authenticator{
		config: cfg,
		oauthCfg: &oauth2.Config{
			ClientID:    cfg.MALClientID,
			Scopes:      []string{"write:users"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  authURL,
				TokenURL: tokenURL,
			},
			RedirectURL: fmt.Sprintf("http://localhost:%d/callback", cfg.RedirectPort),
		},
	}
}

func (a *Authenticator) generatePKCE() (string, string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	verifier := base64.RawURLEncoding.EncodeToString(b)

	h := sha256.Sum256([]byte(verifier))
	challenge := base64.RawURLEncoding.EncodeToString(h[:])

	return verifier, challenge, nil
}

func (a *Authenticator) Login(ctx context.Context) (*models.User, error) {
	verifier, challenge, err := a.generatePKCE()
	if err != nil {
		return nil, fmt.Errorf("failed to generate PKCE: %w", err)
	}
	a.verifier = verifier

	authURL := a.oauthCfg.AuthCodeURL("state",
		oauth2.SetAuthURLParam("code_challenge", challenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	tokenCh := make(chan *oauth2.Token, 1)
	errCh := make(chan error, 1)

	server := &http.Server{Addr: fmt.Sprintf(":%d", a.config.RedirectPort)}
	handler := http.NewServeMux()
	handler.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			errCh <- fmt.Errorf("no code in callback")
			w.WriteHeader(400)
			io.WriteString(w, "Authorization failed: no code received.")
			return
		}

		token, err := a.exchangeCode(code)
		if err != nil {
			errCh <- err
			w.WriteHeader(400)
			io.WriteString(w, "Token exchange failed.")
			return
		}

		tokenCh <- token
		w.WriteHeader(200)
		io.WriteString(w, "Authorization successful! You can close this window.")
	})
	server.Handler = handler

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.RedirectPort))
	if err != nil {
		return nil, fmt.Errorf("failed to start callback server: %w", err)
	}

	go server.Serve(listener)
	defer server.Close()

	runtime.BrowserOpenURL(ctx, authURL)

	select {
	case token := <-tokenCh:
		return a.tokenToUser(token), nil
	case err := <-errCh:
		return nil, err
	case <-time.After(5 * time.Minute):
		return nil, fmt.Errorf("login timeout")
	}
}

func (a *Authenticator) exchangeCode(code string) (*oauth2.Token, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	data := map[string]string{
		"client_id":     a.config.MALClientID,
		"code":          code,
		"code_verifier": a.verifier,
		"grant_type":    "authorization_code",
		"redirect_uri":  fmt.Sprintf("http://localhost:%d/callback", a.config.RedirectPort),
	}

	body, _, err := a.doFormPost(client, tokenURL, data)
	if err != nil {
		return nil, err
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	return &oauth2.Token{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		Expiry:       time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second),
	}, nil
}

func (a *Authenticator) RefreshToken(ctx context.Context, user *models.User) error {
	client := &http.Client{Timeout: 10 * time.Second}

	data := map[string]string{
		"client_id":     a.config.MALClientID,
		"refresh_token": user.RefreshToken,
		"grant_type":    "refresh_token",
	}

	body, _, err := a.doFormPost(client, tokenURL, data)
	if err != nil {
		return fmt.Errorf("token refresh failed: %w", err)
	}

	var tokenResp struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("failed to parse refresh response: %w", err)
	}

	user.AccessToken = tokenResp.AccessToken
	if tokenResp.RefreshToken != "" {
		user.RefreshToken = tokenResp.RefreshToken
	}
	user.TokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)

	return nil
}

func (a *Authenticator) doFormPost(client *http.Client, reqURL string, data map[string]string) ([]byte, int, error) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}

	resp, err := client.Post(reqURL, "application/x-www-form-urlencoded", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	return body, resp.StatusCode, nil
}

func (a *Authenticator) tokenToUser(token *oauth2.Token) *models.User {
	return &models.User{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenExpiry:  token.Expiry,
	}
}

func (a *Authenticator) Logout() {
	a.verifier = ""
}
