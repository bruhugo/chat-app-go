CREATE TABLE chats (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50),
    description VARCHAR(200),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    creator_id BIGINT NOT NULL,
    INDEX idx_chats_creator_id (creator_id),
    FOREIGN KEY (creator_id) REFERENCES users(id)
) 