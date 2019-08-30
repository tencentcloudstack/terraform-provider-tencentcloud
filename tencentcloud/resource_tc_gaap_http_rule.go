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
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Set: func(v interface{}) int {
					return v.(int)
				},
				Description: "Return code of confirmed normal. Available values includes `100`,`200`,`300`,`400` and `500`.",
			},
			"realservers": {
				Type:     schema.TypeSet,
				Required: true,
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					sb := new(strings.Builder)
					sb.WriteString(m["id"].(string))
					sb.WriteString(m["ip"].(string))
					sb.WriteString(fmt.Sprintf("%d", m["port"].(int)))
					sb.WriteString(fmt.Sprintf("%d", m["weight"].(int)))
					return hashcode.String(sb.String())
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
		},
	}
}

func resourceTencentCloudGaapHttpRuleCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_http_rule.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rule := gaapHttpRule{
		listenerId:        d.Get("listener_id").(string),
		domain:            d.Get("domain").(string),
		path:              d.Get("path").(string),
		realserverType:    d.Get("realserver_type").(string),
		scheduler:         d.Get("scheduler").(string),
		healthCheck:       d.Get("health_check").(bool),
		interval:          d.Get("interval").(int),
		connectTimeout:    d.Get("connect_timeout").(int),
		healthCheckPath:   d.Get("health_check_path").(string),
		healthCheckMethod: d.Get("health_check_method").(string),
	}

	if raw, ok := d.GetOk("health_check_status_codes"); ok {
		statusCodeSet := raw.(*schema.Set).List()
		rule.healthCheckStatusCodes = make([]int, 0, len(statusCodeSet))
		for _, c := range statusCodeSet {
			code := c.(int)
			switch code {
			case 100, 200, 300, 400, 500:
				rule.healthCheckStatusCodes = append(rule.healthCheckStatusCodes, code)

			default:
				return fmt.Errorf("invalid health check status code %d", code)
			}
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
	listenerId := d.Get("listener_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	httpRule, realservers, err := service.DescribeHttpRule(ctx, listenerId, id)
	if err != nil {
		return err
	}

	if httpRule == nil {
		d.SetId("")
		return nil
	}

	d.Set("domain", httpRule.domain)
	d.Set("path", httpRule.path)
	d.Set("realserver_type", httpRule.realserverType)
	d.Set("scheduler", httpRule.scheduler)
	d.Set("health_check", httpRule.healthCheck)
	d.Set("interval", httpRule.interval)
	d.Set("connect_timeout", httpRule.connectTimeout)
	d.Set("health_check_path", httpRule.healthCheckPath)
	d.Set("health_check_method", httpRule.healthCheckMethod)

	if _, ok := d.GetOk("health_check_status_codes"); ok || len(httpRule.healthCheckStatusCodes) != 5 {
		d.Set("health_check_status_codes", httpRule.healthCheckStatusCodes)
	}

	realserverSet := make([]map[string]interface{}, 0, len(realservers))
	for _, rs := range realservers {
		realserverSet = append(realserverSet, map[string]interface{}{
			"id":     rs.id,
			"ip":     rs.ip,
			"port":   rs.port,
			"weight": rs.weight,
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
		statusCodeSet := raw.(*schema.Set).List()
		healthCheckStatusCodes = make([]int, 0, len(statusCodeSet))
		for _, code := range statusCodeSet {
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

	if len(realservers) > 0 {
		if err := service.BindHttpRuleRealservers(ctx, listenerId, id, realservers); err != nil {
			return err
		}
		d.SetPartial("realservers")
	}

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
