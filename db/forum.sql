-- SQLite
DROP TABLE IF EXISTS "comment_likes";
DROP TABLE IF EXISTS "post_likes";
DROP TABLE IF EXISTS "post_files";
DROP TABLE IF EXISTS "comments";
DROP TABLE IF EXISTS "post_categories";
DROP TABLE IF EXISTS "posts";
DROP TABLE IF EXISTS "categories";
DROP TABLE IF EXISTS "sessions";
DROP TABLE IF EXISTS "friends";
DROP TABLE IF EXISTS "users";
DROP TABLE IF EXISTS "chats";
DROP TABLE IF EXISTS "chat_members";
DROP TABLE IF EXISTS "messages";
DROP TABLE IF EXISTS "message_files";

CREATE TABLE "categories" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "name" TEXT NOT NULL,
  "color" TEXT NOT NULL,
  "icon" TEXT NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" INTEGER NOT NULL,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (created_by) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "users" (
  "id" INTEGER PRIMARY KEY,
  "uuid" TEXT NOT NULL UNIQUE,
  "type" TEXT NOT NULL CHECK ("type" IN ('admin', 'normal_user', 'test_user')) DEFAULT 'normal_user',
  "name" TEXT,
  "username" TEXT UNIQUE,
  "email" TEXT UNIQUE,
  "profile_photo" TEXT NULL,
  "password" TEXT,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "friends" (
  "id" INTEGER PRIMARY KEY,
  "first_user_id" INTEGER NOT NULL,
  "second_user_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('pending', 'accepted', 'delete')) DEFAULT 'pending',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" INTEGER NOT NULL,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (first_user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (second_user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (created_by) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "posts" (
  "id" INTEGER PRIMARY KEY,
  "uuid" TEXT NOT NULL UNIQUE,
  "title" TEXT NOT NULL,
  "description" TEXT NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "user_id" INTEGER NOT NULL,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "post_files" (
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "file_uploaded_name" TEXT NOT NULL,
  "file_real_name" TEXT NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" INTEGER NOT NULL,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (post_id) REFERENCES "posts" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "post_likes" (
  "id" INTEGER PRIMARY KEY,
  "type" TEXT NOT NULL CHECK ("type" IN ('like', 'dislike')),
  "post_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (post_id) REFERENCES "posts" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

CREATE TABLE "post_categories" (
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "category_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "created_by" INTEGER NOT NULL,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (created_by) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id"),
  FOREIGN KEY (post_id) REFERENCES "posts" ("id"),
  FOREIGN KEY (category_id) REFERENCES "categories" ("id")
);

CREATE TABLE "comments" (
  "id" INTEGER PRIMARY KEY,
  "post_id" INTEGER NOT NULL,
  "description" TEXT NOT NULL,
  "user_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id"),
  FOREIGN KEY (post_id) REFERENCES "posts" ("id")
);

CREATE TABLE "comment_likes" (
  "id" INTEGER PRIMARY KEY,
  "type" TEXT NOT NULL,
  "comment_id" INTEGER NOT NULL,
  "user_id" INTEGER NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  "updated_by" INTEGER,
  FOREIGN KEY (user_id) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id"),
  FOREIGN KEY (comment_id) REFERENCES "comments" ("id")
);

CREATE TABLE "sessions" (
  "id" INTEGER PRIMARY KEY,
  "session_token" TEXT NOT NULL UNIQUE,
  "user_id" INTEGER NOT NULL,
  "expires_at" DATETIME NOT NULL,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES "users" ("id")
);

CREATE TABLE "chats" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "uuid" TEXT NOT NULL UNIQUE,
  "type" TEXT NOT NULL CHECK ("type" IN ('private', 'group')) DEFAULT 'private',
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (created_by) REFERENCES "users" ("id")
);

CREATE TABLE "chat_members" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (chat_id) REFERENCES "chats" ("id") ON DELETE CASCADE,
  FOREIGN KEY (user_id) REFERENCES "users" ("id")
);

CREATE TABLE "messages" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "content" TEXT NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'disable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  FOREIGN KEY (created_by) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id"),
  FOREIGN KEY (chat_id) REFERENCES "chats" ("id") ON DELETE CASCADE
);

CREATE TABLE "message_files" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT,
  "file_uploaded_name" TEXT NOT NULL,
  "file_real_name" TEXT NOT NULL,
  "status" TEXT NOT NULL CHECK ("status" IN ('enable', 'delete')) DEFAULT 'enable',
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" DATETIME,
  FOREIGN KEY (chat_id) REFERENCES "chats" ("id") ON DELETE CASCADE,
  FOREIGN KEY (message_id) REFERENCES "messages" ("id") ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES "users" ("id"),
  FOREIGN KEY (updated_by) REFERENCES "users" ("id")
);

INSERT INTO users(uuid, type,name,username,password, email)
VALUES ('67921bdd-8458-800e-b9d4-065a43242cd3', 'admin', 'admin', 'admin', '$2a$10$DN.v/NkfQjmPaTTz15x0E.u8l2R9.HnB12DpDVMdRPeQZDfMwovSa', 'admin@admin');

INSERT INTO categories (name, color, icon, created_by)
VALUES ('Art', '#5340C8', '<i class="fa-solid fa-palette"></i>', 1), ('Science', '#7F59FE', '<i class="fa-solid fa-atom"></i>', 1),
 ('News', '#0BC4D5', '<i class="fa-solid fa-rss"></i>', 1), ('Technology', '#C34CB4', '<i class="fa-solid fa-microchip"></i>', 1), ('Lifestyle', '#C44343', '<i class="fa-brands fa-pagelines"></i>', 1);

-- INSERT INTO posts(uuid, title, description, user_id)
-- VALUES ('f9edb8d6-c739-4d6f-aaa4-9b298f2e1552', 'first post', 'this is first post of forum that is made by admin', 1);

INSERT INTO posts (uuid, title, description, user_id) VALUES 
('74ead645-b5aa-4622-b32f-ee8a844cce9c', 
'The Evolution of Abstract Art', 
'A look at how abstract art has transformed from traditional styles to modern digital expressions. 
  This post explores key movements like Cubism, Surrealism, and Expressionism, and how artists today use technology to push creative boundaries.',
1),

('dc30b21d-1bad-4dff-907c-9d31dbe54b9f',
'Quantum Computing: The Next Frontier', 
 'Quantum computing leverages the principles of quantum mechanics to perform calculations exponentially faster than classical computers. 
  This post delves into qubits, superposition, and real-world applications in cryptography and AI.', 
 1),

('11c33158-3b99-43d8-8731-38409c6cd87a', 
'Breaking: Major Space Discovery Announced', 
 'Astronomers have discovered a new Earth-like planet within the habitable zone of its star. 
  With the potential to host liquid water, this finding could be a game-changer in the search for extraterrestrial life.', 
 1),

('0c59ce51-dde6-4ff4-aac9-4c8fdc705610', 
'AI and the Future of Digital Art', 
 'Artificial Intelligence is transforming digital art creation. From AI-generated paintings to deep learning-based creative tools, 
  this post explores how algorithms are pushing the boundaries of what we consider original artwork.', 
 1),

('435ac036-77eb-41bc-8124-3d9521d62c19',
'The Role of CRISPR in Genetic Engineering', 
 'CRISPR technology has revolutionized genetic engineering, allowing scientists to edit DNA with unprecedented precision. 
  This breakthrough has significant implications for medicine, agriculture, and even species preservation.', 
 1),

('87166c00-58dc-4780-8df2-1b73ec99bcad', 
'Global Economic Trends in 2025', 
 'As the world recovers from economic shifts, new trends are shaping global markets. This report provides an in-depth analysis of inflation rates, 
  emerging industries, and how digital currencies are impacting financial landscapes.', 
 1);

INSERT INTO post_categories(post_id, category_id, created_by)
VALUES (1, 1, 1), (2, 2, 1), (2, 5, 1), (3, 3, 1), (4, 1, 1), (4, 3, 1), (5, 2, 1), (6, 3, 1);

INSERT INTO comments(post_id, description, user_id)
VALUES (1, 'this is first post comment that is made by admin', 1);