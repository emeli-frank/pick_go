CREATE TABLE users (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    names VARCHAR(64) NOT NULL,
    email VARCHAR (128) NOT NULL,
    password CHAR(60) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY (email)
);

CREATE TABLE products (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR (32) NOT NULL,
    description VARCHAR (512) NOT NULL,
    quantity INTEGER DEFAULT 0,
	regular_price     NUMERIC DEFAULT 0,
	discount_price    NUMERIC DEFAULT 0,

	PRIMARY KEY (id)
);

CREATE TABLE cart_items (
    user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_id INT UNSIGNED NOT NULL,

    UNIQUE KEY (user_id, product_id),
    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,
    FOREIGN KEY (product_id)
        REFERENCES products (id)
        ON DELETE CASCADE
);

CREATE TABLE order_history (
    user_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_id INT UNSIGNED NOT NULL,
    time_ordered DATETIME NOT NULL,

    UNIQUE KEY (user_id, product_id),
    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE,
    FOREIGN KEY (product_id)
        REFERENCES products (id)
        ON DELETE CASCADE
);
