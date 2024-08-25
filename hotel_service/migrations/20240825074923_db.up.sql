CREATE TABLE IF NOT EXISTS hotels(
    id SERIAL PRIMARY KEY,
    name VARCHAR(250),
    location VARCHAR(250),
    rating INT,
    address VARCHAR(250)
);

CREATE TABLE IF NOT EXISTS rooms (
    id SERIAL PRIMARY KEY,
    hotel_id INT,
    room_type VARCHAR(100),
    price_per_night FLOAT,
    availability BOOLEAN DEFAULT TRUE
);