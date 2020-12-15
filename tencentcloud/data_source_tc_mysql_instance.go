/*
Use this data source to get information about a MySQL instance.

Example Usage

```hcl
data "tencentcloud_mysql_instance" "database" {
  mysql_id           = "my-test-database"
  result_output_file = "mytestpath"
}
```
*/
package tencentcloud

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudMysqlInstance() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlInstanceRead,
		Schema: map[string]*schema.Schema{
			"mysql_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance ID, such as `cdb-c1nl9rpv`. It is identical to the instance ID displayed in the database console page.",
			},
			"instance_role": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"master", "ro", "dr"}),
				Description:  "Instance type. Supported values include: `master` - master instance, `dr` - disaster recovery instance, and `ro` - read-only instance.",
			},
			"status": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1, 4, 5}),
				Description:  "Instance status. Available values: `0` - Creating; `1` - Running; `4` - Isolating; `5` - Isolated.",
			},
			"security_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Security groups ID of instance.",
			},
			"pay_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Deprecated:   "It has been deprecated from version 1.36.0. Please use `charge_type` instead.",
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Pay type of instance, `0`: prepay, `1`: postpaid.",
			},
			"charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{MYSQL_CHARGE_TYPE_PREPAID, MYSQL_CHARGE_TYPE_POSTPAID}),
				Description:  "Pay type of instance, valid values are `PREPAID` and `POSTPAID`.",
			},
			"instance_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of mysql instance.",
			},
			"engine_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"5.1", "5.5", "5.6", "5.7", "8.0"}),
				Description:  "The version number of the database engine to use. Supported versions include 5.5/5.6/5.7/8.0.",
			},
			"init_flag": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Initialization mark. Available values: `0` - Uninitialized; `1` - Initialized.",
			},
			"with_dr": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Indicates whether to query disaster recovery instances.",
			},
			"with_ro": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Indicates whether to query read-only instances.",
			},
			"with_master": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateAllowedIntValue([]int{0, 1}),
				Description:  "Indicates whether to query master instances.",
			},
			"offset": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      0,
				ValidateFunc: validateIntegerInRange(0, 1000),
				Description:  "Record offset. Default is 0.",
			},
			"limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      20,
				ValidateFunc: validateIntegerInRange(1, 2000),
				Description:  "Number of results returned for a single request. Default is `20`, and maximum is 2000.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"instance_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of instances. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mysql_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID, such as `cdb-c1nl9rpv`. It is identical to the instance ID displayed in the database console page.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of mysql instance.",
						},
						"instance_role": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type. Supported values include: `master` - master instance, `dr` - disaster recovery instance, and `ro` - read-only instance.",
						},
						"init_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Initialization mark. Available values: `0` - Uninitialized; `1` - Initialized.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance status. Available values: `0` - Creating; `1` - Running; `4` - Isolating; `5` - Isolated.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Information of available zone.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Auto renew flag. NOTES: Only supported prepay instance.",
						},
						"engine_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version number of the database engine to use. Supported versions include `5.5`/`5.6`/`5.7`/`8.0`.",
						},
						"cpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CPU count.",
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size (in MB).",
						},
						"volume_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk capacity (in GB).",
						},
						"internet_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of public network.",
						},
						"internet_host": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Public network domain name.",
						},
						"internet_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Public network port.",
						},
						"intranet_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance IP for internal access.",
						},
						"intranet_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Transport layer port number for internal purpose.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID to which the current instance belongs.",
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of Virtual Private Cloud.",
						},
						"subnet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of subnet to which the current instance belongs.",
						},
						"slave_sync_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Data replication mode. `0` - Async replication; `1` - Semisync replication; `2` - Strongsync replication.",
						},
						"device_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Supported instance model. `HA` - high available version; `Basic` - basic version.",
						},
						"pay_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Pay type of instance, `0`: prepaid, `1`: postpaid.",
						},
						"charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Pay type of instance.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time at which a instance is created.",
						},
						"dead_line_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expire date of instance. NOTES: Only supported prepay instance.",
						},
						"master_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicates the master instance ID of recovery instances.",
						},
						"ro_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "ID list of read-only type associated with the current instance.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"dr_instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "ID list of disaster-recovery type associated with the current instance.",
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
	defer logElapsed("data_source.tencentcloud_mysql_instance.read")()

	logId := getLogId(contextNil)

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
	if chargeType, ok := d.GetOk("charge_type"); ok {
		var payType int
		if chargeType == MYSQL_CHARGE_TYPE_PREPAID {
			payType = MysqlPayByMonth
		} else {
			payType = MysqlPayByUse
		}
		payTypeValue := uint64(payType)
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

	instanceDetails := response.Response.Items
	instanceList := make([]map[string]interface{}, 0, len(instanceDetails))
	ids := make([]string, 0, len(instanceDetails))
	for _, item := range instanceDetails {
		mapping := map[string]interface{}{
			"mysql_id":        item.InstanceId,
			"instance_name":   item.InstanceName,
			"instance_role":   MYSQL_ROLE_MAP[*item.InstanceType],
			"init_flag":       item.InitFlag,
			"status":          item.Status,
			"zone":            item.Zone,
			"auto_renew_flag": item.AutoRenew,
			"engine_version":  item.EngineVersion,
			"cpu_core_count":  item.Cpu,
			"memory_size":     item.Memory,
			"volume_size":     item.Volume,
			"internet_status": item.WanStatus,
			"internet_host":   item.WanDomain,
			"internet_port":   item.WanPort,
			"intranet_ip":     item.Vip,
			"intranet_port":   item.Vport,
			"project_id":      item.ProjectId,
			"vpc_id":          item.UniqVpcId,
			"subnet_id":       item.UniqSubnetId,
			"slave_sync_mode": item.ProtectMode,
			"device_type":     item.DeviceType,
			"pay_type":        item.PayType,
			"create_time":     item.CreateTime,
			"dead_line_time":  item.DeadlineTime,
			"charge_type":     MYSQL_CHARGE_TYPE[int(*item.PayType)],
		}
		if item.MasterInfo != nil {
			mapping["master_instance_id"] = item.MasterInfo.InstanceId
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
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_list", instanceList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance list fail, reason:%s\n ", logId, err.Error())
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		err = writeToFile(output.(string), instanceList)
		if err != nil {
			return err
		}
	}
	return nil
}
