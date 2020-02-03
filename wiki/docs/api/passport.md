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
  "data": {
    "code": "8888"
  },
  "errno": 0,
  "t": 1578466107653492700
}
```

## 短信验证码登录

接口描述: 短信验证码登录

请求路径:  /passport/SmsLogin

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| mobile   | 13805918888         | String | 手机号        | true   |
| code     | 5188                | String | 验证码        | true   |

示例：
```json
{
  "mobile": "13805918888",
  "code": "5188"
}
```

响应数据: 
```json
{
  "data": {
    "token_info": {
      "uid": 354248604778497,
      "token": "954b3779-43cb-4777-9209-900ee78144dc",
      "refresh_token": "80be6212-d964-4b02-9280-30c0130084e8",
      "expired_at": 1587108523
    }
  },
  "errno": 0,
  "t": 1578468523883044600
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
  "data": {
    "token_info": {
      "uid": 354248604778497,
      "token": "a6bb934f-3ca8-4170-bfba-f537de2e55da",
      "refresh_token": "575b83d1-78e1-4c87-91f5-d1567ee82207",
      "expired_at": 1587109011
    }
  },
  "errno": 0,
  "t": 1578469011551142000
}
```

## 密码登录

接口描述: 密码登录

请求路径:  /passport/Login

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| mobile   | 13805918888         | String | 手机号        | true   |
| passwd   | 123456              | String | 密码          | true   |

示例：
```json
{
  "passwd": "123456"
}
```

响应数据: 
```json
{
  "data": {
    "token_info": {
      "uid": 354248604778497,
      "token": "18b17014-748b-47d8-a1ff-8bd9d8bd6c8c",
      "refresh_token": "c9ccdd34-030f-4ead-9656-8fe2fb471996",
      "expired_at": 1587109179
    }
  },
  "errno": 0,
  "t": 1578469179132453000
}
```

## 第三方帐号登录

接口描述: 第三方帐号登录

请求路径:  /passport/OAuthLogin

请求方法:  POST

请求参数 : 参数类型 body
| 参数     | 值                  | 类型   | 说明          | 必须   | 
|:---------|:--------------------|:-------|:--------------|:-------|
| platform | wechat              | String | 第三方平台    | true   |
| code     | xdbNkmeqWn          | String | 授权码        | true   |

示例：
```json
{
  "platform": "weixin",
  "code": "xdbNkmeqWn"
}
```

响应数据: 
```json
{
  "data": {
    "token_info": {
      "uid": 354248604778497,
      "token": "18b17014-748b-47d8-a1ff-8bd9d8bd6c8c",
      "refresh_token": "c9ccdd34-030f-4ead-9656-8fe2fb471996",
      "expired_at": 1587109179
    }
  },
  "errno": 0,
  "t": 1578469179132453000
}
```