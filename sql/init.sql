-- Таблица для хранения оригинального запроса

CREATE TABLE IF NOT EXISTS requests
(
    id       SERIAL NOT NULL PRIMARY KEY,
    request     text   NOT NULL
);

-- Таблица для хранения оригинального ответа
CREATE TABLE IF NOT EXISTS responses
(
    id       SERIAL NOT NULL PRIMARY KEY,
    response     text   default ''
);
