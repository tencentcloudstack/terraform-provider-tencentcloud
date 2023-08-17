/*
Provides a resource to create a mysql reload_balance_proxy_node

Example Usage

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

resource "tencentcloud_mysql_reload_balance_proxy_node" "example" {
  proxy_group_id   = tencentcloud_mysql_proxy.example.proxy_group_id
  proxy_address_id = tencentcloud_mysql_proxy.example.proxy_address_id
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlReloadBalanceProxyNode() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlReloadBalanceProxyNodeCreate,
		Read:   resourceTencentCloudMysqlReloadBalanceProxyNodeRead,
		Delete: resourceTencentCloudMysqlReloadBalanceProxyNodeDelete,

		Schema: map[string]*schema.Schema{
			"proxy_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Proxy id.",
			},

			"proxy_address_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Proxy address id.",
			},
		},
	}
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewReloadBalanceProxyNodeRequest()
		instanceId string
	)
	if v, ok := d.GetOk("proxy_group_id"); ok {
		instanceId = v.(string)
		request.ProxyGroupId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("proxy_address_id"); ok {
		request.ProxyAddressId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().ReloadBalanceProxyNode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql reloadBalanceProxyNode failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlReloadBalanceProxyNodeRead(d, meta)
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlReloadBalanceProxyNodeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_reload_balance_proxy_node.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
