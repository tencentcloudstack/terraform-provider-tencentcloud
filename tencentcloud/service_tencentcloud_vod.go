package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	vod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vod/v20180717"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type VodService struct {
	client *connectivity.TencentCloudClient
}

func (me *VodService) DescribeAdaptiveDynamicStreamingTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.AdaptiveDynamicStreamingTemplate, errRet error) {
	logId := getLogId(ctx)
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeAdaptiveDynamicStreamingTemplates(request)
			if err != nil {
				return retryError(err)
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

func (me *VodService) DescribeAdaptiveDynamicStreamingTemplatesById(ctx context.Context, templateId string) (templateInfo *vod.AdaptiveDynamicStreamingTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
		}
	)

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
	logId := getLogId(ctx)
	request := vod.NewDeleteAdaptiveDynamicStreamingTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteAdaptiveDynamicStreamingTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeProcedureTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.ProcedureTemplate, errRet error) {
	logId := getLogId(ctx)
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeProcedureTemplates(request)
			if err != nil {
				return retryError(err)
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

func (me *VodService) DescribeProcedureTemplatesById(ctx context.Context, templateId string) (templateInfo *vod.ProcedureTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"name": []string{templateId},
		}
	)

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
	logId := getLogId(ctx)
	request := vod.NewDeleteProcedureTemplateRequest()

	request.Name = &templateId
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteProcedureTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeSnapshotByTimeOffsetTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.SnapshotByTimeOffsetTemplate, errRet error) {
	logId := getLogId(ctx)
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeSnapshotByTimeOffsetTemplates(request)
			if err != nil {
				return retryError(err)
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

func (me *VodService) DescribeSnapshotByTimeOffsetTemplatesById(ctx context.Context, templateId string) (templateInfo *vod.SnapshotByTimeOffsetTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
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
	logId := getLogId(ctx)
	request := vod.NewDeleteSnapshotByTimeOffsetTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteSnapshotByTimeOffsetTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeImageSpriteTemplatesByFilter(ctx context.Context, filters map[string]interface{}) (templates []*vod.ImageSpriteTemplate, errRet error) {
	logId := getLogId(ctx)
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeImageSpriteTemplates(request)
			if err != nil {
				return retryError(err)
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

func (me *VodService) DescribeImageSpriteTemplatesById(ctx context.Context, templateId string) (templateInfo *vod.ImageSpriteTemplate, has bool, errRet error) {
	var (
		filter = map[string]interface{}{
			"definitions": []string{templateId},
		}
	)

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
	logId := getLogId(ctx)
	request := vod.NewDeleteImageSpriteTemplateRequest()

	idUint, _ := strconv.ParseUint(templateId, 0, 64)
	request.Definition = helper.Uint64(idUint)
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteImageSpriteTemplate(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}

func (me *VodService) DescribeSuperPlayerConfigsByFilter(ctx context.Context, filters map[string]interface{}) (configs []*vod.PlayerConfig, errRet error) {
	logId := getLogId(ctx)
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
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseVodClient().DescribeSuperPlayerConfigs(request)
			if err != nil {
				return retryError(err)
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
	logId := getLogId(ctx)
	request := vod.NewDeleteSuperPlayerConfigRequest()

	request.Name = &configId
	if subAppid != 0 {
		request.SubAppId = &subAppid
	}

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, errRet = me.client.UseVodClient().DeleteSuperPlayerConfig(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}

	return
}
