package mnsql

/*
#cgo LDFLAGS: -static
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

func cSet(key string, cdata unsafe.Pointer, datalen uint64, ttl int64, settype int, datatype int) int {
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
		res = int(C.Set(ckey, keylen, cdata, C.uint64_t(datalen), C.int(datatype)))
	case 1:
		res = int(C.SetEx(ckey, keylen, cdata, C.uint64_t(datalen), C.int64_t(ttl), C.int(datatype)))
	case 2:
		res = int(C.SetNx(ckey, keylen, cdata, C.uint64_t(datalen), C.int(datatype)))
	case 3:
		res = int(C.SetNex(ckey, keylen, cdata, C.uint64_t(datalen), C.int64_t(ttl), C.int(datatype)))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Set(key string, value interface{}, ttl int64, settype int) int {
	switch data := value.(type) {
	case int:
		return cSet(key, unsafe.Pointer(&data), uint64(unsafe.Sizeof(data)), ttl, settype, 0)
	case uint:
		return cSet(key, unsafe.Pointer(&data), uint64(unsafe.Sizeof(data)), ttl, settype, 1)
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
		res := cSet(key, cdata, uint64(len(data)), ttl, settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cSet(key, unsafe.Pointer(&data[0]), uint64(len(data)), ttl, settype, 16)
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
	datalen := C.uint64_t(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	metex.Lock()
	defer metex.Unlock()
	res := int64(C.Get(ckey, keylen, nil, &datalen, &datatype))
	if res < 0 {
		return nil, int(res)
	}
	datalen = C.uint64_t(res)
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
		var cdata float32
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 12:
		var cdata float64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 13:
		var cdata complex64
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 14:
		var cdata complex128
		C.Get(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 15:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return string(cdata), 0
	case 16:
		cdata := make([]byte, datalen)
		C.Get(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return cdata, 0
	}
	return nil, C.TYPEERROR
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

func cPush(key string, cdata unsafe.Pointer, datalen C.uint64_t, settype int, datatype C.int) int {
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
		return cPush(key, unsafe.Pointer(&data), C.uint64_t(unsafe.Sizeof(data)), settype, 0)
	case uint:
		return cPush(key, unsafe.Pointer(&data), C.uint64_t(unsafe.Sizeof(data)), settype, 1)
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
		res := cPush(key, cdata, C.uint64_t(len(data)), settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cPush(key, unsafe.Pointer(&data[0]), C.uint64_t(len(data)), settype, 16)
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
	var res int64
	var datatype C.int
	datalen := C.uint64_t(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	metex.Lock()
	defer metex.Unlock()
	if settype == 0 {
		res = int64(C.LPop(ckey, keylen, nil, &datalen, &datatype))
	} else if settype == 1 {
		res = int64(C.RPop(ckey, keylen, nil, &datalen, &datatype))
	}
	if res < 0 {
		return nil, int(res)
	}
	datalen = C.uint64_t(res)
	switch datatype {
	case 0:
		var cdata int
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 1:
		var cdata uint
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 2:
		var cdata bool
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 3:
		var cdata int8
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 4:
		var cdata uint8
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 5:
		var cdata int16
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 6:
		var cdata uint16
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 7:
		var cdata int32
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 8:
		var cdata uint32
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 9:
		var cdata int64
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 10:
		var cdata uint64
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 11:
		var cdata float32
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 12:
		var cdata float64
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 13:
		var cdata complex64
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 14:
		var cdata complex128
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype)
		}
		return cdata, 0
	case 15:
		cdata := make([]byte, datalen)
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		}
		return string(cdata), 0
	case 16:
		cdata := make([]byte, datalen)
		if settype == 0 {
			C.LPop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		} else if settype == 1 {
			C.RPop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		}
		return cdata, 0
	}
	return nil, C.TYPEERROR
}

func LPop(key string) (interface{}, int) {
	return _Pop(key, 0)
}

func RPop(key string) (interface{}, int) {
	return _Pop(key, 1)
}

func cHSet(key string, key2 string, cdata unsafe.Pointer, datalen uint64, ttl int64, settype int, datatype int) int {
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
		res = int(C.HSet(ckey, keylen, ckey2, keylen2, cdata, C.uint64_t(datalen), C.int(datatype)))
	case 1:
		res = int(C.HSetEx(ckey, keylen, ckey2, keylen2, cdata, C.uint64_t(datalen), C.int64_t(ttl), C.int(datatype)))
	case 2:
		res = int(C.HSetNx(ckey, keylen, ckey2, keylen2, cdata, C.uint64_t(datalen), C.int(datatype)))
	case 3:
		res = int(C.HSetNex(ckey, keylen, ckey2, keylen2, cdata, C.uint64_t(datalen), C.int64_t(ttl), C.int(datatype)))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey2))
	C.free(unsafe.Pointer(ckey))
	return res
}

func _HSet(key string, key2 string, value interface{}, ttl int64, settype int) int {
	switch data := value.(type) {
	case int:
		return cHSet(key, key2, unsafe.Pointer(&data), uint64(unsafe.Sizeof(data)), ttl, settype, 0)
	case uint:
		return cHSet(key, key2, unsafe.Pointer(&data), uint64(unsafe.Sizeof(data)), ttl, settype, 1)
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
		res := cHSet(key, key2, cdata, uint64(len(data)), ttl, settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cHSet(key, key2, unsafe.Pointer(&data[0]), uint64(len(data)), ttl, settype, 16)
	}
	return C.TYPEERROR
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSet(key string, key2 string, value interface{}) int {
	return _HSet(key, key2, value, -1, 0)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSetEx(key string, key2 string, value interface{}, ttl int64) int {
	return _HSet(key, key2, value, ttl, 1)
}

// 返回定义：0代表成功；-1代表key长度为0；-2代表value长度为0；-3代表value不是支持的类型；
func HSetNx(key string, key2 string, value interface{}) int {
	return _HSet(key, key2, value, -1, 2)
}

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
	datalen := C.uint64_t(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	ckey2 := C.CString(key2)
	defer C.free(unsafe.Pointer(ckey2))
	metex.Lock()
	defer metex.Unlock()
	res := int64(C.HGet(ckey, keylen, ckey2, keylen2, nil, &datalen, &datatype))
	if res < 0 {
		return nil, int(res)
	}
	datalen = C.uint64_t(res)
	switch datatype {
	case 0:
		var cdata int
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 1:
		var cdata uint
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 2:
		var cdata bool
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 3:
		var cdata int8
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 4:
		var cdata uint8
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 5:
		var cdata int16
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 6:
		var cdata uint16
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 7:
		var cdata int32
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 8:
		var cdata uint32
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 9:
		var cdata int64
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 10:
		var cdata uint64
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 11:
		var cdata float32
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 12:
		var cdata float64
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 13:
		var cdata complex64
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 14:
		var cdata complex128
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata), &datalen, &datatype)
		return cdata, 0
	case 15:
		cdata := make([]byte, datalen)
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return string(cdata), 0
	case 16:
		cdata := make([]byte, datalen)
		C.HGet(ckey, keylen, ckey2, keylen2, unsafe.Pointer(&cdata[0]), &datalen, &datatype)
		return cdata, 0
	}
	return nil, C.TYPEERROR
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
				fmt.Println(err.Error())
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
			fmt.Println(err.Error())
			return err
		}
		cmd := strings.Split(string(dat), " ")
		cmdlen := len(cmd)
		if cmdlen > 0 {
			switch cmd[0] {
			case "get":
				if cmdlen > 1 {
					fmt.Println(cmd)
					v, r := Get(cmd[1])
					fmt.Fprintln(conn, v, r)
				}
			case "hget":
				if cmdlen > 2 {
					v, r := HGet(cmd[1], cmd[2])
					fmt.Fprintln(conn, v, r)
				}
			case "listget":
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
			case "exit":
				return nil
			case "closeserver":
				ln.Close()
			}
		}
	}
	return nil
}
