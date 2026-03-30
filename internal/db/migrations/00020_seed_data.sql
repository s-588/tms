-- +goose Up

-- transports
-- +goose StatementBegin
INSERT INTO transports (model, license_plate, payload_capacity, fuel_consumption, created_at) VALUES
('МАЗ-6430', '7890 OP-2', 22000, 32, NOW()),
('Volvo FH16', '1234 QR-3', 28000, 36, NOW()),
('Scania R450', '5678 ST-4', 24000, 31, NOW()),
('MAN TGS', '9012 UV-5', 26000, 33, NOW()),
('Mercedes Actros', '3456 WX-6', 27000, 34, NOW()),
('DAF XF', '7890 YZ-7', 25000, 30, NOW()),
('Iveco Stralis', '2345 AB-8', 21000, 29, NOW()),
('Renault T', '6789 CD-9', 20000, 28, NOW()),
('КАМАЗ-54901', '1234 EF-0', 23000, 31, NOW()),
('ГАЗель NN', '5678 GH-1', 2000, 13, NOW());
-- +goose StatementEnd

-- insurances
-- +goose StatementBegin
INSERT INTO insurances (transport_id, insurance_date, insurance_expiration, payment, coverage, created_at) VALUES
((SELECT transport_id FROM transports WHERE license_plate = '7890 OP-2'), '2026-03-20', '2027-03-20', 1250.00, 52000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '1234 QR-3'), '2026-03-22', '2027-03-22', 1550.00, 75000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '5678 ST-4'), '2026-03-25', '2027-03-25', 1350.00, 62000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '9012 UV-5'), '2026-04-01', '2027-04-01', 1450.00, 68000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '3456 WX-6'), '2026-04-05', '2027-04-05', 1500.00, 70000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '7890 YZ-7'), '2026-04-10', '2027-04-10', 1400.00, 65000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '2345 AB-8'), '2026-04-15', '2027-04-15', 1150.00, 48000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '6789 CD-9'), '2026-04-20', '2027-04-20', 1100.00, 45000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '1234 EF-0'), '2026-04-25', '2027-04-25', 1300.00, 55000.00, NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '5678 GH-1'), '2026-05-01', '2027-05-01', 850.00, 32000.00, NOW());
-- +goose StatementEnd

-- inspections (дополнительные, с датами после 18.03.2026)
-- +goose StatementBegin
INSERT INTO inspections (transport_id, inspection_date, inspection_expiration, status, created_at) VALUES
((SELECT transport_id FROM transports WHERE license_plate = '7890 OP-2'), '2026-03-18', '2027-03-18', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '1234 QR-3'), '2026-03-19', '2026-09-19', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '5678 ST-4'), '2026-03-20', '2026-06-20', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '9012 UV-5'), '2026-03-21', '2026-12-21', 'repair', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '3456 WX-6'), '2026-03-22', '2027-03-22', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '7890 YZ-7'), '2026-03-23', '2026-10-23', 'overdue', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '2345 AB-8'), '2026-03-24', '2026-08-24', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '6789 CD-9'), '2026-03-25', '2026-07-25', 'ready', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '1234 EF-0'), '2026-03-26', '2026-11-26', 'repair', NOW()),
((SELECT transport_id FROM transports WHERE license_plate = '5678 GH-1'), '2026-03-27', '2026-05-27', 'ready', NOW());
-- +goose StatementEnd

-- clients
-- +goose StatementBegin
INSERT INTO clients (name, email, phone, created_at) VALUES
('Иванов Иван', 'ivan.ivanov@example.by', '+375291234567', NOW()),
('Петров Пётр', 'petr.petrov@example.by', '+375331234567', NOW()),
('Сидорова Ольга', 'olga.sidorova@example.by', '+375441234567', NOW()),
('Козлов Андрей', 'andrei.kozlov@example.by', '+375291112233', NOW()),
('Новикова Елена', 'elena.novikova@example.by', '+375331445566', NOW()),
('Морозов Дмитрий', 'dmitry.morozov@example.by', '+375441778899', NOW());
-- +goose StatementEnd

-- prices
-- +goose StatementBegin
INSERT INTO prices (cargo_type, weight, distance, created_at) VALUES
('Сыпучие',              1.05, 1.10, NOW()),
('Жидкости',             1.15, 1.25, NOW()),
('Строительные материалы',1.10, 1.05, NOW()),
('Оборудование',         1.40, 1.60, NOW()),
('Продукты',             1.10, 1.05, NOW()),
('Мебель',               1.20, 1.15, NOW()),
('Химикаты',             1.35, 1.50, NOW()),
('Стекло',               1.45, 1.55, NOW()),
('Металл',               1.05, 1.00, NOW()),
('Древесина',            1.05, 1.05, NOW()),
('Сельхозпродукция',     1.00, 1.00, NOW()),
('Бытовая техника',      1.25, 1.20, NOW()),
('Одежда',               1.15, 1.20, NOW()),
('Автозапчасти',         1.15, 1.10, NOW()),
('Керамика',             1.30, 1.30, NOW());
-- +goose StatementEnd

-- employees
-- +goose StatementBegin
INSERT INTO employees (name, status, job_title, hire_date, salary, license_issued, license_expiration, created_at) VALUES
('Петренко Алексей', 'available', 'driver', '2025-01-15', 1500.00, '2020-05-10', '2030-05-10', NOW()),
('Васильева Мария', 'available', 'dispatcher', '2024-03-01', 1200.00, '2021-02-20', '2031-02-20', NOW()),
('Кузнецов Николай', 'available', 'mechanic', '2023-07-22', 1300.00, '2019-11-01', '2029-11-01', NOW()),
('Леонович Павел', 'assigned', 'driver', '2025-10-10', 1600.00, '2022-01-10', '2032-01-10', NOW()),
('Савицкая Анна', 'available', 'logistics_manager', '2024-05-05', 2000.00, '2020-08-15', '2030-08-15', NOW()),
('Ковалёв Игорь', 'unavailable', 'mechanic', '2022-12-01', 1400.00, '2018-04-20', '2028-04-20', NOW());

-- Update status for some employees (optional, default is available)
UPDATE employees SET status = 'assigned' WHERE name = 'Леонович Павел';
UPDATE employees SET status = 'unavailable' WHERE name = 'Ковалёв Игорь';
-- +goose StatementEnd

-- nodes
-- +goose StatementBegin
INSERT INTO nodes (address, name, geom, created_at) VALUES
('пр-т Независимости 1, Минск', 'Минск-Центр', Point(27.5615, 53.9025), NOW()),
('ул. Тимирязева 65, Минск', 'Северный', Point(27.5312, 53.9389), NOW()),
('ул. Махновича 2, Брест', 'Брест-Авто', Point(23.6877, 52.0964), NOW()),
('ул. Советская 1, Гомель', 'Гомель-Главный', Point(30.9973, 52.4412), NOW()),
('пр-т Черняховского 10, Витебск', 'Витебск-Северный', Point(30.2181, 55.1904), NOW()),
('ул. Первомайская 25, Могилёв', 'Могилёв-Центр', Point(30.3322, 53.9135), NOW()),
('ул. Ожешко 37, Гродно', 'Гродно-Южный', Point(23.8288, 53.6785), NOW());
-- +goose StatementEnd

-- orders are intentionally omitted as per request

-- +goose Down

-- +goose StatementBegin
DELETE FROM orders;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM inspections;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM insurances;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM nodes;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM prices;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM employees;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM clients;
-- +goose StatementEnd

-- +goose StatementBegin
DELETE FROM transports;
-- +goose StatementEnd