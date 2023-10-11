/*
Provides a resource to create a ssl upload_revoke_letter

Example Usage

```hcl
resource "tencentcloud_ssl_upload_revoke_letter_operation" "upload_revoke_letter" {
  certificate_id = "8xRYdDlc"
  revoke_letter = filebase64("./c.pdf")
}
```

Import

ssl upload_revoke_letter can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_upload_revoke_letter_operation.upload_revoke_letter upload_revoke_letter_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslUploadRevokeLetterOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslUploadRevokeLetterCreate,
		Read:   resourceTencentCloudSslUploadRevokeLetterRead,
		Delete: resourceTencentCloudSslUploadRevokeLetterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate ID.",
			},

			"revoke_letter": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The format of the base64-encoded certificate confirmation letter file should be jpg, jpeg, png, or pdf, and the size should be between 1kb and 1.4M.",
			},
		},
	}
}

func resourceTencentCloudSslUploadRevokeLetterCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_upload_revoke_letter_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ssl.NewUploadRevokeLetterRequest()
		response      = ssl.NewUploadRevokeLetterResponse()
		certificateId string
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("revoke_letter"); ok {
		request.RevokeLetter = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().UploadRevokeLetter(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl uploadRevokeLetter failed, reason:%+v", logId, err)
		return err
	}
	if response != nil && response.Response != nil && response.Response.CertificateId != nil {
		certificateId = *response.Response.CertificateId
	}
	d.SetId(certificateId)

	return resourceTencentCloudSslUploadRevokeLetterRead(d, meta)
}

func resourceTencentCloudSslUploadRevokeLetterRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_upload_revoke_letter_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslUploadRevokeLetterDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_upload_revoke_letter_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
