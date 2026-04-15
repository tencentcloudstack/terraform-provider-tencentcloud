## 1. 资源代码实现

- [x] 1.1 创建 `resource_tc_teo_dns_record_v12.go` 文件，定义资源 Schema
- [x] 1.2 实现资源的 `Create` 函数，调用 `CreateDnsRecord` 接口
- [x] 1.3 实现资源的 `Read` 函数，调用 `DescribeDnsRecords` 接口
- [x] 1.4 实现资源的 `Update` 函数，调用 `ModifyDnsRecords` 接口
- [x] 1.5 实现资源的 `Delete` 函数，调用 `DeleteDnsRecords` 接口
- [x] 1.6 实现资源的 ID 解析和生成逻辑（复合 ID 格式：`zone_id#record_id`）
- [x] 1.7 添加资源的 `Timeouts` 配置（虽然 API 是同步接口，但保留以备将来扩展）
- [x] 1.8 在所有 CRUD 函数中添加 `defer tccommon.LogElapsed()` 和 `defer tccommon.InconsistentCheck()`

## 2. 单元测试编写

- [x] 2.1 创建 `resource_tc_teo_dns_record_v12_test.go` 文件
- [x] 2.2 实现测试辅助函数，mock TEO 云 API 客户端
- [x] 2.3 编写 `TestAccTencentCloudTeoDnsRecordV12_basic` 测试用例，测试基本 CRUD 操作
- [x] 2.4 编写 `TestAccTencentCloudTeoDnsRecordV12_withWeight` 测试用例，测试权重配置
- [x] 2.5 编写 `TestAccTencentCloudTeoDnsRecordV12_withLocation` 测试用例，测试解析线路配置
- [x] 2.6 编写 `TestAccTencentCloudTeoDnsRecordV12_withPriority` 测试用例，测试 MX 优先级配置
- [x] 2.7 编写错误处理测试用例，验证 API 调用失败时的错误处理逻辑

## 3. 资源样例文件

- [x] 3.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v12.md` 样例文件
- [x] 3.2 编写 A 类型 DNS 记录的基本配置示例
- [x] 3.3 编写 CNAME 类型 DNS 记录配置示例
- [x] 3.4 编写 MX 类型 DNS 记录配置示例
- [x] 3.5 编写包含权重和解析线路的高级配置示例
- [x] 3.6 编写完整的资源导入、创建、更新和删除操作示例

## 4. 资源注册

- [x] 4.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中导入新的资源文件
- [x] 4.2 在 TEO 服务包的 `resources` map 中注册 `tencentcloud_teo_dns_record_v12` 资源
- [x] 4.3 验证资源注册是否成功，确保 Terraform 可以识别该资源

## 5. 文档生成

- [ ] 5.1 运行 `make doc` 命令生成 `website/docs/r/teo_dns_record_v12.md` 文档（需要在收尾阶段完成）
- [ ] 5.2 验证生成的文档包含所有必需参数和可选参数的说明（需要在收尾阶段完成）
- [ ] 5.3 验证生成的文档包含参数类型、默认值和约束条件说明（需要在收尾阶段完成）
- [ ] 5.4 验证生成的文档包含完整的资源示例（需要在收尾阶段完成）

## 6. 验证和测试

- [ ] 6.1 运行单元测试：`go test -v ./tencentcloud/services/teo/ -run TestAccTencentCloudTeoDnsRecordV12`（禁止执行测试）
- [ ] 6.2 验证所有测试用例通过（禁止执行测试）
- [x] 6.3 检查代码是否符合 Go 语言规范和项目代码风格
- [x] 6.4 验证样例文件语法正确，可以被 Terraform 正确解析
- [ ] 6.5 验证生成的文档格式正确，内容完整（需要在收尾阶段完成）
- [x] 6.6 检查新增文件是否包含在正确的目录结构中
