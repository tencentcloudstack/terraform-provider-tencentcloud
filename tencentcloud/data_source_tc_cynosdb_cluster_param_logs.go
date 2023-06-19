/*
Use this data source to query detailed information of cynosdb cluster_param_logs

Example Usage

```hcl
data "tencentcloud_cynosdb_cluster_param_logs" "cluster_param_logs" {
  cluster_id    = "cynosdbmysql-bws8h88b"
  instance_ids  = ["cynosdbmysql-ins-afqx1hy0"]
  order_by      = "CreateTime"
  order_by_type = "DESC"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbClusterParamLogs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbClusterParamLogsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},
			"instance_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list, used to record specific instances of operations.",
			},
			"order_by": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort field, defining which field to sort based on when returning results.",
			},
			"order_by_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Define specific sorting rules, limited to one of desc, asc, DESC, or ASC.",
			},
			"cluster_param_logs": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter modification record note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter Name.",
						},
						"current_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Current value.",
						},
						"update_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modified value.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "modify state.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Update time.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
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

func dataSourceTencentCloudCynosdbClusterParamLogsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_cluster_param_logs.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId            = getLogId(contextNil)
		ctx              = context.WithValue(context.TODO(), logIdKey, logId)
		service          = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		clusterParamLogs []*cynosdb.ClusterParamModifyLog
		clusterId        string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
		clusterId = v.(string)
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		paramMap["InstanceIds"] = helper.InterfacesStringsPoint(instanceIdsSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		paramMap["OrderBy"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_by_type"); ok {
		paramMap["OrderByType"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbClusterParamLogsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		clusterParamLogs = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(clusterParamLogs))

	if clusterParamLogs != nil {
		for _, clusterParamModifyLog := range clusterParamLogs {
			clusterParamModifyLogMap := map[string]interface{}{}

			if clusterParamModifyLog.ParamName != nil {
				clusterParamModifyLogMap["param_name"] = clusterParamModifyLog.ParamName
			}

			if clusterParamModifyLog.CurrentValue != nil {
				clusterParamModifyLogMap["current_value"] = clusterParamModifyLog.CurrentValue
			}

			if clusterParamModifyLog.UpdateValue != nil {
				clusterParamModifyLogMap["update_value"] = clusterParamModifyLog.UpdateValue
			}

			if clusterParamModifyLog.Status != nil {
				clusterParamModifyLogMap["status"] = clusterParamModifyLog.Status
			}

			if clusterParamModifyLog.CreateTime != nil {
				clusterParamModifyLogMap["create_time"] = clusterParamModifyLog.CreateTime
			}

			if clusterParamModifyLog.UpdateTime != nil {
				clusterParamModifyLogMap["update_time"] = clusterParamModifyLog.UpdateTime
			}

			if clusterParamModifyLog.ClusterId != nil {
				clusterParamModifyLogMap["cluster_id"] = clusterParamModifyLog.ClusterId
			}

			if clusterParamModifyLog.InstanceId != nil {
				clusterParamModifyLogMap["instance_id"] = clusterParamModifyLog.InstanceId
			}

			tmpList = append(tmpList, clusterParamModifyLogMap)
		}

		_ = d.Set("cluster_param_logs", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
