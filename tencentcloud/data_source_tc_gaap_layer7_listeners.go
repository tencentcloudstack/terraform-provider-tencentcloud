package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapLayer7Listeners() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"HTTP", "HTTPS"}),
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
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"forward_protocol": {
							Type:     schema.TypeString,
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
