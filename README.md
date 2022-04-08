# mnsql是memory nosql
仅仅使用mnsql.c与mnsql.h即可。

# 关于datatype
因为存储在内存中仅仅是通过char*的通用格式进行保存的，所以并没有存储数据类型的方法，特意保留一个字段用来记录数据类型，如果使用c进行调用，不是一定要使用，但是由于incr与decr在key不存在时会创建int类型的数据，因此定义int类型为0，其他随意。
