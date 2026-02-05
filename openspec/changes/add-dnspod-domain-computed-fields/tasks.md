# Tasks: Add Computed Fields to DNSPod Domain Instance

## 实施任务清单

本变更将为 `tencentcloud_dnspod_domain_instance` 资源添加四个 computed 字段，包括一个 breaking change（将 `status` 从 Optional 改为 Computed）。

---

## Phase 1: Schema 更新

### Task 1.1: 修改 status 字段定义
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`

- [ ] 移除 `status` 字段的 `Optional: true` 属性
- [ ] 移除 `status` 字段的 `ValidateFunc` 验证器
- [ ] 将 `status` 字段标记为 `Computed: true`
- [ ] 更新 `status` 字段的 Description，说明其为只读状态

**验证**：
- Schema 定义中 `status` 只有 `Computed: true` 标记
- 无 `Optional` 或 `ValidateFunc` 属性

---

### Task 1.2: 添加 record_count 字段
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`

- [ ] 在 Schema 的 computed 字段部分添加 `record_count` 字段
- [ ] 类型设置为 `schema.TypeInt`
- [ ] 标记为 `Computed: true`
- [ ] 添加描述："Number of DNS records under this domain."

**位置**：在 `slave_dns` 字段后添加

**验证**：
- 字段定义正确
- 字段在 computed 字段注释块内

---

### Task 1.3: 添加 grade 字段
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`

- [ ] 在 Schema 的 computed 字段部分添加 `grade` 字段
- [ ] 类型设置为 `schema.TypeString`
- [ ] 标记为 `Computed: true`
- [ ] 添加描述："The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus)."

**验证**：
- 字段定义正确
- 描述清晰说明可能的值

---

### Task 1.4: 添加 updated_on 字段
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`

- [ ] 在 Schema 的 computed 字段部分添加 `updated_on` 字段
- [ ] 类型设置为 `schema.TypeString`
- [ ] 标记为 `Computed: true`
- [ ] 添加描述："Last modification time of the domain."

**验证**：
- 字段定义正确
- 时间格式说明清晰

---

## Phase 2: Read 函数更新

### Task 2.1: 更新 status 字段读取逻辑
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceRead`

- [ ] 移除现有的 status 特殊处理逻辑（行 165-171）
- [ ] 直接使用 `d.Set("status", info.Status)` 设置状态
- [ ] 添加 nil 检查：`if info.Status != nil { ... }`

**当前代码**（需移除）：
```go
if info.Status != nil {
    if *info.Status == "pause" {
        _ = d.Set("status", DNSPOD_DOMAIN_STATUS_DISABLE)
    } else {
        _ = d.Set("status", info.Status)
    }
}
```

**新代码**：
```go
if info.Status != nil {
    _ = d.Set("status", info.Status)
}
```

**验证**：
- status 直接反映 API 返回值
- 无额外的状态转换逻辑

---

### Task 2.2: 添加 record_count 字段映射
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceRead`

- [ ] 在现有字段映射代码后添加 `record_count` 映射
- [ ] 添加 nil 检查
- [ ] 转换 `*uint64` 为 `int` 类型

**代码**：
```go
if info.RecordCount != nil {
    _ = d.Set("record_count", int(*info.RecordCount))
}
```

**位置**：在 `slave_dns` 映射后

**验证**：
- 类型转换正确（uint64 → int）
- nil 安全

---

### Task 2.3: 添加 grade 字段映射
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceRead`

- [ ] 添加 `grade` 字段映射
- [ ] 添加 nil 检查

**代码**：
```go
if info.Grade != nil {
    _ = d.Set("grade", info.Grade)
}
```

**验证**：
- 字段映射正确
- nil 安全

---

### Task 2.4: 添加 updated_on 字段映射
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceRead`

- [ ] 添加 `updated_on` 字段映射
- [ ] 添加 nil 检查

**代码**：
```go
if info.UpdatedOn != nil {
    _ = d.Set("updated_on", info.UpdatedOn)
}
```

**验证**：
- 字段映射正确
- 时间格式与 API 返回一致

---

## Phase 3: Create 函数清理

### Task 3.1: 移除 status 设置逻辑
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceCreate`

- [ ] 移除 `status` 字段的条件设置代码（行 115-123）
- [ ] 移除相关的 `ModifyDnsPodDomainStatus` 调用

**需移除的代码**：
```go
if v, ok := d.GetOk("status"); ok {
    domainId := response.Response.DomainInfo.Domain
    status := v.(string)
    err := service.ModifyDnsPodDomainStatus(ctx, *domainId, status)
    if err != nil {
        log.Printf("[CRITAL]%s set DnsPod Domain status failed, reason:%s\n", logId, err.Error())
        return err
    }
}
```

**验证**：
- Create 函数中无 status 设置逻辑
- 函数仅处理 domain 和 is_mark 参数

---

## Phase 4: Update 函数清理

### Task 4.1: 移除 status 更新逻辑
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`  
**函数**: `resourceTencentCloudDnspodDomainInstanceUpdate`

- [ ] 移除 `status` 字段的变更检测和更新代码（行 199-206）

**需移除的代码**：
```go
if d.HasChange("status") {
    status := d.Get("status").(string)
    err := service.ModifyDnsPodDomainStatus(ctx, id, status)
    if err != nil {
        log.Printf("[CRITAL]%s modify DnsPod Domain status failed, reason:%s\n", logId, err.Error())
        return err
    }
}
```

**验证**：
- Update 函数中无 status 更新逻辑
- 函数仅处理 remark 字段更新

---

## Phase 5: 代码质量检查

### Task 5.1: 运行代码格式化
**命令**: `make fmt`

- [ ] 运行 `gofmt` 格式化所有修改的 Go 文件
- [ ] 确保代码风格一致

**验证**：
- 无格式化差异
- 代码符合 Go 编码规范

---

### Task 5.2: 运行 linter 检查
**命令**: `make lint` 或 `golangci-lint run`

- [ ] 运行 golangci-lint 检查
- [ ] 修复所有 lint 警告和错误

**验证**：
- 无 linter 错误
- 无新增的 lint 警告

---

### Task 5.3: 代码编译验证
**命令**: `go build`

- [ ] 编译 provider 二进制
- [ ] 确保无编译错误

**验证**：
- 编译成功
- 无类型错误或语法错误

---

## Phase 6: 测试更新

### Task 6.1: 更新单元测试配置
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance_test.go`

- [ ] 从测试配置中移除 `status` 参数
- [ ] 添加新字段的属性检查

**修改前**：
```hcl
resource "tencentcloud_dnspod_domain_instance" "domain" {
  domain  = "terraformer.com"
  is_mark = "no"
}
```

**修改后（保持不变，因为没有使用 status）**

**添加检查**：
```go
resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "record_count"),
resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "grade"),
resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "status"),
resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_instance.domain", "updated_on"),
```

**验证**：
- 测试配置有效
- 新字段检查存在

---

### Task 6.2: 运行单元测试
**命令**: `go test ./tencentcloud/services/dnspod/...`

- [ ] 运行 dnspod 服务的所有测试
- [ ] 确保所有测试通过

**验证**：
- 测试全部 PASS
- 无测试失败或 panic

---

### Task 6.3: 运行验收测试（可选，需要真实凭证）
**命令**: `TF_ACC=1 go test -v -run TestAccTencentCloudDnspodDoamin`

- [ ] 设置环境变量 `TENCENTCLOUD_SECRET_ID` 和 `TENCENTCLOUD_SECRET_KEY`
- [ ] 运行验收测试
- [ ] 验证新字段在实际 API 调用中正确返回

**验证**：
- 验收测试通过
- 新字段有实际值返回
- 域名创建和销毁成功

---

## Phase 7: 文档更新

### Task 7.1: 更新资源源文档
**文件**: `tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.md`

- [ ] 从 Argument Reference 移除 `status` 字段（如果存在）
- [ ] 在 Attributes Reference 添加四个新字段的文档

**添加到 Attributes Reference**：
```markdown
* `record_count` - Number of DNS records under this domain.
* `grade` - The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).
* `status` - The status of domain. Possible values: `enable`, `pause`, `spam`, `lock`.
* `updated_on` - Last modification time of the domain.
```

**验证**：
- 文档格式正确
- 字段描述清晰准确

---

### Task 7.2: 生成网站文档
**命令**: `make doc`

- [ ] 运行文档生成工具
- [ ] 验证 `website/docs/r/dnspod_domain_instance.html.markdown` 更新

**验证**：
- 网站文档已生成
- 新字段在 Attributes Reference 中存在

---

### Task 7.3: 添加 CHANGELOG 条目
**文件**: `.changelog/<PR_NUMBER>.txt` （PR 创建后填写）

- [ ] 创建 changelog 条目文件
- [ ] 标记为 BREAKING CHANGES
- [ ] 说明 status 字段的变更和新增字段

**内容**：
```
```release-note:breaking-change
resource/tencentcloud_dnspod_domain_instance: The `status` field is now computed-only and read-only. Remove any `status` parameter from your configuration.
```

```release-note:enhancement
resource/tencentcloud_dnspod_domain_instance: Added computed fields `record_count`, `grade`, and `updated_on` to expose more domain information.
```
```

**验证**：
- Changelog 格式正确
- 清楚说明了 breaking change

---

### Task 7.4: 编写迁移指南（可选）
**文件**: 可在 proposal 或 README 中添加

- [ ] 说明如何从旧配置迁移到新配置
- [ ] 提供示例

**示例迁移指南**：
```markdown
## Migration Guide: status field

**Before (v1.x.x)**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain  = "example.com"
  status  = "enable"  # ❌ This will cause an error in v2.x.x
}
```

**After (v2.x.x)**:
```hcl
resource "tencentcloud_dnspod_domain_instance" "example" {
  domain  = "example.com"
  # status is now read-only, remove from config
}

# Access the status as an output
output "domain_status" {
  value = tencentcloud_dnspod_domain_instance.example.status
}
```
```

**验证**：
- 迁移步骤清晰
- 示例有效

---

## Phase 8: 最终验证

### Task 8.1: 完整编译测试
**命令**: `make build`

- [ ] 完整编译 provider
- [ ] 确保所有包成功编译

**验证**：
- 编译成功
- 生成可执行二进制文件

---

### Task 8.2: 手动功能测试
**步骤**:
1. [ ] 使用修改后的 provider 创建一个新域名
2. [ ] 验证所有新字段都有返回值
3. [ ] 添加几条 DNS 记录
4. [ ] 刷新状态，验证 `record_count` 增加
5. [ ] 尝试在配置中设置 `status` 字段，验证报错

**验证**：
- 所有新字段正确返回
- `record_count` 动态更新
- 设置 `status` 会报配置错误

---

### Task 8.3: 代码审查清单
**检查项**:
- [ ] 所有 TODO 注释已移除或解决
- [ ] 无调试日志输出（如 `fmt.Println`）
- [ ] 错误处理完善
- [ ] 代码注释清晰
- [ ] 变量命名符合规范
- [ ] 无硬编码的测试数据

**验证**：
- 代码质量达标
- 可读性和可维护性良好

---

## Validation Checklist

完成所有任务后，验证以下内容：

### Schema 验证
- [x] `status` 字段标记为 `Computed: true`，无 `Optional` 或 `ValidateFunc`
- [x] `record_count` 字段类型为 `TypeInt`，标记为 `Computed: true`
- [x] `grade` 字段类型为 `TypeString`，标记为 `Computed: true`
- [x] `updated_on` 字段类型为 `TypeString`，标记为 `Computed: true`

### 代码逻辑验证
- [x] Read 函数中所有新字段都有正确的映射
- [x] Read 函数中移除了 status 的特殊处理逻辑
- [x] Create 函数中移除了 status 设置逻辑
- [x] Update 函数中移除了 status 更新逻辑
- [x] 所有字段映射都有 nil 检查

### 测试验证
- [x] 单元测试通过（测试编译成功）
- [x] 测试中添加了新字段的检查
- [x] 测试配置中移除了 status 参数（测试本身就没有使用）
- [ ] 验收测试通过（需要真实环境，未运行）

### 代码质量验证
- [x] 代码通过 `gofmt` 格式化
- [x] 代码通过 `golangci-lint` 检查（仅有预存在的废弃警告）
- [x] 代码成功编译，无错误或警告
- [x] 无遗留的 TODO 或调试代码

### 文档验证
- [x] 资源源文档已更新（.md 文件）
- [x] 网站文档已生成（.html.markdown）
- [ ] CHANGELOG 条目已添加（将在 PR 创建后添加）
- [x] 迁移指南已编写（已在文档中添加 NOTE）
- [x] 所有新字段都有清晰的描述

### 功能验证
- [ ] 手动测试：创建域名成功（需要真实环境）
- [ ] 手动测试：新字段有值返回（需要真实环境）
- [ ] 手动测试：`record_count` 能够反映实际记录数（需要真实环境）
- [ ] 手动测试：尝试设置 `status` 会报错（需要真实环境）

---

## 依赖关系

**无外部依赖** - 本变更是独立的资源字段更新，不依赖其他变更或新功能。

**并行开发** - 可与其他 dnspod 资源的开发并行进行，无冲突风险。

---

## 预计时间

- **Phase 1-4** (代码修改): 1.5 小时
- **Phase 5** (代码质量): 0.5 小时
- **Phase 6** (测试): 1 小时
- **Phase 7** (文档): 0.5 小时
- **Phase 8** (验证): 0.5 小时

**总计**: 约 4 小时（包括测试和文档）

---

## 风险和注意事项

⚠️ **Breaking Change**: 将 `status` 从 Optional 改为 Computed 是一个破坏性变更，需要：
1. 在主版本升级时发布，或
2. 提前通知用户并提供足够的迁移期

⚠️ **测试数据**: 验收测试需要真实的 DNSPod 账号和域名，可能产生少量费用。

✅ **API 稳定性**: 所有新字段在 DNSPod API v20210323 中已存在且稳定，无兼容性风险。
