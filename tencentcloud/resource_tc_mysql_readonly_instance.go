package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func resourceTencentCloudMysqlReadonlyInstance() *schema.Resource {
	readonlyInstanceInfo := map[string]*schema.Schema{
		"master_instance_id": {
			Type:     schema.TypeString,
			ForceNew: true,
			Required: true,
		},
	}

	basic := TencentMsyqlBasicInfo()
	for k, v := range basic {
		readonlyInstanceInfo[k] = v
	}

	return &schema.Resource{
		Create: resourceTencentCloudMysqlReadonlyInstanceCreate,
		Read:   resourceTencentCloudMysqlReadonlyInstanceRead,
		Update: resourceTencentCloudMysqlReadonlyInstanceUpdate,
		Delete: resourceTencentCloudMysqlReadonlyInstanceDelete,

		Schema: readonlyInstanceInfo,
	}
}

func mysqlCreateReadonlyInstancePayByMonth(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(ctx)

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
	logId := GetLogId(ctx)

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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlInfo, err := tencentMsyqlBasicInfoRead(ctx, d, meta)
	if err != nil {
		return err
	}
	d.Set("master_instance_id", *mysqlInfo.MasterInfo.InstanceId)

	return nil
}

func resourceTencentCloudMysqlReadonlyInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
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
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	_, err := mysqlService.IsolateDBInstance(ctx, d.Id())

	if err != nil {
		return err
	}
	return nil
}
