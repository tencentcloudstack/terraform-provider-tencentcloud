package tmp

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

func ResourceTencentCloudMonitorTmpManageGrafanaAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorTmpManageGrafanaAttachmentCreate,
		Read:   resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead,
		Delete: resourceTencentCloudMonitorTmpManageGrafanaAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Prometheus instance ID.",
			},

			"grafana_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance ID.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = monitor.NewBindPrometheusManagedGrafanaRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("grafana_id"); ok {
		request.GrafanaId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().BindPrometheusManagedGrafana(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor manageGrafanaAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead(d, meta)
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	manageGrafanaAttachment, err := service.DescribeMonitorManageGrafanaAttachmentById(ctx, instanceId)
	if err != nil {
		return err
	}

	if manageGrafanaAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorTmpManageGrafanaAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if manageGrafanaAttachment.InstanceId != nil {
		_ = d.Set("instance_id", manageGrafanaAttachment.InstanceId)
	}

	if manageGrafanaAttachment.GrafanaInstanceId != nil {
		_ = d.Set("grafana_id", manageGrafanaAttachment.GrafanaInstanceId)
	}

	return nil
}

func resourceTencentCloudMonitorTmpManageGrafanaAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_manage_grafana_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	instanceId := d.Id()

	if err := service.DeleteMonitorManageGrafanaAttachmentById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
