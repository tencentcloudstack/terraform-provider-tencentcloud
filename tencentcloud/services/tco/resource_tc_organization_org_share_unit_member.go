package tco

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationOrgShareUnitMember() *schema.Resource {
	return &schema.Resource{
		DeprecationMessage: "This resource has been deprecated in Terraform TencentCloud provider version 1.82.28, Please use `tencentcloud_organization_org_share_unit_member_v2` instead.",
		Create:             resourceTencentCloudOrganizationOrgShareUnitMemberCreate,
		Read:               resourceTencentCloudOrganizationOrgShareUnitMemberRead,
		Delete:             resourceTencentCloudOrganizationOrgShareUnitMemberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"unit_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Shared unit ID.",
			},

			"area": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Shared unit region.",
			},

			"members": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    10,
				Description: "Shared member list. Up to 10 items.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_member_uin": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Member uin.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgShareUnitMemberCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = organization.NewAddShareUnitMembersRequest()
		unitId         string
		area           string
		shareMemberUin []string
	)
	if v, ok := d.GetOk("unit_id"); ok {
		unitId = v.(string)
		request.UnitId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("area"); ok {
		area = v.(string)
		request.Area = helper.String(v.(string))
	}

	if v, ok := d.GetOk("members"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			shareMember := organization.ShareMember{}
			if v, ok := dMap["share_member_uin"]; ok {
				shareMemberUin = append(shareMemberUin, helper.IntToStr(v.(int)))
				shareMember.ShareMemberUin = helper.IntInt64(v.(int))
			}
			request.Members = append(request.Members, &shareMember)
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddShareUnitMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgShareUnitMember failed, reason:%+v", logId, err)
		return err
	}
	shareMemberUins := strings.Join(shareMemberUin, tccommon.COMMA_SP)
	d.SetId(strings.Join(append([]string{area, unitId, shareMemberUins}), tccommon.FILED_SP))

	return resourceTencentCloudOrganizationOrgShareUnitMemberRead(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitMemberRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	area := idSplit[0]
	unitId := idSplit[1]
	shareMemberUin := idSplit[2]

	orgShareUnitMember, err := service.DescribeOrganizationOrgShareUnitMemberById(ctx, unitId, area, shareMemberUin)
	if err != nil {
		return err
	}

	if len(orgShareUnitMember) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgShareUnitMember` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("unit_id", unitId)
	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitMemberDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	area := idSplit[0]
	unitId := idSplit[1]
	shareMemberUin := idSplit[2]

	if err := service.DeleteOrganizationOrgShareUnitMemberById(ctx, unitId, area, shareMemberUin); err != nil {
		return err
	}

	return nil
}
