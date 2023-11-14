/*
Provides a resource to create a dcdb open_d_b_extranet_access

Example Usage

```hcl
resource "tencentcloud_dcdb_open_d_b_extranet_access" "open_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag =
}
```

Import

dcdb open_d_b_extranet_access can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_open_d_b_extranet_access.open_d_b_extranet_access open_d_b_extranet_access_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudDcdbOpenDBExtranetAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbOpenDBExtranetAccessCreate,
		Read:   resourceTencentCloudDcdbOpenDBExtranetAccessRead,
		Delete: resourceTencentCloudDcdbOpenDBExtranetAccessDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of instance for which to enable public network access. The ID is in the format of `tdsql-ow728lmc` and can be obtained through the `DescribeDBInstances` API.",
			},

			"ipv6_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether IPv6 is used. Default value: 0.",
			},
		},
	}
}

func resourceTencentCloudDcdbOpenDBExtranetAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_open_d_b_extranet_access.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewOpenDBExtranetAccessRequest()
		response   = dcdb.NewOpenDBExtranetAccessResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("ipv6_flag"); v != nil {
		request.Ipv6Flag = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().OpenDBExtranetAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb openDBExtranetAccess failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbOpenDBExtranetAccessRead(d, meta)
}

func resourceTencentCloudDcdbOpenDBExtranetAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_open_d_b_extranet_access.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbOpenDBExtranetAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_open_d_b_extranet_access.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
