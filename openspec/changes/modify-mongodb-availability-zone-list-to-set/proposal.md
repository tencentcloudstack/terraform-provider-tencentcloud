# Change: 修改 MongoDB 实例 availability_zone_list 字段从 List 改为 Set

## Why

目前 `tencentcloud_mongodb_instance` 和 `tencentcloud_mongodb_sharding_instance` 资源中的 `availability_zone_list` 字段使用 `TypeList` 类型,这存在以下问题:

1. **语义不符**: 可用区列表本质上是一个无序集合,不应该关心顺序。使用 List 类型会导致 Terraform 在可用区顺序改变时产生不必要的 diff。

2. **重复值检测**: List 类型允许重复值,但可用区列表中不应该有重复的可用区。Set 类型会自动去重并防止重复值。

3. **最佳实践**: Terraform 官方建议对于无序、无重复的集合使用 Set 类型,这样可以:
   - 避免因顺序变化导致的误判 diff
   - 自动处理重复值
   - 更好的表达字段的语义

4. **一致性**: 同一项目中其他类似的字段(如 `security_groups`)都使用 Set 类型。

## What Changes

将以下资源中的 `availability_zone_list` 字段从 `schema.TypeList` 改为 `schema.TypeSet`:

### 受影响的资源
1. `tencentcloud_mongodb_instance` - MongoDB 实例
2. `tencentcloud_mongodb_sharding_instance` - MongoDB 分片实例

### Schema 变更

**变更前:**
```go
"availability_zone_list": {
    Type:     schema.TypeList,
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
}
```

**变更后:**
```go
"availability_zone_list": {
    Type:     schema.TypeSet,
    Optional: true,
    Computed: true,
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
}
```

### 代码变更

#### 1. `resource_tc_mongodb_instance.go`

**mongodbAllInstanceReqSet 函数 (第 235-238 行):**
```go
// 变更前
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.([]interface{}))
    value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
}

// 变更后
if v, ok := d.GetOk("availability_zone_list"); ok {
    availabilityZoneList := helper.InterfacesStringsPoint(v.(*schema.Set).List())
    value.FieldByName("AvailabilityZoneList").Set(reflect.ValueOf(availabilityZoneList))
}
```

**resourceTencentCloudMongodbInstanceRead 函数 (第 477-486 行):**
```go
// 变更后添加去重和排序逻辑
if len(replicateSets) > 0 {
    var hiddenZone string
    availabilityZoneMap := make(map[string]bool)
    for _, replicate := range replicateSets[0].Nodes {
        itemZone := *replicate.Zone
        if *replicate.Hidden {
            hiddenZone = itemZone
        }
        availabilityZoneMap[itemZone] = true
    }
    
    availabilityZoneList := make([]string, 0, len(availabilityZoneMap))
    for zone := range availabilityZoneMap {
        availabilityZoneList = append(availabilityZoneList, zone)
    }
    
    _ = d.Set("hidden_zone", hiddenZone)
    _ = d.Set("availability_zone_list", availabilityZoneList)
}
```

#### 2. `resource_tc_mongodb_sharding_instance.go`

类似的变更:
- 第 159-161 行: 从 `v.([]interface{})` 改为 `v.(*schema.Set).List()`
- 第 394 行: 添加去重逻辑

#### 3. `resource_tc_mongodb_readonly_instance.go`

- 第 163 行: 从 `v.([]interface{})` 改为 `v.(*schema.Set).List()`

### 测试变更

所有测试文件中的 `resource.TestCheckResourceAttr` 断言保持不变,因为 Set 类型仍然支持 `.#` 语法来检查元素数量。

示例:
```go
resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "availability_zone_list.#", "5"),
```

### 文档变更

需要在资源文档中说明该字段为 Set 类型(虽然对于用户来说使用方式不变)。

## Impact

### 向后兼容性

⚠️ **这是一个破坏性变更**,会影响现有的 Terraform 状态文件:

1. **状态迁移**: 用户在升级 Provider 版本后首次执行 `terraform plan` 时,会看到 `availability_zone_list` 字段的变更,因为 Terraform 会将状态中的 List 类型转换为 Set 类型。

2. **配置文件兼容**: HCL 配置文件中的写法保持不变,List 和 Set 的语法相同:
   ```hcl
   availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-6"]
   ```

3. **状态文件格式变更**:
   - **变更前 (List)**: 
     ```json
     "availability_zone_list": ["ap-guangzhou-3", "ap-guangzhou-4"]
     ```
   - **变更后 (Set)**:
     ```json
     "availability_zone_list": {
       "value": ["ap-guangzhou-3", "ap-guangzhou-4"],
       "type": "set"
     }
     ```

### 迁移建议

用户升级 Provider 版本后需要:
1. 执行 `terraform plan` 查看变更
2. 如果只是类型转换(List -> Set),执行 `terraform apply` 确认变更
3. 状态文件会自动更新为新格式

### 受影响的用户场景

1. **多可用区部署**: 使用 `availability_zone_list` 部署跨可用区 MongoDB 实例的用户
2. **自动化脚本**: 通过 API 或脚本读取 Terraform 状态文件的用户需要适配新的 Set 格式

### 优势

1. ✅ 更准确的语义表达(无序集合)
2. ✅ 自动去重,避免配置错误
3. ✅ 避免因顺序变化产生的误判 diff
4. ✅ 符合 Terraform 最佳实践
5. ✅ 与项目中其他类似字段保持一致

### 风险

1. ⚠️ 破坏性变更,需要在 CHANGELOG 中明确标注
2. ⚠️ 需要提供迁移指南
3. ⚠️ 建议在主要版本升级时发布(如 v2.0.0)

### 依赖关系

- 不影响其他资源
- 不影响 API 调用
- 只影响 Terraform Schema 和状态文件格式

### 测试影响

- 现有测试用例无需修改(`.#` 语法兼容)
- 需要添加测试用例验证:
  - Set 类型的去重功能
  - 顺序变化不产生 diff
  - 状态迁移正确性

## Implementation Notes

### 实现步骤

1. **修改 Schema 定义**: 将 `TypeList` 改为 `TypeSet`
2. **修改数据读取逻辑**: `v.([]interface{})` → `v.(*schema.Set).List()`
3. **修改数据写入逻辑**: 确保返回的数据可以正确转换为 Set
4. **更新测试**: 验证 Set 行为正常
5. **更新文档**: 说明字段类型和迁移注意事项
6. **执行 go fmt**: 格式化所有修改的文件

### 代码审查要点

- [ ] 所有读取 `availability_zone_list` 的地方都已更新
- [ ] 所有设置 `availability_zone_list` 的地方都已更新
- [ ] 去重逻辑正确实现
- [ ] 测试用例覆盖 Set 特性
- [ ] 文档更新完整
- [ ] CHANGELOG 中标注为破坏性变更

### 发布建议

建议在主要版本更新(如 v2.0.0)中发布此变更,并在发布说明中:
1. 明确标注为 **BREAKING CHANGE**
2. 提供迁移步骤和示例
3. 说明变更原因和优势
4. 提供回滚建议(如果需要保持旧版本)
