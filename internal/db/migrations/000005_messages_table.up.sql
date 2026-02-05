CREATE TABLE messages (
    id BIGINT AUTO_INCREMENT PRIMARY KEY, 
    content VARCHAR(2000) NOT NULL,
    user_id BIGINT NOT NULL,
    chat_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (chat_id) REFERENCES chats(id),
    INDEX idx_user_id (user_id),
    INDEX idx_chat_id (chat_id)
)