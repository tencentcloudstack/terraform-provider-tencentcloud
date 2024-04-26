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

func ResourceTencentCloudOrganizationOrgManagePolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgManagePolicyConfigCreate,
		Read:   resourceTencentCloudOrganizationOrgManagePolicyConfigRead,
		Delete: resourceTencentCloudOrganizationOrgManagePolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"organization_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization ID.",
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

func resourceTencentCloudOrganizationOrgManagePolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request        = organization.NewEnablePolicyTypeRequest()
		organizationId int
		policyType     string
	)
	if v, ok := d.GetOkExists("organization_id"); ok {
		organizationId = v.(int)
		request.OrganizationId = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("policy_type"); ok {
		policyType = v.(string)
		request.PolicyType = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().EnablePolicyType(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization OrgManagePolicyConfig failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{helper.IntToStr(organizationId), policyType}, tccommon.FILED_SP))

	return resourceTencentCloudOrganizationOrgManagePolicyConfigRead(d, meta)
}

func resourceTencentCloudOrganizationOrgManagePolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	organizationId := idSplit[0]

	policyType := idSplit[1]

	OrgManagePolicyConfig, err := service.DescribeOrganizationOrgManagePolicyConfigById(ctx, organizationId, policyType)
	if err != nil {
		return err
	}

	if OrgManagePolicyConfig == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgManagePolicyConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("organization_id", helper.StrToInt(organizationId))
	_ = d.Set("policy_type", policyType)
	return nil
}

func resourceTencentCloudOrganizationOrgManagePolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)

	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	organizationId := idSplit[0]

	policyType := idSplit[1]

	if err := service.DeleteOrganizationOrgManagePolicyConfigById(ctx, organizationId, policyType); err != nil {
		return err
	}

	return nil
}
