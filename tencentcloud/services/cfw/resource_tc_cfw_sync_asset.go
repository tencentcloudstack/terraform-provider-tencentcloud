package cfw

import (
	"fmt"
	"log"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfw "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"
)

func ResourceTencentCloudCfwSyncAsset() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwSyncAssetCreate,
		Read:   resourceTencentCloudCfwSyncAssetRead,
		Delete: resourceTencentCloudCfwSyncAssetDelete,
		Schema: map[string]*schema.Schema{},
	}
}

func resourceTencentCloudCfwSyncAssetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_sync_asset.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = cfw.NewModifyAssetSyncRequest()
		statusRequest = cfw.NewDescribeAssetSyncRequest()
	)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().ModifyAssetSync(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cfw syncAsset failed, reason:%+v", logId, err)
		return err
	}

	// wait
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwClient().DescribeAssetSync(statusRequest)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if *result.Response.Status == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("The fw sync asset status is %d.", *result.Response.Status))
	})

	if err != nil {
		log.Printf("[CRITAL]%s operate cfw syncAsset status failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return resourceTencentCloudCfwSyncAssetRead(d, meta)
}

func resourceTencentCloudCfwSyncAssetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_sync_asset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudCfwSyncAssetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_sync_asset.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
