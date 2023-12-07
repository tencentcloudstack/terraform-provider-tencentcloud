package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoOwnershipVerify() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoOwnershipVerifyCreate,
		Read:   resourceTencentCloudTeoOwnershipVerifyRead,
		Delete: resourceTencentCloudTeoOwnershipVerifyDelete,

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
				Description: "Ownership verification results. `success`: verification successful; `fail`: verification failed.",
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "When the verification result is failed, this field will return the reason.",
			},
		},
	}
}

func resourceTencentCloudTeoOwnershipVerifyCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_ownership_verify.create")()
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
		log.Printf("[CRITAL]%s operate teo ownershipVerify failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(domain)

	if response != nil && response.Response != nil {
		_ = d.Set("status", *response.Response.Status)
		_ = d.Set("result", *response.Response.Result)
	}

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
