/*
Provides a resource to create a mariadb cancel_dcn_job

Example Usage

```hcl
resource "tencentcloud_mariadb_cancel_dcn_job" "cancel_dcn_job" {
  instance_id = ""
}
```

Import

mariadb cancel_dcn_job can be imported using the id, e.g.

```
terraform import tencentcloud_mariadb_cancel_dcn_job.cancel_dcn_job cancel_dcn_job_id
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

func resourceTencentCloudMariadbCancelDcnJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMariadbCancelDcnJobCreate,
		Read:   resourceTencentCloudMariadbCancelDcnJobRead,
		Delete: resourceTencentCloudMariadbCancelDcnJobDelete,
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
		},
	}
}

func resourceTencentCloudMariadbCancelDcnJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_cancel_dcn_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mariadb.NewCancelDcnJobRequest()
		response   = mariadb.NewCancelDcnJobResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMariadbClient().CancelDcnJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mariadb cancelDcnJob failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudMariadbCancelDcnJobRead(d, meta)
}

func resourceTencentCloudMariadbCancelDcnJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_cancel_dcn_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMariadbCancelDcnJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mariadb_cancel_dcn_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
