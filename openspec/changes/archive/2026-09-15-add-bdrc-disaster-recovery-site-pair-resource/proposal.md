## Why

BDRC (Backup and Disaster Recovery Center) provides cross-region/cross-zone/cross-cloud disaster recovery for cloud disks, CFS, and CVM instances. Currently there is no Terraform resource to manage disaster recovery site pairs (`容灾站点对`), forcing users to create and manage them via the console or CLI. Adding a `tencentcloud_bdrc_disaster_recovery_site_pair` resource enables full IaC lifecycle management (create, read, update, delete) for site pairs.

## What Changes

- Add a new Terraform resource `tencentcloud_bdrc_disaster_recovery_site_pair` of type RESOURCE_KIND_GENERAL under `tencentcloud/services/bdrc/`
- Implement CRUD operations backed by four BDRC cloud APIs:
  - **Create**: `CreateDisasterRecoverySitePair` — creates a site pair, returns `SitePairId`
  - **Read**: `DescribeDisasterRecoverySitePairs` — queries the site pair list by `SitePairId` to read back all attributes
  - **Update**: `ModifySitePairAttribute` — modifies the site pair name (`SitePairName`); all other fields are immutable (ForceNew)
  - **Delete**: `DeleteDisasterRecoverySitePairs` — deletes the site pair by `SitePairId`
- Register the resource in `provider.go` and `provider.md`
- Add the BDRC SDK client connection method `UseBdrcV20260330Client()` in `connectivity/client.go`
- Add a service layer file `service_tencentcloud_bdrc.go` with a `DescribeDisasterRecoverySitePairById` helper
- Create the resource documentation `resource_tc_bdrc_disaster_recovery_site_pair.md`
- Create unit tests `resource_tc_bdrc_disaster_recovery_site_pair_test.go` using gomonkey mocks

## Capabilities

### New Capabilities

- `bdrc-disaster-recovery-site-pair-resource`: Provides a Terraform resource to create, read, update, and delete BDRC disaster recovery site pairs, covering the full CRUD lifecycle via `CreateDisasterRecoverySitePair`, `DescribeDisasterRecoverySitePairs`, `ModifySitePairAttribute`, and `DeleteDisasterRecoverySitePairs`.

### Modified Capabilities

<!-- None — this is a brand-new resource with no existing spec changes. -->

## Impact

- **New files**:
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair.go`
  - `tencentcloud/services/bdrc/service_tencentcloud_bdrc.go`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair.md`
  - `tencentcloud/services/bdrc/resource_tc_bdrc_disaster_recovery_site_pair_test.go`
- **Modified files**:
  - `tencentcloud/connectivity/client.go` — add BDRC SDK import + client field + `UseBdrcV20260330Client()` method
  - `tencentcloud/provider.go` — import `bdrc` service package and register `tencentcloud_bdrc_disaster_recovery_site_pair`
  - `tencentcloud/provider.md` — add the new resource entry
- **SDK**: `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bdrc/v20260330` (already vendored)
- **Backward compatibility**: Fully backward compatible — only new resource registration, no changes to existing resources or schemas.
