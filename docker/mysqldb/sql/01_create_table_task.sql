# Schema of table task
# ------------------------------------------------------------

CREATE TABLE `axxon_test`.`task` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `status` VARCHAR(255) NOT NULL,
    `method` VARCHAR(255) NOT NULL,
    `url` VARCHAR(255) NOT NULL,
    `request_headers` JSON NOT NULL,
    `request_body` JSON NOT NULL,
    `response_headers` JSON NULL,
    `response_status_code` INT NULL,
    `response_content_length` INT NULL,
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    KEY `IDX_STATUS_CREATED_AT` (`status`, `created_at`)
    PRIMARY KEY (`id`)
);