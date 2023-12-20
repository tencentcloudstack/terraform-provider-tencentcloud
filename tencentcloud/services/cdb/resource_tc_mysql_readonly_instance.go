package cdb

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlReadonlyInstance() *schema.Resource {
	readonlyInstanceInfo := map[string]*schema.Schema{
		"master_instance_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Indicates the master instance ID of recovery instances.",
		},
		"zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "Zone information, this parameter defaults to, the system automatically selects an Availability Zone.",
		},
		"master_region": {
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			Description: "The zone information of the primary instance is required when you purchase a disaster recovery instance.",
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		readonlyInstanceInfo[k] = v
	}
	delete(readonlyInstanceInfo, "gtid")

	return &schema.Resource{
		Create: resourceTencentCloudMysqlReadonlyInstanceCreate,
		Read:   resourceTencentCloudMysqlReadonlyInstanceRead,
		Update: resourceTencentCloudMysqlReadonlyInstanceUpdate,
		Delete: resourceTencentCloudMysqlReadonlyInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"prepaid_period": 1,
				"force_delete":   false,
			}),
		},
		Schema: readonlyInstanceInfo,
	}
}

func mysqlCreateReadonlyInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := cdb.NewCreateDBInstanceRequest()
	instanceRole := "ro"
	request.InstanceRole = &instanceRole

	payType, ok := d.GetOk("pay_type")
	var period int
	if !ok || payType == -1 {
		period = d.Get("prepaid_period").(int)
	} else {
		period = d.Get("period").(int)
	}
	request.Period = helper.IntInt64(period)

	if v, ok := d.GetOk("mem_size"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("volume_size"); ok {
		request.Volume = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("zone"); ok {
		zone := v.(string)
		request.Zone = &zone
	}
	if v, ok := d.GetOk("master_region"); ok {
		masterRegion := v.(string)
		request.MasterRegion = &masterRegion
	}

	if v, ok := d.GetOk("device_type"); ok {
		request.DeviceType = helper.String(v.(string))
	}

	autoRenewFlag := int64(d.Get("auto_renew_flag").(int))
	request.AutoRenewFlag = &autoRenewFlag

	masterInstanceId := d.Get("master_instance_id").(string)
	request.MasterInstanceId = &masterInstanceId

	// readonly group is not currently supported
	defaultRoGroupMode := "allinone"
	request.RoGroup = &cdb.RoGroup{RoGroupMode: &defaultRoGroupMode}

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
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
	return nil
}

func mysqlCreateReadonlyInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := tccommon.GetLogId(ctx)

	request := cdb.NewCreateDBInstanceHourRequest()
	instanceRole := "ro"
	request.InstanceRole = &instanceRole

	masterInstanceId := d.Get("master_instance_id").(string)
	request.MasterInstanceId = &masterInstanceId

	// readonly group is not currently supported
	defaultRoGroupMode := "allinone"
	request.RoGroup = &cdb.RoGroup{RoGroupMode: &defaultRoGroupMode}

	if v, ok := d.GetOk("mem_size"); ok {
		request.Memory = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("volume_size"); ok {
		request.Volume = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("cpu"); ok {
		request.Cpu = helper.IntInt64(v.(int))
	}

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
		return err
	}

	if v, ok := d.GetOk("zone"); ok {
		zone := v.(string)
		request.Zone = &zone
	}
	if v, ok := d.GetOk("master_region"); ok {
		masterRegion := v.(string)
		request.MasterRegion = &masterRegion
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
	return nil
}

func resourceTencentCloudMysqlReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_readonly_instance.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// the mysql master instance must have a backup before creating a read-only instance
	masterInstanceId := d.Get("master_instance_id").(string)

	err := resource.Retry(2*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		backups, err := mysqlService.DescribeBackupsByMysqlId(ctx, masterInstanceId, 10)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(backups) < 1 {
			return resource.RetryableError(fmt.Errorf("waiting backup creating"))
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	payType := getPayType(d).(int)
	if payType == MysqlPayByMonth {
		err := mysqlCreateReadonlyInstancePayByMonth(ctx, d, meta)
		if err != nil {
			return err
		}
	} else if payType == MysqlPayByUse {
		err := mysqlCreateReadonlyInstancePayByUse(ctx, d, meta)
		if err != nil {
			return err
		}
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
		tagService := &TagService{client: tcClient}
		resourceName := tccommon.BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		log.Printf("[DEBUG]Mysql instance create, resourceName:%s\n", resourceName)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudMysqlReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_readonly_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	mysqlService := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		mysqlInfo, e := tencentMsyqlBasicInfoRead(ctx, d, meta, false)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				return nil
			}
			return tccommon.RetryError(e)
		}
		if mysqlInfo == nil {
			d.SetId("")
			return nil
		}
		_ = d.Set("master_instance_id", *mysqlInfo.MasterInfo.InstanceId)
		_ = d.Set("zone", *mysqlInfo.Zone)
		_ = d.Set("master_region", *mysqlInfo.MasterInfo.Region)

		return nil
	})
	if err != nil {
		return fmt.Errorf("Fail to get basic info from mysql, reaseon %s", err.Error())
	}

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
	tagService := &TagService{client: tcClient}
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

	if *mysqlInfo.CdbError != 0 {
		_ = d.Set("locked", 1)
	} else {
		_ = d.Set("locked", 0)
	}
	_ = d.Set("status", mysqlInfo.Status)
	_ = d.Set("task_status", mysqlInfo.TaskStatus)

	return nil
}

func resourceTencentCloudMysqlReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_readonly_instance.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	payType := getPayType(d).(int)

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

	return resourceTencentCloudMysqlReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudMysqlReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_readonly_instance.delete")()

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
		log.Printf("[WARN]this mysql is readonly instance, it is released asynchronously, and the bound resource is not now fully released now\n")
	}
	return err
}
