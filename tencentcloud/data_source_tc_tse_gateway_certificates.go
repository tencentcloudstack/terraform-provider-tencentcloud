/*
Use this data source to query detailed information of tse gateway_certificates

Example Usage

```hcl

data "tencentcloud_tse_gateway_certificates" "gateway_certificates" {
  gateway_id = "gateway-ddbb709b"
  filters {
    key = "BindDomain"
    value = "example.com"
  }
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayCertificates() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayCertificatesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions, valid value: `BindDomain`, `Name`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count. Note: This field may return null, indicating that a valid value is not available.",
						},
						"certificates_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Certificate list of gateway. Note: This field may return null, indicating that a valid value is not available.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate name. Note: This field may return null, indicating that a valid value is not available.",
									},
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID. Note: This field may return null, indicating that a valid value is not available.",
									},
									"bind_domains": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Domains of the binding. Note: This field may return null, indicating that a valid value is not available.",
									},
									"status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Status of certificate. Reference value:- expired- active. Note: This field may return null, indicating that a valid value is not available.",
									},
									"crt": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Pem format of certificate. Note: This field may return null, indicating that a valid value is not available.",
									},
									"key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Private key of certificate. Note: This field may return null, indicating that a valid value is not available.",
									},
									"expire_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Expiration time of certificate. Note: This field may return null, indicating that a valid value is not available.",
									},
									"create_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Upload time of certificate. Note: This field may return null, indicating that a valid value is not available.",
									},
									"issue_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Issuance time of certificateNote: This field may return null, indicating that a valid value is not available.",
									},
									"cert_source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Source of certificate. Reference value:- native. Source: konga- ssl. Source: ssl platform. Note: This field may return null, indicating that a valid value is not available.",
									},
									"cert_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Certificate ID of ssl platform. Note: This field may return null, indicating that a valid value is not available.",
									},
								},
							},
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

func dataSourceTencentCloudTseGatewayCertificatesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_certificates.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.ListFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			listFilter := tse.ListFilter{}
			listFilterMap := item.(map[string]interface{})

			if v, ok := listFilterMap["key"]; ok {
				listFilter.Key = helper.String(v.(string))
			}
			if v, ok := listFilterMap["value"]; ok {
				listFilter.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &listFilter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result *tse.KongCertificatesList
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := service.DescribeTseGatewayCertificatesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = response
		return nil
	})
	if err != nil {
		return err
	}

	var ids []string
	kongCertificatesListMap := map[string]interface{}{}
	if result != nil {
		ids = make([]string, 0, *result.Total)
		if result.Total != nil {
			kongCertificatesListMap["total"] = result.Total
		}

		if result.CertificatesList != nil {
			certificatesListList := []interface{}{}
			for _, certificatesList := range result.CertificatesList {
				certificatesListMap := map[string]interface{}{}

				if certificatesList.Name != nil {
					certificatesListMap["name"] = certificatesList.Name
				}

				if certificatesList.Id != nil {
					certificatesListMap["id"] = certificatesList.Id
				}

				if certificatesList.BindDomains != nil {
					certificatesListMap["bind_domains"] = certificatesList.BindDomains
				}

				if certificatesList.Status != nil {
					certificatesListMap["status"] = certificatesList.Status
				}

				if certificatesList.Crt != nil {
					certificatesListMap["crt"] = certificatesList.Crt
				}

				if certificatesList.Key != nil {
					certificatesListMap["key"] = certificatesList.Key
				}

				if certificatesList.ExpireTime != nil {
					certificatesListMap["expire_time"] = certificatesList.ExpireTime
				}

				if certificatesList.CreateTime != nil {
					certificatesListMap["create_time"] = certificatesList.CreateTime
				}

				if certificatesList.IssueTime != nil {
					certificatesListMap["issue_time"] = certificatesList.IssueTime
				}

				if certificatesList.CertSource != nil {
					certificatesListMap["cert_source"] = certificatesList.CertSource
				}

				if certificatesList.CertId != nil {
					certificatesListMap["cert_id"] = certificatesList.CertId
				}

				certificatesListList = append(certificatesListList, certificatesListMap)
				ids = append(ids, *certificatesList.Id)
			}

			kongCertificatesListMap["certificates_list"] = certificatesListList
		}

		_ = d.Set("result", []interface{}{kongCertificatesListMap})
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), kongCertificatesListMap); e != nil {
			return e
		}
	}
	return nil
}
