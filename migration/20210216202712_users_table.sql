-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE users (
    id int NOT NULL PRIMARY KEY AUTO_INCREMENT,
    email varchar(255) NOT NULL UNIQUE,
    created_at datetime,
    updated_at datetime,
    deleted_at datetime
);

ALTER TABLE tweets ADD user_id int not null;

ALTER TABLE tweets ADD CONSTRAINT fk_user_id foreign key (user_id) references users (id) ON DELETE RESTRICT ON UPDATE RESTRICT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
