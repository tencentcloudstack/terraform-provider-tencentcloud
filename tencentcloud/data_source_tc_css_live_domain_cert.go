/*
Use this data source to query detailed information of css live_domain_cert

Example Usage

```hcl
data "tencentcloud_css_live_domain_cert" "live_domain_cert" {
  domain_name = "your_domain_name"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCssLiveDomainCert() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssLiveDomainCertRead,
		Schema: map[string]*schema.Schema{
			"domain_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Playback domain name.",
			},

			"domain_cert_info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Certificate information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cert_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Certificate ID.",
						},
						"cert_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate name.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time in UTC format.Note: Beijing time (UTC+8) is used.",
						},
						"https_crt": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Certificate content.",
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Certificate type.0: user-added certificate1: Tencent Cloud-hosted certificate.",
						},
						"cert_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The certificate expiration time in UTC format.Note: Beijing time (UTC+8) is used.",
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name that uses this certificate.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Certificate status.",
						},
						"cert_domains": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "List of domain names in the certificate.[*.x.com] for example.Note: this field may return `null`, indicating that no valid values can be obtained.",
						},
						"cloud_cert_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Tencent Cloud SSL certificate ID.Note: this field may return `null`, indicating that no valid values can be obtained.",
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

func dataSourceTencentCloudCssLiveDomainCertRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_css_live_domain_cert.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain_name"); ok {
		paramMap["DomainName"] = helper.String(v.(string))
	}

	service := CssService{client: meta.(*TencentCloudClient).apiV3Conn}

	var domainCertInfo *css.DomainCertInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssLiveDomainCertByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		domainCertInfo = result
		return nil
	})
	if err != nil {
		return err
	}

	var id string
	domainCertInfoMap := map[string]interface{}{}
	if domainCertInfo != nil {
		if domainCertInfo.CertId != nil {
			domainCertInfoMap["cert_id"] = domainCertInfo.CertId
		}

		if domainCertInfo.CertName != nil {
			domainCertInfoMap["cert_name"] = domainCertInfo.CertName
		}

		if domainCertInfo.Description != nil {
			domainCertInfoMap["description"] = domainCertInfo.Description
		}

		if domainCertInfo.CreateTime != nil {
			domainCertInfoMap["create_time"] = domainCertInfo.CreateTime
		}

		if domainCertInfo.HttpsCrt != nil {
			domainCertInfoMap["https_crt"] = domainCertInfo.HttpsCrt
		}

		if domainCertInfo.CertType != nil {
			domainCertInfoMap["cert_type"] = domainCertInfo.CertType
		}

		if domainCertInfo.CertExpireTime != nil {
			domainCertInfoMap["cert_expire_time"] = domainCertInfo.CertExpireTime
		}

		if domainCertInfo.DomainName != nil {
			domainCertInfoMap["domain_name"] = domainCertInfo.DomainName
		}

		if domainCertInfo.Status != nil {
			domainCertInfoMap["status"] = domainCertInfo.Status
		}

		if domainCertInfo.CertDomains != nil {
			domainCertInfoMap["cert_domains"] = domainCertInfo.CertDomains
		}

		if domainCertInfo.CloudCertId != nil {
			domainCertInfoMap["cloud_cert_id"] = domainCertInfo.CloudCertId
		}

		id = helper.Int64ToStr(*domainCertInfo.CertId)
	}
	_ = d.Set("domain_cert_info", []interface{}{domainCertInfoMap})

	d.SetId(id)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), domainCertInfoMap); e != nil {
			return e
		}
	}
	return nil
}
