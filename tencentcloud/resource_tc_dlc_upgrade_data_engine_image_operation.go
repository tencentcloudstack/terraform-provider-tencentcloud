package tencentcloud

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcUpgradeDataEngineImageOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpgradeDataEngineImageOperationCreate,
		Read:   resourceTencentCloudDlcUpgradeDataEngineImageOperationRead,
		Delete: resourceTencentCloudDlcUpgradeDataEngineImageOperationDelete,
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
		},
	}
}

func resourceTencentCloudDlcUpgradeDataEngineImageOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_upgrade_data_engine_image_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = dlc.NewUpgradeDataEngineImageRequest()
		dataEngineId string
	)
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().UpgradeDataEngineImage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc upgradeDataEngineImageOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineId)

	service := DlcService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"2"}, 5*readRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDlcUpgradeDataEngineImageOperationRead(d, meta)
}

func resourceTencentCloudDlcUpgradeDataEngineImageOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_upgrade_data_engine_image_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpgradeDataEngineImageOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_upgrade_data_engine_image_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
