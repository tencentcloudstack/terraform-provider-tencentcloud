package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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

	result, response, err := me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.GetCIService(ctx)
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

	response, err := me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.CloseCIService(ctx)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "CloseCIService", bucket, response.Status)

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
		return me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.CloseOriginProtect(ctx)
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
		return me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.OpenOriginProtect(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "OpenCIOriginalImageProtection", bucket)

	return
}

func (me *CiService) GetCiOriginalImageProtectionById(ctx context.Context, bucket string) (*ci.OriginProtectResult, error) {
	var errRet error
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIOriginalImageProtection", bucket, errRet.Error())
		}
	}()

	resRaw, err := RetryWithContext(ctx, readRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.GetOriginProtect(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIOriginalImageProtection", bucket)

	return resRaw.(*ci.OriginProtectResult), nil
}

func (me *CiService) CloseCiGuetzliById(ctx context.Context, bucket string) (errRet error) {
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "CloseCIGuetzli", bucket, errRet.Error())
		}
	}()

	_, err := RetryWithContext(ctx, writeRetryTimeout, func(ctx context.Context) (interface{}, error) {
		return me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.DeleteGuetzli(ctx)
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
		return me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.PutGuetzli(ctx)
	})

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "OpenCIGuetzli", bucket)

	return
}

func (me *CiService) GetCiGuetzliById(ctx context.Context, bucket string) (*ci.GetGuetzliResult, error) {
	var errRet error
	logId := getLogId(ctx)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "GetCIGuetzli", bucket, errRet.Error())
		}
	}()

	resRaw, err := RetryWithContext(ctx, readRetryTimeout, func(ctx context.Context) (interface{}, error) {
		res, _, err := me.client.UseCiClient(bucket, connectivity.CI_HOST_PIC).CI.GetGuetzli(ctx)
		return res, err
	})

	if err != nil {
		errRet = err
		return nil, errRet
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s]\n", logId, "GetCIGuetzli", bucket)

	return resRaw.(*ci.GetGuetzliResult), nil
}
