/*
Use this data source to query forward domain of layer7 listeners.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

data "tencentcloud_gaap_http_domains" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "${tencentcloud_gaap_http_domain.foo.domain}"
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapHttpDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapHttpDomainsRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the layer7 listener to be queried.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Forward domain of the layer7 listener to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"domains": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of forward domain of the layer7 listeners. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forward domain of the layer7 listener.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the server certificate.",
						},
						"client_certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the client certificate.",
						},
						"realserver_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether realserver authentication is enable.",
						},
						"realserver_certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CA certificate ID of the realserver.",
						},
						"realserver_certificate_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CA certificate domain of the realserver.",
						},
						"basic_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether basic authentication is enable.",
						},
						"basic_auth_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the basic authentication.",
						},
						"gaap_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether SSL certificate authentication is enable.",
						},
						"gaap_auth_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the SSL certificate.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapHttpDomainsRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_http_domains.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	domainRules, err := service.DescribeDomains(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domainRules))
	domains := make([]map[string]interface{}, 0, len(domainRules))
	for _, dr := range domainRules {
		if dr.CertificateId == nil {
			dr.CertificateId = stringToPointer("default")
		}
		if dr.ClientCertificateId == nil {
			dr.ClientCertificateId = stringToPointer("default")
		}
		if dr.RealServerAuth == nil {
			dr.RealServerAuth = int64ToPointer(0)
		}
		if dr.BasicAuth == nil {
			dr.BasicAuth = int64ToPointer(0)
		}
		if dr.GaapAuth == nil {
			dr.GaapAuth = int64ToPointer(0)
		}

		ids = append(ids, *dr.Domain)

		m := map[string]interface{}{
			"domain":                dr.Domain,
			"certificate_id":        dr.CertificateId,
			"client_certificate_id": dr.ClientCertificateId,
			"realserver_auth":       *dr.RealServerAuth == 1,
			"basic_auth":            *dr.BasicAuth == 1,
			"gaap_auth":             *dr.GaapAuth == 1,
		}

		if dr.RealServerCertificateId != nil {
			m["realserver_certificate_id"] = *dr.RealServerCertificateId
		}
		if dr.RealServerCertificateDomain != nil {
			m["realserver_certificate_domain"] = *dr.RealServerCertificateDomain
		}
		if dr.BasicAuthConfId != nil {
			m["basic_auth_id"] = *dr.BasicAuthConfId
		}
		if dr.GaapCertificateId != nil {
			m["gaap_auth_id"] = *dr.GaapCertificateId
		}

		domains = append(domains, m)
	}

	d.Set("domains", domains)
	d.SetId(dataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), domains); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
