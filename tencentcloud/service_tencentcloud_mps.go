package tencentcloud

import (
	"context"
	"log"

	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MpsService struct {
	client *connectivity.TencentCloudClient
}

func (me *MpsService) DescribeMpsWorkflowById(ctx context.Context, workflowId string) (workflow *mps.WorkflowInfo, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeWorkflowsRequest()
	request.WorkflowIds = []*int64{helper.Int64(helper.StrToInt64(workflowId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeWorkflows(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WorkflowInfoSet) < 1 {
		return
	}

	workflow = response.Response.WorkflowInfoSet[0]
	return
}

func (me *MpsService) DeleteMpsWorkflowById(ctx context.Context, workflowId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteWorkflowRequest()
	request.WorkflowId = helper.Int64(helper.StrToInt64(workflowId))

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) EnableWorkflow(ctx context.Context, workflowId int64) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewEnableWorkflowRequest()
	request.WorkflowId = helper.Int64(workflowId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().EnableWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DisableWorkflow(ctx context.Context, workflowId int64) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDisableWorkflowRequest()
	request.WorkflowId = helper.Int64(workflowId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DisableWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsTranscodeTemplateById(ctx context.Context, definition string) (transcodeTemplate *mps.TranscodeTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeTranscodeTemplatesRequest()
	request.Definitions = []*int64{helper.StrToInt64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeTranscodeTemplates(request)
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

func (me *MpsService) DeleteMpsTranscodeTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteTranscodeTemplateRequest()
	request.Definition = helper.StrToInt64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteTranscodeTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsWatermarkTemplateById(ctx context.Context, definition string) (watermarkTemplate *mps.WatermarkTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeWatermarkTemplatesRequest()
	request.Definitions = []*int64{helper.StrToInt64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeWatermarkTemplates(request)
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

func (me *MpsService) DeleteMpsWatermarkTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteWatermarkTemplateRequest()
	request.Definition = helper.StrToInt64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteWatermarkTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsImageSpriteTemplateById(ctx context.Context, definition string) (imageSpriteTemplate *mps.ImageSpriteTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeImageSpriteTemplatesRequest()
	request.Definitions = []*uint64{helper.StrToUint64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeImageSpriteTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ImageSpriteTemplateSet) < 1 {
		return
	}

	imageSpriteTemplate = response.Response.ImageSpriteTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsImageSpriteTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteImageSpriteTemplateRequest()
	request.Definition = helper.StrToUint64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteImageSpriteTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsSnapshotByTimeoffsetTemplateById(ctx context.Context, definition string) (snapshotByTimeoffsetTemplate *mps.SnapshotByTimeOffsetTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeSnapshotByTimeOffsetTemplatesRequest()
	request.Definitions = []*uint64{helper.StrToUint64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeSnapshotByTimeOffsetTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SnapshotByTimeOffsetTemplateSet) < 1 {
		return
	}

	snapshotByTimeoffsetTemplate = response.Response.SnapshotByTimeOffsetTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsSnapshotByTimeoffsetTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteSnapshotByTimeOffsetTemplateRequest()
	request.Definition = helper.StrToUint64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteSnapshotByTimeOffsetTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsSampleSnapshotTemplateById(ctx context.Context, definition string) (sampleSnapshotTemplate *mps.SampleSnapshotTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeSampleSnapshotTemplatesRequest()
	request.Definitions = []*uint64{helper.StrToUint64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeSampleSnapshotTemplates(request)
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

func (me *MpsService) DeleteMpsSampleSnapshotTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteSampleSnapshotTemplateRequest()
	request.Definition = helper.StrToUint64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteSampleSnapshotTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsAnimatedGraphicsTemplateById(ctx context.Context, definition string) (animatedGraphicsTemplate *mps.AnimatedGraphicsTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeAnimatedGraphicsTemplatesRequest()
	request.Definitions = []*uint64{helper.StrToUint64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeAnimatedGraphicsTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AnimatedGraphicsTemplateSet) < 1 {
		return
	}

	animatedGraphicsTemplate = response.Response.AnimatedGraphicsTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsAnimatedGraphicsTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteAnimatedGraphicsTemplateRequest()
	request.Definition = helper.StrToUint64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteAnimatedGraphicsTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsAdaptiveDynamicStreamingTemplateById(ctx context.Context, definition string) (adaptiveDynamicStreamingTemplate *mps.AdaptiveDynamicStreamingTemplate, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeAdaptiveDynamicStreamingTemplatesRequest()
	request.Definitions = []*uint64{helper.StrToUint64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeAdaptiveDynamicStreamingTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AdaptiveDynamicStreamingTemplateSet) < 1 {
		return
	}

	adaptiveDynamicStreamingTemplate = response.Response.AdaptiveDynamicStreamingTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsAdaptiveDynamicStreamingTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteAdaptiveDynamicStreamingTemplateRequest()
	request.Definition = helper.StrToUint64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteAdaptiveDynamicStreamingTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsAiAnalysisTemplateById(ctx context.Context, definition string) (aiAnalysisTemplate *mps.AIAnalysisTemplateItem, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeAIAnalysisTemplatesRequest()
	request.Definitions = []*int64{helper.StrToInt64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeAIAnalysisTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AIAnalysisTemplateSet) < 1 {
		return
	}

	aiAnalysisTemplate = response.Response.AIAnalysisTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsAiAnalysisTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteAIAnalysisTemplateRequest()
	request.Definition = helper.StrToInt64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteAIAnalysisTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsAiRecognitionTemplateById(ctx context.Context, definition string) (aiRecognitionTemplate *mps.AIRecognitionTemplateItem, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeAIRecognitionTemplatesRequest()
	request.Definitions = []*int64{helper.StrToInt64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeAIRecognitionTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.AIRecognitionTemplateSet) < 1 {
		return
	}

	aiRecognitionTemplate = response.Response.AIRecognitionTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsAiRecognitionTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteAIRecognitionTemplateRequest()
	request.Definition = helper.StrToInt64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteAIRecognitionTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
