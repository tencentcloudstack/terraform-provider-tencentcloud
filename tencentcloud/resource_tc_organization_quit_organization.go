/*
Provides a resource to create a organization quit_organization

Example Usage

```hcl
resource "tencentcloud_organization_quit_organization" "quit_organization" {
  org_id = &lt;nil&gt;
}
```

Import

organization quit_organization can be imported using the id, e.g.

```
terraform import tencentcloud_organization_quit_organization.quit_organization quit_organization_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudOrganizationQuitOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationQuitOrganizationCreate,
		Read:   resourceTencentCloudOrganizationQuitOrganizationRead,
		Delete: resourceTencentCloudOrganizationQuitOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeInt,
				Description: "Organization ID.",
			},
		},
	}
}

func resourceTencentCloudOrganizationQuitOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewQuitOrganizationRequest()
		response = organization.NewQuitOrganizationResponse()
		orgId    int
	)
	if v, _ := d.GetOk("org_id"); v != nil {
		orgId = v.(int)
		request.OrgId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().QuitOrganization(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate organization quitOrganization failed, reason:%+v", logId, err)
		return err
	}

	orgId = *response.Response.OrgId
	d.SetId(helper.Int64ToStr(orgId))

	return resourceTencentCloudOrganizationQuitOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationQuitOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOrganizationQuitOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
