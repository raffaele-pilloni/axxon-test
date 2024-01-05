CREATE TABLE `axxon-test`.`task` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `status` VARCHAR(255) NOT NULL,
    `method` VARCHAR(255) NOT NULL,
    `url` VARCHAR(255) NOT NULL,
    `request_headers` JSON NOT NULL,
    `request_body` JSON NOT NULL,
    `response_status_code` VARCHAR(255) NULL,
    `response_headers` JSON NULL,
    `response_content_length` INT NULL,
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
);