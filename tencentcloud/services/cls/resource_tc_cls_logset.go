package cls

import (
	"context"
	"fmt"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClsLogset() *schema.Resource {
	return &schema.Resource{
		Read:   resourceTencentCloudClsLogsetRead,
		Create: resourceTencentCloudClsLogsetCreate,
		Update: resourceTencentCloudClsLogsetUpdate,
		Delete: resourceTencentCloudClsLogsetDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"logset_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Logset name, which must be unique.",
			},

			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation time.",
			},

			"topic_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of log topics in logset.",
			},

			"role_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "If assumer_uin is not empty, it indicates the service provider who creates the logset.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudClsLogsetCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request  = cls.NewCreateLogsetRequest()
		response *cls.CreateLogsetResponse
	)

	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().CreateLogset(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls logset failed, reason:%+v", logId, err)
		return err
	}

	logsetId := *response.Response.LogsetId

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::cls:%s:uin/:logset/%s", region, logsetId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}
	d.SetId(logsetId)
	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	logsetId := d.Id()

	logset, err := service.DescribeClsLogset(ctx, logsetId)

	if err != nil {
		return err
	}

	if logset == nil {
		d.SetId("")
		return fmt.Errorf("resource `logset` %s does not exist", logsetId)
	}

	if logset.LogsetName != nil {
		_ = d.Set("logset_name", logset.LogsetName)
	}

	if logset.CreateTime != nil {
		_ = d.Set("create_time", logset.CreateTime)
	}

	if logset.TopicCount != nil {
		_ = d.Set("topic_count", logset.TopicCount)
	}

	if logset.RoleName != nil {
		_ = d.Set("role_name", logset.RoleName)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := &TagService{client: tcClient}
	tags, err := tagService.DescribeResourceTags(ctx, "cls", "logset", tcClient.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudClsLogsetUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cls.NewModifyLogsetRequest()

	request.LogsetId = helper.String(d.Id())

	if d.HasChange("logset_name") {
		if v, ok := d.GetOk("logset_name"); ok {
			request.LogsetName = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClsClient().ModifyLogset(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := &TagService{client: tcClient}
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("cls", "logset", tcClient.Region, d.Id())
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudClsLogsetRead(d, meta)
}

func resourceTencentCloudClsLogsetDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cls_logset.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	logsetId := d.Id()

	if err := service.DeleteClsLogsetById(ctx, logsetId); err != nil {
		return err
	}

	return nil
}
