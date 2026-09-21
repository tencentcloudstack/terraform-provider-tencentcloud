---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_api_token"
sidebar_current: "docs-tencentcloud-resource-teo_inference_api_token"
description: |-
  Provides a resource to create a TEO inference API token.
---

# tencentcloud_teo_inference_api_token

Provides a resource to create a TEO inference API token.

## Example Usage

```hcl
resource "tencentcloud_teo_inference_api_token" "example" {
  zone_id = "zone-2qtuhspy7cr6"
  name    = "my-inference-token"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Name of the inference API token. Max length: 30 characters.
* `zone_id` - (Required, String, ForceNew) Site ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `content` - Content of the inference API token.
* `create_time` - Creation time of the inference API token in ISO 8601 format.
* `token_id` - ID of the inference API token.


## Import

TEO inference API token can be imported using the joint id "zone_id#token_id", e.g.

```
terraform import tencentcloud_teo_inference_api_token.example zone-2qtuhspy7cr6#token-abcdefghij
```

