INSERT INTO roles(id, name) VALUES (1, "admin");

INSERT INTO controllers(id, name) VALUES (1, "user");

-- Password: 1234
INSERT INTO usuarios(id, email, password, is_verified, role_id) VALUES (1, "admin@email.com", "$2a$10$uA/rgOrpWK8eWIq5sr6wyu1mDRea6/OBp1HdpFb82U3WDzaLv7bHq", 0, 1); 

-- Admin user rules
-- User controller
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("all", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("all", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("add", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("edit", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("delete", 1, 1, 1);
