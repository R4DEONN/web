USE blog;

CREATE TABLE user
(
    user_id INT NOT NULL AUTO_INCREMENT,
    login VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    PRIMARY KEY (user_id)
) ENGINE=InnoDB
CHARACTER SET = utf8mb4
COLLATE utf8mb4_general_ci
;

