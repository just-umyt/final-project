CREATE TABLE cart(
    id SERIAL NOT NULL PRIMARY KEY,
    user_id INTEGER NOT NULL ,
    sku_id INTEGER NOT NULL ,
    count INTEGER NOT NULL ,
    UNIQUE(user_id, sku_id)
)
