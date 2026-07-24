## 1. Implementation

- [x] 1.1 在 service_tencentcloud_ga2.go 中添加 `isGa2ResourceNotFoundError` 公共函数
- [x] 1.2 修改所有 describe 服务函数的 retry 块，添加 ResourceNotFound 错误处理
- [x] 1.3 修改所有资源 Read 方法，添加 `d.IsNewResource()` 检查
- [x] 1.4 验证代码编译通过
- [ ] 1.5 创建 changelog 文件