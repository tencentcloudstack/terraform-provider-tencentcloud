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

func ResourceTencentCloudOrganizationOrgManagePolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgManagePolicyCreate,
		Read:   resourceTencentCloudOrganizationOrgManagePolicyRead,
		Update: resourceTencentCloudOrganizationOrgManagePolicyUpdate,
		Delete: resourceTencentCloudOrganizationOrgManagePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Policy name.\nThe length is 1~128 characters, which can include Chinese characters, English letters, numbers, and underscores.",
			},

			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Policy content. Refer to the CAM policy syntax.",
			},

			"type": {
				Optional:    true,
				Default:     ServiceControlPolicyType,
				Type:        schema.TypeString,
				Description: "Policy type. Default value is SERVICE_CONTROL_POLICY.\nValid values:\n  - `SERVICE_CONTROL_POLICY`: Service control policy.\n  - `TAG_POLICY`: Tag policy.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Policy description.",
			},

			"policy_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Policy Id.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgManagePolicyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		policyType string
		request    = organization.NewCreatePolicyRequest()
		response   = organization.NewCreatePolicyResponse()
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		policyType = v.(string)
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().CreatePolicy(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization OrgManagePolicy failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{helper.UInt64ToStr(*response.Response.PolicyId), policyType}, tccommon.FILED_SP))
	return resourceTencentCloudOrganizationOrgManagePolicyRead(d, meta)
}

func resourceTencentCloudOrganizationOrgManagePolicyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	policyType := idSplit[1]

	OrgManagePolicy, err := service.DescribeOrganizationOrgManagePolicyById(ctx, policyId, policyType)
	if err != nil {
		return err
	}

	if OrgManagePolicy == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgManagePolicy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if OrgManagePolicy.PolicyName != nil {
		_ = d.Set("name", OrgManagePolicy.PolicyName)
	}

	if OrgManagePolicy.PolicyDocument != nil {
		_ = d.Set("content", OrgManagePolicy.PolicyDocument)
	}

	if OrgManagePolicy.Type != nil {
		_ = d.Set("type", policyType)
	}

	if OrgManagePolicy.Description != nil {
		_ = d.Set("description", OrgManagePolicy.Description)
	}
	_ = d.Set("policy_id", policyId)

	return nil
}

func resourceTencentCloudOrganizationOrgManagePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := organization.NewUpdatePolicyRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]

	request.PolicyId = helper.StrToInt64Point(policyId)

	needChange := false
	mutableArgs := []string{"name", "content", "type", "description"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
		if v, ok := d.GetOk("content"); ok {
			request.Content = helper.String(v.(string))
		}
		if v, ok := d.GetOk("type"); ok {
			request.Type = helper.String(v.(string))
		}
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseOrganizationClient().UpdatePolicy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update organization OrgManagePolicy failed, reason:%+v", logId, err)
			return err
		}

	}
	return resourceTencentCloudOrganizationOrgManagePolicyRead(d, meta)
}

func resourceTencentCloudOrganizationOrgManagePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_organization_org_manage_policy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := OrganizationService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	policyId := idSplit[0]
	policyType := idSplit[1]

	if err := service.DeleteOrganizationOrgManagePolicyById(ctx, policyId, policyType); err != nil {
		return err
	}

	return nil
}
