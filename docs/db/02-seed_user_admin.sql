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
