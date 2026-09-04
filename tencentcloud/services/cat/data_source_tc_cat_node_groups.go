package cat

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCatNodeGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCatNodeGroupsRead,
		Schema: map[string]*schema.Schema{
			"node_type": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Node type. 0: all, 1: IDC, 2: LastMile, 3: Mobile. Defaults to 0 if not set.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"task_category": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Node category. 0: all, 1: PC, 2: Mobile. Defaults to 0 if not set.",
			},

			"ip_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "IP type. 0: all, 1: IPv4, 2: IPv6. Defaults to 0 if not set.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Probe node description keyword.",
			},

			"region_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Region ID. 0: selected probe points, 1: Chinese Mainland, 2: Hong Kong, Macao and Taiwan, 3: Asia Pacific, 4: Europe and America, 5: Africa and Oceania. Defaults to 0 if not set.",
			},

			"district_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Province or country ID. 0 means all. Defaults to 0 if not set.",
			},

			"net_service_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ISP ID. 0: all, 1: China Telecom, 2: China Unicom, 3: China Mobile, 99: others. Defaults to 0 if not set.",
			},

			"node_group_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Node group type. 0: advanced probe group, 1: availability node, 2: my probe group. Defaults to 0 if not set.",
			},

			"task_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Task type, such as 1, 2, 3, 4, 5, 6, 7. 1-page performance, 2-file upload, 3-file download, 4-port performance, 5-network quality, 6-audio and video experience, 7-domain whois. Defaults to 0 if not set, no filtering by task type.",
			},

			"probe_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Test type, including scheduled test and instant test. 0-scheduled probe, others mean instant probe.",
			},

			"node_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Tree node list, two levels in total.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node name.",
						},
						"children": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Child nodes.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Child node ID.",
									},
									"content": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Child node name.",
									},
									"children": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Node list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Node code.",
												},
												"content": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Node name.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"district_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Province or country list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Province (international) or ISP ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name.",
						},
					},
				},
			},

			"net_service_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "ISP list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Province (international) or ISP ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name.",
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

func dataSourceTencentCloudCatNodeGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cat_node_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("node_type"); ok {
		nodeTypeList := v.([]interface{})
		nodeTypePtrList := make([]*int64, 0, len(nodeTypeList))
		for _, item := range nodeTypeList {
			nodeTypePtrList = append(nodeTypePtrList, helper.IntInt64(item.(int)))
		}
		paramMap["node_type"] = nodeTypePtrList
	}

	if v, _ := d.GetOk("task_category"); v != nil {
		paramMap["task_category"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("ip_type"); v != nil {
		paramMap["ip_type"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["name"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("region_id"); v != nil {
		paramMap["region_id"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("district_id"); v != nil {
		paramMap["district_id"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("net_service_id"); v != nil {
		paramMap["net_service_id"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("node_group_type"); v != nil {
		paramMap["node_group_type"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("task_type"); v != nil {
		paramMap["task_type"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("probe_type"); v != nil {
		paramMap["probe_type"] = helper.IntUint64(v.(int))
	}

	catService := CatService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var nodeList []*cat.NodeTree
	var districtList []*cat.DistinctOrNetServiceInfo
	var netServiceList []*cat.DistinctOrNetServiceInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		resultNodeList, resultDistrictList, resultNetServiceList, e := catService.DescribeCatNodeGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if resultNodeList == nil && resultDistrictList == nil && resultNetServiceList == nil {
			e = logError(logId, "DescribeNodeGroups response is nil")
			return resource.NonRetryableError(e)
		}

		nodeList = resultNodeList
		districtList = resultDistrictList
		netServiceList = resultNetServiceList
		return nil
	})
	if err != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		log.Printf("[CRITAL]%s api[%s] fail, reason[%+v]", logId, "DescribeNodeGroups", err)
		return err
	}

	ids := make([]string, 0)
	nodeListTmp := make([]map[string]interface{}, 0)
	if nodeList != nil {
		for _, nodeTree := range nodeList {
			nodeTreeMap := map[string]interface{}{}
			if nodeTree.ID != nil {
				nodeTreeMap["id"] = nodeTree.ID
				ids = append(ids, *nodeTree.ID)
			}
			if nodeTree.Content != nil {
				nodeTreeMap["content"] = nodeTree.Content
			}

			childrenTmp := make([]map[string]interface{}, 0)
			if nodeTree.Children != nil {
				for _, nodeLeaf := range nodeTree.Children {
					nodeLeafMap := map[string]interface{}{}
					if nodeLeaf.ID != nil {
						nodeLeafMap["id"] = nodeLeaf.ID
					}
					if nodeLeaf.Content != nil {
						nodeLeafMap["content"] = nodeLeaf.Content
					}

					innerChildrenTmp := make([]map[string]interface{}, 0)
					if nodeLeaf.Children != nil {
						for _, nodeInfoBase := range nodeLeaf.Children {
							nodeInfoBaseMap := map[string]interface{}{}
							if nodeInfoBase.ID != nil {
								nodeInfoBaseMap["id"] = nodeInfoBase.ID
							}
							if nodeInfoBase.Content != nil {
								nodeInfoBaseMap["content"] = nodeInfoBase.Content
							}
							innerChildrenTmp = append(innerChildrenTmp, nodeInfoBaseMap)
						}
					}
					nodeLeafMap["children"] = innerChildrenTmp
					childrenTmp = append(childrenTmp, nodeLeafMap)
				}
			}
			nodeTreeMap["children"] = childrenTmp
			nodeListTmp = append(nodeListTmp, nodeTreeMap)
		}
	}
	d.SetId(helper.DataResourceIdsHash(ids))

	districtListTmp := make([]map[string]interface{}, 0)
	if districtList != nil {
		for _, district := range districtList {
			districtMap := map[string]interface{}{}
			if district.ID != nil {
				districtMap["id"] = district.ID
			}
			if district.Name != nil {
				districtMap["name"] = district.Name
			}
			districtListTmp = append(districtListTmp, districtMap)
		}
	}

	netServiceListTmp := make([]map[string]interface{}, 0)
	if netServiceList != nil {
		for _, netService := range netServiceList {
			netServiceMap := map[string]interface{}{}
			if netService.ID != nil {
				netServiceMap["id"] = netService.ID
			}
			if netService.Name != nil {
				netServiceMap["name"] = netService.Name
			}
			netServiceListTmp = append(netServiceListTmp, netServiceMap)
		}
	}

	_ = d.Set("node_list", nodeListTmp)
	_ = d.Set("district_list", districtListTmp)
	_ = d.Set("net_service_list", netServiceListTmp)

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), nodeListTmp); e != nil {
			return e
		}
	}

	return nil
}

func logError(logId, msg string) error {
	log.Printf("[CRITAL]%s %s", logId, msg)
	return fmt.Errorf("%s", msg)
}
