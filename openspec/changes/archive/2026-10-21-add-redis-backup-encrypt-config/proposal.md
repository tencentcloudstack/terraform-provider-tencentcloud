## Why

当前 Terraform Provider 缺少管理 Redis 实例备份加密配置的资源，用户无法通过 Terraform 自动化开启/关闭备份加密或切换加密密钥类型（腾讯云托管 / 用户自定义 KMS）。云 API 已提供 `DescribeBackupEncryptConfig`（查询）与 `ModifyBackupEncrypt`（修改）两个接口，需要将其封装为 `RESOURCE_KIND_CONFIG` 资源，使该配置可声明式管理。

## What Changes

- 新增资源 `tencentcloud_redis_backup_encrypt`（RESOURCE_KIND_CONFIG），提供 RU（Read/Update）能力：
  - Create：设置 `instance_id` 为 id 后委托派给 Update（首次 apply 即写入配置）。
  - Read：调用 `DescribeBackupEncryptConfig` 回填 `backup_encrypt`、`encrypt_key_type`、`kms_key_id`、`kms_region`。
  - Update：调用 `ModifyBackupEncrypt` 更新加密配置，并在接口成功后回调 Read 刷新 state。
  - Delete：空操作（配置随实例生命周期存在，不支持回收）。
  - 支持 import（import 时以实例 ID 作为 id）。
- `instance_id`：`ForceNew`、必填。
- `backup_encrypt`：必填（bool），是否开启备份加密。
- `encrypt_key_type`：可选（int），加密密钥类型（0=腾讯云托管默认，1=用户自定义 KMS）。
- `kms_key_id`：可选（string），KMS 密钥 ID。
- `kms_region`：可选（string），KMS 地域。
- 在 `tencentcloud/provider.go` 注册资源，在 `tencentcloud/provider.md` 增加资源文档条目。
- 新增对应 `_test.go` 单元测试（gomonkey mock 云 API）与 `.md` 示例文档。

## Capabilities

### New Capabilities
- `redis-backup-encrypt`: 管理 Redis 实例的备份加密配置，包括加密开关、密钥类型与 KMS 密钥的读取与更新。

### Modified Capabilities
<!-- 无既有 spec 需求变更 -->

## Impact

- 新增文件：
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt.go`
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt_test.go`
  - `tencentcloud/services/crs/resource_tc_redis_backup_encrypt.md`
- 修改文件：
  - `tencentcloud/provider.go`：注册 `tencentcloud_redis_backup_encrypt` 资源。
  - `tencentcloud/provider.md`：新增资源文档条目。
- 依赖云 API（已在 vendor 中）：
  - `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412` 的 `DescribeBackupEncryptConfig` / `ModifyBackupEncrypt`。
- 兼容性：纯新增资源，不影响既有资源与 state，向后兼容。