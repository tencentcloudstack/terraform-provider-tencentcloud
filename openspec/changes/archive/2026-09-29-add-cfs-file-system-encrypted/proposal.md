## Why

`tencentcloud_cfs_file_system` 资源当前未暴露文件系统加密能力。腾讯云 CFS 服务已在 `CreateCfsFileSystem` 接口支持 `Encrypted` 入参（文件系统是否加密，留空默认不加密），并在 `DescribeCfsFileSystems` 接口的 `FileSystemInfo` 响应中返回 `Encrypted` 字段。由于 CFS 不提供更新加密状态的接口，加密属性只能在创建时指定，用户无法通过 Terraform 在创建文件系统时声明是否加密，也无法在 state 中感知加密状态，存在安全隐患与配置不透明问题。

## What Changes

- 在 `tencentcloud_cfs_file_system` 资源 Schema 中新增 `encrypted` 参数（`TypeBool, Optional, ForceNew: true`）
  - 类型为 bool，可选，由于云 API 不支持更新加密状态，标记为 `ForceNew`，修改时触发重建
- 在 Create 函数中读取 `encrypted` 参数并设置到 `CreateCfsFileSystemRequest.Encrypted`
- 在 Read 函数中将 `FileSystemInfo.Encrypted` 同步到 state
- 更新文档示例（通过 `make doc` 生成 `website/docs/r/cfs_file_system.html.markdown`）及源文档 `tencentcloud/services/cfs/resource_tc_cfs_file_system.md`
- 补充单元测试用例

## Capabilities

### New Capabilities
- `cfs-file-system-resource`: CFS 文件系统资源管理，新增文件系统加密（Encrypted）参数支持，覆盖创建时指定加密状态与读取时同步加密状态

### Modified Capabilities
<!-- 无现有 cfs spec，无需修改 -->

## Impact

### 受影响的代码
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go` - 新增 `encrypted` schema 字段，修改 Create/Read 函数
- `tencentcloud/services/cfs/resource_tc_cfs_file_system_test.go` - 补充单元测试用例
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.md` - 源文档更新
- `website/docs/r/cfs_file_system.html.markdown` - 通过 `make doc` 生成

### 云 API 变更
- `CreateCfsFileSystem` 接口新增入参 `Encrypted`（*bool，已存在于 vendor SDK 中）
- `DescribeCfsFileSystems` 接口返回的 `FileSystemInfo.Encrypted`（*bool，已存在于 vendor SDK 中）

### 向后兼容性
- **完全向后兼容** - `encrypted` 为可选参数，现有配置不受影响
- 不修改已有 schema 字段，仅新增可选字段

### 依赖关系
- 无新增依赖，使用现有 vendor SDK 版本（已包含 `Encrypted` 字段）