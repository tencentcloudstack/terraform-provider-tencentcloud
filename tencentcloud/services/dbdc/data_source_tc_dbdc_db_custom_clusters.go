package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomClustersRead,
		Schema: map[string]*schema.Schema{
			"cluster_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Query by one or more Cluster IDs. Maximum 100 IDs per request.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions. Supported filter names: cluster-name (exact match), cluster-status (Creating, Running, Destroying).",
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

			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter by tag Key and Value.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag key.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"cluster_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "DB Custom cluster list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region supported by the cluster.",
						},
						"cluster_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster level. Default value: L500.",
						},
						"cluster_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster status. Values: Creating, Running, Destroying.",
						},
						"cluster_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster version.",
						},
						"cluster_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes in the cluster.",
						},
						"cluster_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster description.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Cluster tag information. Note: This field may return null, indicating that no valid value can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
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

func dataSourceTencentCloudDbdcDbCustomClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_ids"); ok {
		clusterIdsSet := v.([]interface{})
		tmpSet := make([]*string, 0, len(clusterIdsSet))
		for _, item := range clusterIdsSet {
			clusterId := item.(string)
			tmpSet = append(tmpSet, helper.String(clusterId))
		}
		paramMap["ClusterIds"] = tmpSet
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

	if v, ok := d.GetOk("tags"); ok {
		tagsSet := v.([]interface{})
		tmpSet := make([]*dbdcv20201029.Tag, 0, len(tagsSet))
		for _, item := range tagsSet {
			tagsMap := item.(map[string]interface{})
			tag := dbdcv20201029.Tag{}
			if v, ok := tagsMap["key"].(string); ok && v != "" {
				tag.Key = helper.String(v)
			}

			if v, ok := tagsMap["value"].(string); ok {
				tag.Value = helper.String(v)
			}
			tmpSet = append(tmpSet, &tag)
		}
		paramMap["Tags"] = tmpSet
	}

	var respData []*dbdcv20201029.DBCustomCluster
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, _, e := service.DescribeDBCustomClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	clusterSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, cluster := range respData {
			clusterMap := map[string]interface{}{}
			if cluster.ClusterId != nil {
				clusterMap["cluster_id"] = cluster.ClusterId
			}

			if cluster.ClusterName != nil {
				clusterMap["cluster_name"] = cluster.ClusterName
			}

			if cluster.Region != nil {
				clusterMap["region"] = cluster.Region
			}

			if cluster.ClusterLevel != nil {
				clusterMap["cluster_level"] = cluster.ClusterLevel
			}

			if cluster.ClusterStatus != nil {
				clusterMap["cluster_status"] = cluster.ClusterStatus
			}

			if cluster.ClusterVersion != nil {
				clusterMap["cluster_version"] = cluster.ClusterVersion
			}

			if cluster.ClusterNodeNum != nil {
				clusterMap["cluster_node_num"] = cluster.ClusterNodeNum
			}

			if cluster.ClusterDescription != nil {
				clusterMap["cluster_description"] = cluster.ClusterDescription
			}

			if cluster.CreatedTime != nil {
				clusterMap["created_time"] = cluster.CreatedTime
			}

			if cluster.Tags != nil {
				tagsList := make([]map[string]interface{}, 0, len(cluster.Tags))
				for _, tag := range cluster.Tags {
					tagMap := map[string]interface{}{}
					if tag.Key != nil {
						tagMap["key"] = tag.Key
					}
					if tag.Value != nil {
						tagMap["value"] = tag.Value
					}
					tagsList = append(tagsList, tagMap)
				}
				clusterMap["tags"] = tagsList
			}

			clusterSetList = append(clusterSetList, clusterMap)
		}

		_ = d.Set("cluster_set", clusterSetList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
