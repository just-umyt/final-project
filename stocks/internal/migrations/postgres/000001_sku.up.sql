CREATE TABLE sku(
    sku_id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    type TEXT,
    price INT,
    location TEXT,
    count INT,
    user_id INT
);

CREATE TABLE cart_sku(
    cart_id INT NOT NULL,
    sku_id TEXT NOT NULL,

    PRIMARY KEY (sku_id)
);
ALTER TABLE cart_sku
ADD FOREIGN KEY (sku_id) REFERENCES cart_sku(sku_id)
ON UPDATE CASCADE ON DELETE CASCADE;

INSERT INTO sku (sku_id, name, type) VALUES
('SKU001', 't-shirt', 'apparel'),
('SKU002', 'cup', 'accessory'),
('SKU003', 'book', 'stationery'),
('SKU004', 'pen', 'stationery'),
('SKU005', 'powerbank', 'electronics'),
('SKU006', 'hoody', 'apparel'),
('SKU007', 'umbrella', 'accessory'),
('SKU008', 'socks', 'apparel'),
('SKU009', 'wallet', 'accessory'),
('SKU010', 'pink-hoody', 'apparel');