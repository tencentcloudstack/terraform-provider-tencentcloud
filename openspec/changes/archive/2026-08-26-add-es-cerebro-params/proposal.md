## Why

The Tencent Cloud ES `UpdateInstance` API supports Cerebro service configuration (enable/disable, public/private network access, custom private domain), but the `tencentcloud_elasticsearch_instance` Terraform resource does not expose these parameters. Users who need to manage Cerebro service settings through Terraform must currently use the cloud console or API directly, creating a gap in the infrastructure-as-code workflow.

## What Changes

- Add `enable_cerebro` field (TypeBool, Optional, Computed) to enable/disable Cerebro service
- Add `cerebro_public_access` field (TypeString, Optional, Computed) to control Cerebro public network access (`OPEN`/`CLOSE`)
- Add `cerebro_private_access` field (TypeString, Optional, Computed) to control Cerebro private network access (`OPEN`/`CLOSE`)
- Add `cerebro_private_domain` field (TypeString, Optional, Computed) to set Cerebro private network custom domain
- Extend the `UpdateInstance` service layer function to accept and pass through these new parameters
- Handle these fields in the resource Update flow (no changes to Create or Delete flows since Cerebro is not available at creation time)
- Since `InstanceInfo` (DescribeInstances response) does not expose Cerebro fields, these are write-only parameters that preserve user-configured values in state

## Capabilities

### New Capabilities
- `es-instance-cerebro-config`: Configure Cerebro service settings (enable/disable, public/private access, custom domain) on a Tencent Cloud Elasticsearch instance through the `tencentcloud_elasticsearch_instance` resource

### Modified Capabilities
<!-- No existing capabilities are modified -->

## Impact

- **Affected code**: `tencentcloud/services/es/resource_tc_elasticsearch_instance.go`, `tencentcloud/services/es/service_tencentcloud_elasticsearch.go`
- **Affected documentation**: `tencentcloud/services/es/resource_tc_elasticsearch_instance.md`
- **Dependencies**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es v1.1.13` (already available in vendor)
- **Breaking changes**: None (all new fields are Optional)