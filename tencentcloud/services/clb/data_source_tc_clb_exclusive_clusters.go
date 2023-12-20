package clb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClbExclusiveClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbExclusiveClustersRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter to query the list of AZ resources as detailed below: cluster-type - String - Required: No - (Filter condition) Filter by cluster type, such as TGW. cluster-id - String - Required: No - (Filter condition) Filter by cluster ID, such as tgw-xxxxxxxx. cluster-name - String - Required: No - (Filter condition) Filter by cluster name, such as test-xxxxxx. cluster-tag - String - Required: No - (Filter condition) Filter by cluster tag, such as TAG-xxxxx. vip - String - Required: No - (Filter condition) Filter by vip in the cluster, such as x.x.x.x. network - String - Required: No - (Filter condition) Filter by cluster network type, such as Public or Private. zone - String - Required: No - (Filter condition) Filter by cluster zone, such as ap-guangzhou-1. isp - String - Required: No - (Filter condition) Filter by TGW cluster isp type, such as BGP. loadblancer-id - String - Required: No - (Filter condition) Filter by loadblancer-id in the cluste, such as lb-xxxxxxxx.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Required:    true,
							Description: "Filter value array.",
						},
					},
				},
			},

			"cluster_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster ID.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster name.",
						},
						"cluster_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster type: TGW, STGW, VPCGW.",
						},
						"cluster_tag": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated layer-7 tag. Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: ".",
						},
						"network": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "cluster network type.",
						},
						"max_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of connections.",
						},
						"max_in_flow": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum incoming Bandwidth.",
						},
						"max_in_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum incoming packet.",
						},
						"max_out_flow": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum output bandwidth.",
						},
						"max_out_pkg": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum output packet.",
						},
						"max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of new connections.",
						},
						"http_max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of new http connections.",
						},
						"https_max_new_conn": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of new https connections.",
						},
						"http_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Http Qps.",
						},
						"https_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Https Qps.",
						},
						"resource_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of resources in the cluster.",
						},
						"idle_resource_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of free resources in the cluster.",
						},
						"load_balance_director_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of forwarders in the cluster.",
						},
						"isp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Isp: BGP, CMCC,CUCC,CTCC,INTERNAL.",
						},
						"clusters_zone": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Availability zone where the cluster is located.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"master_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Availability master zone where the cluster is located.",
									},
									"slave_zone": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Availability slave zone where the cluster is located.",
									},
								},
							},
						},
						"clusters_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "clusters version.",
						},
						"disaster_recovery_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster disaster recovery type:SINGLE-ZONE, DISASTER-RECOVERY, MUTUAL-DISASTER-RECOVERY.",
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

func dataSourceTencentCloudClbExclusiveClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_clb_exclusive_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*clb.Filter, 0, len(filtersSet))

		for _, item := range filtersSet {
			filter := clb.Filter{}
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
		paramMap["Filters"] = tmpSet
	}

	service := ClbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var clusterSet []*clb.Cluster

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClbExclusiveClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		clusterSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(clusterSet))
	tmpList := make([]map[string]interface{}, 0, len(clusterSet))

	if clusterSet != nil {
		for _, cluster := range clusterSet {
			clusterMap := map[string]interface{}{}

			if cluster.ClusterId != nil {
				clusterMap["cluster_id"] = cluster.ClusterId
			}

			if cluster.ClusterName != nil {
				clusterMap["cluster_name"] = cluster.ClusterName
			}

			if cluster.ClusterType != nil {
				clusterMap["cluster_type"] = cluster.ClusterType
			}

			if cluster.ClusterTag != nil {
				clusterMap["cluster_tag"] = cluster.ClusterTag
			}

			if cluster.Zone != nil {
				clusterMap["zone"] = cluster.Zone
			}

			if cluster.Network != nil {
				clusterMap["network"] = cluster.Network
			}

			if cluster.MaxConn != nil {
				clusterMap["max_conn"] = cluster.MaxConn
			}

			if cluster.MaxInFlow != nil {
				clusterMap["max_in_flow"] = cluster.MaxInFlow
			}

			if cluster.MaxInPkg != nil {
				clusterMap["max_in_pkg"] = cluster.MaxInPkg
			}

			if cluster.MaxOutFlow != nil {
				clusterMap["max_out_flow"] = cluster.MaxOutFlow
			}

			if cluster.MaxOutPkg != nil {
				clusterMap["max_out_pkg"] = cluster.MaxOutPkg
			}

			if cluster.MaxNewConn != nil {
				clusterMap["max_new_conn"] = cluster.MaxNewConn
			}

			if cluster.HTTPMaxNewConn != nil {
				clusterMap["http_max_new_conn"] = cluster.HTTPMaxNewConn
			}

			if cluster.HTTPSMaxNewConn != nil {
				clusterMap["https_max_new_conn"] = cluster.HTTPSMaxNewConn
			}

			if cluster.HTTPQps != nil {
				clusterMap["http_qps"] = cluster.HTTPQps
			}

			if cluster.HTTPSQps != nil {
				clusterMap["https_qps"] = cluster.HTTPSQps
			}

			if cluster.ResourceCount != nil {
				clusterMap["resource_count"] = cluster.ResourceCount
			}

			if cluster.IdleResourceCount != nil {
				clusterMap["idle_resource_count"] = cluster.IdleResourceCount
			}

			if cluster.LoadBalanceDirectorCount != nil {
				clusterMap["load_balance_director_count"] = cluster.LoadBalanceDirectorCount
			}

			if cluster.Isp != nil {
				clusterMap["isp"] = cluster.Isp
			}

			if cluster.ClustersZone != nil {
				clustersZoneMap := map[string]interface{}{}

				if cluster.ClustersZone.MasterZone != nil {
					clustersZoneMap["master_zone"] = cluster.ClustersZone.MasterZone
				}

				if cluster.ClustersZone.SlaveZone != nil {
					clustersZoneMap["slave_zone"] = cluster.ClustersZone.SlaveZone
				}

				clusterMap["clusters_zone"] = []interface{}{clustersZoneMap}
			}

			if cluster.ClustersVersion != nil {
				clusterMap["clusters_version"] = cluster.ClustersVersion
			}

			if cluster.DisasterRecoveryType != nil {
				clusterMap["disaster_recovery_type"] = cluster.DisasterRecoveryType
			}

			ids = append(ids, *cluster.ClusterId)
			tmpList = append(tmpList, clusterMap)
		}

		_ = d.Set("cluster_set", tmpList)
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
