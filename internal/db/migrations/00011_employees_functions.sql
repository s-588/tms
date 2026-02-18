-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION set_employees_deleted_at_trigger_func()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE employees
    SET deleted_at = now()
    WHERE employee_id = OLD.employee_id;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS set_employees_deleted_at_trigger_func();
-- +goose StatementEnd
