package tse

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTseGatewayNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayNodesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "gateway ID.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "gateway group ID.",
			},

			"node_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "nodes information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "gateway node id.",
						},
						"node_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "gateway node ip.",
						},
						"zone_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone idNote: This field may return null, indicating that a valid value is not available.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ZoneNote: This field may return null, indicating that a valid value is not available.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group IDNote: This field may return null, indicating that a valid value is not available.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group nameNote: This field may return null, indicating that a valid value is not available.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "statusNote: This field may return null, indicating that a valid value is not available.",
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

func dataSourceTencentCloudTseGatewayNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tse_gateway_nodes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var result []*tse.CloudNativeAPIGatewayNode
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayNodesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	nodeListList := []interface{}{}
	if result != nil {
		for _, nodeList := range result {
			nodeListMap := map[string]interface{}{}

			if nodeList.NodeId != nil {
				nodeListMap["node_id"] = nodeList.NodeId
			}

			if nodeList.NodeIp != nil {
				nodeListMap["node_ip"] = nodeList.NodeIp
			}

			if nodeList.ZoneId != nil {
				nodeListMap["zone_id"] = nodeList.ZoneId
			}

			if nodeList.Zone != nil {
				nodeListMap["zone"] = nodeList.Zone
			}

			if nodeList.GroupId != nil {
				nodeListMap["group_id"] = nodeList.GroupId
			}

			if nodeList.GroupName != nil {
				nodeListMap["group_name"] = nodeList.GroupName
			}

			if nodeList.Status != nil {
				nodeListMap["status"] = nodeList.Status
			}

			nodeListList = append(nodeListList, nodeListMap)
			ids = append(ids, *nodeList.NodeId)
		}

		_ = d.Set("node_list", nodeListList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), nodeListList); e != nil {
			return e
		}
	}
	return nil
}
