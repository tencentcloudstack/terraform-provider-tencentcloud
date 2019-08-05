package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapLayer4Listeners() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"TCP", "UDP"}),
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"listener_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validatePort,
			},

			// computed
			"listeners": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"realserver_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeInt,
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
						"connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
