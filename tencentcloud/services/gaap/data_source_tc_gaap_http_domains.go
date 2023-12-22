package gaap

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapHttpDomains() *schema.Resource {
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
							Deprecated:  "It has been deprecated from version 1.26.0. Use `client_certificate_ids` instead.",
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the client certificate.",
						},
						"client_certificate_ids": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "ID list of the client certificate.",
						},
						"realserver_auth": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether realserver authentication is enable.",
						},
						"realserver_certificate_id": {
							Deprecated:  "It has been deprecated from version 1.28.0. Use `realserver_certificate_ids` instead.",
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CA certificate ID of the realserver.",
						},
						"realserver_certificate_ids": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "CA certificate ID list of the realserver.",
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
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_http_domains.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	domainRules, err := service.DescribeDomains(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domainRules))
	domains := make([]map[string]interface{}, 0, len(domainRules))
	for _, dr := range domainRules {
		if dr.CertificateId == nil {
			dr.CertificateId = helper.String("default")
		}
		if dr.ClientCertificateId == nil {
			dr.ClientCertificateId = helper.String("default")
		}
		if dr.RealServerAuth == nil {
			dr.RealServerAuth = helper.IntInt64(0)
		}
		if dr.BasicAuth == nil {
			dr.BasicAuth = helper.IntInt64(0)
		}
		if dr.GaapAuth == nil {
			dr.GaapAuth = helper.IntInt64(0)
		}

		ids = append(ids, *dr.Domain)

		var (
			clientCertificateId      *string
			polyClientCertificateIds []*string
			realserverCertificateIds []*string
		)

		clientCertificateId = dr.PolyClientCertificateAliasInfo[0].CertificateId
		for _, poly := range dr.PolyClientCertificateAliasInfo {
			polyClientCertificateIds = append(polyClientCertificateIds, poly.CertificateId)
		}

		realserverCertificateIds = make([]*string, 0, len(dr.PolyRealServerCertificateAliasInfo))
		for _, info := range dr.PolyRealServerCertificateAliasInfo {
			realserverCertificateIds = append(realserverCertificateIds, info.CertificateId)
		}

		var realserverCertificateId *string
		if len(realserverCertificateIds) > 0 {
			realserverCertificateId = realserverCertificateIds[0]
		}

		if dr.RealServerAuth == nil {
			dr.RealServerAuth = helper.Int64(0)
		}

		if dr.BasicAuth == nil {
			dr.BasicAuth = helper.Int64(0)
		}

		if dr.GaapAuth == nil {
			dr.GaapAuth = helper.Int64(0)
		}

		m := map[string]interface{}{
			"domain":                        dr.Domain,
			"certificate_id":                dr.CertificateId,
			"client_certificate_id":         clientCertificateId,
			"client_certificate_ids":        polyClientCertificateIds,
			"realserver_auth":               *dr.RealServerAuth == 1,
			"basic_auth":                    *dr.BasicAuth == 1,
			"basic_auth_id":                 dr.BasicAuthConfId,
			"gaap_auth":                     *dr.GaapAuth == 1,
			"gaap_auth_id":                  dr.GaapCertificateId,
			"realserver_certificate_id":     realserverCertificateId,
			"realserver_certificate_ids":    realserverCertificateIds,
			"realserver_certificate_domain": dr.RealServerCertificateDomain,
		}

		domains = append(domains, m)
	}

	_ = d.Set("domains", domains)

	d.SetId(helper.DataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), domains); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
