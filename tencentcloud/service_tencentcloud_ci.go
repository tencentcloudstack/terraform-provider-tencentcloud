package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentyun/cos-go-sdk-v5"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

type CiService struct {
	client *connectivity.TencentCloudClient
}

func (me *CiService) DescribeCiBucketById(ctx context.Context, bucket string) (serviceResult *ci.CIServiceResult, errRet error) {
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

func (me *CiService) DescribeCiBucketPicStyleById(ctx context.Context, bucket, styleName string) (styleRule *ci.StyleRule, errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetStyle", bucket+"#"+styleName, errRet.Error())
		}
	}()

	styleResult, response, err := me.client.UsePicClient(bucket).CI.GetStyle(ctx, &ci.GetStyleOptions{
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

	response, err := me.client.UsePicClient(bucket).CI.DeleteStyle(ctx, &ci.DeleteStyleOptions{
		StyleName: styleName,
	})
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "DeleteStyle", bucket+"#"+styleName, response.Status)

	return
}

func (me *CiService) DescribeCiHotLinkById(ctx context.Context, bucket string) (hotLink *ci.HotLinkResult, errRet error) {
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

func (me *CiService) DescribeCiMediaTemplateById(ctx context.Context, bucket, templateId string) (mediaSnapshotTemplate *ci.Template, errRet error) {
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
	log.Printf("[DEBUG]%s api[%s] success, request body [%v], response body [%v]\n", logId, "DescribeMediaTemplate", bucket, response)

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
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, "DeleteMediaTemplate", bucket, response)

	return
}
