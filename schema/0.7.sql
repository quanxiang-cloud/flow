test

ALTER TABLE flow  ADD  app_status varchar(20) NOT NULL DEFAULT 'ACTIVE' COMMENT 'app状态' after app_id;
ALTER TABLE flow_instance  ADD  app_status varchar(20) NOT NULL DEFAULT 'ACTIVE' COMMENT 'app状态' after app_id;
