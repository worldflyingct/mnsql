package mnsql

/*
#include "mnsql.h"
*/
import "C"

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"unsafe"
)

var metex sync.Mutex

func typeToInterface(data unsafe.Pointer, datalen C.int, datatype C.int) (interface{}, int) {
	switch datatype {
	case 0:
		cdata := *(*int)(data)
		return cdata, 0
	case 1:
		cdata := *(*uint)(data)
		return cdata, 0
	case 2:
		cdata := *(*bool)(data)
		return cdata, 0
	case 3:
		cdata := *(*int8)(data)
		return cdata, 0
	case 4:
		cdata := *(*uint8)(data)
		return cdata, 0
	case 5:
		cdata := *(*int16)(data)
		return cdata, 0
	case 6:
		cdata := *(*uint16)(data)
		return cdata, 0
	case 7:
		cdata := *(*int32)(data)
		return cdata, 0
	case 8:
		cdata := *(*uint32)(data)
		return cdata, 0
	case 9:
		cdata := *(*int64)(data)
		return cdata, 0
	case 10:
		cdata := *(*uint64)(data)
		return cdata, 0
	case 11:
		cdata := *(*float32)(data)
		return cdata, 0
	case 12:
		cdata := *(*float64)(data)
		return cdata, 0
	case 13:
		cdata := *(*complex64)(data)
		return cdata, 0
	case 14:
		cdata := *(*complex128)(data)
		return cdata, 0
	case 15:
		cdata := C.GoStringN((*C.char)(data), datalen)
		return cdata, 0
	case 16:
		cdata := C.GoBytes(data, datalen)
		return cdata, 0
	}
	return nil, C.TYPEERROR
}

func cSet(key string, cdata unsafe.Pointer, datalen C.int, ttl int64, settype int, datatype C.int) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	if datalen == 0 {
		return C.DATANULL
	}
	var res int
	ckey := C.CString(key)
	metex.Lock()
	switch settype {
	case 0:
		res = int(C.Set(ckey, keylen, cdata, datalen, datatype))
	case 1:
		res = int(C.SetEx(ckey, keylen, cdata, datalen, C.int64_t(ttl), datatype))
	case 2:
		res = int(C.SetNx(ckey, keylen, cdata, datalen, C.int(datatype)))
	case 3:
		res = int(C.SetNex(ckey, keylen, cdata, datalen, C.int64_t(ttl), datatype))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Set(key string, value interface{}, ttl int64, settype int) int {
	switch data := value.(type) {
	case int:
		return cSet(key, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), ttl, settype, 0)
	case uint:
		return cSet(key, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), ttl, settype, 1)
	case bool:
		return cSet(key, unsafe.Pointer(&data), 1, ttl, settype, 2)
	case int8:
		return cSet(key, unsafe.Pointer(&data), 1, ttl, settype, 3)
	case uint8:
		return cSet(key, unsafe.Pointer(&data), 1, ttl, settype, 4)
	case int16:
		return cSet(key, unsafe.Pointer(&data), 2, ttl, settype, 5)
	case uint16:
		return cSet(key, unsafe.Pointer(&data), 2, ttl, settype, 6)
	case int32:
		return cSet(key, unsafe.Pointer(&data), 4, ttl, settype, 7)
	case uint32:
		return cSet(key, unsafe.Pointer(&data), 4, ttl, settype, 8)
	case int64:
		return cSet(key, unsafe.Pointer(&data), 8, ttl, settype, 9)
	case uint64:
		return cSet(key, unsafe.Pointer(&data), 8, ttl, settype, 10)
	case float32:
		return cSet(key, unsafe.Pointer(&data), 4, ttl, settype, 11)
	case float64:
		return cSet(key, unsafe.Pointer(&data), 8, ttl, settype, 12)
	case complex64:
		return cSet(key, unsafe.Pointer(&data), 8, ttl, settype, 13)
	case complex128:
		return cSet(key, unsafe.Pointer(&data), 16, ttl, settype, 14)
	case string:
		cdata := unsafe.Pointer(C.CString(data))
		res := cSet(key, cdata, C.int(len(data)), ttl, settype, 15)
		C.free(cdata)
		return res
	case []byte:
		cdata := C.CBytes(data)
		res := cSet(key, cdata, C.int(len(data)), ttl, settype, 16)
		C.free(cdata)
		return res
	}
	return C.TYPEERROR
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func Set(key string, value interface{}) int {
	return _Set(key, value, -1, 0)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func SetEx(key string, value interface{}, ttl int64) int {
	return _Set(key, value, ttl, 1)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func SetNx(key string, value interface{}) int {
	return _Set(key, value, -1, 2)
}

func SetNex(key string, value interface{}, ttl int64) int {
	return _Set(key, value, ttl, 3)
}

func Get(key string) (interface{}, int) {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return nil, C.KEYLENZERO
	}
	var datatype C.int
	var datalen C.int
	var res C.int
	ckey := C.CString(key)
	metex.Lock()
	data := C.Get(ckey, keylen, &datalen, &datatype, &res)
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	if res < 0 {
		return nil, int(res)
	}
	v, r := typeToInterface(data, datalen, datatype)
	C.free(data)
	return v, r
}

func Del(key string) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.Del(ckey, keylen))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func Incr(key string) int {
	return IncrBy(key, 1)
}

func IncrBy(key string, num int64) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.IncrBy(ckey, keylen, C.int64_t(num)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func Decr(key string) int {
	return DecrBy(key, 1)
}

func DecrBy(key string, num int64) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.DecrBy(ckey, keylen, C.int64_t(num)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func Expire(key string, ttl int64) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.Expire(ckey, keylen, C.int64_t(ttl)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func cPush(key string, cdata unsafe.Pointer, datalen C.int, settype int, datatype C.int) int {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return C.KEYLENZERO
	}
	if datalen == 0 {
		return C.DATANULL
	}
	var res int
	ckey := C.CString(key)
	metex.Lock()
	if settype == 0 {
		res = int(C.LPush(ckey, keylen, cdata, datalen, datatype))
	} else if settype == 1 {
		res = int(C.RPush(ckey, keylen, cdata, datalen, datatype))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Push(key string, value interface{}, settype int) int {
	switch data := value.(type) {
	case int:
		return cPush(key, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), settype, 0)
	case uint:
		return cPush(key, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), settype, 1)
	case bool:
		return cPush(key, unsafe.Pointer(&data), 1, settype, 2)
	case int8:
		return cPush(key, unsafe.Pointer(&data), 1, settype, 3)
	case uint8:
		return cPush(key, unsafe.Pointer(&data), 1, settype, 4)
	case int16:
		return cPush(key, unsafe.Pointer(&data), 2, settype, 5)
	case uint16:
		return cPush(key, unsafe.Pointer(&data), 2, settype, 6)
	case int32:
		return cPush(key, unsafe.Pointer(&data), 4, settype, 7)
	case uint32:
		return cPush(key, unsafe.Pointer(&data), 4, settype, 8)
	case int64:
		return cPush(key, unsafe.Pointer(&data), 8, settype, 9)
	case uint64:
		return cPush(key, unsafe.Pointer(&data), 8, settype, 10)
	case float32:
		return cPush(key, unsafe.Pointer(&data), 4, settype, 11)
	case float64:
		return cPush(key, unsafe.Pointer(&data), 8, settype, 12)
	case complex64:
		return cPush(key, unsafe.Pointer(&data), 8, settype, 13)
	case complex128:
		return cPush(key, unsafe.Pointer(&data), 16, settype, 14)
	case string:
		cdata := unsafe.Pointer(C.CString(data))
		res := cPush(key, cdata, C.int(len(data)), settype, 15)
		C.free(cdata)
		return res
	case []byte:
		cdata := C.CBytes(data)
		res := cPush(key, cdata, C.int(len(data)), settype, 16)
		C.free(cdata)
		return res
	}
	return C.TYPEERROR
}

func LPush(key string, value interface{}) int {
	return _Push(key, value, 0)
}

func RPush(key string, value interface{}) int {
	return _Push(key, value, 1)
}

func _Pop(key string, settype int) (interface{}, int) {
	keylen := C.uint64_t(len(key))
	if keylen == 0 {
		return nil, C.KEYLENZERO
	}
	var datatype C.int
	var datalen C.int
	var res C.int
	var data unsafe.Pointer
	ckey := C.CString(key)
	metex.Lock()
	if settype == 0 {
		data = C.LPop(ckey, keylen, &datalen, &datatype, &res)
	} else if settype == 1 {
		data = C.RPop(ckey, keylen, &datalen, &datatype, &res)
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	if res < 0 {
		return nil, int(res)
	}
	v, r := typeToInterface(data, datalen, datatype)
	C.free(data)
	return v, r
}

func LPop(key string) (interface{}, int) {
	return _Pop(key, 0)
}

func RPop(key string) (interface{}, int) {
	return _Pop(key, 1)
}

func cHSet(key string, key2 string, cdata unsafe.Pointer, datalen C.int, ttl int64, settype int, datatype C.int) int {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return C.KEYLENZERO
	}
	if datalen == 0 {
		return C.DATANULL
	}
	var res int
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	switch settype {
	case 0:
		res = int(C.HSet(ckey, keylen, ckey2, keylen2, cdata, datalen, datatype))
	case 1:
		res = int(C.HSetEx(ckey, keylen, ckey2, keylen2, cdata, datalen, C.int64_t(ttl), datatype))
	case 2:
		res = int(C.HSetNx(ckey, keylen, ckey2, keylen2, cdata, datalen, datatype))
	case 3:
		res = int(C.HSetNex(ckey, keylen, ckey2, keylen2, cdata, datalen, C.int64_t(ttl), datatype))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func _HSet(key string, key2 string, value interface{}, ttl int64, settype int) int {
	switch data := value.(type) {
	case int:
		return cHSet(key, key2, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), ttl, settype, 0)
	case uint:
		return cHSet(key, key2, unsafe.Pointer(&data), C.int(unsafe.Sizeof(data)), ttl, settype, 1)
	case bool:
		return cHSet(key, key2, unsafe.Pointer(&data), 1, ttl, settype, 2)
	case int8:
		return cHSet(key, key2, unsafe.Pointer(&data), 1, ttl, settype, 3)
	case uint8:
		return cHSet(key, key2, unsafe.Pointer(&data), 1, ttl, settype, 4)
	case int16:
		return cHSet(key, key2, unsafe.Pointer(&data), 2, ttl, settype, 5)
	case uint16:
		return cHSet(key, key2, unsafe.Pointer(&data), 2, ttl, settype, 6)
	case int32:
		return cHSet(key, key2, unsafe.Pointer(&data), 4, ttl, settype, 7)
	case uint32:
		return cHSet(key, key2, unsafe.Pointer(&data), 4, ttl, settype, 8)
	case int64:
		return cHSet(key, key2, unsafe.Pointer(&data), 8, ttl, settype, 9)
	case uint64:
		return cHSet(key, key2, unsafe.Pointer(&data), 8, ttl, settype, 10)
	case float32:
		return cHSet(key, key2, unsafe.Pointer(&data), 4, ttl, settype, 11)
	case float64:
		return cHSet(key, key2, unsafe.Pointer(&data), 8, ttl, settype, 12)
	case complex64:
		return cHSet(key, key2, unsafe.Pointer(&data), 8, ttl, settype, 13)
	case complex128:
		return cHSet(key, key2, unsafe.Pointer(&data), 16, ttl, settype, 14)
	case string:
		cdata := unsafe.Pointer(C.CString(data))
		res := cHSet(key, key2, cdata, C.int(len(data)), ttl, settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cHSet(key, key2, unsafe.Pointer(&data[0]), C.int(len(data)), ttl, settype, 16)
	}
	return C.TYPEERROR
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSet(key string, key2 string, value interface{}) int { return _HSet(key, key2, value, -1, 0) }

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSetEx(key string, key2 string, value interface{}, ttl int64) int {
	return _HSet(key, key2, value, ttl, 1)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSetNx(key string, key2 string, value interface{}) int { return _HSet(key, key2, value, -1, 2) }

func HSetNex(key string, key2 string, value interface{}, ttl int64) int {
	return _HSet(key, key2, value, ttl, 3)
}

func HGet(key string, key2 string) (interface{}, int) {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return nil, C.KEYLENZERO
	}
	var datatype C.int
	var datalen C.int
	var res C.int
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	data := C.HGet(ckey, keylen, ckey2, keylen2, &datalen, &datatype, &res)
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	if res < 0 {
		return nil, int(res)
	}
	v, r := typeToInterface(data, datalen, datatype)
	C.free(data)
	return v, r
}

func HDel(key string, key2 string) int {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	res := int(C.HDel(ckey, keylen, ckey2, keylen2))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func HIncr(key string, key2 string) int {
	return HIncrBy(key, key2, 1)
}

func HIncrBy(key string, key2 string, num int64) int {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	res := int(C.HIncrBy(ckey, keylen, ckey2, keylen2, C.int64_t(num)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func HDecr(key string, key2 string) int {
	return HDecrBy(key, key2, 1)
}

func HDecrBy(key string, key2 string, num int64) int {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	res := int(C.HDecrBy(ckey, keylen, ckey2, keylen2, C.int64_t(num)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func HExpire(key string, key2 string, ttl int64) int {
	keylen := C.uint64_t(len(key))
	keylen2 := C.uint64_t(len(key2))
	if keylen == 0 || keylen2 == 0 {
		return C.KEYLENZERO
	}
	ckey := C.CString(key)
	ckey2 := C.CString(key2)
	metex.Lock()
	res := int(C.HExpire(ckey, keylen, ckey2, keylen2, C.int64_t(ttl)))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func StartCmdLineServer(port uint16) (net.Listener, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	go func(ln net.Listener) {
		fmt.Println("CmdLine Server Start")
		for {
			conn, err := ln.Accept()
			if err != nil {
				break
			}
			go cmdLineHandle(conn, ln)
		}
		ln.Close()
		fmt.Println("CmdLine Server Stop")
	}(ln)
	return ln, nil
}

func cmdLineHandle(conn net.Conn, ln net.Listener) error {
	defer conn.Close()
	in := bufio.NewReader(conn)
	for {
		dat, _, err := in.ReadLine()
		if err != nil {
			return err
		}
		cmd := strings.Split(string(dat), " ")
		cmdlen := len(cmd)
		if cmdlen > 0 {
			switch cmd[0] {
			case "keys":
				var datalen C.int
				var res C.int
				metex.Lock()
				data := C.Keys(&datalen, &res)
				metex.Unlock()
				if res < 0 {
					fmt.Fprintln(conn, "malloc fail", res)
					continue
				}
				cdata := C.GoStringN(data, datalen)
				C.free(unsafe.Pointer(data))
				fmt.Fprintln(conn, cdata, 0)
			case "get":
				if cmdlen > 1 {
					v, r := Get(cmd[1])
					fmt.Fprintln(conn, v, r)
				}
			case "hget":
				if cmdlen > 2 {
					v, r := HGet(cmd[1], cmd[2])
					fmt.Fprintln(conn, v, r)
				}
			case "lrange":
				if cmdlen > 1 {
					arr := make([]interface{}, 0)
					for {
						v, r := RPop(cmd[1])
						if r == 0 {
							arr = append(arr, v)
							fmt.Fprintln(conn, v)
						} else {
							break
						}
					}
					arrlen := len(arr)
					for i := 0; i < arrlen; i++ {
						LPush(cmd[1], arr[i])
					}
				}
			case "hkeys":
				if cmdlen > 1 {
					var datalen C.int
					var res C.int
					key := C.CString(cmd[1])
					metex.Lock()
					data := C.HKeys(key, C.uint64_t(len(cmd[1])), &datalen, &res)
					metex.Unlock()
					C.free(unsafe.Pointer(key))
					if res < 0 {
						if res == C.DATANULL {
							fmt.Fprintln(conn, "data is null", res)
						} else if res == C.TYPEERROR {
							fmt.Fprintln(conn, "type is err", res)
						} else {
							fmt.Fprintln(conn, "malloc fail", res)
						}
						continue
					}
					cdata := C.GoStringN(data, datalen)
					C.free(unsafe.Pointer(data))
					fmt.Fprintln(conn, cdata, 0)
				}
			case "exit":
				return nil
			case "closeserver":
				ln.Close()
			}
		}
	}
	return nil
}
