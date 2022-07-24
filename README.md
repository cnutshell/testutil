## Introduction

服务于 DnStore / LogStore 的集成测试框架，通过模拟不同故障或操作，测试集群运行是否满足预期。

故障或操作主要考虑以下几方面：

1. 网络分区故障
2. Store 启动
3. Store 停止

基于以上故障或者操作，验证集群运行是否满足预期以及数据是否完整。

## FAQ

- [x] 除 DnStore/LogStore 外，测试的集群中**还应当包含哪些组件**

  目前的需求，只包含 LogStore 和 DnStore 即可；

  hakeeper 通过 LogStore 启动；

- [x] DnStrore/LogStore 等组件的配置项有哪些

  LogStore 的配置是不是在 logservice.Config 里面

  DnStore 的配置项目前不明确

- [x] 测试框架中启动 Store 还是 Node

  DnStore or LogStore

- [x] 为验证集群状态验证方式，测试集群支持哪些操作以及验证方式

  ```go
  // 支持的操作
  type TestOperation interface {
    // ......
  }
  
  // 支持状态检查
  type TestAwareness interface {
    // ......
  }
  
  // 等待集群状态变化
  type TestAssertState interface {
    // ......
  }
  ```

## TODO

- [x] 增加用于测试验证的方法
- [ ] integrate with real log store
- [ ] integrate with real dn store
- [ ] 是否需要并发安全
- [ ] FIX: get rid of Bootstrap, Service would do this job internally.
- [x] FIX: start store first, then bootstrap
- [ ] integration test framwork directory path and name