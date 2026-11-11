## Why

The CDWDoris `CreateInstanceNew` API supports a `DiskEncrypt` request parameter to enable disk encryption when creating a cluster, and the `DescribeInstance` API returns the encryption status (`Encrypt`) in both `MasterSummary` and `CoreSummary`. However the Terraform resource `tencentcloud_cdwdoris_instance` does not expose disk encryption. Users who want to create encrypted CDWDoris clusters cannot specify disk encryption through Terraform, forcing them to use the console or API directly.

## What Changes

- Add `disk_encrypt` (Optional, immutable after creation) parameter to `tencentcloud_cdwdoris_instance` resource to support enabling disk encryption at creation time. This parameter is passed to the `CreateInstanceNew` API request as `DiskEncrypt`. Changes to `disk_encrypt` after creation are rejected via the immutable args check in the Update function (consistent with the existing `immutableArgs` pattern in this resource).
- Add `encrypt` (Computed) parameter to `tencentcloud_cdwdoris_instance` resource to refresh the disk encryption status returned by the `DescribeInstance` API response. The value is read from `response.InstanceInfo.MasterSummary.Encrypt` (and `CoreSummary.Encrypt`, which carries the same value) and set into state.

## Capabilities

### New Capabilities
- `cdwdoris-instance-disk-encrypt`: Enable the `disk_encrypt` input parameter and `encrypt` computed output on the `tencentcloud_cdwdoris_instance` resource to allow users to create disk-encrypted CDWDoris clusters and refresh the encryption status.

### Modified Capabilities
<!-- No existing specs require modification -->

## Impact

- **Affected files:**
  - `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance.go` — add `disk_encrypt` and `encrypt` schema fields; wire `disk_encrypt` through Create flow (set `request.DiskEncrypt`); add `disk_encrypt` to the `immutableArgs` array in Update; add Read support to refresh `encrypt` from `MasterSummary.Encrypt`/`CoreSummary.Encrypt`
  - `tencentcloud/services/cdwdoris/resource_tc_cdwdoris_instance_test.go` — add unit test cases for the new parameter using gomonkey mocks
- **SDK dependency:** `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdwdoris/v20211228` (already vendored, now contains `DiskEncrypt *int64` on `CreateInstanceNewRequest` and `Encrypt *int64` on `NodesSummary`).
- **Backward compatibility:** fully backward compatible — the new parameter is Optional and defaults to not being set; no existing TF configuration or state is affected.
- **API constraints:** `DiskEncrypt` is only accepted by `CreateInstanceNew` (there is no update API that accepts it), so `disk_encrypt` is immutable after creation. The `Encrypt` field is read from the `DescribeInstance` response (`InstanceInfo.MasterSummary.Encrypt` / `InstanceInfo.CoreSummary.Encrypt`), so it can be refreshed on Read and supports imported resources.