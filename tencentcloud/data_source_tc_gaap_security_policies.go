package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceTencentCloudGaapSecurityPolices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapSecurityPoliciesRead,
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

func dataSourceTencentCloudGaapSecurityPoliciesRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_security_policies.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Get("id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	proxyId, status, action, exist, err := service.DescribeSecurityPolicy(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		d.SetId(id)
		return nil
	}

	d.Set("proxy_id", proxyId)
	d.Set("status", status)
	d.Set("action", action)

	d.SetId(id)

	return nil
}
