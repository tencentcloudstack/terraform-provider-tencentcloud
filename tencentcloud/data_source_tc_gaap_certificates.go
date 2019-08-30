/*
Use this data source to query GAAP certificate.

Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = 0
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}
data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_gaap_certificate.foo.id}"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
)

func dataSourceTencentCloudGaapCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapCertificatesRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the certificate to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the certificate to be queried.",
			},
			"type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(0, 4),
				Description:  "Type of the certificate to be queried. `0` means basic authentication; `1` means client CA certificate; `2` means server SSL certificate; `3` means realserver CA certificate; `4` means proxy SSL certificate.",
			},

			// computed
			"certificates": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of certificate. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the certificate.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the certificate.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Type of the certificate.",
						},
						"create_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Creation time of the certificate, use unix timestamp format.",
						},
						"begin_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Beginning time of the certificate, use unix timestamp format.",
						},
						"end_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Ending time of the certificate, use unix timestamp format.",
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
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_certificates.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var (
		id              *string
		name            *string
		certificateType *int
		ids             []string
		certificates    []map[string]interface{}
	)

	if raw, ok := d.GetOk("id"); ok {
		id = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}
	if raw, ok := d.GetOk("type"); ok {
		certificateType = common.IntPtr(raw.(int))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	respCertificates, err := service.DescribeCertificates(ctx, id, name, certificateType)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(respCertificates))
	certificates = make([]map[string]interface{}, 0, len(respCertificates))
	for _, certificate := range respCertificates {
		ids = append(ids, *certificate.CertificateId)

		m := map[string]interface{}{
			"id":          *certificate.CertificateId,
			"name":        *certificate.CertificateAlias,
			"type":        *certificate.CertificateType,
			"create_time": *certificate.CreateTime,
		}

		if certificate.BeginTime != nil {
			m["begin_time"] = *certificate.BeginTime
		}

		if certificate.EndTime != nil {
			m["end_time"] = *certificate.EndTime
		}

		if certificate.IssuerCN != nil {
			m["issuer_cn"] = *certificate.IssuerCN
		}

		if certificate.SubjectCN != nil {
			m["subject_cn"] = *certificate.SubjectCN
		}

		certificates = append(certificates, m)
	}

	d.Set("certificates", certificates)
	d.SetId(dataResourceIdsHash(ids))

	return nil
}
