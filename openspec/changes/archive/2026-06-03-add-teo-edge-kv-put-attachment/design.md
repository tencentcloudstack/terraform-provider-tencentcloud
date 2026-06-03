## Context

TencentCloud EdgeOne (TEO) 提供 Edge KV 存储服务，允许在边缘节点存储键值对数据。SDK 中已有 `EdgeKVPut`、`EdgeKVGet`、`EdgeKVDelete` 三个接口，分别用于写入、查询和删除 KV 数据。当前 Provider 中已有多个 TEO 资源（如 `tencentcloud_teo_zone`、`tencentcloud_teo_dns_record` 等），本次新增资源遵循相同的代码组织模式。

本资源为 RESOURCE_KIND_ATTACHMENT 类型，表示 KV 数据与命名空间的绑定关系，只有 CRD（Create/Read/Delete）操作，无 Update。

## Goals / Non-Goals

**Goals:**
- 实现 `tencentcloud_teo_edge_k_v_put` 资源，支持通过 Terraform 管理 Edge KV 数据的写入和删除
- 使用联合 ID（zone_id + namespace + key）唯一标识资源
- 支持 import 操作
- 提供完整的单元测试覆盖

**Non-Goals:**
- 不实现 Update 操作（ATTACHMENT 类型资源无需 Update）
- 不实现批量写入/删除（每个资源实例管理单个 KV 对）
- 不新增 service 层方法（直接在资源文件中调用 SDK）

## Decisions

### 1. 资源 ID 设计

**决策**: 使用 `zone_id#namespace#key` 作为联合 ID（`#` 为 `tccommon.FILED_SP`）

**理由**: EdgeKVGet 和 EdgeKVDelete 接口需要 ZoneId、Namespace、Keys 三个参数才能定位到具体的 KV 数据，因此需要三者组合作为唯一标识。这与 Provider 中其他使用联合 ID 的资源保持一致。

### 2. CRD 模式与 immutableArgs 保护

**决策**: 
- `zone_id`、`namespace`、`key` 设置为 ForceNew（构成 ID 的字段）
- `value`、`expiration`、`expiration_ttl` 在 Update 方法中通过 immutableArgs 数组检查，若发生变更则返回错误

**理由**: RESOURCE_KIND_ATTACHMENT 类型资源只有 CRD 接口。ID 字段使用 ForceNew 确保变更时重建资源，其余字段通过 immutableArgs 机制阻止原地更新。

### 3. Read 方法实现

**决策**: 调用 EdgeKVGet 接口时，Keys 参数传入包含单个 key 的数组 `[]*string{&key}`

**理由**: EdgeKVGet 接口设计为批量查询，但本资源每次只管理单个 KV 对。返回的 Data 数组中取第一个元素的 Value 即可。如果返回的 Data 为空或 Value 为空字符串，说明 key 不存在，应调用 `d.SetId("")` 标记资源已被外部删除。

### 4. 不生成 _extension.go 文件

**决策**: 所有逻辑直接写在 `resource_tc_teo_edge_k_v_put_attachment.go` 中

**理由**: 资源逻辑简单，无需拆分文件。

## Risks / Trade-offs

- [Risk] EdgeKVGet 返回空 Value 时无法区分"key 不存在"和"value 为空字符串" → 由于 EdgeKVPut 要求 Value 不能为空，因此空 Value 可以安全地视为 key 不存在
- [Risk] expiration 和 expiration_ttl 为创建时设置的过期参数，Read 接口返回的是 ISO 8601 格式的绝对过期时间 → Read 方法中不回写 expiration 和 expiration_ttl 字段（这两个字段仅在创建时使用），只回写从 Data 中获取的 value
