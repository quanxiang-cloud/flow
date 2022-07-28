alter table flow add cron varchar(20) null;


alter table flow_variables
    modify code varchar(200) default '' not null comment '变量标识';


alter table flow_instance_variables
    modify code varchar(200) default '' not null comment '变量标识';
