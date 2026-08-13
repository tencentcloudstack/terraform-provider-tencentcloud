package tag

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTagAttachmentV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTagAttachmentV2Create,
		Read:   resourceTencentCloudTagAttachmentV2Read,
		Update: resourceTencentCloudTagAttachmentV2Update,
		Delete: resourceTencentCloudTagAttachmentV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"tag_key": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Tag key.",
			},

			"tag_value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Tag value.",
			},

			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "[Six-segment description of resources](https://cloud.tencent.com/document/product/598/10606).",
			},
		},
	}
}

func resourceTencentCloudTagAttachmentV2Create(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment_v2.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId        = tccommon.GetLogId(tccommon.ContextNil)
		ctx          = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request      = tag.NewAddResourceTagRequest()
		tagKey       string
		resourceName string
	)

	if v, ok := d.GetOk("tag_key"); ok {
		tagKey = v.(string)
		request.TagKey = helper.String(v.(string))
	}

	if v, ok := d.GetOk("tag_value"); ok {
		request.TagValue = helper.String(v.(string))
	}

	if v, ok := d.GetOk("resource"); ok {
		resourceName = v.(string)
		request.Resource = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTagClient().AddResourceTagWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tag_attachment_v2 failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tag_attachment_v2 failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(strings.Join([]string{tagKey, resourceName}, tccommon.FILED_SP))
	return resourceTencentCloudTagAttachmentV2Read(d, meta)
}

func resourceTencentCloudTagAttachmentV2Read(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment_v2.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	tagKey := idSplit[0]
	resourceName := idSplit[1]

	tagRes, err := service.DescribeTagAttachmentV2ById(ctx, tagKey, resourceName)
	if err != nil {
		return err
	}

	if tagRes == nil {
		log.Printf("[CRUD] tag_attachment_v2 id=%s", d.Id())
		d.SetId("")
		return nil
	}

	if tagRes.TagKey != nil {
		_ = d.Set("tag_key", tagRes.TagKey)
	}

	if tagRes.TagValue != nil {
		_ = d.Set("tag_value", tagRes.TagValue)
	}

	_ = d.Set("resource", resourceName)

	return nil
}

func resourceTencentCloudTagAttachmentV2Update(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment_v2.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	tagKey := idSplit[0]
	resourceName := idSplit[1]

	if d.HasChange("tag_value") {
		request := tag.NewUpdateResourceTagValueRequest()
		request.TagKey = helper.String(tagKey)
		request.Resource = helper.String(resourceName)
		if v, ok := d.GetOk("tag_value"); ok {
			request.TagValue = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTagClient().UpdateResourceTagValueWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Update tag_attachment_v2 failed, Response is nil."))
			}

			return nil
		})

		if err != nil {
			log.Printf("[CRITAL]%s update tag_attachment_v2 failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTagAttachmentV2Read(d, meta)
}

func resourceTencentCloudTagAttachmentV2Delete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tag_attachment_v2.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		request = tag.NewDeleteResourceTagRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	tagKey := idSplit[0]
	resourceName := idSplit[1]

	request.TagKey = helper.String(tagKey)
	request.Resource = helper.String(resourceName)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTagClient().DeleteResourceTagWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tag_attachment_v2 failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete tag_attachment_v2 failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
