#ifndef __MNSQL_H__
#define __MNSQL_H__

#ifdef __cplusplus
extern "C"
{
#endif

    int Set(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen);
    int SetEx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl);
    int SetNx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen);
    int SetNex(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl);
    int Get(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen);
    int Del(const char *mkey, unsigned int keylen);
    int Incr(const char *mkey, unsigned int keylen);
    int Decr(const char *mkey, unsigned int keylen);

#ifdef __cplusplus
}
#endif

#endif
