/*
Use this data source to query security policies of GAAP proxy.

Example Usage

```hcl
resource "tencentcloud_gaap_proxy" "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource "tencentcloud_gaap_security_policy" "foo" {
  proxy_id = "${tencentcloud_gaap_proxy.foo.id}"
  action   = "ACCEPT"
}

data "tencentcloud_gaap_security_policies" "foo" {
  id = "${tencentcloud_gaap_security_policy.foo.id}"
}
```
*/
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
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the security policy to be queried.",
			},

			// computed
			"proxy_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the GAAP proxy.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Status of the security policy.",
			},
			"action": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default policy.",
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
		d.SetId("")
		return nil
	}

	d.Set("proxy_id", proxyId)
	d.Set("status", status)
	d.Set("action", action)

	d.SetId(id)

	return nil
}
