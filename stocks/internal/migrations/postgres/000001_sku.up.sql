CREATE TABLE sku(
    sku_id BIGINT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT,
    price INT,
    location TEXT,
    count INT,
    user_id INT
);

CREATE TABLE cart_sku(
    cart_id BIGINT NOT NULL,
    sku_id BIGINT NOT NULL,

    PRIMARY KEY (sku_id, cart_id)
);
ALTER TABLE cart_sku
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