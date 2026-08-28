# Reelix Go
Self hosting media server similar to Jellyfin & Plex.

## Technology Stack
- [Go](https://go.dev/)
- [Postgres](https://www.postgresql.org/)
- [Nginx](https://nginx.org/)
- [Docker](https://www.docker.com/)

## Requirements
- Install [Docker/Docker Desktop](https://www.docker.com/products/docker-desktop/).
- Install [VS Code](https://code.visualstudio.com/) and install the [Dev Containers](https://github.com/Microsoft/vscode-remote-release) extension in VS Code.
> **Note for Windows Users:** It's recommended installing and using [WSL Ubuntu](https://learn.microsoft.com/en-us/windows/wsl/install) alongside the WSL extension in VS Code.

## Getting Started w/ Docker Compose

### 1. Setup Environment Variables.
Create .env file with the following values:

- ROOT_PATH (The root location of where your vaults are stored)
- POSTGRES_DB
- POSTGRES_USER
- POSTGRES_PASSWORD
- JWT_SECRET (Used for generating JWTs)

### 2. Start Containers.
```bash
docker-compose up --build
```

## Running Tests

### 1. Start Test Container
```bash
docker-compose -f compose.test.yml up --build
```

### 2. Run Tests
```bash
go test -tags=test ./...
```

```bash
go test -tags=test -v -count=1 ./...
```


## Setting Up Media
Each folder in your root path are considered a vault. Each vault should have a *collections/* and *pictures/* folder with an optional *cover.jpg* to represent the image of the vault.

The *collections/* folder should only contain folders, where each folder will represent a collection for the vault. Each collection folder can then contain folders, where each folder will represent a video for that collection.

The *pictures/* folder can contain an *actors/* folder in addition to other folders with unique names to represent a collection of images. Each image in the *actors/* folder should be in the format of *firstname_lastname.jpg* while other folders the image should be in the format of *001.jpg*, *002.jpg*, etc.

Each video folder can contain the following files:
- *video_title.mp4*
- *video_title.nfo*
- *thumbnail.jpg*
- *backdrop.jpg*

The *.mp4* and *.nfo* file should be the same name, ideally lowercase and `_` used for spaces.

An example of what your file structure could look like:

```
/reelix
  /NASA
    /collections
      /artemis_ii
        /launch_coverage
          launch_coverage.mp4
          launch_coverage.nfo
          thumbnail.jpg
    /pictures
      /artemis_ii
        001.jpg
        002.jpg
        003.jpg
```

The following is an example of what your *.nfo* file could look like:

```xml
<movie>
    <title>Launch Coverage</title>
    <studio>NASA</studio>
    <tag>Artemis II</tag>
    <plot>
        Watch the replay of the Artemis II launch! 
    </plot>

    <actor>
        <name>Reid Wiseman</name>
        <type>Astronaut</type>
    </actor>
    <actor>
        <name>Victor Glover</name>
        <type>Astronaut</type>
    </actor>
    <actor>
        <name>Christina Koch</name>
        <type>Astronaut</type>
    </actor>
    <actor>
        <name>Jeremy Hansen</name>
        <type>Astronaut</type>
    </actor>
</movie>
```

## Project Background
This codebase consists of 3 main parts:
- File Scanner
- PostgreSQL database
- API server

### File Scanner
Walks a structured directory tree rooted at `ROOT_PATH` to discover vaults, collections, videos, galleries, and actor images. Parses `.nfo` XML metadata files to extract title, studio, tags, and actors for each video. Uses `os.ReadDir` and `path/filepath` for filesystem traversal, `encoding/xml` for NFO parsing. Syncs scanned data to PostgreSQL via idempotent upsert operations.

### PostgreSQL Database
PostgreSQL 15 accessed via `pgx/v5` with `pgxpool` connection pooling. Core tables include `vaults`, `collections`, `videos`, `tags`, `video_tags`, `actors`, `video_actors`, `galleries`, and `users`. Uses `INSERT ... ON CONFLICT DO UPDATE` for idempotent writes, `UNNEST()` for bulk inserts, and transactions for atomic video/tag/actor creation. JSON aggregation is used to fetch videos with related tags and actors in single queries.

### API Server
HTTP routing via `gorilla/mux` with JWT authentication (HS256) supporting 15-minute access tokens and 7-day refresh tokens stored in HTTP-only cookies. Password hashing uses `bcrypt` with cost factor 14.

**Public Endpoints (no authentication required):**
- `GET /api/status` - Health check, returns `{"status": "OK"}`
- `POST /api/register` - Register a new user (first registered user becomes admin)
- `POST /api/login` - Authenticate with username/password, sets JWT cookies on success
- `POST /api/refresh` - Refresh expired access token using refresh token cookie
- `POST /api/logout` - Logout user (clears authentication cookies)

**Protected Endpoints (JWT authentication required):**
- `GET /api/vaults` - List all vaults
- `GET /api/vault/{vaultId}` - Get a single vault by ID
- `GET /api/collections/{vaultId}` - List all collections within a vault
- `GET /api/collection/{collectionId}` - Get a single collection by ID
- `GET /api/videos/{collectionId}` - List all videos within a collection
- `GET /api/video/{videoId}` - Get a single video with its associated tags and actors
- `GET /api/galleries/{vaultId}` - List all galleries within a vault
- `GET /api/gallery/{galleryId}` - Get a single gallery by ID
- `GET /api/actors/{vaultId}` - List all actors associated with videos in a vault
- `GET /api/actor/{actorId}` - Get a single actor by ID