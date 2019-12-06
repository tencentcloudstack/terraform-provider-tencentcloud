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
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
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
  listener_id = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain      = "www.qq.com"
}

resource "tencentcloud_gaap_http_rule" "foo" {
  listener_id               = "${tencentcloud_gaap_layer7_listener.foo.id}"
  domain                    = "${tencentcloud_gaap_http_domain.foo.domain}"
  path                      = "/"
  realserver_type           = "IP"
  health_check              = true
  health_check_path         = "/"
  health_check_method       = "GET"
  health_check_status_codes = [200]

  realservers {
    id   = "${tencentcloud_gaap_realserver.foo.id}"
    ip   = "${tencentcloud_gaap_realserver.foo.ip}"
    port = 80
  }

  realservers {
    id   = "${tencentcloud_gaap_realserver.bar.id}"
    ip   = "${tencentcloud_gaap_realserver.bar.ip}"
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

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
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
				Description: "Forward rule domain of the layer7 listener.",
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
				Description:  "Type of the realserver, and the available values include `IP`,`DOMAIN`.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "rr",
				ValidateFunc: validateAllowedStringValue([]string{"rr", "wrr", "lc"}),
				Description:  "Scheduling policy of the layer4 listener, default is `rr`. Available values include `rr`,`wrr` and `lc`.",
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
				Description:  "Interval of the health check, default is 5s.",
			},
			"connect_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validateIntegerInRange(2, 60),
				Description:  "Timeout of the health check response, default is 2s.",
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
				Description:  "Method of the health check. Available values includes `GET` and `HEAD`.",
			},
			"health_check_status_codes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
				Computed:    true,
				Description: "Return code of confirmed normal. Available values includes `100`,`200`,`300`,`400` and `500`.",
			},
			"realservers": {
				Type:     schema.TypeSet,
				Required: true,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return hashcode.String(fmt.Sprintf("%s-%s-%d-%d", m["id"].(string), m["ip"].(string), m["port"].(int), m["weight"].(int)))

				},
				Description: "An information list of GAAP realserver. Each element contains the following attributes:",
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
							Description:  "Scheduling weight, default is 1. The range of values is [1,100].",
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
		},
	}
}

func resourceTencentCloudGaapHttpRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	listenerId := d.Get("listener_id").(string)

	rule := gaapHttpRule{
		listenerId:        listenerId,
		domain:            d.Get("domain").(string),
		path:              d.Get("path").(string),
		realserverType:    d.Get("realserver_type").(string),
		scheduler:         d.Get("scheduler").(string),
		healthCheck:       d.Get("health_check").(bool),
		interval:          d.Get("interval").(int),
		connectTimeout:    d.Get("connect_timeout").(int),
		healthCheckPath:   d.Get("health_check_path").(string),
		healthCheckMethod: d.Get("health_check_method").(string),
		forwardHost:       d.Get("forward_host").(string),
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

	realserverSet := d.Get("realservers").(*schema.Set).List()
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

	if len(realservers) == 0 {
		return errors.New("realserver can't be empty")
	}

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateHttpRule(ctx, rule)
	if err != nil {
		return err
	}

	d.SetId(id)

	if err := service.BindHttpRuleRealservers(ctx, rule.listenerId, id, realservers); err != nil {
		return err
	}

	return resourceTencentCloudGaapHttpRuleRead(d, m)
}

func resourceTencentCloudGaapHttpRuleRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

	d.Set("listener_id", rule.ListenerId)
	d.Set("domain", rule.Domain)
	d.Set("path", rule.Path)
	d.Set("realserver_type", rule.RealServerType)
	d.Set("scheduler", rule.Scheduler)

	if rule.HealthCheck == nil {
		rule.HealthCheck = intToPointer(0)
	}
	d.Set("health_check", *rule.HealthCheck == 1)

	if rule.CheckParams == nil {
		rule.CheckParams = new(gaap.RuleCheckParams)
	}

	d.Set("interval", rule.CheckParams.DelayLoop)
	d.Set("connect_timeout", rule.CheckParams.ConnectTimeout)
	d.Set("health_check_path", rule.CheckParams.Path)
	d.Set("health_check_method", rule.CheckParams.Method)
	d.Set("forward_host", rule.ForwardHost)
	d.Set("health_check_status_codes", rule.CheckParams.StatusCode)

	realserverSet := make([]map[string]interface{}, 0, len(rule.RealServerSet))
	for _, rs := range rule.RealServerSet {
		realserverSet = append(realserverSet, map[string]interface{}{
			"id":     rs.RealServerId,
			"ip":     rs.RealServerIP,
			"port":   rs.RealServerPort,
			"weight": rs.RealServerWeight,
		})
	}
	d.Set("realservers", realserverSet)

	return nil
}

func resourceTencentCloudGaapHttpRuleUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	listenerId := d.Get("listener_id").(string)

	var (
		path       *string
		scheduler  *string
		updateAttr []string
	)

	if d.HasChange("path") {
		updateAttr = append(updateAttr, "path")
		path = stringToPointer(d.Get("path").(string))
	}

	if d.HasChange("scheduler") {
		updateAttr = append(updateAttr, "scheduler")
		scheduler = stringToPointer(d.Get("scheduler").(string))
	}

	if d.HasChange("health_check") {
		updateAttr = append(updateAttr, "health_check")
	}
	healthCheck := d.Get("health_check").(bool)

	if d.HasChange("interval") {
		updateAttr = append(updateAttr, "interval")
	}
	interval := d.Get("interval").(int)

	if d.HasChange("connect_timeout") {
		updateAttr = append(updateAttr, "connect_timeout")
	}
	connectTimeout := d.Get("connect_timeout").(int)

	if d.HasChange("health_check_path") {
		updateAttr = append(updateAttr, "health_check_path")
	}
	healthCheckPath := d.Get("health_check_path").(string)

	if d.HasChange("health_check_method") {
		updateAttr = append(updateAttr, "health_check_method")
	}
	healthCheckMethod := d.Get("health_check_method").(string)

	if d.HasChange("health_check_status_codes") {
		updateAttr = append(updateAttr, "health_check_status_codes")
	}
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
	if d.HasChange("realservers") {
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
		listenerId, id, healthCheckPath, healthCheckMethod,
		path, scheduler, healthCheck, interval, connectTimeout, healthCheckStatusCodes,
	); err != nil {
		return err
	}

	for _, attr := range updateAttr {
		d.SetPartial(attr)
	}

	if d.HasChange("forward_host") {
		forwardHost := d.Get("forward_host").(string)
		if err := service.ModifyHTTPRuleForwardHost(ctx, listenerId, id, forwardHost); err != nil {
			return err
		}

		d.SetPartial("forward_host")
	}

	if len(realservers) > 0 {
		if err := service.BindHttpRuleRealservers(ctx, listenerId, id, realservers); err != nil {
			return err
		}
		d.SetPartial("realservers")
	}

	d.Partial(false)

	return resourceTencentCloudGaapHttpRuleRead(d, m)
}

func resourceTencentCloudGaapHttpRuleDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	listenerId := d.Get("listener_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	return service.DeleteHttpRule(ctx, listenerId, id)
}
