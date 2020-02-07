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
  "data": {
    "token_info": {
      "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwMDQ5MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.ygGAn-8B1z2iVDfBNbSSNqA21-TjMeZBSyb-Hk_MmVZUOtZMllNoddsOuho4pOnQ5RyaqE4dHTdm67LycylSvYbvnhnBHk3qOhBLvcFJCaqgpsklhBf_QcRJTWh8Hy9mRh9AMO6xi9xJncuceOpDmrSahzfSdRp6U_-4rY-oppc",
      "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM1ODk3MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.huN_PYaXSlteC5VPDBuq9sV7Fqw08Fd8DlF1zkZU0CYfULsL6Lt_6lHNmwqvdLeCBtD2ZdI6UZEdmfbljwhXj0jEll5qGmyiHH8C9Yxygceua8-NhfGRAhCW4Nc3vnd50jJcs0xD8fKuApeWunseIvVb7mvFIvMHg_vJiNY6QZU",
      "expired_at": 1581004914
    }
  },
  "errno": 0,
  "t": 1580997714766577200
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
      "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwMDQ5MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.ygGAn-8B1z2iVDfBNbSSNqA21-TjMeZBSyb-Hk_MmVZUOtZMllNoddsOuho4pOnQ5RyaqE4dHTdm67LycylSvYbvnhnBHk3qOhBLvcFJCaqgpsklhBf_QcRJTWh8Hy9mRh9AMO6xi9xJncuceOpDmrSahzfSdRp6U_-4rY-oppc",
      "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM1ODk3MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.huN_PYaXSlteC5VPDBuq9sV7Fqw08Fd8DlF1zkZU0CYfULsL6Lt_6lHNmwqvdLeCBtD2ZdI6UZEdmfbljwhXj0jEll5qGmyiHH8C9Yxygceua8-NhfGRAhCW4Nc3vnd50jJcs0xD8fKuApeWunseIvVb7mvFIvMHg_vJiNY6QZU",
      "expired_at": 1581004914
    }
  },
  "errno": 0,
  "t": 1580997714766577200
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
  "data": {
    "token_info": {
      "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwMDQ5MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.ygGAn-8B1z2iVDfBNbSSNqA21-TjMeZBSyb-Hk_MmVZUOtZMllNoddsOuho4pOnQ5RyaqE4dHTdm67LycylSvYbvnhnBHk3qOhBLvcFJCaqgpsklhBf_QcRJTWh8Hy9mRh9AMO6xi9xJncuceOpDmrSahzfSdRp6U_-4rY-oppc",
      "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM1ODk3MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.huN_PYaXSlteC5VPDBuq9sV7Fqw08Fd8DlF1zkZU0CYfULsL6Lt_6lHNmwqvdLeCBtD2ZdI6UZEdmfbljwhXj0jEll5qGmyiHH8C9Yxygceua8-NhfGRAhCW4Nc3vnd50jJcs0xD8fKuApeWunseIvVb7mvFIvMHg_vJiNY6QZU",
      "expired_at": 1581004914
    }
  },
  "errno": 0,
  "t": 1580997714766577200
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
  "data": {
    "token_info": {
      "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODEwMDQ5MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.ygGAn-8B1z2iVDfBNbSSNqA21-TjMeZBSyb-Hk_MmVZUOtZMllNoddsOuho4pOnQ5RyaqE4dHTdm67LycylSvYbvnhnBHk3qOhBLvcFJCaqgpsklhBf_QcRJTWh8Hy9mRh9AMO6xi9xJncuceOpDmrSahzfSdRp6U_-4rY-oppc",
      "refresh_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODM1ODk3MTQsIkluZm8iOnsiQXBwaWQiOjEsIlVpZCI6MzU2NTI4MDM5MTMzMTg1LCJQbGF0IjoxLCJEZXZpY2VJZCI6IiJ9fQ.huN_PYaXSlteC5VPDBuq9sV7Fqw08Fd8DlF1zkZU0CYfULsL6Lt_6lHNmwqvdLeCBtD2ZdI6UZEdmfbljwhXj0jEll5qGmyiHH8C9Yxygceua8-NhfGRAhCW4Nc3vnd50jJcs0xD8fKuApeWunseIvVb7mvFIvMHg_vJiNY6QZU",
      "expired_at": 1581004914
    }
  },
  "errno": 0,
  "t": 1580997714766577200
}
```