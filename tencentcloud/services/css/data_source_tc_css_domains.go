package css

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	css "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/live/v20180801"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCssDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCssDomainsRead,
		Schema: map[string]*schema.Schema{
			"domain_status": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "domain name status filter. 0-disable, 1-enable.",
			},

			"domain_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain name type filtering. 0-push, 1-play.",
			},

			"is_delay_live": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "0 normal live broadcast 1 slow live broadcast default 0.",
			},

			"domain_prefix": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "domain name prefix.",
			},

			"play_type": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Playing area, this parameter is meaningful only when DomainType=1. 1: Domestic.2: Global.3: Overseas.",
			},

			"domain_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "A list of domain name details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Live domain name.",
						},
						"type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain Type: 0: push stream. 1: Play.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain Status: 0: disable. 1: Enabled.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "add time.Note: This field is Beijing time (UTC+8 time zone).",
						},
						"b_c_name": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Is there a CName to the fixed rule domain name: 0: No. 1: Yes.",
						},
						"target_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name corresponding to the cname.",
						},
						"play_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Playing area, this parameter is meaningful only when Type=1. 1: Domestic. 2: Global. 3: Overseas.",
						},
						"is_delay_live": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to slow live broadcast: 0: normal live broadcast. 1: Slow live broadcast.",
						},
						"current_c_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cname information used by the current client.",
						},
						"rent_tag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "invalid parameter, can be ignored.",
						},
						"rent_expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Failure parameter, can be ignored. Note: This field is Beijing time (UTC+8 time zone).",
						},
						"is_mini_program_live": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "0: Standard live broadcast. 1: Mini program live broadcast. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudCssDomainsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_css_domains.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("domain_status"); ok {
		paramMap["DomainStatus"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("domain_type"); ok {
		paramMap["DomainType"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("is_delay_live"); ok {
		paramMap["IsDelayLive"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("domain_prefix"); ok {
		paramMap["DomainPrefix"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("play_type"); ok {
		paramMap["PlayType"] = helper.IntUint64(v.(int))
	}

	service := CssService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var domainList []*css.DomainInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCssDomainsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		domainList = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(domainList))
	tmpList := make([]map[string]interface{}, 0, len(domainList))

	if domainList != nil {
		for _, domainInfo := range domainList {
			domainInfoMap := map[string]interface{}{}

			if domainInfo.Name != nil {
				domainInfoMap["name"] = domainInfo.Name
			}

			if domainInfo.Type != nil {
				domainInfoMap["type"] = domainInfo.Type
			}

			if domainInfo.Status != nil {
				domainInfoMap["status"] = domainInfo.Status
			}

			if domainInfo.CreateTime != nil {
				domainInfoMap["create_time"] = domainInfo.CreateTime
			}

			if domainInfo.BCName != nil {
				domainInfoMap["b_c_name"] = domainInfo.BCName
			}

			if domainInfo.TargetDomain != nil {
				domainInfoMap["target_domain"] = domainInfo.TargetDomain
			}

			if domainInfo.PlayType != nil {
				domainInfoMap["play_type"] = domainInfo.PlayType
			}

			if domainInfo.IsDelayLive != nil {
				domainInfoMap["is_delay_live"] = domainInfo.IsDelayLive
			}

			if domainInfo.CurrentCName != nil {
				domainInfoMap["current_c_name"] = domainInfo.CurrentCName
			}

			if domainInfo.RentTag != nil {
				domainInfoMap["rent_tag"] = domainInfo.RentTag
			}

			if domainInfo.RentExpireTime != nil {
				domainInfoMap["rent_expire_time"] = domainInfo.RentExpireTime
			}

			if domainInfo.IsMiniProgramLive != nil {
				domainInfoMap["is_mini_program_live"] = domainInfo.IsMiniProgramLive
			}

			ids = append(ids, *domainInfo.Name)
			tmpList = append(tmpList, domainInfoMap)
		}

		_ = d.Set("domain_list", tmpList)
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
