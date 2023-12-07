package tencentcloud

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudWafAttackLogList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafAttackLogListRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain for query, all domain use all.",
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
			"query_count": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     10,
				Description: "Number of queries, default to 10, maximum of 100.",
			},
			"page": {
				Optional:    true,
				Type:        schema.TypeInt,
				Default:     0,
				Description: "Number of pages, starting from 0 by default.",
			},
			"query_string": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Lucene grammar.",
			},
			"sort": {
				Optional:    true,
				Type:        schema.TypeString,
				Default:     "desc",
				Description: "Default desc, support desc, asc.",
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
							Description: "Time string.",
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

func dataSourceTencentCloudWafAttackLogListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_waf_attack_log_list.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		ctx           = context.WithValue(context.TODO(), logIdKey, logId)
		service       = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
		attackLogList []*waf.AttackLogInfo
	)

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

	if v, ok := d.GetOkExists("query_count"); ok {
		paramMap["Count"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("query_string"); ok {
		paramMap["QueryString"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sort"); ok {
		paramMap["Sort"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("page"); ok {
		paramMap["Page"] = helper.IntInt64(v.(int))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafAttackLogListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		attackLogList = result
		return nil
	})

	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(attackLogList))

	if attackLogList != nil {
		for _, attackLogInfo := range attackLogList {
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

			tmpList = append(tmpList, attackLogInfoMap)
		}

		_ = d.Set("data", tmpList)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
