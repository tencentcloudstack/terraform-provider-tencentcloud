package bi

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	bi "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bi/v20220105"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type BiService struct {
	client *connectivity.TencentCloudClient
}

func (me *BiService) DescribeBiDatasourceCloudById(ctx context.Context, projectId, id uint64) (datasourceCloud *bi.DatasourceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDescribeDatasourceListRequest()
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeDatasourceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		for _, v := range response.Response.Data.List {
			if *v.Id == id {
				datasourceCloud = v
				return
			}
		}
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BiService) DeleteBiDatasourceCloudById(ctx context.Context, projectId, id uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDeleteDatasourceRequest()
	request.ProjectId = &projectId
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DeleteDatasource(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *BiService) DescribeBiProjectById(ctx context.Context, projectId uint64) (project *bi.Project, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDescribeProjectInfoRequest()
	request.Id = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DescribeProjectInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	project = response.Response.Data
	return
}

func (me *BiService) DeleteBiProjectById(ctx context.Context, projectId uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDeleteProjectRequest()
	request.Id = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DeleteProject(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *BiService) DescribeBiUserRoleById(ctx context.Context, userId string) (userRole *bi.UserRoleListDataUserRoleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDescribeUserRoleListRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeUserRoleList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		for _, v := range response.Response.Data.List {
			if *v.UserId == userId {
				userRole = v
				return
			}
		}
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BiService) DeleteBiUserRoleById(ctx context.Context, userId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDeleteUserRoleRequest()
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DeleteUserRole(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *BiService) DescribeBiProjectUserRoleById(ctx context.Context, projectId int64, userId string) (projectUserRole *bi.UserRoleListDataUserRoleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDescribeUserRoleProjectListRequest()
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeUserRoleProjectList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		for _, v := range response.Response.Data.List {
			if *v.UserId == userId {
				projectUserRole = v
				return
			}
		}
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BiService) DeleteBiProjectUserRoleById(ctx context.Context, projectId int64, userId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDeleteUserRoleProjectRequest()
	request.ProjectId = &projectId
	request.UserId = &userId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DeleteUserRoleProject(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *BiService) DescribeBiProjectByFilter(ctx context.Context, param map[string]interface{}) (project []*bi.Project, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = bi.NewDescribeProjectListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "PageNo" {
			request.PageNo = v.(*uint64)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "AllPage" {
			request.AllPage = v.(*bool)
		}
		if k == "ModuleCollection" {
			request.ModuleCollection = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeProjectList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		project = append(project, response.Response.Data.List...)
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BiService) DescribeBiDatasourceById(ctx context.Context, projectId uint64, id uint64) (datasource *bi.DatasourceInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDescribeDatasourceListRequest()
	request.ProjectId = &projectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeDatasourceList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		for _, v := range response.Response.Data.List {
			if *v.Id == id {
				datasource = v
				return
			}
		}
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *BiService) DeleteBiDatasourceById(ctx context.Context, projectId uint64, id uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := bi.NewDeleteDatasourceRequest()
	request.ProjectId = &projectId
	request.Id = &id

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseBiClient().DeleteDatasource(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *BiService) DescribeBiUserProjectByFilter(ctx context.Context, param map[string]interface{}) (userProject []*bi.UserIdAndUserName, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = bi.NewDescribeUserProjectListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
		}
		if k == "AllPage" {
			request.AllPage = v.(*bool)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.PageNo = &offset
		request.PageSize = &limit
		response, err := me.client.UseBiClient().DescribeUserProjectList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data.List) < 1 {
			break
		}
		userProject = append(userProject, response.Response.Data.List...)
		if len(response.Response.Data.List) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
