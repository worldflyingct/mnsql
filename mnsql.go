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
		res := cSet(key, cdata, uint(len(data)), ttl, settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cSet(key, unsafe.Pointer(&data[0]), uint(len(data)), ttl, settype, 16)
	}
	return C.TYPEERROR
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

func SetNex(key string, value interface{}, ttl int) int {
	return _Set(key, value, ttl, 3)
}

func Get(key string) (interface{}, int) {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return nil, C.KEYLENZERO
	}
	var datatype C.int
	datalen := C.uint(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	metex.Lock()
	defer metex.Unlock()
	res := int(C.Get(ckey, keylen, nil, &datalen, &datatype))
	if res < 0 {
		return nil, res
	}
	datalen = C.uint(res)
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
	keylen := C.uint(len(key))
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
	keylen := C.uint(len(key))
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
	keylen := C.uint(len(key))
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

func cPush(key string, cdata unsafe.Pointer, datalen C.uint, settype int, datatype C.int) int {
	keylen := C.uint(len(key))
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
		res = int(C.Lpush(ckey, keylen, cdata, datalen, datatype))
	} else if settype == 1 {
		res = int(C.Rpush(ckey, keylen, cdata, datalen, datatype))
	}
	metex.Unlock()
	C.free(unsafe.Pointer(ckey))
	return res
}

func _Push(key string, value interface{}, settype int) int {
	switch data := value.(type) {
	case int:
		return cPush(key, unsafe.Pointer(&data), C.uint(unsafe.Sizeof(data)), settype, 0)
	case uint:
		return cPush(key, unsafe.Pointer(&data), C.uint(unsafe.Sizeof(data)), settype, 1)
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
		res := cPush(key, cdata, C.uint(len(data)), settype, 15)
		C.free(cdata)
		return res
	case []byte:
		return cPush(key, unsafe.Pointer(&data[0]), C.uint(len(data)), settype, 16)
	}
	return C.TYPEERROR
}

func Lpush(key string, value interface{}) int {
	return _Push(key, value, 0)
}

func Rpush(key string, value interface{}) int {
	return _Push(key, value, 1)
}

func _Pop(key string, settype int) (interface{}, int) {
	keylen := C.uint(len(key))
	if keylen == 0 {
		return nil, C.KEYLENZERO
	}
	var res int
	var datatype C.int
	datalen := C.uint(0)
	ckey := C.CString(key)
	defer C.free(unsafe.Pointer(ckey))
	metex.Lock()
	defer metex.Unlock()
	if settype == 0 {
		res = int(C.Lpop(ckey, keylen, nil, &datalen, &datatype))
	} else if settype == 1 {
		res = int(C.Rpop(ckey, keylen, nil, &datalen, &datatype))
	}
	if res < 0 {
		return nil, res
	}
	datalen = C.uint(res)
	switch datatype {
	case 0:
		var cdata int
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 1:
		var cdata uint
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 2:
		var cdata bool
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 3:
		var cdata int8
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 4:
		var cdata uint8
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 5:
		var cdata int16
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 6:
		var cdata uint16
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 7:
		var cdata int32
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 8:
		var cdata uint32
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 9:
		var cdata int64
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 10:
		var cdata uint64
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 11:
		var cdata float32
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 12:
		var cdata float64
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 13:
		var cdata complex64
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 14:
		var cdata complex128
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata), &datalen, &datatype))
		}
		return cdata, 0
	case 15:
		cdata := make([]byte, datalen)
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype))
		}
		return string(cdata), 0
	case 16:
		cdata := make([]byte, datalen)
		if settype == 0 {
			res = int(C.Lpop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype))
		} else if settype == 1 {
			res = int(C.Rpop(ckey, keylen, unsafe.Pointer(&cdata[0]), &datalen, &datatype))
		}
		return cdata, 0
	}
	return nil, C.TYPEERROR
}

func Lpop(key string) (interface{}, int) {
	return _Pop(key, 0)
}

func Rpop(key string) (interface{}, int) {
	return _Pop(key, 1)
}
