package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmRepairTasks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmRepairTasksRead,

		Schema: map[string]*schema.Schema{
			"product": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Product type, optional values: CVM (Cloud Virtual Machine), CDH (Cloud Dedicated Host), CPM2.0 (Cloud Physical Machine 2.0).",
			},
			"task_status": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Task status list. Valid values: 1 (pending authorization), 2 (processing), 3 (ended), 4 (scheduled), 5 (cancelled), 6 (avoided).",
			},
			"task_type_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Task type ID list. Valid values: 101 (instance running hazard), 102 (instance running exception), 103 (instance hard disk exception), 104 (instance network connection exception), 105 (instance running warning), 106 (instance hard disk warning), 107 (instance maintenance upgrade).",
			},
			"task_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Task ID list (e.g., rep-xxxxxxxx). Query specific tasks by task IDs.",
			},
			"instance_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Instance ID list (e.g., ins-xxxxxxxx). Query tasks by instance IDs.",
			},
			"aliases": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Instance name list (alias). Query tasks by instance names.",
			},
			"start_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query start date, format: YYYY-MM-DD hh:mm:ss. Filter tasks created from this date.",
			},
			"end_date": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Query end date, format: YYYY-MM-DD hh:mm:ss. Filter tasks created until this date.",
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting field. Valid values: CreateTime (creation time), AuthTime (authorization time), EndTime (end time).",
			},
			"order": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Sorting order. 0: ascending, 1: descending. Default: 0.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed
			"repair_task_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of repair tasks. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"task_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task ID.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name (alias).",
						},
						"task_type_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task type ID.",
						},
						"task_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Task status.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task creation time.",
						},
						"auth_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task authorization time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task end time.",
						},
						"task_detail": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task detail description.",
						},
						"task_type_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task type name.",
						},
						"task_sub_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Task sub-type.",
						},
						"device_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Device status.",
						},
						"operate_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Operation status.",
						},
						"auth_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Authorization type.",
						},
						"auth_source": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Authorization source.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Availability zone.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC ID.",
						},
						"vpc_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "VPC name.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID.",
						},
						"subnet_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name.",
						},
						"wan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public IP address.",
						},
						"lan_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Private IP address.",
						},
						"product": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Product type.",
						},
					},
				},
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total count of repair tasks that match the filter conditions.",
			},
		},
	}
}

func dataSourceTencentCloudCvmRepairTasksRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_repair_tasks.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = cvm.NewDescribeTaskInfoRequest()
		allRepairTasks []*cvm.RepairTaskInfo
	)

	// Read input parameters
	if v, ok := d.GetOk("product"); ok {
		request.Product = helper.String(v.(string))
	}

	if v, ok := d.GetOk("task_status"); ok {
		taskStatusSet := v.(*schema.Set).List()
		for _, status := range taskStatusSet {
			request.TaskStatus = append(request.TaskStatus, helper.IntInt64(status.(int)))
		}
	}

	if v, ok := d.GetOk("task_type_ids"); ok {
		taskTypeIdsSet := v.(*schema.Set).List()
		for _, typeId := range taskTypeIdsSet {
			request.TaskTypeIds = append(request.TaskTypeIds, helper.IntInt64(typeId.(int)))
		}
	}

	if v, ok := d.GetOk("task_ids"); ok {
		taskIdsSet := v.(*schema.Set).List()
		for _, taskId := range taskIdsSet {
			request.TaskIds = append(request.TaskIds, helper.String(taskId.(string)))
		}
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for _, instanceId := range instanceIdsSet {
			request.InstanceIds = append(request.InstanceIds, helper.String(instanceId.(string)))
		}
	}

	if v, ok := d.GetOk("aliases"); ok {
		aliasesSet := v.(*schema.Set).List()
		for _, alias := range aliasesSet {
			request.Aliases = append(request.Aliases, helper.String(alias.(string)))
		}
	}

	if v, ok := d.GetOk("start_date"); ok {
		request.StartDate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		request.EndDate = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order_field"); ok {
		request.OrderField = helper.String(v.(string))
	}

	if v, ok := d.GetOk("order"); ok {
		request.Order = helper.IntInt64(v.(int))
	}

	// Automatic pagination: fetch all data
	var offset int64 = 0
	var limit int64 = 100
	request.Limit = &limit

	for {
		request.Offset = &offset

		var response *cvm.DescribeTaskInfoResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeTaskInfo(request)
			if e != nil {
				return tccommon.RetryError(e, tccommon.InternalError)
			}
			response = result
			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s read CVM repair tasks failed, reason:%+v", logId, err)
			return err
		}

		if response == nil || response.Response == nil {
			break
		}

		repairTasks := response.Response.RepairTaskInfoSet
		allRepairTasks = append(allRepairTasks, repairTasks...)

		// Check if we need to fetch more data
		if len(repairTasks) < int(limit) {
			break
		}

		offset += limit
	}

	// Process and set output data
	ids := make([]string, 0, len(allRepairTasks))
	repairTaskList := make([]map[string]interface{}, 0, len(allRepairTasks))
	for _, task := range allRepairTasks {
		taskMap := map[string]interface{}{
			"task_id":        task.TaskId,
			"instance_id":    task.InstanceId,
			"alias":          task.Alias,
			"task_type_id":   task.TaskTypeId,
			"task_status":    task.TaskStatus,
			"create_time":    task.CreateTime,
			"auth_time":      task.AuthTime,
			"end_time":       task.EndTime,
			"task_detail":    task.TaskDetail,
			"task_type_name": task.TaskTypeName,
			"task_sub_type":  task.TaskSubType,
			"device_status":  task.DeviceStatus,
			"operate_status": task.OperateStatus,
			"auth_type":      task.AuthType,
			"auth_source":    task.AuthSource,
			"zone":           task.Zone,
			"region":         task.Region,
			"vpc_id":         task.VpcId,
			"vpc_name":       task.VpcName,
			"subnet_id":      task.SubnetId,
			"subnet_name":    task.SubnetName,
			"wan_ip":         task.WanIp,
			"lan_ip":         task.LanIp,
			"product":        task.Product,
		}
		repairTaskList = append(repairTaskList, taskMap)
		if task.TaskId != nil {
			ids = append(ids, *task.TaskId)
		}
	}

	_ = d.Set("repair_task_list", repairTaskList)
	_ = d.Set("total_count", len(allRepairTasks))

	d.SetId(helper.DataResourceIdsHash(ids))

	// Export to file if specified
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), repairTaskList); err != nil {
			return err
		}
	}

	return nil
}
