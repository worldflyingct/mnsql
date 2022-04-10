#ifndef __MNSQL_H__
#define __MNSQL_H__

#include <stdint.h>
#include <stdlib.h>

#define RUNSUCCESS 0
#define KEYLENZERO -1
#define DATANULL -2
#define TYPEERROR -3
#define WRITEDENY -5
#define MALLOCFAIL -6

#ifdef __cplusplus
extern "C"
{
#endif

    int Set(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype);
    int SetEx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl, int datatype);
    int SetNx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype);
    int SetNex(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl, int datatype);
    int Get(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen, int *datatype);
    int Del(const char *mkey, unsigned int keylen);
    int Incr(const char *mkey, unsigned int keylen);
    int IncrBy(const char *mkey, unsigned int keylen, int64_t num);
    int Decr(const char *mkey, unsigned int keylen);
    int DecrBy(const char *mkey, unsigned int keylen, int64_t num);
    int Expire(const char *mkey, unsigned int keylen, int ttl);
    int Lpop(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen, int *datatype);
    int Rpop(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen, int *datatype);
    int Lpush(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype);
    int Rpush(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype);

#ifdef __cplusplus
}
#endif

#endif
