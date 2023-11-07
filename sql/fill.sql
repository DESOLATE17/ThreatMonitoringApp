-- INSERT INTO "users"(login, is_admin, name, password)
-- VALUES ('Ivan123', false, 'Иван Иванов',
--         '781be5ea6620295cbdb249154b840fbe2327d87c666d8e76b29f45f70fcf7d6d');
-- INSERT INTO "users"(login, is_admin, name, password)
-- VALUES ('Dasha2003', true, 'Дарья Такташова',
--         '781be5ea6620295cbdb249154b840fbe2327d87c666d8e76b29f45f70fcf7d6d');

--threats
INSERT INTO threats(name, description, image, count, is_deleted, price)
VALUES ('Фишинг', 'Попытка обманом заставить людей поделиться конфиденциальной информацией, такой как пароли или данные кредитной карты, выдавая себя за заслуживающую доверия организацию посредством электронной почты, сообщений или веб-сайтов.',
        '/image/phishing.jpg', 1850392, false, 30000);
INSERT INTO threats(name, description, image, count, is_deleted, price)
VALUES ( 'Вредоносное ПО', 'Вредоносное программное обеспечение, предназначенное для нарушения работы, повреждения или получения несанкционированного доступа к компьютерным системам, включая вирусы, черви, трояны и программы-вымогатели.',
        '/image/malware.jpeg', 505879385, false, 40000);
INSERT INTO threats(name, description, image, count, is_deleted, price)
VALUES ( 'DDoS Атакаг', 'Перегрузка сети или веб-сайта потоком трафика из нескольких источников, приводящая к его зависанию или сбою.',
        '/image/ddos.jpeg', 384800, false, 70000);
INSERT INTO threats(name, description, image, count, is_deleted, price)
VALUES ( 'Поиск SQL инъекций', 'Использование уязвимостей в базе данных веб-сайта для внедрения вредоносных команд SQL, что потенциально позволяет злоумышленникам получить доступ к конфиденциальным данным или манипулировать ими.',
        '/image/sql_injection.jpeg', 345214, false, 35000);

