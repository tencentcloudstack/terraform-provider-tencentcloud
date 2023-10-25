/*
Provides a resource to create a dlc rollback_data_engine_image

Example Usage

```hcl
data "tencentcloud_dlc_check_data_engine_image_can_be_rollback" "check_data_engine_image_can_be_rollback" {
  data_engine_id = "DataEngine-cgkvbas6"
}
resource "tencentcloud_dlc_rollback_data_engine_image_operation" "rollback_data_engine_image" {
  data_engine_id = "DataEngine-cgkvbas6"
  from_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.from_record_id
  to_record_id = data.tencentcloud_dlc_check_data_engine_image_can_be_rollback.check_data_engine_image_can_be_rollback.to_record_id
}
```

Import

dlc rollback_data_engine_image can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_rollback_data_engine_image_operation.rollback_data_engine_image rollback_data_engine_image_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcRollbackDataEngineImageOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcRollbackDataEngineImageCreateOperation,
		Read:   resourceTencentCloudDlcRollbackDataEngineImageReadOperation,
		Delete: resourceTencentCloudDlcRollbackDataEngineImageDeleteOperation,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine unique id.",
			},

			"from_record_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log record id before rollback.",
			},

			"to_record_id": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Log record id after rollback.",
			},
		},
	}
}

func resourceTencentCloudDlcRollbackDataEngineImageCreateOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_rollback_data_engine_image_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = dlc.NewRollbackDataEngineImageRequest()
		dataEngineId string
	)
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("from_record_id"); ok {
		request.FromRecordId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("to_record_id"); ok {
		request.ToRecordId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().RollbackDataEngineImage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc rollbackDataEngineImage failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineId)

	return resourceTencentCloudDlcRollbackDataEngineImageReadOperation(d, meta)
}

func resourceTencentCloudDlcRollbackDataEngineImageReadOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_rollback_data_engine_image_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcRollbackDataEngineImageDeleteOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_rollback_data_engine_image_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
