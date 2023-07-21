/*
Provides a resource to create a clb replace_cert_for_lbs

Example Usage

Replace Server Cert By Cert ID
```hcl
resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_id = "6vcK02GC"
  }
}
```

Replace Server Cert By Cert Content
```hcl
data "tencentcloud_ssl_certificates" "foo" {
  name = "keep-ssl-ca"
}

resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = data.tencentcloud_ssl_certificates.foo.certificates.0.id
  certificate {
    cert_name    = "tf-test-cert"
    cert_content = <<-EOT
-----BEGIN CERTIFICATE-----
xxxxxxxxxxxxxxxxxxxxxxxxxxx
-----END CERTIFICATE-----
EOT
    cert_key     = <<-EOT
-----BEGIN RSA PRIVATE KEY-----
xxxxxxxxxxxxxxxxxxxxxxxxxxxx
-----END RSA PRIVATE KEY-----
EOT
  }
}
```

Replace Client Cert By Cert Content
```hcl
resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_ca_name = "tf-test-cert"
    cert_ca_content = <<-EOT
-----BEGIN CERTIFICATE-----
xxxxxxxxContentxxxxxxxxxxxxxx
-----END CERTIFICATE-----
EOT
  }
}
```

```
terraform import tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs replace_cert_for_lbs_id
```
*/
package tencentcloud

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClbReplaceCertForLbs() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbReplaceCertForLbsCreate,
		Read:   resourceTencentCloudClbReplaceCertForLbsRead,
		Delete: resourceTencentCloudClbReplaceCertForLbsDelete,
		Schema: map[string]*schema.Schema{
			"old_certificate_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of the certificate to be replaced, which can be a server certificate or a client certificate.",
			},

			"certificate": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Information such as the content of the new certificate.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_mode": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Authentication type. Value range: UNIDIRECTIONAL (unidirectional authentication), MUTUAL (mutual authentication).",
						},
						"cert_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a server certificate. If you leave this parameter empty, you must upload the certificate, including CertContent, CertKey, and CertName.",
						},
						"cert_ca_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of a client certificate. When the listener adopts mutual authentication (i.e., SSLMode = mutual), if you leave this parameter empty, you must upload the client certificate, including CertCaContent and CertCaName.",
						},
						"cert_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the uploaded server certificate. If there is no CertId, this parameter is required.",
						},
						"cert_key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Key of the uploaded server certificate. If there is no CertId, this parameter is required.",
						},
						"cert_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content of the uploaded server certificate. If there is no CertId, this parameter is required.",
						},
						"cert_ca_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Name of the uploaded client CA certificate. When SSLMode = mutual, if there is no CertCaId, this parameter is required.",
						},
						"cert_ca_content": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content of the uploaded client certificate. When SSLMode = mutual, if there is no CertCaId, this parameter is required.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClbReplaceCertForLbsCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request          = clb.NewReplaceCertForLoadBalancersRequest()
		oldCertificateId string
	)
	if v, ok := d.GetOk("old_certificate_id"); ok {
		oldCertificateId = v.(string)
		request.OldCertificateId = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "certificate"); ok {
		certificateInput := clb.CertificateInput{}
		if v, ok := dMap["ssl_mode"]; ok {
			certificateInput.SSLMode = helper.String(v.(string))
		}
		if v, ok := dMap["cert_id"]; ok {
			certificateInput.CertId = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_id"]; ok {
			certificateInput.CertCaId = helper.String(v.(string))
		}
		if v, ok := dMap["cert_name"]; ok {
			certificateInput.CertName = helper.String(v.(string))
		}
		if v, ok := dMap["cert_key"]; ok {
			certificateInput.CertKey = helper.String(v.(string))
		}
		if v, ok := dMap["cert_content"]; ok {
			certificateInput.CertContent = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_name"]; ok {
			certificateInput.CertCaName = helper.String(v.(string))
		}
		if v, ok := dMap["cert_ca_content"]; ok {
			certificateInput.CertCaContent = helper.String(v.(string))
		}
		request.Certificate = &certificateInput
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ReplaceCertForLoadBalancers(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s operate clb replaceCertForLbs failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(oldCertificateId)

	return resourceTencentCloudClbReplaceCertForLbsRead(d, meta)
}

func resourceTencentCloudClbReplaceCertForLbsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.read")()
	defer inconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudClbReplaceCertForLbsDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_clb_replace_cert_for_lbs.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
