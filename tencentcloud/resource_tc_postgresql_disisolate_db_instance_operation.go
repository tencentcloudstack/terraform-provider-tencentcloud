/*
Provides a resource to create a postgresql disisolate_db_instance_operation

Example Usage

```hcl
resource "tencentcloud_postgresql_disisolate_db_instance_operation" "disisolate_db_instance_operation" {
  db_instance_id_set = [local.pgsql_id]
  period = 1
  auto_voucher = false
}
```
*/
package tencentcloud

import (
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlDisisolateDbInstanceOperationCreate,
		Read:   resourceTencentCloudPostgresqlDisisolateDbInstanceOperationRead,
		Delete: resourceTencentCloudPostgresqlDisisolateDbInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id_set": {
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
		},
	}
}

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request         = postgresql.NewDisIsolateDBInstancesRequest()
		ids             []string
		firstInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id_set"); ok {
		dBInstanceIdSetSet := v.(*schema.Set).List()
		for i := range dBInstanceIdSetSet {
			if dBInstanceIdSetSet[i] != nil {
				dBInstanceIdSet := dBInstanceIdSetSet[i].(string)
				request.DBInstanceIdSet = append(request.DBInstanceIdSet, &dBInstanceIdSet)
				ids = append(ids, dBInstanceIdSet)
			}
		}
		firstInstanceId = dBInstanceIdSetSet[0].(string)
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
			if voucherIdsSet[i] != nil {
				voucherIds := voucherIdsSet[i].(string)
				request.VoucherIds = append(request.VoucherIds, &voucherIds)
			}
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UsePostgresqlClient().DisIsolateDBInstances(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql DisisolateDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(ids, FILED_SP))

	service := PostgresqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"running"}, 10*readRetryTimeout, 10*time.Second, service.PostgresqlDbInstanceOperationStateRefreshFunc(firstInstanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresqlDisisolateDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
