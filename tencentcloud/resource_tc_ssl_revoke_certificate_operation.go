/*
Provides a resource to create a ssl revoke_certificate_operation

Example Usage

```hcl
resource "tencentcloud_ssl_revoke_certificate_operation" "revoke_certificate_operation" {
  certificate_id = "a91hoLqi"
  reason = "xx"
}
```

Import

ssl revoke_certificate_operation can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_revoke_certificate_operation.revoke_certificate_operation revoke_certificate_operation_id
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

func resourceTencentCloudSslRevokeCertificateOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslRevokeCertificateOperationCreate,
		Read:   resourceTencentCloudSslRevokeCertificateOperationRead,
		Delete: resourceTencentCloudSslRevokeCertificateOperationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate Id.",
			},

			"reason": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The reason for revoking the certificate.",
			},
		},
	}
}

func resourceTencentCloudSslRevokeCertificateOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate_operation.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request       = ssl.NewRevokeCertificateRequest()
		certificateId uint64
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("reason"); ok {
		request.Reason = helper.String(v.(string))
	}
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().RevokeCertificate(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl revokeCertificateOperation failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(helper.UInt64ToStr(certificateId))

	return resourceTencentCloudSslRevokeCertificateOperationRead(d, meta)
}

func resourceTencentCloudSslRevokeCertificateOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate_operation.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslRevokeCertificateOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_revoke_certificate_operation.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
