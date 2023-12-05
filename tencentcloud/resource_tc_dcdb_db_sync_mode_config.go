package tencentcloud

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDcdbDbSyncModeConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDcdbDbSyncModeConfigCreate,
		Read:   resourceTencentCloudDcdbDbSyncModeConfigRead,
		Update: resourceTencentCloudDcdbDbSyncModeConfigUpdate,
		Delete: resourceTencentCloudDcdbDbSyncModeConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "ID of the instance for which to modify the sync mode. The ID is in the format of `tdsql-ow728lmc`.",
			},

			"sync_mode": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Sync mode. Valid values: `0` (async), `1` (strong sync), `2` (downgradable strong sync).",
			},
		},
	}
}

func resourceTencentCloudDcdbDbSyncModeConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_sync_mode_config.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}
	d.SetId(instanceId)

	return resourceTencentCloudDcdbDbSyncModeConfigUpdate(d, meta)
}

func resourceTencentCloudDcdbDbSyncModeConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_sync_mode_config.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	dbSyncModeConfig, err := service.DescribeDcdbDbSyncModeConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if dbSyncModeConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `DcdbDbSyncModeConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if dbSyncModeConfig.SyncMode != nil {
		_ = d.Set("sync_mode", dbSyncModeConfig.SyncMode)
	}

	return nil
}

func resourceTencentCloudDcdbDbSyncModeConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_sync_mode_config.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := dcdb.NewModifyDBSyncModeRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if d.HasChange("sync_mode") {
		if v, ok := d.GetOkExists("sync_mode"); ok {
			request.SyncMode = helper.IntInt64(v.(int))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDcdbClient().ModifyDBSyncMode(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update dcdb dbSyncModeConfig failed, reason:%+v", logId, err)
		return err
	}

	service := DcdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"0"}, 2*readRetryTimeout, time.Second, service.DcdbDbSyncModeConfigStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudDcdbDbSyncModeConfigRead(d, meta)
}

func resourceTencentCloudDcdbDbSyncModeConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dcdb_db_sync_mode_config.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
