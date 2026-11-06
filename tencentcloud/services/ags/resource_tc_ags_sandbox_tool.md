Provides an AGS SandboxTool configuration for applications that create short-lived sandboxes with the runtime SDK. Terraform does not start, stop or connect to instances.

This is a pre-release integration. The first published Provider version and supported region/type matrix must be recorded after acceptance; no released AGS Provider version is claimed here. The current API must pass the applicable collection reset, permission-safe absence and complete readback checks before production use. For an already managed Resource, a complete successful empty exact-ID query removes its state so a later apply can recreate it, following the existing Provider resource convention. A successful permission-filtered empty response is indistinguishable under this policy; use the intended account and region. Explicit authorization/network errors and incomplete responses retain state. Provisional CREATING identities remain conservative until visibility is established. Import of a missing Tool still fails without creating or adopting anything. After a successful Delete with a valid public request receipt, the same operation waits for either explicit NotFound or a complete successful empty query of that exact ID. An unacknowledged Delete, failed confirmation, cancellation or timeout retains state. Description clearing remains a documented server defect accepted as nonblocking for this local validation; it is not reported as working.

Fixed configuration changes replace the Tool. A same-name replacement can fail with create_before_destroy. Active and paused instances can prevent deletion: release them through their runtime owner first. For a gradual migration, create a Tool with a new name, switch application references, release old instances and then remove the old Tool. Do not manage one Tool or its embedded bindings from multiple resources.

Enable AGS and use a supported region. Reuse the Provider's public authentication settings, environment or role identity; do not configure an AGS API key. The Tool's role_arn is a data-access role and is distinct from the credentials used by Terraform. Custom registry access can require an existing role_arn even when only configuring a Tool; the resource does not create a role or grant IAM permissions. Existing subnets, security groups, storage, roles, images and CLS Topics must have the required permissions and connectivity. Only browser, code-interpreter and custom are in V1 scope. Image construction, third-party registry credentials, DNS, security contexts, independent volumes and instance pools are outside this resource.

Network, role, storage, log, persistent and default-timeout changes replace the Tool in this V1 implementation. Custom image/specification/startup settings, description and user tags update in place after verified readback. Settings are inherited according to the API's instance-creation behavior; Terraform never rolls existing instances. Tag changes can affect quota attribution of existing instances.

Custom startup changes can propagate to instance creation after the Tool API already reports the new configuration. Applications should verify a new instance's image and startup settings before sending work. Terraform verifies Tool readback; it does not establish runtime readiness, launch probe instances or restart existing instances. No fixed startup propagation-time guarantee is claimed.

Omitted Optional/Computed values accept the current remote value. Explicit empty description, tags, env, args and command request reset; the operation fails if the API cannot converge. Removing previously declared ports explicitly clears them. Removing a whole probe replaces the Tool until a public reset contract exists. The immutable qcs:tag:createdBy system tag is not a configurable user tag and remains visible in effective_configuration.tags. Environment values are stored in Terraform state and must not contain secrets.

The timeouts block controls create/update/delete waiting budgets (30m by default) and read waiting budget (5m by default), independently of default_timeout, which controls future instance lifetime. Creation may preserve an accepted ID while returning an error. Verify status and configuration before removing taint or importing; a ClientToken is finite-lived idempotency evidence, not a permanent resource index.

API

Create: [CreateSandboxTool](https://cloud.tencent.com/document/product/1814/124812) and [DescribeSandboxToolList](https://cloud.tencent.com/document/product/1814/124810).

Read: [DescribeSandboxToolList](https://cloud.tencent.com/document/product/1814/124810).

Update: [UpdateSandboxTool](https://cloud.tencent.com/document/api/1814/124809) and DescribeSandboxToolList.

Delete: [DeleteSandboxTool](https://cloud.tencent.com/document/api/1814/124811) and DescribeSandboxToolList. All actions use API version 2025-09-20.

Example Usage

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

Import

Prepare complete HCL for a supported existing Tool, select its owning account and region, then import its actual Tool ID. Import and refresh are read-only and fail if configuration cannot be represented completely. Review the resulting plan before applying. Unsupported custom DNS/security settings, nonempty custom disk storage and non-COS/CFS mounts must not be silently discarded. V1 supports custom CPU/memory inputs only; a remote Resources.Storage override requires a schema that can represent it before adoption.

```shell
terraform import tencentcloud_ags_sandbox_tool.code sdt-xxxxxxxx
terraform plan -detailed-exitcode
```
