CREATE DATABASE device_inventory_dev;
\c device_inventory_dev

CREATE TABLE IF NOT EXISTS devices(
    id SERIAL PRIMARY KEY,
    hostname TEXT,
    ipv4 INET
);

INSERT INTO devices (hostname, ipv4)
VALUES ('labsrx', '10.0.0.60');
INSERT INTO devices (hostname, ipv4)
VALUES ('labpi', '10.0.0.61');
INSERT INTO devices (hostname, ipv4)
VALUES ('throwaway', '1.1.1.1');

-- grabs all rows
SELECT id, hostname, ipv4 FROM devices;

-- grab a specific row
SELECT id, hostname, ipv4 FROM devices WHERE hostname = 'labsrx';

-- update a specific row
UPDATE devices SET ipv4 = '2.2.2.2' WHERE id = 3;

-- delete a row
DELETE FROM devices where id = 3;

-- delete a table
DROP TABLE devices;

-- delete a database
\c postgres;
DROP DATABASE device_inventory_dev;