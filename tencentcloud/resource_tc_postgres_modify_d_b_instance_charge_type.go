/*
Provides a resource to create a postgres modify_d_b_instance_charge_type

Example Usage

```hcl
resource "tencentcloud_postgres_modify_d_b_instance_charge_type" "modify_d_b_instance_charge_type" {
  d_b_instance_id = "postgres-6r233v55"
  instance_charge_type = "PREPAID"
  period = 1
  auto_renew_flag = 0
  auto_voucher = 0
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres modify_d_b_instance_charge_type can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_modify_d_b_instance_charge_type.modify_d_b_instance_charge_type modify_d_b_instance_charge_type_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgres "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudPostgresModifyDBInstanceChargeType() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresModifyDBInstanceChargeTypeCreate,
		Read:   resourceTencentCloudPostgresModifyDBInstanceChargeTypeRead,
		Delete: resourceTencentCloudPostgresModifyDBInstanceChargeTypeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "DbInstance ID.",
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

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresModifyDBInstanceChargeTypeCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_charge_type.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewModifyDBInstanceChargeTypeRequest()
		response     = postgres.NewModifyDBInstanceChargeTypeResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
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
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().ModifyDBInstanceChargeType(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres ModifyDBInstanceChargeType failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresModifyDBInstanceChargeTypeRead(d, meta)
}

func resourceTencentCloudPostgresModifyDBInstanceChargeTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_charge_type.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresModifyDBInstanceChargeTypeDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_modify_d_b_instance_charge_type.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
