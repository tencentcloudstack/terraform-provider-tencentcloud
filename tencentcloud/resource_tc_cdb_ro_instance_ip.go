/*
Provides a resource to create a cdb ro_instance_ip

Example Usage

```hcl
resource "tencentcloud_cdb_ro_instance_ip" "ro_instance_ip" {
  instance_id = ""
  uniq_subnet_id = ""
  uniq_vpc_id = ""
}
```

Import

cdb ro_instance_ip can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_ro_instance_ip.ro_instance_ip ro_instance_ip_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudCdbRoInstanceIp() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRoInstanceIpCreate,
		Read:   resourceTencentCloudCdbRoInstanceIpRead,
		Delete: resourceTencentCloudCdbRoInstanceIpDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Read-only instance ID, in the format: cdbro-3i70uj0k, which is the same as the read-only instance ID displayed on the cloud database console page.",
			},

			"uniq_subnet_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Subnet descriptor, for example: subnet-1typ0s7d.",
			},

			"uniq_vpc_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Vpc descriptor, for example: vpc-a23yt67j, if this field is passed, UniqSubnetId must be passed.",
			},
		},
	}
}

func resourceTencentCloudCdbRoInstanceIpCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_instance_ip.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = cdb.NewCreateRoInstanceIpRequest()
		response   = cdb.NewCreateRoInstanceIpResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_subnet_id"); ok {
		request.UniqSubnetId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("uniq_vpc_id"); ok {
		request.UniqVpcId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().CreateRoInstanceIp(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb roInstanceIp failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudCdbRoInstanceIpRead(d, meta)
}

func resourceTencentCloudCdbRoInstanceIpRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_instance_ip.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRoInstanceIpDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_instance_ip.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
