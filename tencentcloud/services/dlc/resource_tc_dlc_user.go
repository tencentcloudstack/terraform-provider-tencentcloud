package dlc

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUserCreate,
		Read:   resourceTencentCloudDlcUserRead,
		Update: resourceTencentCloudDlcUserUpdate,
		Delete: resourceTencentCloudDlcUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_id": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				Description: "Sub-user UIN that needs to be granted permissions. It can be checked through the upper right corner of Tencent Cloud Console -> Account Information -> Account ID.",
			},

			"user_description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User description, which can make it easy to identify different users.",
			},

			"user_type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Types of users. ADMIN: administrators; COMMON: general users. When the type of user is administrator, the collections of permissions and bound working groups cannot be set. Administrators own all the permissions by default. If the parameter is not filled in, it will be COMMON by default.",
			},

			"user_alias": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "User alias, and its characters are less than 50.",
			},

			"work_group_ids": {
				Computed:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Collection of IDs of working groups bound to users.",
			},
		},
	}
}

func resourceTencentCloudDlcUserCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dlc.NewCreateUserRequest()
		userId  string
	)

	if v, ok := d.GetOk("user_id"); ok {
		userId = v.(string)
		request.UserId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_description"); ok {
		request.UserDescription = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_type"); ok {
		request.UserType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("user_alias"); ok {
		request.UserAlias = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().CreateUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create dlc user failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create dlc user failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(userId)
	return resourceTencentCloudDlcUserRead(d, meta)
}

func resourceTencentCloudDlcUserRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId  = d.Id()
	)

	user, err := service.DescribeDlcUserById(ctx, userId)
	if err != nil {
		return err
	}

	if user == nil {
		log.Printf("[WARN]%s resource `DlcUser` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if user.UserId != nil {
		_ = d.Set("user_id", user.UserId)
	}

	if user.UserDescription != nil {
		_ = d.Set("user_description", user.UserDescription)
	}

	if user.UserType != nil {
		_ = d.Set("user_type", user.UserType)
	}

	if user.UserAlias != nil {
		_ = d.Set("user_alias", user.UserAlias)
	}

	if user.WorkGroupSet != nil {
		workGroups := make([]*int64, len(user.WorkGroupSet))
		for _, workGroup := range user.WorkGroupSet {
			workGroups = append(workGroups, workGroup.WorkGroupId)
		}

		_ = d.Set("work_group_ids", workGroups)
	}

	return nil
}

func resourceTencentCloudDlcUserUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = dlc.NewModifyUserRequest()
		userId  = d.Id()
	)

	immutableArgs := []string{"user_type", "user_alias"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("user_description") {
		if v, ok := d.GetOk("user_description"); ok {
			request.UserDescription = helper.String(v.(string))
		}

		request.UserId = &userId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().ModifyUser(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update dlc user failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudDlcUserRead(d, meta)
}

func resourceTencentCloudDlcUserDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_user.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		userId  = d.Id()
	)

	if err := service.DeleteDlcUserById(ctx, userId); err != nil {
		return err
	}

	return nil
}
