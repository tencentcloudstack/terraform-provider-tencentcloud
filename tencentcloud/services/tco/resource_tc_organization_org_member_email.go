package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationOrgMemberEmail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgMemberEmailCreate,
		Read:   resourceTencentCloudOrganizationOrgMemberEmailRead,
		Update: resourceTencentCloudOrganizationOrgMemberEmailUpdate,
		Delete: resourceTencentCloudOrganizationOrgMemberEmailDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"member_uin": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Member Uin.",
			},

			"email": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Email address.",
			},

			"country_code": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "International region.",
			},

			"phone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Phone number.",
			},

			"bind_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Binding IDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"apply_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Application timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"bind_status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Binding status is not binding: unbound, to be activated: value, successful binding: success, binding failure: failedNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"bind_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Binding timeNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "FailedNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"phone_bind": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Safe mobile phone binding state is not bound: 0, has been binded: 1Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_email.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request   = organization.NewAddOrganizationMemberEmailRequest()
		response  = organization.NewAddOrganizationMemberEmailResponse()
		bindId    uint64
		memberUin int64
	)
	if v, ok := d.GetOkExists("member_uin"); ok {
		memberUin = int64(v.(int))
		request.MemberUin = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone"); ok {
		request.Phone = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddOrganizationMemberEmail(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgMemberEmail failed, reason:%+v", logId, err)
		return err
	}

	bindId = *response.Response.BindId
	d.SetId(helper.Int64ToStr(memberUin) + tccommon.FILED_SP + helper.UInt64ToStr(bindId))

	return resourceTencentCloudOrganizationOrgMemberEmailRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_email.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	memberUin := idSplit[0]
	bindId := idSplit[1]

	orgMemberEmail, err := service.DescribeOrganizationOrgMemberEmailById(ctx, helper.StrToInt64(memberUin), helper.StrToUInt64(bindId))
	if err != nil {
		return err
	}

	if orgMemberEmail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMemberEmail` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("member_uin", helper.StrToInt64(memberUin))

	if orgMemberEmail.Email != nil {
		_ = d.Set("email", orgMemberEmail.Email)
	}

	if orgMemberEmail.CountryCode != nil {
		_ = d.Set("country_code", orgMemberEmail.CountryCode)
	}

	if orgMemberEmail.Phone != nil {
		_ = d.Set("phone", orgMemberEmail.Phone)
	}

	if orgMemberEmail.BindId != nil {
		_ = d.Set("bind_id", orgMemberEmail.BindId)
	}

	if orgMemberEmail.ApplyTime != nil {
		_ = d.Set("apply_time", orgMemberEmail.ApplyTime)
	}

	if orgMemberEmail.BindStatus != nil {
		_ = d.Set("bind_status", orgMemberEmail.BindStatus)
	}

	if orgMemberEmail.BindTime != nil {
		_ = d.Set("bind_time", orgMemberEmail.BindTime)
	}

	if orgMemberEmail.Description != nil {
		_ = d.Set("description", orgMemberEmail.Description)
	}

	if orgMemberEmail.PhoneBind != nil {
		_ = d.Set("phone_bind", orgMemberEmail.PhoneBind)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgMemberEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_email.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := organization.NewUpdateOrganizationMemberEmailBindRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	memberUin := idSplit[0]
	bindId := idSplit[1]

	immutableArgs := []string{"member_uin", "bind_id", "apply_time", "bind_status", "bind_time", "description", "phone_bind"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	request.MemberUin = helper.StrToInt64Point(memberUin)
	request.BindId = helper.StrToInt64Point(bindId)

	if v, ok := d.GetOk("email"); ok {
		request.Email = helper.String(v.(string))
	}

	if v, ok := d.GetOk("country_code"); ok {
		request.CountryCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("phone"); ok {
		request.Phone = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdateOrganizationMemberEmailBind(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update organization orgMemberEmail failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOrganizationOrgMemberEmailRead(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberEmailDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_member_email.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
