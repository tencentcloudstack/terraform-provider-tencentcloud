package postgresql

import (
	"log"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlRenewDbInstanceOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlRenewDbInstanceOperationCreate,
		Read:   resourceTencentCloudPostgresqlRenewDbInstanceOperationRead,
		Delete: resourceTencentCloudPostgresqlRenewDbInstanceOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
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
		},
	}
}

func resourceTencentCloudPostgresqlRenewDbInstanceOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_renew_db_instance_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewRenewInstanceRequest()
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
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
			if voucherIdsSet[i] != nil {
				voucherIds := voucherIdsSet[i].(string)
				request.VoucherIds = append(request.VoucherIds, &voucherIds)
			}
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().RenewInstance(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql RenewDbInstanceOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dBInstanceId)

	service := PostgresqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	conf := tccommon.BuildStateChangeConf([]string{}, []string{"running"}, 10*tccommon.ReadRetryTimeout, 5*time.Second, service.PostgresqlDbInstanceOperationStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudPostgresqlRenewDbInstanceOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlRenewDbInstanceOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_renew_db_instance_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlRenewDbInstanceOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_renew_db_instance_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
