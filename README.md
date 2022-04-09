# mnsql是memory nosql
仅仅使用mnsql.c与mnsql.h即可。  

# 存在如下函数
写入对象，但是如果存在相同的则覆盖写入  
func Set(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
但是如果存在相同的则覆盖写入  
func SetEx(key string, value interface{}) int  

写入对象，但是如果存在相同的key则不写入  
func SetNx(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
如果存在相同的key则不写入  
func SetNex(key string, value interface{}, ttl int) int  

读取对象，如果对象不存在，第2个返回值会返回-2  
func Get(key string) (interface{}, int)  

删除对象  
func Del(key string) int  

对象自动加1，如果对象不存在，会创建一个C.int型的对象，赋值为1  
func Incr(key string) int  

对象自动减1，如果对象不存在，会创建一个C.int型的对象，赋值为-1  
func Decr(key string) int  

# 使用方法  
在需要使用的对象中  
import "github.com/worldflyingct/mnsql"  
然后  
go get github.com/worldflyingct/mnsql  
或是直接使用go mod命令  
go mod tidy  

# 关于datatype
因为存储在内存中仅仅是通过char*的通用格式进行保存的，所以并没有存储数据类型的方法，特意保留一个字段用来记录数据类型，如果使用c进行调用，不是一定要使用，但是由于incr与decr在key不存在时会创建int类型的数据，因此定义int类型为0，其他随意。  
