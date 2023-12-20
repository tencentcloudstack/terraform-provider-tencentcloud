package cynosdb

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCynosdbBinlogSaveDays() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCynosdbBinlogSaveDaysCreate,
		Read:   resourceTencentCloudCynosdbBinlogSaveDaysRead,
		Update: resourceTencentCloudCynosdbBinlogSaveDaysUpdate,
		Delete: resourceTencentCloudCynosdbBinlogSaveDaysDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"binlog_save_days": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Binlog retention days.",
			},
		},
	}
}

func resourceTencentCloudCynosdbBinlogSaveDaysCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_binlog_save_days.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var clusterId string
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	d.SetId(clusterId)

	return resourceTencentCloudCynosdbBinlogSaveDaysUpdate(d, meta)
}

func resourceTencentCloudCynosdbBinlogSaveDaysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_binlog_save_days.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	clusterId := d.Id()
	binlogSaveDays, err := service.DescribeCynosdbBinlogSaveDaysById(ctx, clusterId)
	if err != nil {
		return err
	}

	if binlogSaveDays == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CynosdbBinlogSaveDays` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("cluster_id", clusterId)

	if binlogSaveDays != nil {
		_ = d.Set("binlog_save_days", binlogSaveDays)
	}

	return nil
}

func resourceTencentCloudCynosdbBinlogSaveDaysUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_binlog_save_days.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cynosdb.NewModifyBinlogSaveDaysRequest()

	clusterId := d.Id()
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("binlog_save_days"); ok {
		request.BinlogSaveDays = helper.IntInt64(v.(int))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCynosdbClient().ModifyBinlogSaveDays(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update cynosdb binlogSaveDays failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCynosdbBinlogSaveDaysRead(d, meta)
}

func resourceTencentCloudCynosdbBinlogSaveDaysDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cynosdb_binlog_save_days.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
