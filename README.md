## Introduction

服务于 DnStore / LogStore 的集成测试框架，通过模拟不同故障或操作，测试集群运行是否满足预期。

故障或操作主要考虑以下几方面：

1. 网络分区故障
2. 节点启动
3. 节点停止

基于以上故障或者操作，验证集群运行是否满足预期以及数据是否完整。

## FAQ

- [x] 除 DnStore/LogStore 外，测试的集群中**还应当包含哪些组件**

  目前的需求，只包含 LogStore 和 DnStore 即可；

  hakeeper 通过 LogStore 启动；

- [x] DnStrore/LogStore 等组件的配置项有哪些

  LogStore 的配置是不是在 logservice.Config 里面

  DnStore 的配置项目前不明确

- [ ] 测试框架中启动 Store 还是 Node

- [ ] 为验证集群状态验证方式，测试集群支持哪些验证方式

## TODO

- [ ] 增加用于测试验证的方法
- [ ] 启动带 log store 的集群
- [ ] 启动带 dn store 的集群