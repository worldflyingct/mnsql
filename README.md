# mnsql是memory nosql
仅仅使用mnsql.c与mnsql.h即可。  

# 存在如下函数
写入对象，但是如果存在相同的则覆盖写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func Set(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
但是如果存在相同的则覆盖写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func SetEx(key string, value interface{}) int  

写入对象，但是如果存在相同的key则不写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
func SetNx(key string, value interface{}) int  

写入有生命期的对象，生命期消失后数据自动消失  
如果存在相同的key则不写入  
返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；  
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
