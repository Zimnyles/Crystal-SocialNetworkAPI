CREATE TABLE feedposts (
    id SERIAL PRIMARY KEY,
    login VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    image_path VARCHAR(255) NOT NULL DEFAULT 'no_image',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);