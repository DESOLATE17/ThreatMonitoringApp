INSERT INTO "users"(user_id, login, is_admin, name, password_hash, registration_date)
VALUES (1, 'Dasha2003', true, 'Дарья Такташова',
        '781be5ea6620295cbdb249154b840fbe2327d87c666d8e76b29f45f70fcf7d6d', '2023-02-27 19:10');

--threats
INSERT INTO threats(threat_id, name, description, image, count, is_deleted, price)
VALUES (0, 'Фишинг', 'Попытка обманом заставить людей поделиться конфиденциальной информацией, такой как пароли или данные кредитной карты, выдавая себя за заслуживающую доверия организацию посредством электронной почты, сообщений или веб-сайтов.',
        '/image/phishing.jpg', 1850392, false, 30000);
INSERT INTO threats(threat_id, name, description, image, count, is_deleted, price)
VALUES (1, 'Вредоносное ПО', 'Вредоносное программное обеспечение, предназначенное для нарушения работы, повреждения или получения несанкционированного доступа к компьютерным системам, включая вирусы, черви, трояны и программы-вымогатели.',
        '/image/malware.jpeg', 505879385, false, 40000);
INSERT INTO threats(threat_id, name, description, image, count, is_deleted, price)
VALUES (2, 'DDoS Атакаг', 'Перегрузка сети или веб-сайта потоком трафика из нескольких источников, приводящая к его зависанию или сбою.',
        '/image/ddos.jpeg', 384800, false, 70000);
INSERT INTO threats(threat_id, name, description, image, count, is_deleted, price)
VALUES (3, 'Поиск SQL инъекций', 'Использование уязвимостей в базе данных веб-сайта для внедрения вредоносных команд SQL, что потенциально позволяет злоумышленникам получить доступ к конфиденциальным данным или манипулировать ими.',
        '/image/sql_injection.jpeg', 345214, false, 35000);

