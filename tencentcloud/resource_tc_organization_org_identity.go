/*
Provides a resource to create a organization org_identity

Example Usage

```hcl
resource "tencentcloud_organization_org_identity" "org_identity" {
  identity_alias_name = "example-iac-test"
  identity_policy {
    policy_id = 1
    policy_name = "AdministratorAccess"
    policy_type = 2
  }
  description = "iac-test"
}
```

Import

organization org_identity can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_identity.org_identity org_identity_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationOrgIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrgIdentityCreate,
		Read:   resourceTencentCloudOrganizationOrgIdentityRead,
		Update: resourceTencentCloudOrganizationOrgIdentityUpdate,
		Delete: resourceTencentCloudOrganizationOrgIdentityDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"identity_alias_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Identity name.Supports English letters and numbers, the length cannot exceed 40 characters.",
			},

			"identity_policy": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Identity policy list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "CAM default policy ID. Valid and required when PolicyType is the 2-preset policy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "CAM default policy name. Valid and required when PolicyType is the 2-preset policy.",
						},
						"policy_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Policy type. Value 1-custom policy 2-preset policy; default value 2.",
						},
						"policy_document": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Customize policy content and follow CAM policy syntax. Valid and required when PolicyType is the 1-custom policy.",
						},
					},
				},
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Identity description.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_identity.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = organization.NewCreateOrganizationIdentityRequest()
		response   = organization.NewCreateOrganizationIdentityResponse()
		identityId string
	)
	if v, ok := d.GetOk("identity_alias_name"); ok {
		request.IdentityAliasName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("identity_policy"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			identityPolicy := organization.IdentityPolicy{}
			if v, ok := dMap["policy_id"]; ok {
				identityPolicy.PolicyId = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["policy_name"]; ok {
				identityPolicy.PolicyName = helper.String(v.(string))
			}
			if v, ok := dMap["policy_type"]; ok {
				identityPolicy.PolicyType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["policy_document"]; ok {
				identityPolicy.PolicyDocument = helper.String(v.(string))
			}
			request.IdentityPolicy = append(request.IdentityPolicy, &identityPolicy)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganizationIdentity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization orgIdentity failed, reason:%+v", logId, err)
		return err
	}

	identityId = helper.UInt64ToStr(*response.Response.IdentityId)
	d.SetId(identityId)

	return resourceTencentCloudOrganizationOrgIdentityRead(d, meta)
}

func resourceTencentCloudOrganizationOrgIdentityRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_identity.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgIdentityId := d.Id()

	orgIdentity, err := service.DescribeOrganizationOrgIdentityById(ctx, orgIdentityId)
	if err != nil {
		return err
	}

	if orgIdentity == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgIdentity` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgIdentity.IdentityAliasName != nil {
		_ = d.Set("identity_alias_name", orgIdentity.IdentityAliasName)
	}

	if orgIdentity.IdentityPolicy != nil {
		var identityPolicyList []interface{}
		for _, identityPolicy := range orgIdentity.IdentityPolicy {
			identityPolicyMap := map[string]interface{}{}

			if identityPolicy.PolicyId != nil {
				identityPolicyMap["policy_id"] = identityPolicy.PolicyId
			}

			if identityPolicy.PolicyName != nil {
				identityPolicyMap["policy_name"] = identityPolicy.PolicyName
			}

			if identityPolicy.PolicyType != nil {
				identityPolicyMap["policy_type"] = identityPolicy.PolicyType
			}

			if identityPolicy.PolicyDocument != nil {
				identityPolicyMap["policy_document"] = identityPolicy.PolicyDocument
			}

			identityPolicyList = append(identityPolicyList, identityPolicyMap)
		}

		_ = d.Set("identity_policy", identityPolicyList)

	}

	if orgIdentity.Description != nil {
		_ = d.Set("description", orgIdentity.Description)
	}

	return nil
}

func resourceTencentCloudOrganizationOrgIdentityUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_identity.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := organization.NewUpdateOrganizationIdentityRequest()

	orgIdentityId := d.Id()

	request.IdentityId = helper.StrToUint64Point(orgIdentityId)

	immutableArgs := []string{"identity_alias_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if v, ok := d.GetOk("identity_policy"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			identityPolicy := organization.IdentityPolicy{}
			if v, ok := dMap["policy_id"]; ok {
				identityPolicy.PolicyId = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["policy_name"]; ok {
				identityPolicy.PolicyName = helper.String(v.(string))
			}
			if v, ok := dMap["policy_type"]; ok {
				identityPolicy.PolicyType = helper.IntUint64(v.(int))
			}
			if v, ok := dMap["policy_document"]; ok {
				identityPolicy.PolicyDocument = helper.String(v.(string))
			}
			request.IdentityPolicy = append(request.IdentityPolicy, &identityPolicy)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().UpdateOrganizationIdentity(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update organization orgIdentity failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudOrganizationOrgIdentityRead(d, meta)
}

func resourceTencentCloudOrganizationOrgIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_identity.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	orgIdentityId := d.Id()

	if err := service.DeleteOrganizationOrgIdentityById(ctx, orgIdentityId); err != nil {
		return err
	}

	return nil
}
