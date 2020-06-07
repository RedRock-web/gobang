## 引言

此项目用于红岩后端考核，基本功能是简单的双人五子棋

## 功能

- [x] 用户登陆注册
- [x] 加入房间
- [x] 玩家准备
- [x] 游戏对战
- [x] 房间密码
- [x] 聊天
- [x] 观众观战
- [ ] 落子限时
- [x] 悔棋
- [x] 求和
- [x] 认输
- [ ] 棋局回放
- [ ] 禁手
- [ ] 残局挑战

。。。

## 运行环境

5.6.15-arch1-1，**mar**iadb 10.4.13

这是部署后的 [Demo](http://47.98.57.152:8080/)

## 技术框架

gorm + gin

## 实现思路

五子棋和实时聊天使用 websocket 协议比较方便，但是我并不是很熟悉，所以用的 http, 要达到实时的目的，可以前端轮询。

这个项目是围绕 room 这个 struct 展开的，把一些基本的属性封装到这个 struct 后，如封装一个 password, 之后设置房间密码就围绕这个展开，或如封装一个 board, 这个是房间内棋盘的抽象，然后棋盘相关的，如下棋，判断胜利等功能，又围绕这个展开。大体就是这个思路，然后具体的功能就只是一些业务逻辑的实现，由于时间和实力限制，并没有使用一些比较高级的技巧。下面是一些具体功能的思路:



用户加入房间：创建房间类型，设置一张房间表，用户的加入，退出直接操作这张表

玩家和观众：给房间类型加上这两种属性，再对加入的用户进行数量判断

房间准备：给房间类型加上准备属性，再操作这个属性

游戏对战：给房间加上棋盘属性，再给棋盘加上下棋，悔棋，判断输赢等属性，然后操作棋盘

房间密码：给房间加上密码属性，然后操作这个属性

观战功能：先判断玩家和观众，然后直接分开给两者响应

聊天功能：给房间加上消息功能后，再分别对房间内和游戏内加上聊天功能，把消息请求存于数据库，再从数据库中直接取出所有消息并响应

求和/认输功能：给房间加上两个 flag，然后在游戏中，如果玩家发出请求，直接修改这个 flag，当 flag 满了就直接退出游戏回到房间

棋局回放：

悔棋：

## 项目目录

```
.
├── app
│   ├── gobang
│   └── user
├── cmd
├── configs
├── db
├── jwts
├── logs
├── middleware
├── README.md
├── response
└── router
```

## 接口文档

默认 8080 端口

### 列表

| URL            | http | 功能                     |
| -------------- | ---- | ------------------------ |
| /              | get  | 测试是否部署成功         |
| /user/register | post | 用户注册                 |
| /user/login    | post | 用户登录                 |
| /user/get      | post | 获取用户信息             |
| /user/modify   | post | 修改用户信息             |
| /room/create   | post | 创建房间                 |
| /room/join     | post | 加入房间，成为玩家或观众 |
| /room/password |      |                          |
| /room/exit     |      |                          |
| /room/close    |      |                          |
| /room/chat     |      |                          |
| /game/start    | post | 开始游戏                 |
| /game/play     | post | 玩家下棋/观众观看        |
| /game/peace    |      |                          |
| /game/confess  |      |                          |
| /game/regret   |      |                          |
| /game/history  |      |                          |
| /game/chat     |      |                          |

### 调用流程

用户一共有三个存在状态，普通状态，房间状态，游戏状态。

* 普通状态

用户需要注册登录，这时用户可以修改和查看信息。

* 房间状态

用户可以创建房间。

房间主可以给房间设置密码。

用户可以指定房间号和房间密码加入房间，没有满员则是玩家，满员则是观众。

在房间中用户可以聊天，退出房间。

房间主可以关闭房间。

在房间中可以查看上一次游戏的历史回放。

玩家分为准备和未准备状态。

两个玩家都准备后，任何一个玩家都开始游戏。

* 游戏状态

玩家发送棋子坐标下棋。

玩家轮流下棋。

其中一个玩家胜利后，游戏结束，回到房间。

游戏过程中观众可以观看。

游戏过程中观众可以退出游戏，回到房间中。

游戏过程中用户可以聊天。

游戏过程中玩家可以提出求和/认输/后悔机制，另一个玩家同意后，退出游戏，回到房间中。

### 实例

*  /

```
request

response
看到这个说明网站能用
```



*  /user/register

```
request
{
    "username": "aa",
    "password": "sdsfg"
}
response
{
    "code": 10000,
    "data": "register successful!",
    "message": "ok"
}
```



*  /user/login

```
request
{
    "username": "aa",
    "password": "sdsfg"
}
response
{
    "code": 10000,
    "data": "login successful!",
    "message": "ok"
}
```

* /user/get

```
request
{
    "username": "aa"
}

response
{
    "code": 10000,
    "data": {
        "age": 18,
        "gender": "male",
        "username": "aa"
    },
    "message": "ok"
}
```

* /user/modify

```
request
{
    "username": "aa",
    "age":20
}

response
{
    "code": 10000
}
```

* /room/create

```
request

response
{
    "code": 10000,
    "data": "create room successful,room id is 1591500936",
    "message": "ok"
}
```

* /room/join

```
request
{
    "rid": 1591507233
}

resposne
{
    "code": 10000,
    "data": "successful!",
    "message": "ok"
}
```

* /room/password

```
request
{
    "password": "sdf"
}

response
{
    "code": 10000,
    "data": "set password successfully!",
    "message": "ok"
}
```





* /room/ready

```
request

response
{
    "code": 10000,
    "data": "already prepared, all player has prepared,can start the game!",
    "message": "ok"
}
```

* /game/start

```
request

response
{
    "code": 10000,
    "data": "The game starts successfully, you are a white pawn, please  perform your turn",
    "message": "ok"
}
```

* /game/paly

```
request
{
    "x": 4,
    "y": 4
}

response

# 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 # 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 # 0 0 * 0 0 0 0 0 0 0 0 0 
0 0 0 # 0 * 0 0 0 0 0 0 0 0 0 
0 0 0 0 # 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 * 0 0 * 0 0 0 0 0 0 0 0 0 
0 0 * 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 
0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
```

* /game/