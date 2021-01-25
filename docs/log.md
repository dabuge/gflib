# log 方法扩展

slog 的出现纯粹是为了偷懒，对于常见的业务类型，加上统一的前缀，便于链式调用。

## 配置示例

```yaml
# Logger.
logger:
  Path: "/log/app/example-dir"
  Level: "all"
  Stdout: true
```

配置完全遵循 goframe 官方的 log 配置。

## 调用示例

```go
slog.Init().S3().Error(err)
slog.Init().Mongodb().Error(err)
slog.Init().Redis().Error(err)
slog.Init().Mysql().Error(err)
slog.Init().Cache().Error(err)
//原生调用
slog.Init().Error(err)
```

