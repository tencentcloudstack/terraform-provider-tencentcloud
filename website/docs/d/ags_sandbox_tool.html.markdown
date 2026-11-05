---
subcategory: "Agent Sandbox(AGS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ags_sandbox_tool"
sidebar_current: "docs-tencentcloud-data_source-ags_sandbox_tool"
description: |-
  Reads exactly one existing AGS SandboxTool by Tool ID or exact name without taking ownership. Specify exactly one selector. Zero visible matches, ambiguous matches, incomplete data and authorization errors are reported; a zero-result name query does not prove physical deletion. Name queries never use a prefix match or select the first result.
---

# tencentcloud_ags_sandbox_tool

Reads exactly one existing AGS SandboxTool by Tool ID or exact name without taking ownership. Specify exactly one selector. Zero visible matches, ambiguous matches, incomplete data and authorization errors are reported; a zero-result name query does not prove physical deletion. Name queries never use a prefix match or select the first result.

This is a pre-release integration; pin the accepted published Provider version when one is released. Enable AGS in the selected account and region and reuse the Provider's public authentication settings. Provider credentials and the Tool's data-access role have different purposes. The data source does not require import, mutate the Tool, start instances or expose instance connection tokens.

The configuration output describes supported actual settings and remote defaults, including system tags. Ordinary environment values may be present in state and are not secret storage. Use id as input to the runtime SDK; use the Resource and explicit import only when ongoing configuration ownership is intended. Queries use a bounded five-minute read budget.

## Example Usage

```hcl
variable "existing_tool_id" {
  type = string
}

data "tencentcloud_ags_sandbox_tool" "existing" {
  tool_id = var.existing_tool_id
}

output "sandbox_tool_id" {
  value = data.tencentcloud_ags_sandbox_tool.existing.id
}
```

## Argument Reference

The following arguments are supported:

* `tool_id` - (Optional, String) Exact Tool ID selector. Specify exactly one of tool_id or tool_name; the unused selector is not populated from the result.
* `tool_name` - (Optional, String) Exact Tool name selector, never a prefix or fuzzy match.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `configuration` - Supported read-only actual configuration, including remote defaults and system tags. Environment values are stored in Terraform state; use non-sensitive variables only.
* `id` - Actual Tool ID suitable for the runtime SDK.
* `image_digest` - Resolved custom image digest when available.
* `status_reason` - Service-reported status reason.
* `status` - Actual Tool status, including creating or failed Tools.

## API Reference

Read: [DescribeSandboxToolList](https://cloud.tencent.com/document/product/1814/124810), API version 2025-09-20.
