INSERT INTO roles(id, name) VALUES (1, "admin");

INSERT INTO controllers(id, name) VALUES (1, "user");

-- Admin user rules
-- User controller
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("all", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("all", 0, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("add", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("edit", 1, 1, 1);
INSERT INTO rules(action, permission, role_id, controller_id) VALUES ("delete", 0, 1, 1);
