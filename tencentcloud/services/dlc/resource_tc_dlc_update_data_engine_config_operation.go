package dlc

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcUpdateDataEngineConfigOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpdateDataEngineConfigOperationCreate,
		Read:   resourceTencentCloudDlcUpdateDataEngineConfigOperationRead,
		Delete: resourceTencentCloudDlcUpdateDataEngineConfigOperationDelete,
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
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = dlc.NewUpdateDataEngineConfigRequest()
		dataEngineId string
	)
	if v, ok := d.GetOk("data_engine_id"); ok {
		dataEngineId = v.(string)
		request.DataEngineIds = []*string{&dataEngineId}
	}

	if v, ok := d.GetOk("data_engine_config_command"); ok {
		request.DataEngineConfigCommand = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateDataEngineConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc updateDataEngineConfigOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dataEngineId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"2"}, 5*tccommon.ReadRetryTimeout, time.Second, service.DlcRestartDataEngineStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}
	return resourceTencentCloudDlcUpdateDataEngineConfigOperationRead(d, meta)
}

func resourceTencentCloudDlcUpdateDataEngineConfigOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpdateDataEngineConfigOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_data_engine_config_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
