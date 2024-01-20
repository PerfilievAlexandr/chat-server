-- +goose Up
-- +goose StatementBegin
create table messages
(
    id         serial primary key,
    text       text      not null,
    producer   text      not null,
    created_at timestamp not null default now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table messages;
-- +goose StatementEnd
