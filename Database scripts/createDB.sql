CREATE TABLE users (
    user_id serial PRIMARY KEY,
    email text UNIQUE
);

CREATE TABLE items (
    item_id serial PRIMARY KEY, 
    user_id integer REFERENCES users,
    name varchar (50) NOT NULL,
    dosage varchar (25),
    taken_today boolean DEFAULT false,
    serving_size smallint DEFAULT 1,
    serving_type varchar(5) check (serving_type in ('scoop', 'pill', 'drop')) DEFAULT 'pill',
    timing time
);

-- Seed Data
INSERT INTO users (email) VALUES ('sachadso@gmail.com');
INSERT INTO users (email) VALUES ('nerv.junker@gmail.com');

-- List all user's items
SELECT item_id, user_id, name, dosage, taken_today, serving_size, serving_type FROM items WHERE user_id = 1;