package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapRealservers() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				Default:  -1,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateIp,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// computed
			"realservers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
