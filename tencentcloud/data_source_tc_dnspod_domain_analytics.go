/*
Use this data source to query detailed information of dnspod domain_analytics

Example Usage

```hcl

	data "tencentcloud_dnspod_domain_analytics" "domain_analytics" {
	  domain = "dnspod.cn"
	  start_date = "2023-10-07"
	  end_date = "2023-10-12"
	  dns_format = "HOUR"
	  # domain_id = 123
	}

```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDnspodDomainAnalytics() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodDomainAnalyticsRead,
		Schema: map[string]*schema.Schema{
			"domain": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The domain name to query for resolution volume.",
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

			"dns_format": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "DATE: Statistics by day dimension HOUR: Statistics by hour dimension.",
			},

			"domain_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Domain ID. The parameter DomainId has a higher priority than the parameter Domain. If the parameter DomainId is passed, the parameter Domain will be ignored. You can find all Domains and DomainIds through the DescribeDomainList interface.",
			},

			"data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Subtotal of resolution volume for the current statistical dimension.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Subtotal of resolution volume for the current statistical dimension.",
						},
						"date_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "For daily statistics, it is the statistical date.",
						},
						"hour_key": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "For hourly statistics, it is the hour of the current time (0-23), for example, when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.",
						},
					},
				},
			},

			"info": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain resolution volume statistics query information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_format": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "DATE: Statistics by day dimension HOUR: Statistics by hour dimension.",
						},
						"dns_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total resolution volume for the current statistical period.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The domain name currently being queried.",
						},
						"start_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start time of the current statistical period.",
						},
						"end_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End time of the current statistical period.",
						},
					},
				},
			},

			"alias_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Domain alias resolution volume statistics information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"info": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Domain resolution volume statistics query information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"dns_format": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "DATE: Statistics by day dimension HOUR: Statistics by hour dimension.",
									},
									"dns_total": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Total resolution volume for the current statistical period.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The domain name currently being queried.",
									},
									"start_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Start time of the current statistical period.",
									},
									"end_date": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "End time of the current statistical period.",
									},
								},
							},
						},
						"data": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Subtotal of resolution volume for the current statistical dimension.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"num": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Subtotal of resolution volume for the current statistical dimension.",
									},
									"date_key": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "For daily statistics, it is the statistical date.",
									},
									"hour_key": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "For hourly statistics, it is the hour of the current time (0-23), for example, when HourKey is 23, the statistical period is the resolution volume from 22:00 to 23:00. Note: This field may return null, indicating that no valid value can be obtained.",
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

func dataSourceTencentCloudDnspodDomainAnalyticsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dnspod_domain_analytics.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		domain    string
		aliasData []*dnspod.DomainAliasAnalyticsItem
		data      []*dnspod.DomainAnalyticsDetail
		info      *dnspod.DomainAnalyticsInfo
		err       error
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

	if v, ok := d.GetOk("dns_format"); ok {
		paramMap["DnsFormat"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("domain_id"); ok {
		paramMap["DomainId"] = helper.IntUint64(v.(int))
	}

	service := DnspodService{client: meta.(*TencentCloudClient).apiV3Conn}
	e := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		aliasData, data, info, err = service.DescribeDnspodDomainAnalyticsByFilter(ctx, paramMap)
		if err != nil {
			return retryError(err)
		}
		return nil
	})
	if e != nil {
		return e
	}

	// ids := make([]string, 0, len(data))
	tmpDataList := make([]map[string]interface{}, 0, len(data))
	tmpAliasDataList := make([]map[string]interface{}, 0, len(aliasData))

	if data != nil {
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
			tmpDataList = append(tmpDataList, domainAnalyticsDetailMap)
		}

		_ = d.Set("data", tmpDataList)
	}

	if info != nil {
		domainAnalyticsInfoMap := map[string]interface{}{}

		if info.DnsFormat != nil {
			domainAnalyticsInfoMap["dns_format"] = info.DnsFormat
		}

		if info.DnsTotal != nil {
			domainAnalyticsInfoMap["dns_total"] = info.DnsTotal
		}

		if info.Domain != nil {
			domainAnalyticsInfoMap["domain"] = info.Domain
		}

		if info.StartDate != nil {
			domainAnalyticsInfoMap["start_date"] = info.StartDate
		}

		if info.EndDate != nil {
			domainAnalyticsInfoMap["end_date"] = info.EndDate
		}

		// ids = append(ids, *info.Domain)
		// _ = d.Set("info", domainAnalyticsInfoMap)
		e = helper.SetMapInterfaces(d, "info", domainAnalyticsInfoMap)
		if e != nil {
			return e
		}
	}

	if aliasData != nil {
		for _, domainAliasAnalyticsItem := range aliasData {
			domainAliasAnalyticsItemMap := map[string]interface{}{}

			if domainAliasAnalyticsItem.Info != nil {
				infoMap := map[string]interface{}{}

				if domainAliasAnalyticsItem.Info.DnsFormat != nil {
					infoMap["dns_format"] = domainAliasAnalyticsItem.Info.DnsFormat
				}

				if domainAliasAnalyticsItem.Info.DnsTotal != nil {
					infoMap["dns_total"] = domainAliasAnalyticsItem.Info.DnsTotal
				}

				if domainAliasAnalyticsItem.Info.Domain != nil {
					infoMap["domain"] = domainAliasAnalyticsItem.Info.Domain
				}

				if domainAliasAnalyticsItem.Info.StartDate != nil {
					infoMap["start_date"] = domainAliasAnalyticsItem.Info.StartDate
				}

				if domainAliasAnalyticsItem.Info.EndDate != nil {
					infoMap["end_date"] = domainAliasAnalyticsItem.Info.EndDate
				}

				domainAliasAnalyticsItemMap["info"] = []interface{}{infoMap}
			}

			if domainAliasAnalyticsItem.Data != nil {
				dataList := []interface{}{}
				for _, data := range domainAliasAnalyticsItem.Data {
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

				domainAliasAnalyticsItemMap["data"] = []interface{}{dataList}
			}

			// ids = append(ids, *domainAliasAnalyticsItem.Domain)
			tmpAliasDataList = append(tmpAliasDataList, domainAliasAnalyticsItemMap)
		}

		_ = d.Set("alias_data", tmpAliasDataList)
	}

	d.SetId(helper.DataResourceIdHash(domain))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		e = writeToFile(output.(string), map[string]interface{}{
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
