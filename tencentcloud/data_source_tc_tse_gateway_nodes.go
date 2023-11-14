/*
Use this data source to query detailed information of tse gateway_nodes

Example Usage

```hcl
data "tencentcloud_tse_gateway_nodes" "gateway_nodes" {
  gateway_id = "gateway-xx"
  group_id = "group-xx"
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayNodes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayNodesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"group_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Gateway group ID.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Gateway nodes information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The Number of nodes information.",
						},
						"node_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Nodes information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"node_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway node id.",
									},
									"node_ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Gateway node ip.",
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
										Description: "StatusNote: This field may return null, indicating that a valid value is not available.",
									},
								},
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
		},
	}
}

func dataSourceTencentCloudTseGatewayNodesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_nodes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.DescribeCloudNativeAPIGatewayNodesResult

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGatewayNodesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		describeCloudNativeAPIGatewayNodesResultMap := map[string]interface{}{}

		if result.TotalCount != nil {
			describeCloudNativeAPIGatewayNodesResultMap["total_count"] = result.TotalCount
		}

		if result.NodeList != nil {
			nodeListList := []interface{}{}
			for _, nodeList := range result.NodeList {
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
			}

			describeCloudNativeAPIGatewayNodesResultMap["node_list"] = []interface{}{nodeListList}
		}

		ids = append(ids, *result.GatewayId)
		_ = d.Set("result", describeCloudNativeAPIGatewayNodesResultMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), describeCloudNativeAPIGatewayNodesResultMap); e != nil {
			return e
		}
	}
	return nil
}
