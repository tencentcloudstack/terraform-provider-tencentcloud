/*
Provides a resource to create a tmp tke template attachment

Example Usage

```hcl

resource "tencentcloud_monitor_tmp_tke_template_attachment" "temp_attachment" {
  template_id  = "temp-xxx"

  targets {
    region      = "ap-xxx"
    instance_id = "prom-xxx"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeTemplateAttachment() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tke.NewSyncPrometheusTempRequest()

	if v, ok := d.GetOk("template_id"); ok {
		request.TemplateId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "targets"); ok {
		var prometheusTarget tke.PrometheusTemplateSyncTarget
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

		prometheusTargets := make([]*tke.PrometheusTemplateSyncTarget, 0)
		prometheusTargets = append(prometheusTargets, &prometheusTarget)
		request.Targets = prometheusTargets

	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().SyncPrometheusTemp(request)
		if e != nil {
			return retryError(e)
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
	d.SetId(strings.Join([]string{templateId, instanceId, region}, FILED_SP))

	return resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeTemplateAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TkeService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
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
		return fmt.Errorf("resource `targets` %s does not exist", templateId)
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
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_template_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	request := tke.NewDeletePrometheusTempSyncRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 3 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	templateId := ids[0]
	instanceId := ids[1]
	region := ids[2]

	request.TemplateId = &templateId
	var targets []*tke.PrometheusTemplateSyncTarget
	target := tke.PrometheusTemplateSyncTarget{
		Region:     &region,
		InstanceId: &instanceId,
	}
	targets = append(targets, &target)
	request.Targets = targets

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTkeClient().DeletePrometheusTempSync(request)
		if e != nil {
			return retryError(e)
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
