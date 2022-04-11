#include "mnsql.h"
#include <string.h>
#include <time.h>

struct SINGLE
{
    int datatype;
    void *data;
    uint64_t datalen;
};

struct List
{
    void *data;
    uint64_t datalen;
    int datatype;
    struct List *head;
    struct List *tail;
};

struct Map
{
    char *key;
    uint64_t keylen;
    void *data;
    uint64_t datalen;
    int datatype;
    uint64_t deadline;
    struct Map *tail;
};

struct PARAM
{
    char *key;
    uint64_t keylen;
    int8_t type;
    void *value;
    uint64_t deadline;
    struct PARAM *tail;
};
struct PARAM *globalparam[256];

void FreeItem(struct PARAM *param)
{
    int8_t type = param->type;
    if (type == 0)
    {
        struct SINGLE *item = param->value;
        free(item->data);
        free(param->value);
        free(param->key);
        free(param);
    }
    else if (type == 1)
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
    else if (type == 2)
    {
        struct Map **mapdesc = (struct Map **)param->value;
        for (int i = 0; i < 256; i++)
        {
            struct Map *item = mapdesc[i];
            while (item != NULL)
            {
                struct Map *tmp = item;
                item = item->tail;
                free(tmp->key);
                free(tmp->data);
                free(tmp);
            }
        }
        free(param->value);
        free(param->key);
        free(param);
    }
}

struct PARAM *FindKey(const char *mkey, uint64_t keylen, time_t now)
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
            FreeItem(p);
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

int AddSingleKey(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype,
                 time_t now)
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
        param->deadline = now + ttl;
    }
    param->type = 0;
    param->tail = globalparam[hash];
    globalparam[hash] = param;
    return RUNSUCCESS;
}

int _Set(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype, int nx)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param != NULL)
    {
        if (nx)
        {
            return WRITEDENY;
        }
        uint8_t hash = param->key[0];
        globalparam[hash] = param->tail;
        FreeItem(param);
    }
    return AddSingleKey(mkey, keylen, mdata, datalen, ttl, datatype, now);
}

int Set(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    return _Set(mkey, keylen, mdata, datalen, -1, datatype, 0);
}

int SetEx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype)
{
    return _Set(mkey, keylen, mdata, datalen, ttl, datatype, 0);
}

int SetNx(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    return _Set(mkey, keylen, mdata, datalen, -1, datatype, 1);
}

int SetNex(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl, int datatype)
{
    return _Set(mkey, keylen, mdata, datalen, ttl, datatype, 1);
}

int64_t Get(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
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
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (dlen == 0 || mdata == NULL)
    {
        return item->datalen;
    }
    memcpy(mdata, item->data, dlen);
    *datalen = dlen;
    return item->datalen;
}

int Del(const char *mkey, uint64_t keylen)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
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
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        int64_t n = num;
        return AddSingleKey(mkey, keylen, &n, sizeof(int64_t), -1, 9, now);
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
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        int64_t n = 0 - num;
        return AddSingleKey(mkey, keylen, &n, sizeof(int64_t), -1, 9, now);
    }
    struct SINGLE *item = param->value;
    uint64_t datalen = item->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    return RUNSUCCESS;
}

int Expire(const char *mkey, uint64_t keylen, int64_t ttl)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
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
        param->deadline = now + ttl;
    }
    return RUNSUCCESS;
}

int _Push(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype, int direct)
{
    int8_t isnewparam = 0;
    uint8_t hash = mkey[0];
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        if (keylen == 0)
        {
            return KEYLENZERO;
        }
        char *key = (char *)malloc(keylen);
        if (key == NULL)
        {
            return MALLOCFAIL;
        }
        struct List **listdesc =
            (struct List **)malloc(2 * sizeof(struct List *)); // 创建两个指针，一个指向头部一个指向尾部
        if (listdesc == NULL)
        {
            free(key);
            return MALLOCFAIL;
        }
        param = (struct PARAM *)malloc(sizeof(struct PARAM));
        if (param == NULL)
        {
            free(listdesc);
            free(key);
            return MALLOCFAIL;
        }
        memcpy(key, mkey, keylen);
        listdesc[0] = NULL;
        listdesc[1] = NULL;
        param->key = key;
        param->keylen = keylen;
        param->value = listdesc;
        param->deadline = -1;
        param->type = 1;
        param->tail = globalparam[hash];
        globalparam[hash] = param;
        isnewparam = 1;
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    void *data = (void *)malloc(datalen);
    if (data == NULL)
    {
        if (isnewparam)
        {
            globalparam[hash] = param->tail;
            free(param->key);
            free(param->value);
            free(param);
        }
        return MALLOCFAIL;
    }
    struct List *item = (struct List *)malloc(sizeof(struct List));
    if (item == NULL)
    {
        free(data);
        if (isnewparam)
        {
            globalparam[hash] = param->tail;
            free(param->key);
            free(param->value);
            free(param);
        }
        return MALLOCFAIL;
    }
    memcpy(data, mdata, datalen);
    item->data = data;
    item->datalen = datalen;
    item->datatype = datatype;
    struct List **listdesc = (struct List **)param->value;
    if (direct == 0)
    {
        struct List *listhead = listdesc[0];
        item->tail = listhead;
        if (listhead != NULL)
        {
            listhead->head = item;
        }
        item->head = NULL;
        listdesc[0] = item;
        if (listdesc[1] == NULL)
        {
            listdesc[1] = item;
        }
    }
    else
    {
        struct List *listtail = listdesc[1];
        item->head = listtail;
        if (listtail != NULL)
        {
            listtail->tail = item;
        }
        item->tail = NULL;
        listdesc[1] = item;
        if (listdesc[0] == NULL)
        {
            listdesc[0] = item;
        }
    }
    return RUNSUCCESS;
}

int LPush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    return _Push(mkey, keylen, mdata, datalen, datatype, 0);
}

int RPush(const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int datatype)
{
    return _Push(mkey, keylen, mdata, datalen, datatype, 1);
}

int64_t _Pop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype, int direct)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 1)
    {
        return TYPEERROR;
    }
    struct List **listdesc = (struct List **)param->value;
    struct List *item = direct == 0 ? listdesc[0] : listdesc[1];
    *datatype = item->datatype;
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (dlen == 0 || mdata == NULL)
    {
        return item->datalen;
    }
    *datalen = dlen;
    memcpy(mdata, item->data, dlen);
    if (listdesc[0] == listdesc[1]) // 只有一个对象
    {
        uint8_t hash = mkey[0];
        globalparam[hash] = param->tail;
        free(param->value);
        free(param->key);
        free(param);
    }
    else if (direct == 0)
    {
        listdesc[0] = item->tail;
    }
    else
    {
        listdesc[1] = item->head;
    }
    free(item->data);
    uint64_t len = item->datalen;
    free(item);
    return len;
}

int64_t LPop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    return _Pop(mkey, keylen, mdata, datalen, datatype, 0);
}

int64_t RPop(const char *mkey, uint64_t keylen, void *mdata, uint64_t *datalen, int *datatype)
{
    return _Pop(mkey, keylen, mdata, datalen, datatype, 1);
}

struct Map *FindMap(struct PARAM *param, const char *mkey, uint64_t keylen, time_t now)
{
    if (keylen == 0)
    {
        return NULL;
    }
    struct Map **mapdesc = param->value;
    uint8_t hash = mkey[0];
    struct Map *item = mapdesc[hash];
    struct Map *beforeitem = NULL;
    while (item != NULL)
    {
        if (now > item->deadline)
        {
            struct Map *p = item;
            item = item->tail;
            if (beforeitem != NULL)
            {
                beforeitem->tail = item;
            }
            else
            {
                mapdesc[hash] = item;
            }
            free(p->data);
            free(p->key);
            free(p);
        }
        else if (item->keylen == keylen && !memcmp(item->key, mkey, keylen))
        {
            if (beforeitem != NULL)
            {
                beforeitem->tail = item->tail;
                item->tail = mapdesc[hash];
                mapdesc[hash] = item;
            }
            return item;
        }
        else
        {
            beforeitem = item;
            item = item->tail;
        }
    }
    return NULL;
}

int AddMapItem(struct PARAM *param, const char *mkey, uint64_t keylen, const void *mdata, uint64_t datalen, int64_t ttl,
               int datatype, time_t now)
{
    struct Map **mapdesc = (struct Map **)param->value;
    char *key = (void *)malloc(keylen);
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
    struct Map *item = (struct Map *)malloc(sizeof(struct Map));
    if (item == NULL)
    {
        free(data);
        free(key);
        return MALLOCFAIL;
    }
    memcpy(key, mkey, keylen);
    memcpy(data, mdata, datalen);
    item->data = data;
    item->datalen = datalen;
    item->key = key;
    item->keylen = keylen;
    item->datatype = datatype;
    if (ttl == -1)
    {
        item->deadline = -1;
    }
    else
    {
        item->deadline = now + ttl;
    }
    uint8_t hash = mkey[0];
    item->tail = mapdesc[hash];
    mapdesc[hash] = item;
    return RUNSUCCESS;
}

int _HSet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, uint64_t datalen,
          int64_t ttl, int datatype, int nx)
{
    int8_t isnewparam = 0;
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        if (keylen == 0)
        {
            return KEYLENZERO;
        }
        char *key = (void *)malloc(keylen);
        if (key == NULL)
        {
            return MALLOCFAIL;
        }
        struct Map **mapdesc = (struct Map **)malloc(256 * sizeof(struct Map *));
        if (mapdesc == NULL)
        {
            free(key);
            return MALLOCFAIL;
        }
        param = (struct PARAM *)malloc(sizeof(struct PARAM));
        if (param == NULL)
        {
            free(mapdesc);
            free(key);
            return MALLOCFAIL;
        }
        memcpy(key, mkey, keylen);
        memset(mapdesc, 0, 256);
        param->key = key;
        param->keylen = keylen;
        param->value = mapdesc;
        param->type = 2;
        param->deadline = -1;
        uint8_t hash = mkey[0];
        param->tail = globalparam[hash];
        globalparam[hash] = param;
        isnewparam = 1;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item != NULL)
    {
        if (nx)
        {
            return WRITEDENY;
        }
        uint8_t hash = mkey2[0];
        struct Map **mapdesc = (struct Map **)param->value;
        mapdesc[hash] = item->tail;
        free(item->key);
        free(item->data);
        free(item);
    }
    int res = AddMapItem(param, mkey2, keylen2, mdata, datalen, ttl, datatype, now);
    if (res < 0)
    {
        if (isnewparam)
        {
            uint8_t hash = mkey[0];
            globalparam[hash] = param->tail;
            free(param->key);
            free(param->value);
            free(param);
        }
    }
    return res;
}

int HSet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, uint64_t datalen,
         int datatype)
{
    return _HSet(mkey, keylen, mkey2, keylen2, mdata, datalen, -1, datatype, 0);
}

int HSetEx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, uint64_t datalen,
           int64_t ttl, int datatype)
{
    return _HSet(mkey, keylen, mkey2, keylen2, mdata, datalen, ttl, datatype, 0);
}

int HSetNx(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, uint64_t datalen,
           int datatype)
{
    return _HSet(mkey, keylen, mkey2, keylen2, mdata, datalen, -1, datatype, 1);
}

int HSetNex(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, const void *mdata, uint64_t datalen,
            int64_t ttl, int datatype)
{
    return _HSet(mkey, keylen, mkey2, keylen2, mdata, datalen, ttl, datatype, 1);
}

int64_t HGet(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, void *mdata, uint64_t *datalen,
             int *datatype)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item == NULL)
    {
        if (fp == NULL)
            fp = fopen("logfile.log", "wb");
        fprintf(fp, "in %s,at %d\n", __FILE__, __LINE__);
        fflush(fp);
        return DATANULL;
    }
    *datatype = item->datatype;
    uint64_t dlen = item->datalen > *datalen ? *datalen : item->datalen;
    if (datalen == 0 || mdata == NULL)
    {
        return item->datalen;
    }
    memcpy(mdata, item->data, dlen);
    *datalen = dlen;
    return item->datalen;
}

int HDel(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item == NULL)
    {
        return DATANULL;
    }
    uint8_t hash = mkey2[0];
    struct Map **mapdesc = (struct Map **)param->value;
    mapdesc[hash] = item->tail;
    free(item->key);
    free(item->data);
    free(item);
    return RUNSUCCESS;
}

int HIncr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2)
{
    return HIncrBy(mkey, keylen, mkey2, keylen2, 1);
}

int HIncrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item == NULL)
    {
        int64_t n = num;
        return AddMapItem(param, mkey2, keylen2, &n, sizeof(int64_t), -1, 9, now);
    }
    uint64_t datalen = item->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    return RUNSUCCESS;
}

int HDecr(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2)
{
    return HDecrBy(mkey, keylen, mkey2, keylen2, 1);
}

int HDecrBy(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t num)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item == NULL)
    {
        int64_t n = 0 - num;
        return AddMapItem(param, mkey2, keylen2, &n, sizeof(int64_t), -1, 9, now);
    }
    uint64_t datalen = item->datalen;
    if (datalen == sizeof(char))
    {
        char n = *(char *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(short))
    {
        short n = *(short *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(int))
    {
        int n = *(int *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long))
    {
        long n = *(long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    else if (datalen == sizeof(long long))
    {
        long long n = *(long long *)item->data;
        n -= num;
        memcpy(item->data, &n, datalen);
    }
    return RUNSUCCESS;
}

int HExpire(const char *mkey, uint64_t keylen, const char *mkey2, uint64_t keylen2, int64_t ttl)
{
    time_t now = time(NULL);
    struct PARAM *param = FindKey(mkey, keylen, now);
    if (param == NULL)
    {
        return DATANULL;
    }
    if (param->type != 2)
    {
        return TYPEERROR;
    }
    struct Map *item = FindMap(param, mkey2, keylen2, now);
    if (item == NULL)
    {
        return DATANULL;
    }
    if (ttl == -1)
    {
        item->deadline = -1;
    }
    else
    {
        item->deadline = now + ttl;
    }
    return RUNSUCCESS;
}
