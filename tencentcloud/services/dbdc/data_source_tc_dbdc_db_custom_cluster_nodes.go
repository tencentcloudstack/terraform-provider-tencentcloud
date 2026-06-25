package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomClusterNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomClusterNodesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DB Custom cluster ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported filter names: node-name (DB Custom node name).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter field values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"node_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom cluster node list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node internal IP address.",
						},
						"ssh_endpoint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node SSH access endpoint. Format: IP:Port.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node instance status in the cluster.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node region.",
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node type.",
						},
					},
				},
			},

			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of nodes in the cluster.",
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomClusterNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_cluster_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*dbdcv20201029.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := dbdcv20201029.Filter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.([]interface{})
				for i := range valuesSet {
					value := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(value))
				}
			}
			tmpSet = append(tmpSet, &filter)
		}
		paramMap["Filters"] = tmpSet
	}

	var respData []*dbdcv20201029.DBCustomClusterNode
	var totalCount int64
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, count, e := service.DescribeDBCustomClusterNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		totalCount = count
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	nodeSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, node := range respData {
			nodeMap := map[string]interface{}{}
			if node.NodeId != nil {
				nodeMap["node_id"] = node.NodeId
			}

			if node.NodeName != nil {
				nodeMap["node_name"] = node.NodeName
			}

			if node.LanIP != nil {
				nodeMap["lan_ip"] = node.LanIP
			}

			if node.SSHEndpoint != nil {
				nodeMap["ssh_endpoint"] = node.SSHEndpoint
			}

			if node.Status != nil {
				nodeMap["status"] = node.Status
			}

			if node.Zone != nil {
				nodeMap["zone"] = node.Zone
			}

			if node.NodeType != nil {
				nodeMap["node_type"] = node.NodeType
			}

			nodeSetList = append(nodeSetList, nodeMap)
		}

		_ = d.Set("node_set", nodeSetList)
	}

	_ = d.Set("total_count", totalCount)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
