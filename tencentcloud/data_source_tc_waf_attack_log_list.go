/*
Use this data source to query detailed information of waf attack_log_list

Example Usage

```hcl
data "tencentcloud_waf_attack_log_list" "attack_log_list" {
  domain = ""
  start_time = ""
  end_time = ""
    query_string = ""
  sort = ""
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

func dataSourceTencentCloudWafAttackLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackLogListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain for query , all domain use &amp;amp;#39;all&amp;amp;#39;.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Begin time.",
			},

			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"context": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Write empty string.",
			},

			"query_string": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lucene grammar.",
			},

			"sort": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Default &amp;amp;#39;desc&amp;amp;#39; , &amp;amp;#39;desc&amp;amp;#39; or &amp;amp;#39;asc&amp;amp;#39;.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Attack log array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The detail of attack log.",
						},
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Useless.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Useless.",
						},
						"time_stamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Timestring.",
						},
					},
				},
			},

			"list_over": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Useless.",
			},

			"sql_flag": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Useless.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafAttackLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_attack_log_list.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_time"); ok {
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		paramMap["Sort"] = helper.String(v.(string))
	}

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackLogListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		context = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(context))
	if context != nil {
		_ = d.Set("context", context)
	}

	if data != nil {
		for _, attackLogInfo := range data {
			attackLogInfoMap := map[string]interface{}{}

			if attackLogInfo.Content != nil {
				attackLogInfoMap["content"] = attackLogInfo.Content
			}

			if attackLogInfo.FileName != nil {
				attackLogInfoMap["file_name"] = attackLogInfo.FileName
			}

			if attackLogInfo.Source != nil {
				attackLogInfoMap["source"] = attackLogInfo.Source
			}

			if attackLogInfo.TimeStamp != nil {
				attackLogInfoMap["time_stamp"] = attackLogInfo.TimeStamp
			}

			ids = append(ids, *attackLogInfo.uuid)
			tmpList = append(tmpList, attackLogInfoMap)
		}

		_ = d.Set("data", tmpList)
	}

	if listOver != nil {
		_ = d.Set("list_over", listOver)
	}

	if sqlFlag != nil {
		_ = d.Set("sql_flag", sqlFlag)
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
