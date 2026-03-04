## Implementation Tasks

### 1. Service Layer Implementation
- [x] 1.1 在 `service_tencentcloud_mongodb.go` 中添加 `EnableSRVConnectionUrl` 服务方法（调用 SDK 并等待异步任务完成）
- [x] 1.2 在 `service_tencentcloud_mongodb.go` 中添加 `DescribeSRVConnectionDomain` 服务方法
- [x] 1.3 在 `service_tencentcloud_mongodb.go` 中添加 `ModifySRVConnectionUrl` 服务方法（调用 SDK 并等待异步任务完成）
- [x] 1.4 在 `service_tencentcloud_mongodb.go` 中添加 `DisableSRVConnectionUrl` 服务方法

### 2. Resource Implementation
- [x] 2.1 创建 `resource_tc_mongodb_instance_srv_connection.go` 文件
- [x] 2.2 实现 Resource Schema 定义（instance_id [Required, ForceNew], domain [Optional + Computed]）
- [x] 2.3 实现 Create 函数（调用 EnableSRVConnectionUrl，如果 domain 有值则在创建后调用 ModifySRVConnectionUrl）
- [x] 2.4 实现 Read 函数（调用 DescribeSRVConnectionDomain 查询状态，填充 domain 字段）
- [x] 2.5 实现 Update 函数（当 domain 变化时调用 ModifySRVConnectionUrl 并等待异步任务完成）
- [x] 2.6 实现 Delete 函数（调用 DisableSRVConnectionUrl 关闭功能）
- [x] 2.7 实现 Import 支持（通过 instance_id 导入现有配置，自动填充 domain）
- [x] 2.8 添加日志记录和错误处理
- [x] 2.9 添加 defer LogElapsed 和 InconsistentCheck

### 3. Testing
- [x] 3.1 创建 `resource_tc_mongodb_instance_srv_connection_test.go` 文件
- [x] 3.2 实现基本的 Create/Read/Import 验收测试
- [x] 3.3 实现不带 domain 的测试场景（使用默认域名，验证 domain 字段被自动填充）
- [x] 3.4 实现带 domain 的测试场景（验证自定义域名正确设置）
- [x] 3.5 实现 domain 更新测试场景（从有到有的更新）
- [x] 3.6 实现 Import 功能测试（验证 domain 正确导入）
- [x] 3.7 验证异步任务等待逻辑正确性

### 4. Provider Registration
- [x] 4.1 在 `tencentcloud/provider.go` 的 ResourcesMap 中注册新资源
- [x] 4.2 在 `tencentcloud/provider.md` 的 MongoDB Resource 列表中声明新资源（按字母顺序）

### 5. Documentation
- [x] 5.1 创建 `resource_tc_mongodb_instance_srv_connection.md` 资源文档
- [x] 5.2 添加基本用法示例（不带 domain，说明系统会自动使用默认域名）
- [x] 5.3 添加自定义域名用法示例（指定 domain 参数）
- [x] 5.4 说明 domain 字段的 Optional + Computed 特性
- [x] 5.5 添加 Import 示例
- [x] 5.6 确保文档格式符合 `gendoc` 工具要求（纯文本标题，无 Markdown 标题符号）
- [x] 5.7 运行 `make doc` 生成网站文档到 `website/docs/r/mongodb_instance_srv_connection.html.markdown`

### 6. Validation and Quality
- [x] 6.1 运行 `go build` 确保编译通过
- [x] 6.2 运行 `gofmt` 进行代码格式化
- [x] 6.3 SDK 版本已更新到 v1.3.48
- [x] 6.4 编译成功，无错误
- [ ] 6.5 运行验收测试确保功能正常（需要实际环境测试）

### 7. Git Branch and Pull Request
- [ ] 7.1 创建功能分支 `git checkout -b feat/mongodb-srv-connection`
- [ ] 7.2 提交所有变更 `git add .` 和 `git commit -m "feat: add MongoDB instance SRV connection resource"`
- [ ] 7.3 推送分支到远程仓库 `git push origin feat/mongodb-srv-connection`
- [ ] 7.4 在 GitHub/GitLab 上创建 Pull Request
- [ ] 7.5 记录 PR ID（例如：#3834）

### 8. Changelog
- [ ] 8.1 使用 PR ID 在 `.changelog/` 目录创建 changelog 文件（格式：`.changelog/<PR_ID>.txt`，例如：`.changelog/3834.txt`）
- [ ] 8.2 在 changelog 文件中添加变更说明（enhancement 类型）

## Implementation Notes

### API 实际行为（基于 SDK v1.3.48）

根据实际的 TencentCloud MongoDB SDK，API 的行为与最初需求略有不同：

1. **EnableSRVConnectionUrl**：只接受 `InstanceId` 参数，不支持在创建时指定自定义域名
2. **DescribeSRVConnectionDomain**：只返回 `Domain` 字段，没有 `SRVUrl` 字段
3. **ModifySRVConnectionUrl**：需要 `InstanceId` 和 `CustomDomain` 参数来修改自定义域名
4. **DisableSRVConnectionUrl**：只需要 `InstanceId` 参数

### 实现策略

由于 API 的限制，实现采用了以下策略：

1. **创建流程**：
   - 首先调用 `EnableSRVConnectionUrl` 开启 SRV 连接（使用系统默认域名）
   - 如果用户指定了 `domain`，立即调用 `ModifySRVConnectionUrl` 设置自定义域名
   - 两个异步任务都等待完成

2. **Schema 设计**：
   - 移除了 `srv_url` 字段（API 不返回完整的 SRV URL）
   - 保留 `domain` 字段作为 Optional + Computed 属性
   - 用户可以通过 `domain` 字段查看实际使用的域名（系统默认或自定义）

3. **更新操作**：
   - 当 `domain` 字段变化时，调用 `ModifySRVConnectionUrl` 更新自定义域名

