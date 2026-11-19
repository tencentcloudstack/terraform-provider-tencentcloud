## Context

`tencentcloud_cfs_file_system` 资源用于管理腾讯云 CFS 文件系统，当前 Schema 已包含 `availability_zone`、`access_group_id`、`net_interface`、`protocol`、`storage_type`、`vpc_id`、`subnet_id`、`mount_ip`、`ccn_id`、`cidr_block`、`capacity`、`name`、`tags` 等参数，但未暴露文件系统加密能力。

腾讯云 CFS 服务的 `CreateCfsFileSystem` 接口已支持 `Encrypted` 入参，`DescribeCfsFileSystems` 接口返回的 `FileSystemInfo` 结构体中已包含 `Encrypted` 字段。CFS 没有更新加密状态的接口，因此 `Encrypted` 只能在创建时指定，修改时需触发重建。

**关键文件：**
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go` - 资源 CRUD 实现
- `tencentcloud/services/cfs/resource_tc_cfs_file_system_test.go` - 资源测试

### 云 API 参考信息（从 vendor 摘录，实施阶段无需回读大文件）

vendor 路径：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719/models.go`
包名：`github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfs/v20190719`

**1. CreateCfsFileSystemRequest 中的 Encrypted 字段（models.go:552-553）**

```go
// <p>文件系统是否加密，若留空则默认为不加密</p>
Encrypted *bool `json:"Encrypted,omitnil,omitempty" name:"Encrypted"`
```

**2. FileSystemInfo（DescribeCfsFileSystems 响应元素）中的 Encrypted 字段（models.go:3946-3950）**

```go
// <p>文件系统是否加密,true：代表加密，false：非加密</p>
Encrypted *bool `json:"Encrypted,omitnil,omitempty" name:"Encrypted"`

// <p>加密所使用的密钥，可以为密钥的 ID 或者 ARN</p>
KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
```

> 注：本变更仅新增 `Encrypted` 一个参数，不涉及 `KmsKeyId`。

### 现有 Schema 结构说明（与新增参数关联）

现有 `resource_tc_cfs_file_system.go` 中 Schema 为 `map[string]*schema.Schema`，字段平铺在顶层。Create 函数（第 135 行起）通过 `d.Get()`/`d.GetOk()` 读取参数构造 `cfs.NewCreateCfsFileSystemRequest()`；Read 函数（第 230 行起）通过 `cfsService.DescribeFileSystem()` 获取 `*cfs.FileSystemInfo`，再用 `d.Set()` 回填字段。Update 函数（第 304 行起）维护 `immutableArgs` 列表。新增 `encrypted` 参数将按相同模式平铺到顶层 Schema。

## Goals / Non-Goals

**Goals:**
- 在 `tencentcloud_cfs_file_system` 资源中新增 `encrypted` 参数（`TypeBool, Optional, ForceNew: true`）
- Create 函数中将 `encrypted` 设置到 `CreateCfsFileSystemRequest.Encrypted`
- Read 函数中将 `FileSystemInfo.Encrypted` 同步到 state
- 补充单元测试用例
- 更新文档示例

**Non-Goals:**
- 不支持更新加密状态（云 API 不提供该能力）
- 不新增 `KmsKeyId` 参数（本次变更范围仅限 `Encrypted` 一个参数）

## Decisions

### Decision 1: encrypted 参数标记为 ForceNew
**决策：** `encrypted` 参数设为 `Optional, ForceNew: true`。

**理由：**
- CFS 没有 `UpdateCfsFileSystem` 接口，无法在创建后修改加密状态
- 标记 `ForceNew` 后，用户修改 `encrypted` 会触发资源重建，符合 Terraform 语义
- 与现有 `protocol`、`storage_type` 等 ForceNew 字段保持一致

**替代方案：**
- 标记为不可变并加入 Update 的 immutableArgs → 用户修改时会返回 error，但 ForceNew 更符合 Terraform 惯例（重建而非报错）

### Decision 2: bool 类型参数的读取方式
**决策：** 在 Create 函数中使用 `d.Get("encrypted").(bool)` 直接读取，通过 `helper.Bool()` 转换为 `*bool` 设置到 request。由于该参数为 `Optional` 且无默认值，未设置时 `d.Get()` 返回 `false`，此时仍设置 `request.Encrypted = helper.Bool(false)`。

**理由：**
- 云 API 注释说明"若留空则默认为不加密"，即 `false` 与留空等价
- 始终设置 `Encrypted` 字段可保证 state 与配置一致，避免 drift

### Decision 3: Read 函数 nil 判断
**决策：** 在 Read 函数中设置 `encrypted` 字段前，先判断 `fileSystem.Encrypted != nil`，再调用 `d.Set("encrypted", *fileSystem.Encrypted)`。

**理由：**
- 遵循项目规范：在调用 setXX() 前判断 Response 字段是否为 nil
- 避免因云 API 返回 nil 导致 state 写入错误的零值

## Risks / Trade-offs

### Risk 1: 用户修改 encrypted 触发重建
**风险：** 用户修改 `encrypted` 参数会销毁并重建文件系统，可能导致数据丢失。

**缓解措施：**
- 在文档中明确说明 `encrypted` 修改会触发重建
- ForceNew 是 Terraform 标准行为，用户在 plan 阶段可看到重建提示

### Risk 2: 向后兼容性
**风险：** 新增字段影响现有资源。

**缓解措施：**
- `encrypted` 为可选字段，现有配置不设置该字段时行为不变
- 不修改任何现有 schema 字段，完全向后兼容