package dbdc

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDbdcDbCustomClusterResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbdcDbCustomClusterResourcesRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cluster ID.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"node_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of worker nodes participating in the aggregation (excluding control plane nodes).",
			},

			"capacity": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "The sum of the physical total resource capacity of all nodes in the cluster.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory size in GiB.",
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
				MaxItems:    1,
				Description: "The sum of the allocatable capacity of all nodes in the cluster (= Capacity - system reservation).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory size in GiB.",
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
				MaxItems:    1,
				Description: "The sum of the requests of all non-terminal pods in the cluster (including system pods).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory size in GiB.",
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
				MaxItems:    1,
				Description: "The sum of the limits of all non-terminal pods in the cluster (including system pods, the Pods field has no semantics and is fixed to 0).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory size in GiB.",
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
				MaxItems:    1,
				Description: "Cluster schedulable remaining capacity (the sum of max(0, Allocatable - Requests) for all nodes).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores.",
						},
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory size in GiB.",
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
	}
}

func dataSourceTencentCloudDbdcDbCustomClusterResourcesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dbdc_db_custom_cluster_resources.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(nil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = DbdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		clusterId string
	)

	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
	}

	respData, reqErr := service.DescribeDBCustomClusterResources(ctx, clusterId)
	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return reqErr
	}

	if respData == nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
	}

	if respData != nil {
		if respData.NodeCount != nil {
			_ = d.Set("node_count", int(*respData.NodeCount))
		}

		if respData.Capacity != nil {
			capacityMap := flattenDbdcMetaResource(respData.Capacity)
			if len(capacityMap) > 0 {
				_ = d.Set("capacity", []interface{}{capacityMap})
			}
		}

		if respData.Allocatable != nil {
			allocatableMap := flattenDbdcMetaResource(respData.Allocatable)
			if len(allocatableMap) > 0 {
				_ = d.Set("allocatable", []interface{}{allocatableMap})
			}
		}

		if respData.Requests != nil {
			requestsMap := flattenDbdcMetaResource(respData.Requests)
			if len(requestsMap) > 0 {
				_ = d.Set("requests", []interface{}{requestsMap})
			}
		}

		if respData.Limits != nil {
			limitsMap := flattenDbdcMetaResource(respData.Limits)
			if len(limitsMap) > 0 {
				_ = d.Set("limits", []interface{}{limitsMap})
			}
		}

		if respData.Available != nil {
			availableMap := flattenDbdcMetaResource(respData.Available)
			if len(availableMap) > 0 {
				_ = d.Set("available", []interface{}{availableMap})
			}
		}
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

func flattenDbdcMetaResource(metaResource *dbdcv20201029.MetaResource) map[string]interface{} {
	if metaResource == nil {
		return nil
	}

	result := map[string]interface{}{}
	if metaResource.Cpu != nil {
		result["cpu"] = *metaResource.Cpu
	}

	if metaResource.Memory != nil {
		result["memory"] = *metaResource.Memory
	}

	if metaResource.Pods != nil {
		result["pods"] = int(*metaResource.Pods)
	}

	return result
}
