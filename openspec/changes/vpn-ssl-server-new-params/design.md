## Context

`tencentcloud_vpn_ssl_server` 资源目前支持 VPN SSL Server 的基本配置（网关ID、名称、地址段、协议、端口、加密算法等）。腾讯云 VPC SDK 已更新，`CreateVpnGatewaySslServer` 和 `ModifyVpnGatewaySslServer` 接口新增了5个参数：

1. **SsoEnabled** (bool): 支持 SSO 单点登录认证
2. **AccessPolicyEnabled** (bool): 支持访问策略控制
3. **SamlData** (string): SAML 数据配置
4. **Tags** ([]*Tag): 标签管理
5. **DnsServers** (*DnsServers): 自定义 DNS 配置（包含 PrimaryDns 和 SecondaryDns）

**当前状态**:
- 资源文件: `tencentcloud/services/vpn/resource_tc_vpn_ssl_server.go`
- SDK 已更新: vendor 目录包含最新的 VPC SDK models
- 向后兼容: 所有新参数均为可选参数，默认行为不变

**约束**:
- 必须保持向后兼容性，不能破坏现有 Terraform 配置
- 所有新参数必须为 Optional + Computed
- 需要在 Create、Update、Read 三个函数中同步处理

## Goals / Non-Goals

**Goals:**
- 在 Schema 中新增5个参数，保持与 SDK 定义一致
- 支持 SSO 认证功能（sso_enabled + saml_data）
- 支持访问策略控制（access_policy_enabled）
- 支持标签管理（tags，使用 Provider 通用标签模式）
- 支持自定义 DNS 配置（dns_servers 嵌套对象）
- 所有 CRUD 操作完整支持新参数
- 生成完整的文档和示例

**Non-Goals:**
- 不修改现有参数的行为
- 不涉及 state 迁移或数据转换
- 不处理 SSO 白名单申请流程（由用户自行申请）
- 不实现访问策略的详细配置（AccessPolicyEnabled 仅为开关，具体策略通过其他资源配置）

## Decisions

### 1. Schema 设计

**决策**: 使用 Terraform Plugin SDK v2 标准模式，所有新参数设为 Optional + Computed

**理由**:
- **Optional**: 不强制用户配置，保持向后兼容
- **Computed**: 云端可能返回默认值，避免 diff 误报
- **dns_servers 嵌套对象**: 使用 `TypeList` + `MaxItems: 1` + `Elem: &schema.Resource{}`，符合 Terraform 嵌套对象最佳实践

**替代方案**:
- ❌ 设为 Required: 会破坏现有配置
- ❌ dns_servers 使用 TypeMap: 无法表达明确的字段类型（primary_dns/secondary_dns）

```go
"sso_enabled": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Enable SSO authentication. Default: false. This feature requires whitelist approval.",
},
"access_policy_enabled": {
    Type:        schema.TypeBool,
    Optional:    true,
    Computed:    true,
    Description: "Enable access policy control. Default: false.",
},
"saml_data": {
    Type:        schema.TypeString,
    Optional:    true,
    Description: "SAML-DATA. Required when sso_enabled is true.",
},
"tags": {
    Type:        schema.TypeMap,
    Optional:    true,
    Description: "Tags for resource management.",
    Elem:        &schema.Schema{Type: schema.TypeString},
},
"dns_servers": {
    Type:        schema.TypeList,
    Optional:    true,
    Computed:    true,
    MaxItems:    1,
    Description: "DNS server configuration.",
    Elem: &schema.Resource{
        Schema: map[string]*schema.Schema{
            "primary_dns": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Primary DNS server address.",
            },
            "secondary_dns": {
                Type:        schema.TypeString,
                Optional:    true,
                Description: "Secondary DNS server address.",
            },
        },
    },
},
```

### 2. Tags 处理模式

**决策**: 使用 Provider 的通用标签处理函数 `helper.GetTags` 和状态回设

**理由**:
- 与其他资源（如 CLS Alarm）保持一致
- 自动处理标签的 key-value 映射
- 无需手动构造 `[]*vpc.Tag` 结构

**实现**:
```go
// Create/Update
if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
    request.Tags = make([]*vpc.Tag, 0, len(tags))
    for k, v := range tags {
        request.Tags = append(request.Tags, &vpc.Tag{
            Key:   helper.String(k),
            Value: helper.String(v),
        })
    }
}

// Read
if info.Tags != nil {
    tagMap := make(map[string]string, len(info.Tags))
    for _, tag := range info.Tags {
        tagMap[*tag.Key] = *tag.Value
    }
    _ = d.Set("tags", tagMap)
}
```

### 3. DnsServers 对象处理

**决策**: 使用嵌套 List 处理，手动解析和构造 SDK 对象

**理由**:
- SDK 定义为 `*DnsServers` 结构体，不是 map
- 需要明确的字段名（PrimaryDns/SecondaryDns）
- Terraform 嵌套对象最佳实践

**实现**:
```go
// Create/Update: Terraform List -> SDK Object
if v, ok := d.GetOk("dns_servers"); ok {
    dnsServersList := v.([]interface{})
    if len(dnsServersList) > 0 {
        dnsServersMap := dnsServersList[0].(map[string]interface{})
        request.DnsServers = &vpc.DnsServers{}
        if v, ok := dnsServersMap["primary_dns"]; ok && v.(string) != "" {
            request.DnsServers.PrimaryDns = helper.String(v.(string))
        }
        if v, ok := dnsServersMap["secondary_dns"]; ok && v.(string) != "" {
            request.DnsServers.SecondaryDns = helper.String(v.(string))
        }
    }
}

// Read: SDK Object -> Terraform List
if info.DnsServers != nil {
    dnsServersMap := map[string]interface{}{}
    if info.DnsServers.PrimaryDns != nil {
        dnsServersMap["primary_dns"] = *info.DnsServers.PrimaryDns
    }
    if info.DnsServers.SecondaryDns != nil {
        dnsServersMap["secondary_dns"] = *info.DnsServers.SecondaryDns
    }
    _ = d.Set("dns_servers", []interface{}{dnsServersMap})
}
```

### 4. Update 函数的可变参数列表

**决策**: 在 `mutableArgs` 数组中添加所有新参数

**理由**:
- 现有代码使用 `mutableArgs` 检查哪些字段可以更新
- 所有新参数均支持通过 `ModifyVpnGatewaySslServer` 接口更新
- 保持代码风格一致

```go
mutableArgs := []string{
    "ssl_vpn_server_name", "local_address", "remote_address", "ssl_vpn_protocol",
    "ssl_vpn_port", "integrity_algorithm", "encrypt_algorithm", "compress",
    // 新增
    "sso_enabled", "access_policy_enabled", "saml_data", "tags", "dns_servers",
}
```

### 5. Read 函数中的 DnsServers 读取

**决策**: 检查 `info.DnsServers != nil` 后再处理

**理由**:
- SDK 返回的 `DnsServers` 字段可能为 nil（未配置时）
- 需要避免空指针异常
- 只在有值时设置到 state

## Risks / Trade-offs

### Risk 1: SSO 功能需要白名单
**风险**: 用户设置 `sso_enabled = true` 但未申请白名单，导致 API 报错  
**缓解**: 
- 在文档中明确说明需要申请白名单
- 错误信息会由云 API 返回，用户可根据错误提示申请

### Risk 2: SAML 数据格式验证
**风险**: `saml_data` 格式不正确导致 SSO 配置失败  
**缓解**:
- 不在 Provider 层做格式验证，由云 API 负责
- 在文档中提供正确的 SAML 配置示例
- 用户可参考腾讯云控制台的 SAML 配置指引

### Risk 3: DnsServers 为空值的处理
**风险**: 用户配置空的 dns_servers 块可能导致意外行为  
**缓解**:
- 在 Create/Update 中检查 `dnsServersList` 长度
- 只有当列表非空且至少一个字段有值时才设置 SDK 参数
- 空块会被忽略，保持默认行为

### Risk 4: Tags 的 API 一致性
**风险**: DescribeVpnGatewaySslServers 可能不返回 Tags 字段  
**缓解**:
- 在 Read 函数中添加 `if info.Tags != nil` 检查
- 如果 API 不返回 Tags，state 中会保留用户配置的值（Computed 属性）
- 需要在测试中验证 Tags 的读取行为

## Migration Plan

**部署步骤**:
1. ✅ 确认 vendor 目录已包含最新 SDK（已完成）
2. 修改 Schema 定义，添加5个新参数
3. 修改 Create 函数，映射新参数到 API 请求
4. 修改 Update 函数，添加可变参数和映射逻辑
5. 修改 Read 函数，回设新参数到 state
6. 添加测试用例（basic + sso + dns）
7. 更新 .md 文档模板，添加新参数示例
8. 运行 `make doc` 生成最终文档

**回滚策略**:
- 代码级回滚: 删除新参数，恢复原 Schema
- 无需数据迁移: 所有参数为 Optional，不影响已创建资源

**向后兼容性验证**:
- ✅ 现有测试用例不使用新参数，应继续通过
- ✅ 现有 Terraform 配置不包含新参数，apply 无变更

## Open Questions

1. **DescribeVpnGatewaySslServers 是否返回所有新字段？**  
   需要在测试中验证 API 响应是否包含 `SsoEnabled`、`AccessPolicyEnabled`、`SamlData`、`Tags`、`DnsServers`。如果不返回，需要标记为 Computed 并依赖用户配置。

2. **AccessPolicyEnabled 的具体策略配置方式？**  
   当前仅支持开关（布尔值）。具体的访问策略可能需要通过其他资源或接口配置，超出本次变更范围。

3. **Tags 是否需要 diff suppress？**  
   需要测试 Tags 的大小写敏感性和空格处理。如果 API 规范化了 key/value，可能需要添加 DiffSuppressFunc。
