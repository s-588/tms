-- +goose Up
-- +goose StatementBegin

-- First, ensure we have the correct number of employees and get their IDs
-- Insert 50 employees
INSERT INTO employees (name, created_at) VALUES
('David Chen', '2022-01-10 08:00:00'),
('Priya Sharma', '2022-01-15 09:00:00'),
('Marcus Johnson', '2022-02-01 10:00:00'),
('Sophie Williams', '2022-02-15 11:00:00'),
('Carlos Rodriguez', '2022-03-01 08:30:00'),
('Emma Wilson', '2022-03-15 09:30:00'),
('Kenji Tanaka', '2022-04-01 10:30:00'),
('Aisha Mohammed', '2022-04-15 11:30:00'),
('Peter Schmidt', '2022-05-01 08:45:00'),
('Lisa Anderson', '2022-05-15 09:45:00'),
('Ahmed Hassan', '2022-06-01 10:45:00'),
('Maria Garcia', '2022-06-15 11:45:00'),
('James Miller', '2022-07-01 08:15:00'),
('Sarah Davis', '2022-07-15 09:15:00'),
('Michael Brown', '2022-08-01 10:15:00'),
('Jennifer Taylor', '2022-08-15 11:15:00'),
('Robert Jones', '2022-09-01 08:20:00'),
('Jessica Moore', '2022-09-15 09:20:00'),
('Thomas Martin', '2022-10-01 10:20:00'),
('Emily Thompson', '2022-10-15 11:20:00'),
('Daniel White', '2022-11-01 08:25:00'),
('Amanda Harris', '2022-11-15 09:25:00'),
('Christopher Clark', '2022-12-01 10:25:00'),
('Ashley Lewis', '2022-12-15 11:25:00'),
('Joshua Lee', '2023-01-01 08:35:00'),
('Megan Walker', '2023-01-15 09:35:00'),
('Andrew Hall', '2023-02-01 10:35:00'),
('Samantha Young', '2023-02-15 11:35:00'),
('Brian King', '2023-03-01 08:40:00'),
('Victoria Scott', '2023-03-15 09:40:00'),
('Kevin Green', '2023-04-01 10:40:00'),
('Nicole Adams', '2023-04-15 11:40:00'),
('Jason Baker', '2023-05-01 08:50:00'),
('Rachel Nelson', '2023-05-15 09:50:00'),
('Eric Carter', '2023-06-01 10:50:00'),
('Hannah Mitchell', '2023-06-15 11:50:00'),
('Steven Perez', '2023-07-01 08:55:00'),
('Brittany Roberts', '2023-07-15 09:55:00'),
('Ryan Turner', '2023-08-01 10:55:00'),
('Kelly Phillips', '2023-08-15 11:55:00'),
('Jacob Campbell', '2023-09-01 08:05:00'),
('Lauren Parker', '2023-09-15 09:05:00'),
('Tyler Evans', '2023-10-01 10:05:00'),
('Molly Edwards', '2023-10-15 11:05:00'),
('Scott Collins', '2023-11-01 08:10:00'),
('Vanessa Stewart', '2023-11-15 09:10:00'),
('Justin Sanchez', '2023-12-01 10:10:00'),
('Tiffany Morris', '2023-12-15 11:10:00'),
('Brandon Rogers', '2024-01-01 08:12:00'),
('Christina Reed', '2024-01-15 09:12:00');
select * from employees;

-- Insert 15 fuel types
INSERT INTO fuels (name, supplier, price, created_at) VALUES
('Diesel', 'ExxonMobil', '3.89', '2022-01-01 00:00:00'),
('Gasoline Regular', 'Shell', '3.45', '2022-01-01 00:00:00'),
('Gasoline Premium', 'BP', '4.15', '2022-01-01 00:00:00'),
('Bio-Diesel B20', 'Chevron', '3.95', '2022-01-01 00:00:00'),
('Compressed Natural Gas', 'Total', '2.15', '2022-01-01 00:00:00'),
('Electricity', 'Tesla Supercharger', '0.28', '2022-01-01 00:00:00'),
('Propane', 'AmeriGas', '2.75', '2022-01-01 00:00:00'),
('Ethanol E85', 'Archer Daniels', '3.25', '2022-01-01 00:00:00'),
('Jet Fuel', 'World Fuel', '4.85', '2022-01-01 00:00:00'),
('Marine Diesel', 'Valero', '4.15', '2022-01-01 00:00:00'),
('Hydrogen', 'Air Products', '16.50', '2022-01-01 00:00:00'),
('Kerosene', 'Phillips 66', '4.05', '2022-01-01 00:00:00'),
('Aviation Gasoline', 'ExxonMobil', '6.15', '2022-01-01 00:00:00'),
('Biodiesel B100', 'Renewable Energy', '4.35', '2022-01-01 00:00:00'),
('Liquefied Natural Gas', 'Shell LNG', '2.45', '2022-01-01 00:00:00');

-- Insert 200 clients
INSERT INTO clients (name, email, email_verified, phone, created_at) VALUES
('John Smith', 'john.smith@example.com', true, '+15551234567', '2023-01-15 09:30:00'),
('Maria Garcia', 'maria.garcia@example.com', true, '+15552345678', '2023-01-16 10:15:00'),
('Robert Johnson', 'robert.j@example.com', false, '+15553456789', '2023-01-17 14:20:00'),
('Sarah Miller', 'sarah.miller@example.com', true, '+15554567890', '2023-01-18 11:45:00'),
('James Davis', 'james.davis@example.com', true, '+15555678901', '2023-01-19 16:30:00'),
('Lisa Anderson', 'lisa.anderson@example.com', false, '+15556789012', '2023-01-20 08:15:00'),
('Michael Wilson', 'michael.w@example.com', true, '+15557890123', '2023-01-21 13:40:00'),
('Emily Taylor', 'emily.taylor@example.com', true, '+15558901234', '2023-01-22 15:25:00'),
('David Brown', 'david.brown@example.com', true, '+15559012345', '2023-01-23 09:50:00'),
('Jennifer Lee', 'jennifer.lee@example.com', false, '+15550123456', '2023-01-24 12:10:00'),
('Christopher Clark', 'chris.clark@example.com', true, '+15551122334', '2023-02-01 10:00:00'),
('Amanda Hall', 'amanda.hall@example.com', true, '+15552233445', '2023-02-02 11:30:00'),
('Daniel White', 'daniel.white@example.com', false, '+15553344556', '2023-02-03 14:45:00'),
('Jessica Harris', 'jessica.h@example.com', true, '+15554455667', '2023-02-04 16:20:00'),
('Matthew Martin', 'matt.martin@example.com', true, '+15555566778', '2023-02-05 09:15:00'),
('Ashley Thompson', 'ashley.t@example.com', true, '+15556677889', '2023-02-06 13:40:00'),
('Joshua Martinez', 'josh.m@example.com', false, '+15557788990', '2023-02-07 15:50:00'),
('Megan Robinson', 'megan.r@example.com', true, '+15558899001', '2023-02-08 08:30:00'),
('Andrew Lewis', 'andrew.lewis@example.com', true, '+15559900112', '2023-02-09 12:25:00'),
('Samantha Walker', 'samantha.w@example.com', false, '+15551011123', '2023-02-10 14:35:00'),
('Thomas Young', 'thomas.y@example.com', true, '+15552122234', '2023-03-01 10:30:00'),
('Olivia Allen', 'olivia.allen@example.com', true, '+15553233345', '2023-03-02 11:45:00'),
('William King', 'william.king@example.com', false, '+15554344456', '2023-03-03 14:20:00'),
('Sophia Scott', 'sophia.scott@example.com', true, '+442071234567', '2023-03-04 16:10:00'),
('Benjamin Green', 'ben.green@example.com', true, '+442082345678', '2023-03-05 09:40:00'),
('Isabella Baker', 'isabella.b@example.com', false, '+442093456789', '2023-03-06 13:15:00'),
('Ethan Adams', 'ethan.adams@example.com', true, '+442074567890', '2023-03-07 15:30:00'),
('Mia Nelson', 'mia.nelson@example.com', true, '+442085678901', '2023-03-08 08:45:00'),
('Alexander Carter', 'alex.carter@example.com', false, '+442096789012', '2023-03-09 12:20:00'),
('Charlotte Mitchell', 'charlotte.m@example.com', true, '+442017890123', '2023-03-10 14:50:00');

-- Insert price configurations for different cargo types
INSERT INTO prices (cargo_type, cost, weight, distance, created_at) VALUES
('Electronics', '2.50', 10, 100, '2022-01-01 00:00:00'),
('Furniture', '4.25', 50, 100, '2022-01-01 00:00:00'),
('Clothing', '1.75', 5, 100, '2022-01-01 00:00:00'),
('Food Perishable', '3.50', 20, 100, '2022-01-01 00:00:00'),
('Machinery', '6.00', 100, 100, '2022-01-01 00:00:00'),
('Automotive Parts', '3.25', 30, 100, '2022-01-01 00:00:00'),
('Construction Materials', '5.50', 150, 100, '2022-01-01 00:00:00'),
('Chemicals', '8.75', 25, 100, '2022-01-01 00:00:00'),
('Medical Supplies', '4.50', 15, 100, '2022-01-01 00:00:00'),
('Documents', '1.25', 1, 100, '2022-01-01 00:00:00');

-- Insert 30 orders
INSERT INTO orders (distance, weight, total_price, status, created_at) VALUES
(150, 25, '375.00', 'delivered', '2023-01-05 08:30:00'),
(320, 120, '960.00', 'in_transit', '2023-01-06 14:15:00'),
(75, 8, '112.50', 'delivered', '2023-01-07 11:20:00'),
(450, 200, '2250.00', 'processing', '2023-01-08 16:45:00'),
(220, 45, '990.00', 'delivered', '2023-01-09 09:10:00'),
(180, 65, '1170.00', 'in_transit', '2023-01-10 13:30:00'),
(90, 12, '135.00', 'delivered', '2023-01-11 10:45:00'),
(600, 150, '3600.00', 'pending', '2023-01-12 15:20:00'),
(280, 80, '2240.00', 'delivered', '2023-01-13 08:55:00'),
(120, 30, '360.00', 'in_transit', '2023-01-14 12:10:00'),
(350, 95, '3325.00', 'processing', '2023-01-15 14:40:00'),
(200, 55, '1100.00', 'delivered', '2023-01-16 09:25:00'),
(420, 180, '3780.00', 'in_transit', '2023-01-17 16:15:00'),
(95, 15, '142.50', 'delivered', '2023-01-18 11:05:00'),
(550, 220, '6050.00', 'pending', '2023-01-19 13:50:00'),
(160, 35, '560.00', 'delivered', '2023-01-20 08:40:00'),
(380, 110, '4180.00', 'in_transit', '2023-01-21 15:30:00'),
(110, 22, '242.00', 'delivered', '2023-01-22 10:20:00'),
(480, 190, '4560.00', 'processing', '2023-01-23 14:10:00'),
(240, 75, '1800.00', 'delivered', '2023-01-24 09:15:00'),
(300, 85, '2550.00', 'in_transit', '2023-01-25 12:45:00'),
(130, 28, '364.00', 'delivered', '2023-01-26 08:50:00'),
(520, 210, '5460.00', 'pending', '2023-01-27 16:40:00'),
(170, 40, '680.00', 'delivered', '2023-01-28 11:35:00'),
(400, 130, '5200.00', 'in_transit', '2023-01-29 14:25:00'),
(100, 18, '180.00', 'delivered', '2023-01-30 09:45:00'),
(580, 230, '6670.00', 'processing', '2023-01-31 15:55:00'),
(190, 50, '950.00', 'delivered', '2023-02-01 10:30:00'),
(340, 100, '3400.00', 'in_transit', '2023-02-02 13:20:00'),
(140, 32, '448.00', 'delivered', '2023-02-03 08:35:00');

-- Insert 30 transport vehicles (with valid employee_id references)
-- Now employee_id values 1-30 exist because we inserted 50 employees above
INSERT INTO transports (employee_id, model, license_plate, payload_capacity, fuel_id, fuel_consumption, created_at) VALUES
(1, 'Ford F-150', 'ABC123', 1500, 2, 15, '2022-02-01 00:00:00'),
(2, 'Mercedes Sprinter', 'DEF456', 3500, 1, 12, '2022-02-02 00:00:00'),
(3, 'Volvo FH16', 'GHI789', 20000, 1, 25, '2022-02-03 00:00:00'),
(4, 'Tesla Semi', 'JKL012', 18000, 6, 0, '2022-02-04 00:00:00'),
(5, 'Isuzu NPR', 'MNO345', 6000, 2, 18, '2022-02-05 00:00:00'),
(6, 'Freightliner Cascadia', 'PQR678', 25000, 1, 28, '2022-02-06 00:00:00'),
(7, 'Chevrolet Express', 'STU901', 4000, 2, 16, '2022-02-07 00:00:00'),
(8, 'Peterbilt 389', 'VWX234', 28000, 1, 30, '2022-02-08 00:00:00'),
(9, 'Ram 3500', 'YZA567', 5500, 2, 17, '2022-02-09 00:00:00'),
(10, 'Kenworth T680', 'BCD890', 24000, 1, 26, '2022-02-10 00:00:00'),
(11, 'GMC Sierra 3500', 'EFG123', 6000, 2, 18, '2022-02-11 00:00:00'),
(12, 'International LT', 'HIJ456', 22000, 1, 27, '2022-02-12 00:00:00'),
(13, 'Ford Transit', 'KLM789', 4000, 2, 14, '2022-02-13 00:00:00'),
(14, 'Mack Anthem', 'NOP012', 26000, 1, 29, '2022-02-14 00:00:00'),
(15, 'Toyota Tacoma', 'QRS345', 1200, 2, 19, '2022-02-15 00:00:00'),
(16, 'Western Star 5700', 'TUV678', 30000, 1, 32, '2022-02-16 00:00:00'),
(17, 'Nissan NV3500', 'WXY901', 3500, 2, 16, '2022-02-17 00:00:00'),
(18, 'Volvo VNL', 'ZAB234', 21000, 1, 24, '2022-02-18 00:00:00'),
(19, 'Chevrolet Silverado 2500', 'CDE567', 3500, 2, 17, '2022-02-19 00:00:00'),
(20, 'Freightliner M2', 'FGH890', 12000, 1, 20, '2022-02-20 00:00:00'),
(21, 'Ford E-350', 'IJK123', 4500, 2, 15, '2022-02-21 00:00:00'),
(22, 'Peterbilt 579', 'LMN456', 23000, 1, 26, '2022-02-22 00:00:00'),
(23, 'Ram 2500', 'OPQ789', 3000, 2, 18, '2022-02-23 00:00:00'),
(24, 'Kenworth T880', 'RST012', 32000, 1, 33, '2022-02-24 00:00:00'),
(25, 'GMC Savana', 'UVW345', 3800, 2, 16, '2022-02-25 00:00:00'),
(26, 'International HV', 'XYZ678', 19000, 1, 25, '2022-02-26 00:00:00'),
(27, 'Mercedes Actros', '123ABC', 27000, 1, 28, '2022-02-27 00:00:00'),
(28, 'Ford F-250', '456DEF', 2800, 2, 17, '2022-02-28 00:00:00'),
(29, 'Volvo FE', '789GHI', 11000, 1, 22, '2022-03-01 00:00:00'),
(30, 'Chevrolet Colorado', '012JKL', 1500, 2, 20, '2022-03-02 00:00:00');

-- Insert orders_transport relationships
INSERT INTO orders_transport (order_id, transport_id) VALUES
(1, 1),
(1, 2),
(2, 3),
(3, 4),
(4, 5),
(5, 6),
(6, 7),
(7, 8),
(8, 9),
(9, 10),
(10, 11),
(11, 12),
(12, 13),
(13, 14),
(14, 15),
(15, 16),
(16, 17),
(17, 18),
(18, 19),
(19, 20),
(20, 21),
(21, 22),
(22, 23),
(23, 24),
(24, 25),
(25, 26),
(26, 27),
(27, 28),
(28, 29),
(29, 30),
(30, 1);

-- Insert clients_orders relationships
INSERT INTO clients_orders (client_id, order_id) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5),
(6, 6),
(7, 7),
(8, 8),
(9, 9),
(10, 10),
(11, 11),
(12, 12),
(13, 13),
(14, 14),
(15, 15),
(16, 16),
(17, 17),
(18, 18),
(19, 19),
(20, 20),
(21, 21),
(22, 22),
(23, 23),
(24, 24),
(25, 25),
(26, 26),
(27, 27),
(28, 28),
(29, 29),
(30, 30);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Note: The order matters here due to foreign key constraints
DELETE FROM clients_orders;
DELETE FROM orders_transport;
DELETE FROM transports;
DELETE FROM orders;
DELETE FROM prices;
DELETE FROM fuels;
DELETE FROM employees;
DELETE FROM clients;

-- Reset sequences if needed (optional, but good practice)
ALTER SEQUENCE clients_client_id_seq RESTART;
ALTER SEQUENCE employees_employee_id_seq RESTART;
ALTER SEQUENCE fuels_fuel_id_seq RESTART;
ALTER SEQUENCE orders_order_id_seq RESTART;
ALTER SEQUENCE prices_price_id_seq RESTART;
ALTER SEQUENCE transports_transport_id_seq RESTART;
-- +goose StatementEnd
