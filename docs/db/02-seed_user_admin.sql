USE donation_hub;
INSERT INTO `users` (`id`, `username`, `email`, `password`, `created_at`)
VALUES
(1, 'admin', 'admin@donationhub.com', 'admin123', UNIX_TIMESTAMP()),
(2, 'donor', 'donor@donationhub.com', 'donatur123', UNIX_TIMESTAMP()),
(3, 'requester', 'requester@donationhub.com', 'requester123', UNIX_TIMESTAMP());

INSERT INTO `user_roles` (`user_id`, `role`)
VALUES
(1, 'admin'),
(2, 'donor'),
(3, 'requester');

INSERT  INTO `projects` (`id`, `name`, `description`, `target_amount`, `collection_amount`, `currency`, `status`, `requester_id`, `due_at`, `created_at`, `updated_at`)
VALUES
(1, 'Project Need Review', 'Description 1', 1000, 0, 'IDR', 'need_review', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(2, 'Project Approved', 'Description 2', 2000, 0, 'IDR', 'approved', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(3, 'Project Completed', 'Description 3', 3000, 0, 'IDR', 'completed', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP()),
(4, 'Project Rejected', 'Description 3', 3000, 0, 'IDR', 'rejected', 3, UNIX_TIMESTAMP(), UNIX_TIMESTAMP(), UNIX_TIMESTAMP());

INSERT INTO `project_images` (`id`, `project_id`, `url`)
VALUES
(1, 1, 'https://via.placeholder.com/150'),
(2, 2, 'https://via.placeholder.com/150'),
(3, 3, 'https://via.placeholder.com/150'),
(4, 4, 'https://via.placeholder.com/150');

-- insert donation only on approved project
INSERT INTO `donations` (`id`, `project_id`, `donor_id`, `message`, `amount`, `currency`, `created_at`)
VALUES
(1, 2, 2, 'Donation 1', 1000, 'IDR', UNIX_TIMESTAMP()),
(2, 2, 2, 'Donation 2', 1000, 'IDR', UNIX_TIMESTAMP()),
(3, 2, 2, 'Donation 3', 1000, 'IDR', UNIX_TIMESTAMP());