## Context

当前 TencentCloud Terraform Provider 在 `tencentcloud/services/crs/` 下已有多个 redis 配置类资源（如 `redis_backup_config`、`redis_instance_password_policy_config`），它们遵循 RESOURCE_KIND_CONFIG 的标准模式：以 `instance_id` 为 id，Create 委派 Update，Read 调用查询接口回填，Update 调用修改接口，Delete 为空操作。

本次新增 `tencentcloud_redis_backup_encrypt` 资源，封装云 API `DescribeBackupEncryptConfig`（查询）与 `ModifyBackupEncrypt`（修改），管理 Redis 实例的备份加密配置。相关 SDK struct 已在 vendor 中存在。

### 云 API 关键 struct 定义（摘自 vendor，按行区间读取，供实施阶段直接参考）

#### DescribeBackupEncryptConfig（查询备份加密配置）

Request（`vendor/.../redis/v20180412/models.go`）:
```go
type DescribeBackupEncryptConfigRequestParams struct {
    // 实例 ID，请登录Redis控制台在实例列表复制实例 ID。
    InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
}
```

Response（`vendor/.../redis/v20180412/models.go`）:
```go
type DescribeBackupEncryptConfigResponseParams struct {
    // 备份是否加密。枚举值：true 加密；false 不加密
    BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
    // 加密密钥类型。枚举值：0 腾讯云托管（默认）；1 用户自定义 KMS 密钥
    EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
    // KMS 的密钥 ID。参数格式：xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx
    KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
    // KMS 地域信息。参数格式：ap-xxxx
    KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
    // 唯一请求 ID
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

Client 调用（`vendor/.../redis/v20180412/client.go`）:
```go
func (c *Client) DescribeBackupEncryptConfigWithContext(ctx context.Context, request *DescribeBackupEncryptConfigRequest) (response *DescribeBackupEncryptConfigResponse, err error)
```
可能错误码：`FailedOperation.SystemError`、`InternalError.DbOperationFailed`、`InternalError.InternalError`、`InvalidParameter.PermissionDenied`、`UnauthorizedOperation.NoCAMAuthed`。

#### ModifyBackupEncrypt（修改备份加密配置）

Request（`vendor/.../redis/v20180412/models.go`）:
```go
type ModifyBackupEncryptRequestParams struct {
    // 实例 ID
    InstanceId *string `json:"InstanceId,omitnil,omitempty" name:"InstanceId"`
    // 是否开启备份加密。默认值：false
    BackupEncrypt *bool `json:"BackupEncrypt,omitnil,omitempty" name:"BackupEncrypt"`
    // 加密密钥类型。枚举值：0 腾讯云托管（默认）；1 用户自定义 KMS 密钥。默认值：0
    EncryptKeyType *int64 `json:"EncryptKeyType,omitnil,omitempty" name:"EncryptKeyType"`
    // KMS 的密钥 ID。参数格式：xxxxxxx-xxxx-xxxx-xxxx-xxxxxxx
    KmsKeyId *string `json:"KmsKeyId,omitnil,omitempty" name:"KmsKeyId"`
    // 地域信息。参数格式：ap-xxxx（KMS 支持的地域列表）
    KmsRegion *string `json:"KmsRegion,omitnil,omitempty" name:"KmsRegion"`
}
```

Response（`vendor/.../redis/v20180412/models.go`）:
```go
type ModifyBackupEncryptResponseParams struct {
    // 唯一请求 ID
    RequestId *string `json:"RequestId,omitnil,omitempty" name:"RequestId"`
}
```

Client 调用（`vendor/.../redis/v20180412/client.go`）:
```go
func (c *Client) ModifyBackupEncryptWithContext(ctx context.Context, request *ModifyBackupEncryptRequest) (response *ModifyBackupEncryptResponse, err error)
```

### 参数映射（JsonPath → SchemaName）

| 云API路径 | Schema 字段 | 必填 | 类型 | CRUD 可用 |
|---|---|---|---|---|
| `request.InstanceId` / `response.Response.InstanceId` | `instance_id` | 是 | string | Create(R)/Read/U/Delete |
| `response.Response.BackupEncrypt` / `request.BackupEncrypt` | `backup_encrypt` | 是 | bool | Read / Update |
| `response.Response.EncryptKeyType` / `request.EncryptKeyType` | `encrypt_key_type` | 否 | int | Read / Update |
| `response.Response.KmsKeyId` / `request.KmsKeyId` | `kms_key_id` | 否 | string | Read / Update |
| `response.Response.KmsRegion` / `request.KmsRegion` | `kms_region` | 否 | string | Read / Update |

注意：`DescribeBackupEncryptConfig` 仅返回配置字段，不返回 `InstanceId`；`instance_id` 在 Read 中直接由 `d.Id()` 设置。

### 参考资源结构（代码风格参照）

参照 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config.go`（同包、同 RESOURCE_KIND_CONFIG 模式、同使用 `UseRedisClient()`）：

- Schema：`instance_id`（Required + ForceNew）、其余字段按需求设置 Required/Optional。
- `ResourceTencentCloudRedisBackupEncrypt()` 返回 `*schema.Resource`，声明 Create/Read/Update/Delete + Importer（`schema.ImportStatePassthrough`）。
- Create：`d.SetId(instance_id)` 后 `return resourceTencentCloudRedisBackupEncryptUpdate(d, meta)`。
- Read：`request.InstanceId = helper.String(d.Id())`，`resource.Retry(tccommon.ReadRetryTimeout, ...)` 调用 `UseRedisClient().DescribeBackupEncryptConfigWithContext(ctx, request)`；对 `ResourceNotFound.*` 等实例不存在错误码 `d.SetId("")`（先 `log.Printf("[CRUD] redis_backup_encrypt id=%s", d.Id())` 再清空）；返回空配置时同样处理；逐字段 nil 判断后 `d.Set`。
- Update：构造 `ModifyBackupEncryptRequest`，`resource.Retry(tccommon.WriteRetryTimeout, ...)` 调用 `ModifyBackupEncryptWithContext`，成功后 `return resourceTencentCloudRedisBackupEncryptRead(d, meta)`。
- Delete：空操作（备份加密配置随实例生命周期，无独立删除语义）。
- 客户端获取：`meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseRedisClient()`。
- 导入：`redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"`。

参照 `tencentcloud/services/crs/resource_tc_redis_instance_password_policy_config_test.go` 的测试风格（gomonkey mock `UseRedisClient` 与 `*WithContext` 方法，断言 schema/Read/Create/Update/Delete）。

参照 `tencentcloud/services/crs/resource_tc_redis_backup_config.md` 的文档风格（一句话描述带云产品名 + Example Usage + Import）。

### provider 注册位置

- `tencentcloud/provider.go`：在 redis 资源块（约 1703-1711 行附近）新增一行 `"tencentcloud_redis_backup_encrypt": crs.ResourceTencentCloudRedisBackupEncrypt(),`。
- `tencentcloud/provider.md`：在 redis Resource 列表（约 1025 行起）新增一行 `tencentcloud_redis_backup_encrypt`。

## Goals / Non-Goals

**Goals:**
- 提供 `tencentcloud_redis_backup_encrypt` 资源，声明式管理 Redis 实例备份加密配置（开关、密钥类型、KMS 密钥 ID、KMS 地域）。
- 支持 import（以实例 ID 作为 id）。
- 单元测试覆盖 schema、Read、Create、Update、Delete 路径（gomonkey mock）。
- 在 provider 注册并提供文档示例。

**Non-Goals:**
- 不管理 Redis 实例本身的创建/销毁（由 `tencentcloud_redis_instance` 负责）。
- 不实现真正的删除语义（配置随实例生命周期，Delete 为空操作）。
- 不新增 `_extension.go`。
- 不修改既有 redis 资源的 schema 或行为。

## Decisions

### 1. id 采用单一 `instance_id`（非复合 id）
- **决策**：`d.SetId(instance_id)`，Read 时 `d.Id()` 即实例 ID。
- **理由**：备份加密配置是实例级别的单一配置，一个实例只有一份配置，无第二主键。与同包 `redis_backup_config`、`redis_instance_password_policy_config` 一致。
- **替代方案**：复合 id（如 `instanceId#backup_encrypt`）——无必要，引入解析复杂度且无第二个稳定标识。

### 2. Create 委派 Update（首次 apply 即写入）
- **决策**：Create 中 `d.SetId(instance_id)` 后直接 `return ...Update(d, meta)`。
- **理由**：CONFIG 资源无独立创建 API；配置"存在即生效"，首次 apply 通过 `ModifyBackupEncrypt` 写入。与既有 CONFIG 资源一致。

### 3. Delete 为空操作
- **决策**：Delete 直接 `return nil`，不清空 id、不调用 API。
- **理由**：云 API 未提供"删除备份加密配置"接口；且关闭加密应通过 `backup_encrypt = false` 的 Update 完成，而非 destroy。与 `redis_instance_password_policy_config` 的 Delete 一致。
- **替代方案**：destroy 时自动调用 `ModifyBackupEncrypt(BackupEncrypt=false)`——会改变用户预期（destroy 不应隐式修改云端配置），故不采用。

### 4. Read 处理实例不存在
- **决策**：当查询接口返回实例不存在错误码或返回空配置时，先 `log.Printf("[CRUD] redis_backup_encrypt id=%s", d.Id())` 再 `d.SetId("")`。
- **理由**：遵循仓库规则第 8 条，保留排障现场信息；避免丢失 id 后无法定位。

### 5. 异步处理判定
- **决策**：本资源不涉及异步轮询。
- **理由**：`ModifyBackupEncrypt` 返回仅含 RequestId，云 API 文档未标注为异步接口；Read 接口可直接读到最新配置。Update 成功后直接回调 Read 刷新即可，无需轮询。

### 6. 字段类型与默认值
- `backup_encrypt`：`schema.TypeBool`，Required（控制加密开关，必须显式声明）。
- `encrypt_key_type`：`schema.TypeInt`，Optional（云 API 默认 0）。
- `kms_key_id`：`schema.TypeString`，Optional。
- `kms_region`：`schema.TypeString`，Optional。
- Update 中对 Optional 字段使用 `d.GetOkExists`（针对可能为零值的 int/string）或 `d.GetOk`，参照 password_policy_config 对 bool 用 `GetOk`、对 int 用 `GetOkExists` 的做法。

## Risks / Trade-offs

- [配置与实例生命周期耦合] → Delete 为空操作，用户 destroy 资源不会关闭云端加密；需在文档说明关闭加密应设 `backup_encrypt = false`。
- [Read 返回空配置与实例不存在的区分] → 两种情况均 `d.SetId("")` 并记录 `[CRUD]` 日志；云 API 未提供更细的错误码区分，可接受。
- [KMS 密钥有效性依赖外部资源] → `kms_key_id`/`kms_region` 的有效性由云 API 校验，terraform 仅透传；Update 失败时由 `tccommon.RetryError` 包装返回，用户可见错误信息。
- [向后兼容] → 纯新增资源，无破坏性变更。
