USE donation_hub;

INSERT INTO `users` (`id`, `username`, `email`, `password`, `created_at`)
VALUES
(1, 'admin', 'admin@donationhub.com', 'admin123', UNIX_TIMESTAMP()),
(2, 'donor', 'donor@donationhub.com', 'donor123', UNIX_TIMESTAMP()),
(3, 'requester', 'requester@donationhub.com', 'requester123', UNIX_TIMESTAMP()),
(4, 'user1', 'user1@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(5, 'user2', 'user2@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(6, 'user3', 'user3@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(7, 'user4', 'user4@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(8, 'user5', 'user5@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(9, 'user6', 'user6@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(10, 'user7', 'user7@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(11, 'user8', 'user8@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(12, 'user9', 'user9@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(13, 'user10', 'user10@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(14, 'user11', 'user11@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(15, 'user12', 'user12@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(16, 'user13', 'user13@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(17, 'user14', 'user14@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(18, 'user15', 'user15@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(19, 'user16', 'user16@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(20, 'user17', 'user17@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(21, 'user18', 'user18@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(22, 'user19', 'user19@donationhub.com', 'user123', UNIX_TIMESTAMP()),
(23, 'user20', 'user20@donationhub.com', 'user123', UNIX_TIMESTAMP());

INSERT INTO `user_roles` (`user_id`, `role`)
VALUES
(1, 'admin'),
(2, 'donor'),
(3, 'requester'),
(4, 'donor'),
(5, 'donor'),
(6, 'donor'),
(7, 'donor'),
(8, 'donor'),
(9, 'donor'),
(10, 'donor'),
(11, 'donor'),
(12, 'donor'),
(13, 'donor'),
(14, 'donor'),
(15, 'donor'),
(16, 'donor'),
(17, 'donor'),
(18, 'donor'),
(19, 'donor'),
(20, 'donor'),
(21, 'donor'),
(22, 'donor'),
(23, 'donor');

INSERT  INTO `projects` (`id`, `name`, `description`, `target_amount`, `collection_amount`, `currency`, `status`, `requester_id`, `due_at`, `created_at`, `updated_at`)
VALUES
(1, 'Project Need Review', 'Description 1', 1000, 0, 'IDR', 'need_review', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 'Project Approved', 'Description 2', 2000, 0, 'IDR', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 'Project Completed', 'Description 3', 3000, 0, 'IDR', 'completed', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 'Project Rejected', 'Description 3', 3000, 0, 'IDR', 'rejected', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

INSERT INTO `project_images` (`id`, `project_id`, `url`)
VALUES
(1, 1, 'https://dummyimage.com/600x400/000/fff&text=Project+1'),
(2, 2, 'https://dummyimage.com/600x400/000/fff&text=Project+2'),
(3, 3, 'https://dummyimage.com/600x400/000/fff&text=Project+3'),
(4, 4, 'https://dummyimage.com/600x400/000/fff&text=Project+4');

-- Generate 20 sample projects about animals feeding
INSERT INTO `projects` (`id`, `name`, `description`, `target_amount`, `collection_amount`, `currency`, `status`, `requester_id`, `due_at`, `created_at`, `updated_at`)
VALUES
(5, 'Feeding Lions', 'Feeding lions in the zoo', 500, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(6, 'Feeding Elephants', 'Feeding elephants in the sanctuary', 800, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(7, 'Feeding Tigers', 'Feeding tigers in the wildlife reserve', 700, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(8, 'Feeding Giraffes', 'Feeding giraffes in the safari park', 600, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(9, 'Feeding Penguins', 'Feeding penguins in the aquarium', 400, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(10, 'Feeding Dolphins', 'Feeding dolphins in the marine park', 900, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(11, 'Feeding Koalas', 'Feeding koalas in the wildlife sanctuary', 300, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(12, 'Feeding Gorillas', 'Feeding gorillas in the zoo', 400, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(13, 'Feeding Zebras', 'Feeding zebras in the safari park', 600, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(14, 'Feeding Kangaroos', 'Feeding kangaroos in the wildlife reserve', 500, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(15, 'Feeding Bears', 'Feeding bears in the wildlife sanctuary', 700, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(16, 'Feeding Monkeys', 'Feeding monkeys in the zoo', 400, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(17, 'Feeding Macan', 'Feeding macan in the zoo', 500, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(18, 'Feeding Lumba Lumba', 'Feeding Lumba Lumba in the zoo', 800, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(19, 'Feeding Koala', 'Feeding Koala in the zoo', 300, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(20, 'Feeding Munyuk', 'Feeding Munyuk in the company', 400, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(21, 'Feeding zebra', 'Feeding zebra in the zoo', 600, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(22, 'Feeding kang guru', 'Feeding kang guru in the zoo', 500, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(23, 'Feeding ber uang', 'Feeding ber uang in the zoo', 700, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(24, 'Feeding uu aak', 'Feeding uu aak in the company', 400, 0, 'USD', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

--  generate sample project image for animals
INSERT INTO `project_images` (`id`, `project_id`, `url`)
VALUES
(5, 5, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Lions'),
(6, 6, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Elephants'),
(7, 7, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Tigers'),
(8, 8, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Giraffes'),
(9, 9, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Penguins'),
(10, 10, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Dolphins'),
(11, 11, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Koalas'),
(12, 12, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Gorillas'),
(13, 13, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Zebras'),
(14, 14, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Kangaroos'),
(15, 15, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Bears'),
(16, 16, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Monkeys'),
(17, 17, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Penguins'),
(18, 18, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Dolphins'),
(19, 19, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Koalas'),
(20, 20, 'https://dummyimage.com/600x400/000/fff&text=Feeding+Gorillas'),
(21, 21, 'https://dummyimage.com/600x400/000/fff&text=Feeding+zebra'),
(22, 22, 'https://dummyimage.com/600x400/000/fff&text=Feeding+kangaroo'),
(23, 23, 'https://dummyimage.com/600x400/000/fff&text=Feeding+bear'),
(24, 24, 'https://dummyimage.com/600x400/000/fff&text=Feeding+monkey');