/*
Use this data source to query detailed information of CLB listener rule

Example Usage

```hcl
data "tencentcloud_clb_listener_rules" "foo" {
  clb_id      = "lb-k2zjp9lv"
  listener_id = "lbl-mwr6vbtv"
  rule_id     = "loc-inem40hz"
  domain      = "abc.com"
  url         = "/"
  scheduler   = "WRR"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudClbListenerRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClbListenerRulesRead,

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the CLB to be queried.",
			},
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the CLB listener to be queried.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the forwarding rule to be queried.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Domain name of the forwarding rule to be queried.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Url of the forwarding rule to be queried.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method of the forwarding rule of thr CLB listener, and available values include 'WRR' , 'IP HASH' and 'LEAST_CONN'. The defaule is 'WRR'.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"rule_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of forward rules of listeners. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"clb_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the CLB.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the listener.",
						},
						"domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Domain of the rule.",
						},
						"url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Url of the rule.",
						},
						"rule_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the rule.",
						},
						"health_check_switch": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enabled.",
						},
						"health_check_interval_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval time of health check. The value range is 5-300 sec, and the default is 5 sec.",
						},
						"health_check_health_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Health threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as healthy. The value range is 2-10.",
						},
						"health_check_unhealth_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Unhealth threshold of health check, and the default is 3. If a success result is returned for the health check three consecutive times, the CVM is identified as unhealthy. The value range is 2-10.",
						},
						"health_check_http_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "HTTP Status Code. The default is 31 and value range is 1-31. '0b0001' means the return value '1xx' is health. '0b0010' means the return value '2xx' is health. '0b0100' means the return value '3xx' is health. '0b1000' means the return value 4xx is health. '0b10000' means the return value '5xx' is health. If you want multiple return codes to indicate health, need to add the corresponding values. NOTES: The 'HTTP' health check of the 'TCP' listener only supports specifying one health check status code. NOTES: Only supports listeners of 'HTTP' and 'HTTPS' protocol.",
						},
						"health_check_http_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol.",
						},
						"health_check_http_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain name of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol.",
						},
						"health_check_http_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Methods of health check. NOTES: Only supports listeners of 'HTTPS' and 'HTTP' protocol. The default is 'HEAD', the available value include 'HEAD' and 'GET'.",
						},
						"certificate_ssl_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of SSL Mode, and available values inclue 'UNIDIRECTIONAL', 'MUTUAL'.",
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the server certificate. NOTES: Only supports listeners of 'HTTPS' protocol.",
						},
						"certificate_ca_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the client certificate. NOTES: Only supports listeners of 'HTTPS' protocol.",
						},
						"session_expire_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Time of session persistence within the CLB listener. NOTES: Available when scheduler is specified as 'WRR'.",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling method of the CLB listener, and available values include 'WRR' and 'LEAST_CONN'. The defaule is 'WRR'.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudClbListenerRulesRead(d *schema.ResourceData, meta interface{}) error {
	defer LogElapsed("data_source.tencentcloud_clb_listener_rules.read")()

	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	combinedId := d.Get("listener_id").(string)
	items := strings.Split(combinedId, "#")
	if len(items) != 2 {
		return fmt.Errorf("id of resource.tencentcloud_clb_listener_rules is wrong")
	}

	listenerId := items[0]
	clbId := items[1]
	params := make(map[string]string)
	params["clb_id"] = clbId
	params["listener_id"] = listenerId
	if v, ok := d.GetOk("clb_id"); ok {
		params["clb_id"] = v.(string)
	}
	if v, ok := d.GetOk("scheduler"); ok {
		params["scheduler"] = v.(string)
	}
	if v, ok := d.GetOk("rule_id"); ok {
		params["rule_id"] = v.(string)
	}
	if v, ok := d.GetOk("domain"); ok {
		params["domain"] = v.(string)
	}
	if v, ok := d.GetOk("url"); ok {
		params["url"] = v.(string)
	}

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	rules, err := clbService.DescribeRulesByFilter(ctx, params)
	if err != nil {
		return err
	}

	ruleList := make([]map[string]interface{}, 0, len(rules))
	log.Printf("the length %d", len(rules))
	ids := make([]string, 0, len(rules))
	for _, rule := range rules {
		mapping := map[string]interface{}{
			"clb_id":              clbId,
			"listener_id":         combinedId,
			"rule_id":             *rule.LocationId,
			"domain":              *rule.Domain,
			"url":                 *rule.Url,
			"session_expire_time": *rule.SessionExpireTime,
			"scheduler":           *rule.Scheduler,
		}
		if rule.HealthCheck != nil {
			health_check_switch := false
			if *rule.HealthCheck.HealthSwitch == int64(1) {
				health_check_switch = true
			}
			mapping["health_check_switch"] = health_check_switch
			mapping["health_check_interval_time"] = *rule.HealthCheck.IntervalTime
			mapping["health_check_health_num"] = *rule.HealthCheck.HealthNum
			mapping["health_check_unhealth_num"] = *rule.HealthCheck.UnHealthNum
			mapping["health_check_http_code"] = *rule.HealthCheck.HttpCode
			mapping["health_check_http_method"] = *rule.HealthCheck.HttpCheckMethod
			mapping["health_check_http_domain"] = *rule.HealthCheck.HttpCheckDomain
			mapping["health_check_http_path"] = *rule.HealthCheck.HttpCheckPath
		}
		if rule.Certificate != nil {
			mapping["certificate_ssl_mode"] = *rule.Certificate.SSLMode
			mapping["certificate_id"] = *rule.Certificate.CertId
			mapping["certificate_ca_id"] = *rule.Certificate.CertCaId
		}
		ruleList = append(ruleList, mapping)
		ids = append(ids, *rule.LocationId+"#"+combinedId)
	}

	d.SetId(dataResourceIdsHash(ids))
	if err = d.Set("rule_list", ruleList); err != nil {
		log.Printf("[CRITAL]%s provider set clb listener rule list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), ruleList); err != nil {
			return err
		}
	}

	return nil
}
