---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_proxy"
sidebar_current: "docs-tencentcloud-resource-mysql_proxy"
description: |-
  Provides a resource to create a mysql proxy
---

# tencentcloud_mysql_proxy

Provides a resource to create a mysql proxy

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cdb"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-mysql"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-mysql"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "security_group" {
  name        = "sg-mysql"
  description = "mysql test"
}

resource "tencentcloud_mysql_instance" "example" {
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord123"
  slave_deploy_mode = 1
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  first_slave_zone  = data.tencentcloud_availability_zones_by_product.zones.zones.1.name
  slave_sync_mode   = 1
  instance_name     = "tf-example-mysql"
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  intranet_port     = 3306
  security_groups   = [tencentcloud_security_group.security_group.id]

  tags = {
    name = "test"
  }

  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
}

resource "tencentcloud_mysql_proxy" "example" {
  instance_id    = tencentcloud_mysql_instance.example.id
  uniq_vpc_id    = tencentcloud_vpc.vpc.id
  uniq_subnet_id = tencentcloud_subnet.subnet.id
  proxy_node_custom {
    node_count = 1
    cpu        = 2
    mem        = 4000
    region     = "ap-guangzhou"
    zone       = "ap-guangzhou-3"
  }
  security_group        = [tencentcloud_security_group.security_group.id]
  desc                  = "desc."
  connection_pool_limit = 2
  vip                   = "10.0.0.120"
  vport                 = 3306
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `proxy_node_custom` - (Required, List) Node specification configuration.
* `uniq_subnet_id` - (Required, String) Subnet id.
* `uniq_vpc_id` - (Required, String) Vpc id.
* `connection_pool_limit` - (Optional, Int) Connection Pool Threshold.
* `desc` - (Optional, String) Describe.
* `proxy_version` - (Optional, String) The current version of the database agent. No need to fill in when creating.
* `security_group` - (Optional, Set: [`String`]) Security group.
* `upgrade_time` - (Optional, String) Upgrade time: nowTime (upgrade completed) timeWindow (instance maintenance time), Required when modifying the agent version, No need to fill in when creating.
* `vip` - (Optional, String) IP address.
* `vport` - (Optional, Int) Port.

The `proxy_node_custom` object supports the following:

* `cpu` - (Required, Int) Number of CPU cores.
* `mem` - (Required, Int) Memory size.
* `node_count` - (Required, Int) Number of nodes.
* `region` - (Required, String) Region.
* `zone` - (Required, String) Zone.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `proxy_address_id` - Proxy address id.
* `proxy_group_id` - Proxy group id.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCdbProxy](https://cloud.tencent.com/document/api/236/90587) - 创建数据库代理
* Create: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - 查询异步任务的执行结果
* Create: [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) - 查询数据库代理详情
* Create: [ModifyCdbProxyAddressVipAndVPort](https://cloud.tencent.com/document/api/236/90582) - 修改数据库代理地址VPC
* Read: [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) - 查询数据库代理详情
* Read: [DescribeDBSecurityGroups](https://cloud.tencent.com/document/api/236/15854) - 查询实例安全组信息
* Update: [AdjustCdbProxy](https://cloud.tencent.com/document/api/236/90590) - 调整数据库代理
* Update: [DescribeAsyncRequestInfo](https://cloud.tencent.com/document/api/236/20410) - 查询异步任务的执行结果
* Update: [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) - 查询数据库代理详情
* Update: [ModifyCdbProxyAddressDesc](https://cloud.tencent.com/document/api/236/90583) - 修改代理地址描述
* Update: [ModifyCdbProxyAddressVipAndVPort](https://cloud.tencent.com/document/api/236/90582) - 修改数据库代理地址VPC
* Update: [ModifyCdbProxyParam](https://cloud.tencent.com/document/api/236/90581) - 配置数据库代理参数
* Update: [UpgradeCDBProxyVersion](https://cloud.tencent.com/document/api/236/68160) - 升级数据库代理版本
* Delete: [CloseCDBProxy](https://cloud.tencent.com/document/api/236/68171) - 关闭数据库代理
* Delete: [DescribeCdbProxyInfo](https://cloud.tencent.com/document/api/236/90585) - 查询数据库代理详情

