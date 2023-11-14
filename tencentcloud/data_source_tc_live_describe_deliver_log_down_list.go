/*
Use this data source to query detailed information of live describe_deliver_log_down_list

Example Usage

```hcl
data "tencentcloud_live_describe_deliver_log_down_list" "describe_deliver_log_down_list" {
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	live "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLiveDescribeDeliverLogDownList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLiveDescribeDeliverLogDownListRead,
		Schema: map[string]*schema.Schema{
			"log_info_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of log information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"log_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log name.",
						},
						"log_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log download address.",
						},
						"log_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Log time. UTC format, for example: 2018-11-29T19:00:00Z.Note:Beijing time is UTC time + 8 hours, formatted according to the ISO 8601 standard, see ISO date format description for details.",
						},
						"file_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "File size, in bytes.",
						},
					},
				},
			},

			"total_num": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Total count.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudLiveDescribeDeliverLogDownListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_live_describe_deliver_log_down_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	service := LiveService{client: meta.(*TencentCloudClient).apiV3Conn}

	var logInfoList []*live.PushLogInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLiveDescribeDeliverLogDownListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		logInfoList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(logInfoList))
	tmpList := make([]map[string]interface{}, 0, len(logInfoList))

	if logInfoList != nil {
		for _, pushLogInfo := range logInfoList {
			pushLogInfoMap := map[string]interface{}{}

			if pushLogInfo.LogName != nil {
				pushLogInfoMap["log_name"] = pushLogInfo.LogName
			}

			if pushLogInfo.LogUrl != nil {
				pushLogInfoMap["log_url"] = pushLogInfo.LogUrl
			}

			if pushLogInfo.LogTime != nil {
				pushLogInfoMap["log_time"] = pushLogInfo.LogTime
			}

			if pushLogInfo.FileSize != nil {
				pushLogInfoMap["file_size"] = pushLogInfo.FileSize
			}

			ids = append(ids, *pushLogInfo.LogName)
			tmpList = append(tmpList, pushLogInfoMap)
		}

		_ = d.Set("log_info_list", tmpList)
	}

	if totalNum != nil {
		_ = d.Set("total_num", totalNum)
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
