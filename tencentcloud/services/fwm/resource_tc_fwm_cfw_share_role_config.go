package fwm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	fwm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/fwm/v20250611"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudFwmCfwShareRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudFwmCfwShareRoleConfigCreate,
		Read:   resourceTencentCloudFwmCfwShareRoleConfigRead,
		Update: resourceTencentCloudFwmCfwShareRoleConfigUpdate,
		Delete: resourceTencentCloudFwmCfwShareRoleConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"member_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Target member ID.",
			},
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Target member AppId, used to locate the member in Read.",
			},
			"cfw_share_role": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Cloud firewall share role, valid values: sharer, user, none.",
			},
			"cfw_sharer_app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Sharer AppId, only available when the role is user.",
			},
		},
	}
}

func resourceTencentCloudFwmCfwShareRoleConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	memberId := d.Get("member_id").(string)
	appId := d.Get("app_id").(string)

	d.SetId(appId + tccommon.FILED_SP + memberId)

	return resourceTencentCloudFwmCfwShareRoleConfigUpdate(d, meta)
}

func resourceTencentCloudFwmCfwShareRoleConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	appId := idSplit[0]
	memberId := idSplit[1]

	organMember, err := service.DescribeFwmOrganMemberById(ctx, appId, memberId)
	if err != nil {
		return err
	}

	if organMember == nil {
		log.Printf("[CRUD] fwm_cfw_share_role id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if organMember.CfwShareRole != nil {
		_ = d.Set("cfw_share_role", organMember.CfwShareRole)
	}

	if organMember.CfwSharerAppId != nil {
		_ = d.Set("cfw_sharer_app_id", organMember.CfwSharerAppId)
	}

	if organMember.MemberId != nil {
		_ = d.Set("member_id", organMember.MemberId)
	}

	if organMember.AppId != nil {
		_ = d.Set("app_id", organMember.AppId)
	}

	return nil
}

func resourceTencentCloudFwmCfwShareRoleConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = FwmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	appId := idSplit[0]
	memberId := idSplit[1]

	target := d.Get("cfw_share_role").(string)

	organMember, err := service.DescribeFwmOrganMemberById(ctx, appId, memberId)
	if err != nil {
		return err
	}

	var current string
	if organMember != nil && organMember.CfwShareRole != nil {
		current = *organMember.CfwShareRole
	}

	if current == target {
		return resourceTencentCloudFwmCfwShareRoleConfigRead(d, meta)
	}

	needRelay := isSharerUser(current) && isSharerUser(target)

	if needRelay {
		if err := modifyCfwShareRole(ctx, meta, logId, memberId, "none"); err != nil {
			return err
		}
	}

	if err := modifyCfwShareRole(ctx, meta, logId, memberId, target); err != nil {
		return err
	}

	return resourceTencentCloudFwmCfwShareRoleConfigRead(d, meta)
}

func modifyCfwShareRole(ctx context.Context, meta interface{}, logId, memberId, role string) error {
	request := fwm.NewModifyCfwShareRoleRequest()
	request.CfwShareRole = helper.String(role)
	request.UserMemberId = helper.String(memberId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseFwmClient().ModifyCfwShareRoleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify fwm cfw share role to %s failed, reason:%+v", logId, role, err)
		return err
	}

	return nil
}

func isSharerUser(role string) bool {
	return role == "sharer" || role == "user"
}

func resourceTencentCloudFwmCfwShareRoleConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_fwm_cfw_share_role.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
