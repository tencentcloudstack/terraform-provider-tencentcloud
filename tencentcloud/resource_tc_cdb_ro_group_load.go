/*
Provides a resource to create a cdb ro_group_load

Example Usage

```hcl
resource "tencentcloud_cdb_ro_group_load" "ro_group_load" {
  ro_group_id = ""
}
```

Import

cdb ro_group_load can be imported using the id, e.g.

```
terraform import tencentcloud_cdb_ro_group_load.ro_group_load ro_group_load_id
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

func resourceTencentCloudCdbRoGroupLoad() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCdbRoGroupLoadCreate,
		Read:   resourceTencentCloudCdbRoGroupLoadRead,
		Delete: resourceTencentCloudCdbRoGroupLoadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"ro_group_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The ID of the RO group, in the format: cdbrg-c1nl9rpv.",
			},
		},
	}
}

func resourceTencentCloudCdbRoGroupLoadCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group_load.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = cdb.NewBalanceRoGroupLoadRequest()
		response  = cdb.NewBalanceRoGroupLoadResponse()
		roGroupId string
	)
	if v, ok := d.GetOk("ro_group_id"); ok {
		roGroupId = v.(string)
		request.RoGroupId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCdbClient().BalanceRoGroupLoad(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate cdb roGroupLoad failed, reason:%+v", logId, err)
		return err
	}

	roGroupId = *response.Response.RoGroupId
	d.SetId(roGroupId)

	return resourceTencentCloudCdbRoGroupLoadRead(d, meta)
}

func resourceTencentCloudCdbRoGroupLoadRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group_load.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCdbRoGroupLoadDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cdb_ro_group_load.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
