package tencentcloud

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapHttpRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringPrefix("/"),
			},
			"realserver_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"IP", "DOMAIN"}),
			},
			"scheduler": {
				Type:         schema.TypeString,
				Required:     true,
				Default:      "rr",
				ValidateFunc: validateAllowedStringValue([]string{"rr", "wr", "lc"}),
			},
			"health_check": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"delay_loop": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"health_check_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringPrefix("/"),
			},
			"health_check_method": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateAllowedStringValue([]string{http.MethodGet, http.MethodPost}),
			},
			"health_check_status_codes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     schema.TypeInt,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) { // TODO need confirm
					set := v.(*schema.Set).List()
					for _, v := range set {
						value := v.(int)
						switch value {
						case 100, 200, 300, 400, 500:
						default:
							errs = append(errs, fmt.Errorf("%s has invalid value %d", k, value))
							return
						}
					}
					return
				},
			},
			"health_check_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateIp,
						},
						"port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validatePort,
						},
						"weight": {
							Type:         schema.TypeInt,
							Required:     true,
							Default:      1,
							ValidateFunc: validateIntegerMin(1),
						},
					},
				},
			},
		},
	}
}
