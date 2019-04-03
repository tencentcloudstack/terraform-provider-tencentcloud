package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func TencentMsyqlBasicInfo() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"instance_name": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"pay_type": {
			Type:         schema.TypeInt,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{MysqlPayByMonth, MysqlPayByUse}),
			Default:      0,
		},
		"period": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validateAllowedIntValue(MYSQL_AVAILABLE_PERIOD),
		},
		"auto_renew_flag": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"engine_version": {
			Type:         schema.TypeString,
			ForceNew:     true,
			Optional:     true,
			ValidateFunc: validateAllowedStringValue(MYSQL_SUPPORTS_ENGINE),
			Default:      MYSQL_SUPPORTS_ENGINE[len(MYSQL_SUPPORTS_ENGINE)-1],
		},
		"mem_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"volume_size": {
			Type:     schema.TypeInt,
			Required: true,
		},
		"vpc_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"internet_service": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"gtid": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"project_id": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		"security_groups": {
			Type:     schema.TypeSet,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Set: func(v interface{}) int {
				return hashcode.String(v.(string))
			},
		},

		"parameters": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"tags": {
			Type:     schema.TypeMap,
			Optional: true,
		},

		// Computed values
		"intranet_ip": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"locked": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"task_status": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"internet_host": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"internet_port": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

func resourceTencentCloudMysqlInstance() *schema.Resource {
	specialInfo := map[string]*schema.Schema{
		"availability_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"root_password": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validateMysqlPassword,
		},
		"slave_deploy_mode": {
			Type:         schema.TypeInt,
			Optional:     true,
			ValidateFunc: validateAllowedIntValue([]int{0, 1}),
			Default:      0,
		},
		"first_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"second_slave_zone": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"slave_sync_mode": {
			Type:         schema.TypeInt,
			ValidateFunc: validateAllowedIntValue([]int{0, 1, 2}),
			Optional:     true,
			Default:      0,
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

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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

	engineVersion := d.Get("engine_version").(string)
	if okByMonth {
		requestByMonth.EngineVersion = &engineVersion
	} else {
		requestByUse.EngineVersion = &engineVersion
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

	if intInterface, ok := d.GetOkExists("project_id"); ok {
		intv := int64(intInterface.(int))
		if okByMonth {
			requestByMonth.ProjectId = &intv
		} else {
			requestByUse.ProjectId = &intv
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

	if tagsMap, ok := d.Get("tags").(map[string]interface{}); ok {
		requestResourceTags := make([]*cdb.TagInfo, 0, len(tagsMap))
		for k, v := range tagsMap {
			key := k
			value := v.(string)
			var tagInfo cdb.TagInfo
			tagInfo.TagKey = &key
			tagInfo.TagValue = []*string{&value}
			requestResourceTags = append(requestResourceTags, &tagInfo)
		}
		if okByMonth {
			requestByMonth.ResourceTags = requestResourceTags
		} else {
			requestByUse.ResourceTags = requestResourceTags
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

	logId := GetLogId(ctx)

	request := cdb.NewCreateDBInstanceRequest()

	period := int64(d.Get("period").(int))
	request.Period = &period

	autoRenewFlag := int64(d.Get("auto_renew_flag").(int))
	request.AutoRenewFlag = &autoRenewFlag

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}
	if err := mysqlMasterInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	log.Println(request.ToJsonString(), "2222222222222222222222222222")
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
		fmt.Errorf("mysql CreateDBInstance return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}
	d.SetId(*response.Response.InstanceIds[0])
	return nil
}

func mysqlCreateInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(ctx)
	request := cdb.NewCreateDBInstanceHourRequest()

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	if err := mysqlMasterInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	log.Println(request.ToJsonString(), "1111111111111111111111111111111")

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
		fmt.Errorf("mysql CreateDBInstanceHour return len(InstanceIds) is not 1,but %d", len(response.Response.InstanceIds))
	}
	d.SetId(*response.Response.InstanceIds[0])
	return nil
}

func resourceTencentCloudMysqlInstanceCreate(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	_ = mysqlService
	payType := d.Get("pay_type").(int)

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

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
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
		err = fmt.Errorf("create mysql    task status is %s,we won't wait for it finish", *mysqlInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s create mysql  task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	//初始化

	//开外网

	//开gtid

	return resourceTencentCloudMysqlInstanceRead(d, meta)
}

func tencentMsyqlBasicInfoRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (mysqlInfo *cdb.InstanceInfo,
	errRet error) {

	logId := GetLogId(ctx)

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

	d.Set("instance_name", *mysqlInfo.InstanceName)
	d.Set("pay_type", int(*mysqlInfo.PayType))

	if int(*mysqlInfo.PayType) == MysqlPayByMonth {
		tempInt, _ := d.Get("period").(int)
		if tempInt == 0 {
			d.Set("period", 1)
		}
	}

	if *mysqlInfo.AutoRenew == MYSQL_RENEW_CLOSE {
		*mysqlInfo.AutoRenew = MYSQL_RENEW_NOUSE
	}
	d.Set("auto_renew_flag", int(*mysqlInfo.AutoRenew))

	d.Set("engine_version", *mysqlInfo.EngineVersion)
	d.Set("mem_size", *mysqlInfo.Memory)
	d.Set("volume_size", *mysqlInfo.Volume)
	d.Set("vpc_id", *mysqlInfo.UniqVpcId)
	d.Set("subnet_id", *mysqlInfo.UniqSubnetId)

	if *mysqlInfo.WanStatus == 1 {
		d.Set("internet_service", 1)
		d.Set("internet_host", *mysqlInfo.WanDomain)
		d.Set("internet_port", int(*mysqlInfo.WanPort))
	} else {
		d.Set("internet_service", 0)
		d.Set("internet_host", "")
		d.Set("internet_port", 0)
	}

	isGTIDOpen, err := mysqlService.CheckDBGTIDOpen(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	d.Set("gtid", int(isGTIDOpen))
	d.Set("project_id", int(*mysqlInfo.ProjectId))

	securityGroups, err := mysqlService.DescribeDBSecurityGroups(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	d.Set("security_groups", securityGroups)

	parametersMap, ok := d.Get("parameters").(map[string]interface{})
	if !ok {
		log.Printf("[INFO] %d  config error,parameters is not map[string]interface{}\n", logId)
	} else {
		var cares []string
		for k, _ := range parametersMap {
			cares = append(cares, k)
		}
		caresParameters, err := mysqlService.DescribeCaresParameters(ctx, d.Id(), cares)
		if err != nil {
			errRet = err
			return
		}
		if err := d.Set("parameters", caresParameters); err != nil {
			log.Printf("[CRITAL]%s provider set caresParameters fail, reason:%s\n ", logId, err.Error())
		}
	}
	tags, err := mysqlService.DescribeTagsOfInstanceId(ctx, d.Id())
	if err != nil {
		errRet = err
		return
	}
	if err := d.Set("tags", tags); err != nil {
		log.Printf("[CRITAL]%s provider set tags fail, reason:%s\n ", logId, err.Error())
	}

	d.Set("intranet_ip", *mysqlInfo.Vip)
	d.Set("intranet_port", int(*mysqlInfo.Vport))

	if *mysqlInfo.CdbError != 0 {
		d.Set("locked", 1)
	} else {
		d.Set("locked", 0)
	}
	d.Set("status", *mysqlInfo.Status)
	d.Set("task_status", *mysqlInfo.TaskStatus)
	return
}

func resourceTencentCloudMysqlInstanceRead(d *schema.ResourceData, meta interface{}) error {

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	mysqlInfo, err := tencentMsyqlBasicInfoRead(ctx, d, meta)
	if err != nil {
		return err
	}
	if mysqlInfo == nil {
		d.SetId("")
		return nil
	}
	d.Set("availability_zone", *mysqlInfo.Zone)

	backConfig, err := mysqlService.DescribeDBInstanceConfig(ctx, d.Id())
	if err != nil {
		return err
	}
	d.Set("slave_sync_mode", int(*backConfig.Response.ProtectMode))
	//	d.Set("slave_deploy_mode", int(*backConfig.Response.DeployMode))

	//	if backConfig.Response.SlaveConfig != nil && *backConfig.Response.SlaveConfig.Zone != "" {
	//		d.Set("first_slave_zone", *backConfig.Response.SlaveConfig.Zone)
	//	}
	//	if backConfig.Response.BackupConfig != nil && *backConfig.Response.BackupConfig.Zone != "" {
	//		d.Set("second_slave_zone", *backConfig.Response.BackupConfig.Zone)
	//	}

	return nil
}
func resourceTencentCloudMysqlInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceTencentCloudMysqlInstanceRead(d, meta)
}
func resourceTencentCloudMysqlInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
