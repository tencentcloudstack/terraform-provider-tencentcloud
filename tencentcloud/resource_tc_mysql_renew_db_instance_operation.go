/*
Provides a resource to create a mysql renew_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_mysql_renew_db_instance_operation" "renew_db_instance_operation" {
  instance_id = "cdb-c1nl9rpv"
  time_span = 1
  modify_pay_type = "PREPAID"
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
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudMysqlRenewDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlRenewDbInstanceOperationCreate,
		Read:   resourceTencentCloudMysqlRenewDbInstanceOperationRead,
		Delete: resourceTencentCloudMysqlRenewDbInstanceOperationDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID to be renewed, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, you can use [Query Instance List](https://cloud.tencent.com/document/api/236/ 15872).",
			},

			"time_span": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration, unit: month, optional values include [1,2,3,4,5,6,7,8,9,10,11,12,24,36].",
			},

			"modify_pay_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "If you need to renew the Pay-As-You-Go instance to a Subscription instance, the value of this input parameter needs to be specified as `PREPAID`.",
			},

			"deal_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Deal id.",
			},

			"deadline_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Instance expiration time.",
			},
		},
	}
}

func resourceTencentCloudMysqlRenewDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = mysql.NewRenewDBInstanceRequest()
		response   = mysql.NewRenewDBInstanceResponse()
		instanceId string
	)
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("time_span"); v != nil {
		request.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("modify_pay_type"); ok {
		request.ModifyPayType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMysqlClient().RenewDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql renewDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	dealId := *response.Response.DealId
	d.SetId(instanceId)
	_ = d.Set("deal_id", dealId)

	return resourceTencentCloudMysqlRenewDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudMysqlRenewDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}
	mysqlInfo, errRet := service.DescribeDBInstanceById(ctx, d.Id())
	if errRet != nil {
		return fmt.Errorf("Describe mysql instance fails, reaseon %s", errRet.Error())
	}

	if mysqlInfo == nil {
		d.SetId("")
		return nil
	}

	if mysqlInfo.DeadlineTime != nil {
		_ = d.Set("deadline_time", mysqlInfo.DeadlineTime)
	}

	return nil
}

func resourceTencentCloudMysqlRenewDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mysql_renew_db_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
