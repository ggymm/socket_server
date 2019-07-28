#### 生成消息体

##### 下载编译器
[GITHUB下载地址](https://github.com/protocolbuffers/protobuf/releases/download/v3.9.0-rc1/protoc-3.9.0-rc-1-win64.zip)

##### 安装插件

>​	下载好相关插件需要拷贝到同一目录
>
>```bash
>go install -v github.com/gogo/protobuf/protoc-gen-gogofaster
>go install -v github.com/davyxu/cellnet/protoc-gen-msg
>```

##### 生成消息

>​	执行以下命令生成代码
>
>```bash
>protoc --plugin=protoc-gen-gogofaster=protoc-gen-gogofaster.exe --gogofaster_out=. message.proto
>protoc --plugin=protoc-gen-msg=protoc-gen-msg.exe --msg_out=msgid.go:. message.proto
>```

##### 使用
> 将生成的相关的代码拷贝到项目中即可使用



#### 创建客户端

##### 安装

```bash
npm install -g @vue/cli
```



