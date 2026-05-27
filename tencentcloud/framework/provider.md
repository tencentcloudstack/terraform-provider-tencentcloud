The TencentCloud provider exposes a second stack built on
terraform-plugin-framework. The framework stack lives inside the same
provider binary as the SDKv2 stack (combined via `tf5muxserver`). New
resource types — Resources, Data Sources, Functions, Ephemeral
Resources, List Resources and Actions — are implemented in this stack.

Resources List
