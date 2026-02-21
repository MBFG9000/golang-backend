CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20) UNIQUE,
    city    VARCHAR(100),
    country VARCHAR(100),
    zip     VARCHAR(20),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

INSERT INTO users (first_name, last_name, email, phone, city, country, zip)
VALUES
    ('Александр', 'Иванов',   'a.ivanov@gmail.com',    '+7-900-123-45-67', 'Москва',          'Россия', '101000'),
    ('Мария',     'Петрова',  'm.petrova@yandex.ru',   '+7-911-234-56-78', 'Санкт-Петербург', 'Россия', '190000'),
    ('Дмитрий',   'Сидоров',  'd.sidorov@mail.ru',     '+7-922-345-67-89', 'Казань',          'Россия', '420000'),
    ('Анна',      'Козлова',  'a.kozlova@outlook.com', '+7-933-456-78-90', 'Новосибирск',     'Россия', '630000'),
    ('Сергей',    'Морозов',  's.morozov@gmail.com',   '+7-944-567-89-01', 'Екатеринбург',    'Россия', '620000');