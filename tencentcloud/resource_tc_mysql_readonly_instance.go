/*
Provides a mysql instance resource to create read-only database instances.

~> **NOTE:** Read-only instances can be purchased only for two-node or three-node source instances on MySQL 5.6 or above with the InnoDB engine at a specification of 1 GB memory and 50 GB disk capacity or above.
~> **NOTE:** The terminate operation of read only mysql does NOT take effect immediately, maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

Example Usage

```hcl
resource "tencentcloud_mysql_readonly_instance" "default" {
  master_instance_id = "cdb-dnqksd9f"
  instance_name      = "myTestMysql"
  mem_size           = 128000
  volume_size        = 255
  vpc_id             = "vpc-12mt3l31"
  subnet_id          = "subnet-9uivyb1g"
  intranet_port      = 3306
  security_groups    = ["sg-ot8eclwz"]

  tags = {
    name = "test"
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlReadonlyInstance() *schema.Resource {
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

		Schema: readonlyInstanceInfo,
	}
}

func mysqlCreateReadonlyInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(ctx)

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

func mysqlCreateReadonlyInstancePayByUse(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := getLogId(ctx)

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

func resourceTencentCloudMysqlReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	// the mysql master instance must have a backup before creating a read-only instance
	masterInstanceId := d.Get("master_instance_id").(string)

	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
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

	err = resource.Retry(4*readRetryTimeout, func() *resource.RetryError {
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
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cdb", "instanceId", tcClient.Region, d.Id())
		log.Printf("[DEBUG]Mysql instance create, resourceName:%s\n", resourceName)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudMysqlReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudMysqlReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, e := tencentMsyqlBasicInfoRead(ctx, d, meta, false)
		if e != nil {
			if mysqlService.NotFoundMysqlInstance(e) {
				d.SetId("")
				return nil
			}
			return retryError(e)
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
	return nil
}

func resourceTencentCloudMysqlReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	payType := getPayType(d).(int)

	d.Partial(true)

	if payType == MysqlPayByMonth {
		if d.HasChange("auto_renew_flag") {
			renewFlag := int64(d.Get("auto_renew_flag").(int))
			mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
			if err := mysqlService.ModifyAutoRenewFlag(ctx, d.Id(), renewFlag); err != nil {
				return err
			}

		}
	}
	err := mysqlAllInstanceRoleUpdate(ctx, d, meta)
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
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.delete")()

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
