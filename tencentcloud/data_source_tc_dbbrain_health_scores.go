package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dbbrain "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbbrain/v20210527"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDbbrainHealthScores() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDbbrainHealthScoresRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The ID of the instance whose health score needs to be obtained.",
			},

			"time": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The time to obtain the health score, the time format is as follows: 2019-09-10 12:13:14.",
			},

			"product": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Service product type, supported values include: mysql - cloud database MySQL, cynosdb - cloud database TDSQL-C for MySQL, the default is mysql.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Health score and abnormal deduction items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issue_types": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Exception details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"issue_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Index classification: AVAILABILITY: availability, MAINTAINABILITY: maintainability, PERFORMANCE, performance, RELIABILITY reliability.",
									},
									"events": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "unusual event.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"event_id": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Event ID.",
												},
												"diag_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Diagnostic type.",
												},
												"start_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Starting time.",
												},
												"end_time": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "End Time.",
												},
												"outline": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "overview.",
												},
												"severity": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "severity. The severity is divided into 5 levels, according to the degree of impact from high to low: 1: Fatal, 2: Serious, 3: Warning, 4: Prompt, 5: Healthy.",
												},
												"score_lost": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Points deducted.",
												},
												"metric": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "reserved text.",
												},
												"count": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of alerts.",
												},
											},
										},
									},
									"total_count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The total number of abnormal events.",
									},
								},
							},
						},
						"events_total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total number of abnormal events.",
						},
						"health_score": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health score.",
						},
						"health_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Health level, such as: HEALTH, SUB_HEALTH, RISK, HIGH_RISK.",
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

func dataSourceTencentCloudDbbrainHealthScoresRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dbbrain_health_scores.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	var instanceId string
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["instance_id"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("time"); ok {
		paramMap["time"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("product"); ok {
		paramMap["product"] = helper.String(v.(string))
	}

	service := DbbrainService{client: meta.(*TencentCloudClient).apiV3Conn}

	var data *dbbrain.HealthScoreInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDbbrainHealthScoresByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		data = result
		return nil
	})
	if err != nil {
		return err
	}

	healthScoreInfoMap := map[string]interface{}{}
	if data != nil {
		if data.IssueTypes != nil {
			issueTypesList := []interface{}{}
			for _, issueTypes := range data.IssueTypes {
				issueTypesMap := map[string]interface{}{}

				if issueTypes.IssueType != nil {
					issueTypesMap["issue_type"] = issueTypes.IssueType
				}

				if issueTypes.Events != nil {
					eventsList := []interface{}{}
					for _, events := range issueTypes.Events {
						eventsMap := map[string]interface{}{}

						if events.EventId != nil {
							eventsMap["event_id"] = events.EventId
						}

						if events.DiagType != nil {
							eventsMap["diag_type"] = events.DiagType
						}

						if events.StartTime != nil {
							eventsMap["start_time"] = events.StartTime
						}

						if events.EndTime != nil {
							eventsMap["end_time"] = events.EndTime
						}

						if events.Outline != nil {
							eventsMap["outline"] = events.Outline
						}

						if events.Severity != nil {
							eventsMap["severity"] = events.Severity
						}

						if events.ScoreLost != nil {
							eventsMap["score_lost"] = events.ScoreLost
						}

						if events.Metric != nil {
							eventsMap["metric"] = events.Metric
						}

						if events.Count != nil {
							eventsMap["count"] = events.Count
						}

						eventsList = append(eventsList, eventsMap)
					}

					issueTypesMap["events"] = eventsList
				}

				if issueTypes.TotalCount != nil {
					issueTypesMap["total_count"] = issueTypes.TotalCount
				}

				issueTypesList = append(issueTypesList, issueTypesMap)
			}

			healthScoreInfoMap["issue_types"] = issueTypesList
		}

		if data.EventsTotalCount != nil {
			healthScoreInfoMap["events_total_count"] = data.EventsTotalCount
		}

		if data.HealthScore != nil {
			healthScoreInfoMap["health_score"] = data.HealthScore
		}

		if data.HealthLevel != nil {
			healthScoreInfoMap["health_level"] = data.HealthLevel
		}

		_ = d.Set("data", []interface{}{healthScoreInfoMap})
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), healthScoreInfoMap); e != nil {
			return e
		}
	}
	return nil
}
