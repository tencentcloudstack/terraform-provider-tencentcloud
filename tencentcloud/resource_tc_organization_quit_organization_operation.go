package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudOrganizationQuitOrganizationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudOrganizationQuitOrganizationOperationCreate,
		Read:   resourceTencentCloudOrganizationQuitOrganizationOperationRead,
		Delete: resourceTencentCloudOrganizationQuitOrganizationOperationDelete,
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

func resourceTencentCloudOrganizationQuitOrganizationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request = organization.NewQuitOrganizationRequest()
		orgId   uint64
	)
	if v, _ := d.GetOk("org_id"); v != nil {
		request.OrgId = helper.IntUint64(v.(int))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseOrganizationClient().QuitOrganization(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate organization quitOrganizationOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.UInt64ToStr(orgId))

	return resourceTencentCloudOrganizationQuitOrganizationOperationRead(d, meta)
}

func resourceTencentCloudOrganizationQuitOrganizationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudOrganizationQuitOrganizationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_organization_quit_organization_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
