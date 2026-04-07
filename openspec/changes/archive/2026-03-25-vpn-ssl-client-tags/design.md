## Context

当前 `tencentcloud_vpn_ssl_client` 资源不支持标签管理,而腾讯云 VPC API 已经支持在创建 VPN SSL Client 时设置 Tags 参数。为了与其他资源(如 `tencentcloud_vpn_ssl_server`)保持一致,需要添加标签支持。

**当前状态**:
- `resource_tc_vpn_ssl_client.go` 只有两个字段:`ssl_vpn_server_id` 和 `ssl_vpn_client_name`
- 资源不支持 Update 操作(只有 Create/Read/Delete)
- CreateVpnGatewaySslClient API 支持 Tags 参数

**约束**:
- 必须保持向后兼容,tags 为可选字段
- 遵循项目标准:使用通用 Tag Service 查询标签

## Goals / Non-Goals

**Goals:**
- 在创建 VPN SSL Client 时支持通过 API 设置标签
- 在 Read 操作中使用 Tag Service 同步标签到 state
- **添加 Update 函数,支持使用 Tag Service 原地更新标签**
- 提供完整的文档和测试用例
- 与其他 VPN 资源(vpn_ssl_server)的 tags 实现方式保持一致

**Non-Goals:**
- ~~不支持标签的 Update 操作~~(现在支持)
- 不实现标签的独立管理(使用控制台或其他方式修改的标签会在下次 Read 时同步)

## Decisions

### 决策 1: 使用 API Tags 参数而非 Tag Service 创建

**选择**: 直接使用 `CreateVpnGatewaySslClient` API 的 Tags 参数设置标签

**理由**:
- API 原生支持 Tags 参数,实现简单
- 避免额外的 API 调用
- 与 `vpn_ssl_server` 的实现方式一致

**替代方案**: 
- 创建后调用 Tag Service 的 `ModifyResourceTags` 添加标签 → 增加复杂度和 API 调用次数

### 决策 2: 使用 Tag Service 读取标签

**选择**: 在 Read 函数中使用通用 Tag Service 查询标签

**理由**:
- 与项目中其他资源保持一致的实现方式
- Tag Service 提供了统一的标签查询接口
- `DescribeVpnGatewaySslClients` API 不返回标签信息

**实现**:
```go
// In Read function
tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
tags, err := tagService.DescribeResourceTags(ctx, "vpc", "vpnx", region, sslClientId)
```

### 决策 3: 支持标签原地更新

**选择**: tags 字段设置为 Optional,**添加 Update 函数支持原地更新标签**

**理由**:
- 保持向后兼容,不影响现有配置
- 与 `vpn_ssl_server` 保持一致(其支持 tags 更新)
- 提升用户体验,无需重建资源即可更新标签
- 使用 Tag Service 的 `ModifyTags` 方法实现标签增删改

**实现**:
```go
// In Update function
if d.HasChange("tags") {
    oldInterface, newInterface := d.GetChange("tags")
    replaceTags, deleteTags := svctag.DiffTags(
        oldInterface.(map[string]interface{}), 
        newInterface.(map[string]interface{}),
    )
    tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
    region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
    resourceName := tccommon.BuildTagResourceName("vpc", "vpnx", region, sslClientId)
    err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
    // handle error...
}
```

### 决策 4: 资源定义添加 Update 回调

**选择**: 在资源定义中添加 `Update: resourceTencentCloudVpnSslClientUpdate`

**理由**:
- 当前资源只有 Create/Read/Delete,所有字段都是 ForceNew
- 添加 Update 函数允许 tags 原地更新
- 其他字段保持 ForceNew 不变(ssl_vpn_server_id, ssl_vpn_client_name)
- Update 函数仅处理 tags 变更,不影响其他字段的行为

## Risks / Trade-offs

### ~~[风险] 标签不支持原地更新~~
- **已解决**: 现在添加了 Update 函数,支持标签原地更新

### [风险] 外部标签修改会在 terraform apply 时覆盖
- **影响**: 在控制台或其他方式修改的标签会在 terraform refresh 时同步,在 apply 时会根据配置覆盖
- **缓解**: 这是 Terraform 的预期行为,state 管理标签的唯一真实来源

### [Trade-off] API Tags vs Tag Service
- **选择**: 创建时使用 API Tags,读取和更新时使用 Tag Service
- **好处**: 创建简单,更新统一
- **代价**: 创建和更新使用不同机制,但符合项目惯例
