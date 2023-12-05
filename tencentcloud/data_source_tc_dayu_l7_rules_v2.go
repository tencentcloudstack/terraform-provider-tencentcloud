package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudDayuL7RulesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuL7RulesReadV2,
		Schema: map[string]*schema.Schema{
			"business": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the layer 4 rule works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain of resource.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Protocol of resource, value range [`http`, `https`].",
			},
			"ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Ip of the resource.",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "The page start offset, default is `0`.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Deprecated:  "It has been deprecated from version 1.81.21.",
				Description: "The number of pages, default is `10`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of layer 4 rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keep_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session hold time, in seconds.",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Load balancing mode, the value is [1 (weighted round-robin)].",
						},
						"source_list": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Back-to-source IP or domain name.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Weight value, take value [0,100].",
									},
								},
							},
							Description: "Source list of the rule.",
						},
						"keep_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Session keep switch, value [0 (session keep closed), 1 (session keep open)].",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain of resource.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol of resource, value range [`http`, `https`].",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Back-to-origin method, value [1 (domain name back-to-source), 2 (IP back-to-source)].",
						},
						"https_to_http_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Whether to enable the Https protocol to use Http back-to-source, take the value [0 (off), 1 (on)], default is off.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Rule status, value [0 (rule configuration is successful), 1 (rule configuration is in effect), 2 (rule configuration fails), 3 (rule deletion is in effect), 5 (rule deletion fails), 6 (rule is waiting to be configured), 7 (rule pending deletion), 8 (rule pending configuration certificate)].",
						},
						"cc_level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "CC protection level of HTTPS protocol.",
						},
						"cc_enable": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection status of HTTPS protocol, the value is [0 (off), 1 (on)].",
						},
						"cc_threshold": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection threshold of HTTPS protocol.",
						},
						"region": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The area code.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Rule description.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modify time of resource.",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Virtual port of resource.",
						},
						"cc_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "CC protection status, value [0(off), 1(on)].",
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Ip of the resource.",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "SSL id of the resource.",
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The source of the certificate.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the resource.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuL7RulesReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_dayu_l4_rules_v2.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	business := d.Get("business").(string)
	domain := d.Get("domain").(string)
	protocol := d.Get("protocol").(string)
	ip := d.Get("ip").(string)

	extendParams := make(map[string]interface{})
	extendParams["domain"] = domain
	extendParams["protocol"] = protocol
	extendParams["ip"] = ip

	rules, _, err := service.DescribeL7RulesV2(ctx, business, extendParams)
	if err != nil {
		return err
	}

	resultList := make([]map[string]interface{}, 0)
	for _, rule := range rules {
		tmpResultItem := make(map[string]interface{})
		tmpResultItem["keep_time"] = *rule.KeepTime
		tmpResultItem["lb_type"] = *rule.LbType
		sourceList := make([]map[string]interface{}, 0)
		for _, source := range rule.SourceList {
			tmpSource := make(map[string]interface{})
			tmpSource["source"] = *source.Source
			tmpSource["weight"] = *source.Weight
			sourceList = append(sourceList, tmpSource)
		}
		tmpResultItem["source_list"] = sourceList
		tmpResultItem["keep_enable"] = *rule.KeepEnable
		tmpResultItem["domain"] = *rule.Domain
		tmpResultItem["protocol"] = *rule.Protocol
		tmpResultItem["source_type"] = *rule.SourceType
		tmpResultItem["https_to_http_enable"] = *rule.HttpsToHttpEnable
		tmpResultItem["status"] = *rule.Status
		tmpResultItem["cc_level"] = *rule.CCLevel
		tmpResultItem["cc_enable"] = *rule.CCEnable
		tmpResultItem["cc_threshold"] = *rule.CCThreshold
		tmpResultItem["region"] = *rule.Region
		tmpResultItem["rule_name"] = *rule.RuleName
		tmpResultItem["modify_time"] = *rule.ModifyTime
		tmpResultItem["virtual_port"] = *rule.VirtualPort
		tmpResultItem["cc_status"] = *rule.CCStatus
		tmpResultItem["ip"] = *rule.Ip
		tmpResultItem["cert_type"] = *rule.CertType
		tmpResultItem["id"] = *rule.Id
		resultList = append(resultList, tmpResultItem)
	}
	ids := make([]string, 0, len(resultList))
	for _, listItem := range resultList {
		ids = append(ids, listItem["id"].(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", resultList); e != nil {
		log.Printf("[CRITAL]%s provider set rules fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), resultList)
	}
	return nil
}
