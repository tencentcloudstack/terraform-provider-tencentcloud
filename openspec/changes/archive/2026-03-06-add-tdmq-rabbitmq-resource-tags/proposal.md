# Change: Add Resource Tags Support to TDMQ RabbitMQ VIP Instance

## Why

Currently, the `tencentcloud_tdmq_rabbitmq_vip_instance` resource lacks support for resource tags (资源标签). According to the latest Tencent Cloud API updates:

1. The `CreateRabbitMQVipInstance` API now includes a `ResourceTags` parameter for setting tags during instance creation
2. The `DescribeRabbitMQVipInstances` API returns tags via `rabbitmqVipInstance.ClusterInfo.Tags`
3. The `ModifyRabbitMQVipInstance` API supports updating tags through the `Tags` parameter

This gap prevents users from properly managing their RabbitMQ instances with Tencent Cloud's tag-based resource organization, cost allocation, and access control features.

## What Changes

Add comprehensive resource tags support to `tencentcloud_tdmq_rabbitmq_vip_instance`, including:

- **Schema Addition**: Add `resource_tags` field (TypeMap) to the resource schema
- **Create Operation**: Pass `resource_tags` to `CreateRabbitMQVipInstance` API's `ResourceTags` parameter
- **Read Operation**: Read tags from `DescribeRabbitMQVipInstances` response's `ClusterInfo.Tags` field and populate `resource_tags`
- **Update Operation**: Support tag modifications through `ModifyRabbitMQVipInstance` API's `Tags` parameter
- **Documentation**: Update resource documentation with usage examples and field descriptions

## Impact

### Affected Files
- **Resource Definition**: `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go`
  - Schema update (add `resource_tags` field)
  - Create function (handle ResourceTags)
  - Read function (populate resource_tags from API response)
  - Update function (handle tag modifications)
- **Documentation**: `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown`
  - Add field description
  - Add usage examples
- **Changelog**: `.changelog/<PR_NUMBER>.txt`

### API Mappings
| Operation | TF Field | API Field | API Endpoint |
|-----------|----------|-----------|--------------|
| Create | `resource_tags` | `ResourceTags` (Array of Tag) | `CreateRabbitMQVipInstance` |
| Read | `resource_tags` | `ClusterInfo.Tags` (Array of Tag) | `DescribeRabbitMQVipInstances` |
| Update | `resource_tags` | `Tags` (Array of Tag) | `ModifyRabbitMQVipInstance` |

### Backward Compatibility
- ✅ **Non-breaking**: `resource_tags` is Optional, existing configurations continue to work
- ✅ **Default Behavior**: Empty tags if not specified
- ✅ **State Migration**: Not required (new field)

### Testing Requirements
- Verify tag creation during instance creation
- Verify tag reading during refresh
- Verify tag updates
- Verify tag removal (set to empty map)
- Verify state consistency after CRUD operations

### API Documentation References
- Create API: https://cloud.tencent.com/document/api/1179/88134
- Describe API: https://cloud.tencent.com/document/api/1179/82205
- Modify API: https://cloud.tencent.com/document/api/1179/88450
