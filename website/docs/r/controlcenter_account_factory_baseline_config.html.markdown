---
subcategory: "ControlCenter"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_controlcenter_account_factory_baseline_config"
sidebar_current: "docs-tencentcloud-resource-controlcenter_account_factory_baseline_config"
description: |-
  Provides a resource to create a Controlcenter account factory baseline config
---

# tencentcloud_controlcenter_account_factory_baseline_config

Provides a resource to create a Controlcenter account factory baseline config

## Example Usage

```hcl
resource "tencentcloud_controlcenter_account_factory_baseline_config" "example" {
  name = "default"
  baseline_config_items {
    identifier = "TCC-AF_VPC_SUBNET"
    configuration = jsonencode({
      "VpcName" : "tf-example",
      "CidrBlock" : "10.0.0.0/16",
      "Region" : "1",
      "RegionName" : "ap-guangzhou",
      "Subnets" : [
        {
          "CidrBlock" : "10.0.0.0/24",
          "SubnetName" : "abc",
          "Zone" : "ap-guangzhou-6"
        }
      ]
    })
  }

  baseline_config_items {
    identifier    = "TCC-AF_PRESET_TAG"
    configuration = "{\"TagValuePairs\":[{\"Key\":\"key\",\"Values\":[\"value\"]}]}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Baseline name, which must be unique. Supports only English letters, numbers, Chinese characters, and symbols @, &, _, [], -. Combination of 1-25 Chinese or English characters.
* `baseline_config_items` - (Optional, Set) Baseline configuration, overwrite update. You can query existing baseline configurations via controlcenter:GetAccountFactoryBaseline. You can query supported baseline lists via controlcenter:ListAccountFactoryBaselineItems.

The `baseline_config_items` object supports the following:

* `configuration` - (Optional, String) Account factory baseline item configuration, different baseline items have different configuration parameters.
* `identifier` - (Optional, String) Specifies the unique identifier for account factory baseline item, can only contain `english letters`, `digits`, and `@,._[]-:()()[]+=.`, with a length of 2-128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time.
* `update_time` - Update time.


## Import

Controlcenter account factory baseline config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_controlcenter_account_factory_baseline_config.example nMtrLC9IuQq27wyiICj9bA==
```

