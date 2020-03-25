package tencentcloud

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcaplusdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcaplusdb/v20190823"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TcaplusService struct {
	client *connectivity.TencentCloudClient
}

func (me *TcaplusService) CreateApp(ctx context.Context, idlType, appName, vpcId, subnetId, password string) (applicationId string, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewCreateClusterRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Password = &password
	request.VpcId = &vpcId
	request.SubnetId = &subnetId
	request.IdlType = &idlType
	request.ClusterName = &appName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().CreateCluster(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if response.Response.ClusterId == nil || *response.Response.ClusterId == "" {
		errRet = errors.New("TencentCloud SDK  return empty applicationId")
		return
	}
	applicationId = *response.Response.ClusterId
	return
}

func (me *TcaplusService) DescribeApps(ctx context.Context, applicationId string, applicationName string) (appInfos []*tcaplusdb.ClusterInfo, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeClustersRequest()

	appInfos = make([]*tcaplusdb.ClusterInfo, 0, 100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if applicationId != "" {
		request.ClusterIds = []*string{&applicationId}
	}

	if applicationName != "" {
		request.Filters = []*tcaplusdb.Filter{
			{
				Name:  helper.String("AppName"),
				Value: &applicationName,
			},
		}
	}
	var limit int64 = 20
	var offset int64 = 0
	request.Limit = &limit
	request.Offset = &offset
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTcaplusClient().DescribeClusters(request)
		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		appInfos = append(appInfos, response.Response.Clusters...)
		if len(response.Response.Clusters) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TcaplusService) DescribeApp(ctx context.Context, applicationId string) (appInfo tcaplusdb.ClusterInfo, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeClustersRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterIds = []*string{&applicationId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DescribeClusters(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceNotFound" {
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	has = true
	if len(response.Response.Clusters) == 0 {
		has = false
		return
	}
	if len(response.Response.Clusters) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d appInfo with one applicationId %s",
			len(response.Response.Clusters), applicationId)
		return
	}
	appInfo = *response.Response.Clusters[0]
	return
}

func (me *TcaplusService) DeleteApp(ctx context.Context, applicationId string) (taskId string, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDeleteClusterRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &applicationId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DeleteCluster(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if response.Response.TaskId == nil || *response.Response.TaskId == "" {
		errRet = errors.New("TencentCloud SDK  return empty delete taskId")
		return
	}
	taskId = *response.Response.TaskId
	return
}

func (me *TcaplusService) ModifyAppName(ctx context.Context, applicationId string, applicationName string) (errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewModifyClusterNameRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &applicationId
	request.ClusterName = helper.String(url.QueryEscape(applicationName))

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().ModifyClusterName(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
	}
	return
}

func (me *TcaplusService) ModifyAppPassword(ctx context.Context, applicationId string, oldPassword, newPassword string, oldPasswordExpireLast int64) (errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewModifyClusterPasswordRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &applicationId
	request.OldPassword = &oldPassword
	request.NewPassword = &newPassword
	request.Mode = helper.String("1")

	if oldPasswordExpireLast > 0 {
		expireTime := time.Now().Add(time.Second * time.Duration(oldPasswordExpireLast))
		loc, err := time.LoadLocation("Asia/Shanghai")
		if err != nil {
			errRet = fmt.Errorf("Get Asia/Shanghai time zone fail,%s", err.Error())
			return
		}
		ex := expireTime.In(loc).Format("2006-01-02 15:04:05")
		request.OldPasswordExpireTime = &ex
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().ModifyClusterPassword(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
	}
	return
}

func (me *TcaplusService) DescribeTask(ctx context.Context, applicationId string, taskId string) (taskInfo tcaplusdb.TaskInfoNew, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeTasksRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterIds = []*string{&applicationId}
	request.TaskIds = []*string{&taskId}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTcaplusClient().DescribeTasks(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
	}
	if len(response.Response.TaskInfos) == 0 {
		has = false
		return
	}
	if len(response.Response.TaskInfos) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d taskInfo with one taskId %s",
			len(response.Response.TaskInfos), taskId)
		return
	}
	has = true
	taskInfo = *response.Response.TaskInfos[0]
	return
}

func (me *TcaplusService) CreateZone(ctx context.Context, applicationId string, zoneName string) (zoneId string, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewCreateTableGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.TableGroupName = &zoneName
	request.ClusterId = &applicationId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().CreateTableGroup(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	if response.Response.TableGroupId == nil || *response.Response.TableGroupId == "" {
		errRet = errors.New("TencentCloud SDK  return empty applicationId")
		return
	}
	zoneId = *response.Response.TableGroupId
	return
}

func (me *TcaplusService) DescribeZones(ctx context.Context, applicationId string, zoneId, zoneName string) (infos []*tcaplusdb.TableGroupInfo, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeTableGroupsRequest()

	infos = make([]*tcaplusdb.TableGroupInfo, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if zoneId != "" {
		items := strings.Split(zoneId, ":")
		if len(items) != 2 {
			errRet = fmt.Errorf("zone id is broken,%s", zoneId)
			return
		}
		zoneId = items[1]
	}

	var offset, limit int64 = 0, 20

	request.ClusterId = &applicationId
	if zoneId != "" {
		request.TableGroupIds = []*string{&zoneId}
	}

	if zoneName != "" {
		request.Filters = []*tcaplusdb.Filter{
			{
				Name:  helper.String("ZoneName"),
				Value: &zoneName,
			},
		}
	}
	request.Offset = &offset
	request.Limit = &limit

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTcaplusClient().DescribeTableGroups(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		infos = append(infos, response.Response.TableGroups...)
		if len(response.Response.TableGroups) < int(limit) {
			return
		}
		offset += limit
	}

}

func (me *TcaplusService) DescribeZone(ctx context.Context, applicationId string, zoneId string) (info tcaplusdb.TableGroupInfo, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeTableGroupsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.TableGroupIds = []*string{&zoneId}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DescribeTableGroups(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceNotFound" {
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = errors.New("TencentCloud SDK return nil response")
		return
	}
	has = true
	if len(response.Response.TableGroups) == 0 {
		has = false
		return
	}
	if len(response.Response.TableGroups) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d zone info with one zoneId %s",
			len(response.Response.TableGroups), applicationId)
		return
	}
	info = *response.Response.TableGroups[0]
	return
}

func (me *TcaplusService) DeleteZone(ctx context.Context, applicationId string, zoneId string) (errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDeleteTableGroupRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.TableGroupId = &zoneId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DeleteTableGroup(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	return
}

func (me *TcaplusService) ModifyZoneName(ctx context.Context, applicationId string, zoneId, zoneName string) (errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewModifyTableGroupNameRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.TableGroupId = &zoneId
	request.TableGroupName = &zoneName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().ModifyTableGroupName(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	return
}

func (me *TcaplusService) DeleteIdlFiles(ctx context.Context, tid TcaplusIdlId) (errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDeleteIdlFilesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.ClusterId = &tid.ApplicationId
	request.IdlFiles = []*tcaplusdb.IdlFileInfo{
		{
			FileName:    &tid.FileName,
			FileType:    &tid.FileType,
			FileExtType: &tid.FileExtType,
			FileSize:    &tid.FileSize,
			FileId:      &tid.FileId,
		},
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DeleteIdlFiles(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}
	return
}

func (me *TcaplusService) DesOldIdlFiles(ctx context.Context, tid TcaplusIdlId) (tableInfos []*tcaplusdb.ParsedTableInfoNew, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewVerifyIdlFilesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.ClusterId = &tid.ApplicationId
	request.ExistingIdlFiles = []*tcaplusdb.IdlFileInfo{
		{
			FileName:    &tid.FileName,
			FileType:    &tid.FileType,
			FileExtType: &tid.FileExtType,
			FileSize:    &tid.FileSize,
			FileId:      &tid.FileId,
		},
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().VerifyIdlFiles(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceNotFound" {
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	tableInfos = response.Response.TableInfos
	return
}

func (me *TcaplusService) VerifyIdlFiles(ctx context.Context, tid TcaplusIdlId, zoneId string, fileContent string) (idlId int64, tableInfos []*tcaplusdb.ParsedTableInfoNew, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewVerifyIdlFilesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.ClusterId = &tid.ApplicationId
	request.TableGroupId = &zoneId
	request.NewIdlFiles = []*tcaplusdb.IdlFileInfo{
		{
			FileName:    &tid.FileName,
			FileType:    &tid.FileType,
			FileExtType: &tid.FileExtType,
			FileSize:    &tid.FileSize,
			FileContent: helper.String(base64.StdEncoding.EncodeToString([]byte(url.QueryEscape(fileContent)))),
		},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().VerifyIdlFiles(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.IdlFiles) == 0 {
		errRet = fmt.Errorf("get empty infos from this idl")
		return
	}

	if len(response.Response.IdlFiles) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d idl infos, but we only upload one idl", len(response.Response.IdlFiles))
		return
	}

	idlId = *response.Response.IdlFiles[0].FileId

	tableInfos = response.Response.TableInfos

	return
}

func (me *TcaplusService) CreateTables(ctx context.Context, tid TcaplusIdlId,
	applicationId,
	zoneId,
	tableName,
	tableType,
	description,
	tableIdlType string,
	reservedReadQps,
	reservedWriteQps,
	reservedVolume int64) (taskId string, tableInstanceId string, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewCreateTablesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.IdlFiles = []*tcaplusdb.IdlFileInfo{
		{
			FileName:    &tid.FileName,
			FileType:    &tid.FileType,
			FileExtType: &tid.FileExtType,
			FileSize:    &tid.FileSize,
			FileId:      &tid.FileId,
		},
	}
	request.SelectedTables = []*tcaplusdb.SelectedTableInfoNew{
		{
			TableIdlType:     &tableIdlType,
			ReservedReadQps:  &reservedReadQps,
			ReservedWriteQps: &reservedWriteQps,
			ReservedVolume:   &reservedVolume,
			TableGroupId:     &zoneId,
			TableName:        &tableName,
			TableType:        &tableType,
			Memo:             &description,
		},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().CreateTables(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.TableResults) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK  return empty table create result")
		return
	}

	if len(response.Response.TableResults) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d table create result, but we only upload one idl", len(response.Response.TableResults))
		return
	}

	if response.Response.TableResults[0].Error != nil {
		errRet = fmt.Errorf("TencentCloud SDK  return error,%s", *response.Response.TableResults[0].Error.Message)
		return
	}

	taskId = *response.Response.TableResults[0].TaskId
	tableInstanceId = *response.Response.TableResults[0].TableInstanceId
	return
}
func (me *TcaplusService) DescribeTables(ctx context.Context, applicationId string, zoneId, tableId, tableName string) (infos []*tcaplusdb.TableInfoNew, errRet error) {
	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeTablesRequest()

	infos = make([]*tcaplusdb.TableInfoNew, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if zoneId != "" {
		items := strings.Split(zoneId, ":")
		if len(items) != 2 {
			errRet = fmt.Errorf("zone id is broken,%s", zoneId)
			return
		}
		zoneId = items[1]
	}

	var offset, limit int64 = 0, 20

	request.ClusterId = &applicationId
	if zoneId != "" {
		request.TableGroupIds = []*string{&zoneId}
	}

	if tableId != "" {
		request.Filters = []*tcaplusdb.Filter{
			{
				Name:  helper.String("TableInstanceId"),
				Value: &tableId,
			},
		}
	}

	if tableName != "" {
		filter := &tcaplusdb.Filter{
			Name:  helper.String("TableName"),
			Value: &tableName,
		}

		if len(request.Filters) == 0 {
			request.Filters = append(request.Filters, filter)
		} else {
			request.Filters = []*tcaplusdb.Filter{filter}
		}

	}
	request.Offset = &offset
	request.Limit = &limit
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTcaplusClient().DescribeTables(request)

		if err != nil {
			if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
				if sdkErr.Code == "ResourceNotFound" {
					errRet = nil
					return
				}
			}
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		infos = append(infos, response.Response.TableInfos...)
		if len(response.Response.TableInfos) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TcaplusService) DescribeTable(ctx context.Context, applicationId, tableInstanceId string) (tableInfo tcaplusdb.TableInfoNew, has bool, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeTablesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.ClusterId = &applicationId
	request.Filters = []*tcaplusdb.Filter{
		{
			Name:  helper.String("TableInstanceId"),
			Value: &tableInstanceId,
		},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DescribeTables(request)
	if err != nil {
		if sdkErr, ok := err.(*sdkError.TencentCloudSDKError); ok {
			if sdkErr.Code == "ResourceNotFound" {
				errRet = nil
				return
			}
		}
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.TableInfos) == 0 {
		return
	}

	if len(response.Response.TableInfos) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d table result with one table id %s", len(response.Response.TableInfos), tableInstanceId)
		return
	}

	has = true
	tableInfo = *response.Response.TableInfos[0]

	return

}

func (me *TcaplusService) DeleteTable(ctx context.Context, applicationId, zoneId, tableInstanceId, tableName string) (taskId string, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDeleteTablesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.SelectedTables = []*tcaplusdb.SelectedTableInfoNew{
		{
			TableInstanceId: &tableInstanceId,
			TableGroupId:    &zoneId,
			TableName:       &tableName,
		},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().DeleteTables(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
		return
	}

	if len(response.Response.TableResults) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK return nil taskinfo,%s", request.GetAction())
		return
	}

	taskId = *response.Response.TableResults[0].TaskId
	return
}

func (me *TcaplusService) ModifyTableMemo(ctx context.Context, applicationId, zoneId, tableInstanceId, tableName, newDesc string) (errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewModifyTableMemosRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.TableMemos = []*tcaplusdb.SelectedTableInfoNew{
		{
			TableGroupId:    &zoneId,
			TableName:       &tableName,
			TableInstanceId: &tableInstanceId,
			Memo:            &newDesc,
		},
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcaplusClient().ModifyTableMemos(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
	}
	if len(response.Response.TableResults) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK return nil modify task infos,%s", request.GetAction())
		return
	}
	if len(response.Response.TableResults) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK return %d taskInfo with one op", len(response.Response.TableResults))
		return
	}

	return
}

func (me *TcaplusService) ModifyTables(ctx context.Context, tid TcaplusIdlId,
	applicationId,
	zoneId,
	tableInstanceId,
	tableName,
	tableIdType string) (taskId string, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewModifyTablesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		errRet = fmt.Errorf("zone id is broken,%s", zoneId)
		return
	}
	zoneId = items[1]

	request.ClusterId = &applicationId
	request.IdlFiles = []*tcaplusdb.IdlFileInfo{
		{
			FileName:    &tid.FileName,
			FileType:    &tid.FileType,
			FileExtType: &tid.FileExtType,
			FileSize:    &tid.FileSize,
			FileId:      &tid.FileId,
		},
	}
	request.SelectedTables = []*tcaplusdb.SelectedTableInfoNew{
		{
			TableInstanceId: &tableInstanceId,
			TableGroupId:    &zoneId,
			TableName:       &tableName,
			TableIdlType:    &tableIdType,
		},
	}
	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTcaplusClient().ModifyTables(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK modify task idl return nil response,%s", request.GetAction())
	}
	if len(response.Response.TableResults) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK modify task idl return nil,%s", request.GetAction())
		return
	}
	if len(response.Response.TableResults) != 1 {
		errRet = fmt.Errorf("TencentCloud SDK modify task idl return %d taskInfos with one op", len(response.Response.TableResults))
		return
	}
	taskId = *response.Response.TableResults[0].TaskId
	return
}

func (me *TcaplusService) DescribeIdlFileInfos(ctx context.Context, applicationId string) (infos []*tcaplusdb.IdlFileInfo, errRet error) {

	logId := getLogId(ctx)
	request := tcaplusdb.NewDescribeIdlFileInfosRequest()

	infos = make([]*tcaplusdb.IdlFileInfo, 0, 100)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit int64 = 0, 20

	request.ClusterId = &applicationId

	request.Offset = &offset
	request.Limit = &limit
	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTcaplusClient().DescribeIdlFileInfos(request)

		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response,%s", request.GetAction())
			return
		}
		infos = append(infos, response.Response.IdlFileInfos...)
		if len(response.Response.IdlFileInfos) < int(limit) {
			return
		}
		offset += limit
	}
}
