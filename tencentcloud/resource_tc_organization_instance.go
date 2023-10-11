/*
Provides a resource to create a organization organization

Example Usage

```hcl
resource "tencentcloud_organization_instance" "organization" {
  }
```

Import

organization organization can be imported using the id, e.g.

```
terraform import tencentcloud_organization_instance.organization organization_id
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrganizationCreate,
		Read:   resourceTencentCloudOrganizationOrganizationRead,
		Delete: resourceTencentCloudOrganizationOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Enterprise organization ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"host_uin": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Creator Uin.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"nick_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Creator nickname.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"org_type": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Enterprise organization type.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether to organize an administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"org_policy_type": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Strategy type.Financial Management: FinancialNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"org_policy_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Strategic name.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"org_permission": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "List of membership authority of members.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Permissions ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Permission name.",
						},
					},
				},
			},

			"root_node_id": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Organize the root node ID.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Organize the creation time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"join_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Members join time.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_allow_quit": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Whether the members are allowed to withdraw.Allow: Allow, not allowed: DENIEDNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"pay_uin": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "UIN on behalf of the payer.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"pay_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The name of the payment.Note: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_assign_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether a trusted service administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},

			"is_auth_manager": {
				Computed:    true,
				Type:        schema.TypeBool,
				Description: "Whether the real -name subject administrator.Yes: true, no: falseNote: This field may return NULL, indicating that the valid value cannot be obtained.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewCreateOrganizationRequest()
		response = organization.NewCreateOrganizationResponse()
		orgId    uint64
	)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().CreateOrganization(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create organization organization failed, reason:%+v", logId, err)
		return err
	}

	orgId = *response.Response.OrgId
	d.SetId(helper.UInt64ToStr(orgId))

	return resourceTencentCloudOrganizationOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	organization, err := service.DescribeOrganizationOrganizationById(ctx)
	if err != nil {
		return err
	}

	if organization == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrganization` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if organization.OrgId != nil {
		_ = d.Set("org_id", organization.OrgId)
	}

	if organization.HostUin != nil {
		_ = d.Set("host_uin", organization.HostUin)
	}

	if organization.NickName != nil {
		_ = d.Set("nick_name", organization.NickName)
	}

	if organization.OrgType != nil {
		_ = d.Set("org_type", organization.OrgType)
	}

	if organization.IsManager != nil {
		_ = d.Set("is_manager", organization.IsManager)
	}

	if organization.OrgPolicyType != nil {
		_ = d.Set("org_policy_type", organization.OrgPolicyType)
	}

	if organization.OrgPolicyName != nil {
		_ = d.Set("org_policy_name", organization.OrgPolicyName)
	}

	if organization.OrgPermission != nil {
		var orgPermissionList []interface{}
		for _, orgPermission := range organization.OrgPermission {
			orgPermissionMap := map[string]interface{}{}

			if orgPermission.Id != nil {
				orgPermissionMap["id"] = orgPermission.Id
			}

			if orgPermission.Name != nil {
				orgPermissionMap["name"] = orgPermission.Name
			}

			orgPermissionList = append(orgPermissionList, orgPermissionMap)
		}

		_ = d.Set("org_permission", orgPermissionList)

	}

	if organization.RootNodeId != nil {
		_ = d.Set("root_node_id", organization.RootNodeId)
	}

	if organization.CreateTime != nil {
		_ = d.Set("create_time", organization.CreateTime)
	}

	if organization.JoinTime != nil {
		_ = d.Set("join_time", organization.JoinTime)
	}

	if organization.IsAllowQuit != nil {
		_ = d.Set("is_allow_quit", organization.IsAllowQuit)
	}

	if organization.PayUin != nil {
		_ = d.Set("pay_uin", organization.PayUin)
	}

	if organization.PayName != nil {
		_ = d.Set("pay_name", organization.PayName)
	}

	if organization.IsAssignManager != nil {
		_ = d.Set("is_assign_manager", organization.IsAssignManager)
	}

	if organization.IsAuthManager != nil {
		_ = d.Set("is_auth_manager", organization.IsAuthManager)
	}

	return nil
}

func resourceTencentCloudOrganizationOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	if err := service.DeleteOrganizationOrganizationById(ctx); err != nil {
		return err
	}

	return nil
}
