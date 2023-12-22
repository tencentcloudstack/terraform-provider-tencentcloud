package trabbit

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctdmq "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmq"

	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tdmq "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmq/v20200217"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTdmqRabbitmqNodeList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTdmqRabbitmqNodeListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "rabbitmq cluster ID.",
			},
			"node_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search node name.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "filter parameter name and valueNow there is only one nodeStatusrunning/downArray type, compatible with adding filter parameters later.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the filter parameter.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "value.",
						},
					},
				},
			},
			"sort_element": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by the specified element, now there are only 2cpuUsage/diskUsage.",
			},
			"sort_order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ascend/descend.",
			},
			// computed
			"node_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "cluster listNote: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "node nameNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"node_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "node statusNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"cpu_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CPU usageNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory usage, in MBNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"disk_usage": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "disk usageNote: This field may return null, indicating that no valid value can be obtained.",
						},
						"process_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of Erlang processes for RabbitmqNote: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudTdmqRabbitmqNodeListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tdmq_rabbitmq_node_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = svctdmq.NewTdmqService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		nodeList   []*tdmq.RabbitMQPrivateNode
		instanceId string
		nodeName   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("node_name"); ok {
		paramMap["NodeName"] = helper.String(v.(string))
		nodeName = v.(string)
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tdmq.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := tdmq.Filter{}
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

	if v, ok := d.GetOk("sort_element"); ok {
		paramMap["SortElement"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort_order"); ok {
		paramMap["SortOrder"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTdmqRabbitmqNodeListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		nodeList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0)
	tmpList := make([]map[string]interface{}, 0, len(nodeList))
	if nodeList != nil {
		for _, rabbitMQPrivateNode := range nodeList {
			rabbitMQPrivateNodeMap := map[string]interface{}{}

			if rabbitMQPrivateNode.NodeName != nil {
				rabbitMQPrivateNodeMap["node_name"] = rabbitMQPrivateNode.NodeName
			}

			if rabbitMQPrivateNode.NodeStatus != nil {
				rabbitMQPrivateNodeMap["node_status"] = rabbitMQPrivateNode.NodeStatus
			}

			if rabbitMQPrivateNode.CPUUsage != nil {
				rabbitMQPrivateNodeMap["cpu_usage"] = rabbitMQPrivateNode.CPUUsage
			}

			if rabbitMQPrivateNode.Memory != nil {
				rabbitMQPrivateNodeMap["memory"] = rabbitMQPrivateNode.Memory
			}

			if rabbitMQPrivateNode.DiskUsage != nil {
				rabbitMQPrivateNodeMap["disk_usage"] = rabbitMQPrivateNode.DiskUsage
			}

			if rabbitMQPrivateNode.ProcessNumber != nil {
				rabbitMQPrivateNodeMap["process_number"] = rabbitMQPrivateNode.ProcessNumber
			}

			tmpList = append(tmpList, rabbitMQPrivateNodeMap)
		}

		_ = d.Set("node_list", tmpList)
	}

	ids = append(ids, instanceId)
	ids = append(ids, nodeName)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
