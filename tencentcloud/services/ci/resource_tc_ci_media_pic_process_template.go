package ci

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCiMediaPicProcessTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiMediaPicProcessTemplateCreate,
		Read:   resourceTencentCloudCiMediaPicProcessTemplateRead,
		Update: resourceTencentCloudCiMediaPicProcessTemplateUpdate,
		Delete: resourceTencentCloudCiMediaPicProcessTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},

			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.",
			},

			"pic_process": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "container format.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"is_pic_info": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to return the original image information.",
						},
						"process_rule": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Image processing rules, 1: basic image processing, please refer to the basic image processing document, 2: image compression, please refer to the image compression document, 3: blind watermark, please refer to the blind watermark document.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCiMediaPicProcessTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_pic_process_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var (
		request = cos.CreateMediaPicProcessTemplateOptions{
			Tag: "PicProcess",
		}
		bucket     string
		templateId string
	)

	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "pic_process"); ok {
		picProcess := cos.PicProcess{}
		if v, ok := dMap["is_pic_info"]; ok {
			picProcess.IsPicInfo = v.(string)
		}
		if v, ok := dMap["process_rule"]; ok {
			picProcess.ProcessRule = v.(string)
		}
		request.PicProcess = &picProcess
	}

	var response *cos.CreateMediaTemplateResult
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.CreateMediaPicProcessTemplate(ctx, &request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%+v]\n", logId, "CreateMediaPicProcessTemplate", request, result)
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaPicProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	templateId = response.Template.TemplateId
	d.SetId(bucket + tccommon.FILED_SP + templateId)

	return resourceTencentCloudCiMediaPicProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaPicProcessTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_pic_process_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	mediaPicProcessTemplate, err := service.DescribeCiMediaTemplateById(ctx, bucket, templateId)
	if err != nil {
		return err
	}

	if mediaPicProcessTemplate == nil {
		d.SetId("")
		return fmt.Errorf("resource `track` %s does not exist", d.Id())
	}

	_ = d.Set("bucket", bucket)

	if mediaPicProcessTemplate.Name != "" {
		_ = d.Set("name", mediaPicProcessTemplate.Name)
	}

	if mediaPicProcessTemplate.PicProcess != nil {
		picProcessMap := map[string]interface{}{}

		if mediaPicProcessTemplate.PicProcess.IsPicInfo != "" {
			picProcessMap["is_pic_info"] = mediaPicProcessTemplate.PicProcess.IsPicInfo
		}

		if mediaPicProcessTemplate.PicProcess.ProcessRule != "" {
			picProcessMap["process_rule"] = mediaPicProcessTemplate.PicProcess.ProcessRule
		}

		err = d.Set("pic_process", []interface{}{picProcessMap})
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceTencentCloudCiMediaPicProcessTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_pic_process_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	request := cos.CreateMediaPicProcessTemplateOptions{
		Tag: "PicProcess",
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = v.(string)
	}

	if d.HasChange("pic_process") {
		if dMap, ok := helper.InterfacesHeadMap(d, "pic_process"); ok {
			picProcess := cos.PicProcess{}
			if v, ok := dMap["is_pic_info"]; ok {
				picProcess.IsPicInfo = v.(string)
			}
			if v, ok := dMap["process_rule"]; ok {
				picProcess.ProcessRule = v.(string)
			}
			request.PicProcess = &picProcess
		}
	}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, _, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCiClient(bucket).CI.UpdateMediaPicProcessTemplate(ctx, &request, templateId)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "UpdateMediaPicProcessTemplate", request, result)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci mediaPicProcessTemplate failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiMediaPicProcessTemplateRead(d, meta)
}

func resourceTencentCloudCiMediaPicProcessTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ci_media_pic_process_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := CiService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	bucket := idSplit[0]
	templateId := idSplit[1]

	if err := service.DeleteCiMediaTemplateById(ctx, bucket, templateId); err != nil {
		return err
	}

	return nil
}
