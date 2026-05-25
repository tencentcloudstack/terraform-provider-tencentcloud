The TencentCloud provider exposes a second stack built on
terraform-plugin-framework. The framework stack lives inside the same
provider binary as the SDKv2 stack (combined via `tf5muxserver`). New
resource types — Resources, Data Sources, Functions, Ephemeral
Resources, List Resources and Actions — are implemented in this stack.

Resources List

Provider Meta
Data Source
tencentcloud_provider_runtime

Resource
tencentcloud_local_note

Function
tencentcloud_parse_resource_id

Ephemeral Resource
tencentcloud_temp_credential

List Resource
tencentcloud_region

Cloud Virtual Machine(CVM)
Action
tencentcloud_reboot_instance
