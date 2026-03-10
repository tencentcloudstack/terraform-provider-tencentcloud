# tcr-namespace Specification

## Purpose
TBD - created by archiving change add-tcr-namespace-tag-update. Update Purpose after archive.
## Requirements
### Requirement: 标签 Schema 定义

`tencentcloud_tcr_namespace` 资源 SHALL 在其 schema 定义中包含 `tags` 字段以支持标签管理。

#### Scenario: Schema 包含标签字段

- **假设** 有一个 TCR namespace 资源定义
- **当** 定义 schema 时
- **那么** SHALL 包含 `tags` 字段并具备以下属性:
  - 类型: `schema.TypeMap`
  - Optional: `true`
  - Computed: `true`
  - 元素类型: `schema.TypeString`
  - Description 说明标签键值对的作用

### Requirement: 标签更新检测

资源更新函数 SHALL 检测 Terraform 配置中标签是否发生了变更。

#### Scenario: 检测标签变更

- **假设** 有一个具有现有标签的 TCR namespace 资源
- **当** 用户在 Terraform 配置中修改标签时
- **那么** 更新函数 SHALL 使用 `d.HasChange("tags")` 检测到该变更

### Requirement: 标签差异计算

系统 SHALL 计算新旧标签值之间的差异,确定需要添加、修改或删除哪些标签。

#### Scenario: 计算标签差异

- **假设** 旧标签为 `{"env": "test", "team": "backend"}`,新标签为 `{"env": "prod", "owner": "alice"}`
- **当** 计算标签差异时
- **那么** 系统 SHALL:
  - 识别 `{"env": "prod", "owner": "alice"}` 为需要添加/修改的标签
  - 识别 `{"team"}` 为需要删除的标签

### Requirement: 统一标签服务集成

资源 SHALL 使用腾讯云统一的标签管理服务来应用标签更新。

#### Scenario: 使用 TagService 进行更新

- **假设** 检测到标签变更
- **当** 应用标签更新时
- **那么** 系统 SHALL:
  - 导入 `svctag` 包
  - 创建 `TagService` 实例
  - 使用正确的参数调用 `tagService.ModifyTags()`

### Requirement: TCR 资源命名

系统 SHALL 在调用标签管理 API 时为 TCR namespace 构建正确的资源名称格式。

#### Scenario: 为标签操作构建资源名称

- **假设** instanceId 为 `tcr-abc123`,namespaceName 为 `my-namespace`,区域为 `ap-guangzhou`
- **当** 为标签操作构建资源名称时
- **那么** 系统 SHALL 构建: `qcs::tcr:ap-guangzhou::namespace/{instanceId}/{namespaceName}`
- **并且** SHALL 使用 `tccommon.BuildTagResourceName("tcr", "namespace", region, resourceId)`,其中 resourceId 为 `{instanceId}/{namespaceName}`

### Requirement: 标签更新错误处理

系统 SHALL 正确处理标签更新操作期间的错误,并返回适当的错误消息。

#### Scenario: 处理标签更新失败

- **假设** 正在进行标签更新操作
- **当** `tagService.ModifyTags()` 返回错误时
- **那么** 系统 SHALL:
  - 将错误返回给 Terraform
  - 不继续执行资源读取操作
  - 保留现有状态

### Requirement: 标签读取支持

资源读取函数 SHALL 从 TCR namespace 中读取标签并填充到 state。

#### Scenario: 将标签读取到 state

- **假设** 有一个带标签的 TCR namespace
- **当** 读取资源时
- **那么** 系统 SHALL:
  - 从 namespace 响应中读取标签
  - 将标签列表转换为 map 格式
  - 使用 `d.Set("tags", tagMap)` 在 Terraform state 中设置标签

### Requirement: 文档更新

资源文档 SHALL 包含演示标签用法的示例。

#### Scenario: 记录标签用法

- **假设** 有资源文档文件
- **当** 编写资源文档时
- **那么** 文档 SHALL 包含:
  - 至少一个展示标签定义的示例
  - `tags` 参数的描述
  - 标签格式和约束说明

