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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// max length is 10
var batchSize = 10

func ResourceTencentCloudOrganizationOrgShareUnitMemberV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgShareUnitMemberV2Create,
		Read:   resourceTencentCloudOrganizationOrgShareUnitMemberV2Read,
		Update: resourceTencentCloudOrganizationOrgShareUnitMemberV2Update,
		Delete: resourceTencentCloudOrganizationOrgShareUnitMemberV2Delete,
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
				Type:        schema.TypeSet,
				Description: "Shared member list.",
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

func resourceTencentCloudOrganizationOrgShareUnitMemberV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		request = organization.NewAddShareUnitMembersRequest()
		unitId  string
		area    string
	)

	if v, ok := d.GetOk("unit_id"); ok {
		request.UnitId = helper.String(v.(string))
		unitId = v.(string)
	}

	if v, ok := d.GetOk("area"); ok {
		request.Area = helper.String(v.(string))
		area = v.(string)
	}

	orgShareUnitMembers := make([]*organization.ShareMember, 0, 10)
	if v, ok := d.GetOk("members"); ok {
		for _, item := range v.(*schema.Set).List() {
			if dMap, ok := item.(map[string]interface{}); ok {
				if v, ok := dMap["share_member_uin"]; ok {
					shareMember := organization.ShareMember{}
					shareMember.ShareMemberUin = helper.IntInt64(v.(int))
					orgShareUnitMembers = append(orgShareUnitMembers, &shareMember)
				}
			}
		}
	}

	for i := 0; i < len(orgShareUnitMembers); i += batchSize {
		end := i + batchSize
		if end > len(orgShareUnitMembers) {
			end = len(orgShareUnitMembers)
		}

		batch := orgShareUnitMembers[i:end]
		// clear Members value
		request.Members = nil
		for _, item := range batch {
			shareMember := organization.ShareMember{}
			shareMember.ShareMemberUin = item.ShareMemberUin
			request.Members = append(request.Members, &shareMember)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AddShareUnitMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s create organization share unit member failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(strings.Join([]string{unitId, area}, tccommon.FILED_SP))
	return resourceTencentCloudOrganizationOrgShareUnitMemberV2Read(d, meta)
}

func resourceTencentCloudOrganizationOrgShareUnitMemberV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	unitId := idSplit[0]
	area := idSplit[1]

	orgShareUnitMember, err := service.DescribeOrganizationOrgShareUnitMemberV2ById(ctx, unitId, area)
	if err != nil {
		return err
	}

	if len(orgShareUnitMember) < 1 {
		log.Printf("[WARN]%s resource `tencentcloud_organization_org_share_unit_member_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("unit_id", unitId)
	_ = d.Set("area", area)

	tmpList := make([]interface{}, 0, len(orgShareUnitMember))
	for _, item := range orgShareUnitMember {
		shareMember := map[string]interface{}{
			"share_member_uin": *item.ShareMemberUin,
		}

		tmpList = append(tmpList, shareMember)
	}

	_ = d.Set("members", tmpList)

	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitMemberV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	unitId := idSplit[0]
	area := idSplit[1]

	if d.HasChange("members") {
		oldInterface, newInterface := d.GetChange("members")
		oldInstances := oldInterface.(*schema.Set)
		newInstances := newInterface.(*schema.Set)
		remove := oldInstances.Difference(newInstances).List()
		add := newInstances.Difference(oldInstances).List()

		if len(add) > 0 {
			tmpList := make([]*organization.ShareUnitMember, 0, len(add))
			for _, item := range add {
				if dMap, ok := item.(map[string]interface{}); ok {
					if v, ok := dMap["share_member_uin"]; ok {
						shareMember := organization.ShareUnitMember{}
						shareMember.ShareMemberUin = helper.IntInt64(v.(int))
						tmpList = append(tmpList, &shareMember)
					}
				}
			}

			err := service.AddOrganizationOrgShareUnitMemberV2ById(ctx, unitId, area, tmpList)
			if err != nil {
				return err
			}
		}

		if len(remove) > 0 {
			tmpList := make([]*organization.ShareUnitMember, 0, len(remove))
			for _, item := range remove {
				if dMap, ok := item.(map[string]interface{}); ok {
					if v, ok := dMap["share_member_uin"]; ok {
						shareMember := organization.ShareUnitMember{}
						shareMember.ShareMemberUin = helper.IntInt64(v.(int))
						tmpList = append(tmpList, &shareMember)
					}
				}
			}

			err := service.DeleteOrganizationOrgShareUnitMemberV2ById(ctx, unitId, area, tmpList)
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func resourceTencentCloudOrganizationOrgShareUnitMemberV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_share_unit_member_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	unitId := idSplit[0]
	area := idSplit[1]

	// get all members
	orgShareUnitMembers, err := service.DescribeOrganizationOrgShareUnitMemberV2ById(ctx, unitId, area)
	if err != nil {
		return err
	}

	if len(orgShareUnitMembers) < 1 {
		log.Printf("[WARN]%s resource `tencentcloud_organization_org_share_unit_member_v2` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	// delete all members
	if err := service.DeleteOrganizationOrgShareUnitMemberV2ById(ctx, unitId, area, orgShareUnitMembers); err != nil {
		return err
	}

	return nil
}
