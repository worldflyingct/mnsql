package mnsql

/*
#cgo LDFLAGS: -static
#include "mnsql.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

var metex sync.Mutex

func cSet(key string, cdata unsafe.Pointer, datalen uint, ttl int, settype int, datatype int) int {
	keylen := len(key)
	if keylen == 0 {
		return -1
	}
	if datalen == 0 {
		return -2
	}
	var res int
	ckey := C.CString(key)
	metex.Lock()
	switch settype {
	case 0:
		res = int(C.Set(ckey, C.uint(keylen), cdata, C.uint(datalen), C.int(datatype)))
	case 1:
		res = int(C.SetEx(ckey, C.uint(keylen), cdata, C.uint(datalen), C.int(ttl), C.int(datatype)))
	case 2:
		res = int(C.SetNx(ckey, C.uint(keylen), cdata, C.uint(datalen), C.int(datatype)))
	case 3:
		res = int(C.SetNex(ckey, C.uint(keylen), cdata, C.uint(datalen), C.int(ttl), C.int(datatype)))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Set(key string, value interface{}, ttl int, settype int) int {
	switch data := value.(type) {
	case int:
		return cSet(key, unsafe.Pointer(&data), uint(unsafe.Sizeof(data)), ttl, settype, 0)
	case uint:
		return cSet(key, unsafe.Pointer(&data), uint(unsafe.Sizeof(data)), ttl, settype, 1)
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
	case string:
		cdata := unsafe.Pointer(C.CString(data))
		res := cSet(key, cdata, uint(len(data)), ttl, settype, 11)
		C.free(cdata)
		return res
	case []byte:
		return cSet(key, unsafe.Pointer(&data[0]), uint(len(data)), ttl, settype, 12)
	}
	return -3
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func Set(key string, value interface{}) int {
	return _Set(key, value, -1, 0)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func SetEx(key string, value interface{}, ttl int) int {
	return _Set(key, value, ttl, 1)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func SetNx(key string, value interface{}) int {
	return _Set(key, value, -1, 2)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func SetNex(key string, value interface{}, ttl int) int {
	return _Set(key, value, ttl, 3)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表对象不存在；-3代表未知的数据类型；
func Get(key string) (interface{}, int) {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return nil, -1
	}
	var datatype C.int
	datalen := C.uint(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	metex.Lock()
	defer metex.Unlock()
	datalen = C.Get(ckey, keylen, nil, &datalen, &datatype)
	if datalen == 0 {
		return nil, -2
	}
	switch datatype {
	case 0:
		var cdata int
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 1:
		var cdata uint
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 2:
		var cdata bool
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 3:
		var cdata int8
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 4:
		var cdata uint8
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 5:
		var cdata int16
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 6:
		var cdata uint16
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 7:
		var cdata int32
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 8:
		var cdata uint32
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 9:
		var cdata int64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 10:
		var cdata uint64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 11:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return string(cdata), 0
	case 12:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return cdata, 0
	}
	return nil, -3
}

// 返回定义：0代表成功；-1代表key长度为0；
func Del(key string) int {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return -1
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.Del(ckey, keylen))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

// 返回定义：0代表成功；-1代表key长度为0；
func Incr(key string) int {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return -1
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.Incr(ckey, keylen))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

// 返回定义：0代表成功；-1代表key长度为0；
func Decr(key string) int {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return -1
	}
	ckey := C.CString(key)
	metex.Lock()
	res := int(C.Decr(ckey, keylen))
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}
