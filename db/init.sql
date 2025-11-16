-- Vaults Table
CREATE TABLE IF NOT EXISTS vaults (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE
);

-- Collections Table
CREATE TABLE IF NOT EXISTS collections (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    path       TEXT,
    vault_id   INTEGER NOT NULL,
    FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE,
    UNIQUE (name, vault_id)
);

-- Videos Table
CREATE TABLE IF NOT EXISTS videos (
    id             SERIAL PRIMARY KEY,
    title          TEXT NOT NULL,
    slug           TEXT NOT NULL UNIQUE,
    studio         TEXT,
    collection_id  INTEGER NOT NULL,
    vault_id       INTEGER NOT NULL,
    FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    FOREIGN KEY (vault_id) REFERENCES vaults(id) ON DELETE CASCADE
);

-- Tags Table
CREATE TABLE IF NOT EXISTS tags (
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Video Tags Join Table (Many-to-Many Relationship)
CREATE TABLE IF NOT EXISTS video_tags (
    video_id INTEGER NOT NULL,
    tag_id   INTEGER NOT NULL,
    PRIMARY KEY (video_id, tag_id),
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Actors Table
CREATE TABLE IF NOT EXISTS actors (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    slug       TEXT NOT NULL,
    UNIQUE (name, slug)
);

-- Join Table: video_actors (many-to-many)
CREATE TABLE IF NOT EXISTS video_actors (
    video_id   INTEGER NOT NULL,
    actor_id   INTEGER NOT NULL,
    PRIMARY KEY (video_id, actor_id),
    FOREIGN KEY (video_id) REFERENCES videos(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
);
