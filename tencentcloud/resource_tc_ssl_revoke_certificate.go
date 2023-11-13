/*
Provides a resource to create a ssl revoke_certificate

Example Usage

```hcl
resource "tencentcloud_ssl_revoke_certificate" "revoke_certificate" {
  certificate_id = ""
  reason = ""
}
```

Import

ssl revoke_certificate can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_revoke_certificate.revoke_certificate revoke_certificate_id
```
*/
package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSslRevokeCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslRevokeCertificateCreate,
		Read:   resourceTencentCloudSslRevokeCertificateRead,
		Delete: resourceTencentCloudSslRevokeCertificateDelete,
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

			"reason": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Reasons for revoking certificate.",
			},
		},
	}
}

func resourceTencentCloudSslRevokeCertificateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ssl.NewRevokeCertificateRequest()
		response      = ssl.NewRevokeCertificateResponse()
		certificateId string
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reason"); ok {
		request.Reason = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSslClient().RevokeCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl revokeCertificate failed, reason:%+v", logId, err)
		return err
	}

	certificateId = *response.Response.CertificateId
	d.SetId(certificateId)

	return resourceTencentCloudSslRevokeCertificateRead(d, meta)
}

func resourceTencentCloudSslRevokeCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslRevokeCertificateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
