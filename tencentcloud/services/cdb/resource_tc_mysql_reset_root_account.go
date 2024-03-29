package cdb

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mysql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMysqlResetRootAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMysqlResetRootAccountCreate,
		Read:   resourceTencentCloudMysqlResetRootAccountRead,
		Delete: resourceTencentCloudMysqlResetRootAccountDelete,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The instance ID, in the format: cdb-c1nl9rpv, is the same as the instance ID displayed on the cloud database console page.",
			},
		},
	}
}

func resourceTencentCloudMysqlResetRootAccountCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_reset_root_account.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mysql.NewResetRootAccountRequest()
		instanceId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
		request.InstanceId = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMysqlClient().ResetRootAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate mysql resetRootAccount failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(instanceId)

	return resourceTencentCloudMysqlResetRootAccountRead(d, meta)
}

func resourceTencentCloudMysqlResetRootAccountRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_reset_root_account.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMysqlResetRootAccountDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mysql_reset_root_account.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
