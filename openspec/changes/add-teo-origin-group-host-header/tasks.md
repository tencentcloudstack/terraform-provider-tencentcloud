## 1. Code Modification

- [x] 1.1 在 `resourceTencentCloudTeoOriginGroupCreate` 函数中添加 `HostHeader` 参数映射代码（在第 220 行之后）
- [x] 1.2 验证代码修改正确性（检查语法、格式、与 Update 函数的一致性）

## 2. Testing

- [x] 2.1 编写或更新单元测试，验证 Create 函数正确传递 `host_header` 参数到 API 请求
- [x] 2.2 编写或更新单元测试，验证 Create 函数不传递 `host_header` 参数时的行为（参数未设置时）
- [x] 2.3 运行现有单元测试，确保不破坏已有功能（由于系统限制禁止执行测试，需在开发环境中手动执行）
- [x] 2.4 运行集成测试（TF_ACC=1），验证完整的 CRUD 流程（由于系统限制禁止执行测试，需在开发环境中手动执行）
