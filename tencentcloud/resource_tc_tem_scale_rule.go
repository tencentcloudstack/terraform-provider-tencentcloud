/*
Provides a resource to create a tem scaleRule

Example Usage

```hcl
resource "tencentcloud_tem_scale_rule" "scaleRule" {
  environment_id = "en-853mggjm"
  application_id = "app-3j29aa2p"
  autoscaler {
    autoscaler_name = "test3123"
    description     = "test"
    enabled         = true
    min_replicas    = 1
    max_replicas    = 3
    cron_horizontal_autoscaler {
      name     = "test"
      period   = "* * *"
      priority = 1
      enabled  = true
      schedules {
        start_at        = "03:00"
        target_replicas = 1
      }
    }
    cron_horizontal_autoscaler {
      name     = "test123123"
      period   = "* * *"
      priority = 0
      enabled  = true
      schedules {
        start_at        = "04:13"
        target_replicas = 1
      }
    }
    horizontal_autoscaler {
      metrics      = "CPU"
      enabled      = true
      max_replicas = 3
      min_replicas = 1
      threshold    = 60
    }

  }
}

```
Import

tem scaleRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_tem_scale_rule.scaleRule environmentId#applicationId#scaleRuleId
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
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTemScaleRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTemScaleRuleRead,
		Create: resourceTencentCloudTemScaleRuleCreate,
		Update: resourceTencentCloudTemScaleRuleUpdate,
		Delete: resourceTencentCloudTemScaleRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "environment ID.",
			},

			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "application ID.",
			},

			"autoscaler": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoscaler_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "description.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "enable AutoScaler.",
						},
						"min_replicas": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "minimal replica number.",
						},
						"max_replicas": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "maximal replica number.",
						},
						"cron_horizontal_autoscaler": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "scaler based on cron configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "name.",
									},
									"period": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "period.",
									},
									"priority": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "priority.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "enable scaler.",
									},
									"schedules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "schedule payload.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_at": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "start time.",
												},
												"target_replicas": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "target replica number.",
												},
											},
										},
									},
								},
							},
						},
						"horizontal_autoscaler": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "scaler based on metrics.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "metric name.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "enable scaler.",
									},
									"max_replicas": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "maximal replica number.",
									},
									"min_replicas": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "minimal replica number.",
									},
									"threshold": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "metric threshold.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTemScaleRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_scale_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = tem.NewCreateApplicationAutoscalerRequest()
		response      *tem.CreateApplicationAutoscalerResponse
		environmentId string
		applicationId string
	)

	if v, ok := d.GetOk("environment_id"); ok {
		environmentId = v.(string)
		request.EnvironmentId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("application_id"); ok {
		applicationId = v.(string)
		request.ApplicationId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "autoscaler"); ok {
		autoscaler := tem.Autoscaler{}
		if v, ok := dMap["autoscaler_name"]; ok {
			autoscaler.AutoscalerName = helper.String(v.(string))
		}
		if v, ok := dMap["description"]; ok {
			autoscaler.Description = helper.String(v.(string))
		}
		if v, ok := dMap["enabled"]; ok {
			autoscaler.Enabled = helper.Bool(v.(bool))
		}
		if v, ok := dMap["min_replicas"]; ok {
			autoscaler.MinReplicas = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["max_replicas"]; ok {
			autoscaler.MaxReplicas = helper.IntInt64(v.(int))
		}
		if v, ok := dMap["cron_horizontal_autoscaler"]; ok {
			for _, item := range v.([]interface{}) {
				CronHorizontalAutoscalerMap := item.(map[string]interface{})
				cronHorizontalAutoscaler := tem.CronHorizontalAutoscaler{}
				if v, ok := CronHorizontalAutoscalerMap["name"]; ok {
					cronHorizontalAutoscaler.Name = helper.String(v.(string))
				}
				if v, ok := CronHorizontalAutoscalerMap["period"]; ok {
					cronHorizontalAutoscaler.Period = helper.String(v.(string))
				}
				if v, ok := CronHorizontalAutoscalerMap["priority"]; ok {
					cronHorizontalAutoscaler.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := CronHorizontalAutoscalerMap["enabled"]; ok {
					cronHorizontalAutoscaler.Enabled = helper.Bool(v.(bool))
				}
				if v, ok := CronHorizontalAutoscalerMap["schedules"]; ok {
					for _, item := range v.([]interface{}) {
						SchedulesMap := item.(map[string]interface{})
						cronHorizontalAutoscalerSchedule := tem.CronHorizontalAutoscalerSchedule{}
						if v, ok := SchedulesMap["start_at"]; ok {
							cronHorizontalAutoscalerSchedule.StartAt = helper.String(v.(string))
						}
						if v, ok := SchedulesMap["target_replicas"]; ok {
							cronHorizontalAutoscalerSchedule.TargetReplicas = helper.IntInt64(v.(int))
						}
						cronHorizontalAutoscaler.Schedules = append(cronHorizontalAutoscaler.Schedules, &cronHorizontalAutoscalerSchedule)
					}
				}
				autoscaler.CronHorizontalAutoscaler = append(autoscaler.CronHorizontalAutoscaler, &cronHorizontalAutoscaler)
			}
		}
		if v, ok := dMap["horizontal_autoscaler"]; ok {
			for _, item := range v.([]interface{}) {
				HorizontalAutoscalerMap := item.(map[string]interface{})
				horizontalAutoscaler := tem.HorizontalAutoscaler{}
				if v, ok := HorizontalAutoscalerMap["metrics"]; ok {
					horizontalAutoscaler.Metrics = helper.String(v.(string))
				}
				if v, ok := HorizontalAutoscalerMap["enabled"]; ok {
					horizontalAutoscaler.Enabled = helper.Bool(v.(bool))
				}
				if v, ok := HorizontalAutoscalerMap["max_replicas"]; ok {
					horizontalAutoscaler.MaxReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := HorizontalAutoscalerMap["min_replicas"]; ok {
					horizontalAutoscaler.MinReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := HorizontalAutoscalerMap["threshold"]; ok {
					horizontalAutoscaler.Threshold = helper.IntInt64(v.(int))
				}
				autoscaler.HorizontalAutoscaler = append(autoscaler.HorizontalAutoscaler, &horizontalAutoscaler)
			}
		}
		request.Autoscaler = &autoscaler
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().CreateApplicationAutoscaler(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tem scaleRule failed, reason:%+v", logId, err)
		return err
	}

	scaleRuleId := *response.Response.Result

	d.SetId(environmentId + FILED_SP + applicationId + FILED_SP + scaleRuleId)
	return resourceTencentCloudTemScaleRuleRead(d, meta)
}

func resourceTencentCloudTemScaleRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_scaleRule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	scaleRuleId := idSplit[2]

	scaleRule, err := service.DescribeTemScaleRule(ctx, environmentId, applicationId, scaleRuleId)

	if err != nil {
		return err
	}

	if scaleRule == nil {
		d.SetId("")
		return fmt.Errorf("resource `scaleRule` %s does not exist", scaleRuleId)
	}

	_ = d.Set("environment_id", environmentId)
	_ = d.Set("application_id", applicationId)

	autoscalerMap := map[string]interface{}{}
	if scaleRule.AutoscalerName != nil {
		autoscalerMap["autoscaler_name"] = scaleRule.AutoscalerName
	}
	if scaleRule.Description != nil {
		autoscalerMap["description"] = scaleRule.Description
	}
	if scaleRule.Enabled != nil {
		autoscalerMap["enabled"] = scaleRule.Enabled
	}
	if scaleRule.MinReplicas != nil {
		autoscalerMap["min_replicas"] = scaleRule.MinReplicas
	}
	if scaleRule.MaxReplicas != nil {
		autoscalerMap["max_replicas"] = scaleRule.MaxReplicas
	}
	if scaleRule.CronHorizontalAutoscaler != nil {
		cronHorizontalAutoscalerList := []interface{}{}
		for _, cronHorizontalAutoscaler := range scaleRule.CronHorizontalAutoscaler {
			cronHorizontalAutoscalerMap := map[string]interface{}{}
			if cronHorizontalAutoscaler.Name != nil {
				cronHorizontalAutoscalerMap["name"] = cronHorizontalAutoscaler.Name
			}
			if cronHorizontalAutoscaler.Period != nil {
				cronHorizontalAutoscalerMap["period"] = cronHorizontalAutoscaler.Period
			}
			if cronHorizontalAutoscaler.Priority != nil {
				cronHorizontalAutoscalerMap["priority"] = cronHorizontalAutoscaler.Priority
			}
			if cronHorizontalAutoscaler.Enabled != nil {
				cronHorizontalAutoscalerMap["enabled"] = cronHorizontalAutoscaler.Enabled
			}
			if cronHorizontalAutoscaler.Schedules != nil {
				schedulesList := []interface{}{}
				for _, schedules := range cronHorizontalAutoscaler.Schedules {
					schedulesMap := map[string]interface{}{}
					if schedules.StartAt != nil {
						schedulesMap["start_at"] = schedules.StartAt
					}
					if schedules.TargetReplicas != nil {
						schedulesMap["target_replicas"] = schedules.TargetReplicas
					}

					schedulesList = append(schedulesList, schedulesMap)
				}
				cronHorizontalAutoscalerMap["schedules"] = schedulesList
			}

			cronHorizontalAutoscalerList = append(cronHorizontalAutoscalerList, cronHorizontalAutoscalerMap)
		}
		autoscalerMap["cron_horizontal_autoscaler"] = cronHorizontalAutoscalerList
	}
	if scaleRule.HorizontalAutoscaler != nil {
		horizontalAutoscalerList := []interface{}{}
		for _, horizontalAutoscaler := range scaleRule.HorizontalAutoscaler {
			horizontalAutoscalerMap := map[string]interface{}{}
			if horizontalAutoscaler.Metrics != nil {
				horizontalAutoscalerMap["metrics"] = horizontalAutoscaler.Metrics
			}
			if horizontalAutoscaler.Enabled != nil {
				horizontalAutoscalerMap["enabled"] = horizontalAutoscaler.Enabled
			}
			if horizontalAutoscaler.MaxReplicas != nil {
				horizontalAutoscalerMap["max_replicas"] = horizontalAutoscaler.MaxReplicas
			}
			if horizontalAutoscaler.MinReplicas != nil {
				horizontalAutoscalerMap["min_replicas"] = horizontalAutoscaler.MinReplicas
			}
			if horizontalAutoscaler.Threshold != nil {
				horizontalAutoscalerMap["threshold"] = horizontalAutoscaler.Threshold
			}

			horizontalAutoscalerList = append(horizontalAutoscalerList, horizontalAutoscalerMap)
		}
		autoscalerMap["horizontal_autoscaler"] = horizontalAutoscalerList
	}

	_ = d.Set("autoscaler", []interface{}{autoscalerMap})

	return nil
}

func resourceTencentCloudTemScaleRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_scale_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := tem.NewModifyApplicationAutoscalerRequest()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	scaleRuleId := idSplit[2]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.AutoscalerId = &scaleRuleId

	if d.HasChange("environment_id") {
		return fmt.Errorf("`environment_id` do not support change now.")
	}

	if d.HasChange("application_id") {
		return fmt.Errorf("`application_id` do not support change now.")
	}

	if d.HasChange("autoscaler") {
		if dMap, ok := helper.InterfacesHeadMap(d, "autoscaler"); ok {
			autoscaler := tem.Autoscaler{}
			if v, ok := dMap["autoscaler_name"]; ok {
				autoscaler.AutoscalerName = helper.String(v.(string))
			}
			if v, ok := dMap["description"]; ok {
				autoscaler.Description = helper.String(v.(string))
			}
			if v, ok := dMap["enabled"]; ok {
				autoscaler.Enabled = helper.Bool(v.(bool))
			}
			if v, ok := dMap["min_replicas"]; ok {
				autoscaler.MinReplicas = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["max_replicas"]; ok {
				autoscaler.MaxReplicas = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["cron_horizontal_autoscaler"]; ok {
				for _, item := range v.([]interface{}) {
					CronHorizontalAutoscalerMap := item.(map[string]interface{})
					cronHorizontalAutoscaler := tem.CronHorizontalAutoscaler{}
					if v, ok := CronHorizontalAutoscalerMap["name"]; ok {
						cronHorizontalAutoscaler.Name = helper.String(v.(string))
					}
					if v, ok := CronHorizontalAutoscalerMap["period"]; ok {
						cronHorizontalAutoscaler.Period = helper.String(v.(string))
					}
					if v, ok := CronHorizontalAutoscalerMap["priority"]; ok {
						cronHorizontalAutoscaler.Priority = helper.IntInt64(v.(int))
					}
					if v, ok := CronHorizontalAutoscalerMap["enabled"]; ok {
						cronHorizontalAutoscaler.Enabled = helper.Bool(v.(bool))
					}
					if v, ok := CronHorizontalAutoscalerMap["schedules"]; ok {
						for _, item := range v.([]interface{}) {
							SchedulesMap := item.(map[string]interface{})
							cronHorizontalAutoscalerSchedule := tem.CronHorizontalAutoscalerSchedule{}
							if v, ok := SchedulesMap["start_at"]; ok {
								cronHorizontalAutoscalerSchedule.StartAt = helper.String(v.(string))
							}
							if v, ok := SchedulesMap["target_replicas"]; ok {
								cronHorizontalAutoscalerSchedule.TargetReplicas = helper.IntInt64(v.(int))
							}
							cronHorizontalAutoscaler.Schedules = append(cronHorizontalAutoscaler.Schedules, &cronHorizontalAutoscalerSchedule)
						}
					}
					autoscaler.CronHorizontalAutoscaler = append(autoscaler.CronHorizontalAutoscaler, &cronHorizontalAutoscaler)
				}
			}
			if v, ok := dMap["horizontal_autoscaler"]; ok {
				for _, item := range v.([]interface{}) {
					HorizontalAutoscalerMap := item.(map[string]interface{})
					horizontalAutoscaler := tem.HorizontalAutoscaler{}
					if v, ok := HorizontalAutoscalerMap["metrics"]; ok {
						horizontalAutoscaler.Metrics = helper.String(v.(string))
					}
					if v, ok := HorizontalAutoscalerMap["enabled"]; ok {
						horizontalAutoscaler.Enabled = helper.Bool(v.(bool))
					}
					if v, ok := HorizontalAutoscalerMap["max_replicas"]; ok {
						horizontalAutoscaler.MaxReplicas = helper.IntInt64(v.(int))
					}
					if v, ok := HorizontalAutoscalerMap["min_replicas"]; ok {
						horizontalAutoscaler.MinReplicas = helper.IntInt64(v.(int))
					}
					if v, ok := HorizontalAutoscalerMap["threshold"]; ok {
						horizontalAutoscaler.Threshold = helper.IntInt64(v.(int))
					}
					autoscaler.HorizontalAutoscaler = append(autoscaler.HorizontalAutoscaler, &horizontalAutoscaler)
				}
			}
			request.Autoscaler = &autoscaler
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTemClient().ModifyApplicationAutoscaler(request)
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

	return resourceTencentCloudTemScaleRuleRead(d, meta)
}

func resourceTencentCloudTemScaleRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_scale_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TemService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	environmentId := idSplit[0]
	applicationId := idSplit[1]
	scaleRuleId := idSplit[2]

	if err := service.DisableTemScaleRuleById(ctx, environmentId, applicationId, scaleRuleId); err != nil {
		return err
	}

	if err := service.DeleteTemScaleRuleById(ctx, environmentId, applicationId, scaleRuleId); err != nil {
		return err
	}

	return nil
}
