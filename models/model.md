## 用户集合
```json 
{
    "identity": "唯一标识",
    "account": "账号",
    "password": "密码",
    "nickname": "昵称",
    "sex": 1,
    "email": "邮箱",
    "avatar": "头像",
    "created_at": "2024-04-20 12:11:11",
    "update_at": "2024-04-20 12:11:11"
}
```

## 消息集合
``` json 
{
    "user_identity":"用户的唯一标识",
    "room_identity":"房间的唯一标识",
    "data":"发送的数据",
    "created_at":"2024-04-20 12:11:11", //创建时间
    "updated_at": "2024-04-20 12:11:11"//更新时间
}
```

## 房间集合
``` json 
{
    "number":"房间号",
    "name":"房间名称",
    "info":"房间简介",
    "user_identity":"房间创建者唯一标识",
    "created_at":"2024-04-20 12:11:11", //创建时间
    "updated_at": "2024-04-20 12:11:11"//更新时间
}
```

## 房间消息关联集合
``` json 
{
    "user_identity":"用户唯一标识",
    "room_identity":"房间唯一标识",
    "message":"消息唯一标识",
    "created_at":"2024-04-20 12:11:11", //创建时间
    "updated_at": "2024-04-20 12:11:11"//更新时间
}
```
