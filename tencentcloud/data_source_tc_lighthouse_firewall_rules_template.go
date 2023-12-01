/*
Use this data source to query detailed information of lighthouse firewall_rules_template

Example Usage

```hcl
data "tencentcloud_lighthouse_firewall_rules_template" "firewall_rules_template" {
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudLighthouseFirewallRulesTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudLighthouseFirewallRulesTemplateRead,
		Schema: map[string]*schema.Schema{
			"firewall_rule_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Firewall rule details list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Application type. Valid values are custom, HTTP (80), HTTPS (443), Linux login (22), Windows login (3389), MySQL (3306), SQL Server (1433), all TCP ports, all UDP ports, Ping-ICMP, ALL.",
						},
						"protocol": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Protocol. Valid values are TCP, UDP, ICMP, ALL.",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Port. Valid values are ALL, one single port, multiple ports separated by commas, or port range indicated by a minus sign.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "IP range or IP (mutually exclusive). Default value is 0.0.0.0/0, which indicates all sources.",
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Valid values are (ACCEPT, DROP). Default value is ACCEPT.",
						},
						"firewall_rule_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall rule description.",
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

func dataSourceTencentCloudLighthouseFirewallRulesTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_lighthouse_firewall_rules_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LightHouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var firewallRuleSet []*lighthouse.FirewallRuleInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeLighthouseFirewallRulesTemplateByFilter(ctx)
		if e != nil {
			return retryError(e)
		}
		firewallRuleSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(firewallRuleSet))
	tmpList := make([]map[string]interface{}, 0, len(firewallRuleSet))

	if firewallRuleSet != nil {
		for _, firewallRuleInfo := range firewallRuleSet {
			firewallRuleInfoMap := map[string]interface{}{}

			if firewallRuleInfo.AppType != nil {
				firewallRuleInfoMap["app_type"] = firewallRuleInfo.AppType
			}

			if firewallRuleInfo.Protocol != nil {
				firewallRuleInfoMap["protocol"] = firewallRuleInfo.Protocol
			}

			if firewallRuleInfo.Port != nil {
				firewallRuleInfoMap["port"] = firewallRuleInfo.Port
			}

			if firewallRuleInfo.CidrBlock != nil {
				firewallRuleInfoMap["cidr_block"] = firewallRuleInfo.CidrBlock
			}

			if firewallRuleInfo.Action != nil {
				firewallRuleInfoMap["action"] = firewallRuleInfo.Action
			}

			if firewallRuleInfo.FirewallRuleDescription != nil {
				firewallRuleInfoMap["firewall_rule_description"] = firewallRuleInfo.FirewallRuleDescription
			}
			firewallRuleInfoJson, err := json.Marshal(*firewallRuleInfo)
			if err != nil {
				return err
			}
			ids = append(ids, string(firewallRuleInfoJson))
			tmpList = append(tmpList, firewallRuleInfoMap)
		}

		_ = d.Set("firewall_rule_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
