Тест проекта




1.
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzOTM5OTAsImlhdCI6MTc4MjM5MzA5MCwidXNlcl9pZCI6NCwicm9sZSI6InVzZXIifQ.3UCZg1rSUIVzacyrr_Ga8i21hQqSTOtbIXZW431xMg8

2.
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3ODIzOTUxNTcsImlhdCI6MTc4MjM5NDI1NywidXNlcl9pZCI6NCwicm9sZSI6InVzZXIifQ.A92hkiKm3_kcu_8vvaBQVz2Aw-Z9d-vCaWpVzqSXi5o






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


