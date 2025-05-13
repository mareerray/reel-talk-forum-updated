DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "posts";
DROP TABLE IF EXISTS "comments";
DROP TABLE IF EXISTS "chats";
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "message_read_receipts";


CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    nickname TEXT NOT NULL UNIQUE,
    gender TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    session_token TEXT UNIQUE, 
    session_expiry DATETIME, 
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_activity DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    session_expiry DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    emoji TEXT NOT NULL UNIQUE
);

CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    categories TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    user_name TEXT NOT NULL,
    post_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);

CREATE TABLE chats (
    "id" INTEGER PRIMARY KEY,
    "uuid" TEXT NOT NULL UNIQUE,
    "user_id_1" INTEGER NOT NULL,
    "user_id_2" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" DATETIME,
    "updated_by" INTEGER,
    FOREIGN KEY (user_id_1) REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY (user_id_2) REFERENCES "users" ("id") ON DELETE CASCADE,
    FOREIGN KEY (updated_by) REFERENCES "users" ("id") ON DELETE CASCADE,
    UNIQUE (user_id_1, user_id_2) 
);

CREATE TABLE messages (
    "id" INTEGER PRIMARY KEY,
    "chat_id" INTEGER NOT NULL,
    "user_id_from" INTEGER NOT NULL,
    "content" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" DATETIME,
    FOREIGN KEY (chat_id) REFERENCES "chats" ("id") ON DELETE CASCADE,
    FOREIGN KEY (user_id_from) REFERENCES "users" ("id") ON DELETE CASCADE
);

CREATE TABLE message_read_receipts (
    chat_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    read_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (chat_id, user_id),
    FOREIGN KEY (chat_id) REFERENCES chats(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);


-------------------- Insert sample data into the database------------------------------------------------
-- Insert categories
INSERT INTO categories (name, emoji) VALUES
('Action', '💥'), ('Adventure', '🌄'), ('Animation', '🧚'), 
('Biography', '📚'), ('Comedy', '😂'), ('Crime', '🕵️'), ('Documentary', '🎥'), 
('Drama', '🎭'), ('Fantasy', '🧙'), ('Horror', '👻'), ('Mystery', '🔍'), 
('Romance', '❤️'), ('Sci-Fi', '🚀'), ('Thriller', '😱'), ('Western', '🤠');

-- Insert users
INSERT INTO users (uuid, nickname, email, password_hash, first_name, last_name, gender, age, created_at) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'admin', 'admin@admin.com', '$2a$10$ryPUUMn0CPeuNh.NpQZOwuyoymt1sdzXrePhSeYArwv9puWlg1mF2', 'Admin', 'User', 'Male', 30, CURRENT_TIMESTAMP),
('6ba7b810-9dad-11d1-80b4-00c04fd430c0', 'Mama', 'mama@yahoo.com', '$2a$10$bfVNqrSBscGyfsGMSyEvaOCRbBbC54I2Lht5XuaBLiZKcdgoIRJQO', 'Mama', 'User', 'Female', 35, CURRENT_TIMESTAMP),
('6ba7b811-9dad-11d1-80b4-00c04fd430c1', 'batman', 'batman@batman.com', '$2a$10$1ZAK4MxQuwCJZGqhpBBzPOMoDDeGob..uwEIIO9YsHpqx8qXPNH8u', 'Bruce', 'Wayne', 'Male', 40, CURRENT_TIMESTAMP);


-- Insert posts
INSERT INTO posts (title, content, user_id, categories, created_at) VALUES
('The Thrilling Ride of "Quantum Horizon"', 'Quantum Horizon blends cutting-edge effects with a gripping narrative. The zero-gravity fights are breathtaking.', 2, 'Sci-Fi 🚀,Action 💥,Thriller 😱', CURRENT_TIMESTAMP),
('Laughing Through Time: A Hilarious Adventure', 'A refreshing take on time-travel comedy. Clever writing and impeccable timing had me in stitches.', 1, 'Comedy 😂,Adventure 🌄,Sci-Fi 🚀', CURRENT_TIMESTAMP),
('Whispers of the West: A Haunting Frontier Tale', 'This unconventional Western infuses horror into a frontier setting. Eerie atmosphere keeps viewers on edge.', 3, 'Western 🤠,Horror 👻,Mystery 🔍', CURRENT_TIMESTAMP),
('Brushstrokes of Genius: A Compelling Artist''s Biography', 'A meticulously crafted documentary about painter Isabella Rossi. Balances interviews with stunning visuals of her work.', 2, 'Documentary 🎥,Biography 📚', CURRENT_TIMESTAMP);

-- Insert comments for existing posts
INSERT INTO comments (post_id, content, user_id, user_name, created_at) VALUES
(1, 'This is an amazing post about Quantum Horizon!', 2, 'Mama', CURRENT_TIMESTAMP),
(2, 'Laughing Through Time is a masterpiece!', 1, 'admin', CURRENT_TIMESTAMP),
(3, 'Whispers of the West is so hauntingly beautiful.', 3, 'batman', CURRENT_TIMESTAMP),
(4, 'Brushstrokes of Genius is truly inspiring.', 2, 'Mama', CURRENT_TIMESTAMP);

-- Insert replies to comments (Note: your schema doesn't support nested comments, so these are additional comments on the same posts)
INSERT INTO comments (post_id, content, user_id, user_name, created_at) VALUES
(1, 'I completely agree with you, Mama!', 3, 'batman', CURRENT_TIMESTAMP),
(2, 'Admin, you have a great taste in comedy!', 2, 'Mama', CURRENT_TIMESTAMP),
(3, 'Batman, I felt the same way about the eerie atmosphere.', 1, 'admin', CURRENT_TIMESTAMP),
(4, 'Mama, the visuals were breathtaking indeed.', 3, 'batman', CURRENT_TIMESTAMP);


-- Additional comments with more discussion
INSERT INTO comments (post_id, content, user_id, user_name, created_at) VALUES
(1, 'The special effects in Quantum Horizon were mind-blowing!', 1, 'admin', CURRENT_TIMESTAMP),
(1, 'I thought the plot was a bit confusing in the middle.', 3, 'batman', CURRENT_TIMESTAMP),
(2, 'The time travel paradoxes were actually well-explained.', 3, 'batman', CURRENT_TIMESTAMP),
(2, 'I laughed so hard at the dinosaur scene!', 2, 'Mama', CURRENT_TIMESTAMP),
(3, 'The soundtrack really enhanced the creepy atmosphere.', 2, 'Mama', CURRENT_TIMESTAMP),
(3, 'I could not sleep after watching this!', 1, 'admin', CURRENT_TIMESTAMP),
(4, 'I did not know much about Isabella Rossi before this.', 3, 'batman', CURRENT_TIMESTAMP),
(4, 'Her use of color is revolutionary.', 1, 'admin', CURRENT_TIMESTAMP);


