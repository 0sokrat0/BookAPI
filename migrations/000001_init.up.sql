-- Создание таблицы books
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    year INT,
    isbn VARCHAR,
    genre VARCHAR
);

-- Создание таблицы authors
CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    country VARCHAR
);

-- Создание таблицы readers
CREATE TABLE readers (
    id SERIAL PRIMARY KEY,
    name VARCHAR NOT NULL,
    phone VARCHAR,
    email VARCHAR,
    password VARCHAR,
    admin BOOLEAN
);

-- Создание таблицы book_authors (связующая таблица между books и authors)
CREATE TABLE book_authors (
    book_id INT NOT NULL,
    author_id INT NOT NULL,
    PRIMARY KEY (book_id, author_id),
    FOREIGN KEY (book_id) REFERENCES books(id),
    FOREIGN KEY (author_id) REFERENCES authors(id)
);

-- Создание таблицы reservations
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL,
    reader_id INT NOT NULL,
    start_date DATE,
    end_date DATE,
    FOREIGN KEY (book_id) REFERENCES books(id),
    FOREIGN KEY (reader_id) REFERENCES readers(id)
);
