-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    asin VARCHAR(255) NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    image_url TEXT NOT NULL,
    product_url TEXT NOT NULL,
    is_ku BOOLEAN NOT NULL DEFAULT false,
    added_at TIMESTAMP NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create book_categories junction table
CREATE TABLE IF NOT EXISTS book_categories (
    book_id INTEGER REFERENCES books(id) ON DELETE CASCADE,
    category_id INTEGER REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (book_id, category_id)
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_books_asin ON books(asin);
CREATE INDEX IF NOT EXISTS idx_books_added_at ON books(added_at);
CREATE INDEX IF NOT EXISTS idx_categories_name ON categories(name); 