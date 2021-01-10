-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE tweets (
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    comment varchar(255) NOT NULL,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE tweets;
