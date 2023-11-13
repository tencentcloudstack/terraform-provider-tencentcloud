/*
Provides a resource to create a ssl check_certificate_chain

Example Usage

```hcl
resource "tencentcloud_ssl_check_certificate_chain" "check_certificate_chain" {
  certificate_chain = ""
}
```

Import

ssl check_certificate_chain can be imported using the id, e.g.

```
terraform import tencentcloud_ssl_check_certificate_chain.check_certificate_chain check_certificate_chain_id
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

func resourceTencentCloudSslCheckCertificateChain() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCheckCertificateChainCreate,
		Read:   resourceTencentCloudSslCheckCertificateChainRead,
		Delete: resourceTencentCloudSslCheckCertificateChainDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"certificate_chain": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "The certificate chain to check.",
			},
		},
	}
}

func resourceTencentCloudSslCheckCertificateChainCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_check_certificate_chain.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = ssl.NewCheckCertificateChainRequest()
		response         = ssl.NewCheckCertificateChainResponse()
		certificateChain string
	)
	if v, ok := d.GetOk("certificate_chain"); ok {
		certificateChain = v.(string)
		request.CertificateChain = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSslClient().CheckCertificateChain(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate ssl checkCertificateChain failed, reason:%+v", logId, err)
		return err
	}

	certificateChain = *response.Response.CertificateChain
	d.SetId(certificateChain)

	return resourceTencentCloudSslCheckCertificateChainRead(d, meta)
}

func resourceTencentCloudSslCheckCertificateChainRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_check_certificate_chain.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudSslCheckCertificateChainDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_check_certificate_chain.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
