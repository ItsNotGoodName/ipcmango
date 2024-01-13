-- +goose Up
-- disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- create "new_dahua_afero_files" table
CREATE TABLE `new_dahua_afero_files` (`id` integer NOT NULL PRIMARY KEY AUTOINCREMENT, `file_id` integer NULL, `file_thumbnail_id` integer NULL, `email_attachment_id` integer NULL, `name` text NOT NULL, `ready` boolean NOT NULL DEFAULT false, `size` integer NOT NULL DEFAULT 0, `created_at` datetime NOT NULL, CONSTRAINT `0` FOREIGN KEY (`email_attachment_id`) REFERENCES `dahua_email_attachments` (`id`) ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT `1` FOREIGN KEY (`file_thumbnail_id`) REFERENCES `dahua_file_thumbnails` (`id`) ON UPDATE CASCADE ON DELETE SET NULL, CONSTRAINT `2` FOREIGN KEY (`file_id`) REFERENCES `dahua_files` (`id`) ON UPDATE CASCADE ON DELETE SET NULL);
-- copy rows from old table "dahua_afero_files" to new temporary table "new_dahua_afero_files"
INSERT INTO `new_dahua_afero_files` (`id`, `file_id`, `file_thumbnail_id`, `email_attachment_id`, `name`, `ready`, `created_at`) SELECT `id`, `file_id`, `file_thumbnail_id`, `email_attachment_id`, `name`, `ready`, `created_at` FROM `dahua_afero_files`;
-- drop "dahua_afero_files" table after copying rows
DROP TABLE `dahua_afero_files`;
-- rename temporary table "new_dahua_afero_files" to "dahua_afero_files"
ALTER TABLE `new_dahua_afero_files` RENAME TO `dahua_afero_files`;
-- create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_file_id` ON `dahua_afero_files` (`file_id`);
-- create index "dahua_afero_files_file_thumbnail_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_file_thumbnail_id` ON `dahua_afero_files` (`file_thumbnail_id`);
-- create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_email_attachment_id` ON `dahua_afero_files` (`email_attachment_id`);
-- create index "dahua_afero_files_name" to table: "dahua_afero_files"
CREATE UNIQUE INDEX `dahua_afero_files_name` ON `dahua_afero_files` (`name`);
-- enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;

-- +goose Down
-- reverse: create index "dahua_afero_files_name" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_name`;
-- reverse: create index "dahua_afero_files_email_attachment_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_email_attachment_id`;
-- reverse: create index "dahua_afero_files_file_thumbnail_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_file_thumbnail_id`;
-- reverse: create index "dahua_afero_files_file_id" to table: "dahua_afero_files"
DROP INDEX `dahua_afero_files_file_id`;
-- reverse: create "new_dahua_afero_files" table
DROP TABLE `new_dahua_afero_files`;
