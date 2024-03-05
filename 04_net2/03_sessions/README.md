```
user -> login (username + pass) -> backend (username + pass valid? -> generate session + put in db) -> DB
                    <- 200, Set-Cookie: session_id=172348716238476128374agsdjfgasj
    
     -> POST /image -> backend ( request.Cookie.Get(session_id) ) -> DB (is session_id valid)
                    Cookie: session_id=172348716238476128374agsdjfgasj

                    <- 201 {"id": 1}

     -> GET /image ->
                    Cookie: session_id=172348716238476128374agsdjfgasj

hacker -> POST /image -> backend
                    Cookie: session_id=172348716238476128374agsdjfgasj
                    <- 201 

user -> logout -> backend -> DELETE (db)
                    <- Set-Cookie: session_id=172348716238476128374agsdjfgasj;expires=time.Now()-1
     -> GET /image -> backend
                    Cookie: ""
                    <- 401 

hacker -> POST /image -> backend
                    Cookie: session_id=172348716238476128374agsdjfgasj
                    <- 401 


http packet:
    Method URI
    Headers (Cookie, Authorization)
    Body

user -> login (username + pass) -> backend (username + pass valid? -> generate token -> sign token) -> DB(access token123123123123?)
                    <- 200, {"token": "token123123123123"}
    
     -> POST /image -> backend (token valid?)
                    (HEADER) Authorization: Bearer token123123123123

                    <- 201 {"id": 2}
     -> logout -> DB(token123123123123 -- invalid)
                    <- 401

hacker -> POST /image -> backend
                    (HEADER) Authorization: Bearer token123123123123
                    <- 201
```                