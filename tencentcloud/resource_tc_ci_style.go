/*
Provides a resource to create a ci style

Example Usage

```hcl
resource "tencentcloud_ci_style" "style" {
  style_name = &lt;nil&gt;
  style_body = &lt;nil&gt;
}
```

Import

ci style can be imported using the id, e.g.

```
terraform import tencentcloud_ci_style.style style_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
)

func resourceTencentCloudCiStyle() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiStyleCreate,
		Read:   resourceTencentCloudCiStyleRead,
		Delete: resourceTencentCloudCiStyleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"style_name": {
				Required:    true,
				ForceNew:    true,
				Description: "Style name.",
			},

			"style_body": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Style details.",
			},
		},
	}
}

func resourceTencentCloudCiStyleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_style.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = ci.NewAddStyleRequest()
		response  = ci.NewAddStyleResponse()
		styleName string
	)
	if v, _ := d.GetOk("style_name"); v != nil {
	}

	if v, ok := d.GetOk("style_body"); ok {
		request.StyleBody = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().AddStyle(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci style failed, reason:%+v", logId, err)
		return err
	}

	styleName = *response.Response.StyleName
	d.SetId(styleName)

	return resourceTencentCloudCiStyleRead(d, meta)
}

func resourceTencentCloudCiStyleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_style.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	styleId := d.Id()

	style, err := service.DescribeCiStyleById(ctx, styleName)
	if err != nil {
		return err
	}

	if style == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiStyle` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if style.StyleName != nil {
		_ = d.Set("style_name", style.StyleName)
	}

	if style.StyleBody != nil {
		_ = d.Set("style_body", style.StyleBody)
	}

	return nil
}

func resourceTencentCloudCiStyleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_style.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	styleId := d.Id()

	if err := service.DeleteCiStyleById(ctx, styleName); err != nil {
		return err
	}

	return nil
}
