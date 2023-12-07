package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMonitorTmpTkeRecordRuleYaml() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTkeTmpRecordRuleYamlRead,
		Create: resourceTencentCloudTkeTmpRecordRuleYamlCreate,
		Update: resourceTencentCloudTkeTmpRecordRuleYamlUpdate,
		Delete: resourceTencentCloudTkeTmpRecordRuleYamlDelete,
		//Importer: &schema.ResourceImporter{
		//	State: schema.ImportStatePassthrough,
		//},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instance Id.",
			},

			"content": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateYaml,
				Description:  "Contents of record rules in yaml format.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the instance.",
			},

			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last modified time of record rule.",
			},

			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Used for the argument, if the configuration comes to the template, the template id.",
			},

			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "An ID identify the cluster, like cls-xxxxxx.",
			},
		},
	}
}

func resourceTencentCloudTkeTmpRecordRuleYamlCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewCreatePrometheusRecordRuleYamlRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	tmpRecordRuleName := ""
	if v, ok := d.GetOk("content"); ok {
		m, _ := YamlParser(v.(string))
		metadata := m["metadata"]
		if metadata != nil {
			if metadata.(map[interface{}]interface{})["name"] != nil {
				tmpRecordRuleName = metadata.(map[interface{}]interface{})["name"].(string)
			}
		}
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().CreatePrometheusRecordRuleYaml(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tke tmpRecordRule failed, reason:%+v", logId, err)
		return err
	}

	instanceId := *request.InstanceId
	d.SetId(strings.Join([]string{instanceId, tmpRecordRuleName}, FILED_SP))
	return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
}

func resourceTencentCloudTkeTmpRecordRuleYamlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	name := ids[1]

	recordRuleService := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	request, err := recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			request, err = recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	recordRules := request.Response.Records
	if len(recordRules) == 0 {
		d.SetId("")
		return nil
	}

	recordRule := recordRules[0]
	if recordRule == nil {
		return nil
	}
	_ = d.Set("instance_id", instanceId)
	_ = d.Set("name", recordRule.Name)
	_ = d.Set("update_time", recordRule.UpdateTime)
	_ = d.Set("template_id", recordRule.TemplateId)
	//_ = d.Set("content", recordRule.Content)
	_ = d.Set("cluster_id", recordRule.ClusterId)

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := monitor.NewModifyPrometheusRecordRuleYamlRequest()

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = &ids[0]
	request.Name = &ids[1]

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))

			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseMonitorClient().ModifyPrometheusRecordRuleYaml(request)
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

			return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
		}
	}

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ids := strings.Split(d.Id(), FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}
	if err := service.DeletePrometheusRecordRuleYaml(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
