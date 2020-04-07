CREATE TABLE users (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    names VARCHAR(64) NOT NULL,
    email VARCHAR (128) NOT NULL,
    password CHAR(60) NOT NULL,

    PRIMARY KEY (id),
    UNIQUE KEY (email)
);

CREATE TABLE products (
    id INT UNSIGNED NOT NULL,
    name VARCHAR (32) NOT NULL,
    description VARCHAR (512) NOT NULL,
    quantity INTEGER DEFAULT 0,
	product_status_id INTEGER NOT NULL,
	regular_price     NUMERIC DEFAULT 0,
	discount_price    NUMERIC DEFAULT 0,
	taxable           BOOLEAN DEFAULT false,
	inserted_at       timestamp  not null,
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

CREATE TABLE tags(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	tag_name VARCHAR(64) NOT NULL,
	inserted_at TIMESTAMP NOT NULL
	
	);


CREATE TABLE product_tags(
	tag_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	product_id INT UNSIGNED NOT NULL,
    time_inserted TIMESTAMP NOT NULL,

    UNIQUE KEY (tag_id, product_id),
    FOREIGN KEY (tag_id)
        REFERENCES tags (id)
        ON DELETE CASCADE,
    FOREIGN KEY (product_id)
        REFERENCES products (id)
        ON DELETE CASCADE
);

CREATE TABLE product_status(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	name VARCHAR(32) NOT NULL,
	product_id INT UNSIGNED NOT NULL,
	
	UNIQUE KEY (product_id),
	FOREIGN KEY (product_id)
        REFERENCES products (id)
       ON DELETE CASCADE


);

CREATE TABLE product_categroy (
	category_id INT UNSIGNED NOT NULL AUTO_INCREMENT,
	category_name VARCHAR(100) NOT NULL,
	product_id INT UNSIGNED NOT NULL,
	
	PRIMARY KEY(category_id),
	FOREIGN KEY(product_id) REFERENCES products(id)
	



);

CREATE TABLE roles(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	role_name VARCHAR(64) NOT NULL,
	inserted_at TIMESTAMP NOT NULL
	
	);
	
	
CREATE TABLE user_roles(
	id INT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
	role_id INT UNSIGNED NOT NULL,
	user_id INT UNSIGNED NOT NULL,
	inserted_at TIMESTAMP NOT NULL,
	
	UNIQUE KEY (role_id, user_id),
    FOREIGN KEY (role_id)
        REFERENCES roles (id)
        ON DELETE CASCADE,
    FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON DELETE CASCADE
	
	);

