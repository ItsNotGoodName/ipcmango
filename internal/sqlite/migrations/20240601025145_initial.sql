-- +goose Up
-- create "settings" table
CREATE TABLE `settings` (`id` integer NULL DEFAULT 0, `location` text NOT NULL, `latitude` real NOT NULL, `longitude` real NOT NULL, `sunrise_offset` text NOT NULL, `sunset_offset` text NOT NULL, `sync_video_in_mode` bool NOT NULL, `updated_at` datetime NOT NULL);
-- create index "settings_id" to table: "settings"
CREATE UNIQUE INDEX `settings_id` ON `settings` (`id`);
-- create "users" table
CREATE TABLE `users` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `email` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `disabled_at` datetime NULL);
-- create index "users_email" to table: "users"
CREATE UNIQUE INDEX `users_email` ON `users` (`email`);
-- create index "users_username" to table: "users"
CREATE UNIQUE INDEX `users_username` ON `users` (`username`);
-- create "user_sessions" table
CREATE TABLE `user_sessions` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `user_agent` text NOT NULL, `ip` text NOT NULL, `last_ip` text NOT NULL, `last_used_at` datetime NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `expired_at` datetime NOT NULL, `user_id` integer NOT NULL, `session` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "tokens" table
CREATE TABLE `tokens` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `user_agent` text NOT NULL, `ip` text NOT NULL, `last_ip` text NOT NULL, `last_used_at` datetime NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `expired_at` datetime NOT NULL, `token` text NOT NULL);
-- create "dahua_devices" table
CREATE TABLE `dahua_devices` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `seed` integer NOT NULL, `name` text NOT NULL, `ip` text NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `email` text NULL, `features` json NOT NULL, `location` text NULL, `latitude` real NULL, `longitude` real NULL, `sunrise_offset` text NULL, `sunset_offset` text NULL, `sync_video_in_mode` bool NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "dahua_devices_uuid" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_uuid` ON `dahua_devices` (`uuid`);
-- create index "dahua_devices_seed" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_seed` ON `dahua_devices` (`seed`);
-- create index "dahua_devices_name" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_name` ON `dahua_devices` (`name`);
-- create index "dahua_devices_ip" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_ip` ON `dahua_devices` (`ip`);
-- create index "dahua_devices_email" to table: "dahua_devices"
CREATE UNIQUE INDEX `dahua_devices_email` ON `dahua_devices` (`email`);
-- create "dahua_events" table
CREATE TABLE `dahua_events` (`id` text NOT NULL, `device_id` integer NOT NULL, `code` text NOT NULL, `action` text NOT NULL, `index` integer NOT NULL, `data` json NOT NULL, `created_at` datetime NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_event_rules" table
CREATE TABLE `dahua_event_rules` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `code` text NOT NULL, `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false);
-- create index "dahua_event_rules_code" to table: "dahua_event_rules"
CREATE UNIQUE INDEX `dahua_event_rules_code` ON `dahua_event_rules` (`code`);
-- create "dahua_event_device_rules" table
CREATE TABLE `dahua_event_device_rules` (`device_id` integer NOT NULL, `code` text NOT NULL, `ignore_db` boolean NOT NULL DEFAULT false, `ignore_live` boolean NOT NULL DEFAULT false, `ignore_mqtt` boolean NOT NULL DEFAULT false, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_event_device_rules_device_id_code" to table: "dahua_event_device_rules"
CREATE UNIQUE INDEX `dahua_event_device_rules_device_id_code` ON `dahua_event_device_rules` (`device_id`, `code`);
-- create "dahua_files" table
CREATE TABLE `dahua_files` (`id` text NOT NULL, `device_id` integer NOT NULL, `channel` integer NOT NULL, `start_time` datetime NOT NULL, `end_time` datetime NOT NULL, `length` integer NOT NULL, `type` text NOT NULL, `file_path` text NOT NULL, `duration` integer NOT NULL, `disk` integer NOT NULL, `video_stream` text NOT NULL, `flags` json NOT NULL, `events` json NOT NULL, `cluster` integer NOT NULL, `partition` integer NOT NULL, `pic_index` integer NOT NULL, `repeat` integer NOT NULL, `work_dir` text NOT NULL, `work_dir_sn` boolean NOT NULL, `storage` text NOT NULL, `updated_at` datetime NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_files_start_time" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_start_time` ON `dahua_files` (`start_time`);
-- create index "dahua_files_device_id_file_path" to table: "dahua_files"
CREATE UNIQUE INDEX `dahua_files_device_id_file_path` ON `dahua_files` (`device_id`, `file_path`);
-- create index "dahua_files_device_id_start_time_idx" to table: "dahua_files"
CREATE INDEX `dahua_files_device_id_start_time_idx` ON `dahua_files` (`device_id`, `start_time`);
-- create "dahua_file_cursors" table
CREATE TABLE `dahua_file_cursors` (`device_id` integer NOT NULL, `quick_cursor` datetime NOT NULL, `full_cursor` datetime NOT NULL, `full_epoch` datetime NOT NULL, `full_complete` boolean NOT NULL AS (full_cursor <= full_epoch) STORED, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
CREATE UNIQUE INDEX `dahua_file_cursors_device_id` ON `dahua_file_cursors` (`device_id`);
-- create "dahua_storage_destinations" table
CREATE TABLE `dahua_storage_destinations` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `name` text NOT NULL, `storage` text NOT NULL, `server_address` text NOT NULL, `port` integer NOT NULL, `username` text NOT NULL, `password` text NOT NULL, `remote_directory` text NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL);
-- create index "dahua_storage_destinations_uuid" to table: "dahua_storage_destinations"
CREATE UNIQUE INDEX `dahua_storage_destinations_uuid` ON `dahua_storage_destinations` (`uuid`);
-- create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
CREATE UNIQUE INDEX `dahua_storage_destinations_name` ON `dahua_storage_destinations` (`name`);
-- create "dahua_email_messages" table
CREATE TABLE `dahua_email_messages` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `device_id` integer NOT NULL, `date` datetime NOT NULL, `from` text NOT NULL, `to` json NOT NULL, `subject` text NOT NULL, `text` text NOT NULL, `alarm_event` text NOT NULL, `alarm_input_channel` integer NOT NULL, `alarm_name` text NOT NULL, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create "dahua_email_attachments" table
CREATE TABLE `dahua_email_attachments` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `message_id` integer NULL, `file_name` text NOT NULL, `size` integer NOT NULL, `mime_type` text NOT NULL, CONSTRAINT `0` FOREIGN KEY (`message_id`) REFERENCES `dahua_email_messages` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- create "dahua_email_endpoints" table
CREATE TABLE `dahua_email_endpoints` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `uuid` text NOT NULL, `global` boolean NOT NULL, `expression` text NOT NULL, `title_template` text NOT NULL, `body_template` text NOT NULL, `attachments` boolean NOT NULL, `urls` json NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `disabled_at` datetime NULL);
-- create index "dahua_email_endpoints_uuid" to table: "dahua_email_endpoints"
CREATE UNIQUE INDEX `dahua_email_endpoints_uuid` ON `dahua_email_endpoints` (`uuid`);
-- create "dahua_devices_to_dahua_email_endpoints" table
CREATE TABLE `dahua_devices_to_dahua_email_endpoints` (`device_id` integer NOT NULL, `email_endpoint_id` integer NOT NULL, CONSTRAINT `0` FOREIGN KEY (`email_endpoint_id`) REFERENCES `dahua_email_endpoints` (`id`) ON UPDATE CASCADE ON DELETE CASCADE, CONSTRAINT `1` FOREIGN KEY (`device_id`) REFERENCES `dahua_devices` (`id`) ON UPDATE CASCADE ON DELETE CASCADE);
-- create index "dahua_devices_to_dahua_email_endpoints_device_id_email_endpoint_id" to table: "dahua_devices_to_dahua_email_endpoints"
CREATE UNIQUE INDEX `dahua_devices_to_dahua_email_endpoints_device_id_email_endpoint_id` ON `dahua_devices_to_dahua_email_endpoints` (`device_id`, `email_endpoint_id`);

-- +goose Down
-- reverse: create index "dahua_devices_to_dahua_email_endpoints_device_id_email_endpoint_id" to table: "dahua_devices_to_dahua_email_endpoints"
DROP INDEX `dahua_devices_to_dahua_email_endpoints_device_id_email_endpoint_id`;
-- reverse: create "dahua_devices_to_dahua_email_endpoints" table
DROP TABLE `dahua_devices_to_dahua_email_endpoints`;
-- reverse: create index "dahua_email_endpoints_uuid" to table: "dahua_email_endpoints"
DROP INDEX `dahua_email_endpoints_uuid`;
-- reverse: create "dahua_email_endpoints" table
DROP TABLE `dahua_email_endpoints`;
-- reverse: create "dahua_email_attachments" table
DROP TABLE `dahua_email_attachments`;
-- reverse: create "dahua_email_messages" table
DROP TABLE `dahua_email_messages`;
-- reverse: create index "dahua_storage_destinations_name" to table: "dahua_storage_destinations"
DROP INDEX `dahua_storage_destinations_name`;
-- reverse: create index "dahua_storage_destinations_uuid" to table: "dahua_storage_destinations"
DROP INDEX `dahua_storage_destinations_uuid`;
-- reverse: create "dahua_storage_destinations" table
DROP TABLE `dahua_storage_destinations`;
-- reverse: create index "dahua_file_cursors_device_id" to table: "dahua_file_cursors"
DROP INDEX `dahua_file_cursors_device_id`;
-- reverse: create "dahua_file_cursors" table
DROP TABLE `dahua_file_cursors`;
-- reverse: create index "dahua_files_device_id_start_time_idx" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_start_time_idx`;
-- reverse: create index "dahua_files_device_id_file_path" to table: "dahua_files"
DROP INDEX `dahua_files_device_id_file_path`;
-- reverse: create index "dahua_files_start_time" to table: "dahua_files"
DROP INDEX `dahua_files_start_time`;
-- reverse: create "dahua_files" table
DROP TABLE `dahua_files`;
-- reverse: create index "dahua_event_device_rules_device_id_code" to table: "dahua_event_device_rules"
DROP INDEX `dahua_event_device_rules_device_id_code`;
-- reverse: create "dahua_event_device_rules" table
DROP TABLE `dahua_event_device_rules`;
-- reverse: create index "dahua_event_rules_code" to table: "dahua_event_rules"
DROP INDEX `dahua_event_rules_code`;
-- reverse: create "dahua_event_rules" table
DROP TABLE `dahua_event_rules`;
-- reverse: create "dahua_events" table
DROP TABLE `dahua_events`;
-- reverse: create index "dahua_devices_email" to table: "dahua_devices"
DROP INDEX `dahua_devices_email`;
-- reverse: create index "dahua_devices_ip" to table: "dahua_devices"
DROP INDEX `dahua_devices_ip`;
-- reverse: create index "dahua_devices_name" to table: "dahua_devices"
DROP INDEX `dahua_devices_name`;
-- reverse: create index "dahua_devices_seed" to table: "dahua_devices"
DROP INDEX `dahua_devices_seed`;
-- reverse: create index "dahua_devices_uuid" to table: "dahua_devices"
DROP INDEX `dahua_devices_uuid`;
-- reverse: create "dahua_devices" table
DROP TABLE `dahua_devices`;
-- reverse: create "tokens" table
DROP TABLE `tokens`;
-- reverse: create "user_sessions" table
DROP TABLE `user_sessions`;
-- reverse: create index "users_username" to table: "users"
DROP INDEX `users_username`;
-- reverse: create index "users_email" to table: "users"
DROP INDEX `users_email`;
-- reverse: create "users" table
DROP TABLE `users`;
-- reverse: create index "settings_id" to table: "settings"
DROP INDEX `settings_id`;
-- reverse: create "settings" table
DROP TABLE `settings`;
