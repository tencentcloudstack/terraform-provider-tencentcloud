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

func ResourceTencentCloudOrganizationOrgManagePolicyTarget() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgManagePolicyTargetCreate,
		Read:   resourceTencentCloudOrganizationOrgManagePolicyTargetRead,
		Delete: resourceTencentCloudOrganizationOrgManagePolicyTargetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"target_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Binding target ID of the policy. Member Uin or Department ID.",
			},

			"target_type": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Target type.\nValid values:\n  - `NODE`: Department.\n  - `MEMBER`: Check Member.",
			},

			"policy_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Policy Id.",
			},

			"policy_type": {
				Optional:    true,
				ForceNew:    true,
				Default:     ServiceControlPolicyType,
				Type:        schema.TypeString,
				Description: "Policy type. Default value is SERVICE_CONTROL_POLICY.\nValid values:\n  - `SERVICE_CONTROL_POLICY`: Service control policy.\n  - `TAG_POLICY`: Tag policy.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgManagePolicyTargetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_target.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = organization.NewAttachPolicyRequest()
		policyType string
		policyId   int
		targetType string
		targetId   int
	)
	if v, ok := d.GetOkExists("target_id"); ok {
		targetId = v.(int)
		request.TargetId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("target_type"); ok {
		targetType = v.(string)
		request.TargetType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("policy_id"); ok {
		policyId = v.(int)
		request.PolicyId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy_type"); ok {
		policyType = v.(string)
		request.Type = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().AttachPolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization OrgManagePolicyTarget failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{policyType, helper.Int64ToStr(int64(policyId)), targetType, helper.Int64ToStr(int64(targetId))}, tccommon.FILED_SP))

	return resourceTencentCloudOrganizationOrgManagePolicyTargetRead(d, meta)
}

func resourceTencentCloudOrganizationOrgManagePolicyTargetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_target.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyType := idSplit[0]

	policyId := idSplit[1]

	targetType := idSplit[2]

	targetId := idSplit[3]

	OrgManagePolicyTarget, err := service.DescribeOrganizationOrgManagePolicyTargetById(ctx, policyType, policyId, targetType, targetId)
	if err != nil {
		return err
	}

	if OrgManagePolicyTarget == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgManagePolicyTarget` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("target_id", OrgManagePolicyTarget.Uin)
	if OrgManagePolicyTarget.RelatedType != nil {
		switch *OrgManagePolicyTarget.RelatedType {
		case 1:
			_ = d.Set("target_type", TargetTypeNode)
		case 2:
			_ = d.Set("target_type", TargetTypeMember)
		}

	}
	_ = d.Set("policy_id", helper.StrToInt(policyId))
	_ = d.Set("policy_type", policyType)
	return nil
}

func resourceTencentCloudOrganizationOrgManagePolicyTargetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_target.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 4 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	policyType := idSplit[0]

	policyId := idSplit[1]

	targetType := idSplit[2]

	targetId := idSplit[3]

	if err := service.DeleteOrganizationOrgManagePolicyTargetById(ctx, policyType, policyId, targetType, targetId); err != nil {
		return err
	}

	return nil
}
