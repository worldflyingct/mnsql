#include "mnsql.h"
#include <string.h>
#include <time.h>

struct PARAM
{
    char *key;
    unsigned int keylen;
    int datatype;
    void *data;
    unsigned int datalen;
    unsigned int deadline;
    struct PARAM *tail;
};
struct PARAM *globalparam[256];

struct PARAM *FindKey(const char *mkey, unsigned int keylen)
{
    if (keylen == 0)
    {
        return NULL;
    }
    unsigned int hash = mkey[0];
    struct PARAM *param = globalparam[hash];
    struct PARAM *beforeparam = NULL;
    while (param != NULL)
    {
        time_t now = time(NULL);
        if (now > param->deadline)
        {
            struct PARAM *p = param;
            param = param->tail;
            unsigned int hash = mkey[0];
            if (beforeparam != NULL)
            {
                beforeparam->tail = param;
            }
            else
            {
                globalparam[hash] = param;
            }
            free(p->key);
            free(p->data);
            free(p);
        }
        else if (param->keylen == keylen && !memcmp(param->key, mkey, keylen))
        {
            if (beforeparam != NULL)
            {
                beforeparam->tail = param->tail;
                param->tail = globalparam[hash];
                globalparam[hash] = param;
            }
            return param;
        }
        else
        {
            beforeparam = param;
            param = param->tail;
        }
    }
    return NULL;
}

int AddKey(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl, int datatype)
{
    if (keylen == 0)
    {
        return __LINE__;
    }
    if (datalen == 0)
    {
        return __LINE__;
    }
    unsigned int hash = mkey[0];
    char *key = (char *)malloc(keylen);
    if (key == NULL)
    {
        return __LINE__;
    }
    void *data = (void *)malloc(datalen);
    if (data == NULL)
    {
        free(key);
        return __LINE__;
    }
    struct PARAM *param = (struct PARAM *)malloc(sizeof(struct PARAM));
    if (param == NULL)
    {
        free(key);
        free(data);
        return __LINE__;
    }
    memcpy(key, mkey, keylen);
    memcpy(data, mdata, datalen);
    param->key = key;
    param->keylen = keylen;
    param->data = data;
    param->datalen = datalen;
    time_t now = time(NULL);
    if (ttl == -1)
    {
        param->deadline = -1;
    }
    else
    {
        param->deadline = now + ttl;
    }
    param->datatype = datatype;
    param->tail = globalparam[hash];
    globalparam[hash] = param;
    return 0;
}

int SetEx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl, int datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        unsigned int hash = mkey[0];
        globalparam[hash] = param->tail;
        free(param->key);
        free(param->data);
        free(param);
    }
    return AddKey(mkey, keylen, mdata, datalen, ttl, datatype);
}

int Set(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype)
{
    return SetEx(mkey, keylen, mdata, datalen, -1, datatype);
}

int SetNx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int datatype)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return 0;
    }
    return AddKey(mkey, keylen, mdata, datalen, -1, datatype);
}

int SetNex(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl, int datatype)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return 0;
    }
    return AddKey(mkey, keylen, mdata, datalen, ttl, datatype);
}

unsigned int Get(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen, int *datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return 0;
    }
    unsigned int dlen = param->datalen > *datalen ? *datalen : param->datalen;
    if (dlen > 0 && mdata != NULL)
    {
        memcpy(mdata, param->data, dlen);
        *datalen = dlen;
    }
    *datatype = param->datatype;
    return param->datalen;
}

int Del(const char *mkey, unsigned int keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        unsigned int hash = mkey[0];
        globalparam[hash] = param->tail;
        free(param->key);
        free(param->data);
        free(param);
    }
    return 0;
}

int Incr(const char *mkey, unsigned int keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        int64_t n = 1;
        return AddKey(mkey, keylen, &n, sizeof(int64_t), -1, 9);
    }
    unsigned int datalen = param->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)param->data;
        n++;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)param->data;
        n++;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)param->data;
        n++;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)param->data;
        n++;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)param->data;
        n++;
        memcpy(param->data, &n, datalen);
    }
    return 0;
}

int Decr(const char *mkey, unsigned int keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        int64_t n = -1;
        return AddKey(mkey, keylen, &n, sizeof(int64_t), -1, 9);
    }
    unsigned int datalen = param->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)param->data;
        n--;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)param->data;
        n--;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)param->data;
        n--;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)param->data;
        n--;
        memcpy(param->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)param->data;
        n--;
        memcpy(param->data, &n, datalen);
    }
    return 0;
}
