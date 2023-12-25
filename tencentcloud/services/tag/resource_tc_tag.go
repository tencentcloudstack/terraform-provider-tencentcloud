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

func ResourceTencentCloudTag() *schema.Resource {
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
	defer tccommon.LogElapsed("resource.tencentcloud_tag.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTagClient().CreateTag(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tag tag failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(tagKey + tccommon.FILED_SP + tagValue)

	return resourceTencentCloudTagResourceRead(d, meta)
}

func resourceTencentCloudTagResourceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tag.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tag.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
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
