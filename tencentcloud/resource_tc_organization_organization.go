/*
Provides a resource to create a organization organization

Example Usage

```hcl
resource "tencentcloud_organization_organization" "organization" {
  org_id = &lt;nil&gt;
}
```

Import

organization organization can be imported using the id, e.g.

```
terraform import tencentcloud_organization_organization.organization organization_id
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

func resourceTencentCloudOrganizationOrganization() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationOrganizationCreate,
		Read:   resourceTencentCloudOrganizationOrganizationRead,
		Update: resourceTencentCloudOrganizationOrganizationUpdate,
		Delete: resourceTencentCloudOrganizationOrganizationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"org_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Organization Id.",
			},
		},
	}
}

func resourceTencentCloudOrganizationOrganizationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_organization.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = organization.NewCreateOrganizationRequest()
		response = organization.NewCreateOrganizationResponse()
		orgId    int
	)
	if v, ok := d.GetOkExists("org_id"); ok {
		orgId = v.(int64)
		request.OrgId = helper.IntInt64(v.(int))
	}

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
	d.SetId(helper.Int64ToStr(orgId))

	return resourceTencentCloudOrganizationOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationOrganizationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_organization.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}

	organizationId := d.Id()

	organization, err := service.DescribeOrganizationOrganizationById(ctx, orgId)
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

	return nil
}

func resourceTencentCloudOrganizationOrganizationUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_organization.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"org_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudOrganizationOrganizationRead(d, meta)
}

func resourceTencentCloudOrganizationOrganizationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_organization.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := OrganizationService{client: meta.(*TencentCloudClient).apiV3Conn}
	organizationId := d.Id()

	if err := service.DeleteOrganizationOrganizationById(ctx, orgId); err != nil {
		return err
	}

	return nil
}
