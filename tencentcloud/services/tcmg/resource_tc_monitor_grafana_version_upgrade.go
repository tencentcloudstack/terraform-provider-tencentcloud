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

func ResourceTencentCloudMonitorGrafanaVersionUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaVersionUpgradeCreate,
		Read:   resourceTencentCloudMonitorGrafanaVersionUpgradeRead,
		Update: resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate,
		Delete: resourceTencentCloudMonitorGrafanaVersionUpgradeDelete,
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

			"alias": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Version alias.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate(d, meta)
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	grafanaVersionUpgrade, err := service.DescribeMonitorGrafanaInstance(ctx, instanceId)
	if err != nil {
		return err
	}

	if grafanaVersionUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaVersionUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaVersionUpgrade.Version != nil {
		_ = d.Set("alias", grafanaVersionUpgrade.Version)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := monitor.NewUpgradeGrafanaInstanceRequest()

	instanceId := d.Id()

	request.InstanceId = &instanceId

	if v, ok := d.GetOk("alias"); ok {
		request.Alias = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpgradeGrafanaInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaVersionUpgrade failed, reason:%+v", logId, err)
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
			return resource.NonRetryableError(fmt.Errorf("grafanaInstance status is %v, update version failed.", *instance.InstanceStatus))
		}
		return resource.RetryableError(fmt.Errorf("grafanaInstance status is %v, retry...", *instance.InstanceStatus))
	})
	if err != nil {
		return err
	}

	return resourceTencentCloudMonitorGrafanaVersionUpgradeRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaVersionUpgradeDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_version_upgrade.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
