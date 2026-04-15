## 1. 资源 Schema 定义

- [x] 1.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11.go` 文件，定义资源 Schema
- [x] 1.2 在 Schema 中定义必填字段：zone_id, name, type, content
- [x] 1.3 在 Schema 中定义可选字段：ttl, weight, priority, location
- [x] 1.4 在 Schema 中定义只读字段（Computed）：status, created_on, modified_on
- [x] 1.5 在 Schema 中定义 Timeouts 块，支持 create, update, delete 超时配置

## 2. Service 层实现

- [x] 2.1 在 `tencentcloud/services/teo/service_tencentcloud_teo.go` 中实现 `CreateDnsRecord` 服务函数
- [x] 2.2 在 Service 层中实现 `DescribeDnsRecords` 服务函数，支持通过 RecordId 过滤
- [x] 2.3 在 Service 层中实现 `ModifyDnsRecords` 服务函数
- [x] 2.4 在 Service 层中实现 `DeleteDnsRecords` 服务函数
- [x] 2.5 实现复合 ID（zoneId#recordId）的构建和解析辅助函数

## 3. CRUD 函数实现

- [x] 3.1 实现 `Create` 函数，调用 CreateDnsRecord 接口创建 DNS 记录
- [x] 3.2 在 Create 函数中实现异步操作轮询逻辑（如果接口是异步的）
- [x] 3.3 实现 `Read` 函数，调用 DescribeDnsRecords 接口读取 DNS 记录
- [x] 3.4 在 Read 函数中实现 ResourceNotFound 错误处理逻辑
- [x] 3.5 实现 `Update` 函数，调用 ModifyDnsRecords 接口更新 DNS 记录
- [x] 3.6 在 Update 函数中实现字段变更检测，确保只更新可修改的字段
- [x] 3.7 在 Update 函数中实现不可修改字段（name、type）的错误返回
- [x] 3.8 实现 `Delete` 函数，调用 DeleteDnsRecords 接口删除 DNS 记录
- [x] 3.9 在 Delete 函数中实现删除后轮询确认逻辑

## 4. 错误处理和重试逻辑

- [x] 4.1 在所有 CRUD 函数中添加 `defer tccommon.LogElapsed()` 调用
- [x] 4.2 在所有 CRUD 函数中添加 `defer tccommon.InconsistentCheck()` 调用
- [x] 4.3 在 Create、Update、Delete 函数中使用 `helper.Retry()` 实现重试逻辑
- [x] 4.4 实现网络错误和 API 错误的正确处理和返回

## 5. 单元测试

- [x] 5.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11_test.go` 文件
- [x] 5.2 实现 Create 函数的单元测试，使用 mock 云 API
- [x] 5.3 实现 Read 函数的单元测试，使用 mock 云 API
- [x] 5.4 实现 Update 函数的单元测试，使用 mock 云 API
- [x] 5.5 实现 Delete 函数的单元测试，使用 mock 云 API
- [x] 5.6 实现 ResourceNotFound 场景的单元测试
- [x] 5.7 实现不可修改字段更新错误的单元测试

## 6. 示例文件

- [x] 6.1 创建 `tencentcloud/services/teo/resource_tc_teo_dns_record_v11.md` 示例文件
- [x] 6.2 在示例文件中添加基本创建示例
- [x] 6.3 在示例文件中添加带可选字段的创建示例（ttl、weight、priority、location）
- [x] 6.4 在示例文件中添加更新示例
- [x] 6.5 在示例文件中添加不同记录类型的示例（A、AAAA、CNAME、MX、TXT、NS、CAA、SRV）

## 7. 资源注册

- [x] 7.1 在 Provider 中注册新资源 `tencentcloud_teo_dns_record_v11`
- [x] 7.2 验证资源注册正确性

## 8. 集成测试

- [ ] 8.1 运行 Create 操作的集成测试（TF_ACC=1）
- [ ] 8.2 运行 Read 操作的集成测试（TF_ACC=1）
- [ ] 8.3 运行 Update 操作的集成测试（TF_ACC=1）
- [ ] 8.4 运行 Delete 操作的集成测试（TF_ACC=1）
- [ ] 8.5 运行完整的资源生命周期测试（Create → Read → Update → Delete）（TF_ACC=1）

## 9. 代码验证

- [ ] 9.1 执行 `gofmt` 格式化资源代码文件
- [ ] 9.2 执行 `go build` 验证代码可编译
- [ ] 9.3 执行单元测试确保所有测试通过

## 10. 文档生成

- [ ] 10.1 运行 `make doc` 命令生成 `website/docs/r/teo_dns_record_v11.html.markdown` 文档
- [ ] 10.2 验证生成的文档内容正确性和完整性

## 11. Changelog

- [ ] 11.1 在 `.changelog/` 目录下创建变更日志文件
- [ ] 11.2 在变更日志中记录新增资源的详细信息

## 12. 代码提交和 PR 创建

- [ ] 12.1 检查所有修改的文件是否符合代码规范
- [ ] 12.2 创建特性分支并提交所有变更
- [ ] 12.3 推送到远程仓库
- [ ] 12.4 创建 Pull Request 到 master 分支
