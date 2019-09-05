/*
Use this data source to query SSL certificate.

Example Usage

```hcl
resource "tencentcloud_ssl_certificate" "foo" {
  name       = "test-ssl-certificate"
  type       = "CA"
  project_id = 0
  cert       = "-----BEGIN CERTIFICATE-----\nMIIERzCCAq+gAwIBAgIBAjANBgkqhkiG9w0BAQsFADAoMQ0wCwYDVQQDEwR0ZXN0\nMRcwFQYDVQQKEw50ZXJyYWZvcm0gdGVzdDAeFw0xOTA4MTMwMzE5MzlaFw0yOTA4\nMTAwMzE5MzlaMC4xEzARBgNVBAMTCnNlcnZlciBzc2wxFzAVBgNVBAoTDnRlcnJh\nZm9ybS10ZXN0MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEA1Ryp+DKK\nSNFKZsPtwfR+jzOnQ8YFieIKYgakV688d8YgpolenbmeEPrzT87tunFD7G9f6ALG\nND8rj7npj0AowxhOL/h/v1D9u0UsIaj5i2GWJrqNAhGLaxWiEB/hy5WOiwxDrGei\ngQqJkFM52Ep7G1Yx7PHJmKFGwN9FhIsFi1cNZfVRopZuCe/RMPNusNVZaIi+qcEf\nfsE1cmfmuSlG3Ap0RKOIyR0ajDEzqZn9/0R7VwWCF97qy8TNYk94K/1tq3zyhVzR\nZ83xOSfrTqEfb3so3AU2jyKgYdwr/FZS72VCHS8IslgnqJW4izIXZqgIKmHaRZtM\nN4jUloi6l/6lktt6Lsgh9xECecxziSJtPMaog88aC8HnMqJJ3kScGCL36GYG+Kaw\n5PnDlWXBaeiDe8z/eWK9+Rr2M+rhTNxosAVGfDJyxAXyiX49LQ0v7f9qzwc/0JiD\nbvsUv1cm6OgpoEMP9SXqqBdwGqeKbD2/2jlP48xlYP6l1SoJG3GgZ8dbAgMBAAGj\ndjB0MAwGA1UdEwEB/wQCMAAwEwYDVR0lBAwwCgYIKwYBBQUHAwEwDwYDVR0PAQH/\nBAUDAweAADAdBgNVHQ4EFgQULwWKBQNLL9s3cb3tTnyPVg+mpCMwHwYDVR0jBBgw\nFoAUKwfrmq791mY831S6UHARHtgYnlgwDQYJKoZIhvcNAQELBQADggGBAMo5RglS\nAHdPgaicWJvmvjjexjF/42b7Rz4pPfMjYw6uYO8He/f4UZWv5CZLrbEe7MywaK3y\n0OsfH8AhyN29pv2x8g9wbmq7omZIOZ0oCAGduEXs/A/qY/hFaCohdkz/IN8qi6JW\nVXreGli3SrpcHFchSwHTyJEXgkutcGAsOvdsOuVSmplOyrkLHc8uUe8SG4j8kGyg\nEzaszFjHkR7g1dVyDVUedc588mjkQxYeAamJgfkgIhljWKMa2XzkVMcVfQHfNpM1\nn+bu8SmqRt9Wma2bMijKRG/Blm756LoI+skY+WRZmlDnq8zj95TT0vceGP0FUWh5\nhKyiocABmpQs9OK9HMi8vgSWISP+fYgkm/bKtKup2NbZBoO5/VL2vCEPInYzUhBO\njCbLMjNjtM5KriCaR7wDARgHiG0gBEPOEW1PIjZ9UOH+LtIxbNZ4eEIIINLHnBHf\nL+doVeZtS/gJc4G4Adr5HYuaS9ZxJ0W2uy0eQlOHzjyxR6Mf/rpnilJlcQ==\n-----END CERTIFICATE-----"
}

data "tencentcloud_gaap_certificates" "foo" {
  id = "${tencentcloud_ssl_certificate.foo.id}"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudSslCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSslCertificatesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the SSL certificate to be queried.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Type of the SSL certificate to be queried. Available values includes: `CA` and `SVR`.",
			},
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the SSL certificate to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"certificates": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SSL certificate.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the SSL certificate.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the SSL certificate.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project ID of the SSL certificate.",
						},
						"cert": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Content of the SSL certificate.",
						},
						"product_zh_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate authority.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary domain of the SSL certificate.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the SSL certificate.",
						},
						"begin_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Beginning time of the SSL certificate.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ending time of the SSL certificate.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the SSL certificate.",
						},
						"subject_names": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "ALL domains included in the SSL certificate. Including the primary domain name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSslCertificatesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssl_certificates.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	var (
		name     *string
		certType *string
		id       *string
	)

	if raw, ok := d.GetOk("name"); ok {
		name = stringToPointer(raw.(string))
	}

	if raw, ok := d.GetOk("type"); ok {
		certType = stringToPointer(raw.(string))
	}

	if raw, ok := d.GetOk("id"); ok {
		id = stringToPointer(raw.(string))
	}

	service := SslService{client: m.(*TencentCloudClient).apiV3Conn}

	respCertificates, err := service.DescribeCertificates(ctx, id, name, certType)
	if err != nil {
		return err
	}

	certificates := make([]map[string]interface{}, 0, len(respCertificates))
	ids := make([]string, 0, len(respCertificates))

	for _, certificate := range respCertificates {
		if certificate.Id == nil {
			return errors.New("certificate id is nil")
		}

		ids = append(ids, *certificate.Id)

		if certificate.Alias == nil {
			return errors.New("certificate name is nil")
		}

		if certificate.CertType == nil {
			return errors.New("certificate type is nil")
		}

		if certificate.ProjectId == nil {
			return errors.New("certificate project id is nil")
		}
		projectId, err := strconv.Atoi(*certificate.ProjectId)
		if err != nil {
			return err
		}

		if certificate.Cert == nil {
			return errors.New("certificate cert is nil")
		}

		if certificate.ProductZhName == nil {
			return errors.New("certificate product zh name is nil")
		}

		if certificate.Domain == nil {
			return errors.New("certificate domain is nil")
		}

		if certificate.Status == nil {
			return errors.New("certificate status is nil")
		}

		if certificate.CertBeginTime == nil {
			return errors.New("certificate begin time is nil")
		}

		if certificate.CertEndTime == nil {
			return errors.New("certificate end time is nil")
		}

		if certificate.InsertTime == nil {
			return errors.New("certificate create time is nil")
		}

		m := map[string]interface{}{
			"id":              *certificate.Id,
			"name":            *certificate.Alias,
			"type":            *certificate.CertType,
			"project_id":      projectId,
			"cert":            *certificate.Cert,
			"product_zh_name": *certificate.ProductZhName,
			"domain":          *certificate.Domain,
			"status":          *certificate.Status,
			"begin_time":      *certificate.CertBeginTime,
			"end_time":        *certificate.CertEndTime,
			"create_time":     *certificate.InsertTime,
		}

		if len(certificate.SubjectAltName) > 0 {
			subjectAltNames := make([]string, 0, len(certificate.SubjectAltName))
			for _, name := range certificate.SubjectAltName {
				subjectAltNames = append(subjectAltNames, *name)
			}
			m["subject_names"] = subjectAltNames
		}

		certificates = append(certificates, m)
	}

	d.Set("certificates", certificates)
	d.SetId(dataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), certificates); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
