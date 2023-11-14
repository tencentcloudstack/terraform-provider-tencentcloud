/*
Provides a resource to create a postgres dis_isolate_d_b_instances

Example Usage

```hcl
resource "tencentcloud_postgres_dis_isolate_d_b_instances" "dis_isolate_d_b_instances" {
  d_b_instance_id_set = &lt;nil&gt;
  period = 12
  auto_voucher = false
  voucher_ids = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}
```

Import

postgres dis_isolate_d_b_instances can be imported using the id, e.g.

```
terraform import tencentcloud_postgres_dis_isolate_d_b_instances.dis_isolate_d_b_instances dis_isolate_d_b_instances_id
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

func resourceTencentCloudPostgresDisIsolateDBInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresDisIsolateDBInstancesCreate,
		Read:   resourceTencentCloudPostgresDisIsolateDBInstancesRead,
		Delete: resourceTencentCloudPostgresDisIsolateDBInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"d_b_instance_id_set": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "List of resource IDs. Note that currently you cannot remove multiple instances from isolation at the same time. Only one instance ID can be passed in here.",
			},

			"period": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "The valid period (in months) of the monthly-subscribed instance when removing it from isolation.",
			},

			"auto_voucher": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "Whether to use vouchers. Valid values:true (yes), false (no). Default value:false.",
			},

			"voucher_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Voucher ID list.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudPostgresDisIsolateDBInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_dis_isolate_d_b_instances.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request      = postgres.NewDisIsolateDBInstancesRequest()
		response     = postgres.NewDisIsolateDBInstancesResponse()
		dBInstanceId string
	)
	if v, ok := d.GetOk("d_b_instance_id_set"); ok {
		dBInstanceIdSetSet := v.(*schema.Set).List()
		for i := range dBInstanceIdSetSet {
			dBInstanceIdSet := dBInstanceIdSetSet[i].(string)
			request.DBInstanceIdSet = append(request.DBInstanceIdSet, &dBInstanceIdSet)
		}
	}

	if v, _ := d.GetOk("period"); v != nil {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, _ := d.GetOk("auto_voucher"); v != nil {
		request.AutoVoucher = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("voucher_ids"); ok {
		voucherIdsSet := v.(*schema.Set).List()
		for i := range voucherIdsSet {
			voucherIds := voucherIdsSet[i].(string)
			request.VoucherIds = append(request.VoucherIds, &voucherIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresClient().DisIsolateDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgres DisIsolateDBInstances failed, reason:%+v", logId, err)
		return err
	}

	dBInstanceId = *response.Response.DBInstanceId
	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresDisIsolateDBInstancesRead(d, meta)
}

func resourceTencentCloudPostgresDisIsolateDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_dis_isolate_d_b_instances.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresDisIsolateDBInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgres_dis_isolate_d_b_instances.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
