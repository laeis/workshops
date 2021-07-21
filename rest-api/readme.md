### **Задача: створити сервіс (RESP API) органайзера-калердаря з функціональністю**

-_додавати події, нагадування._     
-_редагувати їх, змінювати назву, час, опис..._
-_видаляти події_   
-_переглядати перелік подій на день, тиждень, місяць, рік (з пітримкою фільтрації по ознакам)_  
``логика должна быть разбита по слоям, т.е. транспортный слой, сервисный слой, слой(https://github.com/bxcodec/go-clean-arch).
Код нужно покрыть тестами + моки(https://github.com/golang/mock)``

### Part 2 #### 
* Add login/logout/signup enpoints (credantials should be stored in memory as well);
* Use JWT for authorization at other endpoints;
* Update logic to Events endpoints to authorize access only to related to user events.
* PUT for update user timezone; Time of events should be returned in a new timezone.
#### Links
- [golang jwt](https://github.com/golang-jwt/jwt)
- [Hands-on with JWT in Go](https://betterprogramming.pub/hands-on-with-jwt-in-golang-8c986d1bb4c0)
- [writing middleware in Go](https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81)

##### SignUp flow: #####
* parse json
* hash pwd    
* store in db 
* response ok

#####  Login flow: ##### 
* Parse payload 
* Validate if user is in our db and pwd matches 
* generate token
* put in the response

#####  Middleware for login ##### 
* GetTokenFrom header
* Validate using jwt lib
* Parse it
* check if token for this email is inside DB

#####  Logout: ##### 
* Parse payload
* Remove token by user id