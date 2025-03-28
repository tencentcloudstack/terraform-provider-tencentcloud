package tmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpTkeTemplateAttachment() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead,
		Create: resourceTencentCloudMonitorTmpTkeTemplateAttachmentCreate,
		Delete: resourceTencentCloudMonitorTmpTkeTemplateAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the template, which is used for the outgoing reference.",
			},

			"targets": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "Sync target details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "target area.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "instance id.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the cluster.",
						},
						"sync_time": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Last sync template time.",
						},
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Template version currently in use.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Cluster type.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the prometheus instance.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name the cluster.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewSyncPrometheusTempRequest()

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "targets"); ok {
		var prometheusTarget monitor.PrometheusTemplateSyncTarget
		if v, ok := dMap["region"]; ok {
			prometheusTarget.Region = helper.String(v.(string))
		}

		if v, ok := dMap["instance_id"]; ok {
			prometheusTarget.InstanceId = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_id"]; ok {
			prometheusTarget.ClusterId = helper.String(v.(string))
		}

		if v, ok := dMap["sync_time"]; ok {
			prometheusTarget.SyncTime = helper.String(v.(string))
		}

		if v, ok := dMap["version"]; ok {
			prometheusTarget.Version = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_type"]; ok {
			prometheusTarget.ClusterType = helper.String(v.(string))
		}

		if v, ok := dMap["instance_name"]; ok {
			prometheusTarget.InstanceName = helper.String(v.(string))
		}

		if v, ok := dMap["cluster_name"]; ok {
			prometheusTarget.ClusterName = helper.String(v.(string))
		}

		prometheusTargets := make([]*monitor.PrometheusTemplateSyncTarget, 0)
		prometheusTargets = append(prometheusTargets, &prometheusTarget)
		request.Targets = prometheusTargets

	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().SyncPrometheusTemp(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s sync tke template failed, reason:%+v", logId, err)
		return err
	}

	templateId := *request.TemplateId
	instanceId := *request.Targets[0].InstanceId
	region := *request.Targets[0].Region
	d.SetId(strings.Join([]string{templateId, instanceId, region}, tccommon.FILED_SP))

	return resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	templateId := ids[0]
	instanceId := ids[1]
	region := ids[2]

	targets, err := service.DescribePrometheusTempSync(ctx, templateId)

	if err != nil {
		return err
	}

	if targets == nil || len(targets) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `targets` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	tempTargets := make([]map[string]interface{}, 0)
	for _, v := range targets {
		if *v.InstanceId == instanceId && *v.Region == region {
			tempTargets = append(tempTargets, map[string]interface{}{
				"region":      v.Region,
				"instance_id": v.InstanceId,
				//"cluster_id":    v.ClusterId,
				//"sync_time":     v.SyncTime,
				//"version":       v.Version,
				//"cluster_type":  v.ClusterType,
				//"instance_name": v.InstanceName,
				//"cluster_name":  v.ClusterName,
			})
		}
	}
	_ = d.Set("targets", tempTargets)

	return nil
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	request := monitor.NewDeletePrometheusTempSyncRequest()

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	templateId := ids[0]
	instanceId := ids[1]
	region := ids[2]

	request.TemplateId = &templateId
	var targets []*monitor.PrometheusTemplateSyncTarget
	target := monitor.PrometheusTemplateSyncTarget{
		Region:     &region,
		InstanceId: &instanceId,
	}
	targets = append(targets, &target)
	request.Targets = targets

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().DeletePrometheusTempSync(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
