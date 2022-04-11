#ifndef __MNSQL_H__
#define __MNSQL_H__

#include <stdint.h>
#include <stdlib.h>

#define RUNSUCCESS 0
#define KEYLENZERO -1
#define DATANULL -2
#define TYPEERROR -3
#define WRITEDENY -4
#define MALLOCFAIL -5

#ifdef __cplusplus
extern "C"
{
#endif

    int Set(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype);
    int SetEx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype);
    int SetNx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype);
    int SetNex(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype);
    int64_t Get(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype);
    int Del(const char *mkey, uint64_t keylen);
    int Incr(const char *mkey, uint64_t keylen);
    int IncrBy(const char *mkey, uint64_t keylen, int64_t num);
    int Decr(const char *mkey, uint64_t keylen);
    int DecrBy(const char *mkey, uint64_t keylen, int64_t num);
    int Expire(const char *mkey, uint64_t keylen, int64_t ttl);
    int LPush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype);
    int64_t LPop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype);
    int RPush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype);
    int64_t RPop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype);
    int HSet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata,
             uint64_t datalen, int datatype);
    int HSetEx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata,
               uint64_t datalen, int64_t ttl, int datatype);
    int HSetNx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata,
               uint64_t datalen, int datatype);
    int HSetNex(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata,
                uint64_t datalen, int64_t ttl, int datatype);
    int64_t HGet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, void *mdata, uint64_t *datalen,
                 int *datatype);
    int HDel(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HIncr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HIncrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num);
    int HDecr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HDecrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num);
    int HExpire(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t ttl);

#ifdef __cplusplus
}
#endif

#endif
