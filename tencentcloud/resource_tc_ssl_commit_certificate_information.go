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

*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ssl "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ssl/v20191205"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudSslCommitCertificateInformation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudSslCommitCertificateInformationCreate,
		Read:   resourceTencentCloudSslCommitCertificateInformationRead,
		Delete: resourceTencentCloudSslCommitCertificateInformationDelete,
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateIntegerInRange(3, 56),
				Description: "Certificate commodity ID. Valid value ranges: (3~42). `3` means SecureSite enhanced Enterprise Edition (EV Pro), " +
					"`4` means SecureSite enhanced (EV), `5` means SecureSite Enterprise Professional Edition (OV Pro), " +
					"`6` means SecureSite Enterprise (OV), `7` means SecureSite Enterprise Type (OV) wildcard, " +
					"`8` means Geotrust enhanced (EV), `9` means Geotrust enterprise (OV), " +
					"`10` means Geotrust enterprise (OV) wildcard, `11` means TrustAsia domain type multi-domain SSL certificate, " +
					"`12` means TrustAsia domain type ( DV) wildcard, `13` means TrustAsia enterprise wildcard (OV) SSL certificate (D3), " +
					"`14` means TrustAsia enterprise (OV) SSL certificate (D3), `15` means TrustAsia enterprise multi-domain (OV) SSL certificate (D3), " +
					"`16` means TrustAsia Enhanced (EV) SSL Certificate (D3), `17` means TrustAsia Enhanced Multiple Domain (EV) SSL Certificate (D3), " +
					"`18` means GlobalSign Enterprise (OV) SSL Certificate, `19` means GlobalSign Enterprise Wildcard (OV) SSL Certificate, " +
					"`20` means GlobalSign Enhanced (EV) SSL Certificate, `21` means TrustAsia Enterprise Wildcard Multiple Domain (OV) SSL Certificate (D3), " +
					"`22` means GlobalSign Enterprise Multiple Domain (OV) SSL Certificate, `23` means GlobalSign Enterprise Multiple Wildcard Domain name (OV) SSL certificate, " +
					"`24` means GlobalSign enhanced multi-domain (EV) SSL certificate, `25` means Wotrus domain type certificate, " +
					"`26` means Wotrus domain type multi-domain certificate, `27` means Wotrus domain type wildcard certificate, " +
					"`28` means Wotrus enterprise type certificate, `29` means Wotrus enterprise multi-domain certificate, " +
					"`30` means Wotrus enterprise wildcard certificate, `31` means Wotrus enhanced certificate, " +
					"`32` means Wotrus enhanced multi-domain certificate, `33` means WoTrus National Secret Domain name Certificate, " +
					"`34` means WoTrus National Secret Domain name Certificate (multiple domain names), `35` WoTrus National Secret Domain name Certificate (wildcard), " +
					"`37` means WoTrus State Secret Enterprise Certificate, `38` means WoTrus State Secret Enterprise Certificate (multiple domain names), " +
					"`39` means WoTrus State Secret Enterprise Certificate (wildcard), `40` means WoTrus National secret enhanced certificate, " +
					"`41` means WoTrus National Secret enhanced Certificate (multiple domain names), `42` means TrustAsia- Domain name Certificate (wildcard multiple domain names), " +
					"`43` means DNSPod Enterprise (OV) SSL Certificate, `44` means DNSPod- Enterprise (OV) wildcard SSL certificate, " +
					"`45` means DNSPod Enterprise (OV) Multi-domain name SSL Certificate, `46` means DNSPod enhanced (EV) SSL certificate, " +
					"`47` means DNSPod enhanced (EV) multi-domain name SSL certificate, `48` means DNSPod Domain name Type (DV) SSL Certificate, " +
					"`49` means DNSPod Domain name Type (DV) wildcard SSL certificate, `50` means DNSPod domain name type (DV) multi-domain name SSL certificate, " +
					"`51` means DNSPod (State Secret) Enterprise (OV) SSL certificate, `52` DNSPod (National Secret) Enterprise (OV) wildcard SSL certificate, " +
					"`53` means DNSPod (National Secret) Enterprise (OV) multi-domain SSL certificate, `54` means DNSPod (National Secret) Domain Name (DV) SSL certificate, " +
					"`55` means DNSPod (National Secret) Domain Name Type (DV) wildcard SSL certificate, `56` means DNSPod (National Secret) Domain Name Type (DV) multi-domain SSL certificate.",
			},
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
