/*
Use this data source to query forward rule of layer7 listeners.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_layer7_listener" "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id     = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain          = "www.qq.com"
  path            = "/"
  realserver_type = "IP"
  health_check    = true

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }
}

data "tencentcloud_gaap_http_rules" "foo" {
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "${tencentcloud_gaap_http_rule.foo.domain}"
}
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapHttpRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapHttpRulesRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the layer7 listener to be queried.",
			},
			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Forward domain of the layer7 listener to be queried.",
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringPrefix("/"),
				Description:  "Path of the forward rule to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"rules": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "An information list of forward rule of the layer7 listeners. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the forward rule.",
						},
						"listener_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the layer7 listener.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Forward domain of the layer7 listener.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of the forward rule.",
						},
						"realserver_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the realserver.",
						},
						"scheduler": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Scheduling policy of the layer4 listener.",
						},
						"health_check": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether health check is enable.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Interval of the health check.",
						},
						"connect_timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Timeout of the health check response.",
						},
						"health_check_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Path of health check.",
						},
						"health_check_method": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Method of the health check.",
						},
						"health_check_status_codes": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Return code of confirmed normal.",
						},
						"realservers": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An information list of GAAP realserver. Each element contains the following attributes:",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "ID of the GAAP realserver.",
									},
									"ip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "IP of the GAAP realserver.",
									},
									"domain": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Domain of the GAAP realserver.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Port of the GAAP realserver.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Scheduling weight.",
									},
									"status": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Status of the GAAP realserver.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudGaapHttpRulesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_http_rules.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	listenerId := d.Get("listener_id").(string)

	var (
		domain string
		path   *string
		ids    []string
		rules  []map[string]interface{}
	)

	if raw, ok := d.GetOk("domain"); ok {
		domain = raw.(string)
	}
	if raw, ok := d.GetOk("path"); ok {
		path = stringToPointer(raw.(string))
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	domainRuleSets, err := service.DescribeDomains(ctx, listenerId, domain)
	if err != nil {
		return err
	}

	ids = make([]string, 0, len(domainRuleSets))
	rules = make([]map[string]interface{}, 0, len(domainRuleSets))

	for _, domainRule := range domainRuleSets {
		for _, rule := range domainRule.RuleSet {
			if rule.RuleId == nil {
				return errors.New("rule id is nil")
			}
			if rule.Path == nil {
				return errors.New("rule path is nil")
			}
			if rule.RealServerType == nil {
				return errors.New("rule realserver type is nil")
			}
			if rule.Scheduler == nil {
				return errors.New("rule scheduler is nil")
			}
			if rule.HealthCheck == nil {
				return errors.New("rule health check is nil")
			}
			if rule.CheckParams == nil {
				return errors.New("rule health check params is nil")
			}
			checkParams := rule.CheckParams

			if checkParams.DelayLoop == nil {
				return errors.New("rule health check interval is nil")
			}
			if checkParams.ConnectTimeout == nil {
				return errors.New("rule health check connect timeout is nil")
			}
			if checkParams.Path == nil {
				return errors.New("rule health check path is nil")
			}
			if checkParams.Method == nil {
				return errors.New("rule health check method is nil")
			}
			if len(checkParams.StatusCode) == 0 {
				return errors.New("rule health check status codes set is empty")
			}

			if path != nil && *rule.Path != *path {
				continue
			}

			ids = append(ids, *rule.RuleId)

			m := map[string]interface{}{
				"id":                  *rule.RuleId,
				"listener_id":         listenerId,
				"domain":              *domainRule.Domain,
				"path":                *rule.Path,
				"realserver_type":     *rule.RealServerType,
				"scheduler":           *rule.Scheduler,
				"health_check":        *rule.HealthCheck == 1,
				"interval":            *checkParams.DelayLoop,
				"connect_timeout":     *checkParams.ConnectTimeout,
				"health_check_path":   *checkParams.Path,
				"health_check_method": *checkParams.Method,
			}
			statusCodes := make([]int, 0, len(checkParams.StatusCode))
			for _, code := range checkParams.StatusCode {
				statusCodes = append(statusCodes, int(*code))
			}
			m["health_check_status_codes"] = statusCodes

			realservers := make([]map[string]interface{}, 0, len(rule.RealServerSet))
			for _, rs := range rule.RealServerSet {
				if rs.RealServerId == nil {
					return errors.New("realserver id is nil")
				}
				if rs.RealServerIP == nil {
					return errors.New("realserver ip or domain is nil")
				}
				if rs.RealServerPort == nil {
					return errors.New("realserver port is nil")
				}
				if rs.RealServerWeight == nil {
					return errors.New("realserver weight is nil")
				}
				if rs.RealServerStatus == nil {
					return errors.New("realserver status is nil")
				}

				realserver := map[string]interface{}{
					"id":     *rs.RealServerId,
					"port":   *rs.RealServerPort,
					"weight": *rs.RealServerWeight,
					"status": *rs.RealServerStatus,
				}

				if net.ParseIP(*rs.RealServerIP) == nil {
					realserver["domain"] = *rs.RealServerIP
				} else {
					realserver["ip"] = *rs.RealServerIP
				}

				realservers = append(realservers, realserver)
			}

			m["realservers"] = realservers

			rules = append(rules, m)
		}
	}

	d.Set("rules", rules)
	d.SetId(dataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), rules); err != nil {
			log.Printf("[CRITAL]%s output file[%s] fail, reason[%s]\n",
				logId, output.(string), err.Error())
			return err
		}
	}

	return nil
}
