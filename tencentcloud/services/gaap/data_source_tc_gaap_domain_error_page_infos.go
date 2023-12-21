package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapDomainErrorPageInfos() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapDomainErrorPageInfosRead,
		Schema: map[string]*schema.Schema{
			"error_page_ids": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Customized error ID list, supporting up to 10.",
			},

			"error_page_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Custom error response configuration setNote: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_page_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Configuration ID for error customization response.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Listener ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "domain name.",
						},
						"error_nos": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							Computed:    true,
							Description: "Original error code.",
						},
						"new_error_no": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "New error codeNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"clear_headers": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed:    true,
							Description: "Response headers that need to be cleanedNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"set_headers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Response header to be setNote: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"header_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP header name.",
									},
									"header_value": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "HTTP header value.",
									},
								},
							},
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Response body set (excluding HTTP header)Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule status, 0 indicates successNote: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudGaapDomainErrorPageInfosRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_domain_error_page_infos.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("error_page_ids"); ok {
		errorPageIdsSet := v.(*schema.Set).List()
		paramMap["ErrorPageIds"] = helper.InterfacesStringsPoint(errorPageIdsSet)
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var errorPageSet []*gaap.DomainErrorPageInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapDomainErrorPageInfosByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		errorPageSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(errorPageSet))
	tmpList := make([]map[string]interface{}, 0, len(errorPageSet))

	if errorPageSet != nil {
		for _, domainErrorPageInfo := range errorPageSet {
			domainErrorPageInfoMap := map[string]interface{}{}

			if domainErrorPageInfo.ErrorPageId != nil {
				domainErrorPageInfoMap["error_page_id"] = domainErrorPageInfo.ErrorPageId
				ids = append(ids, *domainErrorPageInfo.ErrorPageId)
			}

			if domainErrorPageInfo.ListenerId != nil {
				domainErrorPageInfoMap["listener_id"] = domainErrorPageInfo.ListenerId
			}

			if domainErrorPageInfo.Domain != nil {
				domainErrorPageInfoMap["domain"] = domainErrorPageInfo.Domain
			}

			if domainErrorPageInfo.ErrorNos != nil {
				domainErrorPageInfoMap["error_nos"] = domainErrorPageInfo.ErrorNos
			}

			if domainErrorPageInfo.NewErrorNo != nil {
				domainErrorPageInfoMap["new_error_no"] = domainErrorPageInfo.NewErrorNo
			}

			if domainErrorPageInfo.ClearHeaders != nil {
				domainErrorPageInfoMap["clear_headers"] = domainErrorPageInfo.ClearHeaders
			}

			if domainErrorPageInfo.SetHeaders != nil {
				setHeadersList := []interface{}{}
				for _, setHeaders := range domainErrorPageInfo.SetHeaders {
					setHeadersMap := map[string]interface{}{}

					if setHeaders.HeaderName != nil {
						setHeadersMap["header_name"] = setHeaders.HeaderName
					}

					if setHeaders.HeaderValue != nil {
						setHeadersMap["header_value"] = setHeaders.HeaderValue
					}

					setHeadersList = append(setHeadersList, setHeadersMap)
				}

				domainErrorPageInfoMap["set_headers"] = setHeadersList
			}

			if domainErrorPageInfo.Body != nil {
				domainErrorPageInfoMap["body"] = domainErrorPageInfo.Body
			}

			if domainErrorPageInfo.Status != nil {
				domainErrorPageInfoMap["status"] = domainErrorPageInfo.Status
			}

			tmpList = append(tmpList, domainErrorPageInfoMap)
		}

		_ = d.Set("error_page_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
