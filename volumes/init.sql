CREATE TABLE IF NOT EXISTS movies (
    id INT PRIMARY KEY,
    slug TEXT,
    link TEXT,
    title TEXT,
    rating REAL,
    popularity INT,
    poster TEXT,
    backdrop TEXT
);
