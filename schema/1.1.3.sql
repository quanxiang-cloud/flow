CREATE TABLE `flow_process_relation`
(
    `id`          varchar(40) NOT NULL DEFAULT '' COMMENT '流程id',
    `bpmn_text`   text        NOT NULL COMMENT '流程xml文件内容',
    `flow_id`     varchar(40) NOT NULL DEFAULT '' COMMENT 'flowID',
    `process_id`  varchar(40) NOT NULL DEFAULT '' COMMENT 'process中流程的id',
    `creator_id`  varchar(40) NOT NULL DEFAULT '' COMMENT '创建人',
    `create_time` varchar(40)          DEFAULT NULL COMMENT '创建时间',
    `modifier_id` varchar(40) NOT NULL DEFAULT '' COMMENT '更新人',
    `modify_time` varchar(40)          DEFAULT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='流程实例关系表';
