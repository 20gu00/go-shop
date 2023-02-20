# 介绍
go-shop，go开发的商城系统，分api层和rpc层，api层主要是用gin开发，rpc层主要是用grpc开发。
## 技术栈
相关的技术知识和组件等：
- gin 开发api层服务
- grpc 开发rpc层微服务
- protobuf
- mysql
- gorm 操作数据库
- redis
- viper 本地解析配置文件及热更新
- consul 做注册中心，服务发现和负载均衡和健康检查
- nacos，除了viper本地配置解析，也提供了nacos来做分布式配置中心
- zap
- jwt 分布式架构中采用token来传递请求的用户信息比session机制好
- 分布式锁
- validator
- docker 快速部署开发环境

