package tmp

import (
	"context"
	"fmt"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpTkeGlobalNotification() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpTkeGlobalNotificationRead,
		Create: resourceTencentCloudMonitorTmpTkeGlobalNotificationCreate,
		Update: resourceTencentCloudMonitorTmpTkeGlobalNotificationUpdate,
		Delete: resourceTencentCloudMonitorTmpTkeGlobalNotificationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Id.",
			},

			"notification": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "Alarm notification channels.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Alarm notification switch.",
						},
						"type": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"amp", "webhook", "alertmanager"}),
							Description:  "Alarm notification type, Valid values: `amp`, `webhook`, `alertmanager`.",
						},
						"web_hook": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Web hook, if Type is `webhook`, this field is required.",
						},
						"alert_manager": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Alert manager, if Type is `alertmanager`, this field is required.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"url": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Alert manager url.",
									},
									"cluster_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cluster type.",
									},
									"cluster_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Cluster id.",
									},
								},
							},
						},
						"repeat_interval": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Convergence time.",
						},
						"time_range_start": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective start time.",
						},
						"time_range_end": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Effective end time.",
						},
						"notify_way": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"SMS", "EMAIL", "CALL", "WECHAT"}),
							},
							Optional:    true,
							Description: "Alarm notification method, Valid values: `SMS`, `EMAIL`, `CALL`, `WECHAT`.",
						},
						"receiver_groups": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Alarm receiving group(user group).",
						},
						"phone_notify_order": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Optional:    true,
							Description: "Phone alert sequence, NotifyWay is `CALL`, and this parameter is used.",
						},
						"phone_circle_times": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of phone alerts (user group), NotifyWay is `CALL`, and this parameter is used.",
						},
						"phone_inner_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Interval between telephone alarm rounds, NotifyWay is `CALL`, and this parameter is used.",
						},
						"phone_circle_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Telephone alarm off-wheel interval, NotifyWay is `CALL`, and this parameter is used.",
						},
						"phone_arrive_notice": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Phone Alarm Reach Notification, NotifyWay is `CALL`, and this parameter is used.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMonitorTmpTkeGlobalNotificationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_global_notification.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	instanceId := ""
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudMonitorTmpTkeGlobalNotificationUpdate(d, meta)
}

func resourceTencentCloudMonitorTmpTkeGlobalNotificationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_global_notification.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	instanceId := d.Id()

	globalNotification, err := service.DescribeTkeTmpGlobalNotification(ctx, instanceId)

	if err != nil {
		return err
	}

	if globalNotification == nil {
		d.SetId("")
		return fmt.Errorf("resource `global_notification` %s does not exist", instanceId)
	}

	if *globalNotification.Enabled {
		_ = d.Set("instance_id", instanceId)
		alertManager := make([]map[string]interface{}, 0)
		if globalNotification.AlertManager != nil {
			alertManager = append(alertManager, map[string]interface{}{
				"url":          globalNotification.AlertManager.Url,
				"cluster_type": globalNotification.AlertManager.ClusterType,
				"cluster_id":   globalNotification.AlertManager.ClusterId,
			})
		}

		var notifications []map[string]interface{}
		notification := make(map[string]interface{})
		notification["enabled"] = globalNotification.Enabled
		notification["type"] = globalNotification.Type
		notification["web_hook"] = globalNotification.WebHook
		notification["alert_manager"] = alertManager
		notification["repeat_interval"] = globalNotification.RepeatInterval
		notification["time_range_start"] = globalNotification.TimeRangeStart
		notification["time_range_end"] = globalNotification.TimeRangeEnd
		notification["notify_way"] = globalNotification.NotifyWay
		notification["receiver_groups"] = globalNotification.ReceiverGroups
		notification["phone_notify_order"] = globalNotification.PhoneNotifyOrder
		notification["phone_circle_times"] = globalNotification.PhoneCircleTimes
		notification["phone_inner_interval"] = globalNotification.PhoneInnerInterval
		notification["phone_circle_interval"] = globalNotification.PhoneCircleInterval
		notification["phone_arrive_notice"] = globalNotification.PhoneArriveNotice
		notifications = append(notifications, notification)
		_ = d.Set("notification", notifications)
	}

	return nil
}

func resourceTencentCloudMonitorTmpTkeGlobalNotificationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_global_notification.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	if d.HasChange("notification") {
		notification := monitor.PrometheusNotificationItem{}
		if dMap, ok := helper.InterfacesHeadMap(d, "notification"); ok {
			if v, ok := dMap["enabled"]; ok {
				notification.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["type"]; ok {
				notification.Type = helper.String(v.(string))
			}
			if v, ok := dMap["web_hook"]; ok {
				notification.WebHook = helper.String(v.(string))
			}
			if v, ok := helper.InterfacesHeadMap(d, "alert_manager"); ok {
				alertManager := monitor.PrometheusAlertManagerConfig{}
				if vv, ok := v["url"]; ok {
					alertManager.Url = helper.String(vv.(string))
				}
				if vv, ok := v["cluster_type"]; ok {
					alertManager.ClusterType = helper.String(vv.(string))
				}
				if vv, ok := v["cluster_id"]; ok {
					alertManager.ClusterId = helper.String(vv.(string))
				}
				notification.AlertManager = &alertManager
			}

			if v, ok := dMap["repeat_interval"]; ok {
				notification.RepeatInterval = helper.String(v.(string))
			}
			if v, ok := dMap["time_range_start"]; ok {
				notification.TimeRangeStart = helper.String(v.(string))
			}
			if v, ok := dMap["time_range_end"]; ok {
				notification.TimeRangeEnd = helper.String(v.(string))
			}
			if v, ok := dMap["notify_way"]; ok {
				for _, vv := range v.(*schema.Set).List() {
					if vv == "CALL" {
						if v, ok := dMap["receiver_groups"]; ok {
							notification.ReceiverGroups = helper.Strings(v.([]string))
						}
						if v, ok := dMap["phone_notify_order"]; ok {
							notification.PhoneNotifyOrder = helper.InterfacesUint64Point(v.([]interface{}))
						}
						if v, ok := dMap["phone_circle_times"]; ok {
							notification.PhoneCircleTimes = helper.Int64(v.(int64))
						}
						if v, ok := dMap["phone_inner_interval"]; ok {
							notification.PhoneInnerInterval = helper.Int64(v.(int64))
						}
						if v, ok := dMap["phone_circle_interval"]; ok {
							notification.PhoneCircleInterval = helper.Int64(v.(int64))
						}
						if v, ok := dMap["phone_arrive_notice"]; ok {
							notification.PhoneArriveNotice = helper.Bool(v.(bool))
						}
					}
					notification.NotifyWay = append(notification.NotifyWay, helper.String(vv.(string)))
				}
			}
		}

		if _, err := service.ModifyTkeTmpGlobalNotification(ctx, d.Id(), notification); err != nil {
			return err
		}
	}

	return resourceTencentCloudMonitorTmpTkeGlobalNotificationRead(d, meta)
}

func resourceTencentCloudMonitorTmpTkeGlobalNotificationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_global_notification.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	id := d.Id()
	notification := monitor.PrometheusNotificationItem{
		// Turning off the alarm notification function is to delete the alarm notification
		Enabled: helper.Bool(false),
		Type:    helper.String(""),
	}

	if _, err := service.ModifyTkeTmpGlobalNotification(ctx, id, notification); err != nil {
		return err
	}

	return nil
}
