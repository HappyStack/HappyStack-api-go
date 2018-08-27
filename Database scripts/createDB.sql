CREATE TABLE item (
    item_id serial PRIMARY KEY,
    name varchar (50) NOT NULL,
    dosage varchar (25),
    taken_today boolean DEFAULT false,
    serving_size smallint DEFAULT 1,
    serving_type varchar(5) check (serving_type in ('scoop', 'pill', 'drop')) DEFAULT 'pill',
    timing time
);
