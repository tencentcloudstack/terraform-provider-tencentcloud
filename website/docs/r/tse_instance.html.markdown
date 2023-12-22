---
subcategory: "Tencent Cloud Service Engine(TSE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tse_instance"
sidebar_current: "docs-tencentcloud-resource-tse_instance"
description: |-
  Provides a resource to create a tse instance
---

# tencentcloud_tse_instance

Provides a resource to create a tse instance

## Example Usage

### Create zookeeper standard version

```hcl
resource "tencentcloud_tse_instance" "zookeeper_standard" {
  engine_type            = "zookeeper"
  engine_version         = "3.5.9.4"
  engine_product_version = "STANDARD"
  engine_region          = "ap-guangzhou"
  engine_name            = "zookeeper-test"
  trade_type             = 0
  engine_resource_spec   = "spec-qvj6k7t4q"
  engine_node_num        = 3
  vpc_id                 = "vpc-4owdpnwr"
  subnet_id              = "subnet-dwj7ipnc"

  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create zookeeper professional version

```hcl
resource "tencentcloud_tse_instance" "zookeeper_professional" {
  engine_type            = "zookeeper"
  engine_version         = "3.5.9.4"
  engine_product_version = "PROFESSIONAL"
  engine_region          = "ap-guangzhou"
  engine_name            = "zookeeper-test"
  trade_type             = 0
  engine_resource_spec   = "spec-qvj6k7t4q"
  engine_node_num        = 3
  vpc_id                 = "vpc-4owdpnwr"
  subnet_id              = "subnet-dwj7ipnc"

  engine_region_infos {
    engine_region = "ap-guangzhou"
    replica       = 3

    vpc_infos {
      subnet_id = "subnet-dwj7ipnc"
      vpc_id    = "vpc-4owdpnwr"
    }
    vpc_infos {
      subnet_id = "subnet-403mgks4"
      vpc_id    = "vpc-b1puef4z"
    }
  }
  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create nacos standard version

```hcl
resource "tencentcloud_tse_instance" "nacos" {
  enable_client_internet_access = false
  engine_name                   = "test"
  engine_node_num               = 3
  engine_product_version        = "STANDARD"
  engine_region                 = "ap-guangzhou"
  engine_resource_spec          = "spec-1160a35a"
  engine_type                   = "nacos"
  engine_version                = "2.0.3.4"
  subnet_id                     = "subnet-5vpegquy"
  trade_type                    = 0
  vpc_id                        = "vpc-99xmasf9"

  tags = {
    "createdBy" = "terraform"
  }
}
```

### Create polaris base version

```hcl
resource "tencentcloud_tse_instance" "polaris" {
  enable_client_internet_access = false
  engine_name                   = "test"
  engine_node_num               = 2
  engine_product_version        = "BASE"
  engine_region                 = "ap-guangzhou"
  engine_resource_spec          = "spec-c160bas1"
  engine_type                   = "polaris"
  engine_version                = "1.16.0.1"
  subnet_id                     = "subnet-5vpegquy"
  trade_type                    = 0
  vpc_id                        = "vpc-99xmasf9"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_name` - (Required, String) engien name. Reference value: nacos-test.
* `engine_product_version` - (Required, String) Engine product version. Reference value: `Nacos`: `TRIAL`: Development version, optional node num: `1`, optional spec list: `1C1G`; `STANDARD`: Standard versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`. `Zookeeper`: `TRIAL`: Development version, optional node num: `1`, optional spec list: `1C1G`; `STANDARD`: Standard versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`; `PROFESSIONAL`: professional versions, optional node num: `3`, `5`, `7`, optional spec list: `1C2G`, `2C4G`, `4C8G`, `8C16G`, `16C32G`. `Polarismesh`: `BASE`: Base version, optional node num: `1`, optional spec list: `NUM50`; `PROFESSIONAL`: Enterprise versions, optional node num: `2`, `3`, optional spec list: `NUM50`, `NUM100`, `NUM200`, `NUM500`, `NUM1000`, `NUM5000`, `NUM10000`, `NUM50000`.
* `engine_region` - (Required, String) engine deploy region. Reference value: `China area` Reference value: `ap-guangzhou`, `ap-beijing`, `ap-chengdu`, `ap-chongqing`, `ap-nanjing`, `ap-shanghai` `ap-beijing-fsi`, `ap-shanghai-fsi`, `ap-shenzhen-fsi`. `Asia Pacific` area Reference value: `ap-hongkong`, `ap-taipei`, `ap-jakarta`, `ap-singapore`, `ap-bangkok`, `ap-seoul`, `ap-tokyo`. `North America area` Reference value: `na-toronto`, `sa-saopaulo`, `na-siliconvalley`, `na-ashburn`.
* `engine_type` - (Required, String) engine type. Reference value: `zookeeper`, `nacos`, `polaris`.
* `engine_version` - (Required, String) An open source version of the engine. Each engine supports different open source versions, refer to the product documentation or console purchase page.
* `trade_type` - (Required, Int) trade type. Reference value:- 0:postpaid- 1:Prepaid (Interface does not support the creation of prepaid instances yet).
* `enable_client_internet_access` - (Optional, Bool) Client public network access, `true`: on, `false`: off, default: false.
* `engine_node_num` - (Optional, Int) engine node num. see EngineProductVersion.
* `engine_region_infos` - (Optional, List) Details about the regional configuration of the engine in cross-region deployment, only zookeeper professional requires the use of the EngineRegionInfos parameter.
* `engine_resource_spec` - (Optional, String) engine spec ID. see EngineProductVersion.
* `prepaid_period` - (Optional, Int) Prepaid time, in monthly units.
* `prepaid_renew_flag` - (Optional, Int) Automatic renewal mark, prepaid only.  Reference value: `0`: No automatic renewal, `1`: Automatic renewal.
* `subnet_id` - (Optional, String) subnet ID. Assign an IP address to the engine in the VPC subnet. Reference value: subnet-ahde9me9.
* `tags` - (Optional, Map) Tag description list.
* `vpc_id` - (Optional, String) VPC ID. Assign an IP address to the engine in the VPC subnet. Reference value: vpc-conz6aix.

The `engine_region_infos` object supports the following:

* `engine_region` - (Required, String) Engine node region.
* `replica` - (Required, Int) The number of nodes allocated in this region.
* `vpc_infos` - (Required, List) Cluster network information.

The `vpc_infos` object of `engine_region_infos` supports the following:

* `subnet_id` - (Required, String) Subnet ID.
* `vpc_id` - (Required, String) Vpc Id.
* `intranet_address` - (Optional, String) Intranet access addressNote: This field may return null, indicating that a valid value is not available..

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tse instance can be imported using the id, e.g.

```
terraform import tencentcloud_tse_instance.instance instance_id
```

