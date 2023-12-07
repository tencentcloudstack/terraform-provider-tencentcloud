package tencentcloud

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMonitorAlarmBasicAlarms() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMonitorAlarmBasicAlarmsRead,
		Schema: map[string]*schema.Schema{
			"module": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Interface module name, current value monitor.",
			},

			"start_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Start time, default to one day is timestamp.",
			},

			"end_time": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "End time, default to current timestamp.",
			},

			"occur_time_order": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by occurrence time, taking ASC or DESC values.",
			},

			"project_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter based on project ID.",
			},

			"view_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter based on policy type.",
			},

			"alarm_status": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter based on alarm status.",
			},

			"obj_like": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Filter based on alarm objects.",
			},

			"instance_group_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Description: "Filter based on instance group ID.",
			},

			"metric_names": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Filter by indicator name.",
			},

			"alarms": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Alarm List.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The ID of this alarm.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID.",
						},
						"project_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entry name.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm status ID, 0 indicates not recovered; 1 indicates that it has been restored; 2,3,5 indicates insufficient data; 4 indicates it has expired.",
						},
						"alarm_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm status, ALARM indicates not recovered; OK indicates that it has been restored; NO_ DATA indicates insufficient data; NO_ CONF indicates that it has expired.",
						},
						"group_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Policy Group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Group Name.",
						},
						"first_occur_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Time of occurrence.",
						},
						"duration": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Duration in seconds.",
						},
						"last_occur_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time.",
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm content.",
						},
						"obj_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm Object.",
						},
						"obj_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm object ID.",
						},
						"view_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Policy Type.",
						},
						"vpc": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC, only CVM has.",
						},
						"metric_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicator ID.",
						},
						"metric_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicator Name.",
						},
						"alarm_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Alarm type, 0 represents indicator alarm, 2 represents product event alarm, and 3 represents platform event alarm.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"dimensions": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Alarm object dimension information.",
						},
						"notify_way": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Notification method.",
						},
						"instance_group": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Instance Group Information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_group_id": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Instance Group ID.",
									},
									"instance_group_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance Group Name.",
									},
								},
							},
						},
					},
				},
			},

			"warning": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Remarks.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudMonitorAlarmBasicAlarmsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_monitor_alarm_basic_alarms.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("module"); ok {
		paramMap["Module"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("start_time"); ok {
		paramMap["StartTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_time"); ok {
		paramMap["EndTime"] = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("occur_time_order"); ok {
		paramMap["OccurTimeOrder"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsSet := v.(*schema.Set).List()
		paramMap["ProjectIds"] = helper.InterfacesIntInt64Point(projectIdsSet)
	}

	if v, ok := d.GetOk("view_names"); ok {
		viewNamesSet := v.(*schema.Set).List()
		paramMap["ViewNames"] = helper.InterfacesStringsPoint(viewNamesSet)
	}

	if v, ok := d.GetOk("alarm_status"); ok {
		alarmStatusSet := v.(*schema.Set).List()
		paramMap["AlarmStatus"] = helper.InterfacesIntInt64Point(alarmStatusSet)
	}

	if v, ok := d.GetOk("obj_like"); ok {
		paramMap["ObjLike"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_group_ids"); ok {
		instanceGroupIdsSet := v.(*schema.Set).List()
		paramMap["InstanceGroupIds"] = helper.InterfacesIntInt64Point(instanceGroupIdsSet)
	}

	if v, ok := d.GetOk("metric_names"); ok {
		metricNamesSet := v.(*schema.Set).List()
		paramMap["MetricNames"] = helper.InterfacesStringsPoint(metricNamesSet)
	}

	service := MonitorService{client: meta.(*TencentCloudClient).apiV3Conn}

	var alarms []*monitor.DescribeBasicAlarmListAlarms
	var warning *string
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, w, e := service.DescribeMonitorAlarmBasicAlarmsByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		alarms = result
		warning = w
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(alarms))
	tmpList := make([]map[string]interface{}, 0, len(alarms))

	if alarms != nil {
		for _, describeBasicAlarmListAlarms := range alarms {
			describeBasicAlarmListAlarmsMap := map[string]interface{}{}

			if describeBasicAlarmListAlarms.Id != nil {
				describeBasicAlarmListAlarmsMap["id"] = describeBasicAlarmListAlarms.Id
			}

			if describeBasicAlarmListAlarms.ProjectId != nil {
				describeBasicAlarmListAlarmsMap["project_id"] = describeBasicAlarmListAlarms.ProjectId
			}

			if describeBasicAlarmListAlarms.ProjectName != nil {
				describeBasicAlarmListAlarmsMap["project_name"] = describeBasicAlarmListAlarms.ProjectName
			}

			if describeBasicAlarmListAlarms.Status != nil {
				describeBasicAlarmListAlarmsMap["status"] = describeBasicAlarmListAlarms.Status
			}

			if describeBasicAlarmListAlarms.AlarmStatus != nil {
				describeBasicAlarmListAlarmsMap["alarm_status"] = describeBasicAlarmListAlarms.AlarmStatus
			}

			if describeBasicAlarmListAlarms.GroupId != nil {
				describeBasicAlarmListAlarmsMap["group_id"] = describeBasicAlarmListAlarms.GroupId
			}

			if describeBasicAlarmListAlarms.GroupName != nil {
				describeBasicAlarmListAlarmsMap["group_name"] = describeBasicAlarmListAlarms.GroupName
			}

			if describeBasicAlarmListAlarms.FirstOccurTime != nil {
				describeBasicAlarmListAlarmsMap["first_occur_time"] = describeBasicAlarmListAlarms.FirstOccurTime
			}

			if describeBasicAlarmListAlarms.Duration != nil {
				describeBasicAlarmListAlarmsMap["duration"] = describeBasicAlarmListAlarms.Duration
			}

			if describeBasicAlarmListAlarms.LastOccurTime != nil {
				describeBasicAlarmListAlarmsMap["last_occur_time"] = describeBasicAlarmListAlarms.LastOccurTime
			}

			if describeBasicAlarmListAlarms.Content != nil {
				describeBasicAlarmListAlarmsMap["content"] = describeBasicAlarmListAlarms.Content
			}

			if describeBasicAlarmListAlarms.ObjName != nil {
				describeBasicAlarmListAlarmsMap["obj_name"] = describeBasicAlarmListAlarms.ObjName
			}

			if describeBasicAlarmListAlarms.ObjId != nil {
				describeBasicAlarmListAlarmsMap["obj_id"] = describeBasicAlarmListAlarms.ObjId
			}

			if describeBasicAlarmListAlarms.ViewName != nil {
				describeBasicAlarmListAlarmsMap["view_name"] = describeBasicAlarmListAlarms.ViewName
			}

			if describeBasicAlarmListAlarms.Vpc != nil {
				describeBasicAlarmListAlarmsMap["vpc"] = describeBasicAlarmListAlarms.Vpc
			}

			if describeBasicAlarmListAlarms.MetricId != nil {
				describeBasicAlarmListAlarmsMap["metric_id"] = describeBasicAlarmListAlarms.MetricId
			}

			if describeBasicAlarmListAlarms.MetricName != nil {
				describeBasicAlarmListAlarmsMap["metric_name"] = describeBasicAlarmListAlarms.MetricName
			}

			if describeBasicAlarmListAlarms.AlarmType != nil {
				describeBasicAlarmListAlarmsMap["alarm_type"] = describeBasicAlarmListAlarms.AlarmType
			}

			if describeBasicAlarmListAlarms.Region != nil {
				describeBasicAlarmListAlarmsMap["region"] = describeBasicAlarmListAlarms.Region
			}

			if describeBasicAlarmListAlarms.Dimensions != nil {
				describeBasicAlarmListAlarmsMap["dimensions"] = describeBasicAlarmListAlarms.Dimensions
			}

			if describeBasicAlarmListAlarms.NotifyWay != nil {
				describeBasicAlarmListAlarmsMap["notify_way"] = describeBasicAlarmListAlarms.NotifyWay
			}

			if describeBasicAlarmListAlarms.InstanceGroup != nil {
				instanceGroupList := []interface{}{}
				for _, instanceGroup := range describeBasicAlarmListAlarms.InstanceGroup {
					instanceGroupMap := map[string]interface{}{}

					if instanceGroup.InstanceGroupId != nil {
						instanceGroupMap["instance_group_id"] = instanceGroup.InstanceGroupId
					}

					if instanceGroup.InstanceGroupName != nil {
						instanceGroupMap["instance_group_name"] = instanceGroup.InstanceGroupName
					}

					instanceGroupList = append(instanceGroupList, instanceGroupMap)
				}

				describeBasicAlarmListAlarmsMap["instance_group"] = instanceGroupList
			}

			ids = append(ids, strconv.Itoa(int(*describeBasicAlarmListAlarms.Id)))
			tmpList = append(tmpList, describeBasicAlarmListAlarmsMap)
		}

		_ = d.Set("alarms", tmpList)
	}

	if warning != nil {
		_ = d.Set("warning", warning)
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
