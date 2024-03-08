user -> html (0 интерактивности)

user -> backend -> html + execute

Single Page Application
```
user -> webAPP -> index.html (содержит js + css + базовая разметка)
    js -> backend -> db

    js after 30s -> backend (например проверяет нет ли новых сообщений)
 
    js after click -> backend (например отправляет письмо по нажатии кнопки)
```
```
Frontend [html + js] -> Backend [server + database]       
        API = Сетевой протокол(HTTP/WS) + Сериализация(url/form-data/xml)
        RPC (Remote Procedure Call)
        Frontend - [
            https://localhost/rpc?func=login_user&login=user1&password=user2
            /rpc?get_all_users
            /rpc?get_all_users_sorted
        ] -> Backend

        REST (REpresentational State Transfer) -- like Amazon S3 REST API
        Frontend - [
            Resource + Action
            Action -- CRUD (Create Read Update Delete)

            GET /user?sort=name&limit=10 -- read user
            POST /user -- create user
            PUT /user -- update user
            DELETE /user -- delete user

            GET /user/101

            POST /image/123/comment + data {"text": "it's awesome"} <- 201 (Created)
            1хх -- service
            2хх -- ok
            3хх -- redirect / resource 
            4хх -- validation error. 400 -- Входные данные, 401 -- не авторизован, 403 -- запрещено, 404 -- такого нет, 
            5хх -- server error. 500 -- internal error, 502 -- bad gateway, 503 -- дальнейший сервис лежит, 504 -- таймаут дальнешего сервиса

            GET /image/123/comment <- 200 [{"id": 1, "text": "it's awesome"}]
            DELETE /image/123/comment/1 <- 204
        ] -> Backend
        Swagger - https://petstore.swagger.io/
        
        GraphQL 
        Frontend - [
            /graphql?query=
            mutation {
                CreateComment(
                    imageID: 123,
                    input: {
                        Text: "it's awesome"
                    }
                )
                CreateComment(
                    imageID: 123,

                )
            }
        ] -> Backend
        ...
```
```proto
schema:

scalar Content;

type User {
    ID: ID!
    Name: String!
    Image: [Image!]!
}

type Image {
    ID: ID!
    Content: Content!
    Comment: [Comment!]
}

type Comment {
    ID: ID!
    Text: String!
    Like: number
}

type Query {
    GetUsers: [User!]!
}

type Mutation {
    CreateComment(imageID: ID!, input CreateCommentInput): Comment
}

input CreateCommentInput {
    Text: String!
}

query {
    GetUsers {
        ID
        Name
        Image {
            Comment {

            }
        }
    }
} 
-> [{"id": 1, "Name": "name1"}, {"id": 2, "Name": "name2"}]
```
