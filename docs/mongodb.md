# MongoDB 连接文档
MongoDB 连接有以下几个特点：
1. 采用 goframe 框架的全局单例
2. 采用 MongoDB 官方提供的 go 组件
3. 配置文件可以配置官方的连接池大小和单例的连接池大小
4. 配置文件遵循 goframe 框架的惯例

## 配置示例
```yaml
mongodb:
  default:
    uri: "mongodb://127.0.0.1:27017/?connectTimeoutMS=10000&maxPoolSize=100&minPoolSize=10"

    # 这两项不归 options.ClientOptions 管理
    database: "test_mdb"
    poolNum: 30
  second_db:
    uri: "mongodb://127.0.0.2:27017/?connectTimeoutMS=10000&maxPoolSize=100&minPoolSize=10"

    # 这两项不归 options.ClientOptions 管理
    database: "test_mdb2"
    poolNum: 30
```

## 使用示例
```go
//获取默认连接
link := smongodb.Conn()
//获取指定的集合；使用 GetColl() 方法时，将连接配置文件里面指定的数据库
coll := link.GetColl("themex")
opts := options.FindOneAndUpdate().SetUpsert(true)
var updatedDocument bson.M
err := coll.FindOneAndUpdate(
    context.Background(),
    bson.D{
        {"shop_id", gconv.Int64(shopId)},
    },
    bson.D{
        {"$set", bson.D{
            {"content", gconv.String(content)},
        }},
    },
    opts,
).Decode(&updatedDocument)
if err != nil && err != mongo.ErrNoDocuments {
    return nil, err
}
return updatedDocument, nil
```

```go
//获取指定的连接和要操作的集合
coll := smongodb.Conn("second_db").GetColl("themex")
```

```go
//使用配置文件的连接参数，连接其他的数据库
link := smongodb.Conn("default")
sess := link.GetSession()
db := sess.Database("otherDb")
```

## 配置参数说明

- uri

uri 参数请参考这个文档 https://docs.mongodb.com/manual/reference/connection-string/

- database

指定要连接的数据库，直接使用 GetColl() 方法时，会连接配置文件指定的数据库

- pollNum

连接池大小，第一次连接时，程序会建立连接池。go driver 也有独立的连接池，所以最后的连接数是这个参数与uri 参数配置的连接池数量的乘积。

## 健康检查

```go
//Mongodb 检查
err := smongodb.Conn("default").Ping()
if err != nil {
    r.Response.WriteStatusExit(500)
}
err = smongodb.Conn("second_db").Ping()
if err != nil {
    r.Response.WriteStatusExit(500)
}
```

使用 smongodb 请务必配置健康检查，避免 MongoDB 连接失败导致的未知错误