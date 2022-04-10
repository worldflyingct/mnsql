#include "mnsql.h"
#include <string.h>
#include <time.h>

#include <stdio.h>

struct SINGLE
{
    int datatype;
    void *data;
    uint64_t datalen;
};

struct List
{
    char *data;
    uint64_t datalen;
    int datatype;
    struct List *head;
    struct List *tail;
};

struct PARAM
{
    char *key;
    uint64_t keylen;
    int type;
    void *value;
    uint64_t deadline;
    struct PARAM *tail;
};
struct PARAM *globalparam[256];

void FreeSingleItem(struct PARAM *param)
{
    struct SINGLE *item = param->value;
    free(item->data);
    free(param->value);
    free(param->key);
    free(param);
}

void FreeListItem(struct PARAM *param)
{
    struct List **listdesc = (struct List **)param->value;
    struct List *item = listdesc[0];
    while (item != NULL)
    {
        struct List *tmp = item;
        item = item->tail;
        free(tmp->data);
        free(tmp);
    }
    free(param->value);
    free(param->key);
    free(param);
}

void FreeItem(struct PARAM *param)
{
    switch (param->type)
    {
    case 0:
        FreeSingleItem(param);
        break;
    case 1:
        FreeListItem(param);
        break;
    }
}

struct PARAM *FindKey(const char *mkey, uint64_t keylen)
{
    if (keylen == 0)
    {
        return NULL;
    }
    uint8_t hash = mkey[0];
    struct PARAM *param = globalparam[hash];
    struct PARAM *beforeparam = NULL;
    while (param != NULL)
    {
        time_t now = time(NULL);
        if (now > param->deadline)
        {
            struct PARAM *p = param;
            param = param->tail;
            uint8_t hash = mkey[0];
            if (beforeparam != NULL)
            {
                beforeparam->tail = param;
            }
            else
            {
                globalparam[hash] = param;
            }
            FreeItem(param);
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

int AddSingleKey(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype)
{
    if (keylen == 0)
    {
        return KEYLENZERO;
    }
    if (datalen == 0)
    {
        return DATANULL;
    }
    uint8_t hash = mkey[0];
    char *key = (char *)malloc(keylen);
    if (key == NULL)
    {
        return MALLOCFAIL;
    }
    void *data = (void *)malloc(datalen);
    if (data == NULL)
    {
        free(key);
        return MALLOCFAIL;
    }
    struct SINGLE *item = (struct SINGLE *)malloc(sizeof(struct SINGLE));
    if (item == NULL)
    {
        free(data);
        free(key);
        return MALLOCFAIL;
    }
    struct PARAM *param = (struct PARAM *)malloc(sizeof(struct PARAM));
    if (param == NULL)
    {
        free(item);
        free(data);
        free(key);
        return MALLOCFAIL;
    }
    memcpy(key, mkey, keylen);
    memcpy(data, mdata, datalen);
    param->key = key;
    param->keylen = keylen;
    item->data = data;
    item->datalen = datalen;
    item->datatype = datatype;
    param->value = item;
    if (ttl == -1)
    {
        param->deadline = -1;
    }
    else
    {
        time_t now = time(NULL);
        param->deadline = now + ttl;
    }
    param->type = 0;
    param->tail = globalparam[hash];
    globalparam[hash] = param;
    return RUNSUCCESS;
}

int SetEx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        uint8_t hash = param->key[0];
        globalparam[hash] = param->tail;
        FreeItem(param);
    }
    return AddSingleKey(mkey, keylen, mdata, datalen, ttl, datatype);
}

int Set(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    return SetEx(mkey, keylen, mdata, datalen, -1, datatype);
}

int SetNx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return WRITEDENY;
    }
    return AddSingleKey(mkey, keylen, mdata, datalen, -1, datatype);
}

int SetNex(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype)
{
    if (FindKey(mkey, keylen) != NULL)
    {
        return WRITEDENY;
    }
    return AddSingleKey(mkey, keylen, mdata, datalen, ttl, datatype);
}

int64_t Get(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 0)
    {
        return TYPEERROR;
    }
    struct SINGLE *item = param->value;
    *datatype = item->datatype;
    int len = item->datalen;
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (dlen > 0 && mdata != NULL)
    {
        memcpy(mdata, item->data, dlen);
        *datalen = dlen;
    }
    return len;
}

int Del(const char *mkey, uint64_t keylen)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param != NULL)
    {
        return DATANULL;
    }
    uint8_t hash = param->key[0];
    globalparam[hash] = param->tail;
    FreeItem(param);
    return RUNSUCCESS;
}

int Incr(const char *mkey, uint64_t keylen)
{
    return IncrBy(mkey, keylen, 1);
}

int IncrBy(const char *mkey, uint64_t keylen, int64_t num)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        int64_t n = num;
        return AddSingleKey(mkey, keylen, &n, sizeof(int64_t), -1, 9);
    }
    struct SINGLE *item = param->value;
    uint64_t datalen = item->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)item->data;
        n += num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)item->data;
        n += num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)item->data;
        n += num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)item->data;
        n += num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)item->data;
        n += num;
        memcpy(item->data, &n, datalen);
    }
    return RUNSUCCESS;
}

int Decr(const char *mkey, uint64_t keylen)
{
    return DecrBy(mkey, keylen, 1);
}

int DecrBy(const char *mkey, uint64_t keylen, int64_t num)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        int64_t n = 0 - num;
        return AddSingleKey(mkey, keylen, &n, sizeof(int64_t), -1, 9);
    }
    struct SINGLE *single = param->value;
    uint64_t datalen = single->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)single->data;
        n -= num;
        memcpy(single->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)single->data;
        n -= num;
        memcpy(single->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)single->data;
        n -= num;
        memcpy(single->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)single->data;
        n -= num;
        memcpy(single->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)single->data;
        n -= num;
        memcpy(single->data, &n, datalen);
    }
    return RUNSUCCESS;
}

int Expire(const char *mkey, uint64_t keylen, int64_t ttl)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (ttl == -1)
    {
        param->deadline = -1;
    }
    else
    {
        time_t now = time(NULL);
        param->deadline = now + ttl;
    }
    return RUNSUCCESS;
}

struct List *CreateList(const void *mdata, uint64_t datalen, int datatype)
{
    char *data = (char *)malloc(datalen);
    if (data == NULL)
    {
        return NULL;
    }
    struct List *item = (struct List *)malloc(sizeof(struct List));
    if (item == NULL)
    {
        return NULL;
    }
    memcpy(data, mdata, datalen);
    item->data = data;
    item->datalen = datalen;
    item->datatype = datatype;
    return item;
}

int AddListKey(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    if (keylen == 0)
    {
        return KEYLENZERO;
    }
    if (datalen == 0)
    {
        return DATANULL;
    }
    char *key = (char *)malloc(keylen);
    if (key == NULL)
    {
        return MALLOCFAIL;
    }
    struct List *item = CreateList(mdata, datalen, datatype);
    if (item == NULL)
    {
        return MALLOCFAIL;
    }
    item->head = NULL;
    item->tail = NULL;
    struct List **listdesc =
        (struct List **)malloc(2 * sizeof(struct List *)); // 创建两个指针，一个指向头部一个指向尾部
    if (listdesc == NULL)
    {
        free(item->data);
        free(item);
        free(key);
        return MALLOCFAIL;
    }
    struct PARAM *param = (struct PARAM *)malloc(sizeof(struct PARAM));
    if (param == NULL)
    {
        free(item->data);
        free(item);
        free(listdesc);
        return MALLOCFAIL;
    }
    memcpy(key, mkey, keylen);
    listdesc[0] = item;
    listdesc[1] = item;
    param->key = key;
    param->keylen = keylen;
    param->value = listdesc;
    param->deadline = -1;
    param->type = 1;
    uint8_t hash = mkey[0];
    param->tail = globalparam[hash];
    globalparam[hash] = param;
    return RUNSUCCESS;
}

int Lpush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return AddListKey(mkey, keylen, mdata, datalen, datatype);
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    struct List *item = CreateList(mdata, datalen, datatype);
    if (item == NULL)
    {
        return MALLOCFAIL;
    }
    item->head = NULL;
    struct List **listdesc = (struct List **)param->value;
    struct List *listhead = listdesc[0];
    item->tail = listhead;
    listhead->head = item;
    listdesc[0] = item;
    return RUNSUCCESS;
}

int64_t Lpop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    struct List **listdesc = (struct List **)param->value;
    struct List *item = listdesc[0];
    *datatype = item->datatype;
    uint64_t len = item->datalen;
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (dlen > 0 && mdata != NULL)
    {
        memcpy(mdata, item->data, dlen);
        *datalen = dlen;
        if (listdesc[0] == listdesc[1]) // 只有一个对象
        {
            uint8_t hash = mkey[0];
            globalparam[hash] = param->tail;
            free(param->value);
            free(param->key);
            free(param);
        }
        else
        {
            listdesc[0] = item->tail;
        }
        free(item->data);
        free(item);
    }
    return len;
}

int Rpush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return AddListKey(mkey, keylen, mdata, datalen, datatype);
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    struct List *item = CreateList(mdata, datalen, datatype);
    if (item == NULL)
    {
        return MALLOCFAIL;
    }
    struct List **listdesc = (struct List **)param->value;
    struct List *listtail = listdesc[1];
    item->head = listtail;
    listtail->tail = item;
    item->tail = NULL;
    listdesc[1] = item;
    return RUNSUCCESS;
}

int64_t Rpop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    struct PARAM *param = FindKey(mkey, keylen);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    struct List **listdesc = (struct List **)param->value;
    struct List *item = listdesc[1];
    *datatype = item->datatype;
    uint64_t len = item->datalen;
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (dlen > 0 && mdata != NULL)
    {
        memcpy(mdata, item->data, dlen);
        *datalen = dlen;
        if (listdesc[0] == listdesc[1]) // 只有一个对象
        {
            uint8_t hash = mkey[0];
            globalparam[hash] = param->tail;
            free(param->value);
            free(param->key);
            free(param);
        }
        else
        {
            listdesc[1] = item->head;
        }
        free(item->data);
        free(item);
    }
    return len;
}
