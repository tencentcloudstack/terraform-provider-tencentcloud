package cwp

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	response := cwp.NewDescribeLicenseListResponse()
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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCwpClient().DescribeLicenseList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.List == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe license list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

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

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCwpClient().DestroyOrder(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *CwpService) DescribeCwpLicenseBindAttachmentById(ctx context.Context, resourceId, quuid string, licenseId, licenseType uint64) (licenseBindAttachment *cwp.LicenseBindDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cwp.NewDescribeLicenseBindListRequest()
	response := cwp.NewDescribeLicenseBindListResponse()
	request.ResourceId = &resourceId
	request.LicenseId = &licenseId
	request.LicenseType = &licenseType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCwpClient().DescribeLicenseBindList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.List == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe license bind list failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	if len(response.Response.List) < 1 {
		return
	}

	for _, item := range response.Response.List {
		if item != nil && item.Quuid != nil && *item.Quuid == quuid {
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

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCwpClient().ModifyLicenseUnBinds(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	return
}

func (me *CwpService) DescribeCwpMachinesSimpleByFilter(ctx context.Context, param map[string]interface{}) (machinesSimple []*cwp.MachineSimple, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cwp.NewDescribeMachinesSimpleRequest()
		response = cwp.NewDescribeMachinesSimpleResponse()
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

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCwpClient().DescribeMachinesSimple(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Machines == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe machines simple failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Machines) < 1 {
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

func (me *CwpService) DescribeCwpMachinesByFilter(ctx context.Context, param map[string]interface{}) (ret []*cwp.Machine, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = cwp.NewDescribeMachinesRequest()
		response = cwp.NewDescribeMachinesResponse()
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
		if k == "Filters" {
			request.Filters = v.([]*cwp.Filter)
		}
		if k == "MachineRegion" {
			request.MachineRegion = v.(*string)
		}
		if k == "ProjectIds" {
			request.ProjectIds = v.([]*uint64)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseCwpClient().DescribeMachines(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Machines == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe machines failed, Response is nil."))
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if len(response.Response.Machines) < 1 {
			break
		}

		ret = append(ret, response.Response.Machines...)
		if len(response.Response.Machines) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *CwpService) DescribeCwpAutoOpenProversionConfigById(ctx context.Context) (ret *cwp.DescribeLicenseGeneralResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cwp.NewDescribeLicenseGeneralRequest()
	response := cwp.NewDescribeLicenseGeneralResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseCwpClient().DescribeLicenseGeneral(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe license general failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}
