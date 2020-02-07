# passport模块

## 发送短信验证码

接口描述: 发送短信验证码

请求路径:  /passport/Sms

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| mobile   | 13805918888         | String | 手机号        | true   |

示例：
```json
{
  "mobile": "13805918888"
}
```

响应数据: 
```json
{
  "code": "8888"
}
```

## 短信验证码登录

接口描述: 短信验证码登录

请求路径:  /passport/SmsLogin

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| appid    | 1                   | Int    | appid         | true   |
| plat     | 1                   | Int    | 平台类型      | true   |
| mobile   | 13805918888         | String | 手机号        | true   |
| code     | 5188                | String | 验证码        | true   |

示例：
```json
{
  "appid": 1,
  "plat": 1,
  "mobile": "13805918888",
  "code": "5188"
}
```

响应数据: 
```json
{
  "token_info": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODA4ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.Bzfk4xwVaRfm7rpZXT2fsn4tqWC_bd76GWW2tPT5MhJWQcQAP6bzlpx3t2M7GPqH3A9OWPoHxr2bffrnJjuNDacfdMLu_PO8tg8qdoi4e55kCREvEiKyXB9SxyWYadioiyrs00qMt8VfakN6L9PosgS7xtCFADUkoGBTyHdJXzE",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM2NjU2ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.mW09Ga6PJ5_od1KK0HGK03CpLq_cWrims-9zOUioRdjARAuXfg7kCVnhwJwSl8hqXJ7_rq68TO5buhErEVmnD7wLRIFJA7HGkI39sN9SrqwvTP17lWyCySAIjdQy8vtMml5ZhxCMNlMXGaytTeCkV5vEn_lzVyQsr4gwbRfqWwk",
    "expires_at": 1581080881
  }
}
```

## 设置密码

接口描述: 设置密码

请求路径:  /passport/SetPwd

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| passwd   | 123456              | String | 手机号        | true   |

示例：
```json
{
  "passwd": "123456"
}
```

响应数据: 
```json
{
  "token_info": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODA4ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.Bzfk4xwVaRfm7rpZXT2fsn4tqWC_bd76GWW2tPT5MhJWQcQAP6bzlpx3t2M7GPqH3A9OWPoHxr2bffrnJjuNDacfdMLu_PO8tg8qdoi4e55kCREvEiKyXB9SxyWYadioiyrs00qMt8VfakN6L9PosgS7xtCFADUkoGBTyHdJXzE",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM2NjU2ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.mW09Ga6PJ5_od1KK0HGK03CpLq_cWrims-9zOUioRdjARAuXfg7kCVnhwJwSl8hqXJ7_rq68TO5buhErEVmnD7wLRIFJA7HGkI39sN9SrqwvTP17lWyCySAIjdQy8vtMml5ZhxCMNlMXGaytTeCkV5vEn_lzVyQsr4gwbRfqWwk",
    "expires_at": 1581080881
  }
}
```

## 密码登录

接口描述: 密码登录

请求路径:  /passport/Login

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| appid    | 1                   | Int    | appid         | true   |
| plat     | 1                   | Int    | 平台类型      | true   |
| mobile   | 13805918888         | String | 手机号        | true   |
| passwd   | 123456              | String | 密码          | true   |

示例：
```json
{
  "appid": 1,
  "plat": 1,
  "passwd": "123456"
}
```

响应数据: 
```json
{
  "token_info": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODA4ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.Bzfk4xwVaRfm7rpZXT2fsn4tqWC_bd76GWW2tPT5MhJWQcQAP6bzlpx3t2M7GPqH3A9OWPoHxr2bffrnJjuNDacfdMLu_PO8tg8qdoi4e55kCREvEiKyXB9SxyWYadioiyrs00qMt8VfakN6L9PosgS7xtCFADUkoGBTyHdJXzE",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM2NjU2ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.mW09Ga6PJ5_od1KK0HGK03CpLq_cWrims-9zOUioRdjARAuXfg7kCVnhwJwSl8hqXJ7_rq68TO5buhErEVmnD7wLRIFJA7HGkI39sN9SrqwvTP17lWyCySAIjdQy8vtMml5ZhxCMNlMXGaytTeCkV5vEn_lzVyQsr4gwbRfqWwk",
    "expires_at": 1581080881
  }
}
```

## 第三方帐号登录

接口描述: 第三方帐号登录

请求路径:  /passport/OAuthLogin

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| appid    | 1                   | Int    | appid         | true   |
| plat     | 1                   | Int    | 平台类型      | true   |
| platform | wechat              | String | 第三方平台    | true   |
| code     | xdbNkmeqWn          | String | 授权码        | true   |

示例：
```json
{
  "appid": 1,
  "plat": 1,
  "platform": "weixin",
  "code": "xdbNkmeqWn"
}
```

响应数据: 
```json
{
  "token_info": {
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwODA4ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.Bzfk4xwVaRfm7rpZXT2fsn4tqWC_bd76GWW2tPT5MhJWQcQAP6bzlpx3t2M7GPqH3A9OWPoHxr2bffrnJjuNDacfdMLu_PO8tg8qdoi4e55kCREvEiKyXB9SxyWYadioiyrs00qMt8VfakN6L9PosgS7xtCFADUkoGBTyHdJXzE",
    "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM2NjU2ODEsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.mW09Ga6PJ5_od1KK0HGK03CpLq_cWrims-9zOUioRdjARAuXfg7kCVnhwJwSl8hqXJ7_rq68TO5buhErEVmnD7wLRIFJA7HGkI39sN9SrqwvTP17lWyCySAIjdQy8vtMml5ZhxCMNlMXGaytTeCkV5vEn_lzVyQsr4gwbRfqWwk",
    "expires_at": 1581080881
  }
}
```
