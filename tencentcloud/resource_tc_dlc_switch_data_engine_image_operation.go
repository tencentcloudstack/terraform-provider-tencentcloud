/*
Provides a resource to create a dlc switch_data_engine_image_operation

Example Usage

```hcl
resource "tencentcloud_dlc_switch_data_engine_image_operation" "switch_data_engine_image_operation" {
  data_engine_id = "DataEngine-g5ds87d8"
  new_image_version_id = "344ba1c6-b7a9-403a-a255-422fffed6d38"
}
```

Import

dlc switch_data_engine_image_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation switch_data_engine_image_operation_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudDlcSwitchDataEngineImageOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcSwitchDataEngineImageOperationCreate,
		Read:   resourceTencentCloudDlcSwitchDataEngineImageOperationRead,
		Delete: resourceTencentCloudDlcSwitchDataEngineImageOperationDelete,
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

			"new_image_version_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "New image version id.",
			},
		},
	}
}

func resourceTencentCloudDlcSwitchDataEngineImageOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_switch_data_engine_image_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = dlc.NewSwitchDataEngineImageRequest()
		dataEngineId string
	)
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("new_image_version_id"); ok {
		request.NewImageVersionId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().SwitchDataEngineImage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc switchDataEngineImageOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 5*readRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDlcSwitchDataEngineImageOperationRead(d, meta)
}

func resourceTencentCloudDlcSwitchDataEngineImageOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_switch_data_engine_image_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcSwitchDataEngineImageOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_switch_data_engine_image_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
