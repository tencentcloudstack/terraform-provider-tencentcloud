/*
Use this data source to query detailed information of sqlserver instance_param_records

Example Usage

```hcl
data "tencentcloud_sqlserver_instance_param_records" "instance_param_records" {
  instance_id = "mssql-qelbzgwf"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverInstanceParamRecords() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSqlserverInstanceParamRecordsRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of mssql-dj5i29c5n. It is the same as the instance ID displayed in the TencentDB console and the response parameter InstanceId of the DescribeDBInstances API.",
			},
			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter modification records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter name.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter value before modification.",
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Parameter value after modification.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Parameter modification status. Valid values: 1 (initializing and waiting for modification), 2 (modification succeed), 3 (modification failed), 4 (modifying).",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.",
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

func dataSourceTencentCloudSqlserverInstanceParamRecordsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_instance_param_records.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	var items []*sqlserver.ParamRecord

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSqlserverInstanceParamRecordsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		items = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))

	if items != nil {
		for _, paramRecord := range items {
			paramRecordMap := map[string]interface{}{}

			if paramRecord.InstanceId != nil {
				paramRecordMap["instance_id"] = paramRecord.InstanceId
			}

			if paramRecord.ParamName != nil {
				paramRecordMap["param_name"] = paramRecord.ParamName
			}

			if paramRecord.OldValue != nil {
				paramRecordMap["old_value"] = paramRecord.OldValue
			}

			if paramRecord.NewValue != nil {
				paramRecordMap["new_value"] = paramRecord.NewValue
			}

			if paramRecord.Status != nil {
				paramRecordMap["status"] = paramRecord.Status
			}

			if paramRecord.ModifyTime != nil {
				paramRecordMap["modify_time"] = paramRecord.ModifyTime
			}

			tmpList = append(tmpList, paramRecordMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
