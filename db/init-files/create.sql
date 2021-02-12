CREATE TABLE nxlogd_db.`access_log` (
  `datetime` datetime NOT NULL,
  `ipaddress` varchar(15) NOT NULL,
  `port` int(10) unsigned NOT NULL,
  `country` varchar(100) NOT NULL,
  `method` varchar(10) DEFAULT NULL,
  `url` text,
  `status_code` int(10) unsigned NOT NULL,
  `sent_bytes` int(10) unsigned NOT NULL DEFAULT '0',
  `referrer` text,
  `user_agent` text
);
