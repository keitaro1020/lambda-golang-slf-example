CREATE TABLE `cat` (
   `id` varchar(36) NOT NULL,
   `url` varchar(500) NOT NULL,
   `width` int(11) NOT NULL DEFAULT 0,
   `height` int(11) NOT NULL DEFAULT 0,
   `created_at` timestamp NULL DEFAULT NULL,
   PRIMARY KEY (`id`)
);