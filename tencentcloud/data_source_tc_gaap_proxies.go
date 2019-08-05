package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapProxies() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:          schema.TypeSet,
				Optional:      true,
				Elem:          schema.TypeString,
				ConflictsWith: []string{"project_id", "access_region", "realserver_region"},
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"ids"},
			},
			"access_region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
			},
			"realserver_region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"ids"},
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			// computed
			"proxies": {
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
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"concurrent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"access_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
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
						"version": {
							Type:     schema.TypeString,
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
