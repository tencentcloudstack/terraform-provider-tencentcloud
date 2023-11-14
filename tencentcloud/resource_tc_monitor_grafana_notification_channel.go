/*
Provides a resource to create a monitor grafana_notification_channel

Example Usage

```hcl
resource "tencentcloud_monitor_grafana_notification_channel" "grafana_notification_channel" {
  instance_id = &lt;nil&gt;
    channel_name = &lt;nil&gt;
  org_id = 1
  receivers = &lt;nil&gt;
  extra_org_ids = &lt;nil&gt;
}
```

Import

monitor grafana_notification_channel can be imported using the id, e.g.

```
terraform import tencentcloud_monitor_grafana_notification_channel.grafana_notification_channel grafana_notification_channel_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudMonitorGrafanaNotificationChannel() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMonitorGrafanaNotificationChannelCreate,
		Read:   resourceTencentCloudMonitorGrafanaNotificationChannelRead,
		Update: resourceTencentCloudMonitorGrafanaNotificationChannelUpdate,
		Delete: resourceTencentCloudMonitorGrafanaNotificationChannelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Grafana instance id.",
			},

			"channel_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Plugin id.",
			},

			"channel_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Channel name.",
			},

			"org_id": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Grafana organization which channel will be installed, default to 1 representing Main Org.",
			},

			"receivers": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Cloud monitor notification template notice-id list.",
			},

			"extra_org_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Extra grafana organization id list, default to 1 representing Main Org.",
			},
		},
	}
}

func resourceTencentCloudMonitorGrafanaNotificationChannelCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_notification_channel.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request     = monitor.NewCreateGrafanaNotificationChannelRequest()
		response    = monitor.NewCreateGrafanaNotificationChannelResponse()
		channelId   string
		channelName string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("channel_name"); ok {
		channelName = v.(string)
		request.ChannelName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("org_id"); ok {
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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateGrafanaNotificationChannel(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create monitor grafanaNotificationChannel failed, reason:%+v", logId, err)
		return err
	}

	channelId = *response.Response.ChannelId
	d.SetId(strings.Join([]string{channelId, channelName}, FILED_SP))

	return resourceTencentCloudMonitorGrafanaNotificationChannelRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaNotificationChannelRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_notification_channel.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	channelName := idSplit[1]

	grafanaNotificationChannel, err := service.DescribeMonitorGrafanaNotificationChannelById(ctx, channelId, channelName)
	if err != nil {
		return err
	}

	if grafanaNotificationChannel == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MonitorGrafanaNotificationChannel` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if grafanaNotificationChannel.InstanceId != nil {
		_ = d.Set("instance_id", grafanaNotificationChannel.InstanceId)
	}

	if grafanaNotificationChannel.ChannelId != nil {
		_ = d.Set("channel_id", grafanaNotificationChannel.ChannelId)
	}

	if grafanaNotificationChannel.ChannelName != nil {
		_ = d.Set("channel_name", grafanaNotificationChannel.ChannelName)
	}

	if grafanaNotificationChannel.OrgId != nil {
		_ = d.Set("org_id", grafanaNotificationChannel.OrgId)
	}

	if grafanaNotificationChannel.Receivers != nil {
		_ = d.Set("receivers", grafanaNotificationChannel.Receivers)
	}

	if grafanaNotificationChannel.ExtraOrgIds != nil {
		_ = d.Set("extra_org_ids", grafanaNotificationChannel.ExtraOrgIds)
	}

	return nil
}

func resourceTencentCloudMonitorGrafanaNotificationChannelUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_notification_channel.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateGrafanaNotificationChannelRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	channelName := idSplit[1]

	request.ChannelId = &channelId
	request.ChannelName = &channelName

	immutableArgs := []string{"instance_id", "channel_id", "channel_name", "org_id", "receivers", "extra_org_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateGrafanaNotificationChannel(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update monitor grafanaNotificationChannel failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMonitorGrafanaNotificationChannelRead(d, meta)
}

func resourceTencentCloudMonitorGrafanaNotificationChannelDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_grafana_notification_channel.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	channelId := idSplit[0]
	channelName := idSplit[1]

	if err := service.DeleteMonitorGrafanaNotificationChannelById(ctx, channelId, channelName); err != nil {
		return err
	}

	return nil
}
