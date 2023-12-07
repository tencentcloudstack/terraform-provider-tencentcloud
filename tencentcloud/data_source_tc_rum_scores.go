package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	rum "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/rum/v20210622"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudRumScores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRumScoresRead,
		Schema: map[string]*schema.Schema{
			"end_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "End time.",
			},

			"start_time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Start time.",
			},

			"project_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Project ID.",
			},

			"is_demo": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Get data from demo. This parameter is deprecated.",
			},

			"score_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Score list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"static_duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Duration.",
						},
						"page_pv": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pv.",
						},
						"api_fail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of failed api.",
						},
						"api_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of all request api.",
						},
						"static_fail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of failed request static resource.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"page_uv": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User view.",
						},
						"api_duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The mean duration of api request.",
						},
						"score": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The score of project.",
						},
						"page_error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of exception which happened on page.",
						},
						"static_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The number of static resource on page.",
						},
						"record_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of record.",
						},
						"page_duration": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The duration of page load.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project record created time.",
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

func dataSourceTencentCloudRumScoresRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_rum_scores.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var startTime, endTime string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("end_time"); ok {
		endTime = v.(string)
		paramMap["EndTime"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_time"); ok {
		startTime = v.(string)
		paramMap["StartTime"] = helper.String(v.(string))
	}

	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ID"] = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("is_demo"); v != nil {
		paramMap["IsDemo"] = helper.IntInt64(v.(int))
	}

	service := RumService{client: meta.(*TencentCloudClient).apiV3Conn}

	var scoreSet []*rum.ScoreInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRumScoresByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		scoreSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(scoreSet))
	tmpList := make([]map[string]interface{}, 0, len(scoreSet))

	if scoreSet != nil {
		for _, scoreInfo := range scoreSet {
			scoreInfoMap := map[string]interface{}{}

			if scoreInfo.StaticDuration != nil {
				scoreInfoMap["static_duration"] = scoreInfo.StaticDuration
			}

			if scoreInfo.PagePv != nil {
				scoreInfoMap["page_pv"] = scoreInfo.PagePv
			}

			if scoreInfo.ApiFail != nil {
				scoreInfoMap["api_fail"] = scoreInfo.ApiFail
			}

			if scoreInfo.ApiNum != nil {
				scoreInfoMap["api_num"] = scoreInfo.ApiNum
			}

			if scoreInfo.StaticFail != nil {
				scoreInfoMap["static_fail"] = scoreInfo.StaticFail
			}

			if scoreInfo.ProjectID != nil {
				scoreInfoMap["project_id"] = scoreInfo.ProjectID
			}

			if scoreInfo.PageUv != nil {
				scoreInfoMap["page_uv"] = scoreInfo.PageUv
			}

			if scoreInfo.ApiDuration != nil {
				scoreInfoMap["api_duration"] = scoreInfo.ApiDuration
			}

			if scoreInfo.Score != nil {
				scoreInfoMap["score"] = scoreInfo.Score
			}

			if scoreInfo.PageError != nil {
				scoreInfoMap["page_error"] = scoreInfo.PageError
			}

			if scoreInfo.StaticNum != nil {
				scoreInfoMap["static_num"] = scoreInfo.StaticNum
			}

			if scoreInfo.RecordNum != nil {
				scoreInfoMap["record_num"] = scoreInfo.RecordNum
			}

			if scoreInfo.PageDuration != nil {
				scoreInfoMap["page_duration"] = scoreInfo.PageDuration
			}

			if scoreInfo.CreateTime != nil {
				scoreInfoMap["create_time"] = scoreInfo.CreateTime
			}

			ids = append(ids, strconv.FormatInt(*scoreInfo.ProjectID, 10))
			tmpList = append(tmpList, scoreInfoMap)
		}

		_ = d.Set("score_set", tmpList)
	}

	ids = append(ids, startTime)
	ids = append(ids, endTime)
	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
