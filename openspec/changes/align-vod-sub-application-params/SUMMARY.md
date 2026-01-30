# Summary: Align VOD Sub Application Parameters

## ✅ 提案已创建完成

### Change ID
`align-vod-sub-application-params`

### 提案概述
完善 `tencentcloud_vod_sub_application` 资源，支持腾讯云 VOD `CreateSubAppId` API 的所有参数，实现与云 API 的完全对齐。

### 新增参数

| 参数 | 类型 | 是否必填 | ForceNew | 说明 |
|------|------|---------|----------|------|
| `type` | String | 否 | ✅ | 应用类型：`AllInOne`（一体化）或 `Professional`（专业版） |
| `mode` | String | 否 | ✅ | 应用模式：`fileid`（仅FileID）或 `fileid+path`（FileID & Path） |
| `storage_region` | String | 否 | ✅ | 存储地域，如 `ap-guangzhou`、`ap-beijing` |
| `tags` | Map | 否 | ❌ | 标签键值对，最多10个，用于资源管理 |

### 关键特性

#### ✅ 向后兼容
- 所有新参数都是可选的
- 默认值保持与当前行为一致
- 不会破坏现有 Terraform 配置

#### ✅ ForceNew 行为
- `type`、`mode`、`storage_region` 创建后不可修改
- 修改这些参数将触发资源重建
- 符合腾讯云 API 的限制

#### ✅ Tags 支持
- 支持在创建时设置标签
- 支持更新标签（如果 API 支持）
- 标签数量限制：最多10个

### 使用示例

#### 完整配置
```hcl
resource "tencentcloud_vod_sub_application" "complete" {
  name           = "my-professional-app"
  status         = "On"
  description    = "Professional sub application"
  type           = "Professional"
  mode           = "fileid+path"
  storage_region = "ap-guangzhou"
  
  tags = {
    "team"        = "media"
    "environment" = "production"
    "project"     = "video-platform"
  }
}
```

#### 基础配置（向后兼容）
```hcl
resource "tencentcloud_vod_sub_application" "basic" {
  name        = "my-sub-app"
  status      = "On"
  description = "Basic sub application"
  # type 默认为 "AllInOne"
  # mode 默认为 "fileid"
}
```

### 实施阶段

#### Phase 1: Schema 和核心实现 (4 tasks)
- 添加新的 schema 字段定义
- 更新 Create 函数支持新参数
- 添加参数验证逻辑
- 编写基础单元测试

#### Phase 2: Read 和 Update 函数 (3 tasks)
- 更新 Read 函数文档
- 实现 Tags 更新逻辑
- 处理状态一致性

#### Phase 3: 测试 (18 tasks)
- 单元测试（5个）
- 验收测试 - 基础（8个）
- 验收测试 - Tags（8个）
- 验收测试 - 完整场景（4个）

#### Phase 4: 文档 (9 tasks)
- 更新资源文档
- 添加使用示例
- 编写迁移指南

#### Phase 5: 代码质量 (8 tasks)
- 代码格式化和 Lint
- 代码审查准备
- 集成测试

#### Phase 6: 发布准备 (9 tasks)
- Changelog
- 最终验证
- PR 准备

**总计：101 个任务**

### 技术考虑

#### 1. API 限制
- ⚠️ `DescribeSubAppIds` API 不返回 `Type`、`Mode`、`StorageRegion`
- 解决方案：Read 函数中保持这些字段不变，依赖 Terraform state
- 需要在文档中明确说明

#### 2. Tags 更新
- ⚠️ 需要确认 VOD API 是否支持 Tags 更新
- 如果支持：使用 VOD 或统一 Tag Service API
- 如果不支持：将 Tags 标记为 ForceNew

#### 3. StorageRegion 验证
- ⚠️ 有效地域列表未在文档中明确
- 解决方案：不进行客户端验证，依赖 API 错误消息

### 风险和缓解

| 风险 | 影响 | 缓解措施 |
|------|------|---------|
| Tags API 不支持更新 | Tags 修改需要重建资源 | 标记 Tags 为 ForceNew |
| StorageRegion 有效值未知 | 用户可能输入无效值 | 依赖 API 返回错误 |
| Read 函数无法验证 ForceNew 字段 | 状态漂移检测受限 | 文档中明确说明限制 |

### 成功标准
1. ✅ 所有新参数在 Create 时正确传递
2. ✅ ForceNew 参数修改触发资源重建
3. ✅ Tags 支持创建和更新
4. ✅ 所有验收测试通过
5. ✅ 文档完整准确
6. ✅ 向后兼容现有配置

### 下一步行动
1. ✅ 提案已创建并验证
2. ⏳ 等待提案审批
3. ⏳ 开始 Phase 1 实施（审批后）

### 文件位置
- 📄 Proposal: `openspec/changes/align-vod-sub-application-params/proposal.md`
- 📋 Tasks: `openspec/changes/align-vod-sub-application-params/tasks.md`
- 📊 Summary: `openspec/changes/align-vod-sub-application-params/SUMMARY.md`

### 相关资源
- 当前实现：`tencentcloud/services/vod/resource_tc_vod_sub_application.go`
- VOD SDK：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717/models.go`
- CreateSubAppIdRequest：Lines 7486-7507
- ResourceTag：Lines 26295-26300

---

**状态**: ✅ 提案完成，等待审批
**任务进度**: 0/101
**预计影响**: 增强功能，向后兼容，无破坏性变更
