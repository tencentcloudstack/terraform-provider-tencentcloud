package tencentcloud

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapHttpDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapHttpDomainsRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"realserver_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_certificate_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"basic_auth_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gaap_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gaap_auth_id": {
							Type:     schema.TypeString,
							Computed: true,
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
			return errors.New("domain certificate id is nil")
		}
		if dr.ClientCertificateId == nil {
			return errors.New("domain client certificate id is nil")
		}
		if dr.RealServerAuth == nil {
			return errors.New("domain realserver auth is nil")
		}
		if dr.BasicAuth == nil {
			return errors.New("domain basic auth is nil")
		}
		if dr.GaapAuth == nil {
			return errors.New("domain gaap auth is nil")
		}

		ids = append(ids, *dr.Domain)

		m := map[string]interface{}{
			"domain":                *dr.Domain,
			"certificate_id":        *dr.CertificateId,
			"client_certificate_id": *dr.ClientCertificateId,
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

	return nil
}
