# AWS 连接

提供简便和统一的 aws 连接方式

## 配置示例

```yaml
aws:
  # 图片资源的 s3
  default:
    AccessKey: "AKDDDDDDDDDDDDDDD"
    AccessSecret: "Secret2SecretSecret2SecretSecret2SecretSecret2Secret"
    RoleARN: "arn:aws:iam::8888888888888:role/my_dev"
    RoleSessionName: "project-name"
    Region: "us-west-2"
    # 使用 S3 时需要的桶名称
    XshoppyLiquidTemplatesBucket: "s3-bucket-dev"
  # 静态资源的 s3
  static_storage:
    AccessKey: "AKFFFFFFFFFFFFFFF"
    AccessSecret: "Secret2SecretSecret2SecretSecret2SecretSecret2Secret"
    Region: "us-west-2"
    EndPoint: "https://s3-us-west-2.amazonaws.com"
    Bucket: "s3-other-bucket-dev"

```

## 配置说明

支持通过该配置文件直接获取 session，或者获取 sts session。

需要支持 sts 时，`AccessKey`、`AccessSecret`、`RoleARN`、`RoleSessionName`、`Region`参数不能缺少。

所有的配置项支持通过 `saws.Client().GetConfig*("configName") `方法获取，如：

```go
awsClient := saws.Client("static_storage")
Bucket := awsClient.GetConfigString("Bucket")
```

## 使用说明

```go
//获取 STS 后，从 s3 下载文件
awsClient := saws.Client()
var sess = awsClient.GetSessionFromSts()
svc := s3.New(sess)
input := &s3.GetObjectInput{
    Bucket: awsClient.GetConfigString("XshoppyLiquidTemplatesBucket"),
    Key:    aws.String("example/ss.zip"),
}
result, err := svc.GetObject(input)
if err != nil {
    g.log().Error(err)
    return "", err
}
if result != nil {
    defer result.Body.Close()
    content, _ := ioutil.ReadAll(result.Body)
    return gconv.String(content), nil
}
return "", errors.New("文件读取错误")
```

```go
//获取session，上传文件到 s3
awsClient := saws.Client("static_storage")
sess := awsClient.GetSession(&aws.Config{
    Endpoint:         awsClient.GetConfigString("EndPoint"),
    DisableSSL:       aws.Bool(true),
    S3ForcePathStyle: aws.Bool(false), 
})
uploader := s3manager.NewUploader(sess)
ui := &s3manager.UploadInput{
    Bucket: awsClient.GetConfigString("Bucket"),
    Key:    aws.String("example/ss.zip"),
    Body:   bytes.NewReader(file),
}
result, err := uploader.Upload(ui)
if err != nil {
    g.log().Error(err)
    return "", err
}
return result.Location, err
```

