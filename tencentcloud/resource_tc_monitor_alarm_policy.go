/*
Provides a alarm policy resource for monitor.

Example Usage

cvm_device alarm policy

```hcl
resource "tencentcloud_monitor_alarm_policy" "group" {
  policy_name = "hello"
  monitor_type = "MT_QCE"
  enable = 1
  project_id = 1244035
  namespace = "cvm_device"
  conditions {
    is_union_rule = 1
    rules {
      metric_name = "CpuUsage"
      period = 60
      operator = "ge"
      value = "89.9"
      continue_period = 1
      notice_frequency = 3600
      is_power_notice = 0
    }
  }
  event_conditions {
    metric_name = "ping_unreachable"
  }
  event_conditions {
    metric_name = "guest_reboot"
  }
  notice_ids = ["notice-l9ziyxw6"]

  trigger_tasks {
    type = "AS"
    task_config = "{\"Region\":\"ap-guangzhou\",\"Group\":\"asg-0z312312x\",\"Policy\":\"asp-ganig28\"}"
  }
}
```

k8s_cluster alarm policy

```hcl
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "k8s_cluster"
  notice_ids   = [
    "notice-l9ziyxw6",
  ]
  policy_name  = "TkeClusterNew"
  project_id   = 1244035

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 3
      description      = "Allocatable Pods"
      is_power_notice  = 0
      metric_name      = "K8sClusterAllocatablePodsTotal"
      notice_frequency = 3600
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "Count"
      value            = "10"

      filter {
        dimensions = jsonencode(
        [
          [
            {
              Key      = "region"
              Operator = "eq"
              Value    = [
                "ap-guangzhou",
              ]
            },
            {
              Key      = "tke_cluster_instance_id"
              Operator = "in"
              Value    = [
                "cls-czhtobea",
              ]
            },
          ],
        ]
        )
        type       = "DIMENSION"
      }
    }
    rules {
      continue_period  = 3
      description      = "Total CPU Cores"
      is_power_notice  = 0
      metric_name      = "K8sClusterCpuCoreTotal"
      notice_frequency = 3600
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "Core"
      value            = "2"

      filter {
        dimensions = jsonencode(
        [
          [
            {
              Key      = "region"
              Operator = "eq"
              Value    = [
                "ap-guangzhou",
              ]
            },
            {
              Key      = "tke_cluster_instance_id"
              Operator = "in"
              Value    = [
                "cls-czhtobea",
              ]
            },
          ],
        ]
        )
        type       = "DIMENSION"
      }
    }
  }
}
```

cvm_device alarm policy binding cvm by tag
```
resource "tencentcloud_monitor_alarm_policy" "policy" {
  enable       = 1
  monitor_type = "MT_QCE"
  namespace    = "cvm_device"
  notice_ids   = [
    "notice-l9ziyxw6",
  ]
  policy_name  = "policy"
  project_id   = 0

  conditions {
    is_union_rule = 0

    rules {
      continue_period  = 5
      description      = "CPUUtilization"
      is_power_notice  = 0
      metric_name      = "CpuUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "PublicBandwidthUtilization"
      is_power_notice  = 0
      metric_name      = "Outratio"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "MemoryUtilization"
      is_power_notice  = 0
      metric_name      = "MemUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
    rules {
      continue_period  = 5
      description      = "DiskUtilization"
      is_power_notice  = 0
      metric_name      = "CvmDiskUsage"
      notice_frequency = 7200
      operator         = "gt"
      period           = 60
      rule_type        = "STATIC"
      unit             = "%"
      value            = "95"
    }
  }

  event_conditions {
    continue_period  = 0
    description      = "DiskReadonly"
    is_power_notice  = 0
    metric_name      = "disk_readonly"
    notice_frequency = 0
    period           = 0
  }

  policy_tag {
    key   = "test-tag"
    value = "unit-test"
  }
}
```

Import

Alarm policy instance can be imported, e.g.

```
$ terraform import tencentcloud_monitor_alarm_policy.policy policy-id
```

*/
package tencentcloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func AlarmPolicyRule() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"metric_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metric name or event name.",
		},
		"period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Statistical period in seconds.",
		},
		"operator": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Operator.",
		},
		"value": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Threshold.",
		},
		"continue_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Number of periods.",
		},
		"notice_frequency": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Alarm interval in seconds.",
		},
		"is_power_notice": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Whether the alarm frequency increases exponentially.",
		},
		"filter": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Filter condition for one single trigger rule. Must set it when create tke-xxx rules.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"type": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Filter condition type. Valid values: DIMENSION (uses dimensions for filtering).",
					},
					"dimensions": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "JSON string generated by serializing the AlarmPolicyDimension two-dimensional array.",
					},
				},
			},
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Metric display name, which is used in the output parameter.",
		},
		"unit": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Unit, which is used in the output parameter.",
		},
		"rule_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Trigger condition type.",
		},
	}
}

func resourceTencentCloudMonitorAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentMonitorAlarmPolicyCreate,
		Read:   resourceTencentMonitorAlarmPolicyRead,
		Update: resourceTencentMonitorAlarmPolicyUpdate,
		Delete: resourceTencentMonitorAlarmPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"policy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of policy.",
			},
			"monitor_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of monitor.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The type of alarm.",
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringLengthInRange(0, 100),
				Description:  "The remark of policy group.",
			},
			"enable": {
				Type:        schema.TypeInt,
				Default:     1,
				Optional:    true,
				Description: "Whether to enable, default is `1`.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     -1,
				Description: "Project ID. For products with different projects, a value other than -1 must be passed in.",
			},
			"conditon_template_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of trigger condition template.",
			},
			"conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Computed:    true,
				Description: "A list of metric trigger condition.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_union_rule": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validateAllowedIntValue([]int{0, 1}),
							Description:  "The and or relation of indicator alarm rule.",
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: AlarmPolicyRule(),
							},
							Description: "A list of metric trigger condition.",
						},
					},
				},
			},
			"event_conditions": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A list of event trigger condition.",
				Elem: &schema.Resource{
					Schema: AlarmPolicyRule(),
				},
			},
			"notice_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of notification rule IDs.",
				Elem: &schema.Schema{
					Type:        schema.TypeString,
					Description: "ID of the notification rule to be queried.",
				},
			},
			"trigger_tasks": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Triggered task list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Triggered task type.",
						},
						"task_config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration information in JSON format.",
						},
					},
				},
			},
			"policy_tag": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "Policy tag to bind object.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},
			// compute
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm policy create time.",
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The alarm policy update time.",
			},
		},
	}
}

func resourceTencentMonitorAlarmPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy.create")()
	logId := getLogId(contextNil)
	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewCreateAlarmPolicyRequest()
	)
	request.Module = helper.String("monitor")
	request.PolicyName = helper.String(d.Get("policy_name").(string))
	request.MonitorType = helper.String(d.Get("monitor_type").(string))
	request.Namespace = helper.String(d.Get("namespace").(string))

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	enable := d.Get("enable").(int)
	request.Enable = helper.IntInt64(enable)

	//if v, ok := d.GetOk("enable"); ok {
	//	request.Enable = helper.IntInt64(v.(int))
	//}

	projectId := d.Get("project_id").(int)
	if projectId != -1 {
		request.ProjectId = helper.IntInt64(projectId)
	}

	if v, ok := d.GetOk("conditon_template_id"); ok {
		request.ConditionTemplateId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("conditions"); ok {
		conditions := v.([]interface{})

		if len(conditions) != 1 {
			return fmt.Errorf("need only one conditions.")
		}

		condition := conditions[0].(map[string]interface{})
		var policy = monitor.AlarmPolicyCondition{}
		policy.IsUnionRule = helper.IntInt64(condition["is_union_rule"].(int))
		policy.Rules = make([]*monitor.AlarmPolicyRule, 0, 10)

		rules := condition["rules"]

		for _, item := range rules.([]interface{}) {
			m := item.(map[string]interface{})
			alarmPolicyRule := monitor.AlarmPolicyRule{}
			if m["metric_name"] != nil {
				alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
			}
			if m["period"] != nil {
				alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
			}
			if m["operator"] != nil {
				alarmPolicyRule.Operator = helper.String(m["operator"].(string))
			}
			if m["value"] != nil {
				alarmPolicyRule.Value = helper.String(m["value"].(string))
			}
			if m["continue_period"] != nil {
				alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			if m["notice_frequency"] != nil {
				alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
			}
			if m["is_power_notice"] != nil {
				alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
			}
			if v, ok := m["filter"]; ok {
				filters := v.([]interface{})
				if len(filters) > 0 {
					filter := filters[0].(map[string]interface{})
					alarmPolicyFilter := monitor.AlarmPolicyFilter{
						Type:       helper.String(filter["type"].(string)),
						Dimensions: helper.String(filter["dimensions"].(string)),
					}
					alarmPolicyRule.Filter = &alarmPolicyFilter
				}
			}

			if m["description"] != nil {
				alarmPolicyRule.Description = helper.String(m["description"].(string))
			}
			if m["unit"] != nil {
				alarmPolicyRule.Unit = helper.String(m["unit"].(string))
			}
			if m["rule_type"] != nil {
				alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
			}
			policy.Rules = append(policy.Rules, &alarmPolicyRule)
		}
		request.Condition = &policy
	}

	if v, ok := d.GetOk("event_conditions"); ok {
		eventCondition := monitor.AlarmPolicyEventCondition{}
		rules := make([]*monitor.AlarmPolicyRule, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			alarmPolicyRule := monitor.AlarmPolicyRule{}
			if m["metric_name"] != nil {
				alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
			}
			if m["period"] != nil {
				alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
			}
			if m["operator"] != nil {
				alarmPolicyRule.Operator = helper.String(m["operator"].(string))
			}
			if m["value"] != nil {
				alarmPolicyRule.Value = helper.String(m["value"].(string))
			}
			if m["continue_period"] != nil {
				alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
			}
			if m["notice_frequency"] != nil {
				alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
			}
			if m["is_power_notice"] != nil {
				alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
			}
			if m["filter"] != nil {
				filters := m["filter"].([]interface{})
				if len(filters) > 0 {
					filter := filters[0].(map[string]interface{})
					alarmPolicyFilter := monitor.AlarmPolicyFilter{
						Type:       helper.String(filter["type"].(string)),
						Dimensions: helper.String(filter["dimensions"].(string)),
					}
					alarmPolicyRule.Filter = &alarmPolicyFilter
				}
			}
			if m["description"] != nil {
				alarmPolicyRule.Description = helper.String(m["description"].(string))
			}
			if m["unit"] != nil {
				alarmPolicyRule.Unit = helper.String(m["unit"].(string))
			}
			if m["rule_type"] != nil {
				alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
			}
			rules = append(rules, &alarmPolicyRule)
		}
		eventCondition.Rules = rules
		request.EventCondition = &eventCondition
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		notice := make([]*string, 0, 10)
		for _, item := range v.([]interface{}) {
			notice = append(notice, helper.String(item.(string)))
		}
		request.NoticeIds = notice
	}

	if v, ok := d.GetOk("trigger_tasks"); ok {
		tasks := make([]*monitor.AlarmPolicyTriggerTask, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			triggerTask := monitor.AlarmPolicyTriggerTask{}
			triggerTask.Type = helper.String(m["type"].(string))
			triggerTask.TaskConfig = helper.String(m["task_config"].(string))
			tasks = append(tasks, &triggerTask)
		}
		request.TriggerTasks = tasks
	}

	var groupId *string
	var policyId *string
	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().CreateAlarmPolicy(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			return retryError(err, InternalError)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		policyId = response.Response.PolicyId
		groupId = response.Response.OriginId
		return nil
	}); err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s", *policyId))

	// binding tag
	if v, ok := d.GetOk("policy_tag"); ok {
		request := monitor.NewBindingPolicyTagRequest()

		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(*policyId)
		request.ServiceType = helper.String(d.Get("namespace").(string))
		request.GroupId = helper.String(*groupId)
		tagSet := make([]*monitor.PolicyTag, 0, 10)
		for _, item := range v.([]interface{}) {
			m := item.(map[string]interface{})
			tagInfo := monitor.PolicyTag{
				Key:   helper.String(m["key"].(string)),
				Value: helper.String(m["value"].(string)),
			}
			tagSet = append(tagSet, &tagInfo)
		}
		request.Tag = tagSet[0]

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err := monitorService.client.UseMonitorClient().BindingPolicyTag(request)
			if err != nil {
				log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
					logId, request.GetAction(), request.ToJsonString(), err.Error())
				return retryError(err, InternalError)
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
			return nil
		}); err != nil {
			return err
		}
	}

	return resourceTencentMonitorAlarmPolicyRead(d, meta)
}

func resourceTencentMonitorAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy.read")()
	defer inconsistentCheck(d, meta)()

	//logId := getLogId(contextNil)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDescribeAlarmPolicyRequest()
		policy         *monitor.AlarmPolicy
	)

	policyId := d.Id()
	request.PolicyId = &policyId
	request.Module = helper.String("monitor")

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := monitorService.client.UseMonitorClient().DescribeAlarmPolicy(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		policy = response.Response.Policy
		return nil
	}); err != nil {
		return err
	}

	if policy == nil {
		d.SetId("")
		return nil
	}

	var errs []error
	errs = append(errs,
		d.Set("policy_name", policy.PolicyName),
		d.Set("monitor_type", policy.MonitorType),
		d.Set("namespace", policy.Namespace),
		d.Set("remark", policy.Remark),
		d.Set("enable", policy.Enable),
		d.Set("project_id", policy.ProjectId),
	)

	if policy.ConditionTemplateId != nil && *policy.ConditionTemplateId != "" {
		id, err := strconv.ParseInt(*policy.ConditionTemplateId, 10, 64)
		if id != 0 && err == nil {
			errs = append(errs, d.Set("conditon_template_id", id))
		}
	}

	if policy.InsertTime != nil {
		t := time.Unix(*policy.InsertTime, 0)
		tFmt := t.Format("2006-01-02 15:04:05")
		errs = append(errs, d.Set("create_time", tFmt))
	}
	if policy.UpdateTime != nil {
		t := time.Unix(*policy.UpdateTime, 0)
		tFmt := t.Format("2006-01-02 15:04:05")
		errs = append(errs, d.Set("update_time", tFmt))
	}

	var rules = make([]interface{}, 0, 100)
	for _, rule := range policy.Condition.Rules {

		m := map[string]interface{}{
			"metric_name":      rule.MetricName,
			"period":           rule.Period,
			"operator":         rule.Operator,
			"value":            rule.Value,
			"continue_period":  rule.ContinuePeriod,
			"notice_frequency": rule.NoticeFrequency,
			"description":      rule.Description,
			"unit":             rule.Unit,
			"rule_type":        rule.RuleType,
		}
		if rule.Filter != nil {
			if *rule.Filter.Type != "" || *rule.Filter.Dimensions != "" {
				var filter = make([]interface{}, 0, 10)
				alarmPolicyFilter := map[string]interface{}{
					"type":       rule.Filter.Type,
					"dimensions": rule.Filter.Dimensions,
				}
				filter = append(filter, alarmPolicyFilter)
				if len(filter) > 0 {
					m["filter"] = filter
				}
			}
		}

		rules = append(rules, m)
	}

	conditions := map[string]interface{}{
		"is_union_rule": policy.Condition.IsUnionRule,
		"rules":         rules,
	}
	_ = d.Set("conditions", []interface{}{conditions})

	eventConditions := make([]map[string]interface{}, 0, len(policy.EventCondition.Rules))
	for _, eventRule := range policy.EventCondition.Rules {

		m := make(map[string]interface{}, 5)
		m["metric_name"] = eventRule.MetricName
		m["period"] = eventRule.Period
		m["operator"] = eventRule.Operator
		m["value"] = eventRule.Value
		m["continue_period"] = eventRule.ContinuePeriod
		m["notice_frequency"] = eventRule.NoticeFrequency
		m["is_power_notice"] = eventRule.IsPowerNotice
		m["notice_frequency"] = eventRule.NoticeFrequency
		m["description"] = eventRule.Description
		m["unit"] = eventRule.Unit
		m["rule_type"] = eventRule.RuleType
		if eventRule.Filter != nil {
			if *eventRule.Filter.Type != "" || *eventRule.Filter.Dimensions != "" {
				var filter = make([]interface{}, 0, 10)
				alarmPolicyFilter := map[string]interface{}{
					"type":       eventRule.Filter.Type,
					"dimensions": eventRule.Filter.Dimensions,
				}
				filter = append(filter, alarmPolicyFilter)
				if len(filter) > 0 {
					m["filter"] = filter
				}
			}
		}
		eventConditions = append(eventConditions, m)
	}
	_ = d.Set("event_conditions", eventConditions)
	var noticeIds = make([]interface{}, 0, 100)
	for _, notice := range policy.NoticeIds {
		noticeIds = append(noticeIds, notice)
	}
	errs = append(errs, d.Set("notice_ids", noticeIds))

	var triggerTasks = make([]interface{}, 0, 100)
	for _, task := range policy.TriggerTasks {
		m := map[string]interface{}{}
		m["type"] = task.Type
		m["task_config"] = task.TaskConfig
		triggerTasks = append(triggerTasks, m)
	}
	errs = append(errs, d.Set("trigger_tasks", triggerTasks))

	tagSets := make([]map[string]interface{}, 0, len(policy.TagInstances))
	for _, item := range policy.TagInstances {
		tagSets = append(tagSets, map[string]interface{}{
			"key":   item.Key,
			"value": item.Value,
		})
	}
	_ = d.Set("policy_tag", tagSets)

	var errResults *multierror.Error
	for i := range errs {
		err := errs[i]
		if err != nil {
			errResults = multierror.Append(errResults, err)
		}
	}
	return errResults.ErrorOrNil()
}

func resourceTencentMonitorAlarmPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy.update")()
	//logId := getLogId(contextNil)
	//ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	)

	if d.HasChange("policy_name") {
		request := monitor.NewModifyAlarmPolicyInfoRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		request.Key = helper.String("NAME")
		value := d.Get("policy_name").(string)
		request.Value = helper.String(value)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyInfo(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("remark") {
		request := monitor.NewModifyAlarmPolicyInfoRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		request.Key = helper.String("REMARK")
		value := d.Get("remark").(string)
		request.Value = helper.String(value)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyInfo(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("enable") {
		request := monitor.NewModifyAlarmPolicyStatusRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		enable := d.Get("enable").(int)
		request.Enable = helper.IntInt64(enable)

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyStatus(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("conditions") || d.HasChange("event_conditions") {
		request := monitor.NewModifyAlarmPolicyConditionRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		if v, ok := d.GetOk("conditions"); ok {
			conditions := v.([]interface{})

			if len(conditions) != 1 {
				return fmt.Errorf("need only one conditions.")
			}

			condition := conditions[0].(map[string]interface{})
			var policy = monitor.AlarmPolicyCondition{}
			policy.IsUnionRule = helper.IntInt64(condition["is_union_rule"].(int))
			policy.Rules = make([]*monitor.AlarmPolicyRule, 0, 10)

			rules := condition["rules"]

			for _, item := range rules.([]interface{}) {
				m := item.(map[string]interface{})
				alarmPolicyRule := monitor.AlarmPolicyRule{}
				if m["metric_name"] != nil {
					alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
				}
				if m["period"] != nil {
					alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
				}
				if m["value"] != nil {
					alarmPolicyRule.Value = helper.String(m["value"].(string))
				}
				if m["operator"] != nil {
					alarmPolicyRule.Operator = helper.String(m["operator"].(string))
				}
				if m["continue_period"] != nil {
					alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
				}
				if m["notice_frequency"] != nil {
					alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
				}
				if m["is_power_notice"] != nil {
					alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
				}
				if m["filter"] != nil {
					filters := m["filter"].([]interface{})
					// Max Items is 1
					if len(filters) > 0 {
						filter := filters[0].(map[string]interface{})
						alarmPolicyFilter := monitor.AlarmPolicyFilter{
							Type:       helper.String(filter["type"].(string)),
							Dimensions: helper.String(filter["dimensions"].(string)),
						}
						alarmPolicyRule.Filter = &alarmPolicyFilter
					}
				}
				if m["description"] != nil {
					alarmPolicyRule.Description = helper.String(m["description"].(string))
				}
				if m["unit"] != nil {
					alarmPolicyRule.Unit = helper.String(m["unit"].(string))
				}
				if m["rule_type"] != nil {
					alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
				}
				policy.Rules = append(policy.Rules, &alarmPolicyRule)
			}
			request.Condition = &policy
		}
		if v, ok := d.GetOk("event_conditions"); ok {
			eventCondition := monitor.AlarmPolicyEventCondition{}
			rules := make([]*monitor.AlarmPolicyRule, 0, 10)
			for _, item := range v.([]interface{}) {
				m := item.(map[string]interface{})
				alarmPolicyRule := monitor.AlarmPolicyRule{}
				if m["metric_name"] != nil {
					alarmPolicyRule.MetricName = helper.String(m["metric_name"].(string))
				}
				if m["period"] != nil {
					alarmPolicyRule.Period = helper.IntInt64(m["period"].(int))
				}
				if m["operator"] != nil {
					alarmPolicyRule.Operator = helper.String(m["operator"].(string))
				}
				if m["value"] != nil {
					alarmPolicyRule.Value = helper.String(m["value"].(string))
				}
				if m["continue_period"] != nil {
					alarmPolicyRule.ContinuePeriod = helper.IntInt64(m["continue_period"].(int))
				}
				if m["notice_frequency"] != nil {
					alarmPolicyRule.NoticeFrequency = helper.IntInt64(m["notice_frequency"].(int))
				}
				if m["is_power_notice"] != nil {
					alarmPolicyRule.IsPowerNotice = helper.IntInt64(m["is_power_notice"].(int))
				}
				if m["filter"] != nil {
					filters := m["filter"].([]interface{})
					// Max Items is 1
					if len(filters) > 0 {
						filter := filters[0].(map[string]interface{})
						alarmPolicyFilter := monitor.AlarmPolicyFilter{
							Type:       helper.String(filter["type"].(string)),
							Dimensions: helper.String(filter["dimensions"].(string)),
						}
						alarmPolicyRule.Filter = &alarmPolicyFilter
					}
				}
				if m["description"] != nil {
					alarmPolicyRule.Description = helper.String(m["description"].(string))
				}
				if m["unit"] != nil {
					alarmPolicyRule.Unit = helper.String(m["unit"].(string))
				}
				if m["rule_type"] != nil {
					alarmPolicyRule.RuleType = helper.String(m["rule_type"].(string))
				}
				rules = append(rules, &alarmPolicyRule)
			}
			eventCondition.Rules = rules
			request.EventCondition = &eventCondition
		}

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyCondition(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("notice_ids") {
		request := monitor.NewModifyAlarmPolicyNoticeRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())

		if v, ok := d.GetOk("notice_ids"); ok {
			notice := make([]*string, 0, 10)
			for _, item := range v.([]interface{}) {
				notice = append(notice, helper.String(item.(string)))
			}
			request.NoticeIds = notice
		}

		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyNotice(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	if d.HasChange("trigger_tasks") {
		request := monitor.NewModifyAlarmPolicyTasksRequest()
		request.Module = helper.String("monitor")
		request.PolicyId = helper.String(d.Id())
		if v, ok := d.GetOk("trigger_tasks"); ok {
			tasks := make([]*monitor.AlarmPolicyTriggerTask, 0, 10)
			for _, item := range v.([]interface{}) {
				m := item.(map[string]interface{})
				triggerTask := monitor.AlarmPolicyTriggerTask{}
				triggerTask.Type = helper.String(m["type"].(string))
				triggerTask.TaskConfig = helper.String(m["task_config"].(string))
				tasks = append(tasks, &triggerTask)
			}
			request.TriggerTasks = tasks
		}
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			if _, err := monitorService.client.UseMonitorClient().ModifyAlarmPolicyTasks(request); err != nil {
				return retryError(err, InternalError)
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return resourceTencentMonitorAlarmPolicyRead(d, meta)
}

func resourceTencentMonitorAlarmPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_alarm_policy.delete")()

	var (
		monitorService = MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
		request        = monitor.NewDeleteAlarmPolicyRequest()
	)
	request.Module = helper.String("monitor")
	policyIds := []*string{helper.String(d.Id())}
	request.PolicyIds = policyIds

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		if _, err := monitorService.client.UseMonitorClient().DeleteAlarmPolicy(request); err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
