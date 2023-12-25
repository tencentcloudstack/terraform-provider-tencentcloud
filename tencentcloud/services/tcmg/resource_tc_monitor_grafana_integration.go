package tcmg

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorGrafanaIntegration() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorGrafanaIntegrationRead,
		Create: resourceTencentCloudMonitorGrafanaIntegrationCreate,
		Update: resourceTencentCloudMonitorGrafanaIntegrationUpdate,
		Delete: resourceTencentCloudMonitorGrafanaIntegrationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "grafana instance id.",
			},

			"integration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "integration id.",
			},

			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "integration json schema kind.",
			},

			"content": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "generated json string of given integration json schema.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "integration desc.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaIntegrationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_integration.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request       = monitor.NewCreateGrafanaIntegrationRequest()
		response      *monitor.CreateGrafanaIntegrationResponse
		integrationId string
		instanceId    string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("kind"); ok {
		request.Kind = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateGrafanaIntegration(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaIntegration failed, reason:%+v", logId, err)
		return err
	}

	integrationId = *response.Response.IntegrationId

	d.SetId(strings.Join([]string{integrationId, instanceId}, tccommon.FILED_SP))
	return resourceTencentCloudMonitorGrafanaIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_integration.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	instanceId := idSplit[1]

	grafanaIntegration, err := service.DescribeMonitorGrafanaIntegration(ctx, integrationId, instanceId)

	if err != nil {
		return err
	}

	if grafanaIntegration == nil {
		d.SetId("")
		return fmt.Errorf("resource `grafanaIntegration` %s does not exist", integrationId)
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaIntegration.IntegrationId != nil {
		_ = d.Set("integration_id", grafanaIntegration.IntegrationId)
	}

	if grafanaIntegration.Kind != nil {
		_ = d.Set("kind", grafanaIntegration.Kind)
	}

	//if grafanaIntegration.Content != nil {
	//	_ = d.Set("content", grafanaIntegration.Content)
	//}

	return nil
}

func resourceTencentCloudMonitorGrafanaIntegrationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_integration.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateGrafanaIntegrationRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	instanceId := idSplit[1]

	request.IntegrationId = &integrationId
	request.InstanceId = &instanceId

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("kind") {
		return fmt.Errorf("`kind` do not support change now.")
	} else {
		if v, ok := d.GetOk("kind"); ok {
			request.Kind = helper.String(v.(string))
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateGrafanaIntegration(request)
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

	return resourceTencentCloudMonitorGrafanaIntegrationRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaIntegrationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_integration.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	integrationId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteMonitorGrafanaIntegrationById(ctx, integrationId, instanceId); err != nil {
		return err
	}

	return nil
}
