USE tinyid;

CREATE TABLE `tiny_id_info` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `biz_type` varchar(63) NOT NULL,
  `max_id` bigint(20) NOT NULL DEFAULT 0,
  `step` int(11) DEFAULT 0,
  `delta` int(11) NOT NULL DEFAULT 1,
  `remainder` int(11) NOT NULL DEFAULT 0,
  `create_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `version` bigint(20) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq_biz_type` (`biz_type`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO tiny_id_info SET biz_type = "test", max_id = 0, step = 1000, delta = 2, remainder = 0, create_time = now();
