package cwp

import (
	"context"
	"log"
	"strconv"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	cwp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cwp/v20180228"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type CwpService struct {
	client *connectivity.TencentCloudClient
}

func (me *CwpService) DescribeCwpLicenseOrderById(ctx context.Context, resourceId string) (licenseOrder *cwp.LicenseDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cwp.NewDescribeLicenseListRequest()
	request.Filters = []*cwp.Filters{
		{
			Name:       common.StringPtr("ResourceId"),
			Values:     common.StringPtrs([]string{resourceId}),
			ExactMatch: common.BoolPtr(true),
		},
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCwpClient().DescribeLicenseList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.List) < 1 {
		return
	}

	licenseOrder = response.Response.List[0]
	return
}

func (me *CwpService) DeleteCwpLicenseOrderById(ctx context.Context, resourceId string, licenseType *uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cwp.NewDestroyOrderRequest()
	request.ResourceId = &resourceId
	request.LicenseType = licenseType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCwpClient().DestroyOrder(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CwpService) DescribeCwpLicenseBindAttachmentById(ctx context.Context, resourceId, quuid string, licenseId, licenseType uint64) (licenseBindAttachment *cwp.LicenseBindDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cwp.NewDescribeLicenseBindListRequest()
	request.ResourceId = &resourceId
	request.LicenseId = &licenseId
	request.LicenseType = &licenseType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCwpClient().DescribeLicenseBindList(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.List) < 1 {
		return
	}

	for _, item := range response.Response.List {
		if *item.Quuid == quuid {
			licenseBindAttachment = item
			return
		}
	}

	return
}

func (me *CwpService) DeleteCwpLicenseBindAttachmentById(ctx context.Context, resourceId, quuid, licenseType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	licenseTypeInt, _ := strconv.ParseUint(licenseType, 10, 64)

	request := cwp.NewModifyLicenseUnBindsRequest()
	request.ResourceId = common.StringPtr(resourceId)
	request.QuuidList = common.StringPtrs([]string{quuid})
	request.LicenseType = common.Uint64Ptr(licenseTypeInt)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCwpClient().ModifyLicenseUnBinds(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CwpService) DescribeCwpMachinesSimpleByFilter(ctx context.Context, param map[string]interface{}) (machinesSimple []*cwp.MachineSimple, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cwp.NewDescribeMachinesSimpleRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "MachineType" {
			request.MachineType = v.(*string)
		}

		if k == "MachineRegion" {
			request.MachineRegion = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*cwp.Filter)
		}

		if k == "ProjectIds" {
			request.ProjectIds = v.([]*uint64)
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
		response, err := me.client.UseCwpClient().DescribeMachinesSimple(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Machines) < 1 {
			break
		}

		machinesSimple = append(machinesSimple, response.Response.Machines...)
		if len(response.Response.Machines) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
