# AniDesk

A lightweight, fast, modern desktop client for MyAnimeList built with Go, Wails, Svelte, and SQLite.

## Features

- **OAuth2 PKCE Login** — Secure authentication with MyAnimeList
- **Dashboard** — Seasonal anime, top airing, trending, recommendations
- **Search** — Search anime and manga with filters
- **Anime Details** — Synopsis, ratings, genres, studios, characters
- **Library Management** — Track watching, completed, planning, etc.
- **Offline Mode** — Cache anime data and images for offline access
- **Background Sync** — Changes sync automatically when online
- **Dark/Light Theme** — Custom accent colors
- **Cross-Platform** — Windows, macOS, Linux

## Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go, Wails v2 |
| Frontend | Svelte, TypeScript, TailwindCSS |
| Database | SQLite (via sqlx) |
| API Client | Resty |
| Auth | OAuth2 PKCE |
| Icons | Heroicons |

## Project Structure

```
├── app.go              # Wails app bindings
├── main.go             # Entry point
├── wails.json          # Wails config
├── backend/
│   ├── api/            # MAL API client
│   ├── auth/           # OAuth2 PKCE authentication
│   ├── cache/          # Image caching
│   ├── config/         # App configuration
│   ├── database/       # SQLite setup & migrations
│   ├── models/         # Data models
│   ├── notifications/  # System notifications
│   ├── repository/     # Database access layer
│   ├── services/       # Business logic
│   ├── sync/           # Sync engine for offline edits
│   ├── utils/          # Utilities
│   └── workers/        # Background workers
├── frontend/
│   ├── src/
│   │   ├── components/ # Reusable Svelte components
│   │   ├── pages/      # Route pages
│   │   ├── layouts/    # Layout components
│   │   ├── stores/     # Svelte stores
│   │   ├── lib/        # Utilities & helpers
│   │   └── styles/     # Global styles
│   └── index.html
├── resources/          # App resources
└── build/              # Build artifacts
```

## Development

### Prerequisites

- Go 1.23+
- Node.js 20+
- Wails CLI v2

### Setup

```bash
# Install dependencies
go mod tidy
cd frontend && npm install

# Run in development mode
wails dev

# Build for production
wails build
```

### Configuration

Create `~/.anidesk/config.json`:

```json
{
  "mal_client_id": "YOUR_MAL_CLIENT_ID"
}
```

Get a Client ID at [myanimelist.net/apiconfig](https://myanimelist.net/apiconfig).

## Building for Distribution

```bash
# Build for current platform
wails build

# Build with compression
wails build -upx

# Build for specific platforms (cross-compile)
wails build -platform windows/amd64
wails build -platform darwin/amd64
wails build -platform linux/amd64
```

## Architecture

```
UI (Svelte)
    ↓
Service Layer (Go)
    ↓
Repository Layer (Go)
    ↓
SQLite (Local)  ←→  MAL API (Remote)
```

The application uses a layered architecture:
- **UI** never calls the MAL API directly
- **Service Layer** orchestrates business logic
- **Repository Layer** handles database operations
- **Sync Engine** queues offline changes and syncs them when the internet is available

## License

MIT
