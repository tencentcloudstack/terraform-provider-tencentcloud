# Change: 支持 CFS 文件系统 capacity 参数修改

## Why

目前 `tencentcloud_cfs_file_system` 资源的 `capacity` 参数被硬编码为不可修改（在 Update 函数的 immutableArgs 列表中，第 312 行），用户无法通过 Terraform 对 Turbo 系列文件系统进行扩容。这导致：

1. **无法满足业务增长需求**：当文件系统存储容量不足时，用户必须手动通过控制台或 API 扩容，无法通过 IaC 统一管理
2. **配置漂移问题**：手动扩容后，Terraform 状态与实际容量不一致，需要手动同步
3. **运维复杂度增加**：破坏了基础设施即代码的一致性原则

根据腾讯云官方文档：
- **ScaleUpFileSystem API**: https://cloud.tencent.com/document/product/582/90456

该 API 支持对 Turbo 系列文件系统进行扩容（仅扩容，不支持缩容），具有以下特性：
- Turbo 标准型（TB）：扩容步长为 20,480 GiB（20 TiB）
- Turbo 性能型（TP）：扩容步长为 10,240 GiB（10 TiB）
- 操作为异步，需要等待文件系统状态变为 `available`

**业务场景：**
- 用户可以通过修改 `capacity` 参数在线扩容 Turbo 文件系统
- 扩容操作通过 Terraform 完成，保持状态一致性
- 支持自动化扩容策略（如结合监控告警）

## What Changes

在 `tencentcloud_cfs_file_system` 资源中支持 `capacity` 参数的修改：

### Schema 变更
- **移除限制**：从 `immutableArgs` 列表中移除 `capacity`
- **类型不变**：保持为 `TypeInt, Optional, Computed`
- **仅支持扩容**：通过验证逻辑确保新值大于旧值

### Update 逻辑
- 检测 `capacity` 变更
- 调用 `ScaleUpFileSystem` API
- 等待扩容完成（状态从 `updating` 变为 `available`）
- 添加 Update timeout 支持（默认 30 分钟）

### 服务层新增
在 `service_tencentcloud_cfs.go` 中添加：
```go
func (me *CfsService) ScaleUpFileSystem(ctx context.Context, fsId string, targetCapacity int) error
```

### 文档更新
- `website/docs/r/cfs_file_system.html.markdown` - 更新 `capacity` 参数说明，标注支持修改（仅扩容）
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.md` - 同步更新
- 添加 Timeouts 说明中的 `update` 参数

### 变更日志
- `.changelog/<PR_NUMBER>.txt` - 记录 enhancement

## Impact

### 受影响的规范
- 新增规范：`cfs-file-system-resource` - CFS 文件系统资源管理（capacity 更新能力）

### 受影响的代码
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.go` - Update 逻辑，移除 immutableArgs 限制，添加扩容逻辑
- `tencentcloud/services/cfs/service_tencentcloud_cfs.go` - 新增 ScaleUpFileSystem 方法
- `website/docs/r/cfs_file_system.html.markdown` - 文档更新
- `tencentcloud/services/cfs/resource_tc_cfs_file_system.md` - 源文档更新

### 向后兼容性
- ✅ **完全向后兼容** - 现有不修改 capacity 的配置不受影响
- ✅ **不影响现有资源** - 未修改 capacity 时行为与之前完全一致
- ⚠️ **行为变更** - 修改 capacity 从触发重建变为原地更新（破坏性变更，但符合用户预期）

### API 兼容性
- ✅ `ScaleUpFileSystem` API 支持扩容操作（仅 Turbo 系列）
- ✅ 标准型/高性能型文件系统不支持扩容，API 会返回错误
- ✅ 仅支持扩容，缩容会被拒绝

### 验证逻辑
- 新容量必须大于当前容量（不允许缩容或相等）
- 仅 Turbo 系列（storage_type 为 TB 或 TP）支持扩容
- 扩容步长验证：
  - TB: 必须是 20,480 的倍数
  - TP: 必须是 10,240 的倍数

### 测试影响
- 需要验收测试覆盖 Turbo 文件系统的扩容场景
- 测试扩容步长验证逻辑
- 测试状态等待逻辑（创建 → 扩容 → available）

### 依赖关系
- 无新增依赖
- 使用现有的 SDK 版本（已包含 ScaleUpFileSystem API）
- Timeouts 块已在之前的变更中添加，需要扩展支持 Update timeout
