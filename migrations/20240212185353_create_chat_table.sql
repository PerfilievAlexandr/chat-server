-- +goose Up
-- +goose StatementBegin
create table chat
(
    id         uuid primary key,
    owner      text      not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table chat;
-- +goose StatementEnd
