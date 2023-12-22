package dnspod

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodRecordAnalytics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordAnalyticsRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The domain to query for resolution volume.",
			},

			"start_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The start date of the query, format: YYYY-MM-DD.",
			},

			"end_date": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The end date of the query, format: YYYY-MM-DD.",
			},

			"subdomain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The subdomain to query for resolution volume.",
			},

			"dns_format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DATE: Statistics by day dimension, HOUR: Statistics by hour dimension.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "The subtotal of the resolution volume for the current statistical dimension.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The subtotal of the resolution volume for the current statistical dimension.",
						},
						"date_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "For daily statistics, it is the statistical date.",
						},
						"hour_key": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For hourly statistics, it is the hour of the current time for statistics (0-23), e.g., when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Subdomain resolution statistics query information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DATE: Daily statistics, HOUR: Hourly statistics.",
						},
						"dns_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total resolution count for the current statistical period.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain currently being queried.",
						},
						"start_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start date of the current statistical period.",
						},
						"end_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End date of the current statistical period.",
						},
						"subdomain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The subdomain currently being analyzed.",
						},
					},
				},
			},

			"alias_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Subdomain alias resolution statistics information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subdomain resolution statistics query information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DATE: Daily statistics, HOUR: Hourly statistics.",
									},
									"dns_total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total resolution count for the current statistical period.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The domain currently being queried.",
									},
									"start_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start date of the current statistical period.",
									},
									"end_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End date of the current statistical period.",
									},
									"subdomain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The subdomain currently being analyzed.",
									},
								},
							},
						},
						"data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subtotal of the resolution volume for the current statistical dimension.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The subtotal of the resolution volume for the current statistical dimension.",
									},
									"date_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "For daily statistics, it is the statistical date.",
									},
									"hour_key": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "For hourly statistics, it is the hour of the current time for statistics (0-23), e.g., when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudDnspodRecordAnalyticsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_record_analytics.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		domain    string
		aliasData []*dnspod.SubdomainAliasAnalyticsItem
		data      []*dnspod.DomainAnalyticsDetail
		info      *dnspod.SubdomainAnalyticsInfo
		e         error
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain"); ok {
		domain = v.(string)
		paramMap["Domain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("start_date"); ok {
		paramMap["StartDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("end_date"); ok {
		paramMap["EndDate"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("subdomain"); ok {
		paramMap["Subdomain"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dns_format"); ok {
		paramMap["DnsFormat"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	// var data []*dnspod.DomainAnalyticsDetail

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		aliasData, data, info, e = service.DescribeDnspodRecordAnalyticsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		// data = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(data))
	if data != nil {
		tmpList := make([]map[string]interface{}, 0, len(data))
		for _, domainAnalyticsDetail := range data {
			domainAnalyticsDetailMap := map[string]interface{}{}

			if domainAnalyticsDetail.Num != nil {
				domainAnalyticsDetailMap["num"] = domainAnalyticsDetail.Num
			}

			if domainAnalyticsDetail.DateKey != nil {
				domainAnalyticsDetailMap["date_key"] = domainAnalyticsDetail.DateKey
			}

			if domainAnalyticsDetail.HourKey != nil {
				domainAnalyticsDetailMap["hour_key"] = domainAnalyticsDetail.HourKey
			}

			// ids = append(ids, *domainAnalyticsDetail.Domain)
			tmpList = append(tmpList, domainAnalyticsDetailMap)
		}

		_ = d.Set("data", tmpList)
	}

	if info != nil {
		subdomainAnalyticsInfoMap := map[string]interface{}{}

		if info.DnsFormat != nil {
			subdomainAnalyticsInfoMap["dns_format"] = info.DnsFormat
		}

		if info.DnsTotal != nil {
			subdomainAnalyticsInfoMap["dns_total"] = info.DnsTotal
		}

		if info.Domain != nil {
			subdomainAnalyticsInfoMap["domain"] = info.Domain
		}

		if info.StartDate != nil {
			subdomainAnalyticsInfoMap["start_date"] = info.StartDate
		}

		if info.EndDate != nil {
			subdomainAnalyticsInfoMap["end_date"] = info.EndDate
		}

		if info.Subdomain != nil {
			subdomainAnalyticsInfoMap["subdomain"] = info.Subdomain
		}

		// ids = append(ids, *info.Domain)
		e = helper.SetMapInterfaces(d, "info", subdomainAnalyticsInfoMap)
		if e != nil {
			return e
		}
	}

	if aliasData != nil {
		tmpList := make([]map[string]interface{}, 0, len(aliasData))
		for _, subdomainAliasAnalyticsItem := range aliasData {
			subdomainAliasAnalyticsItemMap := map[string]interface{}{}

			if subdomainAliasAnalyticsItem.Info != nil {
				infoMap := map[string]interface{}{}

				if subdomainAliasAnalyticsItem.Info.DnsFormat != nil {
					infoMap["dns_format"] = subdomainAliasAnalyticsItem.Info.DnsFormat
				}

				if subdomainAliasAnalyticsItem.Info.DnsTotal != nil {
					infoMap["dns_total"] = subdomainAliasAnalyticsItem.Info.DnsTotal
				}

				if subdomainAliasAnalyticsItem.Info.Domain != nil {
					infoMap["domain"] = subdomainAliasAnalyticsItem.Info.Domain
				}

				if subdomainAliasAnalyticsItem.Info.StartDate != nil {
					infoMap["start_date"] = subdomainAliasAnalyticsItem.Info.StartDate
				}

				if subdomainAliasAnalyticsItem.Info.EndDate != nil {
					infoMap["end_date"] = subdomainAliasAnalyticsItem.Info.EndDate
				}

				if subdomainAliasAnalyticsItem.Info.Subdomain != nil {
					infoMap["subdomain"] = subdomainAliasAnalyticsItem.Info.Subdomain
				}

				subdomainAliasAnalyticsItemMap["info"] = []interface{}{infoMap}
			}

			if subdomainAliasAnalyticsItem.Data != nil {
				dataList := []interface{}{}
				for _, data := range subdomainAliasAnalyticsItem.Data {
					dataMap := map[string]interface{}{}

					if data.Num != nil {
						dataMap["num"] = data.Num
					}

					if data.DateKey != nil {
						dataMap["date_key"] = data.DateKey
					}

					if data.HourKey != nil {
						dataMap["hour_key"] = data.HourKey
					}

					dataList = append(dataList, dataMap)
				}

				subdomainAliasAnalyticsItemMap["data"] = []interface{}{dataList}
			}

			// ids = append(ids, *subdomainAliasAnalyticsItem.Domain)
			tmpList = append(tmpList, subdomainAliasAnalyticsItemMap)
		}

		_ = d.Set("alias_data", tmpList)
	}

	// d.SetId(helper.DataResourceIdsHash(ids))
	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		e = tccommon.WriteToFile(output.(string), map[string]interface{}{
			"info":       info,
			"data":       data,
			"alias_data": aliasData,
		})
		if e != nil {
			return e
		}
	}
	return nil
}
