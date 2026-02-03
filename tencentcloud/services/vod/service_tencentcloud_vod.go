package vod

import (
	"context"
	"fmt"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewVodService(client *connectivity.TencentCloudClient) VodService {
	return VodService{client: client}
}

type VodService struct {
	client *connectivity.TencentCloudClient
}

func (me *VodService) DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.AdaptiveDynamicStreamingTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeAdaptiveDynamicStreamingTemplatesRequest()

	offset := VOD_DEFAULT_OFFSET
	limit := VOD_MAX_LIMIT
	request.Offset = helper.IntUint64(offset)
	request.Limit = helper.IntUint64(limit)
	if v, ok := filters["type"]; ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := filters["sub_appid"]; ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if v, ok := filters["definitions"]; ok {
		for _, vv := range v.([]string) {
			idUint, _ := strconv.ParseUint(vv, 0, 64)
			request.Definitions = append(request.Definitions, &idUint)
		}
	}

	templates = make([]*vod.AdaptiveDynamicStreamingTemplate, 0)
	for {
		var response *vod.DescribeAdaptiveDynamicStreamingTemplatesResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeAdaptiveDynamicStreamingTemplates(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			templates = append(templates, response.Response.AdaptiveDynamicStreamingTemplateSet...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		if len(response.Response.AdaptiveDynamicStreamingTemplateSet) < VOD_MAX_LIMIT {
			break
		} else {
			offset += limit
		}
	}
	return
}

func (me *VodService) DescribeAdaptiveDynamicStreamingTemplatesById(ctx context.Context, templateId string, subAppId int) (templateInfo *vod.AdaptiveDynamicStreamingTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
		}
	)

	if subAppId != 0 {
		filter["sub_appid"] = subAppId
	}

	templates, errRet := me.DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx, filter)
	if errRet != nil {
		return
	}
	if len(templates) == 0 {
		return
	}
	if len(templates) != 1 {
		errRet = fmt.Errorf("dumplicate template found by id %s", templateId)
		return
	}

	has = true
	templateInfo = templates[0]
	return
}

func (me *VodService) DeleteAdaptiveDynamicStreamingTemplate(ctx context.Context, templateId string, subAppid uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDeleteAdaptiveDynamicStreamingTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteAdaptiveDynamicStreamingTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeProcedureTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.ProcedureTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeProcedureTemplatesRequest()

	offset := VOD_DEFAULT_OFFSET
	limit := VOD_MAX_LIMIT
	request.Offset = helper.IntUint64(offset)
	request.Limit = helper.IntUint64(limit)
	if v, ok := filters["type"]; ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := filters["sub_appid"]; ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if v, ok := filters["name"]; ok {
		for _, vv := range v.([]string) {
			request.Names = append(request.Names, &vv)
		}
	}

	templates = make([]*vod.ProcedureTemplate, 0)
	for {
		var response *vod.DescribeProcedureTemplatesResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeProcedureTemplates(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			templates = append(templates, response.Response.ProcedureTemplateSet...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		if len(response.Response.ProcedureTemplateSet) < VOD_MAX_LIMIT {
			break
		} else {
			offset += limit
		}
	}
	return
}

func (me *VodService) DescribeProcedureTemplatesById(ctx context.Context, templateId string, subAppId int) (templateInfo *vod.ProcedureTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"name": []string{templateId},
		}
	)
	if subAppId != 0 {
		filter["sub_appid"] = subAppId
	}

	templates, errRet := me.DescribeProcedureTemplatesByFilter(ctx, filter)
	if errRet != nil {
		return
	}
	if len(templates) == 0 {
		return
	}
	if len(templates) != 1 {
		errRet = fmt.Errorf("dumplicate template found by id %s", templateId)
		return
	}

	has = true
	templateInfo = templates[0]
	return
}

func (me *VodService) DeleteProcedureTemplate(ctx context.Context, templateId string, subAppid uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDeleteProcedureTemplateRequest()

	request.Name = &templateId
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteProcedureTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.SnapshotByTimeOffsetTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeSnapshotByTimeOffsetTemplatesRequest()

	offset := VOD_DEFAULT_OFFSET
	limit := VOD_MAX_LIMIT
	request.Offset = helper.IntUint64(offset)
	request.Limit = helper.IntUint64(limit)
	if v, ok := filters["type"]; ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := filters["sub_appid"]; ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if v, ok := filters["definitions"]; ok {
		for _, vv := range v.([]string) {
			idUint, _ := strconv.ParseUint(vv, 0, 64)
			request.Definitions = append(request.Definitions, &idUint)
		}
	}

	templates = make([]*vod.SnapshotByTimeOffsetTemplate, 0)
	for {
		var response *vod.DescribeSnapshotByTimeOffsetTemplatesResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeSnapshotByTimeOffsetTemplates(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			templates = append(templates, response.Response.SnapshotByTimeOffsetTemplateSet...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		if len(response.Response.SnapshotByTimeOffsetTemplateSet) < VOD_MAX_LIMIT {
			break
		} else {
			offset += limit
		}
	}
	return
}

func (me *VodService) DescribeSnapshotByTimeOffsetTemplatesById(ctx context.Context, templateId string, subAppId int) (templateInfo *vod.SnapshotByTimeOffsetTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
			"sub_appid":   subAppId,
		}
	)

	templates, errRet := me.DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx, filter)
	if errRet != nil {
		return
	}
	if len(templates) == 0 {
		return
	}
	if len(templates) != 1 {
		errRet = fmt.Errorf("dumplicate template found by id %s", templateId)
		return
	}

	has = true
	templateInfo = templates[0]
	return
}

func (me *VodService) DeleteSnapshotByTimeOffsetTemplate(ctx context.Context, templateId string, subAppid uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDeleteSnapshotByTimeOffsetTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteSnapshotByTimeOffsetTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeImageSpriteTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.ImageSpriteTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeImageSpriteTemplatesRequest()

	offset := VOD_DEFAULT_OFFSET
	limit := VOD_MAX_LIMIT
	request.Offset = helper.IntUint64(offset)
	request.Limit = helper.IntUint64(limit)
	if v, ok := filters["type"]; ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := filters["sub_appid"]; ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if v, ok := filters["definitions"]; ok {
		for _, vv := range v.([]string) {
			idUint, _ := strconv.ParseUint(vv, 0, 64)
			request.Definitions = append(request.Definitions, &idUint)
		}
	}

	templates = make([]*vod.ImageSpriteTemplate, 0)
	for {
		var response *vod.DescribeImageSpriteTemplatesResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeImageSpriteTemplates(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			templates = append(templates, response.Response.ImageSpriteTemplateSet...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		if len(response.Response.ImageSpriteTemplateSet) < VOD_MAX_LIMIT {
			break
		} else {
			offset += limit
		}
	}
	return
}

func (me *VodService) DescribeImageSpriteTemplatesById(ctx context.Context, templateId string, subAppId int) (templateInfo *vod.ImageSpriteTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
		}
	)

	if subAppId != 0 {
		filter["sub_appid"] = subAppId
	}
	templates, errRet := me.DescribeImageSpriteTemplatesByFilter(ctx, filter)
	if errRet != nil {
		return
	}
	if len(templates) == 0 {
		return
	}
	if len(templates) != 1 {
		errRet = fmt.Errorf("dumplicate template found by id %s", templateId)
		return
	}

	has = true
	templateInfo = templates[0]
	return
}

func (me *VodService) DeleteImageSpriteTemplate(ctx context.Context, templateId string, subAppid uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDeleteImageSpriteTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteImageSpriteTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeSuperPlayerConfigsByFilter(ctx context.Context, filters map[string]interface{}) (configs []*vod.PlayerConfig, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeSuperPlayerConfigsRequest()

	offset := VOD_DEFAULT_OFFSET
	limit := VOD_MAX_LIMIT
	request.Offset = helper.IntUint64(offset)
	request.Limit = helper.IntUint64(limit)
	if v, ok := filters["type"]; ok {
		request.Type = helper.String(v.(string))
	}
	if v, ok := filters["sub_appid"]; ok {
		request.SubAppId = helper.IntUint64(v.(int))
	}
	if v, ok := filters["names"]; ok {
		for _, vv := range v.([]string) {
			request.Names = append(request.Names, &vv)
		}
	}

	configs = make([]*vod.PlayerConfig, 0)
	for {
		var response *vod.DescribeSuperPlayerConfigsResponse
		var err error
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeSuperPlayerConfigs(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			configs = append(configs, response.Response.PlayerConfigSet...)
			return nil
		})
		if err != nil {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]", logId, request.GetAction(), request.ToJsonString(), err.Error())
			return
		}
		if len(response.Response.PlayerConfigSet) < VOD_MAX_LIMIT {
			break
		} else {
			offset += limit
		}
	}
	return
}

func (me *VodService) DescribeSuperPlayerConfigsById(ctx context.Context, configId string) (configInfo *vod.PlayerConfig, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"names": []string{configId},
		}
	)

	configs, errRet := me.DescribeSuperPlayerConfigsByFilter(ctx, filter)
	if errRet != nil {
		return
	}
	if len(configs) == 0 {
		return
	}
	if len(configs) != 1 {
		errRet = fmt.Errorf("dumplicate configs found by id %s", configId)
		return
	}

	has = true
	configInfo = configs[0]
	return
}

func (me *VodService) DeleteSuperPlayerConfig(ctx context.Context, configId string, subAppid uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDeleteSuperPlayerConfigRequest()

	request.Name = &configId
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteSuperPlayerConfig(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return tccommon.RetryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeVodSampleSnapshotTemplateById(ctx context.Context, subAppId, definition uint64) (sampleSnapshotTemplate *vod.SampleSnapshotTemplate, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vod.NewDescribeSampleSnapshotTemplatesRequest()
	if subAppId != 0 {
		request.SubAppId = helper.Uint64(subAppId)
	}
	request.Definitions = []*uint64{helper.Uint64(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DescribeSampleSnapshotTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SampleSnapshotTemplateSet) < 1 {
		return
	}

	sampleSnapshotTemplate = response.Response.SampleSnapshotTemplateSet[0]
	return
}

func (me *VodService) DeleteVodSampleSnapshotTemplateById(ctx context.Context, subAppId, definition uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := vod.NewDeleteSampleSnapshotTemplateRequest()
	request.SubAppId = helper.Uint64(subAppId)
	request.Definition = helper.Uint64(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DeleteSampleSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VodService) DescribeVodTranscodeTemplateById(ctx context.Context, subAppId uint64, definition int64) (transcodeTemplate *vod.TranscodeTemplate, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewDescribeTranscodeTemplatesRequest()
	request.SubAppId = helper.Uint64(subAppId)
	request.Definitions = []*int64{helper.Int64(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DescribeTranscodeTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.TranscodeTemplateSet) < 1 {
		return
	}

	transcodeTemplate = response.Response.TranscodeTemplateSet[0]
	return
}

func (me *VodService) DeleteVodTranscodeTemplateById(ctx context.Context, subAppId uint64, definition int64) (errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewDeleteTranscodeTemplateRequest()
	request.SubAppId = helper.Uint64(subAppId)
	request.Definition = helper.Int64(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DeleteTranscodeTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VodService) DescribeVodWatermarkTemplateById(ctx context.Context, subAppId uint64, definition int64) (watermarkTemplate *vod.WatermarkTemplate, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewDescribeWatermarkTemplatesRequest()
	request.SubAppId = helper.Uint64(subAppId)
	request.Definitions = []*int64{helper.Int64(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DescribeWatermarkTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WatermarkTemplateSet) < 1 {
		return
	}

	watermarkTemplate = response.Response.WatermarkTemplateSet[0]
	return
}

func (me *VodService) DeleteVodWatermarkTemplateById(ctx context.Context, subAppId uint64, definition int64) (errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewDeleteWatermarkTemplateRequest()
	request.SubAppId = helper.Uint64(subAppId)
	request.Definition = helper.Int64(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DeleteWatermarkTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *VodService) DescribeVodEventConfig(ctx context.Context, subAppId uint64) (eventConfig *vod.DescribeEventConfigResponseParams, errRet error) {
	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := vod.NewDescribeEventConfigRequest()
	request.SubAppId = &subAppId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseVodClient().DescribeEventConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	eventConfig = response.Response
	return
}

func (me *VodService) DescribeSubApplicationsByFilter(ctx context.Context, filters map[string]interface{}) (subAppInfos []*vod.SubAppIdInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := vod.NewDescribeSubAppIdsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	// Set filter parameters
	if v, ok := filters["name"]; ok {
		request.Name = helper.String(v.(string))
	}
	if v, ok := filters["tags"]; ok {
		tags := v.([]*vod.ResourceTag)
		request.Tags = tags
	}

	// Initialize pagination
	var offset uint64 = 0
	var limit uint64 = 200
	if v, ok := filters["offset"]; ok {
		offset = uint64(v.(int))
	}
	if v, ok := filters["limit"]; ok {
		limit = uint64(v.(int))
	}

	subAppInfos = make([]*vod.SubAppIdInfo, 0)

	// Pagination loop
	for {
		request.Offset = &offset
		request.Limit = &limit

		var response *vod.DescribeSubAppIdsResponse
		var err error

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeSubAppIds(request)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.SubAppIdInfoSet) == 0 {
			break
		}

		subAppInfos = append(subAppInfos, response.Response.SubAppIdInfoSet...)

		// Check if we've retrieved all results
		if response.Response.TotalCount == nil || uint64(len(subAppInfos)) >= *response.Response.TotalCount {
			break
		}

		// Move to next page
		offset += limit
	}

	return
}
