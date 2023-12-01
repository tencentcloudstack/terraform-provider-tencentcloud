/*
Provides a resource to create a monitor tmpAlertRule

Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}


resource "tencentcloud_monitor_tmp_instance" "foo" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_monitor_tmp_cvm_agent" "foo" {
  instance_id = tencentcloud_monitor_tmp_instance.foo.id
  name        = "tf-agent"
}

resource "tencentcloud_monitor_tmp_alert_rule" "foo" {
    duration    = "2m"
    expr        = "avg by (instance) (mysql_global_status_threads_connected) / avg by (instance) (mysql_global_variables_max_connections)  > 0.8"
    instance_id = tencentcloud_monitor_tmp_instance.foo.id
    receivers   = ["notice-f2svbu3w"]
    rule_name   = "MySQL 连接数过多"
    rule_state  = 2
    type        = "MySQL/MySQL 连接数过多"

    annotations {
        key   = "description"
        value = "MySQL 连接数过多, 实例: {{$labels.instance}}，当前值: {{ $value | humanizePercentage }}。"
    }
    annotations {
        key   = "summary"
        value = "MySQL 连接数过多(>80%)"
    }

    labels {
        key   = "severity"
        value = "warning"
    }
}

```
Import

monitor tmpAlertRule can be imported using the id, e.g.
```
$ terraform import tencentcloud_monitor_tmp_alert_rule.tmpAlertRule instanceId#Rule_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpAlertRule() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudMonitorTmpAlertRuleRead,
		Create: resourceTencentCloudMonitorTmpAlertRuleCreate,
		Update: resourceTencentCloudMonitorTmpAlertRuleUpdate,
		Delete: resourceTencentCloudMonitorTmpAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance id.",
			},
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule name.",
			},
			"expr": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Rule expression, reference documentation: `https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/`.",
			},
			"receivers": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Alarm notification template id list.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"rule_state": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Rule state code.",
			},
			"duration": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Rule alarm duration.",
			},
			"labels": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rule alarm duration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
					},
				},
			},
			"annotations": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Rule alarm duration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "value.",
						},
					},
				},
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Alarm Policy Template Classification.",
			},
		},
	}
}

func resourceTencentCloudMonitorTmpAlertRuleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_alert_rule.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = monitor.NewCreateAlertRuleRequest()
		response *monitor.CreateAlertRuleResponse
	)

	var instanceId string

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(instanceId)
	}
	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}
	if v, ok := d.GetOk("expr"); ok {
		request.Expr = helper.String(v.(string))
	}
	if v, ok := d.GetOk("receivers"); ok {
		receivers := v.(*schema.Set).List()
		receiverArr := make([]*string, 0, len(receivers))
		for _, receiver := range receivers {
			receiverArr = append(receiverArr, helper.String(receiver.(string)))
		}
		request.Receivers = receiverArr
	}
	if v, ok := d.GetOk("rule_state"); ok {
		request.RuleState = helper.IntInt64(v.(int))
	}
	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.String(v.(string))
	}
	if v, ok := d.GetOk("labels"); ok {
		labelsList := v.([]interface{})
		prometheusRuleKV := make([]*monitor.PrometheusRuleKV, 0, len(labelsList))
		for _, labels := range labelsList {
			if labels == nil {
				return fmt.Errorf("Invalid `labels` parameter, must not be empty")
			}
			label := labels.(map[string]interface{})
			var kv monitor.PrometheusRuleKV
			kv.Key = helper.String(label["key"].(string))
			kv.Value = helper.String(label["value"].(string))
			prometheusRuleKV = append(prometheusRuleKV, &kv)
		}
		request.Labels = prometheusRuleKV
	}
	if v, ok := d.GetOk("annotations"); ok {
		annotationsList := v.([]interface{})
		prometheusRuleKV := make([]*monitor.PrometheusRuleKV, 0, len(annotationsList))
		for _, annotations := range annotationsList {
			if annotations == nil {
				return fmt.Errorf("Invalid `annotation` parameter, must not be empty")
			}
			annotation := annotations.(map[string]interface{})
			var kv monitor.PrometheusRuleKV
			kv.Key = helper.String(annotation["key"].(string))
			kv.Value = helper.String(annotation["value"].(string))
			prometheusRuleKV = append(prometheusRuleKV, &kv)
		}
		request.Annotations = prometheusRuleKV
	}
	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreateAlertRule(request)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains("FailedOperation", ee.Code) {
				return resource.NonRetryableError(ee)
			}
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create monitor tmpAlertRule failed, reason:%+v", logId, err)
		return err
	}

	ruleId := *response.Response.RuleId

	d.SetId(strings.Join([]string{instanceId, ruleId}, FILED_SP))

	return resourceTencentCloudMonitorTmpAlertRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_alert_rule.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	tmpAlertRule, err := service.DescribeMonitorTmpAlertRuleById(ctx, ids[0], ids[1])

	if err != nil {
		return err
	}

	if tmpAlertRule == nil {
		d.SetId("")
		return fmt.Errorf("resource `tmpAlertRule` %s does not exist", ids[1])
	}

	_ = d.Set("instance_id", ids[0])
	if tmpAlertRule.RuleName != nil {
		_ = d.Set("rule_name", tmpAlertRule.RuleName)
	}
	if tmpAlertRule.Expr != nil {
		_ = d.Set("expr", tmpAlertRule.Expr)
	}
	if tmpAlertRule.Receivers != nil {
		list := tmpAlertRule.Receivers
		result := make([]string, 0, len(list))
		for _, v := range list {
			result = append(result, *v)
		}
		_ = d.Set("receivers", result)
	}
	if tmpAlertRule.RuleState != nil {
		_ = d.Set("rule_state", tmpAlertRule.RuleState)
	}
	if tmpAlertRule.Duration != nil {
		_ = d.Set("duration", tmpAlertRule.Duration)
	}
	if tmpAlertRule.Labels != nil {
		labelsList := tmpAlertRule.Labels
		result := make([]map[string]interface{}, 0, len(labelsList))
		for _, v := range labelsList {
			mapping := map[string]interface{}{
				"key":   v.Key,
				"value": v.Value,
			}
			result = append(result, mapping)
		}
		_ = d.Set("labels", result)
	}
	if tmpAlertRule.Annotations != nil {
		annotationsList := tmpAlertRule.Annotations
		result := make([]map[string]interface{}, 0, len(annotationsList))
		for _, v := range annotationsList {
			mapping := map[string]interface{}{
				"key":   v.Key,
				"value": v.Value,
			}
			result = append(result, mapping)
		}
		_ = d.Set("annotations", result)
	}
	if tmpAlertRule.Type != nil {
		_ = d.Set("type", tmpAlertRule.Type)
	}

	return nil
}

func resourceTencentCloudMonitorTmpAlertRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_alert_rule.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewUpdateAlertRuleRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = helper.String(ids[0])
	request.RuleId = helper.String(ids[1])

	if v, ok := d.GetOk("rule_name"); ok {
		request.RuleName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expr"); ok {
		request.Expr = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receivers"); ok {
		receivers := v.(*schema.Set).List()
		receiverArr := make([]*string, 0, len(receivers))
		for _, receiver := range receivers {
			receiverArr = append(receiverArr, helper.String(receiver.(string)))
		}
		request.Receivers = receiverArr
	}

	if v, ok := d.GetOk("rule_state"); ok {
		request.RuleState = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("duration"); ok {
		request.Duration = helper.String(v.(string))
	}

	if v, ok := d.GetOk("labels"); ok {
		labelsList := v.([]interface{})
		prometheusRuleKV := make([]*monitor.PrometheusRuleKV, 0, len(labelsList))
		for _, labels := range labelsList {
			label := labels.(map[string]interface{})
			var kv monitor.PrometheusRuleKV
			kv.Key = helper.String(label["key"].(string))
			kv.Value = helper.String(label["value"].(string))
			prometheusRuleKV = append(prometheusRuleKV, &kv)
		}
		request.Labels = prometheusRuleKV
	}

	if v, ok := d.GetOk("annotations"); ok {
		annotationsList := v.([]interface{})
		prometheusRuleKV := make([]*monitor.PrometheusRuleKV, 0, len(annotationsList))
		for _, annotations := range annotationsList {
			annotation := annotations.(map[string]interface{})
			var kv monitor.PrometheusRuleKV
			kv.Key = helper.String(annotation["key"].(string))
			kv.Value = helper.String(annotation["value"].(string))
			prometheusRuleKV = append(prometheusRuleKV, &kv)
		}
		request.Annotations = prometheusRuleKV
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().UpdateAlertRule(request)
		if e != nil {
			ee, ok := e.(*sdkErrors.TencentCloudSDKError)
			if ok && IsContains("FailedOperation", ee.Code) {
				return resource.NonRetryableError(ee)
			}
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

	return resourceTencentCloudMonitorTmpAlertRuleRead(d, meta)
}

func resourceTencentCloudMonitorTmpAlertRuleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_alert_rule.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	if err := service.DeleteMonitorTmpAlertRule(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
