package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpGrafanaConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpGrafanaConfigCreate,
		Read:   resourceTencentCloudMonitorTmpGrafanaConfigRead,
		Update: resourceTencentCloudMonitorTmpGrafanaConfigUpdate,
		Delete: resourceTencentCloudMonitorTmpGrafanaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"config": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "JSON encoded string.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpGrafanaConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_grafana_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorTmpGrafanaConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpGrafanaConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_grafana_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	tmpGrafanaConfig, err := service.DescribeMonitorTmpGrafanaConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if tmpGrafanaConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpGrafanaConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if tmpGrafanaConfig.Config != nil {
		_ = d.Set("config", tmpGrafanaConfig.Config)
	}

	return nil
}

func resourceTencentCloudMonitorTmpGrafanaConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_grafana_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateGrafanaConfigRequest()

	instanceId := d.Id()
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("config"); ok {
		request.Config = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateGrafanaConfig(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor tmpGrafanaConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorTmpGrafanaConfigRead(d, meta)
}

func resourceTencentCloudMonitorTmpGrafanaConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_grafana_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
