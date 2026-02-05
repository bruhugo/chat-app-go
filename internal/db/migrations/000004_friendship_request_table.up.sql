CREATE TABLE friendship_requests (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    sender_id BIGINT NOT NULL,
    receiver_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    CONSTRAINT fk_sender
        FOREIGN KEY (sender_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_receiver
        FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,

    UNIQUE KEY uniq_friend_request (sender_id, receiver_id),
    INDEX idx_receiver_id (receiver_id),
    INDEX idx_sender_id (sender_id)
) ENGINE=InnoDB;
