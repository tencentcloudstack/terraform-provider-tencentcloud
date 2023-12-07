package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlSwitchForUpgrade() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlSwitchForUpgradeCreate,
		Read:   resourceTencentCloudMysqlSwitchForUpgradeRead,
		Delete: resourceTencentCloudMysqlSwitchForUpgradeDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of cdb-c1nl9rpv. It is the same as the instance ID displayed on the TencentDB Console page.",
			},
		},
	}
}

func resourceTencentCloudMysqlSwitchForUpgradeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_for_upgrade.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		request    = mysql.NewSwitchForUpgradeRequest()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().SwitchForUpgrade(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s mysql switchForUpgrade failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	err = resource.Retry(7*readRetryTimeout, func() *resource.RetryError {
		mysqlInfo, err := service.DescribeDBInstanceById(ctx, instanceId)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if mysqlInfo == nil {
			err = fmt.Errorf("mysqlid %s instance not exists", instanceId)
			return resource.NonRetryableError(err)
		}
		if *mysqlInfo.Status == MYSQL_STATUS_DELIVING {
			return resource.RetryableError(fmt.Errorf("mysql switchForUpgrade status is MYSQL_STATUS_DELIVING(%d)", MYSQL_STATUS_DELIVING))
		}
		if *mysqlInfo.Status == MYSQL_STATUS_RUNNING {
			return nil
		}
		err = fmt.Errorf("mysql switchForUpgrade status is %v,we won't wait for it finish", *mysqlInfo.Status)
		return resource.NonRetryableError(err)
	})

	if err != nil {
		log.Printf("[CRITAL]%s mysql switchForUpgrade fail, reason:%s\n ", logId, err.Error())
		return err
	}

	return resourceTencentCloudMysqlSwitchForUpgradeRead(d, meta)
}

func resourceTencentCloudMysqlSwitchForUpgradeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_for_upgrade.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	switchForUpgrade, err := service.DescribeDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if switchForUpgrade == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MysqlSwitchForUpgrade` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if switchForUpgrade.InstanceId != nil {
		_ = d.Set("instance_id", switchForUpgrade.InstanceId)
	}

	return nil
}

func resourceTencentCloudMysqlSwitchForUpgradeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_switch_for_upgrade.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
