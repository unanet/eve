drop table if exists database_instance;
drop table if exists database_server;
drop table if exists database_type;

alter table namespace drop column if exists metadata;
alter table artifact drop column if exists metadata;