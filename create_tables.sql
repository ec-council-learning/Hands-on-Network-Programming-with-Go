CREATE DATABASE device_inventory;
\c device_inventory;
CREATE TABLE IF NOT EXISTS vendors(
    id SERIAL PRIMARY KEY,
    name TEXT
);
CREATE TABLE IF NOT EXISTS models(
    id SERIAL PRIMARY KEY,
    vendor_id INT,
    name TEXT,
    FOREIGN KEY (vendor_id) REFERENCES vendors(id)
);
CREATE TABLE IF NOT EXISTS devices(
    id SERIAL PRIMARY KEY,
    hostname TEXT,
    ipv4 INET,
    model_id INT,
    FOREIGN KEY (model_id) REFERENCES models(id)
);

-- seed db
INSERT INTO vendors (name)
VALUES ('juniper');
INSERT INTO vendors (name)
VALUES ('cisco');
INSERT INTO models (name, vendor_id)
VALUES ('srx210', 1);
INSERT INTO models (name, vendor_id)
VALUES ('7206', 2);
INSERT INTO devices (hostname, ipv4, model_id)
VALUES ('labsrx', '10.0.0.60', 1);

-- join
SELECT
    hostname,
    ipv4,
    vendors.name AS vendor,
    models.name AS model
FROM devices
JOIN models
    ON devices.model_id = models.id
JOIN vendors
    ON models.vendor_id = vendors.id;

-- DSN (data source name)
-- export DSN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/device_inventory