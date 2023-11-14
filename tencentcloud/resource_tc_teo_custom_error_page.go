/*
Provides a resource to create a teo custom_error_page

Example Usage

```hcl
resource "tencentcloud_teo_custom_error_page" "custom_error_page" {
  zone_id = &lt;nil&gt;
  entity = &lt;nil&gt;
    name = &lt;nil&gt;
  content = &lt;nil&gt;
}
```

Import

teo custom_error_page can be imported using the id, e.g.

```
terraform import tencentcloud_teo_custom_error_page.custom_error_page custom_error_page_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTeoCustomErrorPage() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoCustomErrorPageCreate,
		Read:   resourceTencentCloudTeoCustomErrorPageRead,
		Update: resourceTencentCloudTeoCustomErrorPageUpdate,
		Delete: resourceTencentCloudTeoCustomErrorPageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Site ID.",
			},

			"entity": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Subdomain.",
			},

			"page_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Page ID.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Page name.",
			},

			"content": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Page content.",
			},
		},
	}
}

func resourceTencentCloudTeoCustomErrorPageCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = teo.NewCreateCustomErrorPageRequest()
		response = teo.NewCreateCustomErrorPageResponse()
		zoneId   string
		entity   string
		pageId   string
	)
	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
		request.Entity = helper.String(v.(string))
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("content"); ok {
		request.Content = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTeoClient().CreateCustomErrorPage(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo customErrorPage failed, reason:%+v", logId, err)
		return err
	}

	zoneId = *response.Response.ZoneId
	d.SetId(strings.Join([]string{zoneId, entity, pageId}, FILED_SP))

	return resourceTencentCloudTeoCustomErrorPageRead(d, meta)
}

func resourceTencentCloudTeoCustomErrorPageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]
	pageId := idSplit[2]

	customErrorPage, err := service.DescribeTeoCustomErrorPageById(ctx, zoneId, entity, pageId)
	if err != nil {
		return err
	}

	if customErrorPage == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TeoCustomErrorPage` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if customErrorPage.ZoneId != nil {
		_ = d.Set("zone_id", customErrorPage.ZoneId)
	}

	if customErrorPage.Entity != nil {
		_ = d.Set("entity", customErrorPage.Entity)
	}

	if customErrorPage.PageId != nil {
		_ = d.Set("page_id", customErrorPage.PageId)
	}

	if customErrorPage.Name != nil {
		_ = d.Set("name", customErrorPage.Name)
	}

	if customErrorPage.Content != nil {
		_ = d.Set("content", customErrorPage.Content)
	}

	return nil
}

func resourceTencentCloudTeoCustomErrorPageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"zone_id", "entity", "page_id", "name", "content"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudTeoCustomErrorPageRead(d, meta)
}

func resourceTencentCloudTeoCustomErrorPageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TeoService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]
	pageId := idSplit[2]

	if err := service.DeleteTeoCustomErrorPageById(ctx, zoneId, entity, pageId); err != nil {
		return err
	}

	return nil
}
