# 变更:为 TCR Namespace 资源添加标签更新支持

## 为什么

`tencentcloud_tcr_namespace` 资源目前在创建和读取操作中支持标签,但缺少标签更新的支持。当用户在 Terraform 配置中修改标签时,这些变更不会应用到实际的 TCR namespace 资源,因为 `ModifyNamespace` API 不支持标签参数。

为了实现标签更新,我们需要在 Update 操作中使用腾讯云统一的标签管理 API(`TagService.ModifyTags`),这与其他资源(如 `tencentcloud_clb_target_group`)使用的模式相同。

### SDK 信息

通过查看 TCR SDK (`tencentcloud-sdk-go/tencentcloud/tcr/v20190924/models.go`):

```go
// CreateNamespaceRequest - 支持 TagSpecification
type CreateNamespaceRequestParams struct {
    RegistryId        *string          `json:"RegistryId,omitnil,omitempty"`
    NamespaceName     *string          `json:"NamespaceName,omitnil,omitempty"`
    IsPublic          *bool            `json:"IsPublic,omitnil,omitempty"`
    TagSpecification  *TagSpecification `json:"TagSpecification,omitnil,omitempty"` // ✅ 支持标签
    IsAutoScan        *bool            `json:"IsAutoScan,omitnil,omitempty"`
    // ...
}

// ModifyNamespaceRequest - 不支持标签参数
type ModifyNamespaceRequestParams struct {
    RegistryId        *string `json:"RegistryId,omitnil,omitempty"`
    NamespaceName     *string `json:"NamespaceName,omitnil,omitempty"`
    IsPublic          *bool   `json:"IsPublic,omitnil,omitempty"`
    IsAutoScan        *bool   `json:"IsAutoScan,omitnil,omitempty"`
    IsPreventVUL      *bool   `json:"IsPreventVUL,omitnil,omitempty"`
    Severity          *string `json:"Severity,omitnil,omitempty"`
    CVEWhitelistItems []*CVEWhitelistItem `json:"CVEWhitelistItems,omitnil,omitempty"`
    // ❌ 无 TagSpecification 字段
}
```

因此需要使用统一的标签管理 API 来实现标签更新。

## 变更内容

- 在 `tencentcloud_tcr_namespace` 资源 schema 中添加 `tags` 字段 (TypeMap, Optional, Computed)
- 在 `resourceTencentCloudTcrNamespaceUpdate` 中实现标签更新逻辑,使用统一的标签服务
- 导入 `svctag` 包用于标签管理
- 使用 `d.HasChange("tags")` 处理标签变更检测
- 使用 `svctag.DiffTags` 计算标签差异
- 使用正确的资源命名调用 `tagService.ModifyTags`: `tcr:namespace:{region}:{instanceId}/{namespaceName}`
- 更新资源文档,包含标签使用示例

## 影响范围

- **影响的规范**: `tcr-namespace` (新增能力)
- **影响的代码**:
  - `tencentcloud/services/tcr/resource_tc_tcr_namespace.go` - 添加 tags schema 和更新逻辑
  - `tencentcloud/services/tcr/resource_tc_tcr_namespace.md` - 添加文档示例
- **破坏性变更**: 无 (添加可选字段)
- **迁移需求**: 不适用 (向后兼容的变更)
