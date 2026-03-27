# Change: Fix CLB Instances Backup Zone Set Data Type

## Why

`tencentcloud_clb_instances` 数据源中存在一个数据类型定义错误，导致 `backup_zone_set` 字段的 schema 定义与 SDK 实际返回的数据结构不匹配。

**当前问题**:
1. Schema 定义 (第360-365行):
   ```go
   "backup_zone_set": {
       Type:        schema.TypeList,
       Computed:    true,
       Elem:        &schema.Schema{Type: schema.TypeMap},  // ❌ 错误: 使用了 TypeMap
       Description: "Backup zone list, each element contains zone_id/zone/zone_name/zone_region/local_zone.",
   },
   ```

2. SDK 返回的实际数据结构 (`[]*clb.ZoneInfo`):
   ```go
   type ZoneInfo struct {
       ZoneId     *int64  `json:"ZoneId"`
       Zone       *string `json:"Zone"`
       ZoneName   *string `json:"ZoneName"`
       ZoneRegion *string `json:"ZoneRegion"`
       LocalZone  *bool   `json:"LocalZone"`
   }
   ```

3. 数据处理代码 (第623-645行) 正确地将 SDK 数据转换为 `[]map[string]interface{}`，但与 schema 定义冲突

**影响**:
- Schema 类型不正确可能导致 Terraform 状态管理问题
- 字段类型校验不严格，不符合 Terraform 最佳实践
- 与腾讯云 API 文档 (https://cloud.tencent.com/document/api/214/30685) 中的 `BackupZoneSet` 字段定义不一致

## What Changes

修复 `tencentcloud_clb_instances` 数据源中 `backup_zone_set` 字段的类型定义:

1. **更新 Schema 定义**:
   - 将 `backup_zone_set` 的 `Elem` 从 `&schema.Schema{Type: schema.TypeMap}` 改为使用 `&schema.Resource` 定义复杂对象
   - 明确定义每个子字段的类型: `zone_id` (TypeInt), `zone` (TypeString), `zone_name` (TypeString), `zone_region` (TypeString), `local_zone` (TypeBool)

2. **保持数据处理逻辑**:
   - 第623-645行的数据转换逻辑已经正确，无需修改
   - 严格按照 SDK 的 `ZoneInfo` 结构进行字段映射和类型转换

3. **代码格式化**:
   - 修改完成后执行 `go fmt` 对代码进行格式化

## Impact

- **影响范围**: `tencentcloud_clb_instances` 数据源
- **破坏性**: 无破坏性变更 (仅修正类型定义，数据结构保持一致)
- **受影响文件**:
  - `tencentcloud/services/clb/data_source_tc_clb_instances.go`
  - `tencentcloud/services/clb/data_source_tc_clb_instances.md` (可能需要更新文档)
- **API 依赖**: 
  - CLB API v20180317: `DescribeLoadBalancers`
  - 文档: https://cloud.tencent.com/document/api/214/30685
  - SDK 字段: `BackupZoneSet []*ZoneInfo`
- **兼容性**: 完全向后兼容，仅修正类型定义
- **测试要求**: 验证读取 CLB 实例时 `backup_zone_set` 字段能正确解析

## Technical Details

### 修改前 (错误的定义):
```go
"backup_zone_set": {
    Type:        schema.TypeList,
    Computed:    true,
    Elem:        &schema.Schema{Type: schema.TypeMap},
    Description: "Backup zone list, each element contains zone_id/zone/zone_name/zone_region/local_zone.",
},
```

### 修改后 (正确的定义):
```go
"backup_zone_set": {
    Type:        schema.TypeList,
    Computed:    true,
    Description: "Backup zone list, each element contains zone_id/zone/zone_name/zone_region/local_zone.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "zone_id": {
                Type:        schema.TypeInt,
                Computed:    true,
                Description: "Availability zone ID (numerical representation).",
            },
            "zone": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Availability zone unique identifier (string representation).",
            },
            "zone_name": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Availability zone name.",
            },
            "zone_region": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Region that this availability zone belongs to.",
            },
            "local_zone": {
                Type:        schema.TypeBool,
                Computed:    true,
                Description: "Whether this is a local zone.",
            },
        },
    },
},
```

### SDK 数据结构参考:
```go
// From tencentcloud-sdk-go/tencentcloud/clb/v20180317/models.go
type ZoneInfo struct {
    ZoneId     *int64  `json:"ZoneId,omitempty" name:"ZoneId"`
    Zone       *string `json:"Zone,omitempty" name:"Zone"`
    ZoneName   *string `json:"ZoneName,omitempty" name:"ZoneName"`
    ZoneRegion *string `json:"ZoneRegion,omitempty" name:"ZoneRegion"`
    LocalZone  *bool   `json:"LocalZone,omitempty" name:"LocalZone"`
}
```

### 数据处理逻辑 (无需修改):
```go
if clbInstance.BackupZoneSet != nil {
    backupZones := make([]map[string]interface{}, 0, len(clbInstance.BackupZoneSet))
    for _, zone := range clbInstance.BackupZoneSet {
        backupZone := make(map[string]interface{})
        if zone.ZoneId != nil {
            backupZone["zone_id"] = *zone.ZoneId  // int64 -> TypeInt
        }
        if zone.Zone != nil {
            backupZone["zone"] = *zone.Zone  // string -> TypeString
        }
        if zone.ZoneName != nil {
            backupZone["zone_name"] = *zone.ZoneName  // string -> TypeString
        }
        if zone.ZoneRegion != nil {
            backupZone["zone_region"] = *zone.ZoneRegion  // string -> TypeString
        }
        if zone.LocalZone != nil {
            backupZone["local_zone"] = *zone.LocalZone  // bool -> TypeBool
        }
        backupZones = append(backupZones, backupZone)
    }
    mapping["backup_zone_set"] = backupZones
}
```
