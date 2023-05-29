/*
Provides a resource to create a postgresql modify_db_instance_charge_type_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_modify_db_instance_charge_type_operation" "modify_db_instance_charge_type_operation" {
  db_instance_id = "postgres-6r233v55"
  instance_charge_type = "PREPAID"
  period = 1
  auto_renew_flag = 0
  auto_voucher = 0
}
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationCreate,
		Read:   resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationRead,
		Delete: resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "dbInstance ID.",
			},

			"instance_charge_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance billing mode. Valid values:PREPAID (monthly subscription), POSTPAID_BY_HOUR (pay-as-you-go).",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Valid period in months of purchased instances. Valid values:1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. This parameter is set to 1 when the pay-as-you-go billing mode is used.",
			},

			"auto_renew_flag": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal flag. Valid values:0 (manual renewal), 1 (auto-renewal). Default value:0.",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically use vouchers.Valid values:1(yes),0(no).Default value:0.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_modify_db_instance_charge_type_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgresql.NewModifyDBInstanceChargeTypeRequest()
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("instance_charge_type"); ok {
		request.InstanceChargeType = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_renew_flag"); v != nil {
		request.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().ModifyDBInstanceChargeType(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql ModifyDbInstanceChargeTypeOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_modify_db_instance_charge_type_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlModifyDbInstanceChargeTypeOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_modify_db_instance_charge_type_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
