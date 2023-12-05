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

func resourceTencentCloudTag() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTagResourceCreate,
		Read:   resourceTencentCloudTagResourceRead,
		Delete: resourceTencentCloudTagResourceDelete,
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
		},
	}
}

func resourceTencentCloudTagResourceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = tag.NewCreateTagRequest()
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
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tag tag failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(tagKey + FILED_SP + tagValue)

	return resourceTencentCloudTagResourceRead(d, meta)
}

func resourceTencentCloudTagResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag.read")()
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

	tagRes, err := service.DescribeTagResourceById(ctx, tagKey, tagValue)
	if err != nil {
		return err
	}

	if tagRes == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TagTag` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if tagRes.TagKey != nil {
		_ = d.Set("tag_key", tagRes.TagKey)
	}

	if tagRes.TagValue != nil {
		_ = d.Set("tag_value", tagRes.TagValue)
	}

	return nil
}

func resourceTencentCloudTagResourceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tag.delete")()
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

	if err := service.DeleteTagResourceById(ctx, tagKey, tagValue); err != nil {
		return err
	}

	return nil
}
