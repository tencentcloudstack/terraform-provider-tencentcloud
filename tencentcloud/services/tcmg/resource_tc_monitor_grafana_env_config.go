package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorGrafanaEnvConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaEnvConfigCreate,
		Read:   resourceTencentCloudMonitorGrafanaEnvConfigRead,
		Update: resourceTencentCloudMonitorGrafanaEnvConfigUpdate,
		Delete: resourceTencentCloudMonitorGrafanaEnvConfigDelete,
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

			"envs": {
				Optional:    true,
				Type:        schema.TypeMap,
				Description: "Environment variables.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaEnvConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_env_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaEnvConfigUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaEnvConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_env_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	grafanaEnvConfig, err := service.DescribeMonitorGrafanaEnvConfigById(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaEnvConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaEnvConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaEnvConfig.Envs != nil {
		v, err := helper.JsonToMap(*grafanaEnvConfig.Envs)
		if err != nil {
			return fmt.Errorf("envs `%v` format error", *grafanaEnvConfig.Envs)
		}
		_ = d.Set("envs", v)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaEnvConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_env_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := monitor.NewUpdateGrafanaEnvironmentsRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("envs"); ok {
		evs, o := helper.MapToString(v.(map[string]interface{}))
		if !o {
			return fmt.Errorf("envs `%s` format error", v)
		}
		request.Envs = &evs
	} else {
		request.Envs = helper.String("{}")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateGrafanaEnvironments(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaEnvConfig failed, reason:%+v", logId, err)
		return err
	}

	time.Sleep(3 * time.Second)
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err = resource.Retry(1*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		instance, errRet := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		if *instance.InstanceStatus == 2 {
			return nil
		}
		if *instance.InstanceStatus == 3 {
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, update envs config failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaEnvConfigRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaEnvConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_env_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
