/*
Use this data source to query GAAP certificate.

Example Usage

```hcl
resource "tencentcloud_gaap_certificate" "foo" {
  type    = "BASIC"
  content = "test:tx2KGdo3zJg/."
  name    = "test_certificate"
}

data "tencentcloud_gaap_certificates" "foo" {
  id = tencentcloud_gaap_certificate.foo.id
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{"BASIC", "CLIENT", "SERVER", "REALSERVER", "PROXY"}),
				Description:  "Type of the certificate to be queried, the available values include `BASIC`, `CLIENT`, `SERVER`, `REALSERVER` and `PROXY`; `BASIC` means basic certificate; `CLIENT` means client CA certificate; `SERVER` means server SSL certificate; `REALSERVER` means realserver CA certificate; `PROXY` means proxy SSL certificate.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
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
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the certificate.",
						},
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
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_certificates.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var (
		id              *string
		name            *string
		certificateType *int
		ids             []string
		certificates    []map[string]interface{}
	)

	if raw, ok := d.GetOk("id"); ok {
		id = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("name"); ok {
		name = helper.String(raw.(string))
	}
	if raw, ok := d.GetOk("type"); ok {
		certificateType = common.IntPtr(gaapCertificateStringMap[raw.(string)])
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

		var (
			certificateType string
			ok              bool
		)
		if certificateType, ok = gaapCertificateIntMap[int(*certificate.CertificateType)]; !ok {
			return fmt.Errorf("unknown certificate type %d", *certificate.CertificateType)
		}

		m := map[string]interface{}{
			"id":          *certificate.CertificateId,
			"name":        *certificate.CertificateAlias,
			"type":        certificateType,
			"create_time": helper.FormatUnixTime(*certificate.CreateTime),
		}

		if certificate.BeginTime != nil {
			m["begin_time"] = helper.FormatUnixTime(*certificate.BeginTime)
		}

		if certificate.EndTime != nil {
			m["end_time"] = helper.FormatUnixTime(*certificate.EndTime)
		}

		if certificate.IssuerCN != nil {
			m["issuer_cn"] = *certificate.IssuerCN
		}

		if certificate.SubjectCN != nil {
			m["subject_cn"] = *certificate.SubjectCN
		}

		certificates = append(certificates, m)
	}

	_ = d.Set("certificates", certificates)
	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), certificates); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
