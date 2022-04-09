package main

/*
#include <stdlib.h>
#include <stdint.h>
#include "mnsql.h"
*/
import "C"

import (
	"sync"
	"unsafe"
)

var metex sync.Mutex

func cSet(key string, cdata unsafe.Pointer, datalen uint, ttl int, settype int, datatype uint8) int {
	res := 0
	ckey := C.CString(key)
	switch settype {
	case 0:
		res = int(C.Set(ckey, C.uint(len(key)), cdata, C.uint(datalen), C.uint8_t(datatype)))
	case 1:
		res = int(C.SetEx(ckey, C.uint(len(key)), cdata, C.uint(datalen), C.int(ttl), C.uint8_t(datatype)))
	case 2:
		res = int(C.SetNx(ckey, C.uint(len(key)), cdata, C.uint(datalen), C.uint8_t(datatype)))
	case 3:
		res = int(C.SetNex(ckey, C.uint(len(key)), cdata, C.uint(datalen), C.int(ttl), C.uint8_t(datatype)))
	}
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Set(key string, value interface{}, ttl int, settype int) int {
	switch data := value.(type) {
	case int:
		clen := unsafe.Sizeof(data)
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, uint(clen), ttl, settype, 0)
	case uint:
		clen := unsafe.Sizeof(data)
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, uint(clen), ttl, settype, 1)
	case bool:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 1, ttl, settype, 2)
	case int8:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 1, ttl, settype, 3)
	case uint8:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 1, ttl, settype, 4)
	case int16:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 2, ttl, settype, 5)
	case uint16:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 2, ttl, settype, 6)
	case int32:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 4, ttl, settype, 7)
	case uint32:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 4, ttl, settype, 8)
	case int64:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 8, ttl, settype, 9)
	case uint64:
		cdata := unsafe.Pointer(&data)
		return cSet(key, cdata, 8, ttl, settype, 10)
	case string:
		d := C.CString(data)
		cdata := unsafe.Pointer(d)
		res := cSet(key, cdata, uint(len(data)), ttl, settype, 11)
		C.free(cdata)
		return res
	case []byte:
		return cSet(key, unsafe.Pointer(&data[0]), uint(len(data)), ttl, settype, 12)
	}
	return -1
}

func Set(key string, value interface{}) int {
	metex.Lock()
	defer metex.Unlock()
	return _Set(key, value, -1, 0)
}

func SetEx(key string, value interface{}, ttl int) int {
	metex.Lock()
	defer metex.Unlock()
	return _Set(key, value, ttl, 1)
}

func SetNx(key string, value interface{}) int {
	metex.Lock()
	defer metex.Unlock()
	return _Set(key, value, -1, 2)
}

func SetNex(key string, value interface{}, ttl int) int {
	metex.Lock()
	defer metex.Unlock()
	return _Set(key, value, ttl, 3)
}

func Get(key string) interface{} {
	metex.Lock()
	defer metex.Unlock()
	var data interface{}
	ckey := C.CString(key)
	keylen := C.uint(len(key))
	var datatype C.uint8_t
	datalen := C.uint(0)
	datalen = C.Get(ckey, keylen, nil, &datalen, &datatype)
	switch datatype {
	case 0:
		var cdata int
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 1:
		var cdata uint
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 2:
		var cdata bool
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 3:
		var cdata int8
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 4:
		var cdata uint8
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 5:
		var cdata int16
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 6:
		var cdata uint16
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 7:
		var cdata int32
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 8:
		var cdata uint32
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 9:
		var cdata int64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 10:
		var cdata uint64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		data = cdata
	case 11:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		data = string(cdata)
	case 12:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		data = cdata
	}
	C.free(unsafe.Pointer(ckey))
	return data
}

func Del(key string) int {
	metex.Lock()
	defer metex.Unlock()
	ckey := C.CString(key)
	res := int(C.Del(ckey, C.uint(len(key))))
	C.free(unsafe.Pointer(ckey))
	return res
}

func Incr(key string) int {
	metex.Lock()
	defer metex.Unlock()
	ckey := C.CString(key)
	res := int(C.Incr(ckey, C.uint(len(key))))
	C.free(unsafe.Pointer(ckey))
	return res
}

func Decr(key string) int {
	metex.Lock()
	defer metex.Unlock()
	ckey := C.CString(key)
	res := int(C.Decr(ckey, C.uint(len(key))))
	C.free(unsafe.Pointer(ckey))
	return res
}
