CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    price VARCHAR(255),
    location VARCHAR(255),
    description VARCHAR(255),
    breed VARCHAR(255)
)