# 接口文档
> 如果需要本地开发调试，请使用Postman
## 一、Restful 接口标准
为了统一接口，我们采用`Restful`进行描述接口
### 1、协议
API 与用户端通讯使用HTTPS 协议
### 2、域名
> - http://localhost:9999 # 开发环境  
> - https://api.dmicro.cn  # 正式环境     
### 3、版本
将API的版本号放入URL,参考 [Github](https://developer.github.com/v3/media/#request-specific-version)
> https://api.dmicro.com/v1/
### 4、路径
路径又称"终点"（endpoint），表示API的具体网址。

在RESTful架构中，每个网址代表一种资源（resource），所以网址中不能有动词，只能有名词，而且所用的名词往往与数据库的表格名对应。一般来说，数据库中的表都是同种记录的"集合"（collection），所以API中的名词也应该使用复数。

例如：

```
- https://api.dmicro.com/v1/users
- https://api.dmicro.com/v1/menus
- https://api.dmicro.com/v1/roles
```
### 5、HTTP(URL)的动词

```
- GET（SELECT）：从服务器取出资源（一项或多项）。
- POST（CREATE）：在服务器新建一个资源。
- PUT（UPDATE）：在服务器更新资源（客户端提供改变后的完整资源）。
- PATCH（UPDATE）：在服务器更新资源（客户端提供改变的属性）。
- DELETE（DELETE）：从服务器删除资源。
- HEAD：获取资源的元数据。
- OPTIONS：获取信息，关于资源的哪些属性是客户端可以改变的。
```

例子：

```
- GET /users 列出所有用户
- POST /users 新建一个用户
- GET /users/:id 获取某个用户的信息
- PUT /users/:id 更新某个指定用户的信息（提供该用户的全部信息）
- PATCH /users/:id 更新某个指定用户的信息（提供该用户的部分信息）
- DELETE /users/:id 删除某个用户
```
### 6、过滤信息（Filtering）
```
- ?page=2&per_page=100：指定第几页，以及每页的记录数。
- ?sortby=name&order=asc：指定返回结果按照哪个属性排序，以及排序顺序。
- ?users_id=1：指定筛选条件
```
### 7、HTTP状态码（Status Codes）
服务器向用户返回的状态码和提示信息，常见的有以下一些（方括号中是该状态码对应的HTTP动词）。
```
- 200 OK - [GET]：服务器成功返回用户请求的数据，该操作是幂等的（Idempotent）。
- 201 CREATED - [POST/PUT/PATCH]：用户新建或修改数据成功。
- 202 Accepted - [*]：表示一个请求已经进入后台排队（异步任务）
- 204 NO CONTENT - [DELETE]：用户删除数据成功。
- 400 INVALID REQUEST - [POST/PUT/PATCH]：用户发出的请求有错误，服务器没有进行新建或修改数据的操作，该操作是幂等的。
- 401 Unauthorized - [*]：表示用户没有权限（令牌、用户名、密码错误）。
- 403 Forbidden - [*] 表示用户得到授权（与401错误相对），但是访问是被禁止的。
- 404 NOT FOUND - [*]：用户发出的请求针对的是不存在的记录，服务器没有进行操作，该操作是幂等的。
- 406 Not Acceptable - [GET]：用户请求的格式不可得（比如用户请求JSON格式，但是只有XML格式）。
- 410 Gone -[GET]：用户请求的资源被永久删除，且不会再得到的。
- 422 Unprocesable entity - [POST/PUT/PATCH] 当创建一个对象时，发生一个验证错误。
- 500 INTERNAL SERVER ERROR - [*]：服务器发生错误，用户将无法判断发出的请求是否成功。
```

### 8、错误处理（Error handling）

```json
{
    "id": "dmicro",
    "code": 1005,
    "detail": "令牌已过期",
    "status": ""
}
```

### 9、返回结果
针对不同操作，服务器向用户返回的结果应该符合以下规范。
```
- GET /collection：返回资源对象的列表（数组）
- GET /collection/resource：返回单个资源对象
- POST /collection：返回新生成的资源对象
- PUT /collection/resource：返回完整的资源对象
- PATCH /collection/resource：返回完整的资源对象
- DELETE /collection/resource：返回一个空文档
```

## 二、公共请求头
| 参数     | 值        | 类型   | 说明     | 必须   | 
|:--------:|:---------:|:------:|:--------:|:------:|
| Token   | token      | String | 令牌     | true   |

## 三、返回结构体

不同业务返回不同结构

示例
```json
{
  "token_info": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODA4ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.Bzfk4xwVaRfm7rpZXT2fsn4tqWC_bd76GWW2tPT5MhJWQcQAP6bzlpx3t2M7GPqH3A9OWPoHxr2bffrnJjuNDacfdMLu_PO8tg8qdoi4e55kCREvEiKyXB9SxyWYadioiyrs00qMt8VfakN6L9PosgS7xtCFADUkoGBTyHdJXzE",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM2NjU2ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.mW09Ga6PJ5_od1KK0HGK03CpLq_cWrims-9zOUioRdjARAuXfg7kCVnhwJwSl8hqXJ7_rq68TO5buhErEVmnD7wLRIFJA7HGkI39sN9SrqwvTP17lWyCySAIjdQy8vtMml5ZhxCMNlMXGaytTeCkV5vEn_lzVyQsr4gwbRfqWwk",
    "expires_at": 1581080881
  }
}
```

## 四、说明
移动端的接口,可不必遵循Restful规范,统一用POST+JSON规范