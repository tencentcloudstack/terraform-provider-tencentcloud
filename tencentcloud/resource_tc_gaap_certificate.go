/*
Provides a resource to create a certificate of GAAP.

Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}
```

Import

GAAP certificate can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_certificate.foo cert-d5y6ei3b
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapCertificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapCertificateCreate,
		Read:   resourceTencentCloudGaapCertificateRead,
		Update: resourceTencentCloudGaapCertificateUpdate,
		Delete: resourceTencentCloudGaapCertificateDelete,
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{"BASIC", "CLIENT", "SERVER", "REALSERVER", "PROXY"}),
				Description:  "Type of the certificate. Valid value: `BASIC`, `CLIENT`, `SERVER`, `REALSERVER` and `PROXY`. `BASIC` means basic certificate; `CLIENT` means client CA certificate; `SERVER` means server SSL certificate; `REALSERVER` means realserver CA certificate; `PROXY` means proxy SSL certificate.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Content of the certificate, and URL encoding. When the certificate is basic authentication, use the `user:xxx password:xxx` format, where the password is encrypted with `htpasswd` or `openssl`; When the certificate is `CA` or `SSL`, the format is `pem`.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Name of the certificate.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "Key of the `SSL` certificate.",
			},

			// computed
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time of the certificate.",
			},
			"begin_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Beginning time of the certificate.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Ending time of the certificate.",
			},
			"issuer_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Issuer name of the certificate.",
			},
			"subject_cn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Subject name of the certificate.",
			},
		},
	}
}

func resourceTencentCloudGaapCertificateCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	certificateType := gaapCertificateStringMap[d.Get("type").(string)]

	content := d.Get("content").(string)

	name := d.Get("name").(string)

	var key *string
	if rawKey, ok := d.GetOk("key"); ok {
		key = helper.String(rawKey.(string))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.createCertificate(ctx, certificateType, content, name, key)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	certificate, err := service.DescribeCertificateById(ctx, id)
	if err != nil {
		return err
	}

	if certificate == nil {
		d.SetId("")
		return nil
	}

	if certificate.CertificateType == nil {
		return errors.New("certificate type is nil")
	}
	if certType, ok := gaapCertificateIntMap[int(*certificate.CertificateType)]; ok {
		_ = d.Set("type", certType)
	} else {
		return fmt.Errorf("unknown certificate type %d", *certificate.CertificateType)
	}

	_ = d.Set("name", certificate.CertificateAlias)

	if certificate.CreateTime == nil {
		return errors.New("certificate create time is nil")
	}
	_ = d.Set("create_time", helper.FormatUnixTime(*certificate.CreateTime))

	if certificate.BeginTime != nil {
		_ = d.Set("begin_time", helper.FormatUnixTime(*certificate.BeginTime))
	}
	if certificate.EndTime != nil {
		_ = d.Set("end_time", helper.FormatUnixTime(*certificate.EndTime))
	}

	_ = d.Set("issuer_cn", certificate.IssuerCN)
	_ = d.Set("subject_cn", certificate.SubjectCN)

	return nil
}

func resourceTencentCloudGaapCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	name := d.Get("name").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if err := service.ModifyCertificateName(ctx, id, name); err != nil {
		return err
	}

	return resourceTencentCloudGaapCertificateRead(d, m)
}

func resourceTencentCloudGaapCertificateDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_certificate.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteCertificate(ctx, id)
}
