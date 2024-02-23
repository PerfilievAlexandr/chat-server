-- +goose Up
-- +goose StatementBegin
create table history
(
    id         uuid primary key,
    text       text         not null,
    status     varchar(256) not null,
    message_id uuid REFERENCES messages (id),
    created_at timestamp    not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table history;
-- +goose StatementEnd