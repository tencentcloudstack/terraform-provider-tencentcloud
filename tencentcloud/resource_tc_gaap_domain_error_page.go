/*
Provide a resource to custom error page info for a GAAP HTTP domain.

Example Usage

```hcl
resource tencentcloud_gaap_proxy "foo" {
  name              = "ci-test-gaap-proxy"
  bandwidth         = 10
  concurrent        = 2
  access_region     = "SouthChina"
  realserver_region = "NorthChina"
}

resource tencentcloud_gaap_layer7_listener "foo" {
  protocol = "HTTP"
  name     = "ci-test-gaap-l7-listener"
  port     = 80
  proxy_id = tencentcloud_gaap_proxy.foo.id
}

resource tencentcloud_gaap_http_domain "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = "www.qq.com"
}

resource tencentcloud_gaap_domain_error_page "foo" {
  listener_id = tencentcloud_gaap_layer7_listener.foo.id
  domain      = tencentcloud_gaap_http_domain.foo.domain
  error_codes = [404, 503]
  body        = "bad request"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapDomainErrorPageInfo() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapDomainErrorPageInfoCreate,
		Read:   resourceTencentCloudGaapDomainErrorPageInfoRead,
		Delete: resourceTencentCloudGaapDomainErrorPageInfoDelete,
		Schema: map[string]*schema.Schema{
			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the layer7 listener.",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "HTTP domain.",
			},
			"error_codes": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Set:         schema.HashInt,
				Description: "Original error codes.",
			},
			"body": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "New response body.",
			},
			"new_error_code": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "New error code.",
			},
			"clear_headers": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "Response headers to be removed.",
			},
			"set_headers": {
				Type:        schema.TypeMap,
				ForceNew:    true,
				Optional:    true,
				Description: "Response headers to be set.",
			},
		},
	}
}

func resourceTencentCloudGaapDomainErrorPageInfoCreate(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_domain_error_page.create")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)
	body := d.Get("body").(string)

	var errorCodes []int
	for _, errorNo := range d.Get("error_codes").(*schema.Set).List() {
		errorCodes = append(errorCodes, errorNo.(int))
	}

	var (
		newErrorCode *int64
		clearHeaders []string
	)

	if raw, ok := d.GetOk("new_error_code"); ok {
		newErrorCode = helper.IntInt64(raw.(int))
	}

	if raw, ok := d.GetOk("clear_headers"); ok {
		for _, clearHeader := range raw.(*schema.Set).List() {
			clearHeaders = append(clearHeaders, clearHeader.(string))
		}
	}

	setHeaders := helper.GetTags(d, "set_headers")

	id, err := service.CreateDomainErrorPageInfo(ctx, listenerId, domain, body, newErrorCode, errorCodes, clearHeaders, setHeaders)
	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceTencentCloudGaapDomainErrorPageInfoRead(d, m)
}

func resourceTencentCloudGaapDomainErrorPageInfoRead(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_domain_error_page.read")()
	defer inconsistentCheck(d, m)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	listenerId := d.Get("listener_id").(string)
	domain := d.Get("domain").(string)

	info, err := service.DescribeDomainErrorPageInfo(ctx, listenerId, domain, id)
	if err != nil {
		return err
	}

	if info == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("listener_id", info.ListenerId)
	_ = d.Set("domain", info.Domain)
	_ = d.Set("error_codes", info.ErrorNos)
	_ = d.Set("body", info.Body)
	_ = d.Set("new_error_code", info.NewErrorNo)
	_ = d.Set("clear_headers", info.ClearHeaders)

	setHeaders := make(map[string]string, len(info.SetHeaders))
	for _, header := range info.SetHeaders {
		setHeaders[*header.HeaderName] = *header.HeaderValue
	}

	_ = d.Set("set_headers", setHeaders)

	return nil
}

func resourceTencentCloudGaapDomainErrorPageInfoDelete(d *schema.ResourceData, m interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_domain_error_page.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	service := GaapService{client: m.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	return service.DeleteDomainErrorPageInfo(ctx, id)
}
