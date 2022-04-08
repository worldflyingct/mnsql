#include <stdlib.h>
#include <string.h>
#include <time.h>

struct PARAM
{
    char *key;
    unsigned int keylen;
    void *data;
    unsigned int datalen;
    unsigned int deadline;
    struct PARAM *head;
    struct PARAM *tail;
};
struct PARAM *globalparam[256];

struct PARAM *FindKey(const char *mkey, unsigned int keylen)
{
    unsigned int hash = mkey[0];
    struct PARAM *param = globalparam[hash];
    while (param != NULL)
    {
        time_t seconds = time(NULL);
        if (seconds > param->deadline)
        {
            struct PARAM *p = param;
            param = param->tail;
            unsigned int hash = mkey[0];
            if (p->head != NULL)
            {
                p->head->tail = p->tail;
            }
            else
            {
                globalparam[hash] = p->tail;
            }
            if (p->tail != NULL)
            {
                p->tail->head = p->head;
            }
            free(p->key);
            free(p->data);
            free(p);
        }
        else if (param->keylen == keylen && !memcmp(param->key, mkey, keylen))
        {
            return param;
        }
        else
        {
            param = param->tail;
        }
    }
    return NULL;
}

void MoveParam(struct PARAM *param, unsigned int hash)
{
    if (param->head != NULL)
    {
        param->head->tail = param->tail;
    }
    else
    {
        globalparam[hash] = param->tail;
    }
    if (param->tail != NULL)
    {
        param->tail->head = param->head;
    }
    struct PARAM *gparam = globalparam[hash];
    param->tail = gparam;
    param->head = NULL;
    if (gparam != NULL)
    {
        gparam->head = param;
    }
    globalparam[hash] = param;
}

int AddKey(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl)
{
    unsigned int hash = mkey[0];
    struct PARAM *gparam = globalparam[hash];
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
    time_t seconds = time(NULL);
    if (ttl == -1)
    {
        param->deadline = -1;
    }
    else
    {
        param->deadline = seconds + ttl;
    }
    param->head = NULL;
    param->tail = gparam;
    if (gparam != NULL)
    {
        gparam->head = param;
    }
    globalparam[hash] = param;
    return 0;
}

int SetEx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        unsigned int hash = mkey[0];
        if (param->head != NULL)
        {
            param->head->tail = param->tail;
        }
        else
        {
            globalparam[hash] = param->tail;
        }
        if (param->tail != NULL)
        {
            param->tail->head = param->head;
        }
        free(param->key);
        free(param->data);
        free(param);
    }
    return AddKey(mkey, keylen, mdata, datalen, ttl);
}

int Set(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen)
{
    return SetEx(mkey, keylen, mdata, datalen, -1);
}

int SetNx(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return 0;
    }
    return AddKey(mkey, keylen, mdata, datalen, -1);
}

int SetNex(const char *mkey, unsigned int keylen, const void *mdata, unsigned int datalen, int ttl)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return 0;
    }
    return AddKey(mkey, keylen, mdata, datalen, ttl);
}

int Get(const char *mkey, unsigned int keylen, void *mdata, unsigned int *datalen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return 0;
    }
    MoveParam(param, mkey[0]);
    unsigned int dlen = param->datalen > *datalen ? *datalen : param->datalen;
    memcpy(mdata, param->data, dlen);
    *datalen = dlen;
    return param->datalen;
}

int Del(const char *mkey, unsigned int keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        unsigned int hash = mkey[0];
        if (param->head != NULL)
        {
            param->head->tail = param->tail;
        }
        else
        {
            globalparam[hash] = param->tail;
        }
        if (param->tail != NULL)
        {
            param->tail->head = param->head;
        }
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
        int n = 1;
        return AddKey(mkey, keylen, &n, sizeof(int), -1);
    }
    MoveParam(param, mkey[0]);
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

int Decr(char *mkey, unsigned int keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        int n = -1;
        return AddKey(mkey, keylen, &n, sizeof(int), -1);
    }
    MoveParam(param, mkey[0]);
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
