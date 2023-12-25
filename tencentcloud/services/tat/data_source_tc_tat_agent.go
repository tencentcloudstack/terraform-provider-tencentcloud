package tat

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTatAgent() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTatAgentRead,
		Schema: map[string]*schema.Schema{
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of instance IDs for the query.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions. agent-status - String - Required: No - (Filter condition) Filter by agent status. Valid values: Online, Offline. environment - String - Required: No - (Filter condition) Filter by the agent environment. Valid value: Linux. instance-id - String - Required: No - (Filter condition) Filter by the instance ID. Up to 10 Filters allowed in one request. For each filter, five Filter.Values can be specified. InstanceIds and Filters cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Field to be filtered.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter values of the field.",
						},
					},
				},
			},

			"automation_agent_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of agent message.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "InstanceId.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent version.",
						},
						"last_heartbeat_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of last heartbeat.",
						},
						"agent_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Agent status.Ranges:&lt;li&gt; Online:Online&lt;li&gt; Offline:Offline.",
						},
						"environment": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Environment for Agent.Ranges:&lt;li&gt; Linux:Linux instance&lt;li&gt; Windows:Windows instance.",
						},
						"support_features": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of feature Agent support.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTatAgentRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tat_agent.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tat.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tat.Filter{}
			filterMap := item.(map[string]interface{})

			if v, ok := filterMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}
			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var automationAgentSet []*tat.AutomationAgentInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTatAgentByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		automationAgentSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(automationAgentSet))
	tmpList := make([]map[string]interface{}, 0, len(automationAgentSet))

	if automationAgentSet != nil {
		for _, automationAgentInfo := range automationAgentSet {
			automationAgentInfoMap := map[string]interface{}{}

			if automationAgentInfo.InstanceId != nil {
				automationAgentInfoMap["instance_id"] = automationAgentInfo.InstanceId
			}

			if automationAgentInfo.Version != nil {
				automationAgentInfoMap["version"] = automationAgentInfo.Version
			}

			if automationAgentInfo.LastHeartbeatTime != nil {
				automationAgentInfoMap["last_heartbeat_time"] = automationAgentInfo.LastHeartbeatTime
			}

			if automationAgentInfo.AgentStatus != nil {
				automationAgentInfoMap["agent_status"] = automationAgentInfo.AgentStatus
			}

			if automationAgentInfo.Environment != nil {
				automationAgentInfoMap["environment"] = automationAgentInfo.Environment
			}

			if automationAgentInfo.SupportFeatures != nil {
				automationAgentInfoMap["support_features"] = automationAgentInfo.SupportFeatures
			}

			ids = append(ids, *automationAgentInfo.InstanceId)
			tmpList = append(tmpList, automationAgentInfoMap)
		}

		_ = d.Set("automation_agent_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
