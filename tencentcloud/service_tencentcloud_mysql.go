package tencentcloud

import (
	"context"
	"fmt"
	"log"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type MysqlService struct {
	client *connectivity.TencentCloudClient
}

//check if the err means the mysql_id is not found
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

	logId := GetLogId(ctx)

	listInitSize := leftNumber
	if listInitSize > 500 {
		listInitSize = 500
	}
	backupInfos = make([]*cdb.BackupInfo, 0, listInitSize)

	request := cdb.NewDescribeBackupsRequest()
	request.InstanceId = &mysqlId

needMoreItems:
	var limit int64 = 10
	if leftNumber > limit {
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

func (me *MysqlService) DescribeDBZoneConfig(ctx context.Context) (sellConfigures []*cdb.RegionSellConf, errRet error) {

	logId := GetLogId(ctx)
	request := cdb.NewDescribeDBZoneConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseMysqlClient().DescribeDBZoneConfig(request)
	if err != nil {
		errRet = err
		return
	}
	sellConfigures = response.Response.Items
	return
}

func (me *MysqlService) DescribeBackupConfigByMysqlId(ctx context.Context, mysqlId string) (desResponse *cdb.DescribeBackupConfigResponse, errRet error) {

	logId := GetLogId(ctx)
	request := cdb.NewDescribeBackupConfigRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

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

	logId := GetLogId(ctx)
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
	logId := GetLogId(ctx)

	request := cdb.NewDescribeDefaultParamsRequest()
	request.EngineVersion = &engineVersion

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

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

	logId := GetLogId(ctx)

	request := cdb.NewDescribeInstanceParamsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseMysqlClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	parameterList = response.Response.Items
	return
}

func (me *MysqlService) DescribeCaresParameters(ctx context.Context, instanceId string, cares []string) (caresKv map[string]interface{}, errRet error) {
	caresKv = make(map[string]interface{})
	parameterList, err := me.DescribeInstanceParameters(ctx, instanceId)
	if err != nil {
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

	logId := GetLogId(ctx)

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

	logId := GetLogId(ctx)

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

	logId := GetLogId(ctx)

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

	logId := GetLogId(ctx)

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

	response, err := me.client.UseMysqlClient().DeleteAccounts(request)
	if err != nil {
		errRet = err
		return
	}
	asyncRequestId = *response.Response.AsyncRequestId
	return
}

func (me *MysqlService) DescribeAccounts(ctx context.Context, mysqlId string) (accountInfos []*cdb.AccountInfo, errRet error) {

	logId := GetLogId(ctx)

	var (
		listInitSize int64 = 100
		limit        int64 = 100
		offset       int64 = 0
		leftNumbers  int64 = 0
		dofirst      bool  = true
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

func (me *MysqlService) DescribeAsyncRequestInfo(ctx context.Context, asyncRequestId string) (status, message string, errRet error) {
	logId := GetLogId(ctx)
	request := cdb.NewDescribeAsyncRequestInfoRequest()
	request.AsyncRequestId = &asyncRequestId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

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

func (me *MysqlService) ModifyAccountPrivileges(ctx context.Context, mysqlId string,
	accountName string, databaseNames []string, privileges []string) (asyncRequestId string, errRet error) {

	logId := GetLogId(ctx)
	request := cdb.NewModifyAccountPrivilegesRequest()
	request.InstanceId = &mysqlId

	var accountInfo = cdb.Account{User: &accountName, Host: &MYSQL_DEFAULT_ACCOUNT_HOST}
	request.Accounts = []*cdb.Account{&accountInfo}

	request.DatabasePrivileges = make([]*cdb.DatabasePrivilege, 0, len(databaseNames))

	for _, databaseName := range databaseNames {
		var temp = databaseName
		var cdbprivileges = cdb.DatabasePrivilege{Database: &temp}
		cdbprivileges.Privileges = make([]*string, len(privileges))

		for i, _ := range privileges {
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

	logId := GetLogId(ctx)

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
	//every exist database all has the privilege
	for privilege, scount := range privilegeCountMap {
		if scount == hasMapSize {
			privileges = append(privileges, privilege)
		}
	}

	log.Printf("[DEBUG]%s we got same privileges is %+v \n", logId, privileges)
	return
}

func (me *MysqlService) DescribeDBInstanceById(ctx context.Context, mysqlId string) (mysqlInfo *cdb.InstanceInfo, errRet error) {

	logId := GetLogId(ctx)
	request := cdb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&mysqlId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

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

	logId := GetLogId(ctx)
	request := cdb.NewDescribeDBInstanceGTIDRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	response, err := me.client.UseMysqlClient().DescribeDBInstanceGTID(request)
	if err != nil {
		errRet = err
		return
	}
	open = *response.Response.IsGTIDOpen
	return
}

func (me *MysqlService) DescribeDBSecurityGroups(ctx context.Context, mysqlId string) (securityGroups []string, errRet error) {
	logId := GetLogId(ctx)
	request := cdb.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &mysqlId
	securityGroups = make([]string, 0, 10)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	response, err := me.client.UseMysqlClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	for _, sg := range response.Response.Groups {
		securityGroups = append(securityGroups, *sg.SecurityGroupId)
	}
	return
}

func (me *MysqlService) ModifyInstanceTags() {}

func (me *MysqlService) DescribeTagsOfInstanceId(ctx context.Context, mysqlId string) (tags map[string]string, errRet error) {

	logId := GetLogId(ctx)
	request := cdb.NewDescribeTagsOfInstanceIdsRequest()
	request.InstanceIds = []*string{&mysqlId}
	tags = make(map[string]string)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

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

	for _, tag := range response.Response.Rows[0].Tags {
		tags[*tag.TagKey] = *tag.TagValue
	}
	return
}

func (me *MysqlService) DescribeDBInstanceConfig(ctx context.Context, mysqlId string) (backupConfig *cdb.DescribeDBInstanceConfigResponse,
	errRet error) {
	logId := GetLogId(ctx)
	request := cdb.NewDescribeDBInstanceConfigRequest()
	request.InstanceId = &mysqlId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
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
