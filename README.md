# mnsql是memory nosql
仅仅使用mnsql.c与mnsql.h即可。  

# 数据类型
提供String、List与Hash类型

# 存在如下函数
写入对象，但是如果存在相同的则覆盖写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func Set(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
但是如果存在相同的则覆盖写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func SetEx(key string, value interface{}, ttl int64) int  

写入对象，但是如果存在相同的key则不写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；-4代表已经存在相同key的对象，写入失败  
func SetNx(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
如果存在相同的key则不写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；-4代表已经存在相同key的对象，写入失败  
func SetNex(key string, value interface{}, ttl int) int  

读取对象  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；-3代表未知的数据类型；  
func Get(key string) (interface{}, int)  

删除对象  
返回定义：0代表成功；-1代表key长度为0；  
func Del(key string) int  

对象自动加1，如果对象不存在，会创建一个int64型的对象，赋值为1  
返回定义：0代表成功；-1代表key长度为0；  
func Incr(key string) int  

对象自动减1，如果对象不存在，会创建一个int64型的对象，赋值为-1  
返回定义：0代表成功；-1代表key长度为0；  
func Decr(key string) int  

对象自动加num，如果对象不存在，会创建一个int64型的对象，赋值为num  
返回定义：0代表成功；-1代表key长度为0；  
func IncrBy(key string, num int64) int  

对象自动减num，如果对象不存在，会创建一个int64型的对象，赋值为-num  
返回定义：0代表成功；-1代表key长度为0；  
func DecrBy(key string, num int64) int  

设置对象的生命期，ttl如果为-1则是永久存在  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；  
func Expire(key string, ttl int64) int  

从左边推入list，如果不存在则创建  
返回定义：0代表成功；-1代表key长度为0；-3代表对象类型错误；  
func LPush(key string, value interface{}) int  

从左边取出list  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；-3代表对象类型错误；  
func LPop(key string) (interface{}, int)  

从右边推入list，如果不存在则创建  
返回定义：0代表成功；-1代表key长度为0；-3代表对象类型错误；  
func RPush(key string, value interface{}) int  

从右边取出list  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；-3代表对象类型错误；  
func RPop(key string) (interface{}, int)  

写入hash对象，但是如果存在相同key但不是hash类的则失败  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func HSet(key string, key2 string, value interface{}) int  

写入有生命期的hash对象，但是如果存在相同key但不是hash类的则失败，生命期消失后数据自动消失  
但是如果存在相同但不是hash类的则失败，如果当前hash存在相同key则覆盖  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func HSetEx(key string, key2 string, value interface{}) int  

写入hash对象，但是如果存在相同key但不是hash类的则失败  
如果当前hash存在相同key则放弃  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；-4代表已经存在相同key的对象，写入失败  
func HSetNx(key string, key2 string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
但是如果存在相同但不是hash类的则失败，如果当前hash存在相同key则放弃，如果存在相同的key则不写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；-4代表已经存在相同key的对象，写入失败  
func HSetNex(key string, key2 string, value interface{}, ttl int) int  

读取hash对象  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；-3代表未知的数据类型；  
func HGet(key string, key2 string) (interface{}, int)  

删除hash对象  
返回定义：0代表成功；-1代表key长度为0；  
func HDel(key string, key2 string) int  

hash对象自动加1，如果对象不存在，会创建一个int64型的对象，赋值为1  
返回定义：0代表成功；-1代表key长度为0；  
func HIncr(key string, key2 string) int  

hash对象自动减1，如果对象不存在，会创建一个int64型的对象，赋值为-1  
返回定义：0代表成功；-1代表key长度为0；  
func HDecr(key string, key2 string) int  

hash对象自动加num，如果对象不存在，会创建一个int64型的对象，赋值为num  
返回定义：0代表成功；-1代表key长度为0；  
func HIncrBy(key string, key2 string, num int64) int  

hash对象自动减num，如果对象不存在，会创建一个int64型的对象，赋值为-num  
返回定义：0代表成功；-1代表key长度为0；  
func HDecrBy(key string, key2 string, num int64) int  

设置hash对象的生命期，ttl如果为-1则是永久存在  
返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；  
func HExpire(key string, key2 string, ttl int64) int  

获取顶层所有的key  
返回定义：0代表成功；-5代表内存不足；  
func Keys() (string, int)  
多个key之间以\r\n作为包间隔，key内部采用key名称，ttl与数据类型作为返回的字符串  

获取hash对象所有的key  
返回定义：0代表成功；-5代表内存不足；  
func HKeys(key string) (string, int)  
多个key之间以\r\n作为包间隔，key内部采用key名称，ttl与数据类型作为返回的字符串  

清空内存数据库  
func FlushDB()  
无参数，无返回  

启动调试服务，port为端口号  
func StartCmdLineServer(port uint16) (net.Listener, error)  
调试服务的使用方法是通过tcp去连接，然后下各种查询命令。已有命令如下  
keys 显示全部的key  
hkeys 显示某个hash对象全部的key  
get 查询string对象  
hget 查询hash对象  
lrange 查询全部list  
exit 退出  
closeserver 关闭调试服务  

# 使用方法  
在需要使用的对象中  
import "github.com/worldflyingct/mnsql"  
然后  
go get github.com/worldflyingct/mnsql  
或是直接使用go mod命令  
go mod tidy  

# 关于datatype
因为存储在内存中仅仅是通过char*的通用格式进行保存的，所以并没有存储数据类型的方法，特意保留一个字段用来记录数据类型，如果使用c进行调用，不是一定要使用。  
但是由于incr与decr在key不存在时会创建int64类型的数据，并保存datatype为9，使用时请注意。  
