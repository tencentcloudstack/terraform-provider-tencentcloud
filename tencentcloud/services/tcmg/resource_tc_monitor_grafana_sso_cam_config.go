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

func ResourceTencentCloudMonitorGrafanaSsoCamConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaSsoCamConfigCreate,
		Read:   resourceTencentCloudMonitorGrafanaSsoCamConfigRead,
		Update: resourceTencentCloudMonitorGrafanaSsoCamConfigUpdate,
		Delete: resourceTencentCloudMonitorGrafanaSsoCamConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},

			"enable_sso_cam_check": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable the CAM authorization: `true` for enabling; `false` for disabling.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaSsoCamConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_cam_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaSsoCamConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoCamConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_cam_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	grafanaSsoCamConfig, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaSsoCamConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaSsoCamConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaSsoCamConfig.EnableSSOCamCheck != nil {
		_ = d.Set("enable_sso_cam_check", grafanaSsoCamConfig.EnableSSOCamCheck)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaSsoCamConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_cam_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewEnableSSOCamCheckRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOkExists("enable_sso_cam_check"); ok {
		request.EnableSSOCamCheck = helper.Bool(v.(bool))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().EnableSSOCamCheck(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaSsoCamConfig failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaSsoCamConfigRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaSsoCamConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_sso_cam_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
