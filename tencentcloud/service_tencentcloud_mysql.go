package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MysqlService struct {
	client *connectivity.TencentCloudClient
}

// check if the err means the mysql_id is not found
func (me *MysqlService) NotFoundMysqlInstance(err error) bool {

	if err == nil {
		return false
	}

	sdkErr, ok := err.(*errors.TencentCloudSDKError)

	if ok {
		if sdkErr.Code == MysqlInstanceIdNotFound || sdkErr.Code == MysqlInstanceIdNotFound2 {
			return true
		}
	}
	return false
}

func (me *MysqlService) DescribeBackupsByMysqlId(ctx context.Context,
	mysqlId string,
	leftNumber int64) (backupInfos []*cdb.BackupInfo, errRet error) {

	logId := getLogId(ctx)

	listInitSize := leftNumber
	if listInitSize > 500 {
		listInitSize = 500
	}
	backupInfos = make([]*cdb.BackupInfo, 0, listInitSize)

	request := cdb.NewDescribeBackupsRequest()
	request.InstanceId = &mysqlId

needMoreItems:
	var limit int64 = 50
	if leftNumber < limit {
		limit = leftNumber
	}
	if leftNumber <= 0 {
		return
	}
	var offset int64 = 0
	request.Limit = &limit
	request.Offset = &offset
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMysqlClient().DescribeBackups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	totalCount := *response.Response.TotalCount
	leftNumber = leftNumber - limit
	offset += limit

	backupInfos = append(backupInfos, response.Response.Items...)
	if leftNumber > 0 && totalCount-offset > 0 {
		goto needMoreItems
	}
	return backupInfos, nil

}

func (me *MysqlService) CreateBackup(ctx context.Context, mysqlId string) (backupId int64, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewCreateBackupRequest()

	backupMethod := "logical"
	request.BackupMethod = &backupMethod
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().CreateBackup(request)
	if err != nil {
		errRet = err
		return
	}
	backupId = int64(*response.Response.BackupId)
	return
}

func (me *MysqlService) DescribeDBZoneConfig(ctx context.Context) (sellConfigures []*cdb.RegionSellConf, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeDBZoneConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBZoneConfig(request)
	if err != nil {
		errRet = err
		return
	}
	sellConfigures = response.Response.Items
	return
}

func (me *MysqlService) DescribeBackupConfigByMysqlId(ctx context.Context, mysqlId string) (desResponse *cdb.DescribeBackupConfigResponse, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeBackupConfigRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeBackupConfig(request)
	if err != nil {
		errRet = err
		return
	}
	desResponse = response
	return
}

func (me *MysqlService) ModifyBackupConfigByMysqlId(ctx context.Context, mysqlId string,
	retentionPeriod int64, backupModel, backupTime string) (errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewModifyBackupConfigRequest()
	request.InstanceId = &mysqlId
	request.ExpireDays = &retentionPeriod
	request.StartTime = &backupTime
	request.BackupMethod = &backupModel

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyBackupConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}
func (me *MysqlService) DescribeDefaultParameters(ctx context.Context, engineVersion string) (parameterList []*cdb.ParameterDetail, errRet error) {
	logId := getLogId(ctx)

	request := cdb.NewDescribeDefaultParamsRequest()
	request.EngineVersion = &engineVersion

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDefaultParams(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	parameterList = response.Response.Items
	return
}

func (me *MysqlService) DescribeInstanceParameters(ctx context.Context, instanceId string) (parameterList []*cdb.ParameterDetail, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewDescribeInstanceParamsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}

	parameterList = response.Response.Items
	return
}

func (me *MysqlService) ModifyInstanceParam(ctx context.Context, instanceId string, params map[string]string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyInstanceParamRequest()
	request.InstanceIds = []*string{&instanceId}
	request.ParamList = make([]*cdb.Parameter, 0, len(params))

	for k, v := range params {
		key := k
		value := v
		var param = cdb.Parameter{Name: &key, CurrentValue: &value}
		request.ParamList = append(request.ParamList, &param)
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyInstanceParam(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	asyncRequestId = *response.Response.AsyncRequestId

	return
}

func (me *MysqlService) DescribeCaresParameters(ctx context.Context, instanceId string, cares []string) (caresKv map[string]interface{}, errRet error) {
	caresKv = make(map[string]interface{})
	parameterList, err := me.DescribeInstanceParameters(ctx, instanceId)
	if err != nil {
		sdkErr, ok := err.(*errors.TencentCloudSDKError)
		if ok && sdkErr.Code == "CdbError" {
			return
		}
		errRet = err
		return
	}

	var inSlice = func(key string) bool {
		for _, care := range cares {
			if key == care {
				return true
			}
		}
		return false
	}

	for _, paramInfo := range parameterList {
		if inSlice(*paramInfo.Name) {
			caresKv[*paramInfo.Name] = *paramInfo.CurrentValue
		}
	}
	return
}

func (me *MysqlService) CreateAccount(ctx context.Context, mysqlId string,
	accountName, accountPassword, accountDescription string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewCreateAccountsRequest()

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	var accountInfos = []*cdb.Account{&accountInfo}

	request.InstanceId = &mysqlId
	request.Password = &accountPassword
	request.Accounts = accountInfos
	request.Description = &accountDescription

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().CreateAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) ModifyAccountPassword(ctx context.Context, mysqlId string,
	accountName, accountPassword string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyAccountPasswordRequest()

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	var accountInfos = []*cdb.Account{&accountInfo}

	request.InstanceId = &mysqlId
	request.NewPassword = &accountPassword
	request.Accounts = accountInfos

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyAccountPassword(request)
	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) ModifyAccountDescription(ctx context.Context, mysqlId string,
	accountName, accountDescription string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyAccountDescriptionRequest()

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	var accountInfos = []*cdb.Account{&accountInfo}

	request.InstanceId = &mysqlId
	request.Description = &accountDescription
	request.Accounts = accountInfos

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyAccountDescription(request)
	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) DeleteAccount(ctx context.Context, mysqlId string,
	accountName string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewDeleteAccountsRequest()

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	var accountInfos = []*cdb.Account{&accountInfo}

	request.InstanceId = &mysqlId
	request.Accounts = accountInfos

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DeleteAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) DescribeAccounts(ctx context.Context, mysqlId string) (accountInfos []*cdb.AccountInfo, errRet error) {

	logId := getLogId(ctx)

	var (
		listInitSize int64 = 100
		limit        int64 = 100
		offset       int64 = 0
		leftNumbers  int64 = 0
		dofirst            = true
	)

	accountInfos = make([]*cdb.AccountInfo, 0, listInitSize)
	request := cdb.NewDescribeAccountsRequest()

	request.InstanceId = &mysqlId
	request.Offset = &offset
	request.Limit = &limit

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

needMoreItems:
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	if dofirst {
		leftNumbers = *response.Response.TotalCount - limit
		dofirst = false
	} else {
		leftNumbers = leftNumbers - limit
	}
	offset = offset + limit

	accountInfos = append(accountInfos, response.Response.Items...)

	if leftNumbers > 0 {
		goto needMoreItems
	} else {
		return
	}

}

func (me *MysqlService) _innerDescribeAsyncRequestInfo(ctx context.Context, asyncRequestId string) (status, message string, errRet error) {
	logId := getLogId(ctx)
	request := cdb.NewDescribeAsyncRequestInfoRequest()
	request.AsyncRequestId = &asyncRequestId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeAsyncRequestInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	status = *response.Response.Status
	message = *response.Response.Info
	return
}

func (me *MysqlService) DescribeAsyncRequestInfo(ctx context.Context, asyncRequestId string) (status, message string, errRet error) {

	// Post https://cdb.tencentcloudapi.com/:  always get "Gateway Time-out"
	status, message, errRet = me._innerDescribeAsyncRequestInfo(ctx, asyncRequestId)
	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			status, message, errRet = me._innerDescribeAsyncRequestInfo(ctx, asyncRequestId)
		}
	}
	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			time.Sleep(2 * time.Second)
			status, message, errRet = me._innerDescribeAsyncRequestInfo(ctx, asyncRequestId)
		}
	}
	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			time.Sleep(5 * time.Second)
			status, message, errRet = me._innerDescribeAsyncRequestInfo(ctx, asyncRequestId)
		}
	}
	return
}

func (me *MysqlService) ModifyAccountPrivileges(ctx context.Context, mysqlId string,
	accountName string, databaseNames []string, privileges []string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewModifyAccountPrivilegesRequest()
	request.InstanceId = &mysqlId

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	request.Accounts = []*cdb.Account{&accountInfo}

	request.DatabasePrivileges = make([]*cdb.DatabasePrivilege, 0, len(databaseNames))

	for _, databaseName := range databaseNames {
		var temp = databaseName
		var cdbprivileges = cdb.DatabasePrivilege{Database: &temp}
		cdbprivileges.Privileges = make([]*string, len(privileges))

		for i := range privileges {
			cdbprivileges.Privileges[i] = &privileges[i]
		}

		request.DatabasePrivileges = append(request.DatabasePrivileges, &cdbprivileges)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyAccountPrivileges(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) DescribeAccountPrivileges(ctx context.Context, mysqlId string,
	accountName string, databaseNames []string) (privileges []string, errRet error) {

	logId := getLogId(ctx)

	privileges = make([]string, 0, len(MYSQL_DATABASE_PRIVILEGE))

	request := cdb.NewDescribeAccountPrivilegesRequest()
	request.InstanceId = &mysqlId
	request.User = &accountName
	request.Host = &MYSQL_DEFAULT_ACCOUNT_HOST

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeAccountPrivileges(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	var inSlice = func(str string, strs []string) bool {
		for _, v := range strs {
			if v == str {
				return true
			}
		}
		return false
	}

	privilegeCountMap := make(map[string]int)

	hasMapSize := 0
	for _, dataPrivilege := range response.Response.DatabasePrivileges {

		if inSlice(*dataPrivilege.Database, databaseNames) {

			hasMapSize++

			for _, privilege := range dataPrivilege.Privileges {
				privilegeCountMap[*privilege]++
			}

		}
	}
	// every exist database all has the privilege
	for privilege, scount := range privilegeCountMap {
		if scount == hasMapSize {
			privileges = append(privileges, privilege)
		}
	}

	log.Printf("[DEBUG]%s we got same privileges is %+v \n", logId, privileges)
	return
}

func (me *MysqlService) DescribeDBInstanceById(ctx context.Context, mysqlId string) (mysqlInfo *cdb.InstanceInfo, errRet error) {

	// Post https://cdb.tencentcloudapi.com/:  always get "Gateway Time-out"
	mysqlInfo, errRet = me._innerDescribeDBInstanceById(ctx, mysqlId)

	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			mysqlInfo, errRet = me._innerDescribeDBInstanceById(ctx, mysqlId)
		}
	}
	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			time.Sleep(3 * time.Second)
			mysqlInfo, errRet = me._innerDescribeDBInstanceById(ctx, mysqlId)

		}
	}
	if errRet != nil {
		if _, ok := errRet.(*errors.TencentCloudSDKError); !ok {
			time.Sleep(5 * time.Second)
			mysqlInfo, errRet = me._innerDescribeDBInstanceById(ctx, mysqlId)
		}
	}
	return
}

func (me *MysqlService) DescribeIsolatedDBInstanceById(ctx context.Context, mysqlId string) (mysqlInfo *cdb.InstanceInfo, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}

	request.Status = []*uint64{helper.Uint64(MYSQL_STATUS_ISOLATED),
		helper.Uint64(MYSQL_STATUS_ISOLATED_1),
		helper.Uint64(MYSQL_STATUS_ISOLATED_2)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) == 0 {
		return
	}
	if len(response.Response.Items) > 1 {
		errRet = fmt.Errorf("One mysql id got %d instance info", len(response.Response.Items))
		return
	}
	mysqlInfo = response.Response.Items[0]

	return
}

func (me *MysqlService) _innerDescribeDBInstanceById(ctx context.Context, mysqlId string) (mysqlInfo *cdb.InstanceInfo, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) == 0 {
		return
	}
	if len(response.Response.Items) > 1 {
		errRet = fmt.Errorf("One mysql id got %d instance info", len(response.Response.Items))
	}
	mysqlInfo = response.Response.Items[0]

	return
}

func (me *MysqlService) DescribeRunningDBInstanceById(ctx context.Context, mysqlId string) (mysqlInfo *cdb.InstanceInfo, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}
	runningStatus := uint64(1)
	request.Status = []*uint64{&runningStatus}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) == 0 {
		return
	}
	if len(response.Response.Items) > 1 {
		errRet = fmt.Errorf("One mysql id got %d instance info", len(response.Response.Items))
	}
	mysqlInfo = response.Response.Items[0]

	return
}

func (me *MysqlService) CheckDBGTIDOpen(ctx context.Context, mysqlId string) (open int64, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeDBInstanceGTIDRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBInstanceGTID(request)
	if err != nil {
		sdkErr, ok := err.(*errors.TencentCloudSDKError)
		if ok && sdkErr.Code == "CdbError" {
			open = 0
			return
		}
		errRet = err
		return
	}
	open = *response.Response.IsGTIDOpen
	return
}

func (me *MysqlService) DescribeDBSecurityGroups(ctx context.Context, mysqlId string) (securityGroups []string, errRet error) {
	logId := getLogId(ctx)
	request := cdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &mysqlId
	securityGroups = make([]string, 0, 10)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, sg := range response.Response.Groups {
		securityGroups = append(securityGroups, *sg.SecurityGroupId)
	}
	return
}

func (me *MysqlService) ModifyInstanceTag(ctx context.Context, mysqlId string, deleteTags, modifyTags map[string]string) (errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyInstanceTagRequest()
	request.InstanceId = &mysqlId

	if len(modifyTags) > 0 {
		request.ReplaceTags = make([]*cdb.TagInfo, 0, len(modifyTags))
		for name, value := range modifyTags {
			tagKey := name
			tagValue := value
			tag := cdb.TagInfo{TagKey: &tagKey, TagValue: []*string{&tagValue}}
			request.ReplaceTags = append(request.ReplaceTags, &tag)
		}
	}
	if len(deleteTags) > 0 {
		request.DeleteTags = make([]*cdb.TagInfo, 0, len(deleteTags))
		for name, value := range deleteTags {
			tagKey := name
			tagValue := value
			tag := cdb.TagInfo{TagKey: &tagKey, TagValue: []*string{&tagValue}}
			request.DeleteTags = append(request.DeleteTags, &tag)
		}
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyInstanceTag(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MysqlService) DescribeTagsOfInstanceId(ctx context.Context, mysqlId string) (tags map[string]string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewDescribeTagsOfInstanceIdsRequest()
	request.InstanceIds = []*string{&mysqlId}
	tags = make(map[string]string)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	var (
		limit  int64 = 10
		offset int64 = 0
	)
	request.Limit = &limit

again:
	if request.Offset == nil {
		request.Offset = &offset
	} else {
		offset = offset + limit
	}
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeTagsOfInstanceIds(request)
	if err != nil {
		errRet = err
		return
	}
	if len(response.Response.Rows) == 0 {
		return
	}
	if len(response.Response.Rows) > 1 {
		errRet = fmt.Errorf("One mysql id got %d tags info rows", len(response.Response.Rows))
	}
	if len(response.Response.Rows[0].Tags) == 0 {
		return
	}
	for _, tag := range response.Response.Rows[0].Tags {
		if _, has := tags[*tag.TagKey]; has {
			return
		}
		tags[*tag.TagKey] = *tag.TagValue
	}

	goto again
}

func (me *MysqlService) DescribeDBInstanceConfig(ctx context.Context, mysqlId string) (backupConfig *cdb.DescribeDBInstanceConfigResponse,
	errRet error) {
	logId := getLogId(ctx)
	request := cdb.NewDescribeDBInstanceConfigRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DescribeDBInstanceConfig(request)

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	backupConfig = response

	return
}

func (me *MysqlService) InitDBInstances(ctx context.Context, mysqlId string, password string) (asyncRequestId string, errRet error) {
	logId := getLogId(ctx)
	request := cdb.NewInitDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}
	request.NewPassword = &password

	paramsMap := map[string]string{
		"character_set_server":   "utf8mb4", // ["utf8","latin1","gbk","utf8mb4"]
		"lower_case_table_names": "1",       // ["0","1"]
	}
	request.Parameters = make([]*cdb.ParamInfo, 0, len(paramsMap))
	for k, v := range paramsMap {
		name := k
		value := v
		param := cdb.ParamInfo{Name: &name, Value: &value}
		request.Parameters = append(request.Parameters, &param)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().InitDBInstances(request)

	if err != nil {
		errRet = err
		return
	}
	if len(response.Response.AsyncRequestIds) != 1 {
		errRet = fmt.Errorf("init one  mysql id got %d async ids", len(response.Response.AsyncRequestIds))
		return
	}

	asyncRequestId = *response.Response.AsyncRequestIds[0]
	return
}

func (me *MysqlService) OpenWanService(ctx context.Context, mysqlId string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewOpenWanServiceRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().OpenWanService(request)

	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) CloseWanService(ctx context.Context, mysqlId string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewCloseWanServiceRequest()
	request.InstanceId = &mysqlId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().CloseWanService(request)

	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) OpenDBInstanceGTID(ctx context.Context, mysqlId string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewOpenDBInstanceGTIDRequest()
	request.InstanceId = &mysqlId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().OpenDBInstanceGTID(request)

	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) ModifyDBInstanceName(ctx context.Context, mysqlId,
	newInstanceName string) (errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewModifyDBInstanceNameRequest()
	request.InstanceId = &mysqlId
	request.InstanceName = &newInstanceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, errRet := me.client.UseMysqlClient().ModifyDBInstanceName(request)

	if errRet != nil {
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *MysqlService) ModifyDBInstanceVipVport(ctx context.Context, mysqlId, vpcId, subnetId string, port int64) (errRet error) {
	logId := getLogId(ctx)
	request := cdb.NewModifyDBInstanceVipVportRequest()
	request.InstanceId = &mysqlId
	request.DstPort = &port
	if vpcId != "" {
		request.UniqVpcId = &vpcId
	}
	if subnetId != "" {
		request.UniqSubnetId = &subnetId
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, errRet := me.client.UseMysqlClient().ModifyDBInstanceVipVport(request)

	if errRet != nil {
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MysqlService) UpgradeDBInstance(ctx context.Context, mysqlId string,
	memSize, volumeSize int64) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)

	var waitSwitch int64 = 0 // 0- switch immediately, 1- time window switch

	request := cdb.NewUpgradeDBInstanceRequest()
	request.InstanceId = &mysqlId
	request.Memory = &memSize
	request.Volume = &volumeSize
	request.WaitSwitch = &waitSwitch

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().UpgradeDBInstance(request)

	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) ModifyDBInstanceProject(ctx context.Context, mysqlId string, newProjectId int64) (errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyDBInstanceProjectRequest()
	request.InstanceIds = []*string{&mysqlId}
	request.NewProjectId = &newProjectId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyDBInstanceProject(request)

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return

}

func (me *MysqlService) ModifyDBInstanceSecurityGroups(ctx context.Context, mysqlId string, securityGroups []string) (errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewModifyDBInstanceSecurityGroupsRequest()
	request.InstanceId = &mysqlId
	request.SecurityGroupIds = make([]*string, 0, len(securityGroups))

	for _, v := range securityGroups {
		value := v
		request.SecurityGroupIds = append(request.SecurityGroupIds, &value)
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyDBInstanceSecurityGroups(request)

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MysqlService) DisassociateSecurityGroup(ctx context.Context, mysqlId string, securityGroup string) (errRet error) {

	logId := getLogId(ctx)

	request := cdb.NewDisassociateSecurityGroupsRequest()
	request.InstanceIds = []*string{&mysqlId}
	request.SecurityGroupId = &securityGroup

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().DisassociateSecurityGroups(request)

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return

}

func (me *MysqlService) ModifyAutoRenewFlag(ctx context.Context, mysqlId string, newRenewFlag int64) (errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewModifyAutoRenewFlagRequest()
	request.InstanceIds = []*string{&mysqlId}
	request.AutoRenew = &newRenewFlag

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().ModifyAutoRenewFlag(request)

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *MysqlService) IsolateDBInstance(ctx context.Context, mysqlId string) (asyncRequestId string, errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewIsolateDBInstanceRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMysqlClient().IsolateDBInstance(request)

	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	//The server returns that AsyncRequestId does not exist
	//asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) OfflineIsolatedInstances(ctx context.Context, mysqlId string) (errRet error) {

	logId := getLogId(ctx)
	request := cdb.NewOfflineIsolatedInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	_, errRet = me.client.UseMysqlClient().OfflineIsolatedInstances(request)

	return
}
