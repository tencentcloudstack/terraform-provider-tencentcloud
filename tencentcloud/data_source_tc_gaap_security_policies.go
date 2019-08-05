package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapSecurityPolices() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed
			"proxy_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
