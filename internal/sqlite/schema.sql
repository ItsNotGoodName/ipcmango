CREATE TABLE settings (
  id INTEGER UNIQUE DEFAULT 0,
  location TEXT NOT NULL,
  latitude REAL NOT NULL,
  longitude REAL NOT NULL,
  sunrise_offset TEXT NOT NULL,
  sunset_offset TEXT NOT NULL,
  sync_video_in_mode BOOL NOT NULL,
  updated_at DATETIME NOT NULL
);

------------
-- Auth
------------
create table goqite (
  id text primary key default ('m_' || lower(hex(randomblob(16)))),
  created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  updated text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  queue text not null,
  body blob not null,
  timeout text not null default (strftime('%Y-%m-%dT%H:%M:%fZ')),
  received integer not null default 0
) strict;

create trigger goqite_updated_timestamp after update on goqite begin
  update goqite set updated = strftime('%Y-%m-%dT%H:%M:%fZ') where id = old.id;
end;

create index goqite_queue_created_idx on goqite (queue, created);

------------
-- Auth
------------
CREATE TABLE users (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  disabled_at DATETIME
);

CREATE TABLE user_sessions (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  ip TEXT NOT NULL,
  last_ip TEXT NOT NULL,
  last_used_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  expired_at DATETIME NOT NULL,

  user_id INTEGER NOT NULL,
  session TEXT NOT NULL,

  FOREIGN KEY (user_id) REFERENCES users (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE tokens (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  ip TEXT NOT NULL,
  last_ip TEXT NOT NULL,
  last_used_at DATETIME NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  expired_at DATETIME NOT NULL,

  token TEXT NOT NULL
);

------------
-- Dahua
------------
CREATE TABLE dahua_devices (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  seed INTEGER NOT NULL UNIQUE,

  name TEXT NOT NULL UNIQUE,
  ip TEXT NOT NULL UNIQUE,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  email TEXT UNIQUE,
  features JSON NOT NULL,

  location TEXT,
  latitude REAL,
  longitude REAL,
  sunrise_offset TEXT,
  sunset_offset TEXT,
  sync_video_in_mode BOOL,

  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL
);

CREATE TABLE dahua_events (
  id TEXT NOT NULL PRIMARY KEY,
  device_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  action TEXT NOT NULL,
  `index` INTEGER NOT NULL,
  data JSON NOT NULL,
  created_at DATETIME NOT NULL,
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_event_rules (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  code TEXT NOT NULL UNIQUE,
  ignore_db BOOLEAN NOT NULL DEFAULT false,
  ignore_live BOOLEAN NOT NULL DEFAULT false,
  ignore_mqtt BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE dahua_event_device_rules (
  device_id INTEGER NOT NULL,
  code TEXT NOT NULL,
  ignore_db BOOLEAN NOT NULL DEFAULT false,
  ignore_live BOOLEAN NOT NULL DEFAULT false,
  ignore_mqtt BOOLEAN NOT NULL DEFAULT false,
  UNIQUE (device_id, code),
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_files (
  id TEXT NOT NULL PRIMARY KEY,
  device_id INTEGER NOT NULL,
  --
  channel INTEGER NOT NULL,
  start_time DATETIME NOT NULL UNIQUE,
  end_time DATETIME NOT NULL,
  length INTEGER NOT NULL,
  type TEXT NOT NULL,
  file_path TEXT NOT NULL,
  duration INTEGER NOT NULL,
  disk INTEGER NOT NULL,
  video_stream TEXT NOT NULL,
  flags JSON NOT NULL,
  events JSON NOT NULL,
  cluster INTEGER NOT NULL,
  partition INTEGER NOT NULL,
  pic_index INTEGER NOT NULL,
  repeat INTEGER NOT NULL,
  work_dir TEXT NOT NULL,
  work_dir_sn BOOLEAN NOT NULL,
  --
  storage TEXT NOT NULL,
  updated_at DATETIME NOT NULL,
  UNIQUE (device_id, file_path),
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE INDEX dahua_files_device_id_start_time_idx ON dahua_files (device_id, start_time);

CREATE TABLE dahua_file_cursors (
  device_id INTEGER NOT NULL UNIQUE,
  quick_cursor DATETIME NOT NULL, -- (scanned) <- quick_cursor -> (not scanned / volatile)
  full_cursor DATETIME NOT NULL, -- (not scanned) <- full_cursor -> (scanned)
  full_epoch DATETIME NOT NULL,
  full_complete BOOLEAN NOT NULL GENERATED ALWAYS AS (full_cursor <= full_epoch) STORED,
  scanning BOOLEAN NOT NULL,
  scan_percent REAL NOT NULL,
  scan_type TEXT NOT NULL,
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_storage_destinations (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL UNIQUE,
  storage TEXT NOT NULL,
  server_address TEXT NOT NULL,
  port INTEGER NOT NULL,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  remote_directory TEXT NOT NULL
);

CREATE TABLE dahua_email_messages (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  device_id INTEGER NOT NULL,
  date DATETIME NOT NULL,
  'from' TEXT NOT NULL,
  `to` JSON NOT NULL,
  subject TEXT NOT NULL,
  `text` TEXT NOT NULL,
  --
  alarm_event TEXT NOT NULL,
  alarm_input_channel INTEGER NOT NULL,
  alarm_name TEXT NOT NULL,
  --
  created_at DATETIME NOT NULL,
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE
);

CREATE TABLE dahua_email_attachments (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL,
  message_id INTEGER,
  file_name TEXT NOT NULL,
  size INTEGER NOT NULL,
  mime_type TEXT NOT NULL,
  FOREIGN KEY (message_id) REFERENCES dahua_email_messages (id) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE dahua_email_endpoints (
  id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  uuid TEXT NOT NULL UNIQUE,
  global BOOLEAN NOT NULL,
  expression TEXT NOT NULL,
  title_template TEXT NOT NULL,
  body_template TEXT NOT NULL,
  attachments BOOLEAN NOT NULL,
  urls JSON NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  disabled_at DATETIME
);

CREATE TABLE dahua_devices_to_email_endpoints (
  device_id INTEGER NOT NULL,
  email_endpoint_id INTEGER NOT NULL,
  UNIQUE (device_id, email_endpoint_id),
  FOREIGN KEY (device_id) REFERENCES dahua_devices (id) ON UPDATE CASCADE ON DELETE CASCADE,
  FOREIGN KEY (email_endpoint_id) REFERENCES dahua_email_endpoints (id) ON UPDATE CASCADE ON DELETE CASCADE
);
