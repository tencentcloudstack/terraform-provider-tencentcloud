/*
Provides a resource to create a ci poster_template

Example Usage

```hcl
resource "tencentcloud_ci_poster_template" "poster_template" {
  input = &lt;nil&gt;
  name = &lt;nil&gt;
  category_ids = &lt;nil&gt;
}
```

Import

ci poster_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_poster_template.poster_template poster_template_id
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

func resourceTencentCloudCiPosterTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiPosterTemplateCreate,
		Read:   resourceTencentCloudCiPosterTemplateRead,
		Delete: resourceTencentCloudCiPosterTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"input": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "PSD files in the COS bucket, the size limit is 100M.",
			},

			"name": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Template name.",
			},

			"category_ids": {
				Optional: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Template category ID, support to pass in more than one, separated by `,` symbol.",
			},
		},
	}
}

func resourceTencentCloudCiPosterTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_poster_template.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request   = ci.NewAddStyleRequest()
		response  = ci.NewAddStyleResponse()
		styleName string
	)
	if v, ok := d.GetOk("input"); ok {
		inputSet := v.(*schema.Set).List()
		for i := range inputSet {
			input := inputSet[i].(string)
			request.Input = append(request.Input, &input)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("category_ids"); ok {
		categoryIdsSet := v.(*schema.Set).List()
		for i := range categoryIdsSet {
			categoryIds := categoryIdsSet[i].(string)
			request.CategoryIds = append(request.CategoryIds, &categoryIds)
		}
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
		log.Printf("[CRITAL]%s create ci posterTemplate failed, reason:%+v", logId, err)
		return err
	}

	styleName = *response.Response.StyleName
	d.SetId(styleName)

	return resourceTencentCloudCiPosterTemplateRead(d, meta)
}

func resourceTencentCloudCiPosterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_poster_template.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	posterTemplateId := d.Id()

	posterTemplate, err := service.DescribeCiPosterTemplateById(ctx, styleName)
	if err != nil {
		return err
	}

	if posterTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiPosterTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if posterTemplate.Input != nil {
		_ = d.Set("input", posterTemplate.Input)
	}

	if posterTemplate.Name != nil {
		_ = d.Set("name", posterTemplate.Name)
	}

	if posterTemplate.CategoryIds != nil {
		_ = d.Set("category_ids", posterTemplate.CategoryIds)
	}

	return nil
}

func resourceTencentCloudCiPosterTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_poster_template.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}
	posterTemplateId := d.Id()

	if err := service.DeleteCiPosterTemplateById(ctx, styleName); err != nil {
		return err
	}

	return nil
}
