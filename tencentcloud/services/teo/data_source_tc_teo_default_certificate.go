package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoDefaultCertificate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoDefaultCertificateRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},

			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions, the upper limit of Filters.Values is 5. The detailed filtering conditions are as follows: zone-id - Filter by zone ID.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter name.",
						},
						"values": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Filter value.",
						},
					},
				},
			},

			"default_server_cert_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Default certificate list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Server certificate ID.",
						},
						"alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate alias.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate type. Valid values: default (default certificate), upload (user uploaded), managed (Tencent Cloud managed).",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate expiration time.",
						},
						"effective_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate effective time.",
						},
						"common_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate common name.",
						},
						"subject_alt_name": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Certificate SAN domains.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Deploy status. Valid values: processing (deploying), deployed (deployed), failed (deploy failed).",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure reason when Status is failed.",
						},
						"sign_algo": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate signing algorithm.",
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoDefaultCertificateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_default_certificate.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*teo.Filter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filter := teo.Filter{}
			filterMap := item.(map[string]interface{})
			if v, ok := filterMap["name"].(string); ok && v != "" {
				filter.Name = helper.String(v)
			}

			if v, ok := filterMap["values"]; ok {
				valuesSet := v.(*schema.Set).List()
				filter.Values = helper.InterfacesStringsPoint(valuesSet)
			}

			tmpSet = append(tmpSet, &filter)
		}

		paramMap["Filters"] = tmpSet
	}

	var certificates []*teo.DefaultServerCertInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoDefaultCertificatesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		certificates = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(certificates))
	certInfoList := make([]map[string]interface{}, 0, len(certificates))
	if certificates != nil {
		for _, certInfo := range certificates {
			certInfoMap := map[string]interface{}{}

			if certInfo.CertId != nil {
				certInfoMap["cert_id"] = certInfo.CertId
				ids = append(ids, *certInfo.CertId)
			}

			if certInfo.Alias != nil {
				certInfoMap["alias"] = certInfo.Alias
			}

			if certInfo.Type != nil {
				certInfoMap["type"] = certInfo.Type
			}

			if certInfo.ExpireTime != nil {
				certInfoMap["expire_time"] = certInfo.ExpireTime
			}

			if certInfo.EffectiveTime != nil {
				certInfoMap["effective_time"] = certInfo.EffectiveTime
			}

			if certInfo.CommonName != nil {
				certInfoMap["common_name"] = certInfo.CommonName
			}

			if certInfo.SubjectAltName != nil {
				certInfoMap["subject_alt_name"] = certInfo.SubjectAltName
			}

			if certInfo.Status != nil {
				certInfoMap["status"] = certInfo.Status
			}

			if certInfo.Message != nil {
				certInfoMap["message"] = certInfo.Message
			}

			if certInfo.SignAlgo != nil {
				certInfoMap["sign_algo"] = certInfo.SignAlgo
			}

			certInfoList = append(certInfoList, certInfoMap)
		}

		_ = d.Set("default_server_cert_info", certInfoList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), certInfoList); e != nil {
			return e
		}
	}

	return nil
}
