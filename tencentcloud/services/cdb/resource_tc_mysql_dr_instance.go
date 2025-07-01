package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlDrInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlDrInstanceCreate,
		Read:   resourceTencentCloudMysqlDrInstanceRead,
		Update: resourceTencentCloudMysqlDrInstanceUpdate,
		Delete: resourceTencentCloudMysqlDrInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"prepaid_period": 1,
				"force_delete":   false,
			}),
		},
		Schema: map[string]*schema.Schema{
			"master_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Indicates the master instance ID of recovery instances.",
			},
			"master_region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The zone information of the primary instance is required when you purchase a disaster recovery instance.",
			},
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
			"device_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specify device type, available values: `UNIVERSAL` (default), `EXCLUSIVE`, `BASIC`.",
			},
			"slave_deploy_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Default:      0,
				Description:  "Availability zone deployment method. Available values: 0 - Single availability zone; 1 - Multiple availability zones.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates which availability zone will be used.",
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
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Project ID, default value is 0.",
			},
			"slave_sync_mode": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1, 2}),
				Default:      0,
				Description:  "Data replication mode. 0 - Async replication; 1 - Semisync replication; 2 - Strongsync replication.",
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate whether to delete instance directly or not. Default is `false`. If set true, the instance will be deleted instead of staying recycle bin. Note: only works for `PREPAID` instance.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Instance tags.",
			},

			// Computed values
			"intranet_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "instance intranet IP.",
			},
		},
	}
}

func resourceTencentCloudMysqlDrInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_dr_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	mysqlService := MysqlService{client: client}
	masterClient := *client
	if v, ok := d.GetOk("master_region"); ok {
		masterClient.Region = v.(string)
	}

	masterInstanceId := d.Get("master_instance_id").(string)
	var masterinstace *cdb.InstanceInfo
	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		masterService := MysqlService{client: &masterClient}
		instace, err := masterService.DescribeDBInstanceById(context.TODO(), masterInstanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if instace == nil {
			return resource.NonRetryableError(fmt.Errorf("Master instance does not exist, %s", masterInstanceId))
		}
		masterinstace = instace
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if *masterinstace.InstanceId != masterInstanceId {
		return fmt.Errorf("Master instance is not mastert, %s", masterInstanceId)
	}

	payType := getPayType(d).(int)
	if payType == MysqlPayByMonth {
		request := cdb.NewCreateDBInstanceRequest()
		if err := mysqlDrInstanceSet(ctx, request, d, meta, *masterinstace); err != nil {
			return err
		}

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDBInstance(request)
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
	} else if payType == MysqlPayByUse {
		request := cdb.NewCreateDBInstanceHourRequest()
		if err := mysqlDrInstanceSet(ctx, request, d, meta, *masterinstace); err != nil {
			return err
		}

		response, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().CreateDBInstanceHour(request)
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
	} else {
		return fmt.Errorf("mysql not support this pay type yet.")
	}

	mysqlID := d.Id()

	err = resource.Retry(4*tccommon.ReadRetryTimeout, func() *resource.RetryError {
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
		err = fmt.Errorf("create mysql task status is %d,we won't wait for it finish", *mysqlInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql  task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		log.Printf("[DEBUG]Mysql instance create, resourceName:%s\n", resourceName)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlDrInstanceRead(d, meta)
}

func resourceTencentCloudMysqlDrInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_dr_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	mysqlInfo, errRet := mysqlService.DescribeDBInstanceById(ctx, d.Id())
	if errRet != nil {
		return fmt.Errorf("Describe mysql instance fails, reaseon %v", errRet.Error())
	}
	if mysqlInfo == nil {
		d.SetId("")
		return nil
	}
	if MysqlDelStates[*mysqlInfo.Status] {
		mysqlInfo = nil
		d.SetId("")
		return nil
	}

	_ = d.Set("master_instance_id", *mysqlInfo.MasterInfo.InstanceId)
	_ = d.Set("master_region", *mysqlInfo.MasterInfo.Region)

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
	_ = d.Set("availability_zone", mysqlInfo.Zone)
	_ = d.Set("slave_deploy_mode", mysqlInfo.DeployMode)
	_ = d.Set("slave_sync_mode", mysqlInfo.ProtectMode)
	_ = d.Set("project_id", mysqlInfo.ProjectId)

	if mysqlInfo.SlaveInfo != nil && mysqlInfo.SlaveInfo.First != nil {
		_ = d.Set("first_slave_zone", mysqlInfo.SlaveInfo.First.Zone)
	}

	if mysqlInfo.SlaveInfo != nil && mysqlInfo.SlaveInfo.First != nil {
		_ = d.Set("first_slave_zone", mysqlInfo.SlaveInfo.First.Zone)
	}

	securityGroups, err := mysqlService.DescribeDBSecurityGroups(ctx, d.Id())
	if err != nil {
		sdkErr, ok := err.(*sdkError.TencentCloudSDKError)
		if ok {
			if sdkErr.Code == MysqlInstanceIdNotFound3 {
				mysqlInfo = nil
				d.SetId("")
				return nil
			}
		}
		return err
	}
	_ = d.Set("security_groups", securityGroups)

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "cdb", "instanceId", tcClient.Region, d.Id())
	if err != nil {
		return err
	}

	if err := d.Set("tags", tags); err != nil {
		log.Printf("[CRITAL]%s provider set tags fail, reason:%s\n ", logId, err.Error())
		return nil
	}

	_ = d.Set("intranet_ip", mysqlInfo.Vip)
	_ = d.Set("intranet_port", int(*mysqlInfo.Vport))

	return nil
}

func resourceTencentCloudMysqlDrInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_dr_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlID := d.Id()
	payType := getPayType(d).(int)

	if d.HasChange("prepaid_period") {
		if v, ok := d.GetOk("charge_type"); ok {
			if v.(string) != MYSQL_CHARGE_TYPE_PREPAID {
				return fmt.Errorf("`prepaid_period` only support prepaid instance.")
			}
		}
	}

	if d.HasChange("charge_type") {
		oldChargeTypeInterface, newChargeTypeInterface := d.GetChange("charge_type")
		oldChargeType := oldChargeTypeInterface.(string)
		newChargeType := newChargeTypeInterface.(string)

		if oldChargeType == MYSQL_CHARGE_TYPE_PREPAID && newChargeType == MYSQL_CHARGE_TYPE_POSTPAID {
			return fmt.Errorf("`charge_type` only supports modification from `POSTPAID` to `PREPAID`.")
		}

		var period int = 1
		if v, ok := d.GetOkExists("prepaid_period"); ok {
			period = v.(int)
		}

		request := cdb.NewRenewDBInstanceRequest()
		request.InstanceId = &mysqlID
		request.ModifyPayType = helper.String(newChargeType)
		request.TimeSpan = helper.IntInt64(period)
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, err := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().RenewDBInstance(request)
			if err != nil {
				return tccommon.RetryError(err, tccommon.InternalError)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Renew DB instance failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	d.Partial(true)

	if payType == MysqlPayByMonth {
		if d.HasChange("auto_renew_flag") {
			renewFlag := int64(d.Get("auto_renew_flag").(int))
			mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
			if err := mysqlService.ModifyAutoRenewFlag(ctx, d.Id(), renewFlag); err != nil {
				return err
			}

		}
	}
	err := mysqlAllInstanceRoleUpdate(ctx, d, meta, true)
	if err != nil {
		return err
	}

	immutableFields := []string{
		"master_instance_id",
		"zone",
		"master_region",
	}

	for _, f := range immutableFields {
		if d.HasChange(f) {
			return fmt.Errorf("argument `%s` cannot be modified for now", f)
		}
	}

	d.Partial(false)

	return resourceTencentCloudMysqlDrInstanceRead(d, meta)
}

func resourceTencentCloudMysqlDrInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_dr_instance.delete")()

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
			if _, ok := err.(*sdkError.TencentCloudSDKError); !ok {
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
	if err == nil {
		log.Printf("[WARN]this mysql is dr instance, it is released asynchronously, and the bound resource is not now fully released now\n")
	}
	return err
}

func mysqlDrInstanceSet(ctx context.Context, requestInter interface{}, d *schema.ResourceData, meta interface{}, instance cdb.InstanceInfo) error {
	requestByMonth, okByMonth := requestInter.(*cdb.CreateDBInstanceRequest)
	requestByUse, _ := requestInter.(*cdb.CreateDBInstanceHourRequest)

	instanceRole := "dr"
	if okByMonth {
		requestByMonth.InstanceRole = &instanceRole
	} else {
		requestByUse.InstanceRole = &instanceRole
	}

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

	payType, ok := d.GetOkExists("pay_type")
	if (!ok || payType == -1) && okByMonth {
		var period int
		if !ok || payType == -1 {
			period = d.Get("prepaid_period").(int)
		} else {
			period = d.Get("period").(int)
		}

		requestByMonth.Period = helper.IntInt64(period)
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

	if v, ok := d.GetOk("master_region"); ok {
		masterRegion := helper.String(v.(string))
		if okByMonth {
			requestByMonth.MasterRegion = masterRegion
		} else {
			requestByUse.MasterRegion = masterRegion
		}
	}

	if v, ok := d.GetOk("auto_renew_flag"); ok {
		autoRenewFlag := helper.IntInt64(v.(int))
		if okByMonth {
			requestByMonth.AutoRenewFlag = autoRenewFlag
		} else {
			requestByUse.AutoRenewFlag = autoRenewFlag
		}
	}

	if v, ok := d.GetOk("master_instance_id"); ok {
		masterInstanceId := helper.String(v.(string))
		if okByMonth {
			requestByMonth.MasterInstanceId = masterInstanceId
		} else {
			requestByUse.MasterInstanceId = masterInstanceId
		}
	}

	if v, ok := d.GetOk("slave_deploy_mode"); ok {
		deployMode := helper.IntInt64(v.(int))
		if okByMonth {
			requestByMonth.DeployMode = deployMode
		} else {
			requestByUse.DeployMode = deployMode
		}
	}

	if v, ok := d.GetOk("availability_zone"); ok {
		zone := helper.String(v.(string))
		if okByMonth {
			requestByMonth.Zone = zone
		} else {
			requestByUse.Zone = zone
		}
	}

	if v, ok := d.GetOk("first_slave_zone"); ok {
		slaveZone := helper.String(v.(string))
		if okByMonth {
			requestByMonth.SlaveZone = slaveZone
		} else {
			requestByUse.SlaveZone = slaveZone
		}
	}

	if v, ok := d.GetOk("second_slave_zone"); ok {
		backupZone := helper.String(v.(string))
		if okByMonth {
			requestByMonth.BackupZone = backupZone
		} else {
			requestByUse.BackupZone = backupZone
		}
	}

	if v, ok := d.GetOk("project_id"); ok {
		projectId := helper.IntInt64(v.(int))
		if okByMonth {
			requestByMonth.ProjectId = projectId
		} else {
			requestByUse.ProjectId = projectId
		}
	}

	if v, ok := d.GetOk("slave_sync_mode"); ok {
		slaveSyncMode := helper.IntInt64(v.(int))
		if okByMonth {
			requestByMonth.ProtectMode = slaveSyncMode
		} else {
			requestByUse.ProtectMode = slaveSyncMode
		}
	}

	if engineType := instance.EngineType; engineType != nil {
		if okByMonth {
			requestByMonth.EngineType = engineType
		} else {
			requestByUse.EngineType = engineType
		}
	}

	if engineVersion := instance.EngineVersion; engineVersion != nil {
		if okByMonth {
			requestByMonth.EngineVersion = engineVersion
		} else {
			requestByUse.EngineVersion = engineVersion
		}
	}

	autoSyncFlag := helper.IntInt64(1)
	if okByMonth {
		requestByMonth.AutoSyncFlag = autoSyncFlag
	} else {
		requestByUse.AutoSyncFlag = autoSyncFlag
	}

	return nil

}
