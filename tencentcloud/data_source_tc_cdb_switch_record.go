/*
Use this data source to query detailed information of cdb switch_record

Example Usage

```hcl
data "tencentcloud_cdb_switch_record" "switch_record" {
  instance_id = ""
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

func dataSourceTencentCloudCdbSwitchRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCdbSwitchRecordRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv or cdbro-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance switching record details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switching time, the format is: 2017-09-03 01:34:31.",
						},
						"switch_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch type, possible return values: TRANSFER - data migration; MASTER2SLAVE - master-standby switch; RECOVERY - master-slave recovery.",
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

func dataSourceTencentCloudCdbSwitchRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cdb_switch_record.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
	}

	service := CdbService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdbSwitchRecordByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		totalCount = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(totalCount))
	if items != nil {
		for _, dBSwitchInfo := range items {
			dBSwitchInfoMap := map[string]interface{}{}

			if dBSwitchInfo.SwitchTime != nil {
				dBSwitchInfoMap["switch_time"] = dBSwitchInfo.SwitchTime
			}

			if dBSwitchInfo.SwitchType != nil {
				dBSwitchInfoMap["switch_type"] = dBSwitchInfo.SwitchType
			}

			ids = append(ids, *dBSwitchInfo.InstanceId)
			tmpList = append(tmpList, dBSwitchInfoMap)
		}

		_ = d.Set("items", tmpList)
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
