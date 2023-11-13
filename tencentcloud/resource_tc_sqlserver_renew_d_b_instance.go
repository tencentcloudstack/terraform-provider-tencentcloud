/*
Provides a resource to create a sqlserver renew_d_b_instance

Example Usage

```hcl
resource "tencentcloud_sqlserver_renew_d_b_instance" "renew_d_b_instance" {
  instance_id = "mssql-i1z41iwd"
  period = 1
  auto_voucher = 1
  voucher_ids =
  auto_renew_flag = 0
}
```

Import

sqlserver renew_d_b_instance can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_renew_d_b_instance.renew_d_b_instance renew_d_b_instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"log"
)

func resourceTencentCloudSqlserverRenewDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSqlserverRenewDBInstanceCreate,
		Read:   resourceTencentCloudSqlserverRenewDBInstanceRead,
		Update: resourceTencentCloudSqlserverRenewDBInstanceUpdate,
		Delete: resourceTencentCloudSqlserverRenewDBInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID.",
			},

			"period": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "How many months to renew, the value range is 1-48, the default is 1.",
			},

			"auto_voucher": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Whether to use the voucher automatically, 0-not use; 1-use; the default is not practical.",
			},

			"voucher_ids": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID array, currently only supports the use of 1 voucher.",
			},

			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Renewal flag 0: Normal renewal 1: Automatic renewal: only valid when subscribing from pay-as-you-go to yearly or monthly.",
			},
		},
	}
}

func resourceTencentCloudSqlserverRenewDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_d_b_instance.create")()
	defer inconsistentCheck(d, meta)()

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	d.SetId(instanceId)

	return resourceTencentCloudSqlserverRenewDBInstanceUpdate(d, meta)
}

func resourceTencentCloudSqlserverRenewDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_d_b_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	renewDBInstanceId := d.Id()

	renewDBInstance, err := service.DescribeSqlserverRenewDBInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if renewDBInstance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `SqlserverRenewDBInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if renewDBInstance.InstanceId != nil {
		_ = d.Set("instance_id", renewDBInstance.InstanceId)
	}

	if renewDBInstance.Period != nil {
		_ = d.Set("period", renewDBInstance.Period)
	}

	if renewDBInstance.AutoVoucher != nil {
		_ = d.Set("auto_voucher", renewDBInstance.AutoVoucher)
	}

	if renewDBInstance.VoucherIds != nil {
		_ = d.Set("voucher_ids", renewDBInstance.VoucherIds)
	}

	if renewDBInstance.AutoRenewFlag != nil {
		_ = d.Set("auto_renew_flag", renewDBInstance.AutoRenewFlag)
	}

	return nil
}

func resourceTencentCloudSqlserverRenewDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_d_b_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := sqlserver.NewRenewDBInstanceRequest()

	renewDBInstanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"instance_id", "period", "auto_voucher", "voucher_ids", "auto_renew_flag"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSqlserverClient().RenewDBInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update sqlserver renewDBInstance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudSqlserverRenewDBInstanceRead(d, meta)
}

func resourceTencentCloudSqlserverRenewDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_sqlserver_renew_d_b_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
