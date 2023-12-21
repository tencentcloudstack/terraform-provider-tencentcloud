package dlc

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcStoreLocationConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcStoreLocationConfigCreate,
		Read:   resourceTencentCloudDlcStoreLocationConfigRead,
		Update: resourceTencentCloudDlcStoreLocationConfigUpdate,
		Delete: resourceTencentCloudDlcStoreLocationConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"store_location": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The calculation results are stored in the cos path, such as: cosn://bucketname/.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Whether to enable advanced settings: 0-no, 1-yes.",
			},
		},
	}
}

func resourceTencentCloudDlcStoreLocationConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_store_location_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	return resourceTencentCloudDlcStoreLocationConfigUpdate(d, meta)
}

func resourceTencentCloudDlcStoreLocationConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_store_location_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	storeLocationConfig, err := service.DescribeDlcStoreLocationConfigById(ctx)
	if err != nil {
		return err
	}

	if storeLocationConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DlcStoreLocationConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if storeLocationConfig.StoreLocation != nil {
		_ = d.Set("store_location", storeLocationConfig.StoreLocation)
	}

	if storeLocationConfig.Enable != nil {
		_ = d.Set("enable", storeLocationConfig.Enable)
	}
	return nil
}

func resourceTencentCloudDlcStoreLocationConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_store_location_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := dlc.NewModifyAdvancedStoreLocationRequest()

	var storeLocation string
	if v, ok := d.GetOk("store_location"); ok {
		storeLocation = v.(string)
		request.StoreLocation = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.IntUint64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().ModifyAdvancedStoreLocation(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dlc storeLocationConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(storeLocation)
	return resourceTencentCloudDlcStoreLocationConfigRead(d, meta)
}

func resourceTencentCloudDlcStoreLocationConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_store_location_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
