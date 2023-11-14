/*
Provides a resource to create a mariadb network

Example Usage

```hcl
resource "tencentcloud_mariadb_network" "network" {
  instance_id = ""
  vpc_id = ""
  subnet_id = ""
  vip = ""
  vipv6 = ""
  vip_release_delay =
}
```

Import

mariadb network can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_network.network network_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMariadbNetwork() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbNetworkCreate,
		Read:   resourceTencentCloudMariadbNetworkRead,
		Delete: resourceTencentCloudMariadbNetworkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"vpc_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "VpcId, ID of the desired VPC network.",
			},

			"subnet_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "SubnetId, subnet ID of the desired VPC network.",
			},

			"vip": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The field is required to specify VIP.",
			},

			"vipv6": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The field is required to specify VIPv6.",
			},

			"vip_release_delay": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "VIP retention period in hours. Value range: 0-168. Default value: `24`. `0` indicates that the VIP will be released immediately, but there will be 1-minute delay.",
			},
		},
	}
}

func resourceTencentCloudMariadbNetworkCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_network.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewModifyInstanceNetworkRequest()
		response   = mariadb.NewModifyInstanceNetworkResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vpc_id"); ok {
		request.VpcId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		request.SubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vip"); ok {
		request.Vip = helper.String(v.(string))
	}

	if v, ok := d.GetOk("vipv6"); ok {
		request.Vipv6 = helper.String(v.(string))
	}

	if v, _ := d.GetOk("vip_release_delay"); v != nil {
		request.VipReleaseDelay = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().ModifyInstanceNetwork(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb network failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbNetworkRead(d, meta)
}

func resourceTencentCloudMariadbNetworkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_network.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbNetworkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_network.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
