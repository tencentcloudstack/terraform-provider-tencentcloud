## Why

The current RabbitMQ VIP instance update logic is overly restrictive, marking many fields as immutable that should be updatable. This prevents users from managing their instances efficiently through Terraform, particularly for common operations like modifying auto-renew settings, adjusting storage, or updating cluster names. The update logic also lacks proper state synchronization and retry mechanisms that match the robust error handling in create operations.

## What Changes

- **Enhance Update Capabilities**: Enable updates for fields that are actually supported by the Tencent Cloud API, including:
  - `auto_renew_flag` - Allow users to toggle auto-renewal settings
  - `time_span` - Allow purchase duration updates for prepaid instances
  - `cluster_name` - Already supported, improve implementation
  - `resource_tags` - Already supported, improve implementation

- **Improve Update Logic**:
  - Add proper state synchronization and wait mechanisms after updates
  - Implement differential updates to only send changed fields to API
  - Add retry mechanisms for update operations
  - Improve error handling with better context messages

- **Add Field-Specific Validation**:
  - Validate that immutable fields (truly immutable per API) are correctly identified
  - Add validation for update constraints (e.g., time_span updates only for prepaid instances)

- **Enhanced Documentation**: Update resource documentation to clearly indicate which fields are updatable and any constraints

## Capabilities

### Modified Capabilities
- `tdmq-rabbitmq-vip-instance`: Update requirements to support additional updatable fields and improve update logic with proper state synchronization

## Impact

**Affected Files**:
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance.go` - Main resource implementation
- `tencentcloud/services/trabbit/resource_tc_tdmq_rabbitmq_vip_instance_test.go` - Unit and acceptance tests
- `website/docs/r/tdmq_rabbitmq_vip_instance.html.markdown` - Resource documentation

**API Interactions**:
- `ModifyRabbitMQVipInstance` - Enhanced to support more updateable fields
- `DescribeRabbitMQVipInstances` - Used for post-update state verification

**User Impact**:
- Users can now manage auto-renewal settings and purchase duration through Terraform
- Improved reliability with better error handling and retry mechanisms
- Better state consistency with post-update verification
- Reduced manual intervention for common instance management tasks

**Backward Compatibility**:
- Fully backward compatible - existing configurations will continue to work
- No breaking changes to resource schema or state
