package cynosdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCynosdbProxyNode() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbProxyNodeRead,
		Schema: map[string]*schema.Schema{
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field, value range:CREATETIME: creation time; PRIODENDTIME: expiration time.",
			},
			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort type, value range:ASC: ascending sort; DESC: descending sort.",
			},
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Search criteria, if there are multiple filters, the relationship between the filters is a logical AND relationship.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"names": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Search String.",
						},
						"values": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Search String.",
						},
						"exact_match": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Exact match or not.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Search Fields. Supported: Status, ProxyNodeId, ClusterId.",
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Operator.",
						},
					},
				},
			},
			"proxy_node_infos": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Database Agent Node List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Agent Node ID.",
						},
						"proxy_node_connections": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The current number of connections of the node. The DescribeProxyNodes interface does not return a value for this field.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database Agent Node CPU.",
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Database Agent Node Memory.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Agent Node Status.",
						},
						"proxy_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Database Agent Group ID.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User AppID.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "region.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability Zone.",
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

func dataSourceTencentCloudCynosdbProxyNodeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cynosdb_proxy_node.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId          = tccommon.GetLogId(tccommon.ContextNil)
		ctx            = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service        = CynosdbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		proxyNodeInfos []*cynosdb.ProxyNodeInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*cynosdb.QueryFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			queryFilter := cynosdb.QueryFilter{}
			queryFilterMap := item.(map[string]interface{})

			if v, ok := queryFilterMap["names"]; ok {
				namesSet := v.(*schema.Set).List()
				queryFilter.Names = helper.InterfacesStringsPoint(namesSet)
			}
			if v, ok := queryFilterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				queryFilter.Values = helper.InterfacesStringsPoint(valuesSet)
			}
			if v, ok := queryFilterMap["exact_match"]; ok {
				queryFilter.ExactMatch = helper.Bool(v.(bool))
			}
			if v, ok := queryFilterMap["name"]; ok {
				queryFilter.Name = helper.String(v.(string))
			}
			if v, ok := queryFilterMap["operator"]; ok {
				queryFilter.Operator = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &queryFilter)
		}
		paramMap["filters"] = tmpSet
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbProxyNodeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		proxyNodeInfos = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(proxyNodeInfos))
	tmpList := make([]map[string]interface{}, 0, len(proxyNodeInfos))

	if proxyNodeInfos != nil {
		for _, proxyNodeInfo := range proxyNodeInfos {
			proxyNodeInfoMap := map[string]interface{}{}

			if proxyNodeInfo.ProxyNodeId != nil {
				proxyNodeInfoMap["proxy_node_id"] = proxyNodeInfo.ProxyNodeId
			}

			if proxyNodeInfo.ProxyNodeConnections != nil {
				proxyNodeInfoMap["proxy_node_connections"] = proxyNodeInfo.ProxyNodeConnections
			}

			if proxyNodeInfo.Cpu != nil {
				proxyNodeInfoMap["cpu"] = proxyNodeInfo.Cpu
			}

			if proxyNodeInfo.Mem != nil {
				proxyNodeInfoMap["mem"] = proxyNodeInfo.Mem
			}

			if proxyNodeInfo.Status != nil {
				proxyNodeInfoMap["status"] = proxyNodeInfo.Status
			}

			if proxyNodeInfo.ProxyGroupId != nil {
				proxyNodeInfoMap["proxy_group_id"] = proxyNodeInfo.ProxyGroupId
			}

			if proxyNodeInfo.ClusterId != nil {
				proxyNodeInfoMap["cluster_id"] = proxyNodeInfo.ClusterId
			}

			if proxyNodeInfo.AppId != nil {
				proxyNodeInfoMap["app_id"] = proxyNodeInfo.AppId
			}

			if proxyNodeInfo.Region != nil {
				proxyNodeInfoMap["region"] = proxyNodeInfo.Region
			}

			if proxyNodeInfo.Zone != nil {
				proxyNodeInfoMap["zone"] = proxyNodeInfo.Zone
			}

			ids = append(ids, *proxyNodeInfo.ProxyNodeId)
			tmpList = append(tmpList, proxyNodeInfoMap)
		}

		_ = d.Set("proxy_node_infos", tmpList)
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
