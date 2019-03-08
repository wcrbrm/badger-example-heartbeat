CREATE TABLE ticks (
  id   TEXT PRIMARY KEY, -- should contain hex(ts) + hash(user, and data)
  u    TEXT NOT NULL,    -- user id
  ts   BIGINT NOT NULL,  -- timestamp
  d    TEXT NULL         -- optional - data for the user at this point
);

CREATE TABLE intervals (
  id      TEXT PRIMARY KEY,
  tsStart BIGINT NOT NULL,
  tsEnd   BIGINT NOT NULL,
  g       TEXT NOT NULL,
  u       TEXT NOT NULL,
  d       TEXT NULL
);

CREATE INDEX idx_intervals_from  ON intervals(tsStart);
CREATE INDEX idx_intervals_to    ON intervals(tsEnd);
CREATE INDEX idx_intervals_g     ON intervals(g);
CREATE INDEX idx_intervals_u     ON intervals(u);
