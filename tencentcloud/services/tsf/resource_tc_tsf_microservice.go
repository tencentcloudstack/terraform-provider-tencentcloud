package tsf

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tsf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tsf/v20180326"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTsfMicroservice() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTsfMicroserviceCreate,
		Read:   resourceTencentCloudTsfMicroserviceRead,
		Update: resourceTencentCloudTsfMicroserviceUpdate,
		Delete: resourceTencentCloudTsfMicroserviceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Namespace ID.",
			},

			"microservice_name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Microservice name.",
			},

			"microservice_desc": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Microservice description information.",
			},

			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tag description list.",
			},
		},
	}
}

func resourceTencentCloudTsfMicroserviceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_microservice.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request          = tsf.NewCreateMicroserviceRequest()
		response         = tsf.NewCreateMicroserviceResponse()
		microserviceName string
		namespaceId      string
		microserviceId   string
	)
	if v, ok := d.GetOk("namespace_id"); ok {
		namespaceId = v.(string)
		request.NamespaceId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_name"); ok {
		microserviceName = v.(string)
		request.MicroserviceName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("microservice_desc"); ok {
		request.MicroserviceDesc = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().CreateMicroservice(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create tsf microservice failed, reason:%+v", logId, err)
		return err
	}

	if *response.Response.Result {
		service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		microservice, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, "", microserviceName)
		if err != nil {
			return err
		}

		microserviceId = *microservice.MicroserviceId
		d.SetId(namespaceId + tccommon.FILED_SP + microserviceId)
	} else {
		return fmt.Errorf("[DEBUG]%s api[%s] Creation failed, and the return result of interface creation is false", logId, request.GetAction())
	}

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := fmt.Sprintf("qcs::tsf:%s:uin/:microservice/%s", region, microserviceId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfMicroserviceRead(d, meta)
}

func resourceTencentCloudTsfMicroserviceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_microservice.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	microservice, err := service.DescribeTsfMicroserviceById(ctx, namespaceId, microserviceId, "")
	if err != nil {
		return err
	}

	if microservice == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `TsfMicroservice` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}
	if microservice.NamespaceId != nil {
		_ = d.Set("namespace_id", microservice.NamespaceId)
	}

	if microservice.MicroserviceName != nil {
		_ = d.Set("microservice_name", microservice.MicroserviceName)
	}

	if microservice.MicroserviceDesc != nil {
		_ = d.Set("microservice_desc", microservice.MicroserviceDesc)
	}

	tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(tcClient)
	tags, err := tagService.DescribeResourceTags(ctx, "tsf", "microservice", tcClient.Region, microserviceId)
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudTsfMicroserviceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_microservice.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := tsf.NewModifyMicroserviceRequest()

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	request.MicroserviceId = &microserviceId

	immutableArgs := []string{"namespace_id", "microservice_name"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("microservice_desc") {
		if v, ok := d.GetOk("microservice_desc"); ok {
			request.MicroserviceDesc = helper.String(v.(string))
		}
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTsfClient().ModifyMicroservice(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update tsf microservice failed, reason:%+v", logId, err)
		return err
	}

	if d.HasChange("tags") {
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		oldTags, newTags := d.GetChange("tags")
		replaceTags, deleteTags := svctag.DiffTags(oldTags.(map[string]interface{}), newTags.(map[string]interface{}))
		resourceName := tccommon.BuildTagResourceName("tsf", "microservice", tcClient.Region, microserviceId)
		if err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags); err != nil {
			return err
		}
	}

	return resourceTencentCloudTsfMicroserviceRead(d, meta)
}

func resourceTencentCloudTsfMicroserviceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_tsf_microservice.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := TsfService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	// namespaceId := idSplit[0]
	microserviceId := idSplit[1]

	if err := service.DeleteTsfMicroserviceById(ctx, microserviceId); err != nil {
		return err
	}

	return nil
}
