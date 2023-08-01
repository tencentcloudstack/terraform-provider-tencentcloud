/*
Provides a resource to create a tag attachment

Example Usage

```hcl

resource "tencentcloud_tag_attachment" "attachment" {
  tag_key = "test3"
  tag_value = "Terraform3"
  resource = "qcs::cvm:ap-guangzhou:uin/100020512675:instance/ins-kfrlvcp4"
}

```

Import

tag attachment can be imported using the id, e.g.

```
terraform import tencentcloud_tag_attachment.attachment attachment_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudTagAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTagAttachmentCreate,
		Read:   resourceTencentCloudTagAttachmentRead,
		Delete: resourceTencentCloudTagAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "tag key.",
			},

			"tag_value": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "tag value.",
			},

			"resource": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "[Six-segment description of resources](https://cloud.tencent.com/document/product/598/10606).",
			},
		},
	}
}

func resourceTencentCloudTagAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = tag.NewAddResourceTagRequest()
		tagKey     string
		tagValue   string
		resourceId string
	)
	if v, ok := d.GetOk("tag_key"); ok {
		tagKey = v.(string)
		request.TagKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_value"); ok {
		tagValue = v.(string)
		request.TagValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource"); ok {
		resourceId = v.(string)
		request.Resource = helper.String(v.(string))
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseTagClient().AddResourceTag(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tag tagAttachment failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(tagKey + FILED_SP + tagValue + FILED_SP + resourceId)

	return resourceTencentCloudTagAttachmentRead(d, meta)
}

func resourceTencentCloudTagAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tagKey := idSplit[0]
	tagValue := idSplit[1]
	resource := idSplit[2]

	tagAttachment, err := service.DescribeTagTagAttachmentById(ctx, tagKey, tagValue, resource)
	if err != nil {
		return err
	}

	if tagAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TagResourceTag` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if len(tagAttachment.Tags) < 1 {
		log.Printf("[WARN]%s resource `TagResourceTag` [%s] Tags is null, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if tagAttachment.Tags[0].TagKey != nil {
		_ = d.Set("tag_key", tagAttachment.Tags[0].TagKey)
	}

	if tagAttachment.Tags[0].TagValue != nil {
		_ = d.Set("tag_value", tagAttachment.Tags[0].TagValue)
	}

	if tagAttachment.Resource != nil {
		_ = d.Set("resource", tagAttachment.Resource)
	}

	return nil
}

func resourceTencentCloudTagAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TagService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 3 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	tagKey := idSplit[0]
	resource := idSplit[2]

	if err := service.DeleteTagTagAttachmentById(ctx, tagKey, resource); err != nil {
		return err
	}

	return nil
}
