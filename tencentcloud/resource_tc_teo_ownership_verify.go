/*
Provides a resource to create a teo ownership_verify

Example Usage

```hcl
resource "tencentcloud_teo_ownership_verify" "ownership_verify" {
  domain = "qq.com"
}
```
*/
package tencentcloud

import (
	// "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	// teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	// "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	// "log"
)

func resourceTencentCloudTeoOwnershipVerify() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoOwnershipVerifyCreate,
		Read:   resourceTencentCloudTeoOwnershipVerifyRead,
		Delete: resourceTencentCloudTeoOwnershipVerifyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Verify domain name.",
			},

			"status": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Ownership verification results. `success`: verification successful; `fail`: verification failed..",
			},
		},
	}
}

func resourceTencentCloudTeoOwnershipVerifyCreate(d *schema.ResourceData, meta interface{}) error {
	// defer logElapsed("resource.tencentcloud_teo_ownership_verify.create")()
	// defer inconsistentCheck(d, meta)()

	// logId := getLogId(contextNil)

	// var (
	// 	request  = teo.NewVerifyOwnershipRequest()
	// 	response = teo.NewVerifyOwnershipResponse()
	// 	domain   string
	// )
	// if v, ok := d.GetOk("domain"); ok {
	// 	request.Domain = helper.String(v.(string))
	// }

	// err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
	// 	result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().VerifyOwnership(request)
	// 	if e != nil {
	// 		return retryError(e)
	// 	} else {
	// 		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
	// 	}
	// 	response = result
	// 	return nil
	// })
	// if err != nil {
	// 	log.Printf("[CRITAL]%s operate teo ownershipVerify failed, reason:%+v", logId, err)
	// 	return err
	// }

	// d.SetId(domain)

	// _ = d.Set("type", *response.Response.Status)

	return resourceTencentCloudTeoOwnershipVerifyRead(d, meta)
}

func resourceTencentCloudTeoOwnershipVerifyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership_verify.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoOwnershipVerifyDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership_verify.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
