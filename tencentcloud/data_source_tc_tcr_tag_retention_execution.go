/*
Use this data source to query detailed information of tcr tag_retention_execution

Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_execution" "tag_retention_execution" {
  registry_id = "tcr-xx"
  retention_id = 1
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTcrTagRetentionExecution() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcrTagRetentionExecutionRead,
		Schema: map[string]*schema.Schema{
			"registry_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance id.",
			},

			"retention_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Retention id.",
			},

			"retention_execution_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of version retention execution records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"execution_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Execution id.",
						},
						"retention_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Retention id.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution start time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution end time.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Execution status: Failed, Succeed, Stopped, InProgress.",
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

func dataSourceTencentCloudTcrTagRetentionExecutionRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcr_tag_retention_execution.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("registry_id"); ok {
		paramMap["RegistryId"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("retention_id"); v != nil {
		paramMap["RetentionId"] = helper.IntInt64(v.(int))
	}

	service := TcrService{client: meta.(*TencentCloudClient).apiV3Conn}

	var retentionExecutionList []*tcr.RetentionExecution

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTcrTagRetentionExecutionByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		retentionExecutionList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(retentionExecutionList))
	tmpList := make([]map[string]interface{}, 0, len(retentionExecutionList))

	if retentionExecutionList != nil {
		for _, retentionExecution := range retentionExecutionList {
			retentionExecutionMap := map[string]interface{}{}

			if retentionExecution.ExecutionId != nil {
				retentionExecutionMap["execution_id"] = retentionExecution.ExecutionId
			}

			if retentionExecution.RetentionId != nil {
				retentionExecutionMap["retention_id"] = retentionExecution.RetentionId
			}

			if retentionExecution.StartTime != nil {
				retentionExecutionMap["start_time"] = retentionExecution.StartTime
			}

			if retentionExecution.EndTime != nil {
				retentionExecutionMap["end_time"] = retentionExecution.EndTime
			}

			if retentionExecution.Status != nil {
				retentionExecutionMap["status"] = retentionExecution.Status
			}

			ids = append(ids, *retentionExecution.RegistryId)
			tmpList = append(tmpList, retentionExecutionMap)
		}

		_ = d.Set("retention_execution_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
