CREATE TABLE articles (
    id serial PRIMARY KEY,
    title VARCHAR (150) UNIQUE NOT NULL,
    body VARCHAR NOT NULL,
    author VARCHAR (150) NOT NULL
)