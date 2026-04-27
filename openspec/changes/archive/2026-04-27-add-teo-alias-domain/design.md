# Design: tencentcloud_teo_alias_domain

## Schema Design

### 资源 ID

```
{zone_id}#{alias_name}
```

`zone_id` 和 `alias_name` 均为 `ForceNew`，变更时需要销毁重建。

### Schema 字段（与 CreateAliasDomain 接口入参严格对齐）

| Terraform 字段 | 类型 | 必填 | ForceNew | Computed | 对应 API 字段 | 说明 |
|---|---|---|---|---|---|---|
| `zone_id` | String | ✅ | ✅ | ❌ | ZoneId | 站点 ID |
| `alias_name` | String | ✅ | ✅ | ❌ | AliasName | 别称域名名称 |
| `target_name` | String | ✅ | ❌ | ❌ | TargetName | 目标域名名称 |
| `cert_type` | String | ❌ | ❌ | ❌ | CertType | 证书配置：none/hosting，默认 none |
| `cert_id` | List(String) | ❌ | ❌ | ❌ | CertId | 证书 ID 列表，CertType=hosting 时填写 |
| `paused` | Bool | ❌ | ❌ | ❌ | （ModifyAliasDomainStatus.Paused） | 是否停用，false=启用，true=停用 |
| `status` | String | ❌ | ❌ | ✅ | Status | 别称域名状态（只读，API 返回） |
| `forbid_mode` | Int | ❌ | ❌ | ✅ | ForbidMode | 封禁模式（只读） |
| `created_on` | String | ❌ | ❌ | ✅ | CreatedOn | 创建时间（只读） |
| `modified_on` | String | ❌ | ❌ | ✅ | ModifiedOn | 修改时间（只读） |

## CRUD 实现设计

### Create

调用 `CreateAliasDomain`，入参：
- `ZoneId`（Required）
- `AliasName`（Required）
- `TargetName`（Required）
- `CertType`（Optional，默认 none）
- `CertId`（Optional，CertType=hosting 时）

成功后若 `paused=true`，继续调用 `ModifyAliasDomainStatus` 将其置为停用状态。

资源 ID 设为 `{zone_id}#{alias_name}`，调用 Read 刷新状态。

### Read

调用 `DescribeAliasDomains`，通过 Filter `alias-name` 精确匹配，遍历结果按 `ZoneId+AliasName` 确认匹配。

在 service 层封装 `DescribeTeoAliasDomainById(ctx, zoneId, aliasName)` 方法，返回 `*teo.AliasDomain`。

注意：证书相关字段（`cert_type`、`cert_id`）API 的查询接口（`DescribeAliasDomains`）返回值中**不包含**这两个字段（AliasDomain 结构体无此字段），Read 时不覆盖这两个字段（保留 state 中的值，防止漂移）。

### Update

检测变更字段：
- `target_name` / `cert_type` / `cert_id` 有变更 → 调用 `ModifyAliasDomain`
- `paused` 有变更 → 调用 `ModifyAliasDomainStatus`

两者均可在同一次 Update 中触发。

### Delete

调用 `DeleteAliasDomain`，入参 `ZoneId` + `AliasNames: [aliasName]`。

## Service 层

在 `service_tencentcloud_teo.go` 中追加：

```go
func (me *TeoService) DescribeTeoAliasDomainById(ctx context.Context, zoneId, aliasName string) (aliasDomain *teo.AliasDomain, errRet error)
```

使用 `DescribeAliasDomains` 接口，Filter `alias-name` 精确查找，分页循环直到找到或耗尽。

## 代码风格参考

严格参考 `resource_tc_igtm_strategy.go` 风格：
- `defer tccommon.LogElapsed(...)()` + `defer tccommon.InconsistentCheck(...)()`
- `resource.Retry(tccommon.WriteRetryTimeout, ...)` 包裹写操作
- `helper.String/IntInt64/IntUint64` 等 helper 函数
- Read 中资源不存在时 `d.SetId("")` 返回 nil
- ID 解析使用 `strings.Split(d.Id(), tccommon.FILED_SP)`

## 文件结构

```
tencentcloud/services/teo/
├── resource_tc_teo_alias_domain.go        # 新增
├── resource_tc_teo_alias_domain.md        # 新增
├── resource_tc_teo_alias_domain_test.go   # 新增
└── service_tencentcloud_teo.go            # 追加 DescribeTeoAliasDomainById
tencentcloud/
└── provider.go                            # 注册新资源
```
