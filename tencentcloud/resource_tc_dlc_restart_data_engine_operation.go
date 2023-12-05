package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcRestartDataEngineOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcRestartDataEngineCreateOperation,
		Read:   resourceTencentCloudDlcRestartDataEngineReadOperation,
		Delete: resourceTencentCloudDlcRestartDataEngineDeleteOperation,
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

			"forced_operation": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to force restart and ignore tasks.",
			},
		},
	}
}

func resourceTencentCloudDlcRestartDataEngineCreateOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_restart_data_engine_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = dlc.NewRestartDataEngineRequest()
		dataEngineId string
	)
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineId = helper.String(v.(string))
	}

	if v, _ := d.GetOkExists("forced_operation"); v != nil {
		request.ForcedOperation = helper.Bool(v.(bool))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().RestartDataEngine(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc restartDataEngine failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 5*readRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDlcRestartDataEngineReadOperation(d, meta)
}

func resourceTencentCloudDlcRestartDataEngineReadOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_restart_data_engine_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcRestartDataEngineDeleteOperation(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_restart_data_engine_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
