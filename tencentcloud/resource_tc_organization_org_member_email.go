/*
Provides a resource to create a organization org_member_email

Example Usage

```hcl
resource "tencentcloud_organization_org_member_email" "org_member_email" {
  member_uin = &lt;nil&gt;
  email = &lt;nil&gt;
  country_code = &lt;nil&gt;
  phone = &lt;nil&gt;
  bind_id = &lt;nil&gt;
}
```

Import

organization org_member_email can be imported using the id, e.g.

```
terraform import tencentcloud_organization_org_member_email.org_member_email org_member_email_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudOrganizationOrgMemberEmail() *schema.Resource {
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
				Description: "Member uin.",
			},

			"email": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Email address.",
			},

			"country_code": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "International dialing code.",
			},

			"phone": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Mobile number.",
			},

			"bind_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Bind Id.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrgMemberEmailCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_email.create")()
	defer inconsistentCheck(d, meta)()

	var bindId uint64
	if v, ok := d.GetOkExists("bind_id"); ok {
		bindId = v.(uint64)
	}

	d.SetId(helper.Int64ToStr(int64(bindId)))

	return resourceTencentCloudOrganizationOrgMemberEmailUpdate(d, meta)
}

func resourceTencentCloudOrganizationOrgMemberEmailRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_email.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	orgMemberEmailId := d.Id()

	orgMemberEmail, err := service.DescribeOrganizationOrgMemberEmailById(ctx, bindId)
	if err != nil {
		return err
	}

	if orgMemberEmail == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `OrganizationOrgMemberEmail` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if orgMemberEmail.MemberUin != nil {
		_ = d.Set("member_uin", orgMemberEmail.MemberUin)
	}

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

	return nil
}

func resourceTencentCloudOrganizationOrgMemberEmailUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_org_member_email.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := organization.NewUpdateOrganizationMemberEmailBindRequest()

	orgMemberEmailId := d.Id()

	request.BindId = &bindId

	immutableArgs := []string{"member_uin", "email", "country_code", "phone", "bind_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("member_uin") {
		if v, ok := d.GetOkExists("member_uin"); ok {
			request.MemberUin = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("email") {
		if v, ok := d.GetOk("email"); ok {
			request.Email = helper.String(v.(string))
		}
	}

	if d.HasChange("country_code") {
		if v, ok := d.GetOk("country_code"); ok {
			request.CountryCode = helper.String(v.(string))
		}
	}

	if d.HasChange("phone") {
		if v, ok := d.GetOk("phone"); ok {
			request.Phone = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().UpdateOrganizationMemberEmailBind(request)
		if e != nil {
			return retryError(e)
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
	defer logElapsed("resource.tencentcloud_organization_org_member_email.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
