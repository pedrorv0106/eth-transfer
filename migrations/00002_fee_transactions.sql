-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE `fee_transactions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `fee_txid` varchar(255) DEFAULT NULL,
  `txid` varchar(255) DEFAULT NULL,
  `amount` varchar(255) DEFAULT NULL,
  `state` varchar(255) DEFAULT NULL,
  `priv_key` varchar(255) DEFAULT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `txid` (`txid`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE fee_transactions;