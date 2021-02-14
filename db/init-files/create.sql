CREATE TABLE nxlogd_db.`access_log` (
  `datetime` datetime NOT NULL,
  `ipaddress` varchar(15) NOT NULL,
  `port` int(10) unsigned NOT NULL,
  `asn` varchar(255) NOT NULL,
  `iso_code` varchar(3) NOT NULL,
  `country` varchar(100) NOT NULL,
  `city`  varchar(255) NOT NULL,
  `long`  FLOAT NOT NULL DEFAULT '0',
  `lat`  FLOAT NOT NULL DEFAULT '0',
  `url` text,
  `user_agent` text
);
