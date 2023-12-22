package teo

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoOwnershipVerify() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ownership_verify.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = teo.NewVerifyOwnershipRequest()
		response = teo.NewVerifyOwnershipResponse()
		domain   string
	)
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		request.Domain = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().VerifyOwnership(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ownership_verify.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudTeoOwnershipVerifyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_ownership_verify.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
