package monitor

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudMonitorAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmPolicyRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Value fixed at monitor.",
			},

			"policy_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Fuzzy search by policy name.",
			},

			"monitor_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by monitor type. Valid values: MT_QCE (Tencent Cloud service monitoring). If this parameter is left empty, all will be queried by default.",
			},

			"namespaces": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by namespace. For the values of different policy types, please see:[Poicy Type List](https://www.tencentcloud.com/document/product/248/39565?has_map=1).",
			},

			"dimensions": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The alarm object list, which is a JSON string. The outer array corresponds to multiple instances, and the inner array is the dimension of an object.For example, 'CVM - Basic Monitor' can be written as: [ {Dimensions: {unInstanceId: ins-qr8d555g}}, {Dimensions: {unInstanceId: ins-qr8d555h}} ]You can also refer to the 'Example 2' below.For more information on the parameter samples of different Tencent Cloud services, see [Product Policy Type and Dimension Information](https://www.tencentcloud.com/document/product/248/39565?has_map=1).Note: If 1 is passed in for NeedCorrespondence, the relationship between a policy and an instance needs to be returned. You can pass in up to 20 alarm object dimensions to avoid request timeout.",
			},

			"receiver_uids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search by recipient. You can get the user list with the API [ListUsers](https://www.tencentcloud.com/document/product/598/34587?from_cn_redirect=1) in 'Cloud Access Management' or query the sub-user information with the API [GetUser](https://www.tencentcloud.com/document/product/598/34590?from_cn_redirect=1). The Uid field in the returned result should be entered here.",
			},

			"receiver_groups": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Search by recipient group. You can get the user group list with the API [ListGroups](https://www.tencentcloud.com/document/product/598/34589?from_cn_redirect=1) in 'Cloud Access Management' or query the user group list where a sub-user is in with the API [ListGroupsForUser](https://www.tencentcloud.com/document/product/598/34588?from_cn_redirect=1). The GroupId field in the returned result should be entered here.",
			},

			"policy_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by default policy. Valid values: DEFAULT (display default policy), NOT_DEFAULT (display non-default policies). If this parameter is left empty, all policies will be displayed.",
			},

			"field": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by field. For example, to sort by the last modification time, use Field: UpdateTime.",
			},

			"order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort order. Valid values: ASC (ascending), DESC (descending).",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "ID array of the policy project, which can be viewed on the following page: [Project Management](https://console.tencentcloud.com/project).",
			},

			"notice_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of the notification template IDs, which can be obtained by querying the notification template list.It can be queried with the API [DescribeAlarmNotices](https://www.tencentcloud.com/document/product/248/39300).",
			},

			"rule_types": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by trigger condition. Valid values: STATIC (display policies with static threshold), DYNAMIC (display policies with dynamic threshold). If this parameter is left empty, all policies will be displayed.",
			},

			"enable": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter by alarm status. Valid values: [1]: enabled; [0]: disabled; [0, 1]: all.",
			},

			"not_binding_notice_rule": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "If 1 is passed in, alarm policies with no notification rules configured are queried. If it is left empty or other values are passed in, all alarm policies are queried.",
			},

			"instance_group_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Instance group ID.",
			},

			"need_correspondence": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether the relationship between a policy and the input parameter filter dimension is required. 1: Yes. 0: No. Default value: 0.",
			},

			"trigger_tasks": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter alarm policy by triggered task (such as auto scaling task). Up to 10 tasks can be specified.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Triggered task type. Valid value: AS (auto scaling)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"task_config": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Configuration information in JSON format, such as {Key1:Value1,Key2:Value2}Note: this field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},

			"one_click_policy_type": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by quick alarm policy. If this parameter is left empty, all policies are displayed. ONECLICK: Display quick alarm policies; NOT_ONECLICK: Display non-quick alarm policies.",
			},

			"not_bind_all": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether the returned result needs to filter policies associated with all objects. Valid values: 1 (Yes), 0 (No).",
			},

			"not_instance_group": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether the returned result needs to filter policies associated with instance groups. Valid values: 1 (Yes), 0 (No).",
			},

			"prom_ins_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "ID of the TencentCloud Managed Service for Prometheus instance, which is used for customizing a metric policy.",
			},

			"receiver_on_call_form_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Search by schedule.",
			},

			"policies": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Policy array.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm policy IDNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm policy nameNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "RemarksNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"monitor_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Monitor type. Valid values: MT_QCE (Tencent Cloud service monitoring)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status. Valid values: 0 (disabled), 1 (enabled)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"use_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of instances bound to policy groupNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID. Valid values: -1 (no project), 0 (default project)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Project nameNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm policy typeNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"condition_template_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger condition template IDNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Metric trigger conditionNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"is_union_rule": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Judgment condition of an alarm trigger condition (0: Any; 1: All; 2: Composite). When the value is set to 2 (i.e., composite trigger conditions), this parameter should be used together with ComplexExpression.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric name or event name. The supported metrics can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322) and the supported events via [DescribeAlarmEvents](https://www.tencentcloud.com/document/product/248/39324).Note: this field may return null, indicating that no valid value is obtained.",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
												},
												"operator": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Threshold. The valid value range can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.",
												},
												"continue_period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via [DescribeAlarmMetrics](https://www.tencentcloud.com/document/product/248/39322)Note: this field may return null, indicating that no valid value is obtained.",
												},
												"notice_frequency": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every dayNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"is_power_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yesNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"filter": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Filter condition for one single trigger rulNote: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Filter condition type. Valid values: DIMENSION (uses dimensions for filteringNote: this field may return null, indicating that no valid values can be obtained.",
															},
															"dimensions": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshiNote: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric display name, which is used in the output parameteNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unit, which is used in the output parameteNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by defaultNote: this field may return null, indicating that no valid value is obtained.",
												},
												"is_advanced": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether it is an advanced metric. 0: No; 1: YesNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"is_open": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether the advanced metric feature is enabled. 0: No; 1: YesNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"product_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Integration center product IDNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"value_max": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Maximum valuNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"value_min": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Minimum valuNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"hierarchical_value": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The configuration of alarm level thresholNote: This field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"remind": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Remind leveNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"warn": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Warn leveNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"serious": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Serious leveNote: This field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
											},
										},
									},
									"complex_expression": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The judgment expression of composite alarm trigger conditions, which is valid when the value of IsUnionRule is 2. This parameter is used to determine that an alarm condition is met only when the expression values are True for multiple trigger conditionsNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"event_condition": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Event trigger conditioNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"rules": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Alarm trigger condition lisNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"metric_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEventsNote: this field may return null, indicating that no valid value is obtained.",
												},
												"period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Statistical period in seconds. The valid values can be queried via DescribeAlarmMetricsNote: this field may return null, indicating that no valid value is obtained.",
												},
												"operator": {
													Type:     schema.TypeString,
													Computed: true,
													Description: "Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.Operator	String	No	Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
												},
												"continue_period": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
												},
												"notice_frequency": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"is_power_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"filter": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"dimensions": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
												"description": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"unit": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"rule_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.",
												},
												"is_advanced": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"is_open": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"product_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"value_max": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Maximum valueNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"value_min": {
													Type:        schema.TypeFloat,
													Computed:    true,
													Description: "Minimum valueNote: This field may return null, indicating that no valid values can be obtained.",
												},
												"hierarchical_value": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"remind": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"warn": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"serious": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.",
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"notice_ids": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Notification rule ID listNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"notices": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Notification rule listNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm notification template IDNote: this field may return null, indicating that no valid values can be obtained.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm notification template nameNote: this field may return null, indicating that no valid values can be obtained.",
									},
									"updated_at": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified timeNote: this field may return null, indicating that no valid values can be obtained.",
									},
									"updated_by": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Last modified byNote: this field may return null, indicating that no valid values can be obtained.",
									},
									"notice_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Alarm notification type. Valid values: ALARM (for unresolved alarms), OK (for resolved alarms), ALL (for all alarms)Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"user_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "User notification listNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"receiver_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Recipient type. Valid values: USER (user), GROUP (user group)Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Notification start time, which is expressed by the number of seconds since 00:00:00. Value range: 0-86399Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"end_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Notification end time, which is expressed by the number of seconds since 00:00:00. Value range: 0-86399Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"notice_way": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "Notification channel list. Valid values: EMAIL (email), SMS (SMS), CALL (phone), WECHAT (WeChat), RTX (WeCom)Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"user_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "User uid listNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"group_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "User group ID listNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"phone_order": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Phone polling listNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"phone_circle_times": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Number of phone pollings. Value range: 1-5Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"phone_inner_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Call interval in seconds within one polling. Value range: 60-900Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"phone_circle_interval": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Polling interval in seconds. Value range: 60-900Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"need_phone_arrive_notice": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether receipt notification is required. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"phone_call_type": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Dial type. SYNC (simultaneous dial), CIRCLE (polled dial). Default value: CIRCLE.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"weekday": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Notification cycle. The values 1-7 indicate Monday to Sunday.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"on_call_form_ids": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
													Computed:    true,
													Description: "List of schedule IDsNote: u200dThis field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"url_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Callback notification listNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"url": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Callback URL, which can contain up to 256 charactersNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"is_valid": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Whether verification is passed. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"validation_code": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Verification codeNote: this field may return null, indicating that no valid values can be obtained.",
												},
												"start_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Start time of the notification in seconds, which is calculated from 00:00:00.Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"end_time": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "End time of the notification in seconds, which is calculated from 00:00:00.Note: this field may return null, indicating that no valid values can be obtained.",
												},
												"weekday": {
													Type: schema.TypeSet,
													Elem: &schema.Schema{
														Type: schema.TypeInt,
													},
													Computed:    true,
													Description: "Notification cycle. The values 1-7 indicate Monday to Sunday.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"is_preset": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Whether it is the system default notification template. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"notice_language": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Notification language. Valid values: zh-CN (Chinese), en-US (English)Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"policy_ids": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "List of IDs of the alarm policies bound to alarm notification templateNote: this field may return null, indicating that no valid values can be obtained.",
									},
									"amp_consumer_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Backend AMP consumer ID.Note: This field may return null, indicating that no valid values can be obtained.",
									},
									"cls_notices": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Channel to push alarm notifications to CLS.Note: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"region": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Region.",
												},
												"log_set_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Logset ID.",
												},
												"topic_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Topic ID.",
												},
												"enable": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Status. Valid values: 0 (disabled), 1 (enabled). Default value: 1 (enabled). This parameter can be left empty.",
												},
											},
										},
									},
									"tags": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Tags bound to a notification templateNote: This field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tag key.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Tag value.",
												},
											},
										},
									},
								},
							},
						},
						"trigger_tasks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Triggered task listNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Triggered task type. Valid value: AS (auto scaling)Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"task_config": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Configuration information in JSON format, such as {Key1:Value1,Key2:Value2}Note: this field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"conditions_temp": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Template policy groupNote: this field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"template_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Template nameNote: u200dThis field may return null, indicating that no valid values can be obtained.",
									},
									"condition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Metric trigger conditionNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"is_union_rule": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Judgment condition of an alarm trigger condition (0: Any; 1: All; 2: Composite). When the value is set to 2 (i.e., composite trigger conditions), this parameter should be used together with ComplexExpression.Note: This field may return null, indicating that no valid values can be obtained.",
												},
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEvents.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"continue_period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"notice_frequency": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"is_power_notice": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"filter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.",
																		},
																		"dimensions": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.",
																		},
																	},
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
															},
															"unit": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
															},
															"rule_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"is_advanced": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"is_open": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"product_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"value_max": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Maximum valueNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"value_min": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Minimum valueNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"hierarchical_value": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"remind": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"warn": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"serious": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																	},
																},
															},
														},
													},
												},
												"complex_expression": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The judgment expression of composite alarm trigger conditions, which is valid when the value of IsUnionRule is 2. This parameter is used to determine that an alarm condition is met only when the expression values are True for multiple trigger conditions.Note: This field may return null, indicating that no valid values can be obtained.",
												},
											},
										},
									},
									"event_condition": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Event trigger conditionNote: this field may return null, indicating that no valid values can be obtained.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"rules": {
													Type:        schema.TypeList,
													Computed:    true,
													Description: "Alarm trigger condition listNote: this field may return null, indicating that no valid values can be obtained.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"metric_name": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric name or event name. The supported metrics can be queried via DescribeAlarmMetrics and the supported events via DescribeAlarmEvents.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Statistical period in seconds. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"operator": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Operatorintelligent = intelligent detection without thresholdeq = equal toge = greater than or equal togt = greater thanle = less than or equal tolt = less thanne = not equal today_increase = day-on-day increaseday_decrease = day-on-day decreaseday_wave = day-on-day fluctuationweek_increase = week-on-week increaseweek_decrease = week-on-week decreaseweek_wave = week-on-week fluctuationcycle_increase = cyclical increasecycle_decrease = cyclical decreasecycle_wave = cyclical fluctuationre = regex matchThe valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"value": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Threshold. The valid value range can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"continue_period": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Number of periods. 1: continue for one period; 2: continue for two periods; and so on. The valid values can be queried via DescribeAlarmMetrics.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"notice_frequency": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Alarm interval in seconds. Valid values: 0 (do not repeat), 300 (alarm once every 5 minutes), 600 (alarm once every 10 minutes), 900 (alarm once every 15 minutes), 1800 (alarm once every 30 minutes), 3600 (alarm once every hour), 7200 (alarm once every 2 hours), 10800 (alarm once every 3 hours), 21600 (alarm once every 6 hours), 43200 (alarm once every 12 hours), 86400 (alarm once every day)Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"is_power_notice": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether the alarm frequency increases exponentially. Valid values: 0 (no), 1 (yes)Note: this field may return null, indicating that no valid values can be obtained.",
															},
															"filter": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "Filter condition for one single trigger ruleNote: this field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"type": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Filter condition type. Valid values: DIMENSION (uses dimensions for filtering)Note: this field may return null, indicating that no valid values can be obtained.",
																		},
																		"dimensions": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "JSON string generated by serializing the AlarmPolicyDimension two-dimensional array. The one-dimensional arrays are in OR relationship, and the elements in a one-dimensional array are in AND relationshipNote: this field may return null, indicating that no valid values can be obtained.",
																		},
																	},
																},
															},
															"description": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Metric display name, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
															},
															"unit": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Unit, which is used in the output parameterNote: this field may return null, indicating that no valid values can be obtained.",
															},
															"rule_type": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Trigger condition type. STATIC: static threshold; dynamic: dynamic threshold. If you do not specify this parameter when creating or editing a policy, STATIC is used by default.Note: this field may return null, indicating that no valid value is obtained.",
															},
															"is_advanced": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether it is an advanced metric. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"is_open": {
																Type:        schema.TypeInt,
																Computed:    true,
																Description: "Whether the advanced metric feature is enabled. 0: No; 1: Yes.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"product_id": {
																Type:        schema.TypeString,
																Computed:    true,
																Description: "Integration center product ID.Note: This field may return null, indicating that no valid values can be obtained.",
															},
															"value_max": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Maximum valueNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"value_min": {
																Type:        schema.TypeFloat,
																Computed:    true,
																Description: "Minimum valueNote: This field may return null, indicating that no valid values can be obtained.",
															},
															"hierarchical_value": {
																Type:        schema.TypeList,
																Computed:    true,
																Description: "The configuration of alarm level thresholdNote: This field may return null, indicating that no valid values can be obtained.",
																Elem: &schema.Resource{
																	Schema: map[string]*schema.Schema{
																		"remind": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Remind levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"warn": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Warn levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																		"serious": {
																			Type:        schema.TypeString,
																			Computed:    true,
																			Description: "Threshold for the Serious levelNote: This field may return null, indicating that no valid values can be obtained.",
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						"last_edit_uin": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Uin of the last modifying userNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"update_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Update timeNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"insert_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation timeNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"region": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "RegionNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"namespace_show_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Namespace display nameNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"is_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is the default policy. Valid values: 1 (yes), 0 (no)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"can_set_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the default policy can be set. Valid values: 1 (yes), 0 (no)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"instance_group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance group IDNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"instance_sum": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total number of instances in instance groupNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"instance_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance group nameNote: this field may return null, indicating that no valid values can be obtained.",
						},
						"rule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Trigger condition type. Valid values: STATIC (static threshold), DYNAMIC (dynamic)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"origin_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy ID for instance/instance group binding and unbinding APIs (BindingPolicyObject, UnBindingAllPolicyObject, UnBindingPolicyObject)Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"tag_instances": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "TagNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag keyNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag valueNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"instance_sum": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Number of instancesNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"service_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service type, for example, CVMNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"region_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Region IDNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"binding_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Binding status. 2: bound; 1: bindingNote: This field may return null, indicating that no valid values can be obtained.",
									},
									"tag_status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Tag status. 2: existent; 1: nonexistentNote: This field may return null, indicating that no valid values can be obtained.",
									},
								},
							},
						},
						"filter_dimensions_param": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information on the filter dimension associated with a policy.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_one_click": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether it is a quick alarm policy.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"one_click_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the quick alarm policy is enabled.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"advanced_metric_number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of advanced metrics.Note: This field may return null, indicating that no valid values can be obtained.",
						},
						"is_bind_all": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether the policy is associated with all objectsNote: This field may return null, indicating that no valid values can be obtained.",
						},
						"tags": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Policy tagNote: This field may return null, indicating that no valid values can be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag key.",
									},
									"value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Tag value.",
									},
								},
							},
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

func dataSourceTencentCloudMonitorAlarmPolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_monitor_alarm_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("policy_name"); ok {
		paramMap["PolicyName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("monitor_types"); ok {
		monitorTypesSet := v.(*schema.Set).List()
		paramMap["MonitorTypes"] = helper.InterfacesStringsPoint(monitorTypesSet)
	}

	if v, ok := d.GetOk("namespaces"); ok {
		namespacesSet := v.(*schema.Set).List()
		paramMap["Namespaces"] = helper.InterfacesStringsPoint(namespacesSet)
	}

	if v, ok := d.GetOk("dimensions"); ok {
		paramMap["Dimensions"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiver_uids"); ok {
		receiverUidsList := []*int64{}
		receiverUidsSet := v.(*schema.Set).List()
		for i := range receiverUidsSet {
			receiverUids := receiverUidsSet[i].(int)
			receiverUidsList = append(receiverUidsList, helper.IntInt64(receiverUids))
		}
		paramMap["ReceiverUids"] = receiverUidsList
	}

	if v, ok := d.GetOk("receiver_groups"); ok {
		receiverGroupsList := []*int64{}
		receiverGroupsSet := v.(*schema.Set).List()
		for i := range receiverGroupsSet {
			receiverGroups := receiverGroupsSet[i].(int)
			receiverGroupsList = append(receiverGroupsList, helper.IntInt64(receiverGroups))
		}
		paramMap["ReceiverGroups"] = receiverGroupsList
	}

	if v, ok := d.GetOk("policy_type"); ok {
		policyTypeSet := v.(*schema.Set).List()
		paramMap["PolicyType"] = helper.InterfacesStringsPoint(policyTypeSet)
	}

	if v, ok := d.GetOk("field"); ok {
		paramMap["Field"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		paramMap["Order"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsList := []*int64{}
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(int)
			projectIdsList = append(projectIdsList, helper.IntInt64(projectIds))
		}
		paramMap["ProjectIds"] = projectIdsList
	}

	if v, ok := d.GetOk("notice_ids"); ok {
		noticeIdsSet := v.(*schema.Set).List()
		paramMap["NoticeIds"] = helper.InterfacesStringsPoint(noticeIdsSet)
	}

	if v, ok := d.GetOk("rule_types"); ok {
		ruleTypesSet := v.(*schema.Set).List()
		paramMap["RuleTypes"] = helper.InterfacesStringsPoint(ruleTypesSet)
	}

	if v, ok := d.GetOk("enable"); ok {
		enableList := []*int64{}
		enableSet := v.(*schema.Set).List()
		for i := range enableSet {
			enable := enableSet[i].(int)
			enableList = append(enableList, helper.IntInt64(enable))
		}
		paramMap["Enable"] = enableList
	}

	if v, ok := d.GetOkExists("not_binding_notice_rule"); ok {
		paramMap["NotBindingNoticeRule"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("instance_group_id"); ok {
		paramMap["InstanceGroupId"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("need_correspondence"); ok {
		paramMap["NeedCorrespondence"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("trigger_tasks"); ok {
		triggerTasksSet := v.([]interface{})
		tmpSet := make([]*monitor.AlarmPolicyTriggerTask, 0, len(triggerTasksSet))

		for _, item := range triggerTasksSet {
			alarmPolicyTriggerTask := monitor.AlarmPolicyTriggerTask{}
			alarmPolicyTriggerTaskMap := item.(map[string]interface{})

			if v, ok := alarmPolicyTriggerTaskMap["type"]; ok {
				alarmPolicyTriggerTask.Type = helper.String(v.(string))
			}
			if v, ok := alarmPolicyTriggerTaskMap["task_config"]; ok {
				alarmPolicyTriggerTask.TaskConfig = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &alarmPolicyTriggerTask)
		}
		paramMap["trigger_tasks"] = tmpSet
	}

	if v, ok := d.GetOk("one_click_policy_type"); ok {
		oneClickPolicyTypeSet := v.(*schema.Set).List()
		paramMap["OneClickPolicyType"] = helper.InterfacesStringsPoint(oneClickPolicyTypeSet)
	}

	if v, ok := d.GetOkExists("not_bind_all"); ok {
		paramMap["NotBindAll"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("not_instance_group"); ok {
		paramMap["NotInstanceGroup"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("prom_ins_id"); ok {
		paramMap["PromInsId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("receiver_on_call_form_ids"); ok {
		receiverOnCallFormIDsSet := v.(*schema.Set).List()
		paramMap["ReceiverOnCallFormIDs"] = helper.InterfacesStringsPoint(receiverOnCallFormIDsSet)
	}

	service := MonitorService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var policies []*monitor.AlarmPolicy
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMonitorAlarmPolicyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		policies = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(policies))
	tmpList := make([]map[string]interface{}, 0, len(policies))

	if policies != nil {
		for _, alarmPolicy := range policies {
			alarmPolicyMap := map[string]interface{}{}

			if alarmPolicy.PolicyId != nil {
				alarmPolicyMap["policy_id"] = alarmPolicy.PolicyId
			}

			if alarmPolicy.PolicyName != nil {
				alarmPolicyMap["policy_name"] = alarmPolicy.PolicyName
			}

			if alarmPolicy.Remark != nil {
				alarmPolicyMap["remark"] = alarmPolicy.Remark
			}

			if alarmPolicy.MonitorType != nil {
				alarmPolicyMap["monitor_type"] = alarmPolicy.MonitorType
			}

			if alarmPolicy.Enable != nil {
				alarmPolicyMap["enable"] = alarmPolicy.Enable
			}

			if alarmPolicy.UseSum != nil {
				alarmPolicyMap["use_sum"] = alarmPolicy.UseSum
			}

			if alarmPolicy.ProjectId != nil {
				alarmPolicyMap["project_id"] = alarmPolicy.ProjectId
			}

			if alarmPolicy.ProjectName != nil {
				alarmPolicyMap["project_name"] = alarmPolicy.ProjectName
			}

			if alarmPolicy.Namespace != nil {
				alarmPolicyMap["namespace"] = alarmPolicy.Namespace
			}

			if alarmPolicy.ConditionTemplateId != nil {
				alarmPolicyMap["condition_template_id"] = alarmPolicy.ConditionTemplateId
			}

			if alarmPolicy.Condition != nil {
				conditionMap := map[string]interface{}{}

				if alarmPolicy.Condition.IsUnionRule != nil {
					conditionMap["is_union_rule"] = alarmPolicy.Condition.IsUnionRule
				}

				if alarmPolicy.Condition.Rules != nil {
					rulesList := []interface{}{}
					for _, rules := range alarmPolicy.Condition.Rules {
						rulesMap := map[string]interface{}{}

						if rules.MetricName != nil {
							rulesMap["metric_name"] = rules.MetricName
						}

						if rules.Period != nil {
							rulesMap["period"] = rules.Period
						}

						if rules.Operator != nil {
							rulesMap["operator"] = rules.Operator
						}

						if rules.Value != nil {
							rulesMap["value"] = rules.Value
						}

						if rules.ContinuePeriod != nil {
							rulesMap["continue_period"] = rules.ContinuePeriod
						}

						if rules.NoticeFrequency != nil {
							rulesMap["notice_frequency"] = rules.NoticeFrequency
						}

						if rules.IsPowerNotice != nil {
							rulesMap["is_power_notice"] = rules.IsPowerNotice
						}

						if rules.Filter != nil {
							filterMap := map[string]interface{}{}

							if rules.Filter.Type != nil {
								filterMap["type"] = rules.Filter.Type
							}

							if rules.Filter.Dimensions != nil {
								filterMap["dimensions"] = rules.Filter.Dimensions
							}

							rulesMap["filter"] = []interface{}{filterMap}
						}

						if rules.Description != nil {
							rulesMap["description"] = rules.Description
						}

						if rules.Unit != nil {
							rulesMap["unit"] = rules.Unit
						}

						if rules.RuleType != nil {
							rulesMap["rule_type"] = rules.RuleType
						}

						if rules.IsAdvanced != nil {
							rulesMap["is_advanced"] = rules.IsAdvanced
						}

						if rules.IsOpen != nil {
							rulesMap["is_open"] = rules.IsOpen
						}

						if rules.ProductId != nil {
							rulesMap["product_id"] = rules.ProductId
						}

						if rules.ValueMax != nil {
							rulesMap["value_max"] = rules.ValueMax
						}

						if rules.ValueMin != nil {
							rulesMap["value_min"] = rules.ValueMin
						}

						if rules.HierarchicalValue != nil {
							hierarchicalValueMap := map[string]interface{}{}

							if rules.HierarchicalValue.Remind != nil {
								hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
							}

							if rules.HierarchicalValue.Warn != nil {
								hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
							}

							if rules.HierarchicalValue.Serious != nil {
								hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
							}

							rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
						}

						rulesList = append(rulesList, rulesMap)
					}

					conditionMap["rules"] = rulesList
				}

				if alarmPolicy.Condition.ComplexExpression != nil {
					conditionMap["complex_expression"] = alarmPolicy.Condition.ComplexExpression
				}

				alarmPolicyMap["condition"] = []interface{}{conditionMap}
			}

			if alarmPolicy.EventCondition != nil {
				eventConditionMap := map[string]interface{}{}

				if alarmPolicy.EventCondition.Rules != nil {
					rulesList := []interface{}{}
					for _, rules := range alarmPolicy.EventCondition.Rules {
						rulesMap := map[string]interface{}{}

						if rules.MetricName != nil {
							rulesMap["metric_name"] = rules.MetricName
						}

						if rules.Period != nil {
							rulesMap["period"] = rules.Period
						}

						if rules.Operator != nil {
							rulesMap["operator"] = rules.Operator
						}

						if rules.Value != nil {
							rulesMap["value"] = rules.Value
						}

						if rules.ContinuePeriod != nil {
							rulesMap["continue_period"] = rules.ContinuePeriod
						}

						if rules.NoticeFrequency != nil {
							rulesMap["notice_frequency"] = rules.NoticeFrequency
						}

						if rules.IsPowerNotice != nil {
							rulesMap["is_power_notice"] = rules.IsPowerNotice
						}

						if rules.Filter != nil {
							filterMap := map[string]interface{}{}

							if rules.Filter.Type != nil {
								filterMap["type"] = rules.Filter.Type
							}

							if rules.Filter.Dimensions != nil {
								filterMap["dimensions"] = rules.Filter.Dimensions
							}

							rulesMap["filter"] = []interface{}{filterMap}
						}

						if rules.Description != nil {
							rulesMap["description"] = rules.Description
						}

						if rules.Unit != nil {
							rulesMap["unit"] = rules.Unit
						}

						if rules.RuleType != nil {
							rulesMap["rule_type"] = rules.RuleType
						}

						if rules.IsAdvanced != nil {
							rulesMap["is_advanced"] = rules.IsAdvanced
						}

						if rules.IsOpen != nil {
							rulesMap["is_open"] = rules.IsOpen
						}

						if rules.ProductId != nil {
							rulesMap["product_id"] = rules.ProductId
						}

						if rules.ValueMax != nil {
							rulesMap["value_max"] = rules.ValueMax
						}

						if rules.ValueMin != nil {
							rulesMap["value_min"] = rules.ValueMin
						}

						if rules.HierarchicalValue != nil {
							hierarchicalValueMap := map[string]interface{}{}

							if rules.HierarchicalValue.Remind != nil {
								hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
							}

							if rules.HierarchicalValue.Warn != nil {
								hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
							}

							if rules.HierarchicalValue.Serious != nil {
								hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
							}

							rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
						}

						rulesList = append(rulesList, rulesMap)
					}

					eventConditionMap["rules"] = rulesList
				}

				alarmPolicyMap["event_condition"] = []interface{}{eventConditionMap}
			}

			if alarmPolicy.NoticeIds != nil {
				alarmPolicyMap["notice_ids"] = alarmPolicy.NoticeIds
			}

			if alarmPolicy.Notices != nil {
				noticesList := []interface{}{}
				for _, notices := range alarmPolicy.Notices {
					noticesMap := map[string]interface{}{}

					if notices.Id != nil {
						noticesMap["id"] = notices.Id
					}

					if notices.Name != nil {
						noticesMap["name"] = notices.Name
					}

					if notices.UpdatedAt != nil {
						noticesMap["updated_at"] = notices.UpdatedAt
					}

					if notices.UpdatedBy != nil {
						noticesMap["updated_by"] = notices.UpdatedBy
					}

					if notices.NoticeType != nil {
						noticesMap["notice_type"] = notices.NoticeType
					}

					if notices.UserNotices != nil {
						userNoticesList := []interface{}{}
						for _, userNotices := range notices.UserNotices {
							userNoticesMap := map[string]interface{}{}

							if userNotices.ReceiverType != nil {
								userNoticesMap["receiver_type"] = userNotices.ReceiverType
							}

							if userNotices.StartTime != nil {
								userNoticesMap["start_time"] = userNotices.StartTime
							}

							if userNotices.EndTime != nil {
								userNoticesMap["end_time"] = userNotices.EndTime
							}

							if userNotices.NoticeWay != nil {
								userNoticesMap["notice_way"] = userNotices.NoticeWay
							}

							if userNotices.UserIds != nil {
								userNoticesMap["user_ids"] = userNotices.UserIds
							}

							if userNotices.GroupIds != nil {
								userNoticesMap["group_ids"] = userNotices.GroupIds
							}

							if userNotices.PhoneOrder != nil {
								userNoticesMap["phone_order"] = userNotices.PhoneOrder
							}

							if userNotices.PhoneCircleTimes != nil {
								userNoticesMap["phone_circle_times"] = userNotices.PhoneCircleTimes
							}

							if userNotices.PhoneInnerInterval != nil {
								userNoticesMap["phone_inner_interval"] = userNotices.PhoneInnerInterval
							}

							if userNotices.PhoneCircleInterval != nil {
								userNoticesMap["phone_circle_interval"] = userNotices.PhoneCircleInterval
							}

							if userNotices.NeedPhoneArriveNotice != nil {
								userNoticesMap["need_phone_arrive_notice"] = userNotices.NeedPhoneArriveNotice
							}

							if userNotices.PhoneCallType != nil {
								userNoticesMap["phone_call_type"] = userNotices.PhoneCallType
							}

							if userNotices.Weekday != nil {
								userNoticesMap["weekday"] = userNotices.Weekday
							}

							if userNotices.OnCallFormIDs != nil {
								userNoticesMap["on_call_form_ids"] = userNotices.OnCallFormIDs
							}

							userNoticesList = append(userNoticesList, userNoticesMap)
						}

						noticesMap["user_notices"] = userNoticesList
					}

					if notices.URLNotices != nil {
						uRLNoticesList := []interface{}{}
						for _, uRLNotices := range notices.URLNotices {
							uRLNoticesMap := map[string]interface{}{}

							if uRLNotices.URL != nil {
								uRLNoticesMap["url"] = uRLNotices.URL
							}

							if uRLNotices.IsValid != nil {
								uRLNoticesMap["is_valid"] = uRLNotices.IsValid
							}

							if uRLNotices.ValidationCode != nil {
								uRLNoticesMap["validation_code"] = uRLNotices.ValidationCode
							}

							if uRLNotices.StartTime != nil {
								uRLNoticesMap["start_time"] = uRLNotices.StartTime
							}

							if uRLNotices.EndTime != nil {
								uRLNoticesMap["end_time"] = uRLNotices.EndTime
							}

							if uRLNotices.Weekday != nil {
								uRLNoticesMap["weekday"] = uRLNotices.Weekday
							}

							uRLNoticesList = append(uRLNoticesList, uRLNoticesMap)
						}

						noticesMap["url_notices"] = uRLNoticesList
					}

					if notices.IsPreset != nil {
						noticesMap["is_preset"] = notices.IsPreset
					}

					if notices.NoticeLanguage != nil {
						noticesMap["notice_language"] = notices.NoticeLanguage
					}

					if notices.PolicyIds != nil {
						noticesMap["policy_ids"] = notices.PolicyIds
					}

					if notices.AMPConsumerId != nil {
						noticesMap["amp_consumer_id"] = notices.AMPConsumerId
					}

					if notices.CLSNotices != nil {
						cLSNoticesList := []interface{}{}
						for _, cLSNotices := range notices.CLSNotices {
							cLSNoticesMap := map[string]interface{}{}

							if cLSNotices.Region != nil {
								cLSNoticesMap["region"] = cLSNotices.Region
							}

							if cLSNotices.LogSetId != nil {
								cLSNoticesMap["log_set_id"] = cLSNotices.LogSetId
							}

							if cLSNotices.TopicId != nil {
								cLSNoticesMap["topic_id"] = cLSNotices.TopicId
							}

							if cLSNotices.Enable != nil {
								cLSNoticesMap["enable"] = cLSNotices.Enable
							}

							cLSNoticesList = append(cLSNoticesList, cLSNoticesMap)
						}

						noticesMap["cls_notices"] = cLSNoticesList
					}

					if notices.Tags != nil {
						tagsList := []interface{}{}
						for _, tags := range notices.Tags {
							tagsMap := map[string]interface{}{}

							if tags.Key != nil {
								tagsMap["key"] = tags.Key
							}

							if tags.Value != nil {
								tagsMap["value"] = tags.Value
							}

							tagsList = append(tagsList, tagsMap)
						}

						noticesMap["tags"] = tagsList
					}

					noticesList = append(noticesList, noticesMap)
				}

				alarmPolicyMap["notices"] = noticesList
			}

			if alarmPolicy.TriggerTasks != nil {
				triggerTasksList := []interface{}{}
				for _, triggerTasks := range alarmPolicy.TriggerTasks {
					triggerTasksMap := map[string]interface{}{}

					if triggerTasks.Type != nil {
						triggerTasksMap["type"] = triggerTasks.Type
					}

					if triggerTasks.TaskConfig != nil {
						triggerTasksMap["task_config"] = triggerTasks.TaskConfig
					}

					triggerTasksList = append(triggerTasksList, triggerTasksMap)
				}

				alarmPolicyMap["trigger_tasks"] = triggerTasksList
			}

			if alarmPolicy.ConditionsTemp != nil {
				conditionsTempMap := map[string]interface{}{}

				if alarmPolicy.ConditionsTemp.TemplateName != nil {
					conditionsTempMap["template_name"] = alarmPolicy.ConditionsTemp.TemplateName
				}

				if alarmPolicy.ConditionsTemp.Condition != nil {
					conditionMap := map[string]interface{}{}

					if alarmPolicy.ConditionsTemp.Condition.IsUnionRule != nil {
						conditionMap["is_union_rule"] = alarmPolicy.ConditionsTemp.Condition.IsUnionRule
					}

					if alarmPolicy.ConditionsTemp.Condition.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range alarmPolicy.ConditionsTemp.Condition.Rules {
							rulesMap := map[string]interface{}{}

							if rules.MetricName != nil {
								rulesMap["metric_name"] = rules.MetricName
							}

							if rules.Period != nil {
								rulesMap["period"] = rules.Period
							}

							if rules.Operator != nil {
								rulesMap["operator"] = rules.Operator
							}

							if rules.Value != nil {
								rulesMap["value"] = rules.Value
							}

							if rules.ContinuePeriod != nil {
								rulesMap["continue_period"] = rules.ContinuePeriod
							}

							if rules.NoticeFrequency != nil {
								rulesMap["notice_frequency"] = rules.NoticeFrequency
							}

							if rules.IsPowerNotice != nil {
								rulesMap["is_power_notice"] = rules.IsPowerNotice
							}

							if rules.Filter != nil {
								filterMap := map[string]interface{}{}

								if rules.Filter.Type != nil {
									filterMap["type"] = rules.Filter.Type
								}

								if rules.Filter.Dimensions != nil {
									filterMap["dimensions"] = rules.Filter.Dimensions
								}

								rulesMap["filter"] = []interface{}{filterMap}
							}

							if rules.Description != nil {
								rulesMap["description"] = rules.Description
							}

							if rules.Unit != nil {
								rulesMap["unit"] = rules.Unit
							}

							if rules.RuleType != nil {
								rulesMap["rule_type"] = rules.RuleType
							}

							if rules.IsAdvanced != nil {
								rulesMap["is_advanced"] = rules.IsAdvanced
							}

							if rules.IsOpen != nil {
								rulesMap["is_open"] = rules.IsOpen
							}

							if rules.ProductId != nil {
								rulesMap["product_id"] = rules.ProductId
							}

							if rules.ValueMax != nil {
								rulesMap["value_max"] = rules.ValueMax
							}

							if rules.ValueMin != nil {
								rulesMap["value_min"] = rules.ValueMin
							}

							if rules.HierarchicalValue != nil {
								hierarchicalValueMap := map[string]interface{}{}

								if rules.HierarchicalValue.Remind != nil {
									hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
								}

								if rules.HierarchicalValue.Warn != nil {
									hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
								}

								if rules.HierarchicalValue.Serious != nil {
									hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
								}

								rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
							}

							rulesList = append(rulesList, rulesMap)
						}

						conditionMap["rules"] = rulesList
					}

					if alarmPolicy.ConditionsTemp.Condition.ComplexExpression != nil {
						conditionMap["complex_expression"] = alarmPolicy.ConditionsTemp.Condition.ComplexExpression
					}

					conditionsTempMap["condition"] = []interface{}{conditionMap}
				}

				if alarmPolicy.ConditionsTemp.EventCondition != nil {
					eventConditionMap := map[string]interface{}{}

					if alarmPolicy.ConditionsTemp.EventCondition.Rules != nil {
						rulesList := []interface{}{}
						for _, rules := range alarmPolicy.ConditionsTemp.EventCondition.Rules {
							rulesMap := map[string]interface{}{}

							if rules.MetricName != nil {
								rulesMap["metric_name"] = rules.MetricName
							}

							if rules.Period != nil {
								rulesMap["period"] = rules.Period
							}

							if rules.Operator != nil {
								rulesMap["operator"] = rules.Operator
							}

							if rules.Value != nil {
								rulesMap["value"] = rules.Value
							}

							if rules.ContinuePeriod != nil {
								rulesMap["continue_period"] = rules.ContinuePeriod
							}

							if rules.NoticeFrequency != nil {
								rulesMap["notice_frequency"] = rules.NoticeFrequency
							}

							if rules.IsPowerNotice != nil {
								rulesMap["is_power_notice"] = rules.IsPowerNotice
							}

							if rules.Filter != nil {
								filterMap := map[string]interface{}{}

								if rules.Filter.Type != nil {
									filterMap["type"] = rules.Filter.Type
								}

								if rules.Filter.Dimensions != nil {
									filterMap["dimensions"] = rules.Filter.Dimensions
								}

								rulesMap["filter"] = []interface{}{filterMap}
							}

							if rules.Description != nil {
								rulesMap["description"] = rules.Description
							}

							if rules.Unit != nil {
								rulesMap["unit"] = rules.Unit
							}

							if rules.RuleType != nil {
								rulesMap["rule_type"] = rules.RuleType
							}

							if rules.IsAdvanced != nil {
								rulesMap["is_advanced"] = rules.IsAdvanced
							}

							if rules.IsOpen != nil {
								rulesMap["is_open"] = rules.IsOpen
							}

							if rules.ProductId != nil {
								rulesMap["product_id"] = rules.ProductId
							}

							if rules.ValueMax != nil {
								rulesMap["value_max"] = rules.ValueMax
							}

							if rules.ValueMin != nil {
								rulesMap["value_min"] = rules.ValueMin
							}

							if rules.HierarchicalValue != nil {
								hierarchicalValueMap := map[string]interface{}{}

								if rules.HierarchicalValue.Remind != nil {
									hierarchicalValueMap["remind"] = rules.HierarchicalValue.Remind
								}

								if rules.HierarchicalValue.Warn != nil {
									hierarchicalValueMap["warn"] = rules.HierarchicalValue.Warn
								}

								if rules.HierarchicalValue.Serious != nil {
									hierarchicalValueMap["serious"] = rules.HierarchicalValue.Serious
								}

								rulesMap["hierarchical_value"] = []interface{}{hierarchicalValueMap}
							}

							rulesList = append(rulesList, rulesMap)
						}

						eventConditionMap["rules"] = rulesList
					}

					conditionsTempMap["event_condition"] = []interface{}{eventConditionMap}
				}

				alarmPolicyMap["conditions_temp"] = []interface{}{conditionsTempMap}
			}

			if alarmPolicy.LastEditUin != nil {
				alarmPolicyMap["last_edit_uin"] = alarmPolicy.LastEditUin
			}

			if alarmPolicy.UpdateTime != nil {
				alarmPolicyMap["update_time"] = alarmPolicy.UpdateTime
			}

			if alarmPolicy.InsertTime != nil {
				alarmPolicyMap["insert_time"] = alarmPolicy.InsertTime
			}

			if alarmPolicy.Region != nil {
				alarmPolicyMap["region"] = alarmPolicy.Region
			}

			if alarmPolicy.NamespaceShowName != nil {
				alarmPolicyMap["namespace_show_name"] = alarmPolicy.NamespaceShowName
			}

			if alarmPolicy.IsDefault != nil {
				alarmPolicyMap["is_default"] = alarmPolicy.IsDefault
			}

			if alarmPolicy.CanSetDefault != nil {
				alarmPolicyMap["can_set_default"] = alarmPolicy.CanSetDefault
			}

			if alarmPolicy.InstanceGroupId != nil {
				alarmPolicyMap["instance_group_id"] = alarmPolicy.InstanceGroupId
			}

			if alarmPolicy.InstanceSum != nil {
				alarmPolicyMap["instance_sum"] = alarmPolicy.InstanceSum
			}

			if alarmPolicy.InstanceGroupName != nil {
				alarmPolicyMap["instance_group_name"] = alarmPolicy.InstanceGroupName
			}

			if alarmPolicy.RuleType != nil {
				alarmPolicyMap["rule_type"] = alarmPolicy.RuleType
			}

			if alarmPolicy.OriginId != nil {
				alarmPolicyMap["origin_id"] = alarmPolicy.OriginId
			}

			if alarmPolicy.TagInstances != nil {
				tagInstancesList := []interface{}{}
				for _, tagInstances := range alarmPolicy.TagInstances {
					tagInstancesMap := map[string]interface{}{}

					if tagInstances.Key != nil {
						tagInstancesMap["key"] = tagInstances.Key
					}

					if tagInstances.Value != nil {
						tagInstancesMap["value"] = tagInstances.Value
					}

					if tagInstances.InstanceSum != nil {
						tagInstancesMap["instance_sum"] = tagInstances.InstanceSum
					}

					if tagInstances.ServiceType != nil {
						tagInstancesMap["service_type"] = tagInstances.ServiceType
					}

					if tagInstances.RegionId != nil {
						tagInstancesMap["region_id"] = tagInstances.RegionId
					}

					if tagInstances.BindingStatus != nil {
						tagInstancesMap["binding_status"] = tagInstances.BindingStatus
					}

					if tagInstances.TagStatus != nil {
						tagInstancesMap["tag_status"] = tagInstances.TagStatus
					}

					tagInstancesList = append(tagInstancesList, tagInstancesMap)
				}

				alarmPolicyMap["tag_instances"] = tagInstancesList
			}

			if alarmPolicy.FilterDimensionsParam != nil {
				alarmPolicyMap["filter_dimensions_param"] = alarmPolicy.FilterDimensionsParam
			}

			if alarmPolicy.IsOneClick != nil {
				alarmPolicyMap["is_one_click"] = alarmPolicy.IsOneClick
			}

			if alarmPolicy.OneClickStatus != nil {
				alarmPolicyMap["one_click_status"] = alarmPolicy.OneClickStatus
			}

			if alarmPolicy.AdvancedMetricNumber != nil {
				alarmPolicyMap["advanced_metric_number"] = alarmPolicy.AdvancedMetricNumber
			}

			if alarmPolicy.IsBindAll != nil {
				alarmPolicyMap["is_bind_all"] = alarmPolicy.IsBindAll
			}

			if alarmPolicy.Tags != nil {
				tagsList := []interface{}{}
				for _, tags := range alarmPolicy.Tags {
					tagsMap := map[string]interface{}{}

					if tags.Key != nil {
						tagsMap["key"] = tags.Key
					}

					if tags.Value != nil {
						tagsMap["value"] = tags.Value
					}

					tagsList = append(tagsList, tagsMap)
				}

				alarmPolicyMap["tags"] = tagsList
			}

			ids = append(ids, *alarmPolicy.PolicyId)
			tmpList = append(tmpList, alarmPolicyMap)
		}

		_ = d.Set("policies", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
