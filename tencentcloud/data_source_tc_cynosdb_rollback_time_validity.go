/*
Use this data source to query detailed information of cynosdb rollback_time_validity

Example Usage

```hcl
data "tencentcloud_cynosdb_rollback_time_validity" "rollback_time_validity" {
  cluster_id = ""
  expect_time = ""
  expect_time_thresh =
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

func dataSourceTencentCloudCynosdbRollbackTimeValidity() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbRollbackTimeValidityRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Cluster ID.",
			},

			"expect_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Expected point in time for rollback.",
			},

			"expect_time_thresh": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "The allowable error range for rollback time points.",
			},

			"pool_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Storage poolID.",
			},

			"query_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The rollback task ID needs to be passed in when rolling back at this time point in the future.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Is the time point valid: pass, the test passed; Fail, detection failed.",
			},

			"suggest_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Suggested time point, this value is only valid when Status is failed.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudCynosdbRollbackTimeValidityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_rollback_time_validity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("expect_time"); ok {
		paramMap["ExpectTime"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("expect_time_thresh"); v != nil {
		paramMap["ExpectTimeThresh"] = helper.IntUint64(v.(int))
	}

	service := CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbRollbackTimeValidityByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		poolId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(poolId))
	if poolId != nil {
		_ = d.Set("pool_id", poolId)
	}

	if queryId != nil {
		_ = d.Set("query_id", queryId)
	}

	if status != nil {
		_ = d.Set("status", status)
	}

	if suggestTime != nil {
		_ = d.Set("suggest_time", suggestTime)
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
