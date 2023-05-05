/*
Provides a resource to create a forward rule of layer7 listener.

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
  proxy_id = tencentcloud_gaap_proxy.foo.id
}

resource "tencentcloud_gaap_realserver" "foo" {
  ip   = "1.1.1.1"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_realserver" "bar" {
  ip   = "8.8.8.8"
  name = "ci-test-gaap-realserver"
}

resource "tencentcloud_gaap_http_domain" "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id               = tencentcloud_gaap_layer7_listener.foo.id
  domain                    = tencentcloud_gaap_http_domain.foo.domain
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = tencentcloud_gaap_realserver.foo.id
    ip   = tencentcloud_gaap_realserver.foo.ip
    port = 80
  }

  realservers {
    id   = tencentcloud_gaap_realserver.bar.id
    ip   = tencentcloud_gaap_realserver.bar.ip
    port = 80
  }
}
```

Import

GAAP http rule can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_http_rule.foo rule-3bsuu01r
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapHttpRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapHttpRuleCreate,
		Read:   resourceTencentCloudGaapHttpRuleRead,
		Update: resourceTencentCloudGaapHttpRuleUpdate,
		Delete: resourceTencentCloudGaapHttpRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the layer7 listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Forward domain of the forward rule.",
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					_, errs = validateStringLengthInRange(1, 80)(v, k)
					if len(errs) > 0 {
						return
					}

					return validateStringPrefix("/")(v, k)
				},
				Description: "Path of the forward rule. Maximum length is 80.",
			},
			"realserver_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"IP", "DOMAIN"}),
				ForceNew:     true,
				Description:  "Type of the realserver. Valid value: `IP` and `DOMAIN`.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rr",
				ValidateFunc: validateAllowedStringValue([]string{"rr", "wrr", "lc"}),
				Description:  "Scheduling policy of the forward rule, default value is `rr`. Valid value: `rr`, `wrr` and `lc`.",
			},
			"health_check": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Indicates whether health check is enable.",
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      5,
				ValidateFunc: validateIntegerInRange(5, 300),
				Description:  "Interval of the health check, default value is 5s.",
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validateIntegerInRange(2, 60),
				Description:  "Timeout of the health check response, default value is 2s.",
			},
			"health_check_path": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/",
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					_, errs = validateStringLengthInRange(1, 80)(v, k)
					if len(errs) > 0 {
						return
					}

					return validateStringPrefix("/")(v, k)
				},
				Description: "Path of health check. Maximum length is 80.",
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      http.MethodHead,
				ValidateFunc: validateAllowedStringValue([]string{http.MethodGet, http.MethodHead}),
				Description:  "Method of the health check. Valid value: `GET` and `HEAD`.",
			},
			"health_check_status_codes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
				Computed:    true,
				Description: "Return code of confirmed normal. Valid value: `100`, `200`, `300`, `400` and `500`.",
			},
			"realservers": {
				Type:     schema.TypeSet,
				Optional: true,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return helper.HashString(fmt.Sprintf("%s-%s-%d-%d", m["id"].(string), m["ip"].(string), m["port"].(int), m["weight"].(int)))

				},
				Description: "An information list of GAAP realserver.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the GAAP realserver.",
						},
						"ip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "IP of the GAAP realserver.",
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validatePort,
							Description:  "Port of the GAAP realserver.",
						},
						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validateIntegerInRange(1, 100),
							Description:  "Scheduling weight, default value is `1`. Valid value ranges: (1~100).",
						},
					},
				},
			},
			"forward_host": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "The default value of requested host which is forwarded to the realserver by the listener is `default`.",
			},
			"sni_switch": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateAllowedStringValue([]string{GAAP_SERVER_NAME_INDICATION_SWITCH_ON, GAAP_SERVER_NAME_INDICATION_SWITCH_OFF}),
				Description:  "ServerNameIndication (SNI) switch. ON means on and OFF means off.",
			},
			"sni": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ServerNameIndication (SNI) is required when the SNI switch is turned on.",
			},
		},
	}
}

func resourceTencentCloudGaapHttpRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.create")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	listenerId := d.Get("listener_id").(string)

	rule := gaapHttpRule{
		listenerId:                 listenerId,
		domain:                     d.Get("domain").(string),
		path:                       d.Get("path").(string),
		realserverType:             d.Get("realserver_type").(string),
		scheduler:                  d.Get("scheduler").(string),
		healthCheck:                d.Get("health_check").(bool),
		interval:                   d.Get("interval").(int),
		connectTimeout:             d.Get("connect_timeout").(int),
		healthCheckPath:            d.Get("health_check_path").(string),
		healthCheckMethod:          d.Get("health_check_method").(string),
		forwardHost:                d.Get("forward_host").(string),
		serverNameIndicationSwitch: d.Get("sni_switch").(string),
		serverNameIndication:       d.Get("sni").(string),
	}

	if raw, ok := d.GetOk("health_check_status_codes"); ok {
		statusCodeSet := raw.(*schema.Set)

		codes := []interface{}{100, 200, 300, 400, 500}
		defaultSet := schema.NewSet(schema.HashInt, codes)
		diff := statusCodeSet.Difference(defaultSet)
		if diff.Len() > 0 {
			return fmt.Errorf("invalid health check status %v", diff.List())
		}

		rule.healthCheckStatusCodes = make([]int, 0, statusCodeSet.Len())
		for _, code := range statusCodeSet.List() {
			rule.healthCheckStatusCodes = append(rule.healthCheckStatusCodes, code.(int))
		}

	} else {
		rule.healthCheckStatusCodes = []int{100, 200, 300, 400, 500}
	}

	if len(rule.healthCheckStatusCodes) == 0 {
		return errors.New("health_check_status_codes can't be empty")
	}

	if rule.serverNameIndicationSwitch == GAAP_SERVER_NAME_INDICATION_SWITCH_ON && rule.serverNameIndication == "" {
		return fmt.Errorf("ServerNameIndication (SNI) is required when the SNI switch is turned on.")
	}
	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateHttpRule(ctx, rule)
	if err != nil {
		return err
	}

	d.SetId(id)

	if raw, ok := d.GetOk("realservers"); ok {
		realserverSet := raw.(*schema.Set).List()
		realservers := make([]gaapRealserverBind, 0, len(realserverSet))

		for _, v := range realserverSet {
			m := v.(map[string]interface{})
			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}

		if err := service.BindHttpRuleRealservers(ctx, rule.listenerId, id, realservers); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapHttpRuleRead(d, m)
}

func resourceTencentCloudGaapHttpRuleRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	rule, err := service.DescribeHttpRule(ctx, id)
	if err != nil {
		return err
	}

	if rule == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("listener_id", rule.ListenerId)
	_ = d.Set("domain", rule.Domain)
	_ = d.Set("path", rule.Path)
	_ = d.Set("realserver_type", rule.RealServerType)
	_ = d.Set("scheduler", rule.Scheduler)

	if rule.HealthCheck == nil {
		rule.HealthCheck = helper.IntUint64(0)
	}
	_ = d.Set("health_check", *rule.HealthCheck == 1)

	if rule.CheckParams == nil {
		rule.CheckParams = new(gaap.RuleCheckParams)
	}

	_ = d.Set("interval", rule.CheckParams.DelayLoop)
	_ = d.Set("connect_timeout", rule.CheckParams.ConnectTimeout)
	_ = d.Set("health_check_path", rule.CheckParams.Path)
	_ = d.Set("health_check_method", rule.CheckParams.Method)
	_ = d.Set("forward_host", rule.ForwardHost)
	_ = d.Set("health_check_status_codes", rule.CheckParams.StatusCode)
	_ = d.Set("sni_switch", rule.ServerNameIndicationSwitch)
	_ = d.Set("sni", rule.ServerNameIndication)

	realserverSet := make([]map[string]interface{}, 0, len(rule.RealServerSet))
	for _, rs := range rule.RealServerSet {
		realserverSet = append(realserverSet, map[string]interface{}{
			"id":     rs.RealServerId,
			"ip":     rs.RealServerIP,
			"port":   rs.RealServerPort,
			"weight": rs.RealServerWeight,
		})
	}
	_ = d.Set("realservers", realserverSet)

	return nil
}

func resourceTencentCloudGaapHttpRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.update")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	listenerId := d.Get("listener_id").(string)

	var (
		path      *string
		scheduler *string
	)
	sniSwitch := d.Get("sni_switch").(string)
	sni := d.Get("sni").(string)
	if sniSwitch == GAAP_SERVER_NAME_INDICATION_SWITCH_ON && sni == "" {
		return fmt.Errorf("ServerNameIndication (SNI) is required when the SNI switch is turned on.")
	}
	if d.HasChange("path") {
		path = helper.String(d.Get("path").(string))
	}

	if d.HasChange("scheduler") {
		scheduler = helper.String(d.Get("scheduler").(string))
	}

	healthCheck := d.Get("health_check").(bool)

	interval := d.Get("interval").(int)

	connectTimeout := d.Get("connect_timeout").(int)

	healthCheckPath := d.Get("health_check_path").(string)

	healthCheckMethod := d.Get("health_check_method").(string)

	var healthCheckStatusCodes []int
	if raw, ok := d.GetOk("health_check_status_codes"); ok {
		statusCodeSet := raw.(*schema.Set)

		codes := []interface{}{100, 200, 300, 400, 500}
		defaultSet := schema.NewSet(schema.HashInt, codes)
		diff := statusCodeSet.Difference(defaultSet)
		if diff.Len() > 0 {
			return fmt.Errorf("invalid health check status %v", diff.List())
		}

		healthCheckStatusCodes = make([]int, 0, statusCodeSet.Len())
		for _, code := range statusCodeSet.List() {
			healthCheckStatusCodes = append(healthCheckStatusCodes, code.(int))
		}
	} else {
		healthCheckStatusCodes = []int{100, 200, 300, 400, 500}
	}

	var realservers []gaapRealserverBind
	realserverUpdate := d.HasChange("realservers")

	if realserverUpdate {
		set := d.Get("realservers").(*schema.Set).List()
		realservers = make([]gaapRealserverBind, 0, len(set))
		for _, v := range set {
			m := v.(map[string]interface{})
			realservers = append(realservers, gaapRealserverBind{
				id:     m["id"].(string),
				ip:     m["ip"].(string),
				port:   m["port"].(int),
				weight: m["weight"].(int),
			})
		}
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	d.Partial(true)

	if err := service.ModifyHTTPRuleAttribute(
		ctx,
		listenerId, id, healthCheckPath, healthCheckMethod, sniSwitch, sni,
		path, scheduler, healthCheck, interval, connectTimeout, healthCheckStatusCodes,
	); err != nil {
		return err
	}

	if d.HasChange("forward_host") {
		forwardHost := d.Get("forward_host").(string)
		if err := service.ModifyHTTPRuleForwardHost(ctx, listenerId, id, forwardHost); err != nil {
			return err
		}

	}

	if realserverUpdate {
		if err := service.BindHttpRuleRealservers(ctx, listenerId, id, realservers); err != nil {
			return err
		}

	}

	d.Partial(false)

	return resourceTencentCloudGaapHttpRuleRead(d, m)
}

func resourceTencentCloudGaapHttpRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.delete")()
	gaapActionMu.Lock()
	defer gaapActionMu.Unlock()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	id := d.Id()
	listenerId := d.Get("listener_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteHttpRule(ctx, listenerId, id)
}
