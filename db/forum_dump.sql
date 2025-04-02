PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "categories" (
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
INSERT INTO categories VALUES(1,'Art','#5340C8','fa-solid fa-palette','enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO categories VALUES(2,'Science','#7F59FE','fa-solid fa-atom','enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO categories VALUES(3,'News','#0BC4D5','fa-solid fa-rss','enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO categories VALUES(4,'Technology','#C34CB4','fa-solid fa-microchip','enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO categories VALUES(5,'Lifestyle','#C44343','fa-brands fa-pagelines','enable','2025-02-09 12:49:34',1,NULL,NULL);
CREATE TABLE IF NOT EXISTS "users" (
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
INSERT INTO users VALUES(1,'67921bdd-8458-800e-b9d4-065a43242cd3','admin','admin','admin','admin@admin',NULL,'$2a$10$DN.v/NkfQjmPaTTz15x0E.u8l2R9.HnB12DpDVMdRPeQZDfMwovSa','enable','2025-02-09 12:49:34',NULL,NULL);
INSERT INTO users VALUES(2,'2dba1545-db5c-451b-857b-b76228f2f704','normal_user','parisa','parisa','parisa_rahimy@yahoo.com',NULL,'$2a$10$DN.v/NkfQjmPaTTz15x0E.u8l2R9.HnB12DpDVMdRPeQZDfMwovSa','enable','2025-02-09 12:54:52',NULL,NULL);
INSERT INTO users VALUES(3,'7e5389a6-6c90-464b-8204-388baa8661fb','normal_user','jedi1','jedi1','jedi1@jedi.com',NULL,'$2a$10$8ytQ3vWiyWMGBn0Tax7BFuAdhb.ootES5B1KYGqVIv6Bl7fQIJab2','enable','2025-03-27 17:12:25',NULL,NULL);
INSERT INTO users VALUES(4,'7bd1cef7-750d-4ebb-92e6-a0757cd87dc4','normal_user','jedi2','jedi2','jedi2@jedi.com',NULL,'$2a$10$AXcCGXvjJmWYYte.lJW0EuWxmItFrN63Z1KzQwGsclC//gdwZdyTK','enable','2025-03-27 17:13:37',NULL,NULL);
INSERT INTO users VALUES(5,'7329dc4e-54ed-4e3c-88a8-b63e7ba0c665','normal_user','parisaNew','parisaNew','parisa.rahimi@gmail.com',NULL,'$2a$10$SXUSRK0J90V7eM5J6mky3OadiRubFPqy6mIBn8nWQPcu7SRvTaK1K','enable','2025-03-28 09:08:23',NULL,NULL);
CREATE TABLE IF NOT EXISTS "friends" (
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
CREATE TABLE IF NOT EXISTS "posts" (
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
INSERT INTO posts VALUES(1,'74ead645-b5aa-4622-b32f-ee8a844cce9c','The Evolution of Abstract Art',replace('A look at how abstract art has transformed from traditional styles to modern digital expressions. \n  This post explores key movements like Cubism, Surrealism, and Expressionism, and how artists today use technology to push creative boundaries.','\n',char(10)),'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO posts VALUES(2,'dc30b21d-1bad-4dff-907c-9d31dbe54b9f','Quantum Computing: The Next Frontier',replace('Quantum computing leverages the principles of quantum mechanics to perform calculations exponentially faster than classical computers. \n  This post delves into qubits, superposition, and real-world applications in cryptography and AI.','\n',char(10)),'delete','2025-02-09 12:49:34',1,'2025-02-09 22:34:16',1);
INSERT INTO posts VALUES(3,'11c33158-3b99-43d8-8731-38409c6cd87a','Breaking: Major Space Discovery Announced',replace('Astronomers have discovered a new Earth-like planet within the habitable zone of its star. \n  With the potential to host liquid water, this finding could be a game-changer in the search for extraterrestrial life.','\n',char(10)),'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO posts VALUES(4,'0c59ce51-dde6-4ff4-aac9-4c8fdc705610','AI and the Future of Digital Art',replace('Artificial Intelligence is transforming digital art creation. From AI-generated paintings to deep learning-based creative tools, \n  this post explores how algorithms are pushing the boundaries of what we consider original artwork.','\n',char(10)),'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO posts VALUES(5,'435ac036-77eb-41bc-8124-3d9521d62c19','The Role of CRISPR in Genetic Engineering',replace('CRISPR technology has revolutionized genetic engineering, allowing scientists to edit DNA with unprecedented precision. \n  This breakthrough has significant implications for medicine, agriculture, and even species preservation.','\n',char(10)),'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO posts VALUES(6,'87166c00-58dc-4780-8df2-1b73ec99bcad','Global Economic Trends in 2025',replace('As the world recovers from economic shifts, new trends are shaping global markets. This report provides an in-depth analysis of inflation rates, \n  emerging industries, and how digital currencies are impacting financial landscapes.','\n',char(10)),'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO posts VALUES(7,'66df25b9-c0da-4518-9aef-c3300cb324d9','my first post','my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post my first post ','enable','2025-02-09 12:55:34',2,NULL,NULL);
INSERT INTO posts VALUES(8,'19807640-997a-4af1-a726-839e33f832f9','test','hello','enable','2025-02-09 22:49:05',2,'2025-02-09 22:50:04',2);
CREATE TABLE IF NOT EXISTS "post_files" (
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
INSERT INTO post_files VALUES(1,7,'75615e0f-ea07-430e-ac39-b74fa0faa36b.jpg','istockphoto-2165761854-1024x1024.jpg','enable','2025-02-09 12:59:05',2,NULL,NULL);
INSERT INTO post_files VALUES(2,2,'ead63826-cbb9-484b-8a2e-1d02f8a6d3a5.jpeg','download.jpeg','delete','2025-02-09 12:59:05',2,'2025-02-09 22:34:16',1);
CREATE TABLE IF NOT EXISTS "post_likes" (
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
INSERT INTO post_likes VALUES(1,'like',6,2,'enable','2025-02-09 12:55:11',NULL,NULL);
CREATE TABLE IF NOT EXISTS "post_categories" (
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
INSERT INTO post_categories VALUES(1,1,1,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(2,2,2,'delete','2025-02-09 12:49:34',1,'2025-02-09 22:34:16',1);
INSERT INTO post_categories VALUES(3,2,5,'delete','2025-02-09 12:49:34',1,'2025-02-09 22:34:16',1);
INSERT INTO post_categories VALUES(4,3,3,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(5,4,1,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(6,4,3,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(7,5,2,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(8,6,3,'enable','2025-02-09 12:49:34',1,NULL,NULL);
INSERT INTO post_categories VALUES(9,7,1,'enable','2025-02-09 12:55:34',2,NULL,NULL);
INSERT INTO post_categories VALUES(10,7,2,'enable','2025-02-09 12:55:34',2,NULL,NULL);
INSERT INTO post_categories VALUES(11,8,4,'delete','2025-02-09 22:49:05',2,'2025-02-09 22:49:31',2);
INSERT INTO post_categories VALUES(12,8,4,'delete','2025-02-09 22:49:31',2,'2025-02-09 22:50:04',2);
INSERT INTO post_categories VALUES(13,8,4,'enable','2025-02-09 22:50:04',2,NULL,NULL);
CREATE TABLE IF NOT EXISTS "comments" (
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
INSERT INTO comments VALUES(1,1,'this is first post comment that is made by admin',1,'enable','2025-02-09 12:49:34',NULL,NULL);
CREATE TABLE IF NOT EXISTS "comment_likes" (
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
CREATE TABLE IF NOT EXISTS "sessions" (
  "id" INTEGER PRIMARY KEY,
  "session_token" TEXT NOT NULL UNIQUE,
  "user_id" INTEGER NOT NULL,
  "expires_at" DATETIME NOT NULL,
  "created_at" DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES "users" ("id")
);
INSERT INTO sessions VALUES(1,'b8566437-3714-40c3-9842-2b97026e1f6a',2,'2025-02-09 16:42:14','2025-02-09 12:54:52');
INSERT INTO sessions VALUES(2,'61c8b387-b41b-4ea6-aac8-dd16c5cc52aa',1,'2025-02-09 21:21:07','2025-02-09 16:42:22');
INSERT INTO sessions VALUES(3,'166972d3-8f6d-4f3e-b924-e2baa5d68885',1,'2025-02-09 22:48:44','2025-02-09 21:21:14');
INSERT INTO sessions VALUES(4,'d7b43188-007f-42a7-8820-5bfb2afd33d3',2,'2025-02-10 12:48:54.213004+02:00','2025-02-09 22:48:54');
INSERT INTO sessions VALUES(5,'b9c159f7-d2a1-4e1e-827e-63863740b3fe',1,'2025-03-28 00:14:43.082605+02:00','2025-03-27 10:14:43');
INSERT INTO sessions VALUES(6,'d2af6fce-922d-46f3-b463-d590bd8df694',2,'2025-03-28 00:15:58.055621+02:00','2025-03-27 10:15:58');
INSERT INTO sessions VALUES(7,'fe2b9027-ad71-48ba-91f4-60a23c223df3',3,'2025-03-27 17:13:50','2025-03-27 17:12:25');
INSERT INTO sessions VALUES(8,'0a48a791-0a93-43e8-b073-f1b4cbad4faf',4,'2025-03-28 07:13:37.661815+02:00','2025-03-27 17:13:37');
INSERT INTO sessions VALUES(9,'e6228d17-7d63-49d2-959e-a92051ed3817',3,'2025-03-28 07:14:09.978385+02:00','2025-03-27 17:14:09');
INSERT INTO sessions VALUES(10,'d28573a5-ef7c-4212-99bd-2fd8da22a102',5,'2025-03-28 12:08:23.419887+02:00','2025-03-28 09:08:23');
INSERT INTO sessions VALUES(11,'3d9cb27d-c65d-42cf-b798-ada0d4f6b7df',5,'2025-03-28 18:46:44.086427+02:00','2025-03-28 15:46:44');
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('categories',6);
COMMIT;
