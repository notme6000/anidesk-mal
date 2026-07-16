<div align="center">
  <br/>
  <img src="frontend/src/assets/images/logo-main.png" alt="AniDesk" width="128" />
  <h1>AniDesk</h1>
  <p><strong>A lightweight, native desktop client for MyAnimeList</strong></p>
  <p>
    Built with Go + Wails + Svelte + SQLite — feels like a native app, not a website wrapped in Electron.
  </p>
  <br/>
</div>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go" alt="Go" />
  <img src="https://img.shields.io/badge/Wails-2.12-DF2B1B?style=flat" alt="Wails" />
  <img src="https://img.shields.io/badge/Svelte-3-FF3E00?style=flat&logo=svelte" alt="Svelte" />
  <img src="https://img.shields.io/badge/license-MIT-blue?style=flat" alt="MIT" />
  <img src="https://img.shields.io/badge/platform-Windows%20|%20macOS%20|%20Linux-lightgrey?style=flat" alt="Platform" />
</p>

---

## ✨ Features

| | |
|---|---|
| 🔐 **OAuth2 PKCE Login** | Secure authentication — no password sharing, no API key exposure |
| 🏠 **Dashboard** | Seasonal anime, top airing, trending, and recommendations at a glance |
| 🔍 **Search** | Search anime and manga with genre/season/rating filters |
| 📄 **Anime Details** | Synopsis, score, rank, genres, studios, episodes, characters |
| 📚 **Library Management** | Track watching, completed, on-hold, dropped, plan-to-watch |
| 📴 **Offline Mode** | Anime data and images cached locally — works without internet |
| 🔄 **Auto Sync** | Queue changes offline, sync automatically when connection returns |
| 🌙 **Dark / Light** | Full theme support with custom accent colors |
| 🖥️ **Cross-Platform** | Native feel on Windows, macOS, and Linux |

## 🚀 Quick Start

### Prerequisites

- Go 1.23+
- Node.js 20+
- Wails CLI v2 (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)

### Setup

```bash
# Clone & enter
git clone https://github.com/yourusername/anidesk.git
cd anidesk

# Install dependencies
go mod tidy
cd frontend && npm install && cd ..

# Run in development mode (hot reload)
wails dev

# Build for production
wails build
```

### Configuration

1. Go to [myanimelist.net/apiconfig](https://myanimelist.net/apiconfig)
2. Create a new application
3. Set **Redirect URI** to `http://localhost:43829/callback`
4. Copy the **Client ID**

Then launch AniDesk, go to **Settings**, paste your Client ID, and click **Save**. Click **Sign In** in the sidebar to authorize.

<details>
<summary><strong>Or configure manually</strong> (via config file)</summary>

```json
// ~/.anidesk/config.json
{
  "mal_client_id": "YOUR_CLIENT_ID_HERE"
}
```

</details>

## 📸 Screenshots

> *Screenshots coming soon.*

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────┐
│                     UI (Svelte)                         │
│  Components · Pages · Layouts · Stores                  │
└──────────────────────┬──────────────────────────────────┘
                       │  Wails Bindings
                       ▼
┌─────────────────────────────────────────────────────────┐
│                 Service Layer (Go)                      │
│  AuthService · AnimeService · LibraryService · Settings │
└──────────────────────┬──────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────┐
│                Repository Layer (Go)                    │
│  UserRepo · AnimeRepo · LibraryRepo · SyncQueueRepo     │
└──────┬───────────────────────────────────┬──────────────┘
       │                                   │
       ▼                                   ▼
┌──────────────┐              ┌──────────────────────────┐
│   SQLite     │              │   MyAnimeList API (v2)   │
│  (Local DB)  │◄────────────►│   api.myanimelist.net    │
└──────────────┘   (Sync)     └──────────────────────────┘
```

### Key Design Decisions

- **Layered architecture** — UI never calls the MAL API directly
- **Offline-first** — All data is cached in SQLite; changes sync via a queue
- **Background workers** — Image downloader, cache cleaner, sync engine run as independent goroutines
- **OAuth PKCE** — Industry-standard secure auth for desktop apps
- **No Electron** — Wails uses the native WebView, keeping memory usage low (≈30MB)

## 📁 Project Structure

```
anidesk/
├── main.go                  # Entry point
├── app.go                   # Wails app bindings (exposed to frontend)
├── wails.json               # Wails configuration
├── backend/
│   ├── api/                 # MyAnimeList REST client (Resty)
│   ├── auth/                # OAuth2 PKCE authentication flow
│   ├── cache/               # Image downloader & cache manager
│   ├── config/              # App configuration (~/.anidesk/config.json)
│   ├── database/            # SQLite setup, migrations, schema
│   ├── models/              # Data models (Go structs with JSON/DB tags)
│   ├── notifications/       # System toast notifications
│   ├── repository/          # Database access layer (sqlx)
│   ├── services/            # Business logic layer
│   ├── sync/                # Offline sync engine with retry queue
│   ├── utils/               # Utility functions
│   └── workers/             # Background goroutine workers
├── frontend/
│   ├── src/
│   │   ├── components/      # Reusable Svelte components
│   │   ├── pages/           # Route pages (Dashboard, Search, etc.)
│   │   ├── layouts/         # Sidebar, MainLayout
│   │   ├── stores/          # Svelte stores (auth, theme, search, etc.)
│   │   ├── lib/             # Icon component, helpers
│   │   └── styles/          # TailwindCSS base styles
│   └── index.html
├── resources/               # App resources (icons, etc.)
├── build/                   # Build artifacts
└── README.md
```

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| **Backend** | [Go](https://go.dev/) + [Wails v2](https://wails.io/) |
| **Frontend** | [Svelte 3](https://svelte.dev/) + [TypeScript](https://www.typescriptlang.org/) |
| **Styling** | [TailwindCSS 3](https://tailwindcss.com/) |
| **Database** | [SQLite](https://www.sqlite.org/) via [sqlx](https://github.com/jmoiron/sqlx) |
| **HTTP Client** | [Resty v2](https://github.com/go-resty/resty) |
| **Auth** | OAuth2 PKCE ([golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2)) |
| **Icons** | [Heroicons](https://heroicons.com/) |
| **Router** | [svelte-spa-router](https://github.com/ItalyPaleAle/svelte-spa-router) |

## 🧪 Running Tests

```bash
go test ./...
```

## 📦 Building for Distribution

```bash
# Current platform
wails build

# With UPX compression
wails build -upx

# Cross-compile
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

## 🔒 Security

- OAuth2 PKCE flow — no client secret required, no password exposure
- Tokens stored in local SQLite database (not exposed to frontend JS)
- Never stores your MAL password
- All API calls go through the Go backend, not the browser
- Input validation on all search queries

## 📄 License

MIT — see [LICENSE](LICENSE) for details.

---

<p align="center">
  Made with ❤️ for anime fans
</p>
