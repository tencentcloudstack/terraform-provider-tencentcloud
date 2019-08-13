/*
Provides a mysql instance resource to create read-only database instances.

~> **NOTE:** The terminate operation of read only mysql does NOT take effect immediately，maybe takes for several hours. so during that time, VPCs associated with that mysql instance can't be terminated also.

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
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func resourceTencentCloudMysqlReadonlyInstance() *schema.Resource {
	readonlyInstanceInfo := map[string]*schema.Schema{
		"master_instance_id": {
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			Description: "Indicates the master instance ID of recovery instances.",
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

	period := int64(d.Get("period").(int))
	request.Period = &period

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

	if err := mysqlAllInstanceRoleSet(ctx, request, d, meta); err != nil {
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

func resourceTencentCloudMysqlReadonlyInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	// the mysql master instance must have a backup before creating a read-only instance
	masterInstanceId := d.Get("master_instance_id").(string)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		backups, err := mysqlService.DescribeBackupsByMysqlId(ctx, masterInstanceId, 10)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(backups) < 1 {
			return resource.RetryableError(fmt.Errorf("waiting backup creating"))
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mysql  task fail, reason:%s\n ", logId, err.Error())
		return err
	}

	payType := d.Get("pay_type").(int)
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

	err = resource.Retry(10*time.Minute, func() *resource.RetryError {
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

	return resourceTencentCloudMysqlReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudMysqlReadonlyInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlInfo, err := tencentMsyqlBasicInfoRead(ctx, d, meta, false)
	if err != nil {
		return err
	}
	if mysqlInfo == nil {
		d.SetId("")
		return nil
	}
	d.Set("master_instance_id", *mysqlInfo.MasterInfo.InstanceId)

	return nil
}

func resourceTencentCloudMysqlReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	payType := d.Get("pay_type").(int)

	d.Partial(true)

	err := mysqlAllInstanceRoleUpdate(ctx, d, meta)
	if err != nil {
		return err
	}

	if payType == MysqlPayByMonth {
		if d.HasChange("auto_renew_flag") {
			renewFlag := int64(d.Get("auto_renew_flag").(int))
			mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
			if err := mysqlService.ModifyAutoRenewFlag(ctx, d.Id(), renewFlag); err != nil {
				return err
			}
			d.SetPartial("auto_renew_flag")
		}
	}

	d.Partial(false)

	return resourceTencentCloudMysqlReadonlyInstanceRead(d, meta)
}

func resourceTencentCloudMysqlReadonlyInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_readonly_instance.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	_, err := mysqlService.IsolateDBInstance(ctx, d.Id())
	if err != nil {
		return err
	}
	var hasDeleted = false

	err = resource.Retry(20*time.Minute, func() *resource.RetryError {
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

	err = mysqlService.OfflineIsolatedInstances(ctx, d.Id())
	if err == nil {
		log.Printf("[WARN]this mysql is readonly instance, it is released asynchronously, and the bound resource is not now fully released now\n")
	}

	return err
}
