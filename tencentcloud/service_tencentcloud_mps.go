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

func (me *MpsService) DescribeMpsPersonSampleById(ctx context.Context, personId string) (personSample *mps.AiSamplePerson, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribePersonSamplesRequest()
	request.PersonIds = []*string{helper.String(personId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribePersonSamples(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.PersonSet) < 1 {
		return
	}

	personSample = response.Response.PersonSet[0]
	return
}

func (me *MpsService) DeleteMpsPersonSampleById(ctx context.Context, personId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeletePersonSampleRequest()
	request.PersonId = &personId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeletePersonSample(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsWordSamplesById(ctx context.Context, keywords []string) (wordSamples []*mps.AiSampleWord, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeWordSamplesRequest()
	request.Keywords = helper.Strings(keywords)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeWordSamples(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WordSet) < 1 {
		return
	}

	wordSamples = response.Response.WordSet
	return
}

func (me *MpsService) DescribeMpsWordSampleById(ctx context.Context, keyword string) (wordSample *mps.AiSampleWord, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeWordSamplesRequest()
	request.Keywords = []*string{
		helper.String(keyword),
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeWordSamples(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.WordSet) < 1 {
		return
	}

	wordSample = response.Response.WordSet[0]
	return
}

func (me *MpsService) DeleteMpsWordSamplesById(ctx context.Context, keywords []string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteWordSamplesRequest()
	request.Keywords = helper.Strings(keywords)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteWordSamples(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsScheduleById(ctx context.Context, scheduleId *string) (schedules []*mps.SchedulesInfo, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeSchedulesRequest()
	if scheduleId != nil {
		request.ScheduleIds = []*int64{helper.StrToInt64Point(*scheduleId)}
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeSchedules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	schedules = response.Response.ScheduleInfoSet
	return
}

func (me *MpsService) DeleteMpsScheduleById(ctx context.Context, scheduleId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteScheduleRequest()
	request.ScheduleId = helper.StrToInt64Point(scheduleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteSchedule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsSchedulesByFilter(ctx context.Context, param map[string]interface{}) (schedules []*mps.SchedulesInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeSchedulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ScheduleIds" {
			request.ScheduleIds = helper.InterfacesIntInt64Point(v.([]interface{}))
		}
		if k == "TriggerType" {
			request.TriggerType = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeSchedules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ScheduleInfoSet) < 1 {
			break
		}
		schedules = append(schedules, response.Response.ScheduleInfoSet...)
		if len(response.Response.ScheduleInfoSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsFlowById(ctx context.Context, flowId string) (flow *mps.DescribeFlow, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeStreamLinkFlowRequest()
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeStreamLinkFlow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flow = response.Response.Info
	return
}

func (me *MpsService) DeleteMpsFlowById(ctx context.Context, flowId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteStreamLinkFlowRequest()
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteStreamLinkFlow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsEventById(ctx context.Context, eventId string) (event *mps.DescribeEvent, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeStreamLinkEventRequest()
	request.EventId = &eventId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeStreamLinkEvent(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	event = response.Response.Info
	return
}

func (me *MpsService) DeleteMpsEventById(ctx context.Context, eventId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteStreamLinkEventRequest()
	request.EventId = &eventId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteStreamLinkEvent(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsTaskDetailById(ctx context.Context, taskId string) (manageTaskOperation *mps.DescribeTaskDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeTaskDetailRequest()
	request.TaskId = &taskId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeTaskDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	manageTaskOperation = response.Response
	return
}

func (me *MpsService) DeleteMpsOutputById(ctx context.Context, flowId, outputId string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteStreamLinkOutputRequest()
	request.FlowId = &flowId
	request.OutputId = &outputId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteStreamLinkOutput(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsInputById(ctx context.Context, flowId, inputId string) (input *mps.DescribeInput, errRet error) {
	logId := getLogId(ctx)

	flow, err := me.DescribeMpsFlowById(ctx, flowId)
	if err != nil {
		return nil, err
	}

	for _, iter := range flow.InputGroup {
		if *iter.InputId == inputId {
			input = iter
			break
		}
	}

	log.Printf("[DEBUG]%s `DescribeMpsInputById` success, inputId: %s, flowId: %s \n", logId, *input.InputId, flowId)
	return
}

func (me *MpsService) DescribeMpsOutputById(ctx context.Context, flowId, outputId string) (output *mps.DescribeOutput, errRet error) {
	logId := getLogId(ctx)

	flow, err := me.DescribeMpsFlowById(ctx, flowId)
	if err != nil {
		return nil, err
	}

	for _, iter := range flow.OutputGroup {
		if *iter.OutputId == outputId {
			output = iter
			break
		}
	}

	log.Printf("[DEBUG]%s `DescribeMpsOutputById` success, outputId: %s, flowId: %s \n", logId, *output.OutputId, flowId)
	return
}

func (me *MpsService) DescribeMpsTasksByFilter(ctx context.Context, param map[string]interface{}) (tasks []*mps.TaskSimpleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "ScrollToken" {
			request.ScrollToken = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeTasks(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}
	tasks = response.Response.TaskSet

	return
}

func (me *MpsService) DescribeMpsContentReviewTemplateById(ctx context.Context, definition string) (contentReviewTemplate *mps.ContentReviewTemplateItem, errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDescribeContentReviewTemplatesRequest()
	request.Definitions = []*int64{helper.StrToInt64Point(definition)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DescribeContentReviewTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.ContentReviewTemplateSet) < 1 {
		return
	}

	contentReviewTemplate = response.Response.ContentReviewTemplateSet[0]
	return
}

func (me *MpsService) DeleteMpsContentReviewTemplateById(ctx context.Context, definition string) (errRet error) {
	logId := getLogId(ctx)

	request := mps.NewDeleteContentReviewTemplateRequest()
	request.Definition = helper.StrToInt64Point(definition)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().DeleteContentReviewTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MpsService) DescribeMpsParseLiveStreamProcessNotificationByFilter(ctx context.Context, param map[string]interface{}) (ret *mps.ParseLiveStreamProcessNotificationResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewParseLiveStreamProcessNotificationRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Content" {
			request.Content = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().ParseLiveStreamProcessNotification(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response

	return
}

func (me *MpsService) DescribeMpsParseNotificationByFilter(ctx context.Context, param map[string]interface{}) (ret *mps.ParseNotificationResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewParseNotificationRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Content" {
			request.Content = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMpsClient().ParseNotification(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response

	return
}

func (me *MpsService) DescribeMpsMediaMetaDataByFilter(ctx context.Context, param map[string]interface{}) (ret *mps.MediaMetaData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeMediaMetaDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InputInfo" {
			request.InputInfo = v.(*mps.MediaInputInfo)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMpsClient().DescribeMediaMetaData(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response.MetaData

	return
}

func (me *MpsService) DescribeMpsAdaptiveDynamicStreamingTemplatesByFilter(ctx context.Context, param map[string]interface{}) (adaptiveDynamicStreamingTemplates []*mps.AdaptiveDynamicStreamingTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeAdaptiveDynamicStreamingTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*uint64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "AdaptiveDynamicStreamingTemplateSet" {
			request.AdaptiveDynamicStreamingTemplateSet = v.([]*mps.AdaptiveDynamicStreamingTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeAdaptiveDynamicStreamingTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AdaptiveDynamicStreamingTemplateSet) < 1 {
			break
		}
		adaptiveDynamicStreamingTemplates = append(adaptiveDynamicStreamingTemplates, response.Response.AdaptiveDynamicStreamingTemplateSet...)
		if len(response.Response.AdaptiveDynamicStreamingTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsAiAnalysisTemplatesByFilter(ctx context.Context, param map[string]interface{}) (aiAnalysisTemplates []*mps.AIAnalysisTemplateItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeAIAnalysisTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*int64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "AIAnalysisTemplateSet" {
			request.AIAnalysisTemplateSet = v.([]*mps.AIAnalysisTemplateItem)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeAIAnalysisTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AIAnalysisTemplateSet) < 1 {
			break
		}
		aiAnalysisTemplates = append(aiAnalysisTemplates, response.Response.AIAnalysisTemplateSet...)
		if len(response.Response.AIAnalysisTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsAiRecognitionTemplatesByFilter(ctx context.Context, param map[string]interface{}) (aiRecognitionTemplates []*mps.AIRecognitionTemplateItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeAIRecognitionTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*int64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "AIRecognitionTemplateSet" {
			request.AIRecognitionTemplateSet = v.([]*mps.AIRecognitionTemplateItem)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeAIRecognitionTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AIRecognitionTemplateSet) < 1 {
			break
		}
		aiRecognitionTemplates = append(aiRecognitionTemplates, response.Response.AIRecognitionTemplateSet...)
		if len(response.Response.AIRecognitionTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsAnimatedGraphicsTemplatesByFilter(ctx context.Context, param map[string]interface{}) (animatedGraphicsTemplates []*mps.AnimatedGraphicsTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeAnimatedGraphicsTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*uint64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "AnimatedGraphicsTemplateSet" {
			request.AnimatedGraphicsTemplateSet = v.([]*mps.AnimatedGraphicsTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeAnimatedGraphicsTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.AnimatedGraphicsTemplateSet) < 1 {
			break
		}
		animatedGraphicsTemplates = append(animatedGraphicsTemplates, response.Response.AnimatedGraphicsTemplateSet...)
		if len(response.Response.AnimatedGraphicsTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsContentReviewTemplatesByFilter(ctx context.Context, param map[string]interface{}) (contentReviewTemplates []*mps.ContentReviewTemplateItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeContentReviewTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeContentReviewTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ContentReviewTemplateSet) < 1 {
			break
		}
		contentReviewTemplates = append(contentReviewTemplates, response.Response.ContentReviewTemplateSet...)
		if len(response.Response.ContentReviewTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsDescribeMediaMetaDataByFilter(ctx context.Context, param map[string]interface{}) (describeMediaMetaData []*mps.MediaMetaData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeMediaMetaDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InputInfo" {
			request.InputInfo = v.(map[string]interface{})
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeMediaMetaData(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MetaData) < 1 {
			break
		}
		describeMediaMetaData = append(describeMediaMetaData, response.Response.MetaData...)
		if len(response.Response.MetaData) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsImageSpriteTemplatesByFilter(ctx context.Context, param map[string]interface{}) (imageSpriteTemplates []*mps.ImageSpriteTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeImageSpriteTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*uint64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "ImageSpriteTemplateSet" {
			request.ImageSpriteTemplateSet = v.([]*mps.ImageSpriteTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeImageSpriteTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ImageSpriteTemplateSet) < 1 {
			break
		}
		imageSpriteTemplates = append(imageSpriteTemplates, response.Response.ImageSpriteTemplateSet...)
		if len(response.Response.ImageSpriteTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsMediaMetaDataByFilter(ctx context.Context, param map[string]interface{}) (mediaMetaData []*mps.MediaMetaData, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeMediaMetaDataRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InputInfo" {
			request.InputInfo = v.(map[string]interface{})
		}
		if k == "MetaData" {
			request.MetaData = v.(map[string]interface{})
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeMediaMetaData(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MetaData) < 1 {
			break
		}
		mediaMetaData = append(mediaMetaData, response.Response.MetaData...)
		if len(response.Response.MetaData) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsParseLiveStreamProcessNotificationByFilter(ctx context.Context, param map[string]interface{}) (parseLiveStreamProcessNotification []*mps.ParseLiveStreamProcessNotificationResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewParseLiveStreamProcessNotificationRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Content" {
			request.Content = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().ParseLiveStreamProcessNotification(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.NotificationType) < 1 {
			break
		}
		parseLiveStreamProcessNotification = append(parseLiveStreamProcessNotification, response.Response.NotificationType...)
		if len(response.Response.NotificationType) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsParseNotificationByFilter(ctx context.Context, param map[string]interface{}) (parseNotification []*mps.ParseNotificationResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewParseNotificationRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Content" {
			request.Content = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().ParseNotification(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EventType) < 1 {
			break
		}
		parseNotification = append(parseNotification, response.Response.EventType...)
		if len(response.Response.EventType) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsPersonSamplesByFilter(ctx context.Context, param map[string]interface{}) (personSamples []*mps.AiSamplePerson, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribePersonSamplesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "PersonIds" {
			request.PersonIds = v.([]*string)
		}
		if k == "Names" {
			request.Names = v.([]*string)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "PersonSet" {
			request.PersonSet = v.([]*mps.AiSamplePerson)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribePersonSamples(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.PersonSet) < 1 {
			break
		}
		personSamples = append(personSamples, response.Response.PersonSet...)
		if len(response.Response.PersonSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsSampleSnapshotTemplatesByFilter(ctx context.Context, param map[string]interface{}) (sampleSnapshotTemplates []*mps.SampleSnapshotTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeSampleSnapshotTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*uint64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "SampleSnapshotTemplateSet" {
			request.SampleSnapshotTemplateSet = v.([]*mps.SampleSnapshotTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeSampleSnapshotTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SampleSnapshotTemplateSet) < 1 {
			break
		}
		sampleSnapshotTemplates = append(sampleSnapshotTemplates, response.Response.SampleSnapshotTemplateSet...)
		if len(response.Response.SampleSnapshotTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsSchedulesByFilter(ctx context.Context, param map[string]interface{}) (schedules []*mps.SchedulesInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeSchedulesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ScheduleIds" {
			request.ScheduleIds = v.([]*int64)
		}
		if k == "TriggerType" {
			request.TriggerType = v.(*string)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeSchedules(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ScheduleInfoSet) < 1 {
			break
		}
		schedules = append(schedules, response.Response.ScheduleInfoSet...)
		if len(response.Response.ScheduleInfoSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsSnapshotByTimeoffsetTemplatesByFilter(ctx context.Context, param map[string]interface{}) (snapshotByTimeoffsetTemplates []*mps.SnapshotByTimeOffsetTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeSnapshotByTimeOffsetTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*uint64)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "SnapshotByTimeOffsetTemplateSet" {
			request.SnapshotByTimeOffsetTemplateSet = v.([]*mps.SnapshotByTimeOffsetTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeSnapshotByTimeOffsetTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.SnapshotByTimeOffsetTemplateSet) < 1 {
			break
		}
		snapshotByTimeoffsetTemplates = append(snapshotByTimeoffsetTemplates, response.Response.SnapshotByTimeOffsetTemplateSet...)
		if len(response.Response.SnapshotByTimeOffsetTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsTasksByFilter(ctx context.Context, param map[string]interface{}) (tasks []*mps.TaskSimpleInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeTasksRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Status" {
			request.Status = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeTasks(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TaskSet) < 1 {
			break
		}
		tasks = append(tasks, response.Response.TaskSet...)
		if len(response.Response.TaskSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsTranscodeTemplatesByFilter(ctx context.Context, param map[string]interface{}) (transcodeTemplates []*mps.TranscodeTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeTranscodeTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "ContainerType" {
			request.ContainerType = v.(*string)
		}
		if k == "TEHDType" {
			request.TEHDType = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "TranscodeType" {
			request.TranscodeType = v.(*string)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "TranscodeTemplateSet" {
			request.TranscodeTemplateSet = v.([]*mps.TranscodeTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeTranscodeTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.TranscodeTemplateSet) < 1 {
			break
		}
		transcodeTemplates = append(transcodeTemplates, response.Response.TranscodeTemplateSet...)
		if len(response.Response.TranscodeTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsWatermarkTemplatesByFilter(ctx context.Context, param map[string]interface{}) (watermarkTemplates []*mps.WatermarkTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeWatermarkTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Definitions" {
			request.Definitions = v.([]*int64)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*uint64)
		}
		if k == "Limit" {
			request.Limit = v.(*uint64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "WatermarkTemplateSet" {
			request.WatermarkTemplateSet = v.([]*mps.WatermarkTemplate)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeWatermarkTemplates(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.WatermarkTemplateSet) < 1 {
			break
		}
		watermarkTemplates = append(watermarkTemplates, response.Response.WatermarkTemplateSet...)
		if len(response.Response.WatermarkTemplateSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MpsService) DescribeMpsWorkflowsByFilter(ctx context.Context, param map[string]interface{}) (workflows []*mps.WorkflowInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mps.NewDescribeWorkflowsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "WorkflowIds" {
			request.WorkflowIds = v.([]*int64)
		}
		if k == "Status" {
			request.Status = v.(*string)
		}
		if k == "Offset" {
			request.Offset = v.(*int64)
		}
		if k == "Limit" {
			request.Limit = v.(*int64)
		}
		if k == "TotalCount" {
			request.TotalCount = v.(*uint64)
		}
		if k == "WorkflowInfoSet" {
			request.WorkflowInfoSet = v.([]*mps.WorkflowInfo)
		}
		if k == "RequestId" {
			request.RequestId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMpsClient().DescribeWorkflows(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.WorkflowInfoSet) < 1 {
			break
		}
		workflows = append(workflows, response.Response.WorkflowInfoSet...)
		if len(response.Response.WorkflowInfoSet) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
