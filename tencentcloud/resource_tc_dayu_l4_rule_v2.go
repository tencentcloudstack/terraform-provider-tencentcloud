/*
Use this resource to create dayu layer 4 rule

~> **NOTE:** This resource only support resource Anti-DDoS of type `bgpip` and `net`

Example Usage

```hcl
resource "tencentcloud_dayu_l4_rule" "test_rule" {
  resource_type             = "bgpip"
  resource_id               = "bgpip-00000294"
  name                      = "rule_test"
  protocol                  = "TCP"
  s_port                    = 80
  d_port                    = 60
  source_type               = 2
  health_check_switch       = true
  health_check_timeout      = 30
  health_check_interval     = 35
  health_check_health_num   = 5
  health_check_unhealth_num = 10
  session_switch            = false
  session_time              = 300

  source_list {
    source = "1.1.1.1"
    weight = 100
  }
  source_list {
    source = "2.2.2.2"
    weight = 50
  }
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"
)

func resourceTencentCloudDayuL4RuleV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDayuL4RuleCreateV2,
		Read:   resourceTencentCloudDayuL4RuleReadV2,
		Delete: resourceTencentCloudDayuL4RuleDeleteV2,

		Schema: map[string]*schema.Schema{
			"business": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue(DAYU_RESOURCE_TYPE_RULE),
				ForceNew:     true,
				Description:  "Business of the resource that the layer 4 rule works for. Valid values: `bgpip` and `net`.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource id.",
			},
			"vpn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Resource vpn.",
			},
			"virtual_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The virtual port of the layer 4 rule.",
			},
			"rules": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "A list of layer 4 rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Protocol of the rule.",
						},
						"source_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The source port of the layer 4 rule.",
						},
						"virtual_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The virtual port of the layer 4 rule.",
						},
						"keeptime": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The keeptime of the layer 4 rule.",
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
							Description: "Source list of the rule.",
						},
						"lb_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "LB type of the rule, `1` for weight cycling and `2` for IP hash.",
						},
						"keep_enable": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "session hold switch.",
						},
						"source_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Source type, `1` for source of host, `2` for source of IP.",
						},
						"rule_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the rule.",
						},
						"remove_switch": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "Remove the watermark state.",
						},
						"region": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Corresponding regional information.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudDayuL4RuleCreateV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule_v2.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	business := d.Get("business").(string)
	vpn := d.Get("vpn").(string)
	virtualPort := d.Get("virtual_port").(int)
	rules := d.Get("rules").([]interface{})
	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := service.CreateNewL4Rules(ctx, business, resourceId, vpn, rules)
	if err != nil {
		return err
	}
	virtualPortStr := strconv.Itoa(virtualPort)
	d.SetId(business + FILED_SP + resourceId + FILED_SP + vpn + FILED_SP + virtualPortStr)
	return resourceTencentCloudDayuL4RuleReadV2(d, meta)
}

func resourceTencentCloudDayuL4RuleReadV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 4 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	business := items[0]
	resourceId := items[1]
	ip := items[2]
	virtualPortStr := items[3]

	extendParams := make(map[string]interface{})
	extendParams["ip"] = ip
	virtualPort, err := strconv.Atoi(virtualPortStr)
	if err != nil {
		log.Printf("virtual_port must be int.")
	}
	extendParams["virtual_port"] = virtualPort

	rules := make([]*dayu.NewL4RuleEntry, 0)
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)

		if err != nil {
			return retryError(err)
		}
		rules = result
		return nil
	})

	if err != nil {
		return err
	}
	posRules := make([]dayu.NewL4RuleEntry, 0)
	for _, rule := range rules {
		if *rule.Id == resourceId {
			posRules = append(posRules, *rule)
		}
	}

	_ = d.Set("rules", posRules)

	return nil
}

func resourceTencentCloudDayuL4RuleDeleteV2(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_dayu_l4_rule.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := DayuService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	items := strings.Split(d.Id(), FILED_SP)
	if len(items) < 4 {
		return fmt.Errorf("broken ID of dayu L4 rule")
	}
	business := items[0]
	resourceId := items[1]
	vpn := items[2]
	virtualPortStr := items[3]
	virtualPort, err := strconv.Atoi(virtualPortStr)
	if err != nil {
		log.Printf("virtual_port must be int.")
	}
	extendParams := make(map[string]interface{})
	extendParams["ip"] = vpn
	extendParams["virtual_port"] = virtualPort

	rules := make([]*dayu.NewL4RuleEntry, 0)
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, err := service.DescribeNewL4Rules(ctx, business, extendParams)

		if err != nil {
			return retryError(err)
		}
		rules = result
		return nil
	})

	if err != nil {
		return err
	}
	var delRuleId string
	if len(rules) > 0 {
		if *rules[0].Ip == vpn && *rules[0].VirtualPort == uint64(virtualPort) {
			delRuleId = *rules[0].RuleId
		}
	}
	err = service.DeleteNewL4Rules(ctx, business, resourceId, vpn, []string{delRuleId})
	if err != nil {
		return err
	}
	return nil
}
