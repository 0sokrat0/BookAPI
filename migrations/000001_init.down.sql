-- Удаляем таблицы в обратном порядке для соблюдения зависимостей
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS book_authors;
DROP TABLE IF EXISTS readers;
DROP TABLE IF EXISTS authors;
DROP TABLE IF EXISTS books;
