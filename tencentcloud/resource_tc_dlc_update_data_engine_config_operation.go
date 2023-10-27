/*
Provides a resource to create a dlc update_data_engine_config_operation

Example Usage

```hcl
resource "tencentcloud_dlc_update_data_engine_config_operation" "update_data_engine_config_operation" {
  data_engine_ids =
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```

Import

dlc update_data_engine_config_operation can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation update_data_engine_config_operation_id
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

func resourceTencentCloudDlcUpdateDataEngineConfigOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpdateDataEngineConfigOperationCreate,
		Read:   resourceTencentCloudDlcUpdateDataEngineConfigOperationRead,
		Delete: resourceTencentCloudDlcUpdateDataEngineConfigOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"data_engine_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Engine unique id.",
			},

			"data_engine_config_command": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Engine configuration command, supports UpdateSparkSQLLakefsPath (update native table configuration), UpdateSparkSQLResultPath (update result path configuration).",
			},
		},
	}
}

func resourceTencentCloudDlcUpdateDataEngineConfigOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = dlc.NewUpdateDataEngineConfigRequest()
		response     = dlc.NewUpdateDataEngineConfigResponse()
		dataEngineId uint64
	)
	if v, ok := d.GetOk("data_engine_ids"); ok {
		dataEngineIdsSet := v.(*schema.Set).List()
		for i := range dataEngineIdsSet {
			dataEngineIds := dataEngineIdsSet[i].(string)
			request.DataEngineIds = append(request.DataEngineIds, &dataEngineIds)
		}
	}

	if v, ok := d.GetOk("data_engine_config_command"); ok {
		request.DataEngineConfigCommand = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().UpdateDataEngineConfig(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc updateDataEngineConfigOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.UInt64ToStr(dataEngineId))

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 5*readRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDlcUpdateDataEngineConfigOperationRead(d, meta)
}

func resourceTencentCloudDlcUpdateDataEngineConfigOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpdateDataEngineConfigOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
