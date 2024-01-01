# slowsql

* show slowsql logs

* show slowsql stats

* kill query

* reset slowsql stats

```
mysql -u test -h 10.177.54.121 -P 3388 -p123456

mysql -u admin -h 10.177.54.121 -P 13388 -paNekZX9CWyve@RzQkY

select sleep(1000) from test;

// 查看慢查询日志
show slowsql logs limit 5;

// id为慢查询的编号对应SqlQueryId字段
kill query id;

// 查看慢查询stats
show slowsql stats;

reset slowsql stats;
```

# show transaction logs

* show transaction logs

* kill query qid in tid

```
select sleep(1000) from test where id = 1;

// 查看慢事务日志
show transaction logs

// qid对应SqlQueryId字段、tid对应的TransId字段
kill query qid in tid
```

# show client stats

```
show slowsql logs limit 5;

select sleep(1000) from test where id = 1;

show client stats;

// 通过id编号中断指定的客户端连接, id取值对应show client stats对应的ClientId字段。
kill connection id;

// 通过指定ip断开所有客户端连接或中断某个数据库的客户端连接
kill client `ip`;

show endpoint pool;
```

# show info

```
show sql info

show sql stats

show sql errors;

show critical logs;

show proxy info;
```