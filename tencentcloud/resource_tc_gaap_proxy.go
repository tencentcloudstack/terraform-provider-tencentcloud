package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func resourceTencentCloudGaapProxy() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateStringLengthInRange(1, 30),
			},
			"project_id": {
				Type:     schema.TypeInt,
				Required: true,
				Default:  0,
			},
			"bandwidth": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{10, 20, 50, 100, 200, 500, 1000}),
			},
			"concurrent": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validateAllowedIntValue([]int{2, 5, 10, 20, 30, 40, 50, 60, 70, 80, 90, 100}),
			},
			"access_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"realserver_region": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
				Default:  true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// computed
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scalarable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"support_protocols": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     schema.TypeString,
			},
			"forward_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
