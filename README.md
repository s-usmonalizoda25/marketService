Тест проекта 
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

