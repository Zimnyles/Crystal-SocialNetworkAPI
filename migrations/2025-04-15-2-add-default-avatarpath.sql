ALTER TABLE users 
ALTER COLUMN avatarpath SET DEFAULT '/static/images/defaultuseravatars/defaultuseravatar.png';

UPDATE users 
SET avatarpath = '/static/images/defaultuseravatars/defaultuseravatar.png'
WHERE avatarpath IS NULL OR avatarpath = '';
