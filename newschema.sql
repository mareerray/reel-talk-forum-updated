-- USERS
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid TEXT UNIQUE NOT NULL,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    nickname TEXT NOT NULL UNIQUE,
    gender TEXT CHECK(gender IN ('Male', 'Female', 'Other')) NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    session_token TEXT UNIQUE,
    session_expiry TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_activity TIMESTAMPTZ DEFAULT NOW()
);

-- SESSIONS
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL,
    session_token TEXT NOT NULL UNIQUE,
    session_expiry TIMESTAMPTZ NOT NULL
);

-- CATEGORIES
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    emoji TEXT NOT NULL UNIQUE
);

-- POSTS
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    categories TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- COMMENTS
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_name TEXT NOT NULL,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- CHATS
CREATE TABLE chats (
    id SERIAL PRIMARY KEY,
    uuid TEXT NOT NULL UNIQUE,
    user_id_1 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    user_id_2 INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    updated_by INTEGER REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id_1, user_id_2)
);

-- MESSAGES
CREATE TABLE messages (
    id SERIAL PRIMARY KEY,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id_from INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- MESSAGE READ RECEIPTS
CREATE TABLE message_read_receipts (
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    read_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (chat_id, user_id)
);

-- SEED DATA
INSERT INTO categories (name, emoji) VALUES
('Action', '💥'), ('Adventure', '🌄'), ('Animation', '🧚'),
('Biography', '📚'), ('Comedy', '😂'), ('Crime', '🕵️'), ('Documentary', '🎥'),
('Drama', '🎭'), ('Fantasy', '🧙'), ('Horror', '👻'), ('Mystery', '🔍'),
('Romance', '❤️'), ('Sci-Fi', '🚀'), ('Thriller', '😱'), ('Western', '🤠');

INSERT INTO users (uuid, nickname, email, password_hash, first_name, last_name, gender, age) VALUES
('550e8400-e29b-41d4-a716-446655440000', 'admin', 'admin@admin.com', '$2a$10$ryPUUMn0CPeuNh.NpQZOwuyoymt1sdzXrePhSeYArwv9puWlg1mF2', 'Admin', 'User', 'Male', 30),
('6ba7b810-9dad-11d1-80b4-00c04fd430c0', 'Mama', 'mama@yahoo.com', '$2a$10$bfVNqrSBscGyfsGMSyEvaOCRbBbC54I2Lht5XuaBLiZKcdgoIRJQO', 'Mama', 'User', 'Female', 35),
('6ba7b811-9dad-11d1-80b4-00c04fd430c1', 'batman', 'batman@batman.com', '$2a$10$1ZAK4MxQuwCJZGqhpBBzPOMoDDeGob..uwEIIO9YsHpqx8qXPNH8u', 'Bruce', 'Wayne', 'Male', 40);

INSERT INTO posts (title, content, user_id, categories) VALUES
('The Thrilling Ride of "Quantum Horizon"', 'Quantum Horizon blends cutting-edge effects with a gripping narrative. The zero-gravity fights are breathtaking.', 2, 'Sci-Fi 🚀,Action 💥,Thriller 😱'),
('Laughing Through Time: A Hilarious Adventure', 'A refreshing take on time-travel comedy. Clever writing and impeccable timing had me in stitches.', 1, 'Comedy 😂,Adventure 🌄,Sci-Fi 🚀'),
('Whispers of the West: A Haunting Frontier Tale', 'This unconventional Western infuses horror into a frontier setting. Eerie atmosphere keeps viewers on edge.', 3, 'Western 🤠,Horror 👻,Mystery 🔍'),
('Brushstrokes of Genius: A Compelling Artist''s Biography', 'A meticulously crafted documentary about painter Isabella Rossi. Balances interviews with stunning visuals of her work.', 2, 'Documentary 🎥,Biography 📚');

INSERT INTO comments (post_id, content, user_id, user_name) VALUES
(1, 'This is an amazing post about Quantum Horizon!', 2, 'Mama'),
(2, 'Laughing Through Time is a masterpiece!', 1, 'admin'),
(3, 'Whispers of the West is so hauntingly beautiful.', 3, 'batman'),
(4, 'Brushstrokes of Genius is truly inspiring.', 2, 'Mama'),
(1, 'I completely agree with you, Mama!', 3, 'batman'),
(2, 'Admin, you have a great taste in comedy!', 2, 'Mama'),
(3, 'Batman, I felt the same way about the eerie atmosphere.', 1, 'admin'),
(4, 'Mama, the visuals were breathtaking indeed.', 3, 'batman'),
(1, 'The special effects in Quantum Horizon were mind-blowing!', 1, 'admin'),
(1, 'I thought the plot was a bit confusing in the middle.', 3, 'batman'),
(2, 'The time travel paradoxes were actually well-explained.', 3, 'batman'),
(2, 'I laughed so hard at the dinosaur scene!', 2, 'Mama'),
(3, 'The soundtrack really enhanced the creepy atmosphere.', 2, 'Mama'),
(3, 'I could not sleep after watching this!', 1, 'admin'),
(4, 'I did not know much about Isabella Rossi before this.', 3, 'batman'),
(4, 'Her use of color is revolutionary.', 1, 'admin');