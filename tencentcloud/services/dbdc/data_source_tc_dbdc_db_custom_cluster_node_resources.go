package dbdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomClusterNodeResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomClusterNodeResourcesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "DB Custom cluster ID.",
			},

			"node_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Node ID list. Up to 50 node IDs per request (enforced by the cloud API).",
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
				Description: "DB Custom cluster node resource list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Node ID.",
						},
						"capacity": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node physical resource total capacity.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU cores.",
									},
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory, in GiB.",
									},
									"pods": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of pods.",
									},
								},
							},
						},
						"allocatable": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node allocatable capacity = Capacity - system reserved.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU cores.",
									},
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory, in GiB.",
									},
									"pods": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of pods.",
									},
								},
							},
						},
						"requests": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Sum of requests of all non-terminal pods on the node (including system pods).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU cores.",
									},
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory, in GiB.",
									},
									"pods": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of pods.",
									},
								},
							},
						},
						"limits": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Sum of limits of all non-terminal pods on the node (including system pods).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU cores.",
									},
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory, in GiB.",
									},
									"pods": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of pods.",
									},
								},
							},
						},
						"available": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Node schedulable remainder = max(0, Allocatable - Requests).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "CPU cores.",
									},
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Memory, in GiB.",
									},
									"pods": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of pods.",
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

func dataSourceTencentCloudDbdcDbCustomClusterNodeResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_cluster_node_resources.read")()
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
		tmpSet := make([]*string, 0, len(nodeIdsSet))
		for _, item := range nodeIdsSet {
			value := item.(string)
			tmpSet = append(tmpSet, helper.String(value))
		}
		paramMap["NodeIds"] = tmpSet
	}

	var respData []*dbdcv20201029.DBCustomClusterNodeResource
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDBCustomClusterNodeResourcesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
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

			nodeMap["capacity"] = flattenMetaResource(node.Capacity)
			nodeMap["allocatable"] = flattenMetaResource(node.Allocatable)
			nodeMap["requests"] = flattenMetaResource(node.Requests)
			nodeMap["limits"] = flattenMetaResource(node.Limits)
			nodeMap["available"] = flattenMetaResource(node.Available)

			nodeSetList = append(nodeSetList, nodeMap)
		}

		_ = d.Set("node_set", nodeSetList)
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

func flattenMetaResource(meta *dbdcv20201029.MetaResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if meta == nil {
		return result
	}

	metaMap := map[string]interface{}{}
	if meta.Cpu != nil {
		metaMap["cpu"] = *meta.Cpu
	}

	if meta.Memory != nil {
		metaMap["memory"] = *meta.Memory
	}

	if meta.Pods != nil {
		metaMap["pods"] = *meta.Pods
	}

	result = append(result, metaMap)
	return result
}
