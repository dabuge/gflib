# Redis 连接文档

由于 goframe 框架底层使用了 **[redigo](https://github.com/gomodule/redigo)** 这个库，这个库原生不支持 redis 的集群模式，因此另外扩展了 redis 的连接，使用了 **[go-redis](https://github.com/go-redis/redis)** 库，以便支持 redis 集群模式。

## 配置示例

```yaml
redis:
  # 集群模式，使用 go-redis 类库连接，使用 sredis.ClusterClient("default") 连接
  default:
    Addrs:
      - "clustercfg-redis.usw2.cache.amazonaws.com:6379"
    PoolSize: 100 # 连接池大小
    MinIdleConns: 20
    Password: "kdkdkiekdkdkdk"
    TLSConfig: { }

  # 非集群模式，注意配置的时候 Addr/DB 参数和集群模式不一样，使用 sredis.Client("second") 连接
  second:
    Addr: "127.0.0.1:6379"
    DB: 0
    PoolSize: 100 # 连接池大小
    MinIdleConns: 20
    Password: "kdkdkiekdkdkdk"

```

## 配置说明

集群模式和普通模式的配置项有所不同，具体可以参考：

集群模式支持的配置参数：https://pkg.go.dev/github.com/go-redis/redis/v8#ClusterOptions

普通模式支持的配置参数：https://pkg.go.dev/github.com/go-redis/redis/v8#Options

## 使用示例

集群模式使用 `sredis.ClusterClient("configName")`连接

普通模式使用 `sredis.Client("configName")`连接

```go
redisClient := sredis.ClusterClient()
cacheKey := "cache_key:" + data.Key
b, err := redisClient.HSetNX(context.TODO(), cacheKey, "ready", "content").Result()
if err != nil {
    g.log().Error(err)
}
```

## 健康检查

```go
//redis 检查
_, err := sredis.ClusterClient().Ping(context.TODO()).Result()
if err != nil {
    r.Response.WriteStatusExit(500)
}
//redis 检查
_, err = sredis.Client("second").Ping(context.TODO()).Result()
if err != nil {
    r.Response.WriteStatusExit(500)
}
r.Response.WriteStatusExit(200)
```



