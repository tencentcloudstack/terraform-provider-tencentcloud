# 设计文档：tencentcloud_cls_dlc_deliver 资源

## 1. 文件结构

```
tencentcloud/services/cls/
├── resource_tc_cls_dlc_deliver.go          # 新增：资源主文件
├── resource_tc_cls_dlc_deliver.md          # 新增：资源文档示例
├── resource_tc_cls_dlc_deliver_test.go     # 新增：验收测试
└── service_tencentcloud_cls.go             # 修改：新增 DescribeClsDlcDeliverById

tencentcloud/
└── provider.go                             # 修改：注册资源
```

---

## 2. Schema 设计

```go
func ResourceTencentCloudClsDlcDeliver() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudClsDlcDeliverCreate,
        Read:   resourceTencentCloudClsDlcDeliverRead,
        Delete: resourceTencentCloudClsDlcDeliverDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{
            // Required
            "topic_id":     TypeString, Required, ForceNew
            "name":         TypeString, Required, ForceNew
            "deliver_type": TypeInt,    Required, ForceNew
            "start_time":   TypeInt,    Required, ForceNew
            "dlc_info":     TypeList MaxItems=1, Required, ForceNew
                ├── "table_info": TypeList MaxItems=1, Required, ForceNew
                │   ├── "data_directory": TypeString, Required, ForceNew
                │   ├── "database_name":  TypeString, Required, ForceNew
                │   └── "table_name":     TypeString, Required, ForceNew
                ├── "field_infos": TypeList, Optional, ForceNew
                │   ├── "cls_field":      TypeString, Required, ForceNew
                │   ├── "dlc_field":      TypeString, Required, ForceNew
                │   ├── "dlc_field_type": TypeString, Required, ForceNew
                │   ├── "fill_field":     TypeString, Optional, ForceNew
                │   └── "disable":        TypeBool,   Optional, ForceNew
                ├── "partition_infos": TypeList, Optional, ForceNew
                │   ├── "cls_field":      TypeString, Required, ForceNew
                │   ├── "dlc_field":      TypeString, Required, ForceNew
                │   └── "dlc_field_type": TypeString, Required, ForceNew
                └── "partition_extra": TypeList MaxItems=1, Optional, ForceNew
                    ├── "time_format": TypeString, Optional, ForceNew
                    └── "time_zone":   TypeString, Optional, ForceNew

            // Optional
            "max_size":        TypeInt, Optional, ForceNew
            "interval":        TypeInt, Optional, ForceNew
            "end_time":        TypeInt, Optional, ForceNew
            "has_services_log": TypeInt, Optional, ForceNew

            // Computed
            "task_id": TypeString, Computed
        },
    }
}
```

**所有字段均设为 ForceNew**，因为没有 Update 接口，任何修改都需要删除重建。

---

## 3. CRUD 函数设计

### Create

```
resourceTencentCloudClsDlcDeliverCreate:
1. 构建 CreateDlcDeliverRequest，填充所有字段
2. dlc_info 嵌套结构需递归构建：
   - DlcInfo.TableInfo = &cls.DlcTableInfo{...}
   - DlcInfo.FieldInfos = []*cls.DlcFiledInfo{...}
   - DlcInfo.PartitionInfos = []*cls.DlcPartitionInfo{...}
   - DlcInfo.PartitionExtra = &cls.DlcPartitionExtra{...}
3. resource.Retry(WriteRetryTimeout) 调用 CreateDlcDeliverWithContext
4. response.Response.TaskId → taskId
5. d.SetId(topicId + "#" + taskId)
6. 调用 Read
```

**注意**：`StartTime`、`EndTime` 等字段类型在 SDK 中是 `*int64`（CreateDlcDeliverRequest），
而在 `DlcDeliverInfo` 返回中是 `*uint64`，读写时需分别处理类型转换：
- Write: `helper.IntInt64(v.(int))` 或 `helper.IntUint64(v.(int))`
- Read: `int(*respData.StartTime)`

### Read

```
resourceTencentCloudClsDlcDeliverRead:
1. idSplit := strings.Split(d.Id(), "#")，解析 topicId, taskId
2. 调用 service.DescribeClsDlcDeliverById(ctx, topicId, taskId)
3. 若 nil → d.SetId(""), return nil
4. 逐字段 d.Set(...)
5. dlc_info 嵌套结构：构建 map 后 d.Set("dlc_info", ...)
```

### Delete

```
resourceTencentCloudClsDlcDeliverDelete:
1. idSplit := strings.Split(d.Id(), "#")
2. 构建 DeleteDlcDeliverRequest{TopicId, TaskId}
3. resource.Retry(WriteRetryTimeout) 调用 DeleteDlcDeliverWithContext
```

---

## 4. Service 层设计

```go
func (me *ClsService) DescribeClsDlcDeliverById(
    ctx context.Context, topicId, taskId string,
) (ret *cls.DlcDeliverInfo, errRet error) {
    request := cls.NewDescribeDlcDeliversRequest()
    request.TopicId = &topicId
    request.Filters = []*cls.Filter{
        {Key: helper.String("taskId"), Values: []*string{&taskId}},
    }

    err := resource.Retry(ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseClsClient().DescribeDlcDelivers(request)
        // ...
    })

    if len(response.Response.Infos) == 0 { return }
    ret = response.Response.Infos[0]
    return
}
```

---

## 5. Provider 注册

在 `tencentcloud/provider.go` 的 `ResourcesMap` 中，紧跟 `tencentcloud_cls_scheduled_sql` 后追加：

```go
"tencentcloud_cls_dlc_deliver": cls.ResourceTencentCloudClsDlcDeliver(),
```

---

## 6. ID 处理

```go
// Create
d.SetId(topicId + "#" + *response.Response.TaskId)

// Read / Delete
idSplit := strings.Split(d.Id(), "#")
if len(idSplit) != 2 {
    return fmt.Errorf("id is broken, %s", d.Id())
}
topicId := idSplit[0]
taskId  := idSplit[1]
```

---

## 7. SDK 关键类型对应

| Schema 字段 | SDK Request 类型 | SDK Response 类型 |
|-------------|-----------------|------------------|
| `deliver_type` | `*int64` | `*uint64` |
| `max_size` | `*int64` | `*uint64` |
| `interval` | `*int64` | `*uint64` |
| `start_time` | `*int64` | `*uint64` |
| `end_time` | `*int64` | `*uint64` |
| `has_services_log` | `*int64` | `*uint64` |

Write 时用 `helper.IntInt64(v.(int))`，Read 时用 `int(*respData.Field)`。

---

## 8. 测试设计

测试函数：`TestAccTencentCloudClsDlcDeliverResource_basic`

Steps：
1. 创建：最小参数（topic_id + name + deliver_type=1 + start_time + dlc_info）
2. ImportState 验证

---

## 9. 文档示例设计（resource_tc_cls_dlc_deliver.md）

```hcl
resource "tencentcloud_cls_dlc_deliver" "example" {
  topic_id     = "715094e3-01b0-4aeb-91f5-ee9f46a4a13c"
  name         = "tf-example-dlc-deliver"
  deliver_type = 1
  start_time   = 1741005340

  dlc_info {
    table_info {
      data_directory = "test_data"
      database_name  = "test_db"
      table_name     = "test_table"
    }

    field_infos {
      cls_field      = "cls_field1"
      dlc_field      = "dlc_field1"
      dlc_field_type = "string"
      disable        = false
    }

    partition_infos {
      cls_field      = "cls_time"
      dlc_field      = "dt"
      dlc_field_type = "string"
    }

    partition_extra {
      time_format = "/%Y/%m/%d/%H"
      time_zone   = "UTC+08:00"
    }
  }
}
```
