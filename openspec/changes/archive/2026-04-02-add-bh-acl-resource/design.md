# 设计文档：tencentcloud_bh_acl 资源

## 1. 文件结构

```
tencentcloud/services/bh/
├── resource_tc_bh_acl.go          # 新增：资源主文件
├── resource_tc_bh_acl.md          # 新增：资源文档示例
├── resource_tc_bh_acl_test.go     # 新增：验收测试
└── service_tencentcloud_bh.go     # 修改：新增 DescribeBhAclById 方法

tencentcloud/
└── provider.go                    # 修改：注册资源
```

---

## 2. Schema 设计

### resource_tc_bh_acl.go — Schema 定义

```go
func ResourceTencentCloudBhAcl() *schema.Resource {
    return &schema.Resource{
        Create: resourceTencentCloudBhAclCreate,
        Read:   resourceTencentCloudBhAclRead,
        Update: resourceTencentCloudBhAclUpdate,
        Delete: resourceTencentCloudBhAclDelete,
        Importer: &schema.ResourceImporter{
            State: schema.ImportStatePassthrough,
        },
        Schema: map[string]*schema.Schema{
            // Required
            "name":               TypeString, Required
            "allow_disk_redirect": TypeBool,   Required
            "allow_any_account":   TypeBool,   Required

            // Optional — 权限开关
            "allow_clip_file_up":    TypeBool, Optional
            "allow_clip_file_down":  TypeBool, Optional
            "allow_clip_text_up":    TypeBool, Optional
            "allow_clip_text_down":  TypeBool, Optional
            "allow_file_up":         TypeBool, Optional
            "allow_file_down":       TypeBool, Optional
            "allow_disk_file_up":    TypeBool, Optional
            "allow_disk_file_down":  TypeBool, Optional
            "allow_shell_file_up":   TypeBool, Optional
            "allow_shell_file_down": TypeBool, Optional
            "allow_file_del":        TypeBool, Optional
            "allow_access_credential": TypeBool, Optional
            "allow_keyboard_logger":   TypeBool, Optional

            // Optional — 数值参数
            "max_file_up_size":              TypeInt, Optional
            "max_file_down_size":            TypeInt, Optional
            "max_access_credential_duration": TypeInt, Optional

            // Optional — 关联集合（TypeSet）
            "user_id_set":          TypeSet of TypeInt,    Optional
            "user_group_id_set":    TypeSet of TypeInt,    Optional
            "device_id_set":        TypeSet of TypeInt,    Optional
            "device_group_id_set":  TypeSet of TypeInt,    Optional
            "app_asset_id_set":     TypeSet of TypeInt,    Optional
            "account_set":          TypeSet of TypeString, Optional
            "cmd_template_id_set":  TypeSet of TypeInt,    Optional
            "ac_template_id_set":   TypeSet of TypeString, Optional

            // Optional — 时间和部门
            "validate_from":  TypeString, Optional, Computed
            "validate_to":    TypeString, Optional, Computed
            "department_id":  TypeString, Optional, Computed

            // Computed
            "acl_id": TypeInt, Computed   // 即资源 ID（uint64）
        },
    }
}
```

**资源 ID 设计**：使用 `CreateAcl` 返回的 `Id`（uint64）直接作为 `d.SetId`，转为字符串。`acl_id` 字段作为 Computed 字段同步写入，方便引用。

---

## 3. CRUD 函数设计

### Create

```
resourceTencentCloudBhAclCreate:
1. 构建 CreateAclRequest，填充所有字段
2. resource.Retry(WriteRetryTimeout) 调用 CreateAcl
3. 从 response.Response.Id 获取 aclId
4. d.SetId(helper.Uint64ToStr(*response.Response.Id))
5. 调用 Read
```

**关键字段处理：**
- `user_id_set`、`device_id_set` 等 TypeSet → `[]*uint64`：遍历 Set 元素，`helper.IntUint64(v.(int))`
- `account_set`、`ac_template_id_set` 等 TypeSet of String → `[]*string`：遍历元素，`helper.String(v.(string))`
- `validate_from` / `validate_to` 直接传字符串

### Read

```
resourceTencentCloudBhAclRead:
1. aclId = d.Id()
2. 调用 service.DescribeBhAclById(ctx, aclId)
3. 若 respData == nil → d.SetId(""), return nil
4. 调用 service.DescribeBhDepartments(ctx) 查询部门功能状态
5. 逐字段 d.Set(...)
6. department_id：仅当 Departments.Enabled == true 时才 set
7. 关联集合字段（UserSet → user_id_set）：从 respData.UserSet 提取 Id 列表
```

**department_id 的条件 set 逻辑：**
```
DescribeDepartments → Departments
    ├── nil 或 Enabled == nil → 跳过 department_id set
    └── *Enabled == true
             └── respData.Department != nil && Department.Id != nil
                     └── d.Set("department_id", ...)
```

### Update

```
resourceTencentCloudBhAclUpdate:
1. mutableArgs = 所有 Optional 字段 + "name"
2. 检测是否有变更
3. 若有变更：构建 ModifyAclRequest，填充 Id（从 d.Id() 解析）和所有字段
4. resource.Retry(WriteRetryTimeout) 调用 ModifyAcl
5. 调用 Read
```

### Delete

```
resourceTencentCloudBhAclDelete:
1. aclId = d.Id() → helper.StrToUint64Point(aclId)
2. 构建 DeleteAclsRequest，IdSet = []*uint64{aclId}
3. resource.Retry(WriteRetryTimeout) 调用 DeleteAcls
```

---

## 4. Service 层设计

### service_tencentcloud_bh.go 新增方法

```go
func (me *BhService) DescribeBhAclById(ctx context.Context, aclId string) (ret *bhv20230418.Acl, errRet error) {
    // ... 已实现
}

func (me *BhService) DescribeBhDepartments(ctx context.Context) (ret *bhv20230418.Departments, errRet error) {
    logId := tccommon.GetLogId(ctx)
    request := bhv20230418.NewDescribeDepartmentsRequest()
    response := bhv20230418.NewDescribeDepartmentsResponse()

    defer func() { /* error log */ }()

    err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
        ratelimit.Check(request.GetAction())
        result, e := me.client.UseBhV20230418Client().DescribeDepartments(request)
        if e != nil { return tccommon.RetryError(e) }
        if result == nil || result.Response == nil {
            return resource.NonRetryableError(...)
        }
        response = result
        return nil
    })

    if err != nil { errRet = err; return }

    ret = response.Response.Departments
    return
}
```

---

## 5. Provider 注册

在 `tencentcloud/provider.go` 的 `ResourcesMap` 中追加：

```go
"tencentcloud_bh_acl": bh.ResourceTencentCloudBhAcl(),
```

---

## 6. 文档示例设计（resource_tc_bh_acl.md）

```hcl
resource "tencentcloud_bh_acl" "example" {
  name                 = "tf-example-acl"
  allow_disk_redirect  = true
  allow_any_account    = false
  allow_clip_file_up   = true
  allow_clip_file_down = false
  allow_clip_text_up   = false
  allow_clip_text_down = false
  allow_file_up        = false
  allow_file_down      = false
  allow_disk_file_up   = false
  allow_disk_file_down = false
  allow_shell_file_up  = false
  allow_shell_file_down = false
  allow_file_del        = false
  allow_access_credential = true
  allow_keyboard_logger   = false
  department_id           = "1"
}
```

---

## 7. 测试设计（resource_tc_bh_acl_test.go）

测试场景：
1. **Basic 创建**：最小参数（name + allow_disk_redirect + allow_any_account）
2. **更新**：修改 name 和部分可选权限字段
3. **Import**：验证 `ImportStatePassthrough` 正常工作

测试函数命名：`TestAccTencentCloudBhAclResource_basic`

---

## 8. 关键实现细节

### TypeSet 字段的写入（Create/Update）

```go
if v, ok := d.GetOk("user_id_set"); ok {
    for _, id := range v.(*schema.Set).List() {
        request.UserIdSet = append(request.UserIdSet, helper.IntUint64(id.(int)))
    }
}
```

### TypeSet 字段的读取（Read）

```go
if respData.UserSet != nil && len(respData.UserSet) > 0 {
    userIdSet := make([]interface{}, 0, len(respData.UserSet))
    for _, user := range respData.UserSet {
        if user.Id != nil {
            userIdSet = append(userIdSet, int(*user.Id))
        }
    }
    _ = d.Set("user_id_set", userIdSet)
}
```

### Bool 字段的读取（注意 BH API 返回 *bool）

```go
if respData.AllowDiskRedirect != nil {
    _ = d.Set("allow_disk_redirect", respData.AllowDiskRedirect)
}
```

### ID 转换

- Create 时：`d.SetId(helper.Uint64ToStr(*response.Response.Id))`
- Delete/Read 时：`helper.StrToUint64Point(d.Id())`
- Update 时：`request.Id = helper.StrToUint64Point(d.Id())`
