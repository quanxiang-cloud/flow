CREATE TABLE `flow_form_field` (
  `id` varchar(40) NOT NULL DEFAULT '' COMMENT '主键',
  `flow_id` varchar(40) NOT NULL DEFAULT '' COMMENT '流程id',
  `form_id` varchar(50) NOT NULL DEFAULT '' COMMENT '表单id',
  `field_name` varchar(50) NOT NULL DEFAULT '' COMMENT '字段名',
  `field_value_path` varchar(100) NOT NULL DEFAULT '' COMMENT '字段值path',
  `creator_id` varchar(40) NOT NULL DEFAULT '',
  `create_time` varchar(40) DEFAULT NULL,
  `modifier_id` varchar(40) NOT NULL DEFAULT '',
  `modify_time` varchar(40) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `flow_id` (`flow_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='工作流表单字段';

ALTER TABLE flow_operation_record  ADD  rel_node_def_key varchar(255) NOT NULL DEFAULT '' COMMENT '关联的node节点';

update flow_variables set field_type ="string" where field_type ='TEXT';
update flow_variables set field_type ="boolean" where field_type ='BOOLEAN';
update flow_variables set field_type ="datetime" where field_type ='DATE';
update flow_variables set field_type ="number" where field_type ='NUMBER' ;

update flow_instance_variables set field_type ="string" where field_type ='TEXT';
update flow_instance_variables set field_type ="boolean" where field_type ='BOOLEAN';
update flow_instance_variables set field_type ="datetime" where field_type ='DATE';
update flow_instance_variables set field_type ="number" where field_type ='NUMBER' ;