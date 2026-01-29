CREATE TABLE friendship_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT,
    receiver_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sender_id) REFERENCES users(id),
    FOREIGN KEY (receiver_id) REFERENCES users(id),
    INDEX (receiver_id) idx_receiver_id,
    INDEX (sender_id) idx_sender_id
)