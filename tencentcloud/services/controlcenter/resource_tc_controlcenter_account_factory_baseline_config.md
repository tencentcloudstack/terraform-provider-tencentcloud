Provides a resource to create a Controlcenter account factory baseline config

Example Usage

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

Import

Controlcenter account factory baseline config can be imported using the customId(like uuid or base64 string), e.g.

```
terraform import tencentcloud_controlcenter_account_factory_baseline_config.example nMtrLC9IuQq27wyiICj9bA==
```
