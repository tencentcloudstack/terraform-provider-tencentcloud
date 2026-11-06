---
subcategory: "Agent Sandbox(AGS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-resource-ags_sandbox_tool"
description: |-
  Provides an AGS SandboxTool configuration for applications that create short-lived sandboxes with the runtime SDK. Terraform does not start, stop or connect to instances.
---

# tencentcloud_ags_sandbox_tool

Provides an AGS SandboxTool configuration for applications that create short-lived sandboxes with the runtime SDK. Terraform does not start, stop or connect to instances.

This is a pre-release integration. The first published Provider version and supported region/type matrix must be recorded after acceptance; no released AGS Provider version is claimed here. The current API must pass the applicable collection reset, permission-safe absence and complete readback checks before production use. For an already managed Resource, a complete successful empty exact-ID query removes its state so a later apply can recreate it, following the existing Provider resource convention. A successful permission-filtered empty response is indistinguishable under this policy; use the intended account and region. Explicit authorization/network errors and incomplete responses retain state. Provisional CREATING identities remain conservative until visibility is established. Import of a missing Tool still fails without creating or adopting anything. After a successful Delete with a valid public request receipt, the same operation waits for either explicit NotFound or a complete successful empty query of that exact ID. An unacknowledged Delete, failed confirmation, cancellation or timeout retains state. Description clearing remains a documented server defect accepted as nonblocking for this local validation; it is not reported as working.

Fixed configuration changes replace the Tool. A same-name replacement can fail with create_before_destroy. Active and paused instances can prevent deletion: release them through their runtime owner first. For a gradual migration, create a Tool with a new name, switch application references, release old instances and then remove the old Tool. Do not manage one Tool or its embedded bindings from multiple resources.

Enable AGS and use a supported region. Reuse the Provider's public authentication settings, environment or role identity; do not configure an AGS API key. The Tool's role_arn is a data-access role and is distinct from the credentials used by Terraform. Custom registry access can require an existing role_arn even when only configuring a Tool; the resource does not create a role or grant IAM permissions. Existing subnets, security groups, storage, roles, images and CLS Topics must have the required permissions and connectivity. Only browser, code-interpreter and custom are in V1 scope. Image construction, third-party registry credentials, DNS, security contexts, independent volumes and instance pools are outside this resource.

Network, role, storage, log, persistent and default-timeout changes replace the Tool in this V1 implementation. Custom image/specification/startup settings, description and user tags update in place after verified readback. Settings are inherited according to the API's instance-creation behavior; Terraform never rolls existing instances. Tag changes can affect quota attribution of existing instances.

Custom startup changes can propagate to instance creation after the Tool API already reports the new configuration. Applications should verify a new instance's image and startup settings before sending work. Terraform verifies Tool readback; it does not establish runtime readiness, launch probe instances or restart existing instances. No fixed startup propagation-time guarantee is claimed.

Omitted Optional/Computed values accept the current remote value. Explicit empty description, tags, env, args and command request reset; the operation fails if the API cannot converge. Removing previously declared ports explicitly clears them. Removing a whole probe replaces the Tool until a public reset contract exists. The immutable qcs:tag:createdBy system tag is not a configurable user tag and remains visible in effective_configuration.tags. Environment values are stored in Terraform state and must not contain secrets.

The timeouts block controls create/update/delete waiting budgets (30m by default) and read waiting budget (5m by default), independently of default_timeout, which controls future instance lifetime. Creation may preserve an accepted ID while returning an error. Verify status and configuration before removing taint or importing; a ClientToken is finite-lived idempotency evidence, not a permanent resource index.

## Example Usage

```hcl
resource "tencentcloud_ags_sandbox_tool" "code" {
  tool_name       = "terraform-code-example"
  tool_type       = "code-interpreter"
  description     = "Code environment managed by Terraform"
  default_timeout = "5m"

  network_configuration {
    network_mode = "PUBLIC"
  }

  tags = {
    environment = "development"
  }

  timeouts {
    create = "30m"
    update = "30m"
    delete = "30m"
  }
}

output "sandbox_tool_id" {
  value = tencentcloud_ags_sandbox_tool.code.id
}
```

## Argument Reference

The following arguments are supported:

* `tool_name` - (Required, String) Unique Tool name, 1-50 letters, digits, underscores or hyphens. Changes replace the Tool.
* `tool_type` - (Required, String) V1 Tool type: browser, code-interpreter or custom. Changes replace the Tool.
* `custom_configuration` - (Optional, List of Object) Exactly one block for custom Tools; forbidden for built-in types. Uses prebuilt, AGS-compatible enterprise/personal registry images. Changes affect newly created instances according to the service contract.
* `default_timeout` - (Optional, String) Default instance lifetime, expressed as a positive whole-second duration up to 24h. Omission accepts the remote default. Changes replace the Tool; equivalent durations do not.
* `description` - (Optional, String) Tool description, up to 200 characters. Empty string requests clearing; omitted values accept the current remote value. The API must support explicit clearing.
* `log_configuration` - (Optional, List of Object) Optional CLS log delivery; changes replace the Tool. Existing Topic and role permissions must be provided.
* `network_configuration` - (Optional, List of Object) Exactly one network block. V1 conservatively replaces the Tool for network changes, including VPC mode and bindings.
* `persistent` - (Optional, Bool) Default persistent mode for future instances. Does not start an instance. Type support must be verified; changes replace the Tool.
* `role_arn` - (Optional, String) Existing CAM data-access role ARN, distinct from Provider credentials. Omission accepts the current role. Changes replace the Tool.
* `storage_mounts` - (Optional, Set of Object) Existing storage mounts. Changes replace the Tool; the Tool never creates or deletes the storage itself.
* `tags` - (Optional, Map) User-managed tags. {} removes all user tags. The documented immutable qcs:tag:createdBy system tag is read-only and remains visible in effective_configuration.tags.
* `timeouts` - (Optional, List of Object) Optional operation waiting budgets; at most one block.

The `custom_configuration` object supports the following:

* `image_registry_type` - (Required, String) enterprise for TCR or personal for CCR. Third-party credential configuration is not supported.
* `image` - (Required, String) Prebuilt image reference. In-place update; immutable references are recommended.
* `args` - (Optional, List) Ordered command arguments. Omission preserves the remote setting; [] explicitly clears the override.
* `command` - (Optional, List) Ordered command override. Omission preserves the remote setting; [] explicitly removes the override when supported.
* `env` - (Optional, Map) Non-sensitive environment variables. Omission preserves remote values; {} clears user overrides. Values are stored in Terraform state.

The `ports` object of `custom_configuration` supports the following:

* `name` - (Required, String) Unique port name.
* `port` - (Required, Int) Business port, 1-65535, excluding platform-reserved ports.
* `protocol` - (Required, String) Supported port protocol.

The `probe` object of `custom_configuration` supports the following:

* `failure_threshold` - (Optional, Int) Probe failure_threshold (positive count). Omission accepts the service default.
* `probe_period_ms` - (Optional, Int) Probe probe_period_ms (positive milliseconds). Omission accepts the service default.
* `probe_timeout_ms` - (Optional, Int) Probe probe_timeout_ms (positive milliseconds). Omission accepts the service default.
* `ready_timeout_ms` - (Optional, Int) Probe ready_timeout_ms (positive milliseconds). Omission accepts the service default.
* `success_threshold` - (Optional, Int) Probe success_threshold (positive count). Omission accepts the service default.

The `resources` object of `custom_configuration` supports the following:

* `cpu` - (Required, String) Positive supported CPU quantity, for example 1 or 1000m; values must satisfy the regional support matrix.
* `memory` - (Required, String) Positive supported memory quantity, for example 1Gi; values must satisfy the regional support matrix.

The `log_configuration` object supports the following:

* `files` - (Required, Set) One to ten files directly under /logs/; subdirectories are unsupported.
* `topic_id` - (Required, String) Existing CLS Topic ID.

The `network_configuration` object supports the following:

* `network_mode` - (Required, String) PUBLIC, SANDBOX or VPC.

The `vpc_config` object of `network_configuration` supports the following:

* `security_group_ids` - (Required, Set) Nonempty set of existing security groups compatible with the selected VPC.
* `subnet_ids` - (Required, Set) Nonempty set of existing subnets in one VPC. Reordering is not a change.

The `storage_mounts` object supports the following:

* `mount_path` - (Required, String) Absolute path inside the instance.
* `name` - (Required, String) Unique mount name.
* `read_only` - (Optional, Bool) Read-only mount; defaults to false when omitted in a declared mount.

The `cfs` object of `storage_mounts` supports the following:

* `file_system_id` - (Required, String) Existing CFS filesystem ID.
* `path` - (Required, String) CFS source path.

The `cos` object of `storage_mounts` supports the following:

* `bucket_name` - (Required, String) Existing bucket name.
* `bucket_path` - (Required, String) Absolute path within the bucket.
* `endpoint` - (Required, String) Existing COS endpoint without embedded credentials.

The `timeouts` object supports the following:

* `create` - (Optional, String) create operation waiting budget; defaults to 30m. Distinct from instance lifetime.
* `delete` - (Optional, String) delete operation waiting budget; defaults to 30m. Distinct from instance lifetime.
* `read` - (Optional, String) read operation waiting budget; defaults to 5m. Distinct from instance lifetime.
* `update` - (Optional, String) update operation waiting budget; defaults to 30m. Distinct from instance lifetime.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `effective_configuration` - Supported actual configuration, including remote defaults and system tags. Does not grant configuration ownership to omitted blocks.
* `id` - Tool ID consumed by the runtime SDK and used for import.
* `image_digest` - Resolved image digest; recalculated when image changes.
* `status_reason` - Service-reported status reason.
* `status` - Actual CREATING, ACTIVE, DELETING or FAILED status.

## API Reference

Create: [CreateSandboxTool](https://cloud.tencent.com/document/product/1814/124812) and [DescribeSandboxToolList](https://cloud.tencent.com/document/product/1814/124810).

Read: [DescribeSandboxToolList](https://cloud.tencent.com/document/product/1814/124810).

Update: [UpdateSandboxTool](https://cloud.tencent.com/document/api/1814/124809) and DescribeSandboxToolList.

Delete: [DeleteSandboxTool](https://cloud.tencent.com/document/api/1814/124811) and DescribeSandboxToolList. All actions use API version 2025-09-20.

## Import

Prepare complete HCL for a supported existing Tool, select its owning account and region, then import its actual Tool ID. Import and refresh are read-only and fail if configuration cannot be represented completely. Review the resulting plan before applying. Unsupported custom DNS/security settings, nonempty custom disk storage and non-COS/CFS mounts must not be silently discarded. V1 supports custom CPU/memory inputs only; a remote Resources.Storage override requires a schema that can represent it before adoption.

```shell
terraform import tencentcloud_ags_sandbox_tool.code sdt-xxxxxxxx
terraform plan -detailed-exitcode
```
