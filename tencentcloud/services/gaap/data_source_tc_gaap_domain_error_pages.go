package gaap

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapDomainErrorPageInfoList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapDomainErrorPageInfoListRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the layer7 listener to be queried.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "HTTP domain to be queried.",
			},
			"ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of the error page info ID to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"error_page_info_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of error page info detail. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the error page info.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the layer7 listener.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "HTTP domain.",
						},
						"error_codes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Original error codes.",
						},
						"body": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "New response body.",
						},
						"new_error_codes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "New error code.",
						},
						"clear_headers": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Response headers to be removed.",
						},
						"set_headers": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Response headers to be set.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapDomainErrorPageInfoListRead(d *schema.ResourceData, m interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_domain_error_pages.read")()
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	var ids []string

	if raw, ok := d.GetOk("ids"); ok {
		ids = helper.InterfacesStrings(raw.(*schema.Set).List())
	}

	service := GaapService{client: m.(tccommon.ProviderMeta).GetAPIV3Conn()}

	respList, err := service.DescribeDomainErrorPageInfoList(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(respList))
	idList := make([]string, 0, len(respList))

	for _, info := range respList {
		if len(ids) > 0 && !helper.StringsContain(ids, *info.ErrorPageId) {
			continue
		}

		idList = append(idList, *info.ErrorPageId)

		m := map[string]interface{}{
			"id":              info.ErrorPageId,
			"listener_id":     info.ListenerId,
			"domain":          info.Domain,
			"error_codes":     info.ErrorNos,
			"body":            info.Body,
			"new_error_codes": info.NewErrorNo,
			"clear_headers":   info.ClearHeaders,
		}

		setHeaders := make(map[string]string, len(info.SetHeaders))
		for _, header := range info.SetHeaders {
			setHeaders[*header.HeaderName] = *header.HeaderValue
		}

		m["set_headers"] = setHeaders

		list = append(list, m)
	}

	_ = d.Set("error_page_info_list", list)

	d.SetId(helper.DataResourceIdsHash(idList))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := tccommon.WriteToFile(output.(string), list); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%v]",
				logId, output.(string), err)
			return err
		}
	}

	return nil
}
