# MongoDB availability_zone_list DiffSuppressFunc 实施任务

## 📋 任务概览

**变更目标**: 为 `availability_zone_list` 字段添加 DiffSuppressFunc,忽略列表顺序差异  
**变更类型**: Enhancement (非破坏性)  
**预计时间**: 2 小时

---

## ✅ 任务清单

### Phase 1: 准备工作

- [ ] 1.1 创建 `diff_suppress_funcs.go` 文件
- [ ] 1.2 实现 `mongodbAvailabilityZoneListDiffSuppress` 函数
- [ ] 1.3 添加辅助函数(如需要)

### Phase 2: 单元测试

- [ ] 2.1 创建 `diff_suppress_funcs_test.go` 文件
- [ ] 2.2 测试场景: 顺序不同但内容相同 → 返回 true
- [ ] 2.3 测试场景: 内容不同 → 返回 false
- [ ] 2.4 测试场景: 长度不同 → 返回 false
- [ ] 2.5 测试场景: 空列表 → 返回 true
- [ ] 2.6 测试场景: nil 值处理
- [ ] 2.7 测试场景: 子元素路径处理

### Phase 3: Resource 集成

#### 3.1 resource_tc_mongodb_instance.go
- [ ] 3.1.1 在 `availability_zone_list` schema 中添加 `DiffSuppressFunc`
- [ ] 3.1.2 导入 diff suppress 函数
- [ ] 3.1.3 验证其他逻辑不变

#### 3.2 resource_tc_mongodb_sharding_instance.go
- [ ] 3.2.1 在 `availability_zone_list` schema 中添加 `DiffSuppressFunc`
- [ ] 3.2.2 导入 diff suppress 函数
- [ ] 3.2.3 验证其他逻辑不变

#### 3.3 resource_tc_mongodb_readonly_instance.go
- [ ] 3.3.1 确认是否有自己的 schema 定义
- [ ] 3.3.2 如有,添加 `DiffSuppressFunc`
- [ ] 3.3.3 如无,确认继承自哪里

### Phase 4: 代码质量

- [ ] 4.1 运行 `go fmt` 格式化所有修改的文件
- [ ] 4.2 运行 `go vet` 检查代码问题
- [ ] 4.3 运行 linter 检查
- [ ] 4.4 修复所有新增的警告和错误
- [ ] 4.5 确保单元测试全部通过

### Phase 5: 集成测试(可选)

- [ ] 5.1 创建测试 MongoDB 实例
- [ ] 5.2 验证顺序不同不触发 diff
- [ ] 5.3 验证内容不同正常触发 diff
- [ ] 5.4 清理测试资源

### Phase 6: 文档更新

- [ ] 6.1 创建 `change.md` 记录详细变更
- [ ] 6.2 创建 `README.md` 快速参考
- [ ] 6.3 更新相关资源文档(如需要)
- [ ] 6.4 准备 CHANGELOG 条目

---

## 📝 详细实施步骤

### Step 1: 创建 DiffSuppressFunc

**文件**: `tencentcloud/services/mongodb/diff_suppress_funcs.go`

```go
package mongodb

import (
    "sort"
    "strings"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// mongodbAvailabilityZoneListDiffSuppress 忽略 availability_zone_list 的顺序差异
// 仅当列表内容相同但顺序不同时,返回 true 抑制 diff
func mongodbAvailabilityZoneListDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
    // 只在列表级别比较,忽略子元素
    if !strings.HasSuffix(k, "availability_zone_list") && !strings.HasSuffix(k, "availability_zone_list.#") {
        return false
    }

    // 获取完整的列表
    oldList, newList := d.GetChange("availability_zone_list")
    
    // 处理 nil 情况
    if oldList == nil && newList == nil {
        return true
    }
    if oldList == nil || newList == nil {
        return false
    }

    // 转换为字符串切片
    oldZones := helper.InterfacesStrings(oldList.([]interface{}))
    newZones := helper.InterfacesStrings(newList.([]interface{}))

    // 长度不同,不抑制 diff
    if len(oldZones) != len(newZones) {
        return false
    }

    // 长度为 0,认为相同
    if len(oldZones) == 0 {
        return true
    }

    // 复制并排序后比较
    oldSorted := make([]string, len(oldZones))
    newSorted := make([]string, len(newZones))
    copy(oldSorted, oldZones)
    copy(newSorted, newZones)
    
    sort.Strings(oldSorted)
    sort.Strings(newSorted)

    // 逐个比较
    for i := range oldSorted {
        if oldSorted[i] != newSorted[i] {
            return false
        }
    }

    // 内容相同,仅顺序不同,抑制 diff
    return true
}
```

---

### Step 2: 创建单元测试

**文件**: `tencentcloud/services/mongodb/diff_suppress_funcs_test.go`

```go
package mongodb

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestMongodbAvailabilityZoneListDiffSuppress(t *testing.T) {
    testCases := []struct {
        name     string
        old      []interface{}
        new      []interface{}
        expected bool
    }{
        {
            name:     "相同顺序",
            old:      []interface{}{"zone-a", "zone-b", "zone-c"},
            new:      []interface{}{"zone-a", "zone-b", "zone-c"},
            expected: true,
        },
        {
            name:     "不同顺序相同内容",
            old:      []interface{}{"zone-a", "zone-b", "zone-c"},
            new:      []interface{}{"zone-c", "zone-a", "zone-b"},
            expected: true,
        },
        {
            name:     "内容不同",
            old:      []interface{}{"zone-a", "zone-b"},
            new:      []interface{}{"zone-a", "zone-c"},
            expected: false,
        },
        {
            name:     "长度不同",
            old:      []interface{}{"zone-a", "zone-b"},
            new:      []interface{}{"zone-a"},
            expected: false,
        },
        {
            name:     "空列表",
            old:      []interface{}{},
            new:      []interface{}{},
            expected: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
                "availability_zone_list": {
                    Type:     schema.TypeList,
                    Optional: true,
                    Elem:     &schema.Schema{Type: schema.TypeString},
                },
            }, map[string]interface{}{
                "availability_zone_list": tc.old,
            })

            d.Set("availability_zone_list", tc.new)

            result := mongodbAvailabilityZoneListDiffSuppress(
                "availability_zone_list",
                "",
                "",
                d,
            )

            if result != tc.expected {
                t.Errorf("expected %v, got %v", tc.expected, result)
            }
        })
    }
}
```

---

### Step 3: 修改 Resource Schema

**示例**: `resource_tc_mongodb_instance.go`

```go
// 在 schema 定义中添加 DiffSuppressFunc
"availability_zone_list": {
    Type:             schema.TypeList,
    Optional:         true,
    Computed:         true,
    DiffSuppressFunc: mongodbAvailabilityZoneListDiffSuppress, // 添加这行
    Elem: &schema.Schema{
        Type: schema.TypeString,
    },
    Description: "...",
},
```

---

### Step 4: 验证步骤

#### 4.1 运行单元测试
```bash
cd tencentcloud/services/mongodb
go test -v -run TestMongodbAvailabilityZoneListDiffSuppress
```

#### 4.2 格式化代码
```bash
go fmt ./tencentcloud/services/mongodb/...
```

#### 4.3 Lint 检查
```bash
golangci-lint run ./tencentcloud/services/mongodb/...
```

#### 4.4 编译检查
```bash
go build ./tencentcloud/services/mongodb/...
```

---

## 🎯 验收标准

### 功能性
- [ ] 顺序不同但内容相同时,`terraform plan` 不显示 diff
- [ ] 内容不同时,`terraform plan` 正常显示 diff
- [ ] 创建时顺序保持不变
- [ ] 更新操作正常工作

### 代码质量
- [ ] 所有单元测试通过
- [ ] 代码覆盖率 > 80%
- [ ] 无 linter 警告
- [ ] 代码格式符合规范

### 文档完整性
- [ ] DiffSuppressFunc 有清晰注释
- [ ] 测试用例覆盖所有场景
- [ ] change.md 记录完整
- [ ] README.md 说明清晰

---

## 📊 预估工作量

| 阶段 | 预估时间 | 说明 |
|------|---------|------|
| 准备工作 | 30 分钟 | 创建函数和基础结构 |
| 单元测试 | 30 分钟 | 编写全面的测试用例 |
| Resource 集成 | 30 分钟 | 修改 3 个资源文件 |
| 代码质量检查 | 15 分钟 | 格式化、lint、编译 |
| 文档更新 | 15 分钟 | 创建文档和说明 |
| **总计** | **2 小时** | |

---

## 🚨 注意事项

1. **不要修改 CRUD 逻辑**: 只添加 DiffSuppressFunc,不改变创建/读取/更新逻辑
2. **保持 List 类型**: 不要改为 Set,保持 TypeList
3. **子元素处理**: DiffSuppressFunc 会被列表的每个元素调用,需要正确处理
4. **nil 安全**: 注意处理 nil 值和空列表的情况
5. **性能考虑**: 避免在 DiffSuppressFunc 中执行重量级操作

---

## 📚 参考实现

可以参考 Terraform AWS Provider 中类似的实现:
- `aws_instance` 的 `security_groups` 字段
- `aws_lb_target_group` 的 tags 处理

---

**任务状态**: 待开始  
**负责人**: 待分配  
**优先级**: High
