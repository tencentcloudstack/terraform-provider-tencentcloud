// Code generated by iacg; DO NOT EDIT.
package postgresql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPostgresqlDedicatedClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPostgresqlDedicatedClustersRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Querying based on one or more filtering criteria, the currently supported filtering criteria are: dedicated-cluster-id: filtering by dedicated cluster ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Filter values.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			// computed
			"dedicated_cluster_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Dedicated cluster set info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Dedicated cluster ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name.",
						},
						"zone": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Zone.",
						},
						"standby_dedicated_cluster_set": {
							Type:        schema.TypeSet,
							Optional:    true,
							Description: "Disaster recovery cluster.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"instance_count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Instance count.",
						},
						"cpu_total": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Cpu total.",
						},
						"cpu_available": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of available CPUs.",
						},
						"mem_total": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Total amount of memory.",
						},
						"mem_available": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Available Memory.",
						},
						"disk_total": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Total number of disks.",
						},
						"disk_available": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Disk availability.",
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

func dataSourceTencentCloudPostgresqlDedicatedClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_postgresql_dedicated_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(nil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service  = PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		respData []*postgresql.DedicatedCluster
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*postgresql.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			filter := postgresql.Filter{}
			if v, ok := filtersMap["name"]; ok {
				filter.Name = helper.String(v.(string))
			}

			if v, ok := filtersMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				for i := range valuesSet {
					values := valuesSet[i].(string)
					filter.Values = append(filter.Values, helper.String(values))
				}
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePostgresqlDedicatedClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	dedicatedClusterSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, dedicatedClusterSet := range respData {
			dedicatedClusterSetMap := map[string]interface{}{}

			if dedicatedClusterSet.DedicatedClusterId != nil {
				dedicatedClusterSetMap["dedicated_cluster_id"] = dedicatedClusterSet.DedicatedClusterId
			}

			if dedicatedClusterSet.Name != nil {
				dedicatedClusterSetMap["name"] = dedicatedClusterSet.Name
			}

			if dedicatedClusterSet.Zone != nil {
				dedicatedClusterSetMap["zone"] = dedicatedClusterSet.Zone
			}

			if dedicatedClusterSet.StandbyDedicatedClusterSet != nil {
				dedicatedClusterSetMap["standby_dedicated_cluster_set"] = helper.PStrings(dedicatedClusterSet.StandbyDedicatedClusterSet)
			}

			if dedicatedClusterSet.InstanceCount != nil {
				dedicatedClusterSetMap["instance_count"] = dedicatedClusterSet.InstanceCount
			}

			if dedicatedClusterSet.CpuTotal != nil {
				dedicatedClusterSetMap["cpu_total"] = dedicatedClusterSet.CpuTotal
			}

			if dedicatedClusterSet.CpuAvailable != nil {
				dedicatedClusterSetMap["cpu_available"] = dedicatedClusterSet.CpuAvailable
			}

			if dedicatedClusterSet.MemTotal != nil {
				dedicatedClusterSetMap["mem_total"] = dedicatedClusterSet.MemTotal
			}

			if dedicatedClusterSet.MemAvailable != nil {
				dedicatedClusterSetMap["mem_available"] = dedicatedClusterSet.MemAvailable
			}

			if dedicatedClusterSet.DiskTotal != nil {
				dedicatedClusterSetMap["disk_total"] = dedicatedClusterSet.DiskTotal
			}

			if dedicatedClusterSet.DiskAvailable != nil {
				dedicatedClusterSetMap["disk_available"] = dedicatedClusterSet.DiskAvailable
			}

			ids = append(ids, *dedicatedClusterSet.DedicatedClusterId)
			dedicatedClusterSetList = append(dedicatedClusterSetList, dedicatedClusterSetMap)
		}

		_ = d.Set("dedicated_cluster_set", dedicatedClusterSetList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dedicatedClusterSetList); e != nil {
			return e
		}
	}

	return nil
}
