/*
Provides a resource to create a tem scale_rule

Example Usage

```hcl
resource "tencentcloud_tem_scale_rule" "scale_rule" {
  environment_id = "en-xxx"
  application_id = "app-xxx"
  autoscaler {
		autoscaler_name = "test"
		description = "test"
		enabled = true
		min_replicas = 1
		max_replicas = 2
		cron_horizontal_autoscaler {
			name = "test"
			period = "test"
			priority = 1
			enabled = true
			schedules {
				start_at = "03:00"
				target_replicas = 1
			}
		}
		horizontal_autoscaler {
			metrics = "test"
			enabled = true
			max_replicas = 2
			min_replicas = 1
			threshold = 60
		}

  }
}
```

Import

tem scale_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tem_scale_rule.scale_rule scale_rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tem "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tem/v20210701"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTemScaleRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTemScaleRuleCreate,
		Read:   resourceTencentCloudTemScaleRuleRead,
		Update: resourceTencentCloudTemScaleRuleUpdate,
		Delete: resourceTencentCloudTemScaleRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Environment ID.",
			},

			"application_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Application ID.",
			},

			"autoscaler": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"autoscaler_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Description.",
						},
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Enable AutoScaler.",
						},
						"min_replicas": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minimal replica number.",
						},
						"max_replicas": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Maximal replica number.",
						},
						"cron_horizontal_autoscaler": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Scaler based on cron configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Name.",
									},
									"period": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Period.",
									},
									"priority": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Priority.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Enable scaler.",
									},
									"schedules": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "Schedule payload.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"start_at": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Start time.",
												},
												"target_replicas": {
													Type:        schema.TypeInt,
													Required:    true,
													Description: "Target replica number.",
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
							Description: "Scaler based on metrics.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"metrics": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Metric name.",
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Enable scaler.",
									},
									"max_replicas": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Maximal replica number.",
									},
									"min_replicas": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Minimal replica number.",
									},
									"threshold": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Metric threshold.",
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
		response      = tem.NewCreateApplicationAutoscalerResponse()
		environmentId string
		applicationId string
		autoscalerId  string
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
				cronHorizontalAutoscalerMap := item.(map[string]interface{})
				cronHorizontalAutoscaler := tem.CronHorizontalAutoscaler{}
				if v, ok := cronHorizontalAutoscalerMap["name"]; ok {
					cronHorizontalAutoscaler.Name = helper.String(v.(string))
				}
				if v, ok := cronHorizontalAutoscalerMap["period"]; ok {
					cronHorizontalAutoscaler.Period = helper.String(v.(string))
				}
				if v, ok := cronHorizontalAutoscalerMap["priority"]; ok {
					cronHorizontalAutoscaler.Priority = helper.IntInt64(v.(int))
				}
				if v, ok := cronHorizontalAutoscalerMap["enabled"]; ok {
					cronHorizontalAutoscaler.Enabled = helper.Bool(v.(bool))
				}
				if v, ok := cronHorizontalAutoscalerMap["schedules"]; ok {
					for _, item := range v.([]interface{}) {
						schedulesMap := item.(map[string]interface{})
						cronHorizontalAutoscalerSchedule := tem.CronHorizontalAutoscalerSchedule{}
						if v, ok := schedulesMap["start_at"]; ok {
							cronHorizontalAutoscalerSchedule.StartAt = helper.String(v.(string))
						}
						if v, ok := schedulesMap["target_replicas"]; ok {
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
				horizontalAutoscalerMap := item.(map[string]interface{})
				horizontalAutoscaler := tem.HorizontalAutoscaler{}
				if v, ok := horizontalAutoscalerMap["metrics"]; ok {
					horizontalAutoscaler.Metrics = helper.String(v.(string))
				}
				if v, ok := horizontalAutoscalerMap["enabled"]; ok {
					horizontalAutoscaler.Enabled = helper.Bool(v.(bool))
				}
				if v, ok := horizontalAutoscalerMap["max_replicas"]; ok {
					horizontalAutoscaler.MaxReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := horizontalAutoscalerMap["min_replicas"]; ok {
					horizontalAutoscaler.MinReplicas = helper.IntInt64(v.(int))
				}
				if v, ok := horizontalAutoscalerMap["threshold"]; ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tem scaleRule failed, reason:%+v", logId, err)
		return err
	}

	environmentId = *response.Response.EnvironmentId
	d.SetId(strings.Join([]string{environmentId, applicationId, autoscalerId}, FILED_SP))

	return resourceTencentCloudTemScaleRuleRead(d, meta)
}

func resourceTencentCloudTemScaleRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tem_scale_rule.read")()
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
	autoscalerId := idSplit[2]

	scaleRule, err := service.DescribeTemScaleRuleById(ctx, environmentId, applicationId, autoscalerId)
	if err != nil {
		return err
	}

	if scaleRule == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TemScaleRule` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if scaleRule.EnvironmentId != nil {
		_ = d.Set("environment_id", scaleRule.EnvironmentId)
	}

	if scaleRule.ApplicationId != nil {
		_ = d.Set("application_id", scaleRule.ApplicationId)
	}

	if scaleRule.Autoscaler != nil {
		autoscalerMap := map[string]interface{}{}

		if scaleRule.Autoscaler.AutoscalerName != nil {
			autoscalerMap["autoscaler_name"] = scaleRule.Autoscaler.AutoscalerName
		}

		if scaleRule.Autoscaler.Description != nil {
			autoscalerMap["description"] = scaleRule.Autoscaler.Description
		}

		if scaleRule.Autoscaler.Enabled != nil {
			autoscalerMap["enabled"] = scaleRule.Autoscaler.Enabled
		}

		if scaleRule.Autoscaler.MinReplicas != nil {
			autoscalerMap["min_replicas"] = scaleRule.Autoscaler.MinReplicas
		}

		if scaleRule.Autoscaler.MaxReplicas != nil {
			autoscalerMap["max_replicas"] = scaleRule.Autoscaler.MaxReplicas
		}

		if scaleRule.Autoscaler.CronHorizontalAutoscaler != nil {
			cronHorizontalAutoscalerList := []interface{}{}
			for _, cronHorizontalAutoscaler := range scaleRule.Autoscaler.CronHorizontalAutoscaler {
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

					cronHorizontalAutoscalerMap["schedules"] = []interface{}{schedulesList}
				}

				cronHorizontalAutoscalerList = append(cronHorizontalAutoscalerList, cronHorizontalAutoscalerMap)
			}

			autoscalerMap["cron_horizontal_autoscaler"] = []interface{}{cronHorizontalAutoscalerList}
		}

		if scaleRule.Autoscaler.HorizontalAutoscaler != nil {
			horizontalAutoscalerList := []interface{}{}
			for _, horizontalAutoscaler := range scaleRule.Autoscaler.HorizontalAutoscaler {
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

			autoscalerMap["horizontal_autoscaler"] = []interface{}{horizontalAutoscalerList}
		}

		_ = d.Set("autoscaler", []interface{}{autoscalerMap})
	}

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
	autoscalerId := idSplit[2]

	request.EnvironmentId = &environmentId
	request.ApplicationId = &applicationId
	request.AutoscalerId = &autoscalerId

	immutableArgs := []string{"environment_id", "application_id", "autoscaler"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
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
					cronHorizontalAutoscalerMap := item.(map[string]interface{})
					cronHorizontalAutoscaler := tem.CronHorizontalAutoscaler{}
					if v, ok := cronHorizontalAutoscalerMap["name"]; ok {
						cronHorizontalAutoscaler.Name = helper.String(v.(string))
					}
					if v, ok := cronHorizontalAutoscalerMap["period"]; ok {
						cronHorizontalAutoscaler.Period = helper.String(v.(string))
					}
					if v, ok := cronHorizontalAutoscalerMap["priority"]; ok {
						cronHorizontalAutoscaler.Priority = helper.IntInt64(v.(int))
					}
					if v, ok := cronHorizontalAutoscalerMap["enabled"]; ok {
						cronHorizontalAutoscaler.Enabled = helper.Bool(v.(bool))
					}
					if v, ok := cronHorizontalAutoscalerMap["schedules"]; ok {
						for _, item := range v.([]interface{}) {
							schedulesMap := item.(map[string]interface{})
							cronHorizontalAutoscalerSchedule := tem.CronHorizontalAutoscalerSchedule{}
							if v, ok := schedulesMap["start_at"]; ok {
								cronHorizontalAutoscalerSchedule.StartAt = helper.String(v.(string))
							}
							if v, ok := schedulesMap["target_replicas"]; ok {
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
					horizontalAutoscalerMap := item.(map[string]interface{})
					horizontalAutoscaler := tem.HorizontalAutoscaler{}
					if v, ok := horizontalAutoscalerMap["metrics"]; ok {
						horizontalAutoscaler.Metrics = helper.String(v.(string))
					}
					if v, ok := horizontalAutoscalerMap["enabled"]; ok {
						horizontalAutoscaler.Enabled = helper.Bool(v.(bool))
					}
					if v, ok := horizontalAutoscalerMap["max_replicas"]; ok {
						horizontalAutoscaler.MaxReplicas = helper.IntInt64(v.(int))
					}
					if v, ok := horizontalAutoscalerMap["min_replicas"]; ok {
						horizontalAutoscaler.MinReplicas = helper.IntInt64(v.(int))
					}
					if v, ok := horizontalAutoscalerMap["threshold"]; ok {
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tem scaleRule failed, reason:%+v", logId, err)
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
	autoscalerId := idSplit[2]

	if err := service.DeleteTemScaleRuleById(ctx, environmentId, applicationId, autoscalerId); err != nil {
		return err
	}

	return nil
}
