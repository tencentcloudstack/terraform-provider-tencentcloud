package tco

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudOrganizationServiceAssign() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationServiceAssignCreate,
		Read:   resourceTencentCloudOrganizationServiceAssignRead,
		Update: resourceTencentCloudOrganizationServiceAssignUpdate,
		Delete: resourceTencentCloudOrganizationServiceAssignDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"service_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Organization service ID.",
			},
			"member_uins": {
				Required:    true,
				Type:        schema.TypeList,
				MaxItems:    20,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Uin list of the delegated admins, Including up to 20 items.",
			},
			"management_scope": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Management scope of the delegated admin. Valid values: 1 (all members), 2 (partial members). Default value: `1`.",
			},
			"management_scope_uins": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Uin list of the managed members. This parameter is valid when `management_scope` is `2`.",
			},
			"management_scope_node_ids": {
				Optional:    true,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "ID list of the managed departments. This parameter is valid when `management_scope` is `2`.",
			},
		},
	}
}

func resourceTencentCloudOrganizationServiceAssignCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_service_assign.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		request   = organization.NewCreateOrgServiceAssignRequest()
		serviceId string
	)

	if v, _ := d.GetOkExists("service_id"); v != nil {
		request.ServiceId = helper.IntUint64(v.(int))
		serviceId = strconv.Itoa(v.(int))
	}

	if v, _ := d.GetOk("member_uins"); v != nil {
		memberUins := v.([]interface{})
		tmpList := make([]*int64, 0, len(memberUins))
		for i := range memberUins {
			memberUin := memberUins[i].(int)
			tmpList = append(tmpList, helper.IntInt64(memberUin))
		}

		request.MemberUins = tmpList
	}

	if v, _ := d.GetOkExists("management_scope"); v != nil {
		request.ManagementScope = helper.IntUint64(v.(int))
	}

	if v, _ := d.GetOk("management_scope_uins"); v != nil {
		managementScopeUins := v.([]interface{})
		tmpList := make([]*int64, 0, len(managementScopeUins))
		for i := range managementScopeUins {
			managementScopeUin := managementScopeUins[i].(int)
			tmpList = append(tmpList, helper.IntInt64(managementScopeUin))
		}

		request.ManagementScopeUins = tmpList
	}

	if v, _ := d.GetOk("management_scope_node_ids"); v != nil {
		managementScopeNodeIds := v.([]interface{})
		tmpList := make([]*int64, 0, len(managementScopeNodeIds))
		for i := range managementScopeNodeIds {
			managementScopeNodeId := managementScopeNodeIds[i].(int)
			tmpList = append(tmpList, helper.IntInt64(managementScopeNodeId))
		}

		request.ManagementScopeNodeIds = tmpList
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreateOrgServiceAssign(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create organization service assign failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(serviceId)

	return resourceTencentCloudOrganizationServiceAssignRead(d, meta)
}

func resourceTencentCloudOrganizationServiceAssignRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_service_assign.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId = d.Id()
	)

	items, err := service.DescribeOrganizationServiceAssignMemberById(ctx, serviceId)
	if err != nil {
		return err
	}

	if len(items) == 0 {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationServiceAssignMember` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if items[0].ServiceId != nil {
		_ = d.Set("service_id", items[0].ServiceId)
	}

	if items[0].ManagementScope != nil {
		_ = d.Set("management_scope", items[0].ManagementScope)
	}

	tmpMemberUinList := make([]int64, 0)
	for _, item := range items {
		if item.MemberUin != nil {
			tmpMemberUinList = append(tmpMemberUinList, *item.MemberUin)
		}
	}

	if len(tmpMemberUinList) != 0 {
		_ = d.Set("member_uins", tmpMemberUinList)
	}

	if items[0].ManagementScopeMembers != nil {
		tmpList := make([]*int64, 0, len(items[0].ManagementScopeMembers))
		for _, v := range items[0].ManagementScopeMembers {
			if v.MemberUin != nil {
				tmpList = append(tmpList, v.MemberUin)
			}
		}

		_ = d.Set("management_scope_uins", tmpList)
	}

	if items[0].ManagementScopeNodes != nil {
		tmpList := make([]*int64, 0, len(items[0].ManagementScopeNodes))
		for _, v := range items[0].ManagementScopeNodes {
			if v.NodeId != nil {
				tmpList = append(tmpList, v.NodeId)
			}
		}

		_ = d.Set("management_scope_node_ids", tmpList)
	}

	return nil
}

func resourceTencentCloudOrganizationServiceAssignUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_service_assign.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	immutableArgs := []string{"service_id", "member_uins", "management_scope", "management_scope_uins", "management_scope_node_ids"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	return resourceTencentCloudOrganizationServiceAssignRead(d, meta)
}

func resourceTencentCloudOrganizationServiceAssignDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_service_assign.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service   = OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		serviceId = d.Id()
	)

	memberUinList := make([]*int64, 0)
	if v, _ := d.GetOk("member_uins"); v != nil {
		memberUins := v.([]interface{})
		for i := range memberUins {
			memberUin := memberUins[i].(int)
			memberUinList = append(memberUinList, helper.IntInt64(memberUin))
		}
	}

	if err := service.DeleteOrganizationServiceAssignMemberById(ctx, serviceId, memberUinList); err != nil {
		return err
	}

	return nil
}
