# Proposal: Add Computed Fields to DNSPod Domain Instance Resource

## 摘要

为 `tencentcloud_dnspod_domain_instance` 资源添加四个新的 computed 参数：`record_count`、`grade`、`status`（从 Optional+Computed 改为完全 Computed）、`updated_on`，以便用户能够获取更完整的域名信息。

## 问题描述

当前 `tencentcloud_dnspod_domain_instance` 资源仅暴露了部分域名信息作为 computed 字段（`domain_id`、`create_time`、`slave_dns`），而腾讯云 DNSPod API 的 `DomainInfo` 结构体返回了更多有用的信息。用户需要这些额外的字段来：

1. **记录数量 (RecordCount)**: 了解域名下的解析记录总数，用于监控和管理
2. **套餐等级 (Grade)**: 查看域名当前的套餐等级（免费版、专业版等）
3. **域名状态 (Status)**: 获取真实的域名状态（当前 status 字段是 Optional 可写的，会导致混淆）
4. **最后更新时间 (UpdatedOn)**: 追踪域名配置的最后修改时间

当前存在的问题：
- `status` 字段当前是 **Optional + Computed**，但实际上它应该是**只读的状态字段**
- 用户无法通过 Terraform 获取域名的记录数量、套餐等级和更新时间
- 这些信息在 API 响应中已经存在，只是未被映射到资源的 schema 中

## 解决方案

### 方案概述

将 `status` 字段从 Optional 改为完全 Computed（只读），并添加三个新的 computed 字段，从 DNSPod API 返回的 `DomainInfo` 结构中读取并映射这些值。

### 变更范围

**单一资源变更**：仅涉及 `tencentcloud_dnspod_domain_instance` 资源的 schema 和 read 逻辑更新。

### 设计选择

#### 1. **将 status 从 Optional 改为 Computed**

**理由**：
- DNSPod API 的 `Status` 字段表示域名的真实运行状态（`enable`/`pause`/`spam`/`lock`），这是**系统状态而非用户配置**
- 当前的 `status` 作为 Optional 字段存在，实际上是通过另外的 API 调用 `ModifyDnsPodDomainStatus` 来设置启用/暂停状态
- 将其改为 Computed 可以准确反映域名的实际状态，用户可以通过新增的专门资源（如未来可能的 `tencentcloud_dnspod_domain_status` 资源）来控制状态

**兼容性影响**：
- ⚠️ **BREAKING CHANGE**: 现有配置中使用 `status` 参数的用户需要移除该字段
- 建议在下一个主版本变更时实施，或提供迁移指南

**替代方案**：保持 `status` 为 Optional + Computed，但这会导致混淆（状态是用户设置的还是系统返回的？）

#### 2. **字段映射策略**

所有新增字段均为 **Computed Only**，因为它们是系统返回的只读信息：

| Schema 字段     | SDK 字段          | 类型     | 说明                                           |
|----------------|------------------|----------|------------------------------------------------|
| `record_count` | `RecordCount`    | `int`    | 域名下的解析记录总数                            |
| `grade`        | `Grade`          | `string` | 域名套餐等级（如 "DP_Free", "DP_Plus" 等）      |
| `status`       | `Status`         | `string` | 域名状态（enable/pause/spam/lock）**改为只读**  |
| `updated_on`   | `UpdatedOn`      | `string` | 域名最后更新时间（RFC3339 格式）                 |

### 实施细节

#### Schema 变更

**修改前**：
```go
"status": {
    Type:         schema.TypeString,
    Optional:     true,  // 当前是可写的
    ValidateFunc: tccommon.ValidateAllowedStringValue(DNSPOD_DOMAIN_STATUS_TYPE),
    Description:  "The status of Domain.",
},
```

**修改后**：
```go
"status": {
    Type:        schema.TypeString,
    Computed:    true,  // 改为只读
    Description: "The status of domain. Possible values: `enable`, `pause`, `spam`, `lock`.",
},
// 新增字段
"record_count": {
    Type:        schema.TypeInt,
    Computed:    true,
    Description: "Number of DNS records under this domain.",
},
"grade": {
    Type:        schema.TypeString,
    Computed:    true,
    Description: "The DNS plan/package grade of the domain (e.g., DP_Free, DP_Plus).",
},
"updated_on": {
    Type:        schema.TypeString,
    Computed:    true,
    Description: "Last modification time of the domain.",
},
```

#### Read 函数变更

在 `resourceTencentCloudDnspodDomainInstanceRead` 中添加字段映射：

```go
// 现有代码
_ = d.Set("domain_id", info.DomainId)
_ = d.Set("domain", info.Domain)
_ = d.Set("create_time", info.CreatedOn)
_ = d.Set("is_mark", info.IsMark)
_ = d.Set("slave_dns", info.SlaveDNS)

// 新增映射
if info.Status != nil {
    _ = d.Set("status", info.Status)  // 直接使用 API 返回的状态
}

if info.RecordCount != nil {
    _ = d.Set("record_count", int(*info.RecordCount))
}

if info.Grade != nil {
    _ = d.Set("grade", info.Grade)
}

if info.UpdatedOn != nil {
    _ = d.Set("updated_on", info.UpdatedOn)
}
```

#### Create/Update 函数变更

- **移除** `Create` 函数中设置 `status` 的逻辑（行 115-123）
- **移除** `Update` 函数中修改 `status` 的逻辑（行 199-206）

### 测试策略

1. **单元测试更新**：
   - 更新现有测试配置，移除 `status` 参数
   - 添加新字段的 `TestCheckResourceAttrSet` 验证

2. **验收测试**：
   ```hcl
   resource "tencentcloud_dnspod_domain_instance" "test" {
     domain  = "terraform-test.com"
     is_mark = "no"
   }
   
   # 验证新字段
   output "domain_info" {
     value = {
       record_count = tencentcloud_dnspod_domain_instance.test.record_count
       grade        = tencentcloud_dnspod_domain_instance.test.grade
       status       = tencentcloud_dnspod_domain_instance.test.status
       updated_on   = tencentcloud_dnspod_domain_instance.test.updated_on
     }
   }
   ```

3. **手动测试**：
   - 创建域名并验证所有新字段有值
   - 添加解析记录后验证 `record_count` 增加
   - 验证 `status` 字段不能被手动设置（配置错误提示）

### 文档更新

在 `website/docs/r/dnspod_domain_instance.html.markdown` 中：

1. **移除** Argument Reference 中的 `status`
2. **添加** 到 Attributes Reference：
   ```markdown
   * `record_count` - Number of DNS records under this domain.
   * `grade` - The DNS plan/package grade of the domain.
   * `status` - The status of domain. Possible values: `enable`, `pause`, `spam`, `lock`.
   * `updated_on` - Last modification time of the domain.
   ```
3. **添加迁移指南**（CHANGELOG）：
   ```markdown
   BREAKING CHANGES:
   * **resource/tencentcloud_dnspod_domain_instance**: The `status` field is now computed-only. 
     Remove any `status` parameter from your configuration. Use the read-only `status` attribute 
     to get the actual domain status.
   ```

## 影响分析

### 用户影响

**Breaking Change**:
- 使用 `status = "enable"` 或 `status = "disable"` 的用户需要移除该配置
- 影响范围：仅影响显式设置 `status` 参数的配置（预计较少）

**Benefits**:
- ✅ 获取域名的完整状态信息
- ✅ 能够监控域名的解析记录数量
- ✅ 了解域名的套餐等级
- ✅ 追踪域名配置的更新时间
- ✅ `status` 字段的语义更清晰（只读状态 vs. 可写配置）

### 技术风险

- **低风险**：仅添加 computed 字段，不影响现有的 Create/Delete 逻辑
- **API 兼容性**：所有字段在 DNSPod API v20210323 中已存在且稳定
- **测试覆盖**：需要更新测试用例以验证新字段

## 替代方案

### 方案 A：保持 status 为 Optional（不推荐）

**优点**：
- 无 breaking change
- 保持向后兼容

**缺点**：
- 语义混淆：`status` 既是配置项又是状态输出
- API 行为不一致：设置 status 需要额外 API 调用，但读取时返回的是真实状态
- 用户体验差：不清楚 `status` 字段的真实含义

### 方案 B：添加新的 domain_status 字段（不推荐）

**优点**：
- 完全向后兼容
- 区分配置的 `status` 和实际的 `domain_status`

**缺点**：
- 造成字段冗余和混淆
- 不符合 Terraform 最佳实践（状态字段应该是 computed）

## 实施计划

**阶段 1：代码实现**（预计 2 小时）
1. 修改 schema 定义
2. 更新 Read 函数添加字段映射
3. 移除 Create/Update 中的 status 设置逻辑
4. 运行 `gofmt` 和 `golangci-lint`

**阶段 2：测试**（预计 1 小时）
1. 更新单元测试
2. 运行验收测试
3. 手动测试验证

**阶段 3：文档**（预计 30 分钟）
1. 更新资源文档
2. 添加 CHANGELOG 条目
3. 编写迁移指南

**总预计时间**：3.5 小时

## 验收标准

- [ ] `record_count`、`grade`、`updated_on` 字段添加到 schema 并标记为 Computed
- [ ] `status` 字段从 Optional 改为 Computed
- [ ] Read 函数正确映射所有新字段
- [ ] Create 和 Update 函数中移除了 status 设置逻辑
- [ ] 所有测试通过（包括验收测试）
- [ ] 代码通过 lint 检查
- [ ] 文档已更新（包括迁移指南）
- [ ] CHANGELOG 已添加 BREAKING CHANGES 条目

## 参考资料

- DNSPod API 文档：`DescribeDomain` 接口
- SDK 源码：`vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323/models.go#L5171-L5270`
- 现有资源代码：`tencentcloud/services/dnspod/resource_tc_dnspod_domain_instance.go`
