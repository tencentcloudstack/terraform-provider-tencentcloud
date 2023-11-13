/*
Use this data source to query detailed information of tcr replication_instance

Example Usage

```hcl
data "tencentcloud_tcr_replication_instance" "replication_instance" {
  registry_id = "tcr-xx"
  replication_registry_id = "tcr-xx-1"
  replication_region_id = 1
  show_replication_log = false
        tags = {
    "createdBy" = "terraform"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrReplicationInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrReplicationInstanceRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Master registry id.",
			},

			"replication_registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Synchronization instance id.",
			},

			"replication_region_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Synchronization instance region id.",
			},

			"show_replication_log": {
				Optional:    true,
				Type:        schema.TypeBool,
				Description: "Whether to display the synchronization log.",
			},

			"replication_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sync status.",
			},

			"replication_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Sync complete time.",
			},

			"replication_log": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Sync log. Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource type. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Source image. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Destination resource. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sync status. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time. Note: This field may return null, indicating that no valid value can be obtained.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTcrReplicationInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_replication_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("replication_registry_id"); ok {
		paramMap["ReplicationRegistryId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("replication_region_id"); v != nil {
		paramMap["ReplicationRegionId"] = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("show_replication_log"); v != nil {
		paramMap["ShowReplicationLog"] = helper.Bool(v.(bool))
	}

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrReplicationInstanceByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		replicationStatus = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(replicationStatus))
	if replicationStatus != nil {
		_ = d.Set("replication_status", replicationStatus)
	}

	if replicationTime != nil {
		_ = d.Set("replication_time", replicationTime)
	}

	if replicationLog != nil {
		replicationLogMap := map[string]interface{}{}

		if replicationLog.ResourceType != nil {
			replicationLogMap["resource_type"] = replicationLog.ResourceType
		}

		if replicationLog.Source != nil {
			replicationLogMap["source"] = replicationLog.Source
		}

		if replicationLog.Destination != nil {
			replicationLogMap["destination"] = replicationLog.Destination
		}

		if replicationLog.Status != nil {
			replicationLogMap["status"] = replicationLog.Status
		}

		if replicationLog.StartTime != nil {
			replicationLogMap["start_time"] = replicationLog.StartTime
		}

		if replicationLog.EndTime != nil {
			replicationLogMap["end_time"] = replicationLog.EndTime
		}

		ids = append(ids, *replicationLog.ReplicationRegistryId)
		_ = d.Set("replication_log", replicationLogMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
