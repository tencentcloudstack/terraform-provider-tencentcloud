## Why

TencentCloud AGS (Agent Sandbox) service requires Terraform support to manage sandbox tools as infrastructure-as-code. The `CreateSandboxTool` API is already available in the AGS SDK (`v20250920`), but no corresponding Terraform resource exists. Adding `tencentcloud_ags_sandbox_tool` resource enables users to provision, update, and delete sandbox tools declaratively through Terraform, aligning with IaC best practices.

## What Changes

- Add new Terraform resource `tencentcloud_ags_sandbox_tool` for managing AGS sandbox tools
- Implement CRUD operations using AGS SDK v20250920 APIs: `CreateSandboxTool`, `DescribeSandboxToolList`, `UpdateSandboxTool`, `DeleteSandboxTool`
- Support all sandbox tool configuration parameters including network configuration, storage mounts, custom configuration, computer configuration, and log configuration
- Support tag management via the AGS Tag API
- Register the resource in `provider.go` and `provider.md`

## Capabilities

### New Capabilities
- `ags-sandbox-tool-resource`: Terraform resource for managing AGS sandbox tools with full CRUD lifecycle support, including network configuration (PUBLIC/VPC modes), storage mounts (COS, CFS, Image, AgentBucket), custom sandbox configuration (image, resources, ports, env, probe, DNS, security context), computer configuration (WAA, OSWorld), log configuration (CLS), and tag management.

### Modified Capabilities
<!-- None -->

## Impact

- **New files**: `tencentcloud/services/ags/resource_tc_ags_sandbox_tool.go`, `tencentcloud/services/ags/resource_tc_ags_sandbox_tool.md`, `tencentcloud/services/ags/resource_tc_ags_sandbox_tool_test.go`
- **Modified files**: `tencentcloud/provider.go` (resource registration), `tencentcloud/provider.md` (resource documentation entry)
- **Dependencies**: Uses existing `github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920` SDK package (already vendored)