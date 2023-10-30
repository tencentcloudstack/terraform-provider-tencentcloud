package tencentcloud

import (
	"context"
	"log"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	wedata "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/wedata/v20210820"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type WedataService struct {
	client *connectivity.TencentCloudClient
}

func (me *WedataService) DescribeWedataRuleTemplateById(ctx context.Context, projectId string, ruleTemplateId string) (ruleTemplate *wedata.RuleTemplate, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
	request.TemplateId = helper.StrToUint64Point(ruleTemplateId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response.Data != nil {

		ruleTemplate = response.Response.Data
	}

	return
}

func (me *WedataService) DeleteWedataRuleTemplateById(ctx context.Context, projectId, ruleTemplateId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteRuleTemplateRequest()
	request.ProjectId = helper.String(projectId)
	request.Ids = []*uint64{helper.StrToUint64Point(ruleTemplateId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteRuleTemplate(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *WedataService) DescribeWedataRuleTemplatesByFilter(ctx context.Context, param map[string]interface{}) (ruleTemplates []*wedata.RuleTemplate, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeRuleTemplatesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*uint64)
		}
		if k == "SourceObjectType" {
			request.SourceObjectType = v.(*uint64)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}
		if k == "SourceEngineTypes" {
			request.SourceEngineTypes = v.([]*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeRuleTemplates(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ruleTemplates = response.Response.Data

	return
}

func (me *WedataService) DescribeWedataDataSourceListByFilter(ctx context.Context, param map[string]interface{}) (dataSourceList []*wedata.DataSourceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedata.OrderField)
		}

		if k == "Filters" {
			request.Filters = v.([]*wedata.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNum  uint64 = 0
		pageSize uint64 = 20
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataClient().DescribeDataSourceList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.Rows) < 1 {
			break
		}

		dataSourceList = append(dataSourceList, response.Response.Data.Rows...)
		if len(response.Response.Data.Rows) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataDataSourceInfoListByFilter(ctx context.Context, param map[string]interface{}) (dataSourceInfoList []*wedata.DatasourceBaseInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceInfoListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectId" {
			request.ProjectId = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.(*wedata.Filter)
		}

		if k == "OrderFields" {
			request.OrderFields = v.(*wedata.OrderField)
		}

		if k == "Type" {
			request.Type = v.(*string)
		}

		if k == "DatasourceName" {
			request.DatasourceName = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		pageNum  uint64 = 0
		pageSize uint64 = 20
	)
	for {
		request.PageNumber = &pageNum
		request.PageSize = &pageSize
		response, err := me.client.UseWedataClient().DescribeDataSourceInfoList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DatasourceSet) < 1 {
			break
		}

		dataSourceInfoList = append(dataSourceInfoList, response.Response.DatasourceSet...)
		if len(response.Response.DatasourceSet) < int(pageSize) {
			break
		}

		pageNum += pageSize
	}

	return
}

func (me *WedataService) DescribeWedataDataSourceWithoutInfoByFilter(ctx context.Context, param map[string]interface{}) (dataSourceWithoutInfo []*wedata.DataSourceInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = wedata.NewDescribeDataSourceWithoutInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "OrderFields" {
			request.OrderFields = v.([]*wedata.OrderField)
		}

		if k == "Filters" {
			request.Filters = v.([]*wedata.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeDataSourceWithoutInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Data) < 1 {
		return
	}

	dataSourceWithoutInfo = response.Response.Data
	return
}

func (me *WedataService) DescribeWedataDatasourceById(ctx context.Context, datasourceId string) (datasource *wedata.DataSourceInfo, errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDescribeDatasourceRequest()
	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.Id = &Id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DescribeDatasource(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	datasource = response.Response.Data
	return
}

func (me *WedataService) DeleteWedataDatasourceById(ctx context.Context, datasourceId string) (errRet error) {
	logId := getLogId(ctx)

	request := wedata.NewDeleteDataSourcesRequest()
	Id, _ := strconv.ParseUint(datasourceId, 10, 64)
	request.Ids = common.Uint64Ptrs([]uint64{Id})

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseWedataClient().DeleteDataSources(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
