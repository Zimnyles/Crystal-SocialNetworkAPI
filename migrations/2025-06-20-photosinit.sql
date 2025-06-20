CREATE TABLE photos (
    id SERIAL PRIMARY KEY,              
    user_id INTEGER NOT NULL,        
    created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,  
    file_path VARCHAR(255) NOT NULL,     
    is_public BOOLEAN DEFAULT TRUE,    
    
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);