/*
Provides a resource to create a security policy of GAAP proxy.

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
  proxy_id = tencentcloud_gaap_proxy.foo.id
  action   = "DROP"
}
```

Import

GAAP security policy can be imported using the id, e.g.

```
  $ terraform import tencentcloud_gaap_security_policy.foo pl-xxxx
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTencentCloudGaapSecurityPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapSecurityPolicyCreate,
		Read:   resourceTencentCloudGaapSecurityPolicyRead,
		Update: resourceTencentCloudGaapSecurityPolicyUpdate,
		Delete: resourceTencentCloudGaapSecurityPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"proxy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the GAAP proxy.",
			},
			"action": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAllowedStringValue([]string{"ACCEPT", "DROP"}),
				ForceNew:     true,
				Description:  "Default policy, the available values include `ACCEPT` and `DROP`.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether policy is enable, default value is `true`.",
			},
		},
	}
}

func resourceTencentCloudGaapSecurityPolicyCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_policy.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	proxyId := d.Get("proxy_id").(string)
	action := d.Get("action").(string)
	enable := d.Get("enable").(bool)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id, err := service.CreateSecurityPolicy(ctx, proxyId, action)
	if err != nil {
		return err
	}

	d.SetId(id)

	if enable {
		if err := service.EnableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapSecurityPolicyRead(d, m)
}

func resourceTencentCloudGaapSecurityPolicyRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_policy.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	proxyId, status, action, exist, err := service.DescribeSecurityPolicy(ctx, id)
	if err != nil {
		return err
	}

	if !exist {
		d.SetId("")
		return nil
	}

	_ = d.Set("proxy_id", proxyId)
	_ = d.Set("action", action)
	_ = d.Set("enable", status == GAAP_SECURITY_POLICY_BOUND)

	return nil
}

func resourceTencentCloudGaapSecurityPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_policy.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	proxyId := d.Get("proxy_id").(string)
	enable := d.Get("enable").(bool)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if enable {
		if err := service.EnableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	} else {
		if err := service.DisableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return resourceTencentCloudGaapSecurityPolicyRead(d, m)
}

func resourceTencentCloudGaapSecurityPolicyDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_security_policy.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	id := d.Id()
	enable := d.Get("enable").(bool)
	proxyId := d.Get("proxy_id").(string)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	if enable {
		if err := service.DisableSecurityPolicy(ctx, proxyId, id); err != nil {
			return err
		}
	}

	return service.DeleteSecurityPolicy(ctx, id)
}
