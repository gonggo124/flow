#ifndef DA_H
#define DA_H

#define da_append(xs, x) do {                                           \
                if (xs.size >= xs.cap) {                                \
                        if (xs.cap == 0) xs.cap=64;                     \
                        else xs.cap*=2;                                 \
                        xs.data=realloc(xs.data,xs.cap*sizeof(*xs.data)); \
                }                                                       \
                xs.data[xs.size++]=x;                                   \
        } while(0)


#define da_clear(xs) do {                               \
                xs.size = 0;                            \
        } while(0)

#define da_destroy(xs) do {                     \
                free(xs.data);                  \
                xs.size = 0;                    \
                xs.cap = 0;                     \
                xs.data = NULL;                 \
        } while(0)

#endif // DA_H
