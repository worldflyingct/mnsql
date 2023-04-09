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

    int Set(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int datatype);
    int SetEx(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int64_t ttl, int datatype);
    int SetNx(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int datatype);
    int SetNex(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int64_t ttl, int datatype);
    void *Get(const char *mkey, uint64_t keylen, int *datalen, int *datatype, int *res);
    int Del(const char *mkey, uint64_t keylen);
    int Incr(const char *mkey, uint64_t keylen);
    int IncrBy(const char *mkey, uint64_t keylen, int64_t num);
    int Decr(const char *mkey, uint64_t keylen);
    int DecrBy(const char *mkey, uint64_t keylen, int64_t num);
    int Expire(const char *mkey, uint64_t keylen, int64_t ttl);
    int LPush(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int datatype);
    void *LPop(const char *mkey, uint64_t keylen, int *datalen, int *datatype, int *res);
    int RPush(const char *mkey, uint64_t keylen, const void *mdata, int datalen, int datatype);
    void *RPop(const char *mkey, uint64_t keylen, int *datalen, int *datatype, int *res);
    int HSet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, int datalen,
             int datatype);
    int HSetEx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, int datalen,
               int64_t ttl, int datatype);
    int HSetNx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, int datalen,
               int datatype);
    int HSetNex(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, int datalen,
                int64_t ttl, int datatype);
    void *HGet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int *datalen, int *datatype,
               int *res);
    int HDel(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HIncr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HIncrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num);
    int HDecr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2);
    int HDecrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num);
    int HExpire(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t ttl);
    char *Keys(int *datalen, int *res);
    char *Lrange(char *mkey, uint64_t keylen, int *datalen, int *res);
    char *HKeys(char *mkey, uint64_t keylen, int *datalen, int *res);
    void FlushDB();

#ifdef __cplusplus
}
#endif

#endif
