CREATE DATABASE device_inventory;
\c device_inventory;
CREATE TABLE IF NOT EXISTS devices(
    hostname TEXT
);