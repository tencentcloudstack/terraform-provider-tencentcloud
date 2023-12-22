package postgresql

import (
	"log"
	"strings"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlDisisolateDbInstanceOperation() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().DisIsolateDBInstances(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql DisisolateDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join(ids, tccommon.FILED_SP))

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 10*tccommon.ReadRetryTimeout, 10*time.Second, service.PostgresqlDbInstanceOperationStateRefreshFunc(firstInstanceId, []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresqlDisisolateDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlDisisolateDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_disisolate_db_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
