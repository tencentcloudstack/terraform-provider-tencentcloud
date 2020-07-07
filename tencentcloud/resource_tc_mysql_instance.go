/*
Provides a mysql instance resource to create master database instances.

~> **NOTE:** If this mysql has readonly instance, the terminate operation of the mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

Example Usage

```hcl
resource "tencentcloud_mysql_instance" "default" {
  internet_service = 1
  engine_version   = "5.7"
  charge_type = "POSTPAID"
  root_password     = "********"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-4"
  second_slave_zone = "ap-guangzhou-4"
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-4"
  project_id        = 201901010001
  instance_name     = "myTestMysql"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-12mt3l31"
  subnet_id         = "subnet-9uivyb1g"
  intranet_port     = 3306
  security_groups   = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }

  parameters = {
    max_connections = "1000"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TencentMsyqlBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: validateStringLengthInRange(1, 100),
			Description:  "The name of a mysql instance.",
		},
		"pay_type": {
			Type:          schema.TypeInt,
			Deprecated:    "It has been deprecated from version 1.36.0. Please use `charge_type` instead.",
			Optional:      true,
			ValidateFunc:  validateAllowedIntValue([]int{MysqlPayByMonth, MysqlPayByUse}),
			ConflictsWith: []string{"charge_type", "prepaid_period"},
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				return true
			},
			Default:     -1,
			Description: "Pay type of instance, 0: prepaid, 1: postpaid.",
		},
		"charge_type": {
			Type:          schema.TypeString,
			ForceNew:      true,
			Optional:      true,
			ValidateFunc:  validateAllowedStringValue([]string{MYSQL_CHARGE_TYPE_PREPAID, MYSQL_CHARGE_TYPE_POSTPAID}),
			ConflictsWith: []string{"pay_type", "period"},
			Default:       MYSQL_CHARGE_TYPE_POSTPAID,
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				if (olds == "" && news == MYSQL_CHARGE_TYPE_POSTPAID) ||
					(olds == MYSQL_CHARGE_TYPE_POSTPAID && news == "") {
					if v, ok := d.GetOkExists("pay_type"); ok && v.(int) == MysqlPayByUse {
						return true
					}
				} else if (olds == "" && news == MYSQL_CHARGE_TYPE_PREPAID) ||
					(olds == MYSQL_CHARGE_TYPE_PREPAID && news == "") {
					if v, ok := d.GetOkExists("pay_type"); ok && v.(int) == MysqlPayByMonth {
						return true
					}
				}
				return olds == news
			},
			Description: "Pay type of instance, valid values are `PREPAID`, `POSTPAID`. Default is `POSTPAID`.",
		},
		"period": {
			Type:          schema.TypeInt,
			Deprecated:    "It has been deprecated from version 1.36.0. Please use `prepaid_period` instead.",
			Optional:      true,
			Default:       -1,
			ConflictsWith: []string{"charge_type", "prepaid_period"},
			ValidateFunc:  validateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				return true
			},
			Description: "Period of instance. NOTES: Only supported prepaid instance.",
		},
		"prepaid_period": {
			Type:          schema.TypeInt,
			Optional:      true,
			Default:       1,
			ConflictsWith: []string{"pay_type", "period"},
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				return true
			},
			ValidateFunc: validateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
			Description:  "Period of instance. NOTES: Only supported prepaid instance.",
		},
		"auto_renew_flag": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Auto renew flag. NOTES: Only supported prepaid instance.",
		},
		"intranet_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateIntegerInRange(1024, 65535),
			Default:      3306,
			Description:  "Public access port, rang form 1024 to 65535 and default value is 3306.",
		},
		"mem_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory size (in MB).",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Disk size (in GB).",
		},
		"vpc_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringLengthInRange(1, 100),
			Description:  "ID of VPC, which can be modified once every 24 hours and can't be removed.",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: validateStringLengthInRange(1, 100),
			Description:  "Private network ID. If vpc_id is set, this value is required.",
		},

		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return hashcode.String(v.(string))
			},
			Description: "Security groups to use.",
		},

		"tags": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Instance tags.",
		},

		"force_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Indicate whether to delete instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance. When the main mysql instance set true, this para of the readonly mysql instance will not take effect.",
		},
		// Computed values
		"intranet_ip": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "instance intranet IP.",
		},

		"locked": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether the instance is locked. 0 - No; 1 - Yes.",
		},
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Instance status. Available values: 0 - Creating; 1 - Running; 4 - Isolating; 5 - Isolated.",
		},
		"task_status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates which kind of operations is being executed.",
		},

		"gtid": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether GTID is enable. 0 - Not enabled; 1 - Enabled.",
		},
	}
}

func resourceTencentCloudMysqlInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"parameters": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "List of parameters to use.",
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Indicates whether to enable the access to an instance from public network: 0 - No, 1 - Yes.",
		},
		"engine_version": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
			Default:      MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-1],
			Description:  "The version number of the database engine to use. Supported versions include 5.5/5.6/5.7, and default is 5.7.",
		},

		"availability_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Computed:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"root_password": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validateMysqlPassword,
			Description:  "Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.",
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ForceNew:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.",
		},
		"first_slave_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Zone information about first slave instance.",
		},
		"second_slave_zone": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
			Description: "Zone information about second slave instance.",
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			Default:      0,
			Description:  "Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.",
		},
		"project_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "Project ID, default value is 0.",
		},

		// Computed values
		"internet_host": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "host for public access.",
		},
		"internet_port": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Access port for public access.",
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		specialInfo[k] = v
	}
	return &schema.Resource{
		Create: resourceTencentCloudMysqlInstanceCreate,
		Read:   resourceTencentCloudMysqlInstanceRead,
		Update: resourceTencentCloudMysqlInstanceUpdate,
		Delete: resourceTencentCloudMysqlInstanceDelete,

		Schema: specialInfo,
	}
}

/*
   [master] and [dr] and [ro] all need set
*/
func mysqlAllInstanceRoleSet(ctx context.Context, requestInter interface{}, d *schema.ResourceData, meta interface{}) error {
	requestByMonth, okByMonth := requestInter.(*cdb.CreateDBInstanceRequest)
	requestByUse, _ := requestInter.(*cdb.CreateDBInstanceHourRequest)

	var goodsNum int64 = 1
	if okByMonth {
		requestByMonth.GoodsNum = &goodsNum
	} else {
		requestByUse.GoodsNum = &goodsNum
	}

	if instanceNameInterface, ok := d.GetOk("instance_name"); ok {
		instanceName := instanceNameInterface.(string)
		if okByMonth {
			requestByMonth.InstanceName = &instanceName
		} else {
			requestByUse.InstanceName = &instanceName
		}
	}

	intranetPort := int64(d.Get("intranet_port").(int))
	if okByMonth {
		requestByMonth.Port = &intranetPort
	} else {
		requestByUse.Port = &intranetPort
	}

	memSize := int64(d.Get("mem_size").(int))
	if okByMonth {
		requestByMonth.Memory = &memSize
	} else {
		requestByUse.Memory = &memSize
	}

	volumeSize := int64(d.Get("volume_size").(int))
	if okByMonth {
		requestByMonth.Volume = &volumeSize
	} else {
		requestByUse.Volume = &volumeSize
	}

	if strInterface, ok := d.GetOk("vpc_id"); ok {
		str := strInterface.(string)
		if okByMonth {
			requestByMonth.UniqVpcId = &str
		} else {
			requestByUse.UniqVpcId = &str
		}

	}
	if strInterface, ok := d.GetOk("subnet_id"); ok {
		str := strInterface.(string)
		if okByMonth {
			requestByMonth.UniqSubnetId = &str
		} else {
			requestByUse.UniqSubnetId = &str
		}
	}
	err := fmt.Errorf("You have to set both vpc_id and subnet_id")
	if okByMonth {
		if requestByMonth.UniqVpcId != nil && requestByMonth.UniqSubnetId == nil {
			return err
		}
		if requestByMonth.UniqVpcId == nil && requestByMonth.UniqSubnetId != nil {
			return err
		}
	} else {
		if requestByUse.UniqVpcId != nil && requestByUse.UniqSubnetId == nil {
			return err
		}
		if requestByUse.UniqVpcId == nil && requestByUse.UniqSubnetId != nil {
			return err
		}
	}

	if temp, ok := d.GetOkExists("security_groups"); ok {
		securityGroups := temp.(*schema.Set).List()
		requestSecurityGroup := make([]*string, 0, len(securityGroups))
		for _, v := range securityGroups {
			str := v.(string)
			requestSecurityGroup = append(requestSecurityGroup, &str)
		}
		if okByMonth {
			requestByMonth.SecurityGroup = requestSecurityGroup
		} else {
			requestByUse.SecurityGroup = requestSecurityGroup
		}
	}
	return nil

}

/*
 [master] need set
*/
func mysqlMasterInstanceRoleSet(ctx context.Context, requestInter interface{}, d *schema.ResourceData, meta interface{}) error {
	requestByMonth, okByMonth := requestInter.(*cdb.CreateDBInstanceRequest)
	requestByUse, _ := requestInter.(*cdb.CreateDBInstanceHourRequest)

	if parametersMap, ok := d.Get("parameters").(map[string]interface{}); ok {
		requestParamList := make([]*cdb.ParamInfo, 0, len(parametersMap))
		for k, v := range parametersMap {
			key := k
			value := v.(string)
			var paramInfo = cdb.ParamInfo{Name: &key, Value: &value}
			requestParamList = append(requestParamList, &paramInfo)
		}
		if okByMonth {
			requestByMonth.ParamList = requestParamList
		} else {
			requestByUse.ParamList = requestParamList
		}
	}

	if intInterface, ok := d.GetOkExists("project_id"); ok {
		intv := int64(intInterface.(int))
		if okByMonth {
			requestByMonth.ProjectId = &intv
		} else {
			requestByUse.ProjectId = &intv
		}
	}

	engineVersion := d.Get("engine_version").(string)
	if okByMonth {
		requestByMonth.EngineVersion = &engineVersion
	} else {
		requestByUse.EngineVersion = &engineVersion
	}
	if stringInterface, ok := d.GetOk("availability_zone"); ok {
		str := stringInterface.(string)
		if okByMonth {
			requestByMonth.Zone = &str
		} else {
			requestByUse.Zone = &str
		}
	}

	if stringInterface, ok := d.GetOk("root_password"); ok {
		str := stringInterface.(string)
		if okByMonth {
			requestByMonth.Password = &str
		} else {
			requestByUse.Password = &str
		}
	}
	slaveDeployMode := int64(d.Get("slave_deploy_mode").(int))
	if okByMonth {
		requestByMonth.DeployMode = &slaveDeployMode
	} else {
		requestByUse.DeployMode = &slaveDeployMode
	}
	if stringInterface, ok := d.GetOk("first_slave_zone"); ok {
		str := stringInterface.(string)
		if okByMonth {
			requestByMonth.SlaveZone = &str
		} else {
			requestByUse.SlaveZone = &str
		}
	}

	if stringInterface, ok := d.GetOk("second_slave_zone"); ok {
		str := stringInterface.(string)
		if okByMonth {
			requestByMonth.BackupZone = &str
		} else {
			requestByUse.BackupZone = &str
		}
	}
	slaveSyncMode := int64(d.Get("slave_sync_mode").(int))
	if okByMonth {
		requestByMonth.ProtectMode = &slaveSyncMode
	} else {
		requestByUse.ProtectMode = &slaveSyncMode
	}
	return nil
}

func mysqlCreateInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	logId := getLogId(ctx)

	request := cdb.NewCreateDBInstanceRequest()

	payType, oldOk := d.GetOkExists("pay_type")
	var period int
	if !oldOk || payType == -1 {
		period = d.Get("prepaid_period").(int)
	} else {
		period = d.Get("period").(int)
	}

	request.Period = helper.IntInt64(period)

	autoRenewFlag := int64(d.Get("auto_renew_flag").(int))
	request.AutoRenewFlag = &autoRenewFlag

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateDBInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	if len(response.Response.InstanceIds) != 1 {
		return fmt.Errorf("mysql CreateDBInstance return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}
	d.SetId(*response.Response.InstanceIds[0])
	return nil
}

func mysqlCreateInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	logId := getLogId(ctx)
	request := cdb.NewCreateDBInstanceHourRequest()

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	if err := mysqlMasterInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	response, err := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().CreateDBInstanceHour(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}
	if len(response.Response.InstanceIds) != 1 {
		return fmt.Errorf("mysql CreateDBInstanceHour return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}
	d.SetId(*response.Response.InstanceIds[0])
	return nil
}

func resourceTencentCloudMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	payType := getPayType(d).(int)

	if payType == MysqlPayByMonth {
		err := mysqlCreateInstancePayByMonth(ctx, d, meta)
		if err != nil {
			return err
		}
	} else if payType == MysqlPayByUse {
		err := mysqlCreateInstancePayByUse(ctx, d, meta)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("mysql not support this pay type yet.")
	}

	mysqlID := d.Id()

	err := resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, mysqlID)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if mysqlInfo == nil {
			err = fmt.Errorf("mysqlid %s instance not exists", mysqlID)
			return resource.NonRetryableError(err)
		}
		if *mysqlInfo.Status == MYSQL_STATUS_DELIVING {
			return resource.RetryableError(fmt.Errorf("create mysql task  status is MYSQL_STATUS_DELIVING(%d)", MYSQL_STATUS_DELIVING))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return nil
		}
		err = fmt.Errorf("create mysql    task status is %v,we won't wait for it finish", *mysqlInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql  task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	//internet service
	internetService := d.Get("internet_service").(int)
	if internetService == 1 {
		asyncRequestId, err := mysqlService.OpenWanService(ctx, d.Id())
		if err != nil {
			return err
		}
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("create account task  status is %s", taskStatus))
			}
			err = fmt.Errorf("open internet service task status is %s,we won't wait for it finish ,it show message:%s", ",",
				message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s open internet service   fail, reason:%s\n ", logId, err.Error())
			return err
		}
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlInstanceRead(d, meta)
}

func tencentMsyqlBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}, master bool) (mysqlInfo *cdb.InstanceInfo,
	errRet error) {

	if d.Id() == "" {
		return
	}

	logId := getLogId(ctx)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	mysqlInfo, errRet = mysqlService.DescribeDBInstanceById(ctx, d.Id())
	if errRet != nil {
		errRet = fmt.Errorf("Describe mysql instance fails, reaseon %s", errRet.Error())
		return
	}
	if mysqlInfo == nil {
		d.SetId("")
		return
	}
	if MysqlDelStates[*mysqlInfo.Status] {
		mysqlInfo = nil
		d.SetId("")
		return
	}

	_ = d.Set("instance_name", *mysqlInfo.InstanceName)

	_ = d.Set("charge_type", MYSQL_CHARGE_TYPE[int(*mysqlInfo.PayType)])
	_ = d.Set("pay_type", -1)
	_ = d.Set("period", -1)
	if int(*mysqlInfo.PayType) == MysqlPayByMonth {
		tempInt, _ := d.Get("prepaid_period").(int)
		if tempInt == 0 {
			_ = d.Set("prepaid_period", 1)
		}
	}

	if *mysqlInfo.AutoRenew == MYSQL_RENEW_CLOSE {
		*mysqlInfo.AutoRenew = MYSQL_RENEW_NOUSE
	}
	_ = d.Set("auto_renew_flag", int(*mysqlInfo.AutoRenew))
	_ = d.Set("mem_size", mysqlInfo.Memory)
	_ = d.Set("volume_size", mysqlInfo.Volume)
	if d.Get("vpc_id").(string) != "" {
		errRet = d.Set("vpc_id", mysqlInfo.UniqVpcId)
	}
	if d.Get("subnet_id").(string) != "" {
		errRet = d.Set("subnet_id", mysqlInfo.UniqSubnetId)
	}

	securityGroups, err := mysqlService.DescribeDBSecurityGroups(ctx, d.Id())
	if err != nil {
		sdkErr, ok := err.(*errors.TencentCloudSDKError)
		if ok {
			if sdkErr.Code == MysqlInstanceIdNotFound3 {
				mysqlInfo = nil
				d.SetId("")
				return
			}
		}
		errRet = err
		return
	}
	_ = d.Set("security_groups", securityGroups)
	if master {
		isGTIDOpen, err := mysqlService.CheckDBGTIDOpen(ctx, d.Id())
		if err != nil {
			errRet = err
			return
		}
		_ = d.Set("gtid", int(isGTIDOpen))
	}

	tcClient := meta.(*TencentCloudClient).apiV3Conn
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cdb", "instanceId", tcClient.Region, d.Id())
	if err != nil {
		errRet = err
		return
	}

	if err := d.Set("tags", tags); err != nil {
		log.Printf("[CRITAL]%s provider set tags fail, reason:%s\n ", logId, err.Error())
		return
	}

	_ = d.Set("intranet_ip", mysqlInfo.Vip)
	_ = d.Set("intranet_port", int(*mysqlInfo.Vport))

	if *mysqlInfo.CdbError != 0 {
		_ = d.Set("locked", 1)
	} else {
		_ = d.Set("locked", 0)
	}
	_ = d.Set("status", mysqlInfo.Status)
	_ = d.Set("task_status", mysqlInfo.TaskStatus)
	return
}

func resourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	var mysqlInfo *cdb.InstanceInfo
	var e error
	var onlineHas = true
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, e = tencentMsyqlBasicInfoRead(ctx, d, meta, true)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		if mysqlInfo == nil {
			d.SetId("")
			onlineHas = false
			return nil
		}
		_ = d.Set("project_id", int(*mysqlInfo.ProjectId))
		_ = d.Set("engine_version", mysqlInfo.EngineVersion)
		if *mysqlInfo.WanStatus == 1 {
			_ = d.Set("internet_service", 1)
			_ = d.Set("internet_host", mysqlInfo.WanDomain)
			_ = d.Set("internet_port", int(*mysqlInfo.WanPort))
		} else {
			_ = d.Set("internet_service", 0)
			_ = d.Set("internet_host", "")
			_ = d.Set("internet_port", 0)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Fail to get basic info from mysql, reaseon %s", err.Error())
	}
	if !onlineHas {
		return nil
	}
	parametersMap, ok := d.Get("parameters").(map[string]interface{})
	if !ok {
		log.Printf("[INFO] %v  config error,parameters is not map[string]interface{}\n", logId)
	} else {
		var cares []string
		for k := range parametersMap {
			cares = append(cares, k)
		}

		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			caresParameters, e := mysqlService.DescribeCaresParameters(ctx, d.Id(), cares)
			if e != nil {
				if mysqlService.NotFoundMysqlInstance(e) {
					d.SetId("")
					onlineHas = false
					return nil
				}
				return retryError(e)
			}
			if e := d.Set("parameters", caresParameters); e != nil {
				log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, e.Error())
				return resource.NonRetryableError(e)
			}
			_ = d.Set("availability_zone", mysqlInfo.Zone)
			return nil
		})
		if err != nil {
			return fmt.Errorf("Describe CaresParameters Fail, reason:%s", err.Error())
		}
		if !onlineHas {
			return nil
		}
	}
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		backConfig, e := mysqlService.DescribeDBInstanceConfig(ctx, d.Id())
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return retryError(e)
		}
		_ = d.Set("slave_sync_mode", int(*backConfig.Response.ProtectMode))
		_ = d.Set("slave_deploy_mode", int(*backConfig.Response.DeployMode))
		if backConfig.Response.SlaveConfig != nil && *backConfig.Response.SlaveConfig.Zone != "" {
			//if you set ,i set
			if _, ok := d.GetOk("first_slave_zone"); ok {
				_ = d.Set("first_slave_zone", *backConfig.Response.SlaveConfig.Zone)
			}
		}
		if backConfig.Response.BackupConfig != nil && *backConfig.Response.BackupConfig.Zone != "" {
			//if you set ,i set
			if _, ok := d.GetOk("second_slave_zone"); ok {
				_ = d.Set("second_slave_zone", *backConfig.Response.BackupConfig.Zone)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Describe DBInstanceConfig Fail, reason:%s", err.Error())
	}
	return nil
}

/*
   [master] and [dr] and [ro] all need update
*/
func mysqlAllInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	logId := getLogId(ctx)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	if d.HasChange("instance_name") {
		if err := mysqlService.ModifyDBInstanceName(ctx, d.Id(), d.Get("instance_name").(string)); err != nil {
			return err
		}
		d.SetPartial("instance_name")
	}

	if d.HasChange("intranet_port") || d.HasChange("vpc_id") || d.HasChange("subnet_id") {
		var (
			intranetPort = int64(d.Get("intranet_port").(int))
			vpcId        = d.Get("vpc_id").(string)
			subnetId     = d.Get("subnet_id").(string)
		)
		if d.HasChange("vpc_id") {
			vpcId = d.Get("vpc_id").(string)
			if vpcId == "" {
				return fmt.Errorf("[vpc_id]Once a setting cannot be deleted,it can only be modified")
			}
		}
		if d.HasChange("subnet_id") {
			subnetId = d.Get("subnet_id").(string)
			if vpcId == "" {
				return fmt.Errorf("[subnet_id]Once a setting cannot be deleted,it can only be modified")
			}
		}

		if err := mysqlService.ModifyDBInstanceVipVport(ctx, d.Id(), vpcId, subnetId, intranetPort); err != nil {
			return err
		}
		if d.HasChange("intranet_port") {
			d.SetPartial("intranet_port")
		}
		if d.HasChange("vpc_id") {
			d.SetPartial("vpc_id")
		}
		if d.HasChange("subnet_id") {
			d.SetPartial("subnet_id")
		}
	}

	if d.HasChange("mem_size") || d.HasChange("volume_size") {

		memSize := int64(d.Get("mem_size").(int))
		volumeSize := int64(d.Get("volume_size").(int))

		asyncRequestId, err := mysqlService.UpgradeDBInstance(ctx, d.Id(), memSize, volumeSize)

		if err != nil {
			return err
		}

		err = resource.Retry(6*time.Hour, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)

			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}

			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("update mysql  mem_size/volume_size status is %s", taskStatus))
			}
			err = fmt.Errorf("update mysql  mem_size/volume_size task status is %s,we won't wait for it finish ,it show message:%s",
				",", message)
			return resource.NonRetryableError(err)
		})

		if err != nil {
			log.Printf("[CRITAL]%s update mysql  mem_size/volume_size  fail, reason:%s\n ", logId, err.Error())
			return err
		}
		if d.HasChange("mem_size") {
			d.SetPartial("mem_size")
		}
		if d.HasChange("volume_size") {
			d.SetPartial("volume_size")
		}
	}

	if d.HasChange("security_groups") {

		oldValue, newValue := d.GetChange("security_groups")

		oldSecuritygroups := oldValue.(*schema.Set).List()
		newSecuritygroups := newValue.(*schema.Set).List()

		isDelete := false

		if len(newSecuritygroups) == 0 && len(oldSecuritygroups) != 0 {
			isDelete = true
			newSecuritygroups = append(newSecuritygroups, oldSecuritygroups[0])
		}

		var newStrs = make([]string, 0, len(newSecuritygroups))
		for _, v := range newSecuritygroups {
			newStrs = append(newStrs, v.(string))
		}

		if err := mysqlService.ModifyDBInstanceSecurityGroups(ctx, d.Id(), newStrs); err != nil {
			return err
		}
		if isDelete {
			oldFirst := oldSecuritygroups[0].(string)
			if err := mysqlService.DisassociateSecurityGroup(ctx, d.Id(), oldFirst); err != nil {
				return err
			}
		}
		d.SetPartial("security_groups")
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("cdb", "instanceId", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		d.SetPartial("tags")
	}

	return nil
}

/*
 [master] need set
*/
func mysqlMasterInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(ctx)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	if d.HasChange("project_id") {
		newProjectId := int64(d.Get("project_id").(int))
		if err := mysqlService.ModifyDBInstanceProject(ctx, d.Id(), newProjectId); err != nil {
			return err
		}
		d.SetPartial("project_id")
	}

	if d.HasChange("parameters") {

		oldValue, newValue := d.GetChange("parameters")

		oldParameters := oldValue.(map[string]interface{})
		newParameters := newValue.(map[string]interface{})

		//set(oldParameters-newParameters)need set to Default
		var oldMinusNew = make(map[string]interface{}, len(oldParameters))
		for k, v := range oldParameters {
			if _, has := newParameters[k]; !has {
				oldMinusNew[k] = v
			}
		}

		supportsParameters := make(map[string]*cdb.ParameterDetail)
		parameterList, err := mysqlService.DescribeInstanceParameters(ctx, d.Id())
		if err != nil {
			return err
		}
		for _, parameter := range parameterList {
			supportsParameters[*parameter.Name] = parameter
		}

		for parameName := range newParameters {
			if _, has := supportsParameters[parameName]; !has {
				return fmt.Errorf("this mysql not support param %s set", parameName)
			}
		}

		modifyParameters := make(map[string]string)
		for parameName, detail := range supportsParameters {
			//set to Default
			if old, has := oldMinusNew[parameName]; has {
				modifyParameters[parameName] = *detail.Default
				log.Printf("[DEBUG] %s mysql need set param  %+v to default:%+v, old:%+v\n", logId, parameName, *detail.Default, old)
				continue
			}
			//set(newParameters) need set add or modify
			if v, has := newParameters[parameName]; has {
				modifyParameters[parameName] = v.(string)
				continue
			}
		}

		log.Printf("[DEBUG] %s mysql need set params:%+v\n", logId, modifyParameters)

		tag := "modify param"
		if len(modifyParameters) > 0 {
			asyncRequestId, err := mysqlService.ModifyInstanceParam(ctx, d.Id(), modifyParameters)
			if err != nil {
				log.Printf("[CRITAL]%s update mysql %s fail, reason:%s\n ", logId, tag, err.Error())
				return err
			}
			err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
				taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
				if err != nil {
					if _, ok := err.(*errors.TencentCloudSDKError); !ok {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}
				if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
					return nil
				}
				if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
					return resource.RetryableError(fmt.Errorf("update mysql  %s status is %s", tag, taskStatus))
				}
				err = fmt.Errorf("update mysql   task status is %s,we won't wait for it finish ,it show message:%s",
					tag, message)
				return resource.NonRetryableError(err)
			})
			if err != nil {
				log.Printf("[CRITAL]%s update mysql  %s  fail, reason:%s\n ", logId, tag, err.Error())
				return err
			}
		}
		d.SetPartial("parameters")
	}

	if d.HasChange("internet_service") {
		internetService := d.Get("internet_service").(int)
		var (
			asyncRequestId       = ""
			err            error = nil
			tag                  = "close internet service"
		)
		if internetService == 0 {
			asyncRequestId, err = mysqlService.CloseWanService(ctx, d.Id())
		} else {
			asyncRequestId, err = mysqlService.OpenWanService(ctx, d.Id())
			tag = "open internet service"
		}

		if err != nil {
			log.Printf("[CRITAL]%s update mysql %s fail, reason:%s\n ", logId, tag, err.Error())
			return err
		}
		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("update mysql  %s status is %s", tag, taskStatus))
			}
			err = fmt.Errorf("update mysql task status is %s,we won't wait for it finish ,it show message:%s",
				tag, message)
			return resource.NonRetryableError(err)
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mysql  %s  fail, reason:%s\n ", logId, tag, err.Error())
			return err
		}
		d.SetPartial("internet_service")
	}

	if d.HasChange("root_password") {

		var (
			newPassword = d.Get("root_password").(string)
			userName    = "root"
		)

		asyncRequestId, err := mysqlService.ModifyAccountPassword(ctx, d.Id(), userName, MYSQL_DEFAULT_ACCOUNT_HOST, newPassword)

		if err != nil {
			return err
		}

		err = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
			taskStatus, message, err := mysqlService.DescribeAsyncRequestInfo(ctx, asyncRequestId)
			if err != nil {
				if _, ok := err.(*errors.TencentCloudSDKError); !ok {
					return resource.RetryableError(err)
				} else {
					return resource.NonRetryableError(err)
				}
			}
			if taskStatus == MYSQL_TASK_STATUS_SUCCESS {
				return nil
			}
			if taskStatus == MYSQL_TASK_STATUS_INITIAL || taskStatus == MYSQL_TASK_STATUS_RUNNING {
				return resource.RetryableError(fmt.Errorf("change root password status is %s", taskStatus))
			}
			err = fmt.Errorf("change root password task status is %s,we won't wait for it finish ,it show message:%s",
				taskStatus, message)
			return resource.NonRetryableError(err)
		})
		if err != nil {
			log.Printf("[CRITAL]%s change root password   fail, reason:%s\n ", logId, err.Error())
			return err
		}
		d.SetPartial("root_password")
	}
	return nil
}

func mysqlUpdateInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	if err := mysqlAllInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}

	if d.HasChange("auto_renew_flag") {
		renewFlag := int64(d.Get("auto_renew_flag").(int))
		mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
		if err := mysqlService.ModifyAutoRenewFlag(ctx, d.Id(), renewFlag); err != nil {
			return err
		}
		d.SetPartial("auto_renew_flag")
	}

	if d.HasChange("period") || d.HasChange("prepaid_period") {
		return fmt.Errorf("After the initialization of the field[%s] to set does not make sense", "period or prepaid_period")
	}
	return nil
}

func mysqlUpdateInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	if err := mysqlAllInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	payType := getPayType(d).(int)

	d.Partial(true)
	if payType == MysqlPayByMonth {
		err := mysqlUpdateInstancePayByMonth(ctx, d, meta)
		if err != nil {
			return err
		}
	} else if payType == MysqlPayByUse {
		err := mysqlUpdateInstancePayByUse(ctx, d, meta)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("mysql not support this pay type yet.")
	}
	d.Partial(false)
	time.Sleep(7 * time.Second)

	return resourceTencentCloudMysqlInstanceRead(d, meta)
}

func resourceTencentCloudMysqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, err := mysqlService.IsolateDBInstance(ctx, d.Id())
		if err != nil {
			//for the pay order wait
			return retryError(err, InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var hasDeleted = false

	payType := getPayType(d).(int)
	forceDelete := d.Get("force_delete").(bool)
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, d.Id())

		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if mysqlInfo == nil {
			hasDeleted = true
			return nil
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATING || *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return resource.RetryableError(fmt.Errorf("mysql isolating."))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_ISOLATED {
			return nil
		}
		return resource.NonRetryableError(fmt.Errorf("after IsolateDBInstance mysql Status is %d", *mysqlInfo.Status))
	})

	if hasDeleted {
		return nil
	}
	if err != nil {
		return err
	}

	if payType == MysqlPayByMonth && !forceDelete {
		return nil
	}

	err = mysqlService.OfflineIsolatedInstances(ctx, d.Id())
	if err != nil {
		return err
	}

	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeIsolatedDBInstanceById(ctx, d.Id())
		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if mysqlInfo == nil {
			return nil
		} else {
			if mysqlInfo.RoGroups != nil && len(mysqlInfo.RoGroups) > 0 {
				log.Printf("[WARN]this mysql has RoGroups , RoGroups is released asynchronously, and the bound resource is not now fully released now\n")
				return nil
			}
			return resource.RetryableError(fmt.Errorf("after OfflineIsolatedInstances mysql Status is %d", *mysqlInfo.Status))
		}
	})
	return err
}

func getPayType(d *schema.ResourceData) (payType interface{}) {
	chargeType := d.Get("charge_type")
	payType, oldOk := d.GetOkExists("pay_type")

	if !oldOk || payType == -1 {
		if chargeType == MYSQL_CHARGE_TYPE_PREPAID {
			payType = MysqlPayByMonth
		} else {
			payType = MysqlPayByUse
		}
	}
	return
}
