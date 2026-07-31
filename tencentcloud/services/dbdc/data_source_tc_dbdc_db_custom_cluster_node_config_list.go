package dbdc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomClusterNodeConfigList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomClusterNodeConfigListRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DB Custom cluster ID.",
			},

			"node_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Specifies the NodeId list to query. Up to 100 NodeIds per request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
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
				Description: "DB Custom cluster node config list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"labels": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node labels.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Label key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Label value.",
									},
								},
							},
						},
						"taints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node taints.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Taint key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Taint value.",
									},
									"effect": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Taint effect. Valid values: NoSchedule, PreferNoSchedule, NoExecute.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDbdcDbCustomClusterNodeConfigListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_cluster_node_config_list.read")()
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

	if v, ok := d.GetOk("node_ids"); ok {
		nodeIdsSet := v.([]interface{})
		nodeIds := make([]*string, 0, len(nodeIdsSet))
		for i := range nodeIdsSet {
			nodeIds = append(nodeIds, helper.String(nodeIdsSet[i].(string)))
		}
		paramMap["NodeIds"] = nodeIds
	}

	respData, errRet := service.DescribeDBCustomClusterNodeConfigByFilter(ctx, paramMap)
	if errRet != nil {
		return errRet
	}

	if respData == nil || len(respData) == 0 {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return nil
	}

	nodeSetList := make([]map[string]interface{}, 0, len(respData))
	for _, nodeConfig := range respData {
		nodeMap := map[string]interface{}{}
		if nodeConfig.NodeId != nil {
			nodeMap["node_id"] = nodeConfig.NodeId
		}

		labelsList := make([]map[string]interface{}, 0)
		if nodeConfig.Labels != nil {
			for _, label := range nodeConfig.Labels {
				labelMap := map[string]interface{}{}
				if label.Key != nil {
					labelMap["key"] = label.Key
				}
				if label.Value != nil {
					labelMap["value"] = label.Value
				}
				labelsList = append(labelsList, labelMap)
			}
		}
		nodeMap["labels"] = labelsList

		taintsList := make([]map[string]interface{}, 0)
		if nodeConfig.Taints != nil {
			for _, taint := range nodeConfig.Taints {
				taintMap := map[string]interface{}{}
				if taint.Key != nil {
					taintMap["key"] = taint.Key
				}
				if taint.Value != nil {
					taintMap["value"] = taint.Value
				}
				if taint.Effect != nil {
					taintMap["effect"] = taint.Effect
				}
				taintsList = append(taintsList, taintMap)
			}
		}
		nodeMap["taints"] = taintsList

		nodeSetList = append(nodeSetList, nodeMap)
	}

	_ = d.Set("node_set", nodeSetList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
