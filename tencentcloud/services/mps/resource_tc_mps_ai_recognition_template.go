package mps

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsAiRecognitionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsAiRecognitionTemplateCreate,
		Read:   resourceTencentCloudMpsAiRecognitionTemplateRead,
		Update: resourceTencentCloudMpsAiRecognitionTemplateUpdate,
		Delete: resourceTencentCloudMpsAiRecognitionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai recognition template name, length limit: 64 characters.",
			},

			"comment": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Ai recognition template description information, length limit: 256 characters.",
			},

			"face_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Face recognition control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ai face recognition task switch, optional value:ON/OFF.",
						},
						"score": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Description: "Face recognition filter score, when the recognition result reaches the score above, the recognition result will be returned. The default is 95 points. Value range: 0 - 100.",
						},
						"default_library_label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Default face filter tag, specify the tag of the default face that needs to be returned. If not filled or empty, all default face results will be returned. Label optional value:entertainment, sport, politician.",
						},
						"user_define_library_label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "User-defined face filter tag, specify the tag of the user-defined face that needs to be returned. If not filled or empty, all custom face results will be returned.The maximum number of tags is 100, and the length of each tag is up to 16 characters.",
						},
						"face_library": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Face library selection, optional value:Default, UserDefine, AllDefault value: All, use the system default face library and user-defined face library.",
						},
					},
				},
			},

			"ocr_full_text_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ocr full text control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ocr full text recognition task switch, optional value:ON/OFF.",
						},
					},
				},
			},

			"ocr_words_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Ocr words recognition control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Ocr words recognition task switch, optional value:ON/OFF.",
						},
						"label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
						},
					},
				},
			},

			"asr_full_text_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Asr full text recognition control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Asr full text recognition task switch, optional value:ON/OFF.",
						},
						"subtitle_format": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Generated subtitle file format, if left blank or blank string means no subtitle file will be generated, optional value:vtt: Generate WebVTT subtitle files.",
						},
					},
				},
			},

			"asr_words_configure": {
				Optional:    true,
				Type:        schema.TypeList,
				MaxItems:    1,
				Description: "Asr word recognition control parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Asr word recognition task switch, optional value:ON/OFF.",
						},
						"label_set": {
							Type: schema.TypeSet,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Optional:    true,
							Description: "Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudMpsAiRecognitionTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		request    = mps.NewCreateAIRecognitionTemplateRequest()
		response   = mps.NewCreateAIRecognitionTemplateResponse()
		definition int64
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("comment"); ok {
		request.Comment = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "face_configure"); ok {
		faceConfigureInfo := mps.FaceConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			faceConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["score"]; ok {
			faceConfigureInfo.Score = helper.Float64(v.(float64))
		}
		if v, ok := dMap["default_library_label_set"]; ok {
			defaultLibraryLabelSetSet := v.(*schema.Set).List()
			for i := range defaultLibraryLabelSetSet {
				defaultLibraryLabelSet := defaultLibraryLabelSetSet[i].(string)
				faceConfigureInfo.DefaultLibraryLabelSet = append(faceConfigureInfo.DefaultLibraryLabelSet, &defaultLibraryLabelSet)
			}
		}
		if v, ok := dMap["user_define_library_label_set"]; ok {
			userDefineLibraryLabelSetSet := v.(*schema.Set).List()
			for i := range userDefineLibraryLabelSetSet {
				userDefineLibraryLabelSet := userDefineLibraryLabelSetSet[i].(string)
				faceConfigureInfo.UserDefineLibraryLabelSet = append(faceConfigureInfo.UserDefineLibraryLabelSet, &userDefineLibraryLabelSet)
			}
		}
		if v, ok := dMap["face_library"]; ok {
			faceConfigureInfo.FaceLibrary = helper.String(v.(string))
		}
		request.FaceConfigure = &faceConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ocr_full_text_configure"); ok {
		ocrFullTextConfigureInfo := mps.OcrFullTextConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			ocrFullTextConfigureInfo.Switch = helper.String(v.(string))
		}
		request.OcrFullTextConfigure = &ocrFullTextConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "ocr_words_configure"); ok {
		ocrWordsConfigureInfo := mps.OcrWordsConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			ocrWordsConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["label_set"]; ok {
			labelSetSet := v.(*schema.Set).List()
			for i := range labelSetSet {
				labelSet := labelSetSet[i].(string)
				ocrWordsConfigureInfo.LabelSet = append(ocrWordsConfigureInfo.LabelSet, &labelSet)
			}
		}
		request.OcrWordsConfigure = &ocrWordsConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "asr_full_text_configure"); ok {
		asrFullTextConfigureInfo := mps.AsrFullTextConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			asrFullTextConfigureInfo.Switch = helper.String(v.(string))
		}
		if v := dMap["subtitle_format"]; v != "" {
			asrFullTextConfigureInfo.SubtitleFormat = helper.String(v.(string))
		}
		request.AsrFullTextConfigure = &asrFullTextConfigureInfo
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "asr_words_configure"); ok {
		asrWordsConfigureInfo := mps.AsrWordsConfigureInfo{}
		if v, ok := dMap["switch"]; ok {
			asrWordsConfigureInfo.Switch = helper.String(v.(string))
		}
		if v, ok := dMap["label_set"]; ok {
			labelSetSet := v.(*schema.Set).List()
			for i := range labelSetSet {
				labelSet := labelSetSet[i].(string)
				asrWordsConfigureInfo.LabelSet = append(asrWordsConfigureInfo.LabelSet, &labelSet)
			}
		}
		request.AsrWordsConfigure = &asrWordsConfigureInfo
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().CreateAIRecognitionTemplate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps aiRecognitionTemplate failed, reason:%+v", logId, err)
		return err
	}

	definition = *response.Response.Definition
	d.SetId(helper.Int64ToStr(definition))

	return resourceTencentCloudMpsAiRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiRecognitionTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	definition := d.Id()

	aiRecognitionTemplate, err := service.DescribeMpsAiRecognitionTemplateById(ctx, definition)
	if err != nil {
		return err
	}

	if aiRecognitionTemplate == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsAiRecognitionTemplate` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if aiRecognitionTemplate.Name != nil {
		_ = d.Set("name", aiRecognitionTemplate.Name)
	}

	if aiRecognitionTemplate.Comment != nil {
		_ = d.Set("comment", aiRecognitionTemplate.Comment)
	}

	if aiRecognitionTemplate.FaceConfigure != nil {
		faceConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.FaceConfigure.Switch != nil {
			faceConfigureMap["switch"] = aiRecognitionTemplate.FaceConfigure.Switch
		}

		if aiRecognitionTemplate.FaceConfigure.Score != nil {
			faceConfigureMap["score"] = aiRecognitionTemplate.FaceConfigure.Score
		}

		if aiRecognitionTemplate.FaceConfigure.DefaultLibraryLabelSet != nil {
			faceConfigureMap["default_library_label_set"] = aiRecognitionTemplate.FaceConfigure.DefaultLibraryLabelSet
		}

		if aiRecognitionTemplate.FaceConfigure.UserDefineLibraryLabelSet != nil {
			faceConfigureMap["user_define_library_label_set"] = aiRecognitionTemplate.FaceConfigure.UserDefineLibraryLabelSet
		}

		if aiRecognitionTemplate.FaceConfigure.FaceLibrary != nil {
			faceConfigureMap["face_library"] = aiRecognitionTemplate.FaceConfigure.FaceLibrary
		}

		_ = d.Set("face_configure", []interface{}{faceConfigureMap})
	}

	if aiRecognitionTemplate.OcrFullTextConfigure != nil {
		ocrFullTextConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.OcrFullTextConfigure.Switch != nil {
			ocrFullTextConfigureMap["switch"] = aiRecognitionTemplate.OcrFullTextConfigure.Switch
		}

		_ = d.Set("ocr_full_text_configure", []interface{}{ocrFullTextConfigureMap})
	}

	if aiRecognitionTemplate.OcrWordsConfigure != nil {
		ocrWordsConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.OcrWordsConfigure.Switch != nil {
			ocrWordsConfigureMap["switch"] = aiRecognitionTemplate.OcrWordsConfigure.Switch
		}

		if aiRecognitionTemplate.OcrWordsConfigure.LabelSet != nil {
			ocrWordsConfigureMap["label_set"] = aiRecognitionTemplate.OcrWordsConfigure.LabelSet
		}

		_ = d.Set("ocr_words_configure", []interface{}{ocrWordsConfigureMap})
	}

	if aiRecognitionTemplate.AsrFullTextConfigure != nil {
		asrFullTextConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.AsrFullTextConfigure.Switch != nil {
			asrFullTextConfigureMap["switch"] = aiRecognitionTemplate.AsrFullTextConfigure.Switch
		}

		if aiRecognitionTemplate.AsrFullTextConfigure.SubtitleFormat != nil {
			asrFullTextConfigureMap["subtitle_format"] = aiRecognitionTemplate.AsrFullTextConfigure.SubtitleFormat
		}

		_ = d.Set("asr_full_text_configure", []interface{}{asrFullTextConfigureMap})
	}

	if aiRecognitionTemplate.AsrWordsConfigure != nil {
		asrWordsConfigureMap := map[string]interface{}{}

		if aiRecognitionTemplate.AsrWordsConfigure.Switch != nil {
			asrWordsConfigureMap["switch"] = aiRecognitionTemplate.AsrWordsConfigure.Switch
		}

		if aiRecognitionTemplate.AsrWordsConfigure.LabelSet != nil {
			asrWordsConfigureMap["label_set"] = aiRecognitionTemplate.AsrWordsConfigure.LabelSet
		}

		_ = d.Set("asr_words_configure", []interface{}{asrWordsConfigureMap})
	}

	return nil
}

func resourceTencentCloudMpsAiRecognitionTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := mps.NewModifyAIRecognitionTemplateRequest()

	definition := d.Id()

	request.Definition = helper.StrToInt64Point(definition)

	mutableArgs := []string{"name", "comment", "face_configure", "ocr_full_text_configure", "ocr_words_configure", "asr_full_text_configure", "asr_words_configure"}

	needChange := false

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("comment"); ok {
			request.Comment = helper.String(v.(string))
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "face_configure"); ok {
			faceConfigureInfo := mps.FaceConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				faceConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["score"]; ok {
				faceConfigureInfo.Score = helper.Float64(v.(float64))
			}
			if v, ok := dMap["default_library_label_set"]; ok {
				defaultLibraryLabelSetSet := v.(*schema.Set).List()
				for i := range defaultLibraryLabelSetSet {
					defaultLibraryLabelSet := defaultLibraryLabelSetSet[i].(string)
					faceConfigureInfo.DefaultLibraryLabelSet = append(faceConfigureInfo.DefaultLibraryLabelSet, &defaultLibraryLabelSet)
				}
			}
			if v, ok := dMap["user_define_library_label_set"]; ok {
				userDefineLibraryLabelSetSet := v.(*schema.Set).List()
				for i := range userDefineLibraryLabelSetSet {
					userDefineLibraryLabelSet := userDefineLibraryLabelSetSet[i].(string)
					faceConfigureInfo.UserDefineLibraryLabelSet = append(faceConfigureInfo.UserDefineLibraryLabelSet, &userDefineLibraryLabelSet)
				}
			}
			if v, ok := dMap["face_library"]; ok {
				faceConfigureInfo.FaceLibrary = helper.String(v.(string))
			}
			request.FaceConfigure = &faceConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ocr_full_text_configure"); ok {
			ocrFullTextConfigureInfo := mps.OcrFullTextConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				ocrFullTextConfigureInfo.Switch = helper.String(v.(string))
			}
			request.OcrFullTextConfigure = &ocrFullTextConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "ocr_words_configure"); ok {
			ocrWordsConfigureInfo := mps.OcrWordsConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				ocrWordsConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					ocrWordsConfigureInfo.LabelSet = append(ocrWordsConfigureInfo.LabelSet, &labelSet)
				}
			}
			request.OcrWordsConfigure = &ocrWordsConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "asr_full_text_configure"); ok {
			asrFullTextConfigureInfo := mps.AsrFullTextConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				asrFullTextConfigureInfo.Switch = helper.String(v.(string))
			}
			if v := dMap["subtitle_format"]; v != "" {
				asrFullTextConfigureInfo.SubtitleFormat = helper.String(v.(string))
			}
			request.AsrFullTextConfigure = &asrFullTextConfigureInfo
		}

		if dMap, ok := helper.InterfacesHeadMap(d, "asr_words_configure"); ok {
			asrWordsConfigureInfo := mps.AsrWordsConfigureInfoForUpdate{}
			if v, ok := dMap["switch"]; ok {
				asrWordsConfigureInfo.Switch = helper.String(v.(string))
			}
			if v, ok := dMap["label_set"]; ok {
				labelSetSet := v.(*schema.Set).List()
				for i := range labelSetSet {
					labelSet := labelSetSet[i].(string)
					asrWordsConfigureInfo.LabelSet = append(asrWordsConfigureInfo.LabelSet, &labelSet)
				}
			}
			request.AsrWordsConfigure = &asrWordsConfigureInfo
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().ModifyAIRecognitionTemplate(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update mps aiRecognitionTemplate failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudMpsAiRecognitionTemplateRead(d, meta)
}

func resourceTencentCloudMpsAiRecognitionTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_ai_recognition_template.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := MpsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	definition := d.Id()

	if err := service.DeleteMpsAiRecognitionTemplateById(ctx, definition); err != nil {
		return err
	}

	return nil
}
