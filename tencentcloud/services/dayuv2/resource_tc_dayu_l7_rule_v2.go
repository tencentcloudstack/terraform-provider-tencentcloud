package dayuv2

import (
	"context"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdayu "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dayu"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func ResourceTencentCloudDayuL7RuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL7RuleCreateV2,
		Read:   resourceTencentCloudDayuL7RuleReadV2,
		Update: resourceTencentCloudDayuL7RuleUpdateV2,
		Delete: resourceTencentCloudDayuL7RuleDeleteV2,

		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the resource that the layer 7 rule works for.",
			},
			"resource_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Ip of the resource that the layer 7 rule works for.",
			},
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(svcdayu.DAYU_RESOURCE_TYPE),
				ForceNew:     true,
				Description:  "Type of the resource that the layer 7 rule works for, valid value is `bgpip`.",
			},
			"rule": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Description: "A list of layer 7 rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keeptime": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The keeptime of the layer 4 rule.",
						},
						"domain": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain of the rule.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol of the rule.",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Source type, `1` for source of host, `2` for source of IP.",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "LB type of the rule, `1` for weight cycling and `2` for IP hash.",
						},
						"cert_type": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "The source of the certificate must be filled in when the forwarding protocol is https, the value [2 (Tencent Cloud Hosting Certificate)], and 0 when the forwarding protocol is http.",
						},
						"ssl_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "When the certificate source is a Tencent Cloud managed certificate, this field must be filled in with the managed certificate ID.",
						},
						"source_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Source IP or domain.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Weight of the source.",
									},
								},
							},
						},
						"keep_enable": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "session hold switch.",
						},
						"cc_enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "HTTPS protocol CC protection status, value [0 (off), 1 (on)], defaule is 0.",
						},
						"https_to_http_enable": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Whether to enable the Https protocol to use Http back-to-source, take the value [0 (off), 1 (on)], do not fill in the default is off, defaule is 0.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDayuL7RuleCreateV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	resourceIp := d.Get("resource_ip").(string)
	resourceType := d.Get("resource_type").(string)
	rule := d.Get("rule").([]interface{})
	ruleItem := rule[0].(map[string]interface{})
	domain := ruleItem["domain"].(string)
	protocol := ruleItem["protocol"].(string)
	dayuService := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	err := dayuService.CreateL7RuleV2(ctx, resourceType, resourceId, resourceIp, rule)
	if err != nil {
		return err
	}
	d.SetId(resourceType + tccommon.FILED_SP + domain + tccommon.FILED_SP + protocol)
	return resourceTencentCloudDayuL7RuleReadV2(d, meta)
}

func resourceTencentCloudDayuL7RuleUpdateV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	dayuService := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	business := items[0]
	domain := items[1]
	protocol := items[2]
	extendParams := make(map[string]interface{})
	extendParams["domain"] = domain
	extendParams["protocol"] = protocol
	rules, _, err := dayuService.DescribeL7RulesV2(ctx, business, extendParams)
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		err := fmt.Errorf("Create l7 rule failed.")
		return err
	}
	ruleItem := rules[0]
	resourceId := *ruleItem.Id
	if d.HasChange("rule.0.protocol") {
		protocol = d.Get("protocol").(string)
		ruleItem.Protocol = &protocol
	}

	if d.HasChange("rule.0.source_type") {
		sourceType := uint64(d.Get("source_type").(int))
		ruleItem.SourceType = &sourceType
	}
	if d.HasChange("rule.0.ssl_id") {
		ruleConfig := d.Get("rule").([]interface{})
		ruleConfigItem := ruleConfig[0].(map[string]interface{})
		sslId := ruleConfigItem["ssl_id"].(string)
		ruleItem.SSLId = &sslId
	}
	if d.HasChange("rule.0.cert_type") {
		ruleConfig := d.Get("rule").([]interface{})
		ruleConfigItem := ruleConfig[0].(map[string]interface{})
		certType := uint64(ruleConfigItem["cert_type"].(int))
		ruleItem.CertType = &certType
	}
	if d.HasChange("rule.0.https_to_http_enable") {
		ruleConfig := d.Get("rule").([]interface{})
		ruleConfigItem := ruleConfig[0].(map[string]interface{})
		httpsToHttpEnable := uint64(ruleConfigItem["https_to_http_enable"].(int))
		ruleItem.HttpsToHttpEnable = &httpsToHttpEnable
	}
	if d.HasChange("rule.0.cc_enable") {
		ruleConfig := d.Get("rule").([]interface{})
		ruleConfigItem := ruleConfig[0].(map[string]interface{})
		ccEnable := uint64(ruleConfigItem["cc_enable"].(int))
		ruleItem.CCEnable = &ccEnable
	}
	if d.HasChange("rule.0.source_list") {
		ruleConfig := d.Get("rule").([]interface{})
		ruleConfigItem := ruleConfig[0].(map[string]interface{})
		sourceList := ruleConfigItem["source_list"].([]interface{})
		sources := make([]*dayu.L4RuleSource, 0)
		for _, source := range sourceList {
			sourceItem := source.(map[string]interface{})
			weight := uint64(sourceItem["weight"].(int))
			subSource := sourceItem["source"].(string)
			tmpSource := dayu.L4RuleSource{
				Source: &subSource,
				Weight: &weight,
			}
			sources = append(sources, &tmpSource)
		}
		ruleItem.SourceList = sources
	}

	err = dayuService.ModifyL7RuleV2(ctx, business, resourceId, ruleItem)
	if err != nil {
		return err
	}
	d.SetId(business + tccommon.FILED_SP + domain + tccommon.FILED_SP + protocol)
	return resourceTencentCloudDayuL7RuleReadV2(d, meta)
}

func resourceTencentCloudDayuL7RuleReadV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of dayu L7 rule")
	}
	business := items[0]
	domain := items[1]
	protocol := items[2]
	extendParams := make(map[string]interface{})
	extendParams["domain"] = domain
	extendParams["protocol"] = protocol
	dayuService := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	for {
		rules, _, err := dayuService.DescribeL7RulesV2(ctx, business, extendParams)
		if err != nil {
			return err
		}
		if len(rules) == 0 {
			err := fmt.Errorf("Create l7 rule failed.")
			return err
		}
		if *rules[0].Status == uint64(0) {
			_ = d.Set("resource_id", *rules[0].Id)

			return nil
		} else {
			continue
		}
	}
}

func resourceTencentCloudDayuL7RuleDeleteV2(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dayu_l7_rule.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	items := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(items) < 3 {
		return fmt.Errorf("broken ID of L7 rule")
	}
	business := items[0]
	domain := items[1]
	protocol := items[2]

	dayuService := svcdayu.NewDayuService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	extendParams := make(map[string]interface{})
	extendParams["domain"] = domain
	extendParams["protocol"] = protocol
	rules, _, err := dayuService.DescribeL7RulesV2(ctx, business, extendParams)
	if err != nil {
		return err
	}
	if len(rules) == 0 {
		err := fmt.Errorf("Create l7 rule failed.")
		return err
	}
	ruleItem := rules[0]
	resourceId := *ruleItem.Id
	resourceIp := *ruleItem.Ip
	ruleId := *ruleItem.RuleId
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := dayuService.DeleteL7RulesV2(ctx, business, resourceId, resourceIp, ruleId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
