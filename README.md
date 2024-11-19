# edu-project


docker run --name edu-mysql -p 3307:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql:8.0.40 --init-connect="SET collation_connedtion=utf8mb4_0900 ai ci" --init-connect="SET NAMES utf8mb4" --skip-character-set-client-handshake

alter user root@localhost identified with mysql_native_password by '123456';
FLUSH PRIVILEGES；
select host,user,plugin,authentication string from rmysql.user;