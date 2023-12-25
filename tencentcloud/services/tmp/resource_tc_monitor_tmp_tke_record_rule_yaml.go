package tmp

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMonitorTmpTkeRecordRuleYaml() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateYaml,
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewCreatePrometheusRecordRuleYamlRequest()

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
	}

	tmpRecordRuleName := ""
	if v, ok := d.GetOk("content"); ok {
		m, _ := tccommon.YamlParser(v.(string))
		metadata := m["metadata"]
		if metadata != nil {
			if metadata.(map[interface{}]interface{})["name"] != nil {
				tmpRecordRuleName = metadata.(map[interface{}]interface{})["name"].(string)
			}
		}
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().CreatePrometheusRecordRuleYaml(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	d.SetId(strings.Join([]string{instanceId, tmpRecordRuleName}, tccommon.FILED_SP))
	return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
}

func resourceTencentCloudTkeTmpRecordRuleYamlRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	instanceId := ids[0]
	name := ids[1]

	recordRuleService := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	request, err := recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			request, err = recordRuleService.DescribePrometheusRecordRuleByName(ctx, instanceId, name)
			if err != nil {
				return tccommon.RetryError(err)
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
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := monitor.NewModifyPrometheusRecordRuleYamlRequest()

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	request.InstanceId = &ids[0]
	request.Name = &ids[1]

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))

			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMonitorClient().ModifyPrometheusRecordRuleYaml(request)
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

			return resourceTencentCloudTkeTmpRecordRuleYamlRead(d, meta)
		}
	}

	return nil
}

func resourceTencentCloudTkeTmpRecordRuleYamlDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_monitor_tmp_tke_record_rule_yaml.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	ids := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(ids) != 2 {
		return fmt.Errorf("id is broken, id is %s", d.Id())
	}

	service := svcmonitor.NewMonitorService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	if err := service.DeletePrometheusRecordRuleYaml(ctx, ids[0], ids[1]); err != nil {
		return err
	}

	return nil
}
