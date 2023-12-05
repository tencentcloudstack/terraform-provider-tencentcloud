package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudDlcModifyUserTypOperation() *schema.Resource {
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
	defer logElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

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

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseDlcClient().ModifyUserType(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcModifyUserTypOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dlc_modify_user_typ_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
