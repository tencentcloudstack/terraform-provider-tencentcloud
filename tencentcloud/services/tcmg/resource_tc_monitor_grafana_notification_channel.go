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

func ResourceTencentCloudMonitorGrafanaNotificationChannel() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorGrafanaNotificationChannelRead,
		Create: resourceTencentCloudMonitorGrafanaNotificationChannelCreate,
		Update: resourceTencentCloudMonitorGrafanaNotificationChannelUpdate,
		Delete: resourceTencentCloudMonitorGrafanaNotificationChannelDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "grafana instance id.",
			},

			"channel_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "plugin id.",
			},

			"channel_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "channel name.",
			},

			"org_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Grafana organization which channel will be installed, default to 1 representing Main Org.",
			},

			"receivers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Computed:    true,
				Description: "cloud monitor notification template notice-id list.",
			},

			"extra_org_ids": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "extra grafana organization id list, default to 1 representing Main Org.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaNotificationChannelCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_notification_channel.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = monitor.NewCreateGrafanaNotificationChannelRequest()
		response   *monitor.CreateGrafanaNotificationChannelResponse
		channelId  string
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("channel_name"); ok {
		request.ChannelName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("org_id"); ok {
		request.OrgId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("receivers"); ok {
		receiversSet := v.(*schema.Set).List()
		for i := range receiversSet {
			receivers := receiversSet[i].(string)
			request.Receivers = append(request.Receivers, &receivers)
		}
	}

	if v, ok := d.GetOk("extra_org_ids"); ok {
		extraOrgIdsSet := v.(*schema.Set).List()
		for i := range extraOrgIdsSet {
			extraOrgIds := extraOrgIdsSet[i].(string)
			request.ExtraOrgIds = append(request.ExtraOrgIds, &extraOrgIds)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreateGrafanaNotificationChannel(request)
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
		log.Printf("[CRITAL]%s create monitor grafanaNotificationChannel failed, reason:%+v", logId, err)
		return err
	}

	channelId = *response.Response.ChannelId

	d.SetId(channelId + tccommon.FILED_SP + instanceId)
	return resourceTencentCloudMonitorGrafanaNotificationChannelRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_notification_channel.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	instanceId := idSplit[1]

	grafanaNotificationChannel, err := service.DescribeMonitorGrafanaNotificationChannel(ctx, channelId, instanceId)

	if err != nil {
		return err
	}

	if grafanaNotificationChannel == nil {
		d.SetId("")
		return fmt.Errorf("resource `grafanaNotificationChannel` %s does not exist", channelId)
	}

	_ = d.Set("instance_id", instanceId)

	if grafanaNotificationChannel.ChannelId != nil {
		_ = d.Set("channel_id", grafanaNotificationChannel.ChannelId)
	}

	if grafanaNotificationChannel.ChannelName != nil {
		_ = d.Set("channel_name", grafanaNotificationChannel.ChannelName)
	}

	_ = d.Set("receivers", grafanaNotificationChannel.Receivers)

	return nil
}

func resourceTencentCloudMonitorGrafanaNotificationChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_notification_channel.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewUpdateGrafanaNotificationChannelRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	instanceId := idSplit[1]

	request.ChannelId = &channelId
	request.InstanceId = &instanceId

	if v, ok := d.GetOk("channel_name"); ok {
		request.ChannelName = helper.String(v.(string))
	}

	if d.HasChange("instance_id") {
		return fmt.Errorf("`instance_id` do not support change now.")
	}

	if d.HasChange("channel_name") {
		return fmt.Errorf("`channel_name` do not support change now.")
	}

	if d.HasChange("org_id") {
		return fmt.Errorf("`org_id` do not support change now.")
	}

	if d.HasChange("receivers") {
		return fmt.Errorf("`receivers` do not support change now.")
	}

	if d.HasChange("extra_org_ids") {
		return fmt.Errorf("`extra_org_ids` do not support change now.")
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().UpdateGrafanaNotificationChannel(request)
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

	return resourceTencentCloudMonitorGrafanaNotificationChannelRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaNotificationChannelDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_grafana_notification_channel.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteMonitorGrafanaNotificationChannelById(ctx, channelId, instanceId); err != nil {
		return err
	}

	return nil
}
