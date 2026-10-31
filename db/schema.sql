CREATE TABLE IF NOT EXISTS wanted_criminals (
    id               BIGSERIAL PRIMARY KEY,
    name             TEXT,                                   -- real name, may be unknown
    nickname         TEXT             NOT NULL,
    reward           BIGINT           NOT NULL CHECK (reward >= 0),
    latitude         DOUBLE PRECISION CHECK (latitude  BETWEEN -90  AND 90),
    longitude        DOUBLE PRECISION CHECK (longitude BETWEEN -180 AND 180),
    last_known_note  TEXT,                                   -- text alternative to coordinates
    difficulty       SMALLINT         NOT NULL CHECK (difficulty BETWEEN 1 AND 10),
    armed            BOOLEAN          NOT NULL DEFAULT FALSE,
    has_crew         BOOLEAN          NOT NULL DEFAULT FALSE,
    escape_master    BOOLEAN          NOT NULL DEFAULT FALSE,
    CHECK ((latitude IS NULL) = (longitude IS NULL))
);

CREATE TABLE IF NOT EXISTS hunt_logs (
    id           BIGSERIAL PRIMARY KEY,
    criminal_id  BIGINT      NOT NULL REFERENCES wanted_criminals (id) ON DELETE CASCADE,
    success      BOOLEAN     NOT NULL,
    notes        TEXT        NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_hunt_logs_criminal_id ON hunt_logs (criminal_id);
