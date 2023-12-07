package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrWebhookTriggerLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrWebhookTriggerLogsRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "instance Id.",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "namespace.",
			},

			"trigger_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "trigger id.",
			},

			"logs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "log list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "log id.",
						},
						"trigger_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "trigger Id.",
						},
						"event_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "event type.",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "notification type.",
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "webhook trigger detail.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "update time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "status.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTcrWebhookTriggerLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_webhook_trigger_logs.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["registry_id"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["namespace"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("trigger_id"); ok {
		paramMap["trigger_id"] = helper.IntInt64(v.(int))
	}

	service := TCRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var logs []*tcr.WebhookTriggerLog

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrWebhookTriggerLogByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		logs = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(logs))
	tmpList := make([]map[string]interface{}, 0, len(logs))

	if logs != nil {
		for _, webhookTriggerLog := range logs {
			webhookTriggerLogMap := map[string]interface{}{}

			if webhookTriggerLog.Id != nil {
				webhookTriggerLogMap["id"] = webhookTriggerLog.Id
			}

			if webhookTriggerLog.TriggerId != nil {
				webhookTriggerLogMap["trigger_id"] = webhookTriggerLog.TriggerId
			}

			if webhookTriggerLog.EventType != nil {
				webhookTriggerLogMap["event_type"] = webhookTriggerLog.EventType
			}

			if webhookTriggerLog.NotifyType != nil {
				webhookTriggerLogMap["notify_type"] = webhookTriggerLog.NotifyType
			}

			if webhookTriggerLog.Detail != nil {
				webhookTriggerLogMap["detail"] = webhookTriggerLog.Detail
			}

			if webhookTriggerLog.CreationTime != nil {
				webhookTriggerLogMap["creation_time"] = webhookTriggerLog.CreationTime
			}

			if webhookTriggerLog.UpdateTime != nil {
				webhookTriggerLogMap["update_time"] = webhookTriggerLog.UpdateTime
			}

			if webhookTriggerLog.Status != nil {
				webhookTriggerLogMap["status"] = webhookTriggerLog.Status
			}

			ids = append(ids, helper.Int64ToStr(*webhookTriggerLog.Id))
			tmpList = append(tmpList, webhookTriggerLogMap)
		}

		_ = d.Set("logs", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
