/*
Provides a resource to create a postgres renew_instance

Example Usage

```hcl
resource "tencentcloud_postgres_renew_instance" "renew_instance" {
  d_b_instance_id = "postgres-6fego161"
  period = 12
  auto_voucher = 0
  voucher_ids = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres renew_instance can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_renew_instance.renew_instance renew_instance_id
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

func resourceTencentCloudPostgresRenewInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresRenewInstanceCreate,
		Read:   resourceTencentCloudPostgresRenewInstanceRead,
		Delete: resourceTencentCloudPostgresRenewInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of postgres-6fego161.",
			},

			"period": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Renewal duration in months.",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Whether to automatically use vouchers. 1:yes, 0:no. Default value:0.",
			},

			"voucher_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID list (only one voucher can be specified currently).",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresRenewInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_renew_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewRenewInstanceRequest()
		response     = postgres.NewRenewInstanceResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id"); ok {
		dBInstanceId = v.(string)
		request.DBInstanceId = helper.String(v.(string))
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().RenewInstance(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres RenewInstance failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresRenewInstanceRead(d, meta)
}

func resourceTencentCloudPostgresRenewInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_renew_instance.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresRenewInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_renew_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
