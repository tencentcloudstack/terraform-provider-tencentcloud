package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	elasticsearch "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/es/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudElasticsearchViews() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudElasticsearchViewsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"cluster_view": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Cluster view.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"health": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Cluster health status.",
						},
						"visible": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Whether the cluster is visible.",
						},
						"break": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Whether the cluster is broken or not.",
						},
						"avg_disk_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average disk utilization.",
						},
						"avg_mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average memory utilization.",
						},
						"avg_cpu_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average cpu utilization.",
						},
						"total_disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total storage size of cluster.",
						},
						"target_node_types": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Client request node.",
						},
						"node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of online nodes.",
						},
						"total_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of nodes.",
						},
						"data_node_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of data nodes.",
						},
						"index_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Index number.",
						},
						"doc_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of documents.",
						},
						"disk_used_in_bytes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bytes used on disk.",
						},
						"shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Shard number.",
						},
						"primary_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Primary shard number.",
						},
						"relocating_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Relocating shard number.",
						},
						"initializing_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Initializing shard number.",
						},
						"unassigned_shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unassigned shard number.",
						},
						"total_cos_storage": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Storage capacity of COS Enterprise Edition (in GB).",
						},
						"searchable_snapshot_cos_bucket": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise cluster searchable bucket name stored in snapshot cos.",
						},
						"searchable_snapshot_cos_app_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Enterprise cluster can search the appid to which snapshot cos belongs.",
						},
					},
				},
			},

			"nodes_view": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Node View.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node id.",
						},
						"node_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ip.",
						},
						"visible": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Whether the node is visible.",
						},
						"break": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Whether or not to break.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total disk size of node.",
						},
						"disk_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Disk usage.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node memory size (in GB).",
						},
						"mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory usage.",
						},
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU number.",
						},
						"cpu_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "CPU usage.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Zone.",
						},
						"node_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node role.",
						},
						"node_http_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node HTTP IP.",
						},
						"jvm_mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "JVM memory usage.",
						},
						"shard_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of node fragments.",
						},
						"disk_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of disk ID on the node.",
						},
						"hidden": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether it is a hidden availability zone.",
						},
						"is_coordination_node": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether to act as a coordinator node or not.",
						},
					},
				},
			},

			"kibanas_view": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Kibanas view.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Kibana node ip.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size.",
						},
						"disk_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Disk usage.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Node memory size.",
						},
						"mem_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory usage.",
						},
						"cpu_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU number.",
						},
						"cpu_usage": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "cpu usage.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "zone.",
						},
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node id.",
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

func dataSourceTencentCloudElasticsearchViewsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_elasticsearch_views.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var instanceId string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		paramMap["InstanceId"] = helper.String(instanceId)
	}

	service := ElasticsearchService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		clusterView  *elasticsearch.ClusterView
		nodesViews   []*elasticsearch.NodeView
		kibanasViews []*elasticsearch.KibanaView
	)

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		var e error
		clusterView, nodesViews, kibanasViews, e = service.DescribeElasticsearchViewsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	clusterViewMapList := []interface{}{}
	nodesViewMapList := []interface{}{}
	kibanasViewMapList := []interface{}{}

	if clusterView != nil {
		clusterViewMap := map[string]interface{}{}
		if clusterView.Health != nil {
			clusterViewMap["health"] = clusterView.Health
		}

		if clusterView.Visible != nil {
			clusterViewMap["visible"] = clusterView.Visible
		}

		if clusterView.Break != nil {
			clusterViewMap["break"] = clusterView.Break
		}

		if clusterView.AvgDiskUsage != nil {
			clusterViewMap["avg_disk_usage"] = clusterView.AvgDiskUsage
		}

		if clusterView.AvgMemUsage != nil {
			clusterViewMap["avg_mem_usage"] = clusterView.AvgMemUsage
		}

		if clusterView.AvgCpuUsage != nil {
			clusterViewMap["avg_cpu_usage"] = clusterView.AvgCpuUsage
		}

		if clusterView.TotalDiskSize != nil {
			clusterViewMap["total_disk_size"] = clusterView.TotalDiskSize
		}

		if clusterView.TargetNodeTypes != nil {
			clusterViewMap["target_node_types"] = clusterView.TargetNodeTypes
		}

		if clusterView.NodeNum != nil {
			clusterViewMap["node_num"] = clusterView.NodeNum
		}

		if clusterView.TotalNodeNum != nil {
			clusterViewMap["total_node_num"] = clusterView.TotalNodeNum
		}

		if clusterView.DataNodeNum != nil {
			clusterViewMap["data_node_num"] = clusterView.DataNodeNum
		}

		if clusterView.IndexNum != nil {
			clusterViewMap["index_num"] = clusterView.IndexNum
		}

		if clusterView.DocNum != nil {
			clusterViewMap["doc_num"] = clusterView.DocNum
		}

		if clusterView.DiskUsedInBytes != nil {
			clusterViewMap["disk_used_in_bytes"] = clusterView.DiskUsedInBytes
		}

		if clusterView.ShardNum != nil {
			clusterViewMap["shard_num"] = clusterView.ShardNum
		}

		if clusterView.PrimaryShardNum != nil {
			clusterViewMap["primary_shard_num"] = clusterView.PrimaryShardNum
		}

		if clusterView.RelocatingShardNum != nil {
			clusterViewMap["relocating_shard_num"] = clusterView.RelocatingShardNum
		}

		if clusterView.InitializingShardNum != nil {
			clusterViewMap["initializing_shard_num"] = clusterView.InitializingShardNum
		}

		if clusterView.UnassignedShardNum != nil {
			clusterViewMap["unassigned_shard_num"] = clusterView.UnassignedShardNum
		}

		if clusterView.TotalCosStorage != nil {
			clusterViewMap["total_cos_storage"] = clusterView.TotalCosStorage
		}

		if clusterView.SearchableSnapshotCosBucket != nil {
			clusterViewMap["searchable_snapshot_cos_bucket"] = clusterView.SearchableSnapshotCosBucket
		}

		if clusterView.SearchableSnapshotCosAppId != nil {
			clusterViewMap["searchable_snapshot_cos_app_id"] = clusterView.SearchableSnapshotCosAppId
		}
		clusterViewMapList = append(clusterViewMapList, clusterViewMap)
		_ = d.Set("cluster_view", clusterViewMapList)
	}

	if nodesViews != nil {
		for _, nodeView := range nodesViews {
			nodeViewMap := map[string]interface{}{}

			if nodeView.NodeId != nil {
				nodeViewMap["node_id"] = nodeView.NodeId
			}

			if nodeView.NodeIp != nil {
				nodeViewMap["node_ip"] = nodeView.NodeIp
			}

			if nodeView.Visible != nil {
				nodeViewMap["visible"] = nodeView.Visible
			}

			if nodeView.Break != nil {
				nodeViewMap["break"] = nodeView.Break
			}

			if nodeView.DiskSize != nil {
				nodeViewMap["disk_size"] = nodeView.DiskSize
			}

			if nodeView.DiskUsage != nil {
				nodeViewMap["disk_usage"] = nodeView.DiskUsage
			}

			if nodeView.MemSize != nil {
				nodeViewMap["mem_size"] = nodeView.MemSize
			}

			if nodeView.MemUsage != nil {
				nodeViewMap["mem_usage"] = nodeView.MemUsage
			}

			if nodeView.CpuNum != nil {
				nodeViewMap["cpu_num"] = nodeView.CpuNum
			}

			if nodeView.CpuUsage != nil {
				nodeViewMap["cpu_usage"] = nodeView.CpuUsage
			}

			if nodeView.Zone != nil {
				nodeViewMap["zone"] = nodeView.Zone
			}

			if nodeView.NodeRole != nil {
				nodeViewMap["node_role"] = nodeView.NodeRole
			}

			if nodeView.NodeHttpIp != nil {
				nodeViewMap["node_http_ip"] = nodeView.NodeHttpIp
			}

			if nodeView.JvmMemUsage != nil {
				nodeViewMap["jvm_mem_usage"] = nodeView.JvmMemUsage
			}

			if nodeView.ShardNum != nil {
				nodeViewMap["shard_num"] = nodeView.ShardNum
			}

			if nodeView.DiskIds != nil {
				nodeViewMap["disk_ids"] = nodeView.DiskIds
			}

			if nodeView.Hidden != nil {
				nodeViewMap["hidden"] = nodeView.Hidden
			}

			if nodeView.IsCoordinationNode != nil {
				nodeViewMap["is_coordination_node"] = nodeView.IsCoordinationNode
			}

			nodesViewMapList = append(nodesViewMapList, nodeViewMap)
		}

		_ = d.Set("nodes_view", nodesViewMapList)
	}

	if kibanasViews != nil {
		for _, kibanaView := range kibanasViews {
			kibanaViewMap := map[string]interface{}{}

			if kibanaView.Ip != nil {
				kibanaViewMap["ip"] = kibanaView.Ip
			}

			if kibanaView.DiskSize != nil {
				kibanaViewMap["disk_size"] = kibanaView.DiskSize
			}

			if kibanaView.DiskUsage != nil {
				kibanaViewMap["disk_usage"] = kibanaView.DiskUsage
			}

			if kibanaView.MemSize != nil {
				kibanaViewMap["mem_size"] = kibanaView.MemSize
			}

			if kibanaView.MemUsage != nil {
				kibanaViewMap["mem_usage"] = kibanaView.MemUsage
			}

			if kibanaView.CpuNum != nil {
				kibanaViewMap["cpu_num"] = kibanaView.CpuNum
			}

			if kibanaView.CpuUsage != nil {
				kibanaViewMap["cpu_usage"] = kibanaView.CpuUsage
			}

			if kibanaView.Zone != nil {
				kibanaViewMap["zone"] = kibanaView.Zone
			}

			if kibanaView.NodeId != nil {
				kibanaViewMap["node_id"] = kibanaView.NodeId
			}

			kibanasViewMapList = append(kibanasViewMapList, kibanaViewMap)
		}

		_ = d.Set("kibanas_view", kibanasViewMapList)
	}

	d.SetId(instanceId)
	result := make(map[string]interface{})
	result["cluster_view"] = clusterViewMapList
	result["nodes_view"] = nodesViewMapList
	result["kibanas_view"] = kibanasViewMapList
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), result); e != nil {
			return e
		}
	}
	return nil
}
