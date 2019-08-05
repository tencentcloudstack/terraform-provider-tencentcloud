package tencentcloud

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceTencentCloudGaapHttpRules() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"path": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateStringPrefix("/"),
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"client_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"basic_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"basic_auth_config_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"basic_auth_config_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realserver_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"realserver_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realserver_certificate_Name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gaap_auth": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"gaap_certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gaap_certificate_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"realserver_certificate_domain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// computed
			"rules": {
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
						"listener_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scheduler": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay_loop": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"connect_timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"health_check_path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_method": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"health_check_status_codes": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     schema.TypeInt,
						},
						"health_check_domain": {
							Type:     schema.TypeString,
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
						"basic_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"basic_auth_config_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"basic_auth_config_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"realserver_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_certificate_Name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gaap_auth": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"gaap_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gaap_certificate_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"realserver_certificate_domain": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"realservers": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
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
									"weight": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"status": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
