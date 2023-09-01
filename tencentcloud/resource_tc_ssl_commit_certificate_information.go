/*
Provides a resource to create a ssl commit_certificate_information

Example Usage

```hcl
resource "tencentcloud_ssl_pay_certificate" "example" {
  product_id = 33
  domain_num = 1
  alias      = "example-ssl-update"
  project_id = 0
  wait_commit_flag = true
  information {
    csr_type              = "online"
    certificate_domain    = "www.domain.com"
    organization_name     = "test-update"
    organization_division = "test"
    organization_address  = "test"
    organization_country  = "CN"
    organization_city     = "test"
    organization_region   = "test"
    postal_code           = "0755"
    phone_area_code       = "0755"
    phone_number          = "12345678901"
    verify_type           = "DNS"
    admin_first_name      = "test"
    admin_last_name       = "test"
    admin_phone_num       = "12345678901"
    admin_email           = "test@tencent.com"
    admin_position        = "dev"
    contact_first_name    = "test"
    contact_last_name     = "test"
    contact_email         = "test@tencent.com"
    contact_number        = "12345678901"
    contact_position      = "dev"
  }
}
resource "tencentcloud_ssl_commit_certificate_information" "example" {
  product_id = 33
  certificate_id           = tencentcloud_ssl_pay_certificate.example.certificate_id
}
```

ssl commit_certificate_information can be imported using the id, e.g.

*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudSslCommitCertificateInformation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCommitCertificateInformationCreate,
		Read:   resourceTencentCloudSslCommitCertificateInformationRead,
		Delete: resourceTencentCloudSslCommitCertificateInformationDelete,
		Schema: map[string]*schema.Schema{
			"certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Certificate Id.",
			},

			"verify_type": {
				Optional:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Domain name verification method.",
			},
			"confirm_letter": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The base64-encoded certificate confirmation file should be in jpg, jpeg, png, pdf, and the size should be between 1kb and 1.4M. Note: it only works when product_id is set to 8, 9 or 10.",
			},
		},
	}
}

func resourceTencentCloudSslCommitCertificateInformationCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_commit_certificate_information.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	sslService := SSLService{client: meta.(*TencentCloudClient).apiV3Conn}

	var (
		request       = ssl.NewCommitCertificateInformationRequest()
		productId     = int64(d.Get("product_id").(int))
		certificateId string
	)
	if v, ok := d.GetOk("certificate_id"); ok {
		certificateId = v.(string)
		request.CertificateId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("verify_type"); ok {
		request.VerifyType = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseSSLCertificateClient().CommitCertificateInformation(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ssl commitCertificateInformation failed, reason:%+v", logId, err)
		return err
	}

	if IsContainProductId(productId, GEOTRUST_OV_EV_TYPE) {
		confirmLetter := d.Get("confirm_letter").(string)
		uploadConfirmLetterRequest := ssl.NewUploadConfirmLetterRequest()
		uploadConfirmLetterRequest.CertificateId = helper.String(certificateId)
		uploadConfirmLetterRequest.ConfirmLetter = helper.String(confirmLetter)
		if err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			if err = sslService.UploadConfirmLetter(ctx, uploadConfirmLetterRequest); err != nil {
				if sdkError, ok := err.(*errors.TencentCloudSDKError); ok {
					code := sdkError.GetCode()
					if code == InvalidParam || code == CertificateNotFound {
						return resource.NonRetryableError(sdkError)
					}
				}
				return retryError(err)
			}
			return nil
		}); err != nil {
			return err
		}
	}

	d.SetId(certificateId)

	return nil
}
func resourceTencentCloudSslCommitCertificateInformationRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_pay_certificate.read")()
	defer inconsistentCheck(d, meta)()
	return nil
}
func resourceTencentCloudSslCommitCertificateInformationDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ssl_commit_certificate_information.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
