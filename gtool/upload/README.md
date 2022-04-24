# 文件上传工具 - Cli

## 目录：

- [1、特点](#1特点)
- [2、快速开始](#2快速开始)
- [3、配合Typora](#4配合Typora)

## 1、特点

- 利用阿里云Oss，上传图片
- `Typora` 配合使用，写一些markdown，很方便，不需要本地保存图片
- 支持多配置路径，适合上传多个文件
- 个人使用一般是将博客的图片全部上传到阿里云上，个人的一些资料也是，会把url保存住

## 2、快速开始

### 1、下载

```shell
go get -v github.com/anthony-dong/go-sdk/gtool
```

### 2、使用帮助

> ​	配置文件来自于 `go-tool --config 参数`，所以变更配置文件需要指定这个

```bash
➜  gtool upload --help
Name: File upload tool

Usage: gtool upload [flags]

Options:
  -d, --decode string   Set the upload file name decode method (uuid|url|md5) (default "uuid")
  -f, --file string     Set the local path of upload file
  -h, --help            help for upload
  -t, --type string     Set the upload config type (default "default")

Global Options:
      --config-file string   set the config file (default "/Users/bytedance/.go-tool.yaml")
      --log-level string     set the log level in "fatal|error|warn|info|debug" (default "debug")

To get more help with gtool, check out our guides at https://github.com/Anthony-Dong/go-sdk
```

### 3、快速开始

```shell
➜  gtool upload --file ./go.mod --decode md5
2022/04/23 22:06:53.102363 cli.go:53: [INFO] [upload] start config:
{
  "decode": "md5",
  "file": "/Users/bytedance/go/src/github.com/anthony-dong/golib/go.mod",
  "type": "default"
}
2022/04/23 22:06:53.224842 cli.go:84: [INFO] [upload] end success, url: https://xxxx.xxx-xxxx.xxxx.com/image/2022/4-23/34d1f91fb2e514b8576fab1a75a89a6b.mod
```

### 4、配置文件

> 支持json 或者 yaml， json需要key命令是下划线模式

```yaml
Upload:
  Bucket:
    default:
      AccessKeyId: xxxx
      AccessKeySecret: xxx
      Endpoint: oss-accelerate.xxxxx.com
      UrlEndpoint: xxxx.oss-accelerate.xxxx.com
      Bucket: xxxx
      PathPrefix: file # 前缀其实就是 UrlEndpoint/{PathPrefix}/y/m/d/${filename}
    image:
      AccessKeyId: xxxx
      AccessKeySecret: xxxx
      Endpoint: oss-accelerate.xxxxx.com
      UrlEndpoint: xxxx.oss-accelerate.xxxx.com
      Bucket: tyut
      PathPrefix: image
```

阿里云Oss端配置大概就是这些：

![image-20200914135934215](https://tyut.oss-accelerate.aliyuncs.com/image/2020/9-14/42cdf58e904e4dbeac06028639db9d40.png)

如果参数不输入 `-t`，则默认走 `default`标签！，所以一般推荐都设置一个default标签

```shell
➜  gtool upload -f ./go.mod --decode md5 --type pdf
2022/04/23 22:09:30.637767 cli.go:53: [INFO] [upload] start config:
{
  "decode": "md5",
  "file": "/Users/bytedance/go/src/github.com/anthony-dong/golib/go.mod",
  "type": "pdf"
}
2022/04/23 22:09:30.882701 cli.go:84: [INFO] [upload] end success, url: https://xxxx-xxx.oss-xxxxx.aliyuncs.com/pdf/2022/4-23/34d1f91fb2e514b8576fab1a75a89a6b.mod
```

## 4、配合Typora

只需要设置如下： 记得`go-tool`写成绝对路径，最后验证一下即可

我配置的命令是 `/Users/fanhaodong/go/bin/gtool  --log-level fatal upload --file  `

![image-20210116173613955](https://tyut.oss-accelerate.aliyuncs.com/image/2021/1-16/e60e7865de434e37a78856ad91a00c8a.png)

