typedef struct Cookie {
        const char* name;
        const char* value;
        const char* domain;
        const char* path;
        long expires;
        int httpOnly;
        int secure;
} Cookie;


Cookie* getCookies(char* url, int* len);
void setCookies(char* url, Cookie* cookies, int len);
