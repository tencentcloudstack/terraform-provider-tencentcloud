/*
Provides a resource to create a dcdb cancel_dcn_job

Example Usage

```hcl
resource "tencentcloud_dcdb_cancel_dcn_job" "cancel_dcn_job" {
  instance_id = ""
}
```

Import

dcdb cancel_dcn_job can be imported using the id, e.g.

```
terraform import tencentcloud_dcdb_cancel_dcn_job.cancel_dcn_job cancel_dcn_job_id
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

func resourceTencentCloudDcdbCancelDcnJob() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbCancelDcnJobCreate,
		Read:   resourceTencentCloudDcdbCancelDcnJobRead,
		Delete: resourceTencentCloudDcdbCancelDcnJobDelete,
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

func resourceTencentCloudDcdbCancelDcnJobCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = dcdb.NewCancelDcnJobRequest()
		response   = dcdb.NewCancelDcnJobResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().CancelDcnJob(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dcdb cancelDcnJob failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudDcdbCancelDcnJobRead(d, meta)
}

func resourceTencentCloudDcdbCancelDcnJobRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDcdbCancelDcnJobDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_cancel_dcn_job.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
