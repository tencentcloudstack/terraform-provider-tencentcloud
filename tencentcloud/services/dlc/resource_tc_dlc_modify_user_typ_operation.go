package dlc

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcModifyUserTypOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcModifyUserTypOperationCreate,
		Read:   resourceTencentCloudDlcModifyUserTypOperationRead,
		Delete: resourceTencentCloudDlcModifyUserTypOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "User id (uin), if left blank, it defaults to the caller's sub-uin.",
			},

			"user_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "User type, only support: ADMIN: ddministrator/COMMON: ordinary user.",
			},
		},
	}
}

func resourceTencentCloudDlcModifyUserTypOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request = dlc.NewModifyUserTypeRequest()
		userId  string
	)
	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		request.UserType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().ModifyUserType(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate dlc modifyUserTypOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)

	return resourceTencentCloudDlcModifyUserTypOperationRead(d, meta)
}

func resourceTencentCloudDlcModifyUserTypOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcModifyUserTypOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
