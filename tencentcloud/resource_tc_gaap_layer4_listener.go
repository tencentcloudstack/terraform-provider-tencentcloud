package tencentcloud

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapLayer4Listener() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP"}),
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validatePort,
			},
			"scheduler": {
				Type:         schema.TypeString,
				Required:     true,
				Default:      "rr",
				ValidateFunc: validateAllowedStringValue([]string{"rr", "wr", "lc"}),
			},
			"realserver_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"IP", "DOMAIN"}),
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"health_check": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"delay_loop": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"realserver_bind_set": {
				Type:     schema.TypeSet,
				Optional: true,
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

			// computed
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
