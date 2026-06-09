# Tasks: 修改 MongoDB availability_zone_list 字段为 Set 类型

## Phase 1: Schema 和代码修改

### Task 1.1: 修改 resource_tc_mongodb_instance.go
- [ ] 修改 Schema 定义 (第 105-116 行)
  - 将 `Type: schema.TypeList` 改为 `Type: schema.TypeSet`
- [ ] 修改 mongodbAllInstanceReqSet 函数 (第 235-238 行)
  - 将 `v.([]interface{})` 改为 `v.(*schema.Set).List()`
- [ ] 修改 resourceTencentCloudMongodbInstanceRead 函数 (第 477-486 行)
  - 添加去重逻辑,使用 map 去重后转换为 slice
  - 确保返回的数据兼容 Set 类型
- [ ] 执行 `go fmt tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`

**预期变更文件:**
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.go`

### Task 1.2: 修改 resource_tc_mongodb_sharding_instance.go
- [ ] 修改 Schema 定义 (第 38 行附近)
  - 将 `Type: schema.TypeList` 改为 `Type: schema.TypeSet`
- [ ] 修改创建函数中的数据读取逻辑 (第 159-161 行)
  - 将 `v.([]interface{})` 改为 `v.(*schema.Set).List()`
- [ ] 修改 Read 函数中的数据写入逻辑 (第 394 行附近)
  - 添加去重逻辑
- [ ] 执行 `go fmt tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`

**预期变更文件:**
- `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.go`

### Task 1.3: 修改 resource_tc_mongodb_readonly_instance.go
- [ ] 修改创建函数中的数据读取逻辑 (第 163 行)
  - 将 `v.([]interface{})` 改为 `v.(*schema.Set).List()`
- [ ] 执行 `go fmt tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

**预期变更文件:**
- `tencentcloud/services/mongodb/resource_tc_mongodb_readonly_instance.go`

## Phase 2: 测试验证

### Task 2.1: 代码编译验证
- [ ] 执行 `go build` 验证代码无编译错误
  ```bash
  cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
  go build
  ```

### Task 2.2: 运行单元测试
- [ ] 运行 MongoDB 相关测试
  ```bash
  go test -v ./tencentcloud/services/mongodb/...
  ```

### Task 2.3: 运行验收测试 (可选,需要云账号)
- [ ] 设置环境变量
  ```bash
  export TF_ACC=1
  export TENCENTCLOUD_SECRET_ID="your-secret-id"
  export TENCENTCLOUD_SECRET_KEY="your-secret-key"
  ```
- [ ] 运行 mongodb_instance 验收测试
  ```bash
  go test -v -run TestAccTencentCloudMongodbInstance ./tencentcloud/services/mongodb/
  ```
- [ ] 运行 mongodb_sharding_instance 验收测试
  ```bash
  go test -v -run TestAccTencentCloudMongodbShardingInstance ./tencentcloud/services/mongodb/
  ```

### Task 2.4: 状态迁移测试
- [ ] 创建测试用 Terraform 配置
- [ ] 使用旧版本 Provider 创建资源并生成状态文件
- [ ] 切换到新版本 Provider
- [ ] 执行 `terraform plan` 验证状态迁移
- [ ] 确认只有类型转换,无实际资源变更

## Phase 3: 文档更新

### Task 3.1: 更新资源文档
- [ ] 更新 `tencentcloud/services/mongodb/resource_tc_mongodb_instance.md`
  - 说明 `availability_zone_list` 为 Set 类型
  - 添加使用示例
- [ ] 更新 `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.md`
  - 说明 `availability_zone_list` 为 Set 类型
  - 添加使用示例
- [ ] 执行 `make doc` 生成网站文档
  ```bash
  cd /Users/yanxiang/Tencent/Golang/terraform-provider-tencentcloud
  make doc
  ```

**预期变更文件:**
- `tencentcloud/services/mongodb/resource_tc_mongodb_instance.md`
- `tencentcloud/services/mongodb/resource_tc_mongodb_sharding_instance.md`
- `website/docs/r/mongodb_instance.html.markdown` (自动生成)
- `website/docs/r/mongodb_sharding_instance.html.markdown` (自动生成)

### Task 3.2: 更新 CHANGELOG
- [ ] 在 `CHANGELOG.md` 中添加条目
  ```markdown
  ## [Unreleased]
  
  ### BREAKING CHANGES
  
  * **resource/tencentcloud_mongodb_instance**: `availability_zone_list` 字段类型从 List 改为 Set。升级后首次执行 plan 会显示状态变更,这是正常的类型转换,不会影响实际资源。
  * **resource/tencentcloud_mongodb_sharding_instance**: `availability_zone_list` 字段类型从 List 改为 Set。升级后首次执行 plan 会显示状态变更,这是正常的类型转换,不会影响实际资源。
  
  ### IMPROVEMENTS
  
  * **resource/tencentcloud_mongodb_instance**: 使用 Set 类型避免可用区顺序变化导致的误判 diff,并自动去重 ([#xxxx](链接))
  * **resource/tencentcloud_mongodb_sharding_instance**: 使用 Set 类型避免可用区顺序变化导致的误判 diff,并自动去重 ([#xxxx](链接))
  ```

**预期变更文件:**
- `CHANGELOG.md`

### Task 3.3: 创建迁移指南 (可选)
- [ ] 创建 `docs/guides/mongodb-availability-zone-list-migration.md`
  - 说明变更原因
  - 提供迁移步骤
  - 列出常见问题和解决方案

## Phase 4: 代码审查检查清单

### Task 4.1: 代码审查
- [ ] 所有 `availability_zone_list` 的读取逻辑已更新为 Set 类型
- [ ] 所有 `availability_zone_list` 的写入逻辑已添加去重
- [ ] 没有遗漏的 List 类型引用
- [ ] 代码风格符合项目规范
- [ ] 所有修改的文件都已执行 `go fmt`

### Task 4.2: 测试覆盖检查
- [ ] 编译通过
- [ ] 单元测试通过
- [ ] 验收测试通过 (如果可以运行)
- [ ] 手动测试状态迁移场景

### Task 4.3: 文档完整性检查
- [ ] 资源文档已更新
- [ ] 网站文档已生成
- [ ] CHANGELOG 已更新
- [ ] 迁移指南已创建 (如果需要)

## Phase 5: 发布准备

### Task 5.1: 版本标记
- [ ] 确定发布版本号 (建议主版本升级,如 v2.0.0)
- [ ] 在 Git 中创建标签

### Task 5.2: 发布说明
- [ ] 准备发布说明,重点强调:
  - BREAKING CHANGE 标识
  - 变更原因和优势
  - 迁移步骤
  - 影响范围
  - 回滚建议

### Task 5.3: 社区通知
- [ ] 在 GitHub Releases 中发布说明
- [ ] 通知用户关注状态迁移
- [ ] 准备回答用户可能的问题

## 估算工作量

- **Phase 1 (代码修改)**: 2-3 小时
- **Phase 2 (测试验证)**: 2-4 小时
- **Phase 3 (文档更新)**: 1-2 小时
- **Phase 4 (代码审查)**: 1 小时
- **Phase 5 (发布准备)**: 1-2 小时

**总计**: 7-12 小时

## 风险和注意事项

1. ⚠️ **破坏性变更**: 必须在主版本更新时发布
2. ⚠️ **状态迁移**: 需要充分测试状态文件格式转换
3. ⚠️ **用户通知**: 必须清晰告知用户此变更的影响
4. ⚠️ **回滚计划**: 如果出现问题,用户可能需要降级到旧版本

## 依赖和前置条件

- Go 1.17+ 开发环境
- Terraform Plugin SDK v2
- (可选) 腾讯云账号用于运行验收测试
- (可选) 已有的 MongoDB 实例用于测试状态迁移
