package saws

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/frame/gins"
	"github.com/gogf/gf/os/gcache"
	"github.com/gogf/gf/os/gtimer"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"time"
)

// 获取临时token请求参数，用于前后端交互参数格式约定
type saws struct {
	Config map[string]interface{}
}

const (
	// Default group name for instance usage.
	DEFAULT_NAME   = "default"
	gAWS_NODE_NAME = "aws"
)

var (
	instances = gmap.NewStrAnyMap(true)
)

type Credentials struct {
	AccessKeyId     *string
	SecretAccessKey *string
	SessionToken    *string
	Expiration      int64
}

//Client 返回一个 go-redis 的集群单例
func Client(name ...string) *saws {
	config := gins.Config()
	key := DEFAULT_NAME
	if len(name) > 0 && name[0] != "" {
		key = name[0]
	}
	return instances.GetOrSetFuncLock(key, func() interface{} {
		var m map[string]interface{}
		if _, v := gutil.MapPossibleItemByKey(gins.Config().GetMap("."), gAWS_NODE_NAME); v != nil {
			m = gconv.Map(v)
		}
		if len(m) > 0 {
			if v, ok := m[key]; ok {
				saws := new(saws)
				saws.Config = gconv.Map(v)
				return saws
			} else {
				panic(fmt.Sprintf(`configuration for aws not found for group "%s"`, key))
			}
		} else {
			panic(fmt.Sprintf(`incomplete configuration for aws: "aws" node not found in config file "%s"`, config.FilePath()))
		}
		return nil
	}).(*saws)
}

//assumeRole 获取STS权限，后续需对接接口
func (s *saws) assumeRole() (cre *Credentials, err error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      s.GetConfigString("Region"),
		Credentials: credentials.NewStaticCredentials(gconv.String(s.Config["AccessKey"]), gconv.String(s.Config["AccessSecret"]), ""),
	})
	svc := sts.New(sess)
	result, err := svc.AssumeRole(&sts.AssumeRoleInput{
		RoleArn:         s.GetConfigString("RoleARN"),
		RoleSessionName: s.GetConfigString("RoleSessionName"),
	})
	if err != nil {
		g.Log().Async(true).Cat("s3").Error(err)
		return nil, err
	}
	var c Credentials
	c.AccessKeyId = result.Credentials.AccessKeyId
	c.SecretAccessKey = result.Credentials.SecretAccessKey
	c.SessionToken = result.Credentials.SessionToken
	//过期时间转换成时间戳，和 php 版本的保持一直
	c.Expiration = result.Credentials.Expiration.Unix()
	//STS token 的有效期一个小时(3600秒)，我们缓存3500秒
	gcache.Set("cache_key:aws:sts_token", c, 3500*time.Second)
	return &c, nil
}

//GetSts 从缓存中读取 STS token
func (s *saws) GetSts() (Credentials, error) {
	cache, err := gcache.Get("cache_key:aws:sts_token")
	if err != nil {
		g.Log().SetPrefix("[cache]")
		g.Log().Async(true).Error(err)
	}
	var c Credentials
	if cache == nil {
		c, _ := s.assumeRole()
		//每隔3000秒，重新设置一次缓存
		gtimer.AddSingleton(3000*time.Second, func() {
			_, _ = s.assumeRole()
		})
		return *c, nil
	}
	err = gconv.Struct(cache, &c)
	if err != nil {
		g.Log().SetPrefix("[s3]")
		g.Log().Async(true).Error(err)
		return Credentials{}, nil
	}
	return c, nil
}

//GetSessionFromSts 获取 AWS 的 session，以 sts 的方式
func (s *saws) GetSessionFromSts(cfgs ...*aws.Config) *session.Session {
	stsInfo, _ := s.GetSts()
	opts := session.Options{}
	opts.Config.Credentials = credentials.NewStaticCredentials(*stsInfo.AccessKeyId, *stsInfo.SecretAccessKey, *stsInfo.SessionToken)
	opts.Config.Region = s.GetConfigString("Region")
	opts.Config.MergeIn(cfgs...)
	sess := session.Must(session.NewSession(&opts.Config))
	return sess
}

//GetSession 获取 AWS 的 session
func (s *saws) GetSession(cfgs ...*aws.Config) *session.Session {
	opts := session.Options{}
	opts.Config.Credentials = credentials.NewStaticCredentials(gconv.String(s.Config["AccessKey"]), gconv.String(s.Config["AccessSecret"]), "")
	opts.Config.Region = s.GetConfigString("Region")
	opts.Config.MergeIn(cfgs...)
	sess := session.Must(session.NewSession(&opts.Config))
	return sess
}

//GetConfigString 以 string 类型获取配置文件的值
func (s *saws) GetConfigString(name string) *string {
	return aws.String(gconv.String(s.Config[name]))
}

//GetConfigBool 以 Bool 类型获取配置文件的值
func (s *saws) GetConfigBool(name string) *bool {
	return aws.Bool(gconv.Bool(s.Config[name]))
}

//GetConfigInt 以 Int 类型获取配置文件的值
func (s *saws) GetConfigInt(name string) *int {
	return aws.Int(gconv.Int(s.Config[name]))
}

//GetConfigInt64 以 Int64 类型获取配置文件的值
func (s *saws) GetConfigInt64(name string) *int64 {
	return aws.Int64(gconv.Int64(s.Config[name]))
}

//GetConfigFloat64 以 Float64 类型获取配置文件的值
func (s *saws) GetConfigFloat64(name string) *float64 {
	return aws.Float64(gconv.Float64(s.Config[name]))
}
