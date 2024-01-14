# redis

## 常见面试题

### Redis数据类型

* string类型的应用场景：缓存对象、常规计数、分布式锁、共享session信息
* list类型的应用场景：消息队列
* hash类型：缓存对象、购物车等
* set类型：聚合计算（并集、交集、差集），比如点赞、共同关注、抽奖活动
* zset类型：排序场景，比如排行榜、电话和姓名排序

1. string的内部实现：简单动态字符串(sds)。sds保存二进制数据，不以空字符判断字符串是否结束，而是使用len属性的值；sds获取字符串长度的时间复杂度为O(1)；sds是api安全的，拼接字符串不会导致缓冲区溢出。
2. list类型的内部实现：双向链表和压缩列表。元素个数小于512个，每个元素的值都小于64字节，采用压缩列表，否则采用双向链表。
3. hash的内部实现：元素个数小于512个，每个元素的值都小于64字节，采用压缩列表，否则采用哈希表。
4. set的内部实现：元素都是整数且元素个数小于 512采用整数集合实现，否则采用哈希表
5. zset的内部实现：元素个数小于128个，每个元素值小于64字节，采用压缩列表作为zset类型的底层数据结构，否则采用跳表实现

### Redis的线程模型
