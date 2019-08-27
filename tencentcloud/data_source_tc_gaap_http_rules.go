package tencentcloud

import (
	"context"
	"errors"
	"net"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapHttpRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapHttpRulesRead,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringPrefix("/"),
			},

			// computed
			"rules": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_status_codes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
							Set: func(v interface{}) int {
								return v.(int)
							},
						},
						"realservers": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ip": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"domain": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Computed: true,
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
				return errors.New("rule health check delay loop is nil")
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
				"delay_loop":          *checkParams.DelayLoop,
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

	return nil
}
