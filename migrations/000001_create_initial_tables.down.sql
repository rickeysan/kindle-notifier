-- Drop indexes
DROP INDEX IF EXISTS idx_categories_name;
DROP INDEX IF EXISTS idx_books_added_at;
DROP INDEX IF EXISTS idx_books_asin;
 
-- Drop tables
DROP TABLE IF EXISTS book_categories;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS books; 