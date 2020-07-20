ALTER TABLE feed ADD alias varchar(25) default '';
UPDATE feed SET (alias) = 'int' WHERE id in (1,2);
UPDATE feed SET (alias) = 'stage' WHERE id in (5,6);
UPDATE feed SET (alias) = 'prod' WHERE id in (7,8);

