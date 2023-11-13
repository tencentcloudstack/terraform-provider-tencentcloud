/*
Provides a resource to create a teo ownership

Example Usage

```hcl
resource "tencentcloud_teo_ownership" "ownership" {
  domain = "qq.com"
}
```

Import

teo ownership can be imported using the id, e.g.

```
terraform import tencentcloud_teo_ownership.ownership ownership_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudTeoOwnership() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoOwnershipCreate,
		Read:   resourceTencentCloudTeoOwnershipRead,
		Delete: resourceTencentCloudTeoOwnershipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Site or acceleration domain name.",
			},
		},
	}
}

func resourceTencentCloudTeoOwnershipCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewVerifyOwnershipRequest()
		response = teo.NewVerifyOwnershipResponse()
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().VerifyOwnership(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate teo ownership failed, reason:%+v", logId, err)
		return err
	}

	domain = *response.Response.Domain
	d.SetId(domain)

	return resourceTencentCloudTeoOwnershipRead(d, meta)
}

func resourceTencentCloudTeoOwnershipRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoOwnershipDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
