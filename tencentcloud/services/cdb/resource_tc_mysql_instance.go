package cdb

import (
	"context"
	"fmt"
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

//internal version: replace import begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace import end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

var importMysqlFlag = false

const InWindow = 1

func TencentMsyqlBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:         schema.TypeString,
			Required:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
			Description:  "The name of a mysql instance.",
		},
		"pay_type": {
			Type:          schema.TypeInt,
			Deprecated:    "It has been deprecated from version 1.36.0. Please use `charge_type` instead.",
			Optional:      true,
			ValidateFunc:  tccommon.ValidateAllowedIntValue([]int{MysqlPayByMonth, MysqlPayByUse}),
			ConflictsWith: []string{"charge_type", "prepaid_period"},
			DiffSuppressFunc: func(k, olds, news string, d *schema.ResourceData) bool {
				return true
			},
			Default:     -1,
			Description: "Pay type of instance. Valid values: `0`, `1`. `0`: prepaid, `1`: postpaid.",
		},
		"charge_type": {
			Type:          schema.TypeString,
			ForceNew:      true,
			Optional:      true,
			ValidateFunc:  tccommon.ValidateAllowedStringValue([]string{MYSQL_CHARGE_TYPE_PREPAID, MYSQL_CHARGE_TYPE_POSTPAID}),
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
			Description: "Pay type of instance. Valid values:`PREPAID`, `POSTPAID`. Default is `POSTPAID`.",
		},
		"period": {
			Type:          schema.TypeInt,
			Deprecated:    "It has been deprecated from version 1.36.0. Please use `prepaid_period` instead.",
			Optional:      true,
			Default:       -1,
			ConflictsWith: []string{"charge_type", "prepaid_period"},
			ValidateFunc:  tccommon.ValidateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
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
			ValidateFunc:  tccommon.ValidateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
			Description:   "Period of instance. NOTES: Only supported prepaid instance.",
		},
		"auto_renew_flag": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Auto renew flag. NOTES: Only supported prepaid instance.",
		},
		"intranet_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateIntegerInRange(1024, 65535),
			Default:      3306,
			Description:  "Public access port. Valid value ranges: [1024~65535]. The default value is `3306`.",
		},
		"mem_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Memory size (in MB).",
		},
		"cpu": {
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
			Description: "CPU cores.",
		},
		"volume_size": {
			Type:        schema.TypeInt,
			Required:    true,
			Description: "Disk size (in GB).",
		},
		"vpc_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
			Description:  "ID of VPC, which can be modified once every 24 hours and can't be removed.",
		},
		"subnet_id": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: tccommon.ValidateStringLengthInRange(1, 100),
			Description:  "Private network ID. If `vpc_id` is set, this value is required.",
		},

		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return helper.HashString(v.(string))
			},
			Description: "Security groups to use.",
		},
		"param_template_id": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specify parameter template id.",
		},
		"fast_upgrade": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Specify whether to enable fast upgrade when upgrade instance spec, available value: `1` - enabled, `0` - disabled.",
		},
		"device_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			Description: "Specify device type, available values:\n" +
				"	- `UNIVERSAL` (default): universal instance,\n" +
				"	- `EXCLUSIVE`: exclusive instance,\n" +
				"	- `BASIC_V2`: ONTKE single-node instance,\n" +
				"	- `CLOUD_NATIVE_CLUSTER`: cluster version standard type,\n" +
				"	- `CLOUD_NATIVE_CLUSTER_EXCLUSIVE`: cluster version enhanced type.\n" +
				"If it is not specified, it defaults to a universal instance.",
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
			Description: "Indicate whether to delete instance directly or not. Default is `false`. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance. When the main mysql instance set true, this para of the readonly mysql instance will not take effect.",
		},
		"wait_switch": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Switch the method of accessing new instances, default is `0`. Supported values include: `0` - switch immediately, `1` - switch in time window.",
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
			Description: "Indicates whether the instance is locked. Valid values: `0`, `1`. `0` - No; `1` - Yes.",
		},
		"status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Instance status. Valid values: `0`, `1`, `4`, `5`. `0` - Creating; `1` - Running; `4` - Isolating; `5` - Isolated.",
		},
		"task_status": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates which kind of operations is being executed.",
		},

		"gtid": {
			Type:        schema.TypeInt,
			Computed:    true,
			Description: "Indicates whether GTID is enable. `0` - Not enabled; `1` - Enabled.",
		},
	}
}

func ResourceTencentCloudMysqlInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"parameters": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "List of parameters to use.",
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Indicates whether to enable the access to an instance from public network: 0 - No, 1 - Yes.",
		},
		"engine_version": {
			Type:         schema.TypeString,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
			Default:      MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-2],
			Description:  "The version number of the database engine to use. Supported versions include 5.5/5.6/5.7/8.0, and default is 5.7. Upgrade the instance engine version to support 5.6/5.7 and switch immediately.",
		},
		"engine_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Instance engine type. The default value is `InnoDB`. Supported values include `InnoDB` and `RocksDB`.",
		},
		"upgrade_subversion": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Whether it is a kernel subversion upgrade, supported values: 1 - upgrade the kernel subversion; 0 - upgrade the database engine version. Only need to fill in when upgrading kernel subversion and engine version.",
		},
		"max_deay_time": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Latency threshold. Value range 1~10. Only need to fill in when upgrading kernel subversion and engine version.",
		},

		"availability_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Indicates which availability zone will be used.",
		},
		"root_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: tccommon.ValidateMysqlPassword,
			Description:  "Password of root account. This parameter can be specified when you purchase master instances, but it should be ignored when you purchase read-only instances or disaster recovery instances.",
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
			Default:      0,
			Description:  "Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.",
		},
		"first_slave_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Zone information about first slave instance.",
		},
		"second_slave_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Zone information about second slave instance.",
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
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
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importMysqlFlag = true
				defaultValues := map[string]interface{}{
					"charge_type":       MYSQL_CHARGE_TYPE_POSTPAID,
					"prepaid_period":    1,
					"auto_renew_flag":   0,
					"intranet_port":     3306,
					"force_delete":      false,
					"internet_service":  0,
					"engine_version":    MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-2],
					"slave_deploy_mode": 0,
					"slave_sync_mode":   0,
					"project_id":        0,
				}

				for k, v := range defaultValues {
					_ = d.Set(k, v)
				}
				return []*schema.ResourceData{d}, nil
			},
		},
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

	cpu := int64(d.Get("cpu").(int))
	if okByMonth {
		requestByMonth.Cpu = &cpu
	} else {
		requestByUse.Cpu = &cpu
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

	if v, ok := d.GetOk("param_template_id"); ok {
		paramTemplateId := helper.IntInt64(v.(int))
		if okByMonth {
			requestByMonth.ParamTemplateId = paramTemplateId
		} else {
			requestByUse.ParamTemplateId = paramTemplateId
		}
	}

	if v, ok := d.GetOk("device_type"); ok {
		deviceType := helper.String(v.(string))
		if okByMonth {
			requestByMonth.DeviceType = deviceType
		} else {
			requestByUse.DeviceType = deviceType
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

	isBasic := isBasicDevice(d)
	if parametersMap, ok := d.Get("parameters").(map[string]interface{}); ok && !isBasic {
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

	if stringInterface, ok := d.GetOk("engine_type"); ok {
		engineType := stringInterface.(string)
		if okByMonth {
			requestByMonth.EngineType = &engineType
		} else {
			requestByUse.EngineType = &engineType
		}
	}

	if stringInterface, ok := d.GetOk("availability_zone"); ok {
		str := stringInterface.(string)
		if okByMonth {
			requestByMonth.Zone = &str
		} else {
			requestByUse.Zone = &str
		}
	}

	if v, ok := d.GetOk("root_password"); ok && v.(string) != "" && !isBasic {
		str := v.(string)
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

	logId := tccommon.GetLogId(ctx)

	request := cdb.NewCreateDBInstanceRequest()
	clientToken := helper.BuildToken()
	request.ClientToken = &clientToken
	//internal version: replace var begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	//internal version: replace var end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

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

	var response *cdb.CreateDBInstanceResponse
	err := resource.Retry(2*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		// shadowed response will not pass to outside
		r, inErr := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDBInstance(request)
		if inErr != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), inErr.Error())
			//internal version: replace bpass begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			//internal version: replace bpass end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
			return tccommon.RetryError(inErr)
		}

		if len(r.Response.InstanceIds) < 1 && clientToken != "" {
			return resource.RetryableError(fmt.Errorf("%s returns nil instanceIds but client token provided, retrying", request.GetAction()))
		}

		response = r
		//internal version: replace instanceId begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace instanceId end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		return nil
	})

	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.InstanceIds) != 1 {
		return fmt.Errorf("mysql CreateDBInstance return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}
	//internal version: replace setId begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	d.SetId(*response.Response.InstanceIds[0])
	//internal version: replace setId end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	return nil
}

func mysqlCreateInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	logId := tccommon.GetLogId(ctx)
	request := cdb.NewCreateDBInstanceHourRequest()
	clientToken := helper.BuildToken()
	request.ClientToken = &clientToken

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	if err := mysqlMasterInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	var response *cdb.CreateDBInstanceHourResponse
	err := resource.Retry(2*tccommon.WriteRetryTimeout, func() *resource.RetryError {
		// shadowed response will not pass to outside
		r, inErr := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDBInstanceHour(request)
		if inErr != nil {
			return tccommon.RetryError(inErr)
		}
		if len(r.Response.InstanceIds) < 1 && clientToken != "" {
			return resource.RetryableError(fmt.Errorf("%s returns nil instanceIds but client token provided, retrying", request.GetAction()))
		}
		response = r
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceIds) != 1 {
		return fmt.Errorf("mysql CreateDBInstanceHour return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}

	id := *response.Response.InstanceIds[0]

	d.SetId(id)

	return nil
}

func resourceTencentCloudMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	//internal version: replace mysqlServer begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	//internal version: replace mysqlServer end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

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

	//internal version: replace setTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	//internal version: replace setTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	err := resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

	return resourceTencentCloudMysqlInstanceRead(d, meta)
}

func tencentMsyqlBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}, master bool) (mysqlInfo *cdb.InstanceInfo,
	errRet error) {

	if d.Id() == "" {
		return
	}

	logId := tccommon.GetLogId(ctx)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
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
	_ = d.Set("cpu", mysqlInfo.Cpu)
	_ = d.Set("volume_size", mysqlInfo.Volume)
	_ = d.Set("vpc_id", mysqlInfo.UniqVpcId)
	_ = d.Set("subnet_id", mysqlInfo.UniqSubnetId)
	_ = d.Set("device_type", mysqlInfo.DeviceType)

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

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var mysqlInfo *cdb.InstanceInfo
	var e error
	var onlineHas = true
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mysqlInfo, e = tencentMsyqlBasicInfoRead(ctx, d, meta, true)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return tccommon.RetryError(e)
		}
		if mysqlInfo == nil {
			d.SetId("")
			onlineHas = false
			return nil
		}
		_ = d.Set("project_id", int(*mysqlInfo.ProjectId))
		_ = d.Set("engine_version", mysqlInfo.EngineVersion)
		if mysqlInfo.EngineType != nil {
			_ = d.Set("engine_type", *mysqlInfo.EngineType)
		}
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
	if importMysqlFlag {
		// import logic:
		log.Printf("[INFO] %v ,begin to import parameter...\n", logId)
		parameterList, err := mysqlService.DescribeInstanceParameters(ctx, *mysqlInfo.InstanceId)
		if err != nil {
			return err
		}

		parameters := make(map[string]string)
		for _, v := range parameterList {
			parameters[*v.Name] = *v.CurrentValue
		}
		if e := d.Set("parameters", parameters); e != nil {
			log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, e.Error())
			return e
		}
		_ = d.Set("availability_zone", mysqlInfo.Zone)
		importMysqlFlag = false
	} else if parametersMap, ok := d.Get("parameters").(map[string]interface{}); ok {
		// read logic:
		var cares []string
		for k := range parametersMap {
			cares = append(cares, k)
		}

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			caresParameters, e := mysqlService.DescribeCaresParameters(ctx, d.Id(), cares)
			if e != nil {
				if mysqlService.NotFoundMysqlInstance(e) {
					d.SetId("")
					onlineHas = false
					return nil
				}
				return tccommon.RetryError(e)
			}
			if e := d.Set("parameters", caresParameters); e != nil {
				log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, e.Error())
				return resource.NonRetryableError(e)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Describe CaresParameters Fail, reason:%s", err.Error())
		}
		if !onlineHas {
			return nil
		}
	}
	err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		backConfig, e := mysqlService.DescribeDBInstanceConfig(ctx, d.Id())
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				onlineHas = false
				return nil
			}
			return tccommon.RetryError(e)
		}
		_ = d.Set("slave_sync_mode", int(*backConfig.Response.ProtectMode))
		_ = d.Set("slave_deploy_mode", int(*backConfig.Response.DeployMode))
		if backConfig.Response.SlaveConfig != nil && *backConfig.Response.SlaveConfig.Zone != "" {
			_ = d.Set("first_slave_zone", *backConfig.Response.SlaveConfig.Zone)
		}
		if backConfig.Response.BackupConfig != nil && *backConfig.Response.BackupConfig.Zone != "" {
			_ = d.Set("second_slave_zone", *backConfig.Response.BackupConfig.Zone)
		}

		if backConfig.Response.Zone != nil {
			_ = d.Set("availability_zone", *backConfig.Response.Zone)
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
func mysqlAllInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}, isReadonly bool) error {

	logId := tccommon.GetLogId(ctx)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("instance_name") {
		if err := mysqlService.ModifyDBInstanceName(ctx, d.Id(), d.Get("instance_name").(string)); err != nil {
			return err
		}

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
	}

	if isReadonly {
		if d.HasChange("mem_size") || d.HasChange("cpu") || d.HasChange("volume_size") || d.HasChange("device_type") {

			memSize := int64(d.Get("mem_size").(int))
			cpu := int64(d.Get("cpu").(int))
			volumeSize := int64(d.Get("volume_size").(int))
			deviceType := ""
			var waitSwitch int64 = 0

			fastUpgrade := int64(0)
			if v, ok := d.GetOk("fast_upgrade"); ok {
				fastUpgrade = int64(v.(int))
			}
			if v, ok := d.GetOk("device_type"); ok {
				deviceType = v.(string)
			}

			if v, ok := d.GetOkExists("wait_switch"); ok {
				waitSwitch = int64(v.(int))
			}

			asyncRequestId, err := mysqlService.UpgradeDBInstance(ctx, d.Id(), memSize, cpu, volumeSize, fastUpgrade, deviceType, -1, -1, "", "", waitSwitch)

			if err != nil {
				return err
			}

			if waitSwitch != InWindow {
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
						return resource.RetryableError(fmt.Errorf("update mysql mem_size/volume_size status is %s", taskStatus))
					}
					err = fmt.Errorf("update mysql mem_size/volume_size task status is %s, we won't wait for it finish, it show message:%s",
						taskStatus, message)
					return resource.NonRetryableError(err)
				})

				if err != nil {
					log.Printf("[CRITAL]%s update mysql mem_size/volume_size fail, reason:%s\n", logId, err.Error())
					return err
				}
			} else {
				err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
					mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, d.Id())

					if err != nil {
						if _, ok := err.(*errors.TencentCloudSDKError); !ok {
							return resource.RetryableError(err)
						} else {
							return resource.NonRetryableError(err)
						}
					}

					if *mysqlInfo.TaskStatus == 15 {
						return nil
					}
					return resource.RetryableError(fmt.Errorf("update mysql mem_size/volume_size task status is %v", mysqlInfo.TaskStatus))
				})

				if err != nil {
					log.Printf("[CRITAL]%s update mysql mem_size/volume_size fail, reason:%s\n", logId, err.Error())
					return err
				}
			}
		}
	} else {
		if d.HasChange("mem_size") || d.HasChange("cpu") || d.HasChange("volume_size") || d.HasChange("device_type") || d.HasChange("slave_deploy_mode") || d.HasChange("first_slave_zone") || d.HasChange("second_slave_zone") || d.HasChange("slave_sync_mode") {

			memSize := int64(d.Get("mem_size").(int))
			cpu := int64(d.Get("cpu").(int))
			volumeSize := int64(d.Get("volume_size").(int))
			var slaveDeployMode int64 = -1
			slaveSyncMode := int64(d.Get("slave_sync_mode").(int))
			deviceType := ""
			firstSlaveZone := ""
			secondSlaveZone := ""
			var waitSwitch int64 = 0

			fastUpgrade := int64(0)
			if v, ok := d.GetOk("fast_upgrade"); ok {
				fastUpgrade = int64(v.(int))
			}
			if v, ok := d.GetOk("device_type"); ok {
				deviceType = v.(string)
			}

			if d.HasChange("first_slave_zone") || d.HasChange("second_slave_zone") {
				if v, ok := d.GetOk("first_slave_zone"); ok {
					firstSlaveZone = v.(string)
				}

				if v, ok := d.GetOk("second_slave_zone"); ok {
					secondSlaveZone = v.(string)
				}
			} else {
				mysqlInfo, e := tencentMsyqlBasicInfoRead(ctx, d, meta, true)
				if e != nil {
					return e
				}
				if mysqlInfo != nil && mysqlInfo.SlaveInfo != nil && mysqlInfo.SlaveInfo.First != nil && mysqlInfo.SlaveInfo.First.Zone != nil {
					firstSlaveZone = *mysqlInfo.SlaveInfo.First.Zone
				}
				if mysqlInfo != nil && mysqlInfo.SlaveInfo != nil && mysqlInfo.SlaveInfo.Second != nil && mysqlInfo.SlaveInfo.Second.Zone != nil {
					firstSlaveZone = *mysqlInfo.SlaveInfo.Second.Zone
				}
			}

			if v, ok := d.GetOkExists("wait_switch"); ok {
				waitSwitch = int64(v.(int))
			}

			if d.HasChange("slave_deploy_mode") {
				if v, ok := d.GetOkExists("slave_deploy_mode"); ok {
					slaveDeployMode = int64(v.(int))
				}
			}

			asyncRequestId, err := mysqlService.UpgradeDBInstance(ctx, d.Id(), memSize, cpu, volumeSize, fastUpgrade, deviceType, slaveDeployMode, slaveSyncMode, firstSlaveZone, secondSlaveZone, waitSwitch)

			if err != nil {
				return err
			}

			if waitSwitch != InWindow {
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
						return resource.RetryableError(fmt.Errorf("update mysql mem_size/volume_size status is %s", taskStatus))
					}
					err = fmt.Errorf("update mysql mem_size/volume_size task status is %s, we won't wait for it finish, it show message:%s",
						taskStatus, message)
					return resource.NonRetryableError(err)
				})

				if err != nil {
					log.Printf("[CRITAL]%s update mysql mem_size/volume_size fail, reason:%s\n ", logId, err.Error())
					return err
				}
			} else {
				err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
					mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, d.Id())

					if err != nil {
						if _, ok := err.(*errors.TencentCloudSDKError); !ok {
							return resource.RetryableError(err)
						} else {
							return resource.NonRetryableError(err)
						}
					}

					if *mysqlInfo.TaskStatus == 15 {
						return nil
					}
					return resource.RetryableError(fmt.Errorf("update mysql engineVersion task status is %v", mysqlInfo.TaskStatus))
				})

				if err != nil {
					log.Printf("[CRITAL]%s update mysql engineVersion fail, reason:%s\n", logId, err.Error())
					return err
				}
			}
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

	}

	if d.HasChange("engine_version") || d.HasChange("upgrade_subversion") || d.HasChange("max_deay_time") {
		engineVersion := ""
		var upgradeSubversion int64
		var maxDelayTime int64
		var waitSwitch int64 = 0
		if v, ok := d.GetOk("engine_version"); ok {
			engineVersion = v.(string)
		}
		if v, ok := d.GetOk("upgrade_subversion"); ok {
			upgradeSubversion = int64(v.(int))
		}
		if v, ok := d.GetOk("max_deay_time"); ok {
			maxDelayTime = int64(v.(int))
		}
		if v, ok := d.GetOkExists("wait_switch"); ok {
			waitSwitch = int64(v.(int))
		}

		asyncRequestId, err := mysqlService.UpgradeDBInstanceEngineVersion(ctx, d.Id(), engineVersion, upgradeSubversion, maxDelayTime, waitSwitch)
		if err != nil {
			return err
		}

		if waitSwitch != InWindow {
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
					return resource.RetryableError(fmt.Errorf("update mysql engineVersion status is %s", taskStatus))
				}
				err = fmt.Errorf("update mysql engineVersion task status is %s, we won't wait for it finish, it show message:%s",
					taskStatus, message)
				return resource.NonRetryableError(err)
			})

			if err != nil {
				log.Printf("[CRITAL]%s update mysql engineVersion fail, reason:%s\n", logId, err.Error())
				return err
			}
		} else {
			err = resource.Retry(tccommon.ReadRetryTimeout*5, func() *resource.RetryError {
				mysqlInfo, err := mysqlService.DescribeDBInstanceById(ctx, d.Id())

				if err != nil {
					if _, ok := err.(*errors.TencentCloudSDKError); !ok {
						return resource.RetryableError(err)
					} else {
						return resource.NonRetryableError(err)
					}
				}

				if *mysqlInfo.TaskStatus == 15 {
					return nil
				}
				return resource.RetryableError(fmt.Errorf("update mysql engineVersion task status is %v", mysqlInfo.TaskStatus))
			})

			if err != nil {
				log.Printf("[CRITAL]%s update mysql engineVersion fail, reason:%s\n", logId, err.Error())
				return err
			}

		}
	}

	if d.HasChange("tags") {

		oldValue, newValue := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldValue.(map[string]interface{}), newValue.(map[string]interface{}))

		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cdb", "instanceId", region, d.Id())
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
		//internal version: replace waitTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
		//internal version: replace waitTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
	}

	if d.HasChange("param_template_id") {
		return fmt.Errorf("argument `param_template_id` cannot be modified for now")
	}

	if d.HasChange("availability_zone") {
		return fmt.Errorf("argument `availability_zone` cannot be modified for now")
	}

	if d.HasChange("engine_version") {
		return fmt.Errorf("argument `engine_version` cannot be modified for now")
	}
	return nil
}

/*
[master] need set
*/
func mysqlMasterInstanceRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	if d.HasChange("project_id") {
		newProjectId := int64(d.Get("project_id").(int))
		if err := mysqlService.ModifyDBInstanceProject(ctx, d.Id(), newProjectId); err != nil {
			return err
		}

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

		version := d.Get("engine_version").(string)
		if version == "8.0" && oldParameters["lower_case_table_names"] != newParameters["lower_case_table_names"] {
			return fmt.Errorf("this mysql 8.0 not support param `lower_case_table_names` set")
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
			err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
		err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

		err = resource.Retry(10*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

	}
	return nil
}

func mysqlUpdateInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	if err := mysqlAllInstanceRoleUpdate(ctx, d, meta, false); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}

	if d.HasChange("auto_renew_flag") {
		renewFlag := int64(d.Get("auto_renew_flag").(int))
		mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		if err := mysqlService.ModifyAutoRenewFlag(ctx, d.Id(), renewFlag); err != nil {
			return err
		}

	}

	if d.HasChange("period") || d.HasChange("prepaid_period") {
		return fmt.Errorf("After the initialization of the field[%s] to set does not make sense", "period or prepaid_period")
	}
	return nil
}

func mysqlUpdateInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	if err := mysqlAllInstanceRoleUpdate(ctx, d, meta, false); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleUpdate(ctx, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTencentCloudMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

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
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_instance.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		_, err := mysqlService.IsolateDBInstance(ctx, d.Id())
		if err != nil {
			//for the pay order wait
			return tccommon.RetryError(err, tccommon.InternalError)
		}
		return nil
	})

	if err != nil {
		return err
	}

	var hasDeleted = false

	payType := getPayType(d).(int)
	forceDelete := d.Get("force_delete").(bool)
	err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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

	err = resource.Retry(7*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := mysqlService.DescribeIsolatedDBInstanceById(ctx, d.Id())
		if err != nil {
			if _, ok := err.(*errors.TencentCloudSDKError); !ok {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(err)
			}
		}
		if mysqlInfo == nil || *mysqlInfo.Status == 6 {
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

func isBasicDevice(d *schema.ResourceData) bool {
	v, ok := d.GetOk("device_type")
	if !ok {
		return false
	}
	return v.(string) == "BASIC"
}
