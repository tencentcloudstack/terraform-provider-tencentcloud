/*
Provides a resource to create a teo custom_error_page

Example Usage

```hcl
resource "tencentcloud_teo_custom_error_page" "error_page_0" {
  zone_id = data.tencentcloud_teo_zone_ddos_policy.zone_policy.zone_id
  entity  = data.tencentcloud_teo_zone_ddos_policy.zone_policy.shield_areas[0].application[0].host

  name    = "test"
  content = "<html lang='en'><body><div><p>test content</p></div></body></html>"
}

```
*/
package tencentcloud

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTeoCustomErrorPage() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudTeoCustomErrorPageRead,
		Create: resourceTencentCloudTeoCustomErrorPageCreate,
		Update: resourceTencentCloudTeoCustomErrorPageUpdate,
		Delete: resourceTencentCloudTeoCustomErrorPageDelete,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"entity": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Subdomain.",
			},

			"page_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Page ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Page name.",
			},

			"content": {
				Type:        schema.TypeString,
				Required:    true,
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
		response *teo.CreateCustomErrorPageResponse
		zoneId   string
		entity   string
		pageId   int64
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
		request.ZoneId = &zoneId
	}

	if v, ok := d.GetOk("entity"); ok {
		entity = v.(string)
		request.Entity = &entity
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
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo customErrorPage failed, reason:%+v", logId, err)
		return err
	}

	pageId = *response.Response.PageId

	d.SetId(zoneId + FILED_SP + entity + FILED_SP + helper.Int64ToStr(pageId))
	return resourceTencentCloudTeoCustomErrorPageRead(d, meta)
}

func resourceTencentCloudTeoCustomErrorPageRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.read")()
	defer inconsistentCheck(d, meta)()

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	entity := idSplit[1]
	pageId := idSplit[2]

	_ = d.Set("zone_id", zoneId)
	_ = d.Set("entity", entity)
	_ = d.Set("page_id", pageId)

	if v, ok := d.GetOk("name"); ok {
		_ = d.Set("name", helper.String(v.(string)))
	}

	if v, ok := d.GetOk("content"); ok {
		_ = d.Set("content", helper.String(v.(string)))
	}

	return nil
}

func resourceTencentCloudTeoCustomErrorPageUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.update")()
	defer inconsistentCheck(d, meta)()

	var change bool

	if d.HasChange("zone_id") {
		return fmt.Errorf("`zone_id` do not support change now.")
	}

	if d.HasChange("entity") {
		return fmt.Errorf("`entity` do not support change now.")
	}

	if d.HasChange("name") {
		change = true
	}

	if d.HasChange("content") {
		change = true
	}

	if change {
		err := resourceTencentCloudTeoCustomErrorPageCreate(d, meta)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudTeoCustomErrorPageDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_teo_custom_error_page.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
