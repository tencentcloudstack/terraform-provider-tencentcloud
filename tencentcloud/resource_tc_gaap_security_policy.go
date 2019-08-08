package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func resourceTencentCloudGaapSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
				Default:  true,
			},

			// computed
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
