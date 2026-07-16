package notifications

import (
	"context"
	"log"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Manager struct {
	ctx context.Context
}

func New(ctx context.Context) *Manager {
	return &Manager{ctx: ctx}
}

func (m *Manager) Show(title, message string) {
	runtime.LogInfo(m.ctx, title+": "+message)
	runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
		Title:   title,
		Message: message,
	})
}

func (m *Manager) Error(title, message string) {
	runtime.LogError(m.ctx, title+": "+message)
	runtime.MessageDialog(m.ctx, runtime.MessageDialogOptions{
		Title:   title,
		Message: message,
		Type:    runtime.ErrorDialog,
	})
}

var _ = log.Println
