# goblack_session

本项目演示如何在 Go 项目中集成基于 Redis 的 Session 管理。

## 依赖
- [goblack](https://github.com/varunmbet/goblack)
- [redis](https://github.com/varunmbet/redis)

请确保已正确安装上述依赖。

## 快速开始

```go
package main

import (
    "github.com/varunmbet/goblack"
    "github.com/varunmbet/redis"
    "your_module_path/session" // 替换为你的 session 包路径
    "your_module_path/models"  // 替换为你的 models 包路径
    "encoding/gob"
)

func main() {
    app := goblack.Instance("demo")
    redissession := session.NewRedisProvider()
    redissession.Init(goblack.SessionOptions{
        Session_lifetime: 1800, // session 有效期（秒）
        IDLength: 32,           // session id 长度
        Providerconfig: session.RedisOptions{
            Prefix: "session:cms:", // redis key 前缀
            Client: redis.Redisdb,  // 你的 redis 客户端实例
        },
    })
    app.SetSession(true, redissession)
    gob.Register(models.Sys_admin{}) // 注册需要序列化的结构体
    // ... 你的其他业务代码
}
```

## 说明
- `session.NewRedisProvider()` 创建一个 Redis session provider。
- `Init` 方法用于初始化 session 配置。
- `SetSession` 启用 session 管理。
- `gob.Register` 注册需要存储到 session 的自定义结构体。

## 注意事项
- 请确保 `redis.Redisdb` 已正确初始化并连接到你的 Redis 服务。
- 替换 `your_module_path` 为你实际的 Go module 路径。
- 如需自定义 session 结构体，请记得用 `gob.Register` 注册。

## License
MIT 