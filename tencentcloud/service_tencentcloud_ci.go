package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type CiService struct {
	client *connectivity.TencentCloudClient
}

func (me *CiService) DescribeCiBucketById(ctx context.Context, bucket string) (serviceResult *cos.CIServiceResult, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIService", bucket, errRet.Error())
		}
	}()

	result, response, err := me.client.UsePicClient(bucket).CI.GetCIService(ctx)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s], http status [%s]\n", logId, "GetCIService", bucket, result, response.Status)

	serviceResult = result
	return
}

func (me *CiService) DeleteCiBucketById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "CloseCIService", bucket, errRet.Error())
		}
	}()

	response, err := me.client.UsePicClient(bucket).CI.CloseCIService(ctx)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "CloseCIService", bucket, response.Status)

	return
}

func (me *CiService) DescribeCiBucketPicStyleById(ctx context.Context, bucket, styleName string) (styleRule *cos.StyleRule, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetStyle", bucket+"#"+styleName, errRet.Error())
		}
	}()

	styleResult, response, err := me.client.UsePicClient(bucket).CI.GetStyle(ctx, &cos.GetStyleOptions{
		StyleName: styleName,
	})
	if err != nil {
		// if response.StatusCode == 400 {
		// 	log.Printf("[CRITAL]%s api[%s] success, request body [%s], response status [%v]\n", logId, "GetStyle", bucket+"#"+styleName, response.StatusCode)
		// 	return
		// }
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "GetStyle", bucket+"#"+styleName, response.Status)

	if len(styleResult.StyleRule) < 1 {
		return
	}

	styleRule = &styleResult.StyleRule[0]

	return
}

func (me *CiService) DeleteCiBucketPicStyleById(ctx context.Context, bucket, styleName string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "DeleteStyle", bucket+"#"+styleName, errRet.Error())
		}
	}()

	response, err := me.client.UsePicClient(bucket).CI.DeleteStyle(ctx, &cos.DeleteStyleOptions{
		StyleName: styleName,
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "DeleteStyle", bucket+"#"+styleName, response.Status)

	return
}

func (me *CiService) DescribeCiHotLinkById(ctx context.Context, bucket string) (hotLink *cos.HotLinkResult, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetHotLink", bucket, errRet.Error())
		}
	}()

	hotLinkResult, response, err := me.client.UsePicClient(bucket).CI.GetHotLink(ctx)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "GetHotLink", bucket, response.Status)

	if len(hotLinkResult.Url) < 1 {
		return
	}

	hotLink = hotLinkResult

	return
}

func (me *CiService) DescribeCiMediaTemplateById(ctx context.Context, bucket, templateId string) (mediaSnapshotTemplate *cos.Template, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "DescribeMediaTemplate", bucket, errRet.Error())
		}
	}()

	response, _, err := me.client.UseCiClient(bucket).CI.DescribeMediaTemplate(ctx, &cos.DescribeMediaTemplateOptions{
		Ids: templateId,
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%v, %v], response body [%+v]\n", logId, "DescribeMediaTemplate", bucket, templateId, response)

	if len(response.TemplateList) < 1 {
		return
	}

	mediaSnapshotTemplate = &response.TemplateList[0]
	return
}

func (me *CiService) DeleteCiMediaTemplateById(ctx context.Context, bucket, templateId string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "DeleteMediaTemplate", bucket, errRet.Error())
		}
	}()

	response, _, err := me.client.UseCiClient(bucket).CI.DeleteMediaTemplate(ctx, templateId)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%+v]\n", logId, "DeleteMediaTemplate", bucket, response)

	return
}

func (me *CiService) CloseCiOriginalImageProtectionById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "CloseCIOriginalImageProtection", bucket, errRet.Error())
		}
	}()

	_, err := RetryWithContext(ctx, writeRetryTimeout, func(ctx context.Context) (interface{}, error) {
		return me.client.UsePicClient(bucket).CI.CloseOriginProtect(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "CloseCIOriginalImageProtection", bucket)

	return
}

func (me *CiService) OpenCiOriginalImageProtectionById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "OpenCIOriginalImageProtection", bucket, errRet.Error())
		}
	}()

	_, err := RetryWithContext(ctx, writeRetryTimeout, func(ctx context.Context) (interface{}, error) {
		return me.client.UsePicClient(bucket).CI.OpenOriginProtect(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "OpenCIOriginalImageProtection", bucket)

	return
}

func (me *CiService) GetCiOriginalImageProtectionById(ctx context.Context, bucket string) (*cos.OriginProtectResult, error) {
	var errRet error
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIOriginalImageProtection", bucket, errRet.Error())
		}
	}()

	resRaw, err := RetryWithContext(ctx, readRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UsePicClient(bucket).CI.GetOriginProtect(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIOriginalImageProtection", bucket)

	return resRaw.(*cos.OriginProtectResult), nil
}

func (me *CiService) CloseCiGuetzliById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "CloseCIGuetzli", bucket, errRet.Error())
		}
	}()

	_, err := RetryWithContext(ctx, writeRetryTimeout, func(ctx context.Context) (interface{}, error) {
		return me.client.UsePicClient(bucket).CI.DeleteGuetzli(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "CloseCIGuetzli", bucket)

	return
}

func (me *CiService) OpenCiGuetzliById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "OpenCIGuetzli", bucket, errRet.Error())
		}
	}()

	_, err := RetryWithContext(ctx, writeRetryTimeout, func(ctx context.Context) (interface{}, error) {
		return me.client.UsePicClient(bucket).CI.PutGuetzli(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "OpenCIGuetzli", bucket)

	return
}

func (me *CiService) GetCiGuetzliById(ctx context.Context, bucket string) (*cos.GetGuetzliResult, error) {
	var errRet error
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIGuetzli", bucket, errRet.Error())
		}
	}()

	resRaw, err := RetryWithContext(ctx, readRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UsePicClient(bucket).CI.GetGuetzli(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIGuetzli", bucket)

	return resRaw.(*cos.GetGuetzliResult), nil
}

func (me *CiService) DescribeCiStyleById(ctx context.Context, styleName string) (style *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewGetStyleRequest()
	request.StyleName = &styleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().GetStyle(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	style = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiStyleById(ctx context.Context, styleName string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteStyleRequest()
	request.StyleName = &styleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteStyle(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaAnimationTemplateById(ctx context.Context, templateId string) (mediaAnimationTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaAnimationTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaAnimationTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaConcatTemplateById(ctx context.Context, templateId string) (mediaConcatTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaConcatTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaConcatTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaPicProcessTemplateById(ctx context.Context, templateId string) (mediaPicProcessTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaPicProcessTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaPicProcessTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaSmartCoverTemplateById(ctx context.Context, templateId string) (mediaSmartCoverTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaSmartCoverTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaSmartCoverTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaSnapshotTemplateById(ctx context.Context, templateId string) (mediaSnapshotTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaSnapshotTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaSnapshotTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaSpeechRecognitionTemplateById(ctx context.Context, templateId string) (mediaSpeechRecognitionTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaSpeechRecognitionTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaSpeechRecognitionTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaSuperResolutionTemplateById(ctx context.Context, templateId string) (mediaSuperResolutionTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaSuperResolutionTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaSuperResolutionTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaTranscodeTemplateById(ctx context.Context, templateId string) (mediaTranscodeTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaTranscodeTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaTranscodeTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaTtsTemplateById(ctx context.Context, templateId string) (mediaTtsTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaTtsTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaTtsTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaVideoProcessTemplateById(ctx context.Context, templateId string) (mediaVideoProcessTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaVideoProcessTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaVideoProcessTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaVoiceSeparateTemplateById(ctx context.Context, templateId string) (mediaVoiceSeparateTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	mediaVoiceSeparateTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiMediaVoiceSeparateTemplateById(ctx context.Context, templateId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaTemplateRequest()
	request.TemplateId = &templateId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiMediaWorkflowById(ctx context.Context, workflowId string) (mediaWorkflow *ci.MediaWorkflow, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDescribeMediaWorkflowRequest()
	request.WorkflowId = &workflowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DescribeMediaWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.MediaWorkflow) < 1 {
		return
	}

	mediaWorkflow = response.Response.MediaWorkflow[0]
	return
}

func (me *CiService) DeleteCiMediaWorkflowById(ctx context.Context, workflowId string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteMediaWorkflowRequest()
	request.WorkflowId = &workflowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteMediaWorkflow(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CiService) DescribeCiPosterTemplateById(ctx context.Context, styleName string) (posterTemplate *ci.GetStyleResult, errRet error) {
	logId := getLogId(ctx)

	request := ci.NewGetStyleRequest()
	request.StyleName = &styleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().GetStyle(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GetStyleResult) < 1 {
		return
	}

	posterTemplate = response.Response.GetStyleResult[0]
	return
}

func (me *CiService) DeleteCiPosterTemplateById(ctx context.Context, styleName string) (errRet error) {
	logId := getLogId(ctx)

	request := ci.NewDeleteStyleRequest()
	request.StyleName = &styleName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCiClient().DeleteStyle(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
