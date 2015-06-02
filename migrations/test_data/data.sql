INSERT INTO users (id, login, name, password, role, mail, active) VALUES 
(700, 'manager', 'SpiderMan', '$2a$10$ru8Moh9lU6ndSvxpCvNU8uvhKDb3LKrNO41vwBjN44vuCZS6m.I1G', 2, 'grr89@mail.ru', true),
(701, 'admin', 'SpongeBob', '$2a$10$QRaH873FNSpjkqX0QnhuFOLfpt2iKGq1EV3WAsanUphUl6z2ZCKHq', 1, 'mirolim777@gmail.com', true),
(702, 'user', 'Vasya', '$2a$10$ru8Moh9lU6ndSvxpCvNU8uvhKDb3LKrNO41vwBjN44vuCZS6m.I1G', 1, 'zak01011996@gmail.com', true);

INSERT INTO departments (id, title) VALUES (700, 'Test department');

INSERT INTO departments (id, title) VALUES (701, 'Test department for criteria 1'), (702, 'Test department for criteria 2'); 

INSERT INTO criteria (id, title, ref_dep) VALUES (700, 'test 1', 701), (701, 'test 2', 701), (702, 'test 3', 701), (703, 'test 4', 702), (704, 'test 5', 702), (705, 'test 6', 702);

INSERT INTO candidates 
       	(id, name, lname, note, active, fname, addr, married, ref_dep, salary, currency, phone, email) 
VALUES 	(700, 'Ivan', 'Ivanov', 'Some long note', true, 'Ivanovich', 'Some boring address', false, 700, 1000, '$', '1235654', 'ex@mail.com'), 
	(701, 'Palonchi', 'Pistonchiev', 'Some long note', true, 'Palonchi ogli', 'Some boring address', false, 701, 800, '$', '1235654', 'ex@mail.com'); 

INSERT INTO ratings (ref_crit, ref_cand, value, ref_user) VALUES (700, 700, 5, 700), (701, 700, 3, 701), (701, 701, 3, 701), (701, 701, 4, 701);

INSERT INTO contacts (id, ref_cand, name, value, active) VALUES (700, 700, 'mail','some@gmail.com', true), (701, 701, 'skype','skypelogin', true) 
