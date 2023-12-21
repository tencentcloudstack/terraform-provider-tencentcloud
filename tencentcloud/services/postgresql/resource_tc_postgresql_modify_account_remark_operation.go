package postgresql

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudPostgresqlModifyAccountRemarkOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudPostgresqlModifyAccountRemarkOperationCreate,
		Read:   resourceTencentCloudPostgresqlModifyAccountRemarkOperationRead,
		Delete: resourceTencentCloudPostgresqlModifyAccountRemarkOperationDelete,
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance ID in the format of postgres-4wdeb0zv.",
			},

			"user_name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Instance username.",
			},

			"remark": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "New remarks corresponding to user `UserName`.",
			},
		},
	}
}

func resourceTencentCloudPostgresqlModifyAccountRemarkOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_account_remark_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request      = postgresql.NewModifyAccountRemarkRequest()
		dBInstanceId string
	)
	if v, ok := d.GetOk("db_instance_id"); ok {
		request.DBInstanceId = helper.String(v.(string))
		dBInstanceId = v.(string)
	}

	if v, ok := d.GetOk("user_name"); ok {
		request.UserName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("remark"); ok {
		request.Remark = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UsePostgresqlClient().ModifyAccountRemark(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate postgresql ModifyAccountRemarkOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(dBInstanceId)

	return resourceTencentCloudPostgresqlModifyAccountRemarkOperationRead(d, meta)
}

func resourceTencentCloudPostgresqlModifyAccountRemarkOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_account_remark_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudPostgresqlModifyAccountRemarkOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_postgresql_modify_account_remark_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
