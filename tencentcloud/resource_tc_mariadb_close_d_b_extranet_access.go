/*
Provides a resource to create a mariadb close_d_b_extranet_access

Example Usage

```hcl
resource "tencentcloud_mariadb_close_d_b_extranet_access" "close_d_b_extranet_access" {
  instance_id = ""
  ipv6_flag =
}
```

Import

mariadb close_d_b_extranet_access can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_close_d_b_extranet_access.close_d_b_extranet_access close_d_b_extranet_access_id
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

func resourceTencentCloudMariadbCloseDBExtranetAccess() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCloseDBExtranetAccessCreate,
		Read:   resourceTencentCloudMariadbCloseDBExtranetAccessRead,
		Delete: resourceTencentCloudMariadbCloseDBExtranetAccessDelete,
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

func resourceTencentCloudMariadbCloseDBExtranetAccessCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_d_b_extranet_access.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewCloseDBExtranetAccessRequest()
		response   = mariadb.NewCloseDBExtranetAccessResponse()
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CloseDBExtranetAccess(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb closeDBExtranetAccess failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbCloseDBExtranetAccessRead(d, meta)
}

func resourceTencentCloudMariadbCloseDBExtranetAccessRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_d_b_extranet_access.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCloseDBExtranetAccessDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_close_d_b_extranet_access.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
