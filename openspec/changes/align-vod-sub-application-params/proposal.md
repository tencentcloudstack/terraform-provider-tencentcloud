# Proposal: Align VOD Sub Application Parameters

## Change ID
`align-vod-sub-application-params`

## Summary
完善 `tencentcloud_vod_sub_application` 资源，支持 `CreateSubAppId` API 的所有参数，包括 `Type`、`Mode`、`StorageRegion` 和 `Tags`，实现与腾讯云 VOD API 的完全对齐。

## Motivation
当前 `tencentcloud_vod_sub_application` 资源仅支持 `CreateSubAppId` API 的部分参数（`Name` 和 `Description`），缺少以下重要参数：

1. **Type**（应用类型）- 支持一体化（AllInOne）和专业版（Professional）两种模式
2. **Mode**（应用模式）- 支持仅 FileID 模式和 FileID & Path 模式
3. **StorageRegion**（存储地域）- 指定媒体文件存储的地域
4. **Tags**（标签）- 用于资源管理和成本分配

缺少这些参数限制了用户在创建子应用时的灵活性，无法满足以下场景：
- 选择专业版应用类型以获得更多功能
- 指定 FileID & Path 模式以支持路径访问
- 选择特定存储地域以优化成本和性能
- 为子应用打标签以进行资源分类和成本管理

## Background

### VOD 子应用概念
腾讯云点播（VOD）的子应用（Sub Application）是主应用下的独立单元，用于实现多租户隔离和资源管理。每个子应用拥有独立的：
- 媒体资源存储空间
- 资源配置和管理
- 数据统计和分析
- API 调用权限

### API 参数说明
根据 SDK 定义（`tencentcloud-sdk-go/tencentcloud/vod/v20180717/models.go`），`CreateSubAppIdRequest` 支持以下参数：

| 参数 | 类型 | 是否必填 | 说明 | 当前支持情况 |
|------|------|---------|------|------------|
| Name | String | 是 | 应用名称，长度限制：40个字符 | ✅ 已支持 |
| Description | String | 否 | 应用简介，长度限制：300个字符 | ✅ 已支持 |
| Type | String | 否 | 应用类型：AllInOne（一体化）或 Professional（专业版），默认 AllInOne | ❌ 未支持 |
| Mode | String | 否 | 应用模式：fileid（仅FileID模式）或 fileid+path（FileID & Path模式），默认 fileid | ❌ 未支持 |
| StorageRegion | String | 否 | 存储地域，如：ap-guangzhou、ap-beijing 等 | ❌ 未支持 |
| Tags | []ResourceTag | 否 | 标签列表，最多10个，用于资源管理 | ❌ 未支持 |

### ResourceTag 结构
```go
type ResourceTag struct {
    TagKey   *string `json:"TagKey,omitnil,omitempty" name:"TagKey"`
    TagValue *string `json:"TagValue,omitnil,omitempty" name:"TagValue"`
}
```

## Proposed Solution

### Schema Changes

#### New Input Arguments
1. **type** (Optional, String)
   - 描述：应用类型
   - 可选值：`AllInOne`（一体化）、`Professional`（专业版）
   - 默认值：`AllInOne`
   - ForceNew：true（修改需要重建资源）

2. **mode** (Optional, String)
   - 描述：应用模式
   - 可选值：`fileid`（仅FileID模式）、`fileid+path`（FileID & Path模式）
   - 默认值：`fileid`
   - ForceNew：true（修改需要重建资源）

3. **storage_region** (Optional, String)
   - 描述：存储地域
   - 示例：`ap-guangzhou`、`ap-beijing`、`ap-shanghai` 等
   - ForceNew：true（修改需要重建资源）

4. **tags** (Optional, Map of String)
   - 描述：标签键值对，用于资源管理
   - 最大数量：10个
   - 键长度限制：128个字符
   - 值长度限制：256个字符
   - 支持更新（通过 Update 函数）

### Implementation Details

#### 1. Schema Definition
```go
Schema: map[string]*schema.Schema{
    // ... existing fields ...
    
    "type": {
        Type:         schema.TypeString,
        Optional:     true,
        ForceNew:     true,
        Default:      "AllInOne",
        ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"AllInOne", "Professional"}),
        Description:  "Sub application type. Valid values: `AllInOne` (all-in-one), `Professional` (professional edition). Default: `AllInOne`.",
    },
    
    "mode": {
        Type:         schema.TypeString,
        Optional:     true,
        ForceNew:     true,
        Default:      "fileid",
        ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"fileid", "fileid+path"}),
        Description:  "Sub application mode. Valid values: `fileid` (FileID only), `fileid+path` (FileID & Path). Default: `fileid`.",
    },
    
    "storage_region": {
        Type:        schema.TypeString,
        Optional:    true,
        ForceNew:    true,
        Description: "Storage region for media files, e.g., `ap-guangzhou`, `ap-beijing`.",
    },
    
    "tags": {
        Type:        schema.TypeMap,
        Optional:    true,
        Elem:        &schema.Schema{Type: schema.TypeString},
        Description: "Tag key-value pairs for resource management. Maximum 10 tags.",
    },
}
```

#### 2. Create Function Enhancement
```go
func resourceTencentCloudVodSubApplicationCreate(d *schema.ResourceData, meta interface{}) error {
    // ... existing code ...
    
    request := vod.NewCreateSubAppIdRequest()
    
    // Existing parameters
    if v, ok := d.GetOk("name"); ok {
        request.Name = helper.String(v.(string))
    }
    if v, ok := d.GetOk("description"); ok {
        request.Description = helper.String(v.(string))
    }
    
    // New parameters
    if v, ok := d.GetOk("type"); ok {
        request.Type = helper.String(v.(string))
    }
    if v, ok := d.GetOk("mode"); ok {
        request.Mode = helper.String(v.(string))
    }
    if v, ok := d.GetOk("storage_region"); ok {
        request.StorageRegion = helper.String(v.(string))
    }
    if v, ok := d.GetOk("tags"); ok {
        tags := v.(map[string]interface{})
        for key, value := range tags {
            tag := vod.ResourceTag{
                TagKey:   helper.String(key),
                TagValue: helper.String(value.(string)),
            }
            request.Tags = append(request.Tags, &tag)
        }
    }
    
    // ... API call ...
}
```

#### 3. Read Function Enhancement
```go
func resourceTencentCloudVodSubApplicationRead(d *schema.ResourceData, meta interface{}) error {
    // ... existing code ...
    
    // Note: DescribeSubAppIds 返回的 SubAppIdInfo 结构中不包含 Type, Mode, StorageRegion
    // 这些字段仅在创建时设置，无法通过 API 查询
    // 因此在 Read 函数中保持这些字段的状态不变（从 Terraform state 读取）
    
    // Existing fields
    _ = d.Set("name", appInfo.Name)
    _ = d.Set("description", appInfo.Description)
    _ = d.Set("status", d.Get("status").(string))
    _ = d.Set("create_time", appInfo.CreateTime)
    
    // Type, Mode, StorageRegion, Tags 保持不变（ForceNew 字段不需要更新）
}
```

#### 4. Update Function Enhancement
```go
func resourceTencentCloudVodSubApplicationUpdate(d *schema.ResourceData, meta interface{}) error {
    // ... existing code for name and description ...
    
    // Handle tags update
    if d.HasChange("tags") {
        // Note: VOD API 可能不支持单独更新 Tags
        // 需要确认是否需要通过 ModifySubAppIdInfo API 或专门的标签 API
        // 暂时假设需要通过统一的标签服务更新
        
        oldTags, newTags := d.GetChange("tags")
        
        // 构造 QCS 资源名称
        resourceName := fmt.Sprintf("qcs::vod:%s:uin/:subapplication/%s", 
            meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region, 
            d.Id())
        
        // 使用 tag service 更新标签
        tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
        replaceTags, deleteTags := svctag.DiffTags(
            oldTags.(map[string]interface{}), 
            newTags.(map[string]interface{}))
        
        if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
            return err
        }
    }
    
    // ... rest of update logic ...
}
```

### Usage Example

#### Basic Example (Existing)
```hcl
resource "tencentcloud_vod_sub_application" "basic" {
  name        = "my-sub-app"
  status      = "On"
  description = "Basic sub application"
}
```

#### Complete Example (New)
```hcl
resource "tencentcloud_vod_sub_application" "complete" {
  name           = "my-professional-app"
  status         = "On"
  description    = "Professional sub application with custom settings"
  type           = "Professional"
  mode           = "fileid+path"
  storage_region = "ap-guangzhou"
  
  tags = {
    "team"        = "media"
    "environment" = "production"
    "project"     = "video-platform"
  }
}
```

## Implementation Tasks

### Phase 1: Schema and Create Function (Core)
1. [ ] 添加新的 schema 字段（type, mode, storage_region, tags）
2. [ ] 更新 Create 函数以支持新参数
3. [ ] 添加参数验证逻辑
4. [ ] 编写基础单元测试

### Phase 2: Read and Update Functions
5. [ ] 更新 Read 函数文档说明（Type/Mode/StorageRegion 不可查询）
6. [ ] 实现 Tags 更新逻辑（如果 API 支持）
7. [ ] 处理状态一致性问题

### Phase 3: Testing and Documentation
8. [ ] 编写验收测试（包含所有新参数）
9. [ ] 更新资源文档（.md 文件）
10. [ ] 添加使用示例
11. [ ] 更新 Changelog

## Testing Strategy

### Unit Tests
- 验证 schema 定义正确性
- 验证参数验证逻辑
- 验证 tags 转换逻辑

### Acceptance Tests
```go
func TestAccTencentCloudVodSubApplication_complete(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:     func() { tcacctest.AccPreCheck(t) },
        Providers:    tcacctest.AccProviders,
        CheckDestroy: testAccCheckVodSubApplicationDestroy,
        Steps: []resource.TestStep{
            {
                Config: testAccVodSubApplicationComplete,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "type", "Professional"),
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "mode", "fileid+path"),
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "storage_region", "ap-guangzhou"),
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.%", "3"),
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.team", "media"),
                ),
            },
            // Test tags update
            {
                Config: testAccVodSubApplicationCompleteUpdate,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.%", "2"),
                    resource.TestCheckResourceAttr("tencentcloud_vod_sub_application.complete", "tags.team", "media-updated"),
                ),
            },
        },
    })
}
```

## Considerations and Risks

### 1. ForceNew Parameters
- `type`, `mode`, `storage_region` 标记为 ForceNew
- 原因：这些参数在创建后无法修改（API 限制）
- 影响：修改这些参数将导致资源重建

### 2. Tags Update
- 需要确认 VOD API 是否支持 Tags 更新
- 如果不支持，Tags 也应标记为 ForceNew
- 可能需要使用统一的标签服务 API

### 3. Read Function Limitation
- `DescribeSubAppIds` 响应不包含 Type, Mode, StorageRegion
- 解决方案：Read 函数中保持这些字段不变（依赖 Terraform state）
- 文档中需要明确说明这个限制

### 4. Backward Compatibility
- 所有新参数都是 Optional
- 默认行为与现有实现一致（AllInOne + fileid 模式）
- 不会破坏现有 Terraform 配置

### 5. API Validation
- 需要通过实际 API 测试验证参数有效性
- 特别是 StorageRegion 的有效值列表
- Tags 的数量和长度限制

## Success Criteria
1. ✅ 所有新参数在 Create 时正确传递给 API
2. ✅ ForceNew 参数修改时触发资源重建
3. ✅ Tags 支持创建和更新（如果 API 支持）
4. ✅ 验收测试全部通过
5. ✅ 文档完整且清晰
6. ✅ 向后兼容现有配置

## Documentation Updates
- [ ] 更新 `resource_tc_vod_sub_application.md`
- [ ] 添加完整的参数说明
- [ ] 添加使用示例（基础和完整）
- [ ] 说明 ForceNew 参数的限制
- [ ] 说明 Read 函数的限制

## References
- VOD SDK Source: `vendor/github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717/models.go`
- Current Implementation: `tencentcloud/services/vod/resource_tc_vod_sub_application.go`
- CreateSubAppIdRequest Structure: Lines 7486-7507
- ResourceTag Structure: Lines 26295-26300
