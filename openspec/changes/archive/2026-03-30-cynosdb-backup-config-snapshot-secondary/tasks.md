# Tasks: Add snapshot_secondary_backup_config to tencentcloud_cynosdb_backup_config

## Task 1: 更新 Schema 定义

**文件**：`tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

在 `ResourceTencentCloudCynosdbBackupConfig()` 的 `Schema` map 中，在 `logic_backup_config` 字段之后添加 `snapshot_secondary_backup_config` 字段：

```go
"snapshot_secondary_backup_config": {
    Type:        schema.TypeList,
    Optional:    true,
    Computed:    true,
    MaxItems:    1,
    Description: "Secondary snapshot backup configuration.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "backup_custom_auto_time": {
                Type:        schema.TypeBool,
                Optional:    true,
                Computed:    true,
                Description: "Whether to use system auto time.",
            },
            "backup_time_beg": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "Backup start time. Range: [0-24*3600]. E.g. 0:00, 1:00, 2:00 are 0, 3600, 7200.",
            },
            "backup_time_end": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "Backup end time. Range: [0-24*3600]. E.g. 0:00, 1:00, 2:00 are 0, 3600, 7200.",
            },
            "backup_week_days": {
                Type:        schema.TypeList,
                Optional:    true,
                Computed:    true,
                Description: "Backup week days array (length 7, Sunday to Saturday). Values: full, increment, none.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
            "backup_interval_time": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "Backup interval time.",
            },
            "reserve_duration": {
                Type:        schema.TypeInt,
                Optional:    true,
                Computed:    true,
                Description: "Backup retention period in seconds. 7 days = 604800. Max: 158112000.",
            },
            "backup_trigger_strategy": {
                Type:        schema.TypeString,
                Optional:    true,
                Computed:    true,
                Description: "Backup trigger strategy. Values: periodically (periodic auto backup), frequent (high-frequency backup).",
            },
            "cross_regions_enable": {
                Type:        schema.TypeString,
                Computed:    true,
                Description: "Whether cross-region backup is enabled. Values: yes, no.",
            },
            "cross_regions": {
                Type:        schema.TypeList,
                Computed:    true,
                Description: "Cross-region backup target regions.",
                Elem: &schema.Schema{
                    Type: schema.TypeString,
                },
            },
        },
    },
},
```

- [ ] 在 logic_backup_config 字段之后添加上述 schema 定义
- [ ] 执行 `go fmt tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

## Task 2: 更新 Read 模块

**文件**：`tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

在 `resourceTencentCloudCynosdbBackupConfigRead` 函数中，在 `logic_backup_config` 的 set 逻辑之后，添加 `SnapshotSecondaryBackupConfig` 的读取和 set 逻辑：

```go
if respData.SnapshotSecondaryBackupConfig != nil {
    snapshotSecondaryBackupConfigMap := map[string]interface{}{}
    cfg := respData.SnapshotSecondaryBackupConfig

    if cfg.BackupCustomAutoTime != nil {
        snapshotSecondaryBackupConfigMap["backup_custom_auto_time"] = cfg.BackupCustomAutoTime
    }
    if cfg.BackupTimeBeg != nil {
        snapshotSecondaryBackupConfigMap["backup_time_beg"] = cfg.BackupTimeBeg
    }
    if cfg.BackupTimeEnd != nil {
        snapshotSecondaryBackupConfigMap["backup_time_end"] = cfg.BackupTimeEnd
    }
    if cfg.BackupWeekDays != nil {
        weekDays := make([]string, 0, len(cfg.BackupWeekDays))
        for _, v := range cfg.BackupWeekDays {
            weekDays = append(weekDays, *v)
        }
        snapshotSecondaryBackupConfigMap["backup_week_days"] = weekDays
    }
    if cfg.BackupIntervalTime != nil {
        snapshotSecondaryBackupConfigMap["backup_interval_time"] = cfg.BackupIntervalTime
    }
    if cfg.ReserveDuration != nil {
        snapshotSecondaryBackupConfigMap["reserve_duration"] = cfg.ReserveDuration
    }
    if cfg.BackupTriggerStrategy != nil {
        snapshotSecondaryBackupConfigMap["backup_trigger_strategy"] = cfg.BackupTriggerStrategy
    }
    if cfg.CrossRegionsEnable != nil {
        snapshotSecondaryBackupConfigMap["cross_regions_enable"] = cfg.CrossRegionsEnable
    }
    if cfg.CrossRegions != nil {
        regions := make([]string, 0, len(cfg.CrossRegions))
        for _, v := range cfg.CrossRegions {
            regions = append(regions, *v)
        }
        snapshotSecondaryBackupConfigMap["cross_regions"] = regions
    }

    _ = d.Set("snapshot_secondary_backup_config", []interface{}{snapshotSecondaryBackupConfigMap})
}
```

注意：DescribeBackupConfig 返回的 `SnapshotSecondaryBackupConfig` 类型是 `BackupConfigInfo`，其中含 `CrossRegions` 和 `CrossRegionsEnable` 字段。

- [ ] 在 read 函数中添加上述 SnapshotSecondaryBackupConfig 读取逻辑
- [ ] 执行 `go fmt tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

## Task 3: 更新 Update 模块

**文件**：`tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

**3.1** 将 `snapshot_secondary_backup_config` 加入 `mutableArgs` 列表：

```go
mutableArgs := []string{"backup_time_beg", "backup_time_end", "reserve_duration", "backup_freq", "backup_type", "logic_backup_config", "snapshot_secondary_backup_config"}
```

**3.2** 在 `logic_backup_config` 的构建逻辑之后，添加 `snapshot_secondary_backup_config` 的构建逻辑：

```go
if snapshotSecondaryBackupConfigList, ok := d.GetOk("snapshot_secondary_backup_config"); ok {
    snapshotSecondaryBackupConfigArr := snapshotSecondaryBackupConfigList.([]interface{})
    if len(snapshotSecondaryBackupConfigArr) > 0 {
        cfgMap := snapshotSecondaryBackupConfigArr[0].(map[string]interface{})
        snapshotCfg := &cynosdbv20190107.SnapshotBackupConfig{}

        if v, ok := cfgMap["backup_custom_auto_time"].(bool); ok {
            snapshotCfg.BackupCustomAutoTime = helper.Bool(v)
        }
        if v, ok := cfgMap["backup_time_beg"].(int); ok {
            snapshotCfg.BackupTimeBeg = helper.IntUint64(v)
        }
        if v, ok := cfgMap["backup_time_end"].(int); ok {
            snapshotCfg.BackupTimeEnd = helper.IntUint64(v)
        }
        if v, ok := cfgMap["backup_week_days"].([]interface{}); ok && len(v) > 0 {
            for _, day := range v {
                snapshotCfg.BackupWeekDays = append(snapshotCfg.BackupWeekDays, helper.String(day.(string)))
            }
        }
        if v, ok := cfgMap["backup_interval_time"].(int); ok {
            snapshotCfg.BackupIntervalTime = helper.Int64(int64(v))
        }
        if v, ok := cfgMap["reserve_duration"].(int); ok {
            snapshotCfg.ReserveDuration = helper.IntUint64(v)
        }
        if v, ok := cfgMap["backup_trigger_strategy"].(string); ok && v != "" {
            snapshotCfg.BackupTriggerStrategy = helper.String(v)
        }

        request.SnapshotSecondaryBackupConfig = snapshotCfg
    }
}
```

- [ ] 将 snapshot_secondary_backup_config 加入 mutableArgs
- [ ] 在 update 函数中添加上述构建逻辑
- [ ] 执行 `go fmt tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.go`

## Task 4: 更新 .md 示例文件

**文件**：`tencentcloud/services/cynosdb/resource_tc_cynosdb_backup_config.md`

在现有示例基础上，增加 `snapshot_secondary_backup_config` 用法示例：

```hcl
resource "tencentcloud_cynosdb_backup_config" "example" {
  cluster_id       = "cynosdbmysql-xxxxxx"
  backup_time_beg  = 7200
  backup_time_end  = 21600
  reserve_duration = 604800

  snapshot_secondary_backup_config {
    backup_time_beg          = 7200
    backup_time_end          = 21600
    reserve_duration         = 604800
    backup_trigger_strategy  = "periodically"
  }
}
```

- [ ] 更新 .md 示例文件

## Task 5: 编译验证

- [ ] 执行 `go build ./tencentcloud/services/cynosdb/` 确认编译无误
