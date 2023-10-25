/*
Provides a resource to create a dlc modify_data_engine_description_operation

Example Usage

```hcl
resource "tencentcloud_dlc_modify_data_engine_description_operation" "modify_data_engine_description_operation" {
  data_engine_name = "testEngine"
  message = "test"
}
```

Import

dlc modify_data_engine_description_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_modify_data_engine_description_operation.modify_data_engine_description_operation modify_data_engine_description_operation_id
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

func resourceTencentCloudDlcModifyDataEngineDescriptionOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcModifyDataEngineDescriptionOperationCreate,
		Read:   resourceTencentCloudDlcModifyDataEngineDescriptionOperationRead,
		Delete: resourceTencentCloudDlcModifyDataEngineDescriptionOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The name of the engine to modify.",
			},

			"message": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine description information, the maximum length is 250.",
			},
		},
	}
}

func resourceTencentCloudDlcModifyDataEngineDescriptionOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_modify_data_engine_description_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request        = dlc.NewModifyDataEngineDescriptionRequest()
		dataEngineName string
	)
	if v, ok := d.GetOk("data_engine_name"); ok {
		dataEngineName = v.(string)
		request.DataEngineName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("message"); ok {
		request.Message = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().ModifyDataEngineDescription(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc modifyDataEngineDescriptionOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineName)

	return resourceTencentCloudDlcModifyDataEngineDescriptionOperationRead(d, meta)
}

func resourceTencentCloudDlcModifyDataEngineDescriptionOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_modify_data_engine_description_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcModifyDataEngineDescriptionOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_modify_data_engine_description_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
