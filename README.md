Тест проекта




1.
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzOTM5OTAsImlhdCI6MTc4MjM5MzA5MCwidXNlcl9pZCI6NCwicm9sZSI6InVzZXIifQ.3UCZg1rSUIVzacyrr_Ga8i21hQqSTOtbIXZW431xMg8

2.
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzOTUxNTcsImlhdCI6MTc4MjM5NDI1NywidXNlcl9pZCI6NCwicm9sZSI6InVzZXIifQ.A92hkiKm3_kcu_8vvaBQVz2Aw-Z9d-vCaWpVzqSXi5o

3.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI0NzAxMjUsImlhdCI6MTc4MjQ2OTIyNSwidXNlcl9pZCI6NCwicm9sZSI6InVzZXIifQ.VP79zaqv8ChX1WuiIcZ97d_qGYSgaqxccTUXlTN4fYI

4.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI0NzA4OTIsImlhdCI6MTc4MjQ2OTk5MiwidXNlcl9pZCI6Niwicm9sZSI6InVzZXIifQ.xamLcCraCetsvmK0szkDAJ30z9tqaYHlj0wrabTCCis

5.eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODI0NzExNjQsImlhdCI6MTc4MjQ3MDI2NCwidXNlcl9pZCI6Niwicm9sZSI6InVzZXIifQ.wmUmGa9lnZtBdUC5KYFZTZcwl2xqDRq2kqmH0vAsH-4





http://localhost:8080/auth/register:

1. 
{
  "name": "Suhrob",
  "email": "suhrob.test@gmail.com",
  "password": "password1234",
  "phone": "992900000000"
}  --> {
    "user_id": 1
}  , 201 created

2. 
 {
  "name": "Suhrob",
  "email": "",
  "password": "password1234",
  "phone": "992900000000"
}  --> invalid email format 400 bad req

3. 
{
  "name": "Suhrob",
  "email": "test2@gmail.com",
  "password": "123",
  "phone": "992900000000"
} --> password must be at least 6 characters long  400 bad req

4.
{
  "name": "Suhrob",
  "email": "suhrob.test@gmail.com",
  "password": "password1234",
  "phone": "992900000000"
} --> email already exists 409 conflict

http://localhost:8080/auth/login:

1. 
{
  "email": "suhrob.test@gmail.com",
  "password": "password1234"
} -->   {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzODc2MjgsImlhdCI6MTc4MjM4NjcyOCwidXNlcl9pZCI6MSwicm9sZSI6InVzZXIifQ.FF1hDXZ3ziOVwby1Haiu06U-XjGKhWOuy_LfpC42wPA",
    "refresh_token": "d660645cb03b3123034fd724b3f045645f7d0ab1bfda6680d58b432e24e334de",
    "expires_at": "2026-06-25T16:40:28.107538716+05:00"
}   200, ok

2. 

{
  "email": "suhrob.test@gmail.com",
  "password": "wrongpassword"
} --> {
    "error": "invalid email or password"
} 401 , unauthorizied

3. 

{
  "email": "wrong@gmail.com",
  "password": "password1234"
} --> user not founf 404 not found



http://localhost:8080/auth/refresh:

1. 
{
  "refresh_token": "d660645cb03b3123034fd724b3f045645f7d0ab1bfda6680d58b432e24e334de"
} -->  {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzODgxNjEsImlhdCI6MTc4MjM4NzI2MSwidXNlcl9pZCI6MSwicm9sZSI6InVzZXIifQ.j-SL3dpBRVyud0dS9XFydLw7LL5_ln5v0MEvO-Or2kE",
    "refresh_token": "e79842149acd94bf79e9098e80db3a4fc1d68401851ddf06e93fe8af0bc70185",
    "expires_at": "2026-06-25T16:49:21.954294451+05:00"
}  200 ok


http://localhost:8080/users/me  GET:

1. 
без access token ---> invalid authorization header format  , 401 unauth

c неправильным acc token ---> invalid or expired token , 401 unauth

c access token ---> Показывает мой профиль



http://localhost:8080/users/me PUT: 

1. 
{
    "name":"Suhrob",
    "email":"suhrobusmonalizoda0625@gmail.com",
    "password":"admin12345",
    "phone":"992071055225"
}  
c токеном --> {
    "message": "profile updated successfully"
}  200 , OK

 без access token ---> invalid authorization header format , 401 unauth
 c неправильным токен ---> invalid or expired token, 401 unauth

http://localhost:8080/users/me DELETE: 
1. 
c токеном --> {
    "message": "account deleted successfully"
}  

с неправильным токеном --> invalid or expired token 401 unauth

http://localhost:8080/orders POST
1. 
{
    "product":"IPhone 15",
    "price":3000
} ---> {
    "order_id": 2
} ---> 201 Created

2. 
{
    "product":"Mac",
    "price":-9000
} ---> {
    "error": "order price must be greater than zero"
} ---> 400 bad req

3. 
без токена ---> authorization header is required

с неправильным токеном ---> invalid or expired token

http://localhost:8080/orders GET 

1. 
без токена ---> authorization header is required
c неправильнвм токеном ---> invalid or expired token

c праивильным токеном ----> [
    {
        "id": 2,
        "product": "IPhone 15",
        "price": 3000,
        "user_id": 4,
        "status": "new",
        "created_at": "2026-06-26T15:21:44.416909+05:00"
    },
    {
        "id": 1,
        "product": "Go Courses",
        "price": 150.5,
        "user_id": 4,
        "status": "new",
        "created_at": "2026-06-25T20:00:01.560241+05:00"
    }
]

200, ok

2. 
Заходим из другого пользователя у которого нет закозов
вывод ---> []

200 , ok





http://localhost:8080/orders/1 GET

1. 
Вставляем свой токен но выбираем айди чужого заказа 
---->   {
    "error": "access denied"
} ---> 403 Forbidden

2. 


http://localhost:8080/users/change-password :

1. 
{
    "old_password":"admin12345",
    "new_password":"admin25"
} ---> {
    "message": "password changed successfully"
} ---> 200, ok

2. 
{
    "old_password":"admin9999999999",
    "new_password":"admin25"
} ---> {
    "error": "old password is incorrect"
}  ----> 401 unauth

http://localhost:8080/admin/users:

1. 
Если вставить токен простого юзера ----> forbidden: admin access required 403

2. 
с токеном админа ---> показывает юзеров 200, ok

3. 
инвалидный токен или без токена ---> invalid or expired token

http://localhost:8080/admin/orders:

1. 
с токеном админа ---> показывает все заказы 200, ok

2. 
Если вставить токен простого юзера ----> forbidden: admin access required 403

3. 
инвалидный токен или без токена ---> invalid or expired token  401 unauth
