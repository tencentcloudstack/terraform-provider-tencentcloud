/*
Use this data source to query detailed information of mysql instance_param_record

Example Usage

```hcl
data "tencentcloud_mysql_instance_param_record" "instance_param_record" {
  instance_id = ""
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlInstanceParamRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceParamRecordRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Parameter modification record.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "instance id.",
						},
						"param_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "parameter name.",
						},
						"old_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The value of the parameter before modification.",
						},
						"new_value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The modified value of the parameter.",
						},
						"is_success": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the parameter is modified successfully.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Change the time.",
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

func dataSourceTencentCloudMysqlInstanceParamRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_mysql_instance_param_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var instanceParamRecord []*cdb.ParamRecord
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlInstanceParamRecordByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		instanceParamRecord = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(instanceParamRecord))
	tmpList := make([]map[string]interface{}, 0, len(instanceParamRecord))
	if instanceParamRecord != nil {
		for _, paramRecord := range instanceParamRecord {
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

			if paramRecord.IsSucess != nil {
				paramRecordMap["is_success"] = paramRecord.IsSucess
			}

			if paramRecord.ModifyTime != nil {
				paramRecordMap["modify_time"] = paramRecord.ModifyTime
			}

			ids = append(ids, *paramRecord.InstanceId)
			tmpList = append(tmpList, paramRecordMap)
		}

		_ = d.Set("items", tmpList)
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
