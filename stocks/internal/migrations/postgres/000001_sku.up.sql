CREATE TABLE sku(
    sku_id BIGINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT
);

CREATE TABLE stock(
    id  SERIAL NOT NULL PRIMARY KEY,
    sku_id BIGINT UNIQUE, 
    price INT NOT NULL,
    location VARCHAR(255),
    count INT NOT NULL,
    user_id BIGINT
);

ALTER TABLE stock
ADD FOREIGN KEY (sku_id) REFERENCES sku(sku_id)
ON UPDATE CASCADE ON DELETE CASCADE;

INSERT INTO sku (sku_id, name, type) VALUES
(1001, 't-shirt', 'apparel'),
(2020, 'cup', 'accessory'),
(3033, 'book', 'stationery'),
(4044, 'pen', 'stationery'),
(5055, 'powerbank', 'electronics'),
(6066, 'hoody', 'apparel'),
(7077, 'umbrella', 'accessory'),
(8088, 'socks', 'apparel'),
(9099, 'wallet', 'accessory'),
(10101, 'pink-hoody', 'apparel');