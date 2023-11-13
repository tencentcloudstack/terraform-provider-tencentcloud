/*
Use this data source to query detailed information of tcr webhook_trigger_log

Example Usage

```hcl
data "tencentcloud_tcr_webhook_trigger_log" "webhook_trigger_log" {
  registry_id = "tcr-xxx"
  namespace = "nginx"
    tags = {
    "createdBy" = "terraform"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrWebhookTriggerLog() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrWebhookTriggerLogRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance Id.",
			},

			"namespace": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace.",
			},

			"logs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Log list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Log id.",
						},
						"trigger_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Trigger Id.",
						},
						"event_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event type.",
						},
						"notify_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Notification type.",
						},
						"detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Webhook trigger detail.",
						},
						"creation_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status.",
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

func dataSourceTencentCloudTcrWebhookTriggerLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_webhook_trigger_log.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["Namespace"] = helper.String(v.(string))
	}

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

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

			ids = append(ids, *webhookTriggerLog.RegistryId)
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
