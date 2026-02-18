-- +goose Up
-- +goose StatementBegin
create trigger trigger_orders_updated_at
before update
on orders
for each row
execute function set_updated_at_trigger_func();
-- +goose StatementEnd

-- +goose StatementBegin
create trigger trigger_orders_deleted_at
before delete
on orders
for each row
execute function set_orders_deleted_at_trigger_func();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop trigger trigger_orders_updated_at on orders;
-- +goose StatementEnd

-- +goose StatementBegin
drop trigger trigger_orders_deleted_at on orders;
-- +goose StatementEnd
