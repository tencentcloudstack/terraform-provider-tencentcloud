## Why

The `tencentcloud_tdmq_namespace` resource currently does not support tagging namespaces during creation or updating tags after creation. The TDMQ `CreateEnvironment` API already supports a `Tags` parameter, and the `DescribeEnvironments` API returns tags in the `Environment` response. Users need the ability to assign and update tags on TDMQ namespaces for resource organization and management purposes.

## What Changes

- Add a `Tags` parameter (TypeList of key/value pairs) to the `tencentcloud_tdmq_namespace` resource schema
- The `Tags` parameter will be Optional (mutable, supports updates)
- Pass `Tags` to the `CreateEnvironment` API request during resource creation
- Read `Tags` from the `DescribeEnvironments` API response during resource read
- In the Update method, if tags change, use the tag service's `ModifyTags` (which calls `TagResources`/`UnTagResources` APIs) to add/remove/update tags
- Update the service layer `CreateTdmqNamespace` function to accept and pass Tags
- Update the resource `.md` documentation example

## Capabilities

### New Capabilities
- `tdmq-namespace-tags`: Add Tags parameter support to tencentcloud_tdmq_namespace resource, enabling tag assignment during namespace creation and tag updates via TagResources/UnTagResources APIs

### Modified Capabilities

## Impact

- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.go` - Add Tags schema, create/read/update logic (tag update uses svctag.ModifyTags)
- `tencentcloud/services/tdmq/service_tencentcloud_tdmq.go` - Update `CreateTdmqNamespace` to accept Tags parameter
- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace.md` - Update documentation with Tags example
- `tencentcloud/services/tpulsar/resource_tc_tdmq_namespace_test.go` - Add unit tests for Tags parameter
