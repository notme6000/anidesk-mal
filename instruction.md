You are an experienced Go software engineer, UI/UX designer, and desktop application architect.

Your task is to build a production-quality cross-platform desktop application for MyAnimeList using Go and Wails.

==========================
PROJECT DETAILS
==========================

Project Name:
AniDesk

Description:
AniDesk is a lightweight, fast, modern desktop client for MyAnimeList.

It should feel like a native desktop application instead of a web page wrapped inside Electron.

The application must be cleanly architected, scalable, maintainable, and well documented.

The code should follow Go best practices.

Never generate spaghetti code.

==========================
TECH STACK
==========================

Backend
- Go (latest stable)
- Wails
- SQLite
- sqlx (preferred) or GORM
- Resty HTTP client
- OAuth2 PKCE
- Go routines
- Context package

Frontend
- Svelte
- TypeScript
- TailwindCSS
- Vite
- Heroicons
- Svelte Stores

Packaging
- Windows
- Linux
- macOS

==========================
PROJECT STRUCTURE
==========================

Organize the project like this

backend/

    api/

    auth/

    cache/

    database/

    models/

    repository/

    services/

    sync/

    workers/

    notifications/

    config/

    utils/

frontend/

    src/

        components/

        pages/

        layouts/

        stores/

        lib/

        assets/

        styles/

resources/

build/

==========================
ARCHITECTURE
==========================

Use a layered architecture.

UI

↓

Service Layer

↓

Repository Layer

↓

Database

↓

MyAnimeList API

Never call the API directly from UI components.

==========================
MAIN FEATURES
==========================

Authentication

- OAuth Login
- OAuth Logout
- Token Refresh
- Secure Token Storage

Dashboard

Show

- Continue Watching
- Seasonal Anime
- Top Airing
- Trending
- Recommendations
- Recently Updated

Anime

Each anime page should display

- Poster
- Banner
- Synopsis
- Rating
- Studios
- Genres
- Related Anime
- Characters
- Voice Actors
- Recommendations
- Reviews
- Statistics

Search

Search by

- Anime
- Manga
- Character
- Studio

Filters

- Year
- Genre
- Rating
- Status
- Airing
- Season

Library

Categories

Watching

Completed

Dropped

On Hold

Plan To Watch

User can

Update score

Update progress

Add notes

Favorite

Remove

Settings

Dark Mode

Light Mode

Accent Colors

Notifications

Cache Settings

About Page

==========================
OFFLINE MODE
==========================

Cache all viewed anime.

Cache posters.

Cache banners.

Cache user library.

Cache search history.

Application should continue working offline.

Automatically synchronize changes once internet is available.

==========================
BACKGROUND WORKERS
==========================

Implement independent goroutines for

Image Downloader

Background Synchronizer

Cache Cleaner

Update Checker

Notification Scheduler

Workers must communicate safely.

==========================
DATABASE
==========================

SQLite

Tables

users

anime_cache

manga_cache

library

favorites

history

settings

sync_queue

images

==========================
IMAGE CACHE
==========================

Images should be downloaded once.

Store under

cache/images/

Never redownload existing images.

Load images locally whenever possible.

==========================
API CLIENT
==========================

Create a reusable API package.

Functions

SearchAnime()

SearchManga()

GetAnime()

GetUser()

GetSeason()

GetRecommendations()

UpdateAnimeStatus()

UpdateAnimeScore()

No duplicated HTTP code.

==========================
SYNC ENGINE
==========================

Maintain a queue.

Local edits

↓

SQLite

↓

Queue

↓

Background Sync

↓

MyAnimeList

Retry failed updates.

Never lose data.

==========================
UI DESIGN
==========================

Use modern minimal UI.

Design inspiration

- Steam
- Spotify
- Arc Browser
- Zen Browser
- Discord
- Fluent Design
- Material 3

Rounded cards.

Smooth animations.

Blur effects where appropriate.

Responsive layout.

No clutter.

Dark mode first.

==========================
HOME PAGE
==========================

Hero banner

Continue Watching

Seasonal Anime

Trending

Recommendations

Recent Activity

Statistics

==========================
ANIME PAGE
==========================

Large banner

Poster

Title

Score

Episodes

Description

Genres

Characters

Reviews

Recommendations

Related Anime

==========================
STATE MANAGEMENT
==========================

Use Svelte Stores.

Stores

UserStore

AnimeStore

ThemeStore

SettingsStore

SearchStore

==========================
ERROR HANDLING
==========================

Proper logging.

Toast notifications.

Retry failed requests.

Graceful error messages.

==========================
PERFORMANCE
==========================

Lazy loading.

Infinite scrolling.

Pagination.

Efficient image cache.

Minimal RAM usage.

Fast startup.

==========================
SECURITY
==========================

Use OAuth PKCE.

Never expose secrets.

Validate all inputs.

Sanitize API responses.

Secure token storage.

==========================
CODE QUALITY
==========================

Requirements

- Modular code
- SOLID principles
- Dependency Injection where useful
- Interfaces
- Small reusable functions
- Proper comments
- Unit tests
- Meaningful variable names
- Idiomatic Go
- Type-safe TypeScript

==========================
DEVELOPMENT PROCESS
==========================

Do NOT generate the entire project in one response.

Instead work incrementally.

Each milestone must

1. Explain what will be built.

2. Generate complete code.

3. Explain the code.

4. Wait for approval.

Milestones

1.
Initialize Wails project.

2.
Setup Tailwind + Svelte.

3.
Create routing.

4.
Build layout.

5.
Implement OAuth.

6.
Create API client.

7.
SQLite integration.

8.
Dashboard.

9.
Search.

10.
Anime page.

11.
Library.

12.
Sync engine.

13.
Caching.

14.
Settings.

15.
Notifications.

16.
Testing.

17.
Packaging.

==========================
DOCUMENTATION
==========================

Generate

README.md

Installation Guide

Developer Guide

Architecture Diagram

API Documentation

Folder Documentation

==========================
GOAL
==========================

The final application should be polished enough to be released as an open-source desktop client for MyAnimeList.

Every decision should prioritize:

- Speed
- Simplicity
- Native desktop feel
- Clean architecture
- Low memory usage
- Maintainability
- Scalability

Do not skip any implementation details. Build the project step by step and ensure each feature is complete before moving to the next milestone.
