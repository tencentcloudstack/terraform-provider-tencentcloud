package tencentcloud

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceTencentCloudGaapSecurityRule() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cidr_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errs []error) {
					if _, err := validateIp(v, k); len(err) == 0 {
						return
					}

					if _, err := validateCIDRNetworkAddress(v, k); len(err) != 0 {
						errs = append(errs, fmt.Errorf("%s %v is not valid IP address or valid CIDR IP address",
							k, v))
					}
					return
				},
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				Default:      "ALL",
				ValidateFunc: validateAllowedStringValue([]string{"ALL", "TCP", "UDP"}),
				ForceNew:     true,
			},
			"port": {
				Type:     schema.TypeString,
				Required: true,
				Default:  "ALL",
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value == "ALL" {
						return
					}
					match, _ := regexp.MatchString("^(\\d{1,5},)*\\d{1,5}$|^\\d{1,5}-\\d{1,5}$", value)
					if !match {
						errors = append(errors, fmt.Errorf("%s example: 53、80,443、80-90, Not configured to represent all ports", k))
					}
					return
				},
			},
		},
	}
}
