/*
Provides a resource to create a tag tag

Example Usage

```hcl
resource "tencentcloud_tag_tag" "tag" {
  tag_key = "test"
  tag_value = "Terraform"
}
```

Import

tag tag can be imported using the id, e.g.

```
terraform import tencentcloud_tag_tag.tag tag_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"strings"
)

func resourceTencentCloudTagTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTagTagCreate,
		Read:   resourceTencentCloudTagTagRead,
		Update: resourceTencentCloudTagTagUpdate,
		Delete: resourceTencentCloudTagTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Tag key.",
			},

			"tag_value": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Tag value.",
			},
		},
	}
}

func resourceTencentCloudTagTagCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_tag.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tag.NewCreateTagRequest()
		response = tag.NewCreateTagResponse()
		tagKey   string
		tagValue string
	)
	if v, ok := d.GetOk("tag_key"); ok {
		tagKey = v.(string)
		request.TagKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_value"); ok {
		tagValue = v.(string)
		request.TagValue = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTagClient().CreateTag(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tag tag failed, reason:%+v", logId, err)
		return err
	}

	tagKey = *response.Response.TagKey
	d.SetId(strings.Join([]string{tagKey, tagValue}, FILED_SP))

	return resourceTencentCloudTagTagRead(d, meta)
}

func resourceTencentCloudTagTagRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_tag.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tagKey := idSplit[0]
	tagValue := idSplit[1]

	tag, err := service.DescribeTagTagById(ctx, tagKey, tagValue)
	if err != nil {
		return err
	}

	if tag == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TagTag` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tag.TagKey != nil {
		_ = d.Set("tag_key", tag.TagKey)
	}

	if tag.TagValue != nil {
		_ = d.Set("tag_value", tag.TagValue)
	}

	return nil
}

func resourceTencentCloudTagTagUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_tag.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	immutableArgs := []string{"tag_key", "tag_value"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}
	return resourceTencentCloudTagTagRead(d, meta)
}

func resourceTencentCloudTagTagDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_tag.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tagKey := idSplit[0]
	tagValue := idSplit[1]

	if err := service.DeleteTagTagById(ctx, tagKey, tagValue); err != nil {
		return err
	}

	return nil
}
