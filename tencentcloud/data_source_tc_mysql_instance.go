package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentCloudMysqlInstanceDetail() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"mysql_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_role": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"init_flag": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"zone": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"auto_renew_flag": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"engine_version": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cpu_core_count": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"memory_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"internet_status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"internet_domain": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"internet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"intranet_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"project_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"slave_sync_mode": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"device_type": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pay_type": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"create_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"dead_line_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"master_instance_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ro_instance_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"dr_instance_ids": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"master", "ro", "dr"}),
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 4, 5}),
			},
			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pay_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7"}),
			},
			"init_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"with_dr": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"with_ro": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"with_master": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			},
			"offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 1000),
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: validateIntegerInRange(1, 2000),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values
			"instance_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mysql_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_role": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"init_flag": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auto_renew_flag": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_core_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"volume_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"internet_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"intranet_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"intranet_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_sync_mode": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"device_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pay_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"dead_line_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ro_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dr_instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_mysql_instance.read")()

	logId := GetLogId(nil)

	request := cdb.NewDescribeDBInstancesRequest()

	if mysqlId, ok := d.GetOk("mysql_id"); ok {
		mysqlIdValue := mysqlId.(string)
		request.InstanceIds = []*string{&mysqlIdValue}
	}
	if instanceType, ok := d.GetOk("instance_role"); ok {
		instanceTypeValue := instanceType.(string)
		var instanceRole uint64 = 1
		for k, v := range MYSQL_ROLE_MAP {
			if instanceTypeValue == v {
				instanceRole = uint64(k)
			}
		}
		request.InstanceTypes = []*uint64{&instanceRole}
	}
	if status, ok := d.GetOk("status"); ok {
		statusValue := uint64(status.(int))
		request.Status = []*uint64{&statusValue}
	}
	if securityGroupId, ok := d.GetOk("security_group_id"); ok {
		securityGroupIdValue := securityGroupId.(string)
		request.SecurityGroupId = &securityGroupIdValue
	}
	if payType, ok := d.GetOk("pay_type"); ok {
		payTypeValue := uint64(payType.(int))
		request.PayTypes = []*uint64{&payTypeValue}
	}
	if instanceName, ok := d.GetOk("instance_name"); ok {
		instanceNameValue := instanceName.(string)
		request.InstanceNames = []*string{&instanceNameValue}
	}
	if taskStatus, ok := d.GetOk("task_status"); ok {
		taskStatusValue := uint64(taskStatus.(int))
		request.TaskStatus = []*uint64{&taskStatusValue}
	}
	if engineVersion, ok := d.GetOk("engine_version"); ok {
		engineVersionValue := engineVersion.(string)
		request.EngineVersions = []*string{&engineVersionValue}
	}
	if initFlag, ok := d.GetOk("init_flag"); ok {
		initFlagValue := int64(initFlag.(int))
		request.InitFlag = &initFlagValue
	}
	if withDr, ok := d.GetOk("with_dr"); ok {
		withDrValue := int64(withDr.(int))
		request.WithDr = &withDrValue
	}
	if withRo, ok := d.GetOk("with_ro"); ok {
		withRoValue := int64(withRo.(int))
		request.WithRo = &withRoValue
	}
	if withMaster, ok := d.GetOk("with_master"); ok {
		withMasterValue := int64(withMaster.(int))
		request.WithMaster = &withMasterValue
	}
	offset := d.Get("offset")
	offsetValue := uint64(offset.(int))
	request.Offset = &offsetValue
	limit := d.Get("limit")
	limitValue := uint64(limit.(int))
	request.Limit = &limitValue

	client := meta.(*TencentCloudClient).apiV3Conn
	response, err := client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return fmt.Errorf("api[DescribeDBInstances]fail, return %s", err.Error())
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceDetails := response.Response.Items
	instanceList := make([]map[string]interface{}, 0, len(instanceDetails))
	ids := make([]string, 0, len(instanceDetails))
	for _, item := range instanceDetails {
		mapping := map[string]interface{}{
			"mysql_id":        *item.InstanceId,
			"instance_name":   *item.InstanceName,
			"instance_role":   MYSQL_ROLE_MAP[*item.InstanceType],
			"init_flag":       *item.InitFlag,
			"status":          *item.Status,
			"zone":            *item.Zone,
			"auto_renew_flag": *item.AutoRenew,
			"engine_version":  *item.EngineVersion,
			"cpu_core_count":  *item.Cpu,
			"memory_size":     *item.Memory,
			"volume_size":     *item.Volume,
			"internet_status": *item.WanStatus,
			"internet_domain": *item.WanDomain,
			"internet_port":   *item.WanPort,
			"intranet_ip":     *item.Vip,
			"intranet_port":   *item.Vport,
			"project_id":      *item.ProjectId,
			"vpc_id":          *item.UniqVpcId,
			"subnet_id":       *item.UniqSubnetId,
			"slave_sync_mode": *item.ProtectMode,
			"device_type":     *item.DeviceType,
			"pay_type":        *item.PayType,
			"create_time":     *item.CreateTime,
			"dead_line_time":  *item.DeadlineTime,
		}
		if item.MasterInfo != nil {
			mapping["master_instance_id"] = *item.MasterInfo.InstanceId
		} else {
			mapping["master_instance_id"] = ""
		}
		if len(item.RoGroups) > 0 {
			roInstanceIds := make([]string, 0)
			for _, roGroupInfo := range item.RoGroups {
				for _, roInfo := range roGroupInfo.RoInstances {
					roInstanceIds = append(roInstanceIds, *roInfo.InstanceId)
				}
			}
			mapping["ro_instance_ids"] = roInstanceIds
		}
		if len(item.DrInfo) > 0 {
			drInstanceIds := make([]string, 0)
			for _, drInfo := range item.DrInfo {
				drInstanceIds = append(drInstanceIds, *drInfo.InstanceId)
			}
			mapping["dr_instance_ids"] = drInstanceIds
		}

		ids = append(ids, *item.InstanceId)
		instanceList = append(instanceList, mapping)
	}
	d.SetId(dataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		writeToFile(output.(string), instanceList)
	}
	return nil
}
