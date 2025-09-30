1)
```
user (alice)
    req vk.com/

server 
    resp html + csrf_token

user (alice)
    req vk.com/rate?id=1&vote=up&csrf_token=abcd

server
    check get_csrf(base, alice) == csrf_token
```


2) double submit 
```
user (alice)
    req vk.com/

server
    resp html + csrf_token + set_cookie(csrf_token)

user (alice)
    req vk.com/rate?id=1&vote=up&csrf_token=abcd
        cookie csrf_token

server
    check csrf_token_cookie == csrf_token_param
```