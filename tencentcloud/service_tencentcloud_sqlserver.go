package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pkg/errors"
	SDKErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	sqlserver "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sqlserver/v20180328"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type SqlserverService struct {
	client *connectivity.TencentCloudClient
}

func (me *SqlserverService) DescribeZones(ctx context.Context) (zoneInfoList []*sqlserver.ZoneInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeZonesRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.DescribeZonesResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DescribeZones(request)
		if e != nil {
			log.Printf("[CRITAL]%s DescribeZones fail, reason:%s\n", logId, e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		zoneInfoList = response.Response.ZoneSet
	}
	return
}

func (me *SqlserverService) DescribeProductConfig(ctx context.Context, zone string) (specInfoList []*sqlserver.SpecInfo, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeProductConfigRequest()
	request.Zone = &zone

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.DescribeProductConfigResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DescribeProductConfig(request)
		if e != nil {
			log.Printf("[CRITAL]%s DescribeProductConfig fail, reason:%s\n", logId, e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}
	if response != nil && response.Response != nil {
		specInfoList = response.Response.SpecInfoList
	}
	return
}

func (me *SqlserverService) CreateSqlserverInstance(ctx context.Context, request *sqlserver.CreateDBInstancesRequest) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().CreateDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	if len(response.Response.DealNames) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK returns empty SQL Server ID")
		return
	} else if len(response.Response.DealNames) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK returns more than one SQL Server ID")
		return
	}

	dealId := *response.Response.DealNames[0]
	instanceId, err = me.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		errRet = err
	}
	return
}

func (me *SqlserverService) ModifySqlserverInstanceName(ctx context.Context, instanceId string, name string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyDBInstanceNameRequest()
	request.InstanceId = &instanceId
	request.InstanceName = &name
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ModifyDBInstanceName(request)
	return err
}

func (me *SqlserverService) ModifySqlserverInstanceProjectId(ctx context.Context, instanceId string, projectId int) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyDBInstanceProjectRequest()
	request.InstanceIdSet = []*string{&instanceId}
	request.ProjectId = helper.IntInt64(projectId)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ModifyDBInstanceProject(request)
	return err
}

func (me *SqlserverService) UpgradeSqlserverInstance(ctx context.Context, instanceId string, memory, storage, autoVoucher int, voucherIds []*string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewUpgradeDBInstanceRequest()
	request.InstanceId = &instanceId
	request.Memory = helper.IntInt64(memory)
	request.Storage = helper.IntInt64(storage)
	request.WaitSwitch = helper.IntInt64(0)
	if autoVoucher != 0 {
		request.AutoVoucher = helper.IntInt64(autoVoucher)
	}
	if len(voucherIds) > 0 {
		request.VoucherIds = voucherIds
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().UpgradeDBInstance(request)
	if err != nil {
		return err
	}

	startPending := false
	//check status not expanding
	errRet = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribeSqlserverInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cannot find SQL Server instance %s", instanceId))
		}
		if *instance.Status != 2 {
			startPending = true
			return resource.RetryableError(fmt.Errorf("expanding , SQL Server instance ID %s, status %d.... ", instanceId, *instance.Status))
		} else if !startPending {
			return resource.RetryableError(fmt.Errorf("ready for expanding, SQL Server instance ID %s, status %d.... ", instanceId, *instance.Status))
		}
		return nil
	})

	return
}

func (me *SqlserverService) RemoveSecurityGroup(ctx context.Context, instanceId string, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDisassociateSecurityGroupsRequest()
	request.InstanceIdSet = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().DisassociateSecurityGroups(request)

	return err
}

func (me *SqlserverService) AddSecurityGroup(ctx context.Context, instanceId string, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewAssociateSecurityGroupsRequest()
	request.InstanceIdSet = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().AssociateSecurityGroups(request)
	time.Sleep(10 * time.Second)
	return err
}

func (me *SqlserverService) ModifySqlserverInstanceMaintenanceSpan(ctx context.Context, instanceId string, weekSet []int, startTime string, timeSpan int) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyMaintenanceSpanRequest()
	request.InstanceId = &instanceId
	if len(weekSet) > 0 {
		request.Weekly = make([]*int64, 0)
		for _, i := range weekSet {
			request.Weekly = append(request.Weekly, helper.IntInt64(i))
		}
	}
	if startTime != "" {
		request.StartTime = &startTime
	}
	if timeSpan != 0 {
		request.Span = helper.IntUint64(timeSpan)
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ModifyMaintenanceSpan(request)

	return err
}

func (me *SqlserverService) TerminateSqlserverInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewTerminateDBInstanceRequest()
	request.InstanceIdSet = []*string{&instanceId}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().TerminateDBInstance(request)
	return err
}

func (me *SqlserverService) DeleteSqlserverInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().DeleteDBInstance(request)
	return err
}

func (me *SqlserverService) DescribeSqlserverInstances(ctx context.Context, instanceId, instanceName string, projectId int, vpcId, subnetId string, netType int) (instanceList []*sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if instanceId != "" {
		request.InstanceIdSet = []*string{&instanceId}
	}
	if instanceName != "" {
		request.InstanceNameSet = []*string{&instanceName}
	}
	if projectId != -1 {
		request.ProjectId = helper.IntUint64(projectId)
	}
	if subnetId != "" && netType != BASIC_NETWORK {
		request.SubnetId = &subnetId
	}
	if vpcId != "" && netType != BASIC_NETWORK {
		request.VpcId = &vpcId
	}

	if netType == BASIC_NETWORK {
		//basic network
		request.VpcId = helper.String("")
		request.SubnetId = helper.String("")
	}
	var offset, limit int64 = 0, 20

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		instanceList = append(instanceList, response.Response.DBInstances...)
		if len(response.Response.DBInstances) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeSqlserverInstanceById(ctx context.Context, instanceId string) (instance *sqlserver.DBInstance, has bool, errRet error) {
	instanceList, err := me.DescribeSqlserverInstances(ctx, instanceId, "", -1, "", "", 1)
	if err != nil {
		errRet = err
		return
	}
	if len(instanceList) == 0 {
		return
	} else if len(instanceList) > 1 {
		errRet = fmt.Errorf("[DescribeDBInstanceById]SDK returns more than one instance with instanceId %s", instanceId)
	}

	instance = instanceList[0]
	if instance != nil && *instance.Status != 8 && *instance.Status != 4 && *instance.Status != 6 {
		has = true
	}
	return
}

func (me *SqlserverService) DescribeMaintenanceSpan(ctx context.Context, instanceId string) (weekSet []int, startTime string, timeSpan int, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeMaintenanceSpanRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeMaintenanceSpan(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	for _, v := range response.Response.Weekly {
		weekSet = append(weekSet, int(*v))
	}
	startTime = *response.Response.StartTime
	timeSpan = int(*response.Response.Span)

	return
}

func (me *SqlserverService) DescribeInstanceSecurityGroups(ctx context.Context, instanceId string) (securityGroups []string, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	securityGroups = make([]string, 0, len(response.Response.SecurityGroupSet))
	for _, v := range response.Response.SecurityGroupSet {
		securityGroups = append(securityGroups, *v.SecurityGroupId)
	}

	return
}

func (me *SqlserverService) DescribeSqlserverBackups(ctx context.Context, instanceId, backupName string, startTime string, endTime string) (backupList []*sqlserver.Backup, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeBackupsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.StartTime = &startTime
	request.EndTime = &endTime
	if backupName != "" {
		request.BackupName = &backupName
	}

	var offset, limit int64 = 0, 20

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeBackups(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		backupList = append(backupList, response.Response.Backups...)
		if len(response.Response.Backups) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeReadonlyGroupList(ctx context.Context, instanceId string) (groupList []*sqlserver.ReadOnlyGroup, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeReadOnlyGroupListRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeReadOnlyGroupList(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	groupList = response.Response.ReadOnlyGroupSet

	return
}

func (me *SqlserverService) CreateSqlserverReadonlyInstance(ctx context.Context, request *sqlserver.CreateReadOnlyDBInstancesRequest) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().CreateReadOnlyDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	if len(response.Response.DealNames) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK returns empty SQL Server ID")
		return
	} else if len(response.Response.DealNames) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK returns more than one SQL Server ID")
		return
	}

	dealId := *response.Response.DealNames[0]

	instanceId, err = me.GetInfoFromDeal(ctx, dealId)

	if err != nil {
		errRet = err
	}
	return
}

func (me *SqlserverService) DescribeReadonlyGroupListByReadonlyInstanceId(ctx context.Context, instanceId string) (readonlyInstance *sqlserver.DescribeReadOnlyGroupByReadOnlyInstanceResponseParams, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeReadOnlyGroupByReadOnlyInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeReadOnlyGroupByReadOnlyInstance(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	readonlyInstance = response.Response
	return
}

func (me *SqlserverService) CreateSqlserverAccount(ctx context.Context, instanceId string, userName string, password string, remark string, isAdmin bool) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewCreateAccountRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	account := sqlserver.AccountCreateInfo{UserName: &userName, Password: &password, IsAdmin: &isAdmin}
	if remark != "" {
		account.Remark = &remark
	}
	request.Accounts = []*sqlserver.AccountCreateInfo{&account}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().CreateAccount(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	flowId := *response.Response.FlowId
	err = me.WaitForTaskFinish(ctx, flowId)
	if err != nil {
		errRet = err
	}
	return
}

func (me *SqlserverService) DescribeSqlserverAccounts(ctx context.Context, instanceId string) (accounts []*sqlserver.AccountDetail, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeAccountsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId

	var offset, limit uint64 = 0, 20

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeAccounts(request)
		if err != nil {
			ee, ok := err.(*SDKErrors.TencentCloudSDKError)
			if !ok || ee.Code != "ResourceNotFound.InstanceNotFound" {
				errRet = err
				return
			}

			return

		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		accounts = append(accounts, response.Response.Accounts...)
		if len(response.Response.Accounts) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeSqlserverAccountById(ctx context.Context, instanceId string, userName string) (account *sqlserver.AccountDetail, has bool, errRet error) {
	accountList, err := me.DescribeSqlserverAccounts(ctx, instanceId)
	if err != nil {
		errRet = err
		return
	}
	if len(accountList) == 0 {
		return
	}

	for _, v := range accountList {
		if *v.Name == userName {
			account = v
			has = true
			return
		}
	}
	return
}

func (me *SqlserverService) ModifySqlserverAccountRemark(ctx context.Context, instanceId string, userName string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyAccountRemarkRequest()
	request.InstanceId = &instanceId
	request.Accounts = []*sqlserver.AccountRemark{{UserName: &userName, Remark: &remark}}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ModifyAccountRemark(request)
	return err
}

func (me *SqlserverService) ResetSqlserverAccountPassword(ctx context.Context, instanceId string, userName string, password string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewResetAccountPasswordRequest()
	request.InstanceId = &instanceId
	request.Accounts = []*sqlserver.AccountPassword{{UserName: &userName, Password: &password}}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ResetAccountPassword(request)
	if err != nil {
		errRet = err
		return
	}

	//check status not resetting
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribeSqlserverAccountById(ctx, instanceId, userName)
		if err != nil {
			return resource.NonRetryableError(errors.WithStack(err))
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cannot find SQL Server account %s%s%s", instanceId, FILED_SP, userName))
		}
		if int(*instance.Status) == 4 {
			return resource.RetryableError(fmt.Errorf("resetting , SQL Server instance ID %s, name %s, status %d.... ", instanceId, userName, *instance.Status))
		} else {
			return nil
		}
	})

	return
}

func (me *SqlserverService) DeleteSqlserverAccount(ctx context.Context, instanceId string, userName string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteAccountRequest()
	request.UserNames = []*string{&userName}
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().DeleteAccount(request)
	if err != nil {
		ee, ok := err.(*SDKErrors.TencentCloudSDKError)
		if !ok || ee.Code != "ResourceNotFound.InstanceNotFound" {
			errRet = err
			return
		}
		return
	}

	//check status not deleting
	errRet = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribeSqlserverAccountById(ctx, instanceId, userName)
		if err != nil {
			return resource.NonRetryableError(errors.WithStack(err))
		}
		if !has {
			return nil
		}
		if int(*instance.Status) == -1 {
			return resource.RetryableError(fmt.Errorf("deleting , SQL Server instance ID %s, name %s, status %d.... ", instanceId, userName, *instance.Status))
		} else {
			return resource.NonRetryableError(fmt.Errorf("invalid, SQL Server instance ID %s, name %s, status %d...", instanceId, userName, *instance.Status))
		}
	})

	return
}

func (me *SqlserverService) ModifyAccountDBAttachment(ctx context.Context, instanceId, accountName, dbName, privilege string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyAccountPrivilegeRequest()
	request.InstanceId = &instanceId
	request.Accounts = []*sqlserver.AccountPrivilegeModifyInfo{{UserName: &accountName, DBPrivileges: []*sqlserver.DBPrivilegeModifyInfo{{DBName: &dbName, Privilege: &privilege}}}}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	//check account exists
	_, has, err := me.DescribeSqlserverAccountById(ctx, instanceId, accountName)

	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("SQL Server account %s , instance ID %s is not exist", accountName, instanceId)
	}

	//check db exists
	_, has, err = me.DescribeDBDetailsById(ctx, fmt.Sprintf("%s%s%s", instanceId, FILED_SP, dbName))
	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("SQL Server DB %s , instance ID %s is not exist", dbName, instanceId)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().ModifyAccountPrivilege(request)
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	if err != nil {
		return err
	}

	flowId := int64(*response.Response.FlowId)
	err = me.WaitForTaskFinish(ctx, flowId)
	return err
}

func (me *SqlserverService) DescribeAccountDBAttachments(ctx context.Context, instanceId, accountName, dbName string) (attachments []map[string]string, errRet error) {

	if accountName != "" {
		//check account exists
		accounts, has, err := me.DescribeSqlserverAccountById(ctx, instanceId, accountName)
		if err != nil {
			errRet = err
			return
		}
		if !has {
			return
		}

		for _, v := range accounts.Dbs {

			if (dbName != "" && *v.DBName == dbName) || dbName == "" {
				mapping := make(map[string]string)
				mapping["db_name"] = *v.DBName
				mapping["account_name"] = accountName
				mapping["privilege"] = *v.Privilege
				attachments = append(attachments, mapping)
			}
		}
		return
	} else {
		dbInfos, err := me.DescribeDBsOfInstance(ctx, instanceId)
		if err != nil {
			errRet = err
			return
		}
		if len(dbInfos) == 0 {
			return
		}
		for _, v := range dbInfos {
			if (dbName != "" && *v.Name == dbName) || dbName == "" {
				for _, vv := range v.Accounts {
					mapping := make(map[string]string)
					mapping["db_name"] = *v.Name
					mapping["account_name"] = *vv.UserName
					mapping["privilege"] = *vv.Privilege
					attachments = append(attachments, mapping)
				}
			}
		}
		return
	}
}

func (me *SqlserverService) DescribeAccountDBAttachmentById(ctx context.Context, instanceId, accountName, dbName string) (attachment map[string]string, has bool, errRet error) {
	attachments, err := me.DescribeAccountDBAttachments(ctx, instanceId, accountName, dbName)
	if err != nil {
		errRet = err
		return
	}

	if len(attachments) == 0 {
		return
	}

	attachment = attachments[0]
	has = true
	return
}

func (me *SqlserverService) GetInfoFromDeal(ctx context.Context, dealId string, timeout ...time.Duration) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeOrdersRequest()
	request.DealNames = []*string{&dealId}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var flowId int64
	var retryTimeout time.Duration
	if timeout != nil {
		retryTimeout = timeout[0]
	} else {
		retryTimeout = readRetryTimeout * 20
	}
	outErr := resource.Retry(retryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeOrders(request)
		if err != nil {
			return resource.RetryableError(err)
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK returns nil response, %s", request.GetAction())
			return resource.RetryableError(errRet)
		}
		if len(response.Response.Deals) == 0 {
			errRet = fmt.Errorf("TencentCloud SDK returns empty deal")
			return resource.RetryableError(errRet)
		} else if len(response.Response.Deals) > 1 {
			errRet = fmt.Errorf("TencentCloud SDK returns more than one deal")
			return resource.RetryableError(errRet)
		}
		if len(response.Response.Deals[0].InstanceIdSet) == 0 {
			err = fmt.Errorf("TencentCloud SDK returns empty InstanceIdSet")
			return resource.RetryableError(err)
		}
		instanceId = *response.Response.Deals[0].InstanceIdSet[0]
		flowId = *response.Response.Deals[0].FlowId
		if flowId == 0 {
			err = fmt.Errorf("TencentCloud SDK returns empty flowId")
			return resource.RetryableError(err)
		}
		return nil
	})
	if outErr != nil {
		return instanceId, outErr
	}
	outErr = me.WaitForTaskFinish(ctx, flowId)
	if outErr != nil {
		errRet = outErr
	}

	return
}

func (me *SqlserverService) WaitForTaskFinish(ctx context.Context, flowId int64) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeFlowStatusRequest()
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(6*writeRetryTimeout, func() *resource.RetryError {
		taskResponse, err := me.client.UseSqlserverClient().DescribeFlowStatus(request)
		ratelimit.Check(request.GetAction())
		if err != nil {
			return resource.NonRetryableError(errors.WithStack(err))
		}
		if *taskResponse.Response.Status == int64(SQLSERVER_TASK_RUNNING) {
			return resource.RetryableError(errors.WithStack(fmt.Errorf("SQLSERVER task status is %d(task running), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		} else if *taskResponse.Response.Status == int64(SQLSERVER_TASK_FAIL) {
			return resource.NonRetryableError(errors.WithStack(fmt.Errorf("SQLSERVER task status is %d(failed), requestId is %s", *taskResponse.Response.Status, *taskResponse.Response.RequestId)))
		}
		return nil
	})
	return
}

func (me *SqlserverService) CreateSqlserverDB(ctx context.Context, instanceID string, dbname string, charset string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewCreateDBRequest()

	// set instance id
	request.InstanceId = &instanceID
	// set DBs
	var dbCreateInfo = sqlserver.DBCreateInfo{
		DBName:  &dbname,
		Charset: &charset,
		Remark:  &remark,
	}
	var dbInfoList []*sqlserver.DBCreateInfo
	dbInfoList = append(dbInfoList, &dbCreateInfo)
	request.DBs = dbInfoList

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *sqlserver.CreateDBResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().CreateDB(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response != nil && response.Response != nil && *response.Response.FlowId != 0 {
		return me.WaitForTaskFinish(ctx, *response.Response.FlowId)
	}
	return
}

func (me *SqlserverService) DescribeDBsOfInstance(ctx context.Context, instanceId string) (instanceDBList []*sqlserver.DBDetail, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	if instanceId != "" {
		request.InstanceIdSet = []*string{&instanceId}
	}
	var offset, limit uint64 = SQLSERVER_DEFAULT_OFFSET, SQLSERVER_DEFAULT_LIMIT

	for {
		request.Offset = &offset
		request.Limit = &limit
		var response *sqlserver.DescribeDBsResponse
		err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseSqlserverClient().DescribeDBs(request)
			if e != nil {
				log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
				return retryError(e)
			}
			response = result
			return nil
		})
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK returns nil response for api[%s]", request.GetAction())
			return
		}
		if len(response.Response.DBInstances) == 0 {
			return
		} else if len(response.Response.DBInstances) > 1 {
			errRet = fmt.Errorf("[CRITAL]%s api[%s] returned multiple DB lists for one instance", logId, request.GetAction())
			return
		}
		instanceDBList = append(instanceDBList, response.Response.DBInstances[0].DBDetails...)
		if len(response.Response.DBInstances) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeDBDetailsById(ctx context.Context, dbId string) (dbInfo *sqlserver.DBDetail, has bool, errRet error) {
	idItem := strings.Split(dbId, FILED_SP)
	if len(idItem) < 2 {
		errRet = fmt.Errorf("broken ID of SQLServer DB %s", dbId)
		return
	}
	instanceId := idItem[0]
	dbName := idItem[1]

	instanceDBList, err := me.DescribeDBsOfInstance(ctx, instanceId)
	if err != nil {
		errRet = err
		return
	}
	if len(instanceDBList) == 0 {
		return
	}

	for _, dbDetail := range instanceDBList {
		if *dbDetail.Name == dbName {
			dbInfo = dbDetail
			if *dbDetail.Status != SQLSERVER_DB_DELETING {
				has = true
			}
			break
		}
	}
	return
}

func (me *SqlserverService) ModifySqlserverDBRemark(ctx context.Context, instanceId string, dbName string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyDBRemarkRequest()
	request.InstanceId = &instanceId
	request.DBRemarks = []*sqlserver.DBRemark{{Name: &dbName, Remark: &remark}}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseSqlserverClient().ModifyDBRemark(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s\n", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return
}

func (me *SqlserverService) DeleteSqlserverDB(ctx context.Context, instanceId string, names []*string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteDBRequest()
	request.InstanceId = &instanceId
	request.Names = names
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var response *sqlserver.DeleteDBResponse
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseSqlserverClient().DeleteDB(request)
		if e != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		response = result
		return nil
	})
	if err != nil {
		return err
	}

	if response != nil && response.Response != nil {
		return me.WaitForTaskFinish(ctx, *response.Response.FlowId)
	}
	return
}

func (me *SqlserverService) CreateSqlserverPublishSubscribe(ctx context.Context, publishInstanceId, subscribeInstanceId, publishSubscribeName string, databaseTuples []interface{}) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewCreatePublishSubscribeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.PublishInstanceId = &publishInstanceId
	request.SubscribeInstanceId = &subscribeInstanceId
	request.PublishSubscribeName = &publishSubscribeName
	for _, inst_ := range databaseTuples {
		inst := inst_.(map[string]interface{})
		request.DatabaseTupleSet = append(request.DatabaseTupleSet, sqlServerNewDatabaseTuple(inst["publish_database"], inst["publish_database"]))
	}

	var response *sqlserver.CreatePublishSubscribeResponse
	errRet = resource.Retry(2*writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, errRet = me.client.UseSqlserverClient().CreatePublishSubscribe(request)
		if errRet != nil {
			if ee, ok := errRet.(*SDKErrors.TencentCloudSDKError); ok {
				if ee.Code == INTERNALERROR_DBERROR || ee.Code == INSTANCE_STATUS_INVALID {
					return resource.RetryableError(errRet)
				} else {
					return resource.NonRetryableError(errRet)
				}
			}
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), errRet.Error())
			return resource.NonRetryableError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	flowId := int64(*response.Response.FlowId)
	errRet = me.WaitForTaskFinish(ctx, flowId)
	return
}

func sqlServerNewDatabaseTuple(publishDatabase, subscribeDatabase interface{}) *sqlserver.DatabaseTuple {
	publish := publishDatabase.(string)
	subscribe := subscribeDatabase.(string)
	databaseTuple := sqlserver.DatabaseTuple{
		PublishDatabase:   &publish,
		SubscribeDatabase: &subscribe,
	}
	return &databaseTuple
}

func (me *SqlserverService) DescribeSqlserverPublishSubscribeById(ctx context.Context, instanceId, pubOrSubInstanceId string) (instance *sqlserver.PublishSubscribe, has bool, errRet error) {
	paramMap := make(map[string]interface{})
	paramMap["instanceId"] = instanceId
	paramMap["pubOrSubInstanceId"] = pubOrSubInstanceId
	paramMap["publishSubscribeId"] = *helper.IntUint64(0)
	instanceList, err := me.DescribeSqlserverPublishSubscribes(ctx, paramMap)
	if err != nil {
		errRet = err
		return
	}

	if len(instanceList) == 0 {
		return
	}
	instance = instanceList[0]
	if instance != nil {
		has = true
	}
	return
}

func (me *SqlserverService) DescribeSqlserverPublishSubscribes(ctx context.Context, paramMap map[string]interface{}) (publishSubscribeList []*sqlserver.PublishSubscribe, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribePublishSubscribeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	if v, ok := paramMap["instanceId"]; ok {
		instanceId := v.(string)
		request.InstanceId = &instanceId
	}
	if v, ok := paramMap["pubOrSubInstanceId"]; ok {
		pubOrSubInstanceId := v.(string)
		request.PubOrSubInstanceId = &pubOrSubInstanceId
	}
	if v, ok := paramMap["pubOrSubInstanceIp"]; ok {
		pubOrSubInstanceIp := v.(string)
		request.PubOrSubInstanceIp = &pubOrSubInstanceIp
	}
	if v, ok := paramMap["publishSubscribeId"]; ok {
		publishSubscribeId := v.(uint64)
		request.PublishSubscribeId = &publishSubscribeId
	}
	if v, ok := paramMap["publishSubscribeName"]; ok {
		publishSubscribeName := v.(string)
		request.PublishSubscribeName = &publishSubscribeName
	}
	if v, ok := paramMap["publishDBName"]; ok {
		publishDBName := v.(string)
		request.PublishDBName = &publishDBName
		request.SubscribeDBName = &publishDBName
	}
	var offset, limit uint64 = SQLSERVER_DEFAULT_OFFSET, SQLSERVER_DEFAULT_LIMIT

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		var response *sqlserver.DescribePublishSubscribeResponse
		var err error
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			response, err = me.client.UseSqlserverClient().DescribePublishSubscribe(request)
			if err != nil {
				log.Printf("[CRITAL]%s DescribePublishSubscribe fail, reason:%s", logId, err.Error())
				return retryError(err)
			}
			return nil
		})

		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil || response.Response.PublishSubscribeSet == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		publishSubscribeList = append(publishSubscribeList, response.Response.PublishSubscribeSet...)
		if len(response.Response.PublishSubscribeSet) < int(limit) {
			break
		}
		offset += limit
	}
	return
}
func (me *SqlserverService) ModifyPublishSubscribeName(ctx context.Context, id uint64, name string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewModifyPublishSubscribeNameRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.PublishSubscribeId = &id
	request.PublishSubscribeName = &name
	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, e := me.client.UseSqlserverClient().ModifyPublishSubscribeName(request)
		if e != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason:%s", logId, request.GetAction(), e.Error())
			return retryError(e)
		}
		return nil
	})
	return errRet
}

func (me *SqlserverService) DeletePublishSubscribe(ctx context.Context, publishSubscribe *sqlserver.PublishSubscribe, deleteDatabaseTuples []*sqlserver.DatabaseTuple) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeletePublishSubscribeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.PublishSubscribeId = publishSubscribe.Id
	request.DatabaseTupleSet = deleteDatabaseTuples
	var response *sqlserver.DeletePublishSubscribeResponse
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		response, err = me.client.UseSqlserverClient().DeletePublishSubscribe(request)
		if err != nil {
			log.Printf("[CRITAL]%s %s fail, reason:%s", logId, request.GetAction(), err.Error())
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *SqlserverService) RecycleDBInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewRecycleDBInstanceRequest()
	request.InstanceId = &instanceId
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var err error
	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err = me.client.UseSqlserverClient().RecycleDBInstance(request)
		if err != nil {
			// FIXME: if action offline then kill this function
			code := err.(*SDKErrors.TencentCloudSDKError).Code
			if code == "InvalidAction" {
				return nil
			}
			return retryError(err)
		}
		return nil
	})
	if err != nil {
		return err
	}
	return
}

func (me *SqlserverService) CreateSqlserverBasicInstance(ctx context.Context, paramMap map[string]interface{}, weekSet []int, voucherIds, securityGroups []string) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewCreateBasicDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	var (
		cpu         = paramMap["cpu"].(int)
		memory      = paramMap["memory"].(int)
		storage     = paramMap["storage"].(int)
		subnetId    = paramMap["subnetId"].(string)
		vpcId       = paramMap["vpcId"].(string)
		machineType = paramMap["machineType"].(string)
		payType     = paramMap["payType"].(string)
		dbVersion   = paramMap["engineVersion"].(string)
		period      = paramMap["period"].(int)
		autoRenew   = paramMap["autoRenew"].(int)
		autoVoucher = paramMap["autoVoucher"].(int)
		zone        = paramMap["availabilityZone"].(string)
		collation   = paramMap["collation"].(string)
	)
	request.Cpu = helper.IntUint64(cpu)
	request.Memory = helper.IntUint64(memory)
	request.Storage = helper.IntUint64(storage)
	request.SubnetId = &subnetId
	request.VpcId = &vpcId
	request.MachineType = &machineType
	request.InstanceChargeType = &payType
	request.GoodsNum = helper.IntUint64(1)
	request.DBVersion = &dbVersion
	request.Period = helper.IntInt64(period)
	request.AutoRenewFlag = helper.IntInt64(autoRenew)

	request.AutoVoucher = helper.IntInt64(autoVoucher)
	request.Zone = &zone
	request.Collation = &collation
	if v, ok := paramMap["projectId"]; ok {
		projectId := v.(int)
		request.ProjectId = helper.IntUint64(projectId)
	}
	if v, ok := paramMap["startTime"]; ok {
		startTime := v.(string)
		request.StartTime = &startTime
	}
	if v, ok := paramMap["timeSpan"]; ok {
		timeSpan := v.(int)
		request.Span = helper.IntInt64(timeSpan)
	}

	if len(weekSet) > 0 {
		request.Weekly = make([]*int64, 0)
		for _, i := range weekSet {
			request.Weekly = append(request.Weekly, helper.IntInt64(i))
		}
	}
	request.VoucherIds = make([]*string, 0, len(voucherIds))
	for _, v := range voucherIds {
		request.VoucherIds = append(request.VoucherIds, &v)
	}

	request.SecurityGroupList = make([]*string, 0, len(securityGroups))
	for _, v := range securityGroups {
		request.SecurityGroupList = append(request.SecurityGroupList, &v)
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().CreateBasicDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	dealId := *response.Response.DealName
	instanceId, err = me.GetInfoFromDeal(ctx, dealId)
	if err != nil {
		errRet = err
	}
	return
}

func (me *SqlserverService) UpgradeSqlserverBasicInstance(ctx context.Context, instanceId string, memory int, storage, cpu, autoVoucher int, voucherIds []string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewUpgradeDBInstanceRequest()
	request.InstanceId = &instanceId
	request.Memory = helper.IntInt64(memory)
	request.Storage = helper.IntInt64(storage)
	request.Cpu = helper.IntInt64(cpu)
	request.AutoVoucher = helper.IntInt64(autoVoucher)
	if autoVoucher == 1 {
		request.VoucherIds = make([]*string, 0, len(voucherIds))
	}
	for _, v := range voucherIds {
		request.VoucherIds = append(request.VoucherIds, &v)
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().UpgradeDBInstance(request)
	if err != nil {
		return err
	}

	startPending := false

	errRet = resource.Retry(10*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribeSqlserverInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		}
		if !has {
			return resource.NonRetryableError(fmt.Errorf("cannot find SQL Server basic instance %s", instanceId))
		}
		// Status == 9, expanding
		if *instance.Status == 9 {
			startPending = true
			return resource.RetryableError(fmt.Errorf("expanding, SQL Server basic instance ID %s, status %d.... ", instanceId, *instance.Status))
		} else if !startPending && *instance.Status == 2 {
			return resource.RetryableError(fmt.Errorf("expanding, SQL Server basic Prepaid instance ID %s, status %d.... ", instanceId, *instance.Status))
		}
		return nil
	})

	return
}

func (me *SqlserverService) NewModifyDBInstanceRenewFlag(ctx context.Context, instanceId string, renewFlag int) (errRet error) {
	logId := getLogId(ctx)
	var instanceRenewInfo = make([]*sqlserver.InstanceRenewInfo, 1)
	instanceRenewInfo[0] = &sqlserver.InstanceRenewInfo{
		InstanceId: &instanceId,
		RenewFlag:  helper.IntInt64(renewFlag),
	}
	request := sqlserver.NewModifyDBInstanceRenewFlagRequest()
	request.RenewFlags = instanceRenewInfo

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseSqlserverClient().ModifyDBInstanceRenewFlag(request)

	return err
}

func (me *SqlserverService) DescribeSqlserverMigrationById(ctx context.Context, migrateId string) (migration *sqlserver.DescribeMigrationDetailResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeMigrationDetailRequest()
	request.MigrateId = helper.StrToUint64Point(migrateId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeMigrationDetail(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	migration = response.Response
	return
}

func (me *SqlserverService) DeleteSqlserverMigrationById(ctx context.Context, migrateId string) (errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDeleteMigrationRequest()
	request.MigrateId = helper.StrToUint64Point(migrateId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DeleteMigration(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverMigrationsByFilter(ctx context.Context, param map[string]interface{}) (migrateTasks []*sqlserver.MigrateTask, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeMigrationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "migrate_name" {
			request.MigrateName = v.(*string)
		}

		if k == "status_set" {
			request.StatusSet = v.([]*int64)
		}

		if k == "order_by" {
			request.OrderBy = v.(*string)
		}

		if k == "order_by_type" {
			request.OrderByType = v.(*string)
		}

	}
	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 20

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeMigrations(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MigrateTaskSet) < 1 {
			break
		}
		migrateTasks = append(migrateTasks, response.Response.MigrateTaskSet...)
		if len(response.Response.MigrateTaskSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}
	return
}

func (me *SqlserverService) DescribeSqlserverConfigBackupStrategyById(ctx context.Context, instanceId string) (configBackupStrategy *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DBInstances) < 1 {
		return
	}

	configBackupStrategy = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverBackupByBackupId(ctx context.Context, instanceId string, startTime string, endTime string, backupId uint64) (backupList []*sqlserver.Backup, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeBackupsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.StartTime = &startTime
	request.EndTime = &endTime
	request.BackupId = &backupId

	var offset, limit int64 = 0, 20

	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseSqlserverClient().DescribeBackups(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
			return
		}
		backupList = append(backupList, response.Response.Backups...)
		if len(response.Response.Backups) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *SqlserverService) DescribeBackupByFlowId(ctx context.Context, instanceId, flowId string) (BackupInfo *sqlserver.DescribeBackupByFlowIdResponse, errRet error) {

	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeBackupByFlowIdRequest()
	)

	request.InstanceId = &instanceId
	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeBackupByFlowId(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	BackupInfo = response
	return
}

func (me *SqlserverService) DescribeSqlserverBackupsById(ctx context.Context, instanceId, groupId string) (generalBackups *sqlserver.BackupFile, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeBackupFilesRequest()
	request.InstanceId = &instanceId
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeBackupFiles(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	generalBackups = response.Response.BackupFiles[0]
	return
}

func (me *SqlserverService) DeleteSqlserverGeneralBackupsById(ctx context.Context, instanceId, backupName string) (errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewRemoveBackupsRequest()
	request.InstanceId = &instanceId
	request.BackupNames = []*string{&backupName}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().RemoveBackups(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverBackupCommand(ctx context.Context, param map[string]interface{}) (datasourceBackupCommand []*sqlserver.DescribeBackupCommandResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeBackupCommandRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "BackupFileType" {
			request.BackupFileType = v.(*string)
		}
		if k == "DataBaseName" {
			request.DataBaseName = v.(*string)
		}
		if k == "IsRecovery" {
			request.IsRecovery = v.(*string)
		}
		if k == "LocalPath" {
			request.LocalPath = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeBackupCommand(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	datasourceBackupCommand = append(datasourceBackupCommand, response.Response)
	return
}

func (me *SqlserverService) DescribeCloneStatusByFlowId(ctx context.Context, flowId int64) (cloneStatus *sqlserver.DescribeFlowStatusResponseParams, errRet error) {

	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeFlowStatusRequest()
	)

	request.FlowId = &flowId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeFlowStatus(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	cloneStatus = response.Response
	return
}

func (me *SqlserverService) DescribeSqlserverGeneralCloneById(ctx context.Context, instanceId string) (generalCommunication []*sqlserver.DbNormalDetail, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBsNormalRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBsNormal(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	generalCommunication = response.Response.DBList
	return
}

func (me *SqlserverService) DeleteSqlserverGeneralCloneDB(ctx context.Context, instanceId, dbName string) (deleteResp *sqlserver.DeleteDBResponse, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteDBRequest()
	request.InstanceId = &instanceId
	request.Names = []*string{&dbName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DeleteDB(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	deleteResp = response
	return
}

func (me *SqlserverService) DescribeSqlserverFullBackupMigrationById(ctx context.Context, instanceId, backupMigrationId string) (fullBackupMigration *sqlserver.Migration, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeBackupMigrationRequest()
	request.InstanceId = &instanceId
	request.BackupMigrationId = &backupMigrationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeBackupMigration(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	fullBackupMigration = response.Response.BackupMigrationSet[0]
	return
}

func (me *SqlserverService) DeleteSqlserverFullBackupMigrationById(ctx context.Context, instanceId, backupMigrationId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteBackupMigrationRequest()
	request.InstanceId = &instanceId
	request.BackupMigrationId = &backupMigrationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DeleteBackupMigration(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverIncreBackupMigrationById(ctx context.Context, instanceId, backupMigrationId, incrementalMigrationId string) (increBackupMigration *sqlserver.Migration, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeIncrementalMigrationRequest()
	request.InstanceId = &instanceId
	request.BackupMigrationId = &backupMigrationId
	request.IncrementalMigrationId = &incrementalMigrationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeIncrementalMigration(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	increBackupMigration = response.Response.IncrementalMigrationSet[0]
	return
}

func (me *SqlserverService) DeleteSqlserverIncreBackupMigrationById(ctx context.Context, instanceId, backupMigrationId, incrementalMigrationId string) (errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDeleteIncrementalMigrationRequest()
	request.InstanceId = &instanceId
	request.BackupMigrationId = &backupMigrationId
	request.IncrementalMigrationId = &incrementalMigrationId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DeleteIncrementalMigration(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverBusinessIntelligenceFileById(ctx context.Context, instanceId, fileName string) (businessIntelligenceFile *sqlserver.BusinessIntelligenceFile, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeBusinessIntelligenceFileRequest()
	request.InstanceId = &instanceId
	request.FileName = &fileName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeBusinessIntelligenceFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	businessIntelligenceFile = response.Response.BackupMigrationSet[0]
	return
}

func (me *SqlserverService) DeleteSqlserverBusinessIntelligenceFileById(ctx context.Context, instanceId, fileName string) (errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDeleteBusinessIntelligenceFileRequest()
	request.InstanceId = &instanceId
	request.FileNameSet = []*string{&fileName}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DeleteBusinessIntelligenceFile(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverBusinessIntelligenceInstanceById(ctx context.Context, instanceId string) (businessIntelligenceInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	businessIntelligenceInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeMaintenanceSpanById(ctx context.Context, instanceId string) (maintenanceSpan *sqlserver.DescribeMaintenanceSpanResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeMaintenanceSpanRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeMaintenanceSpan(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	maintenanceSpan = response.Response

	return
}

func (me *SqlserverService) TerminateSqlserverInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewTerminateDBInstanceRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().TerminateDBInstance(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DeleteSqlserverInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDeleteDBInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DeleteDBInstance(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverGeneralCommunicationById(ctx context.Context, instanceId string) (generalCommunication *sqlserver.InterInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstanceInterRequest()
	request.InstanceId = &instanceId
	limit := int64(1)
	request.Limit = &limit

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstanceInter(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	generalCommunication = response.Response.InterInstanceSet[0]
	return
}

func (me *SqlserverService) DeleteSqlserverGeneralCommunicationById(ctx context.Context, instanceId string) (flowId int64, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewCloseInterCommunicationRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().CloseInterCommunication(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	flowId = *response.Response.InterInstanceFlowSet[0].FlowId

	return
}

func (me *SqlserverService) DescribeSqlserverBackupUploadSizeByFilter(ctx context.Context, param map[string]interface{}) (datasourceBackupUploadSize []*sqlserver.CosUploadBackupFile, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeBackupUploadSizeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}

		if k == "BackupMigrationId" {
			request.BackupMigrationId = v.(*string)
		}

		if k == "IncrementalMigrationId" {
			request.IncrementalMigrationId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeBackupUploadSize(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.CosUploadBackupFileSet) < 1 {
		return
	}

	datasourceBackupUploadSize = append(datasourceBackupUploadSize, response.Response.CosUploadBackupFileSet...)
	return
}

func (me *SqlserverService) DescribeSqlserverCrossRegionZoneByFilter(ctx context.Context, param map[string]interface{}) (datasourceCrossRegionZone *sqlserver.DescribeCrossRegionZoneResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeCrossRegionZoneRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeCrossRegionZone(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	datasourceCrossRegionZone = response.Response

	return
}

func (me *SqlserverService) DescribeSqlserverDatasourceDBCharsetsByFilter(ctx context.Context, param map[string]interface{}) (databaseCharsets []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeDBCharsetsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseSqlserverClient().DescribeDBCharsets(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	databaseCharsets = response.Response.DatabaseCharsets
	return
}

func (me *SqlserverService) DescribeSqlserverInstanceParamRecordsByFilter(ctx context.Context, param map[string]interface{}) (instanceParamRecords []*sqlserver.ParamRecord, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeInstanceParamRecordsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSqlserverClient().DescribeInstanceParamRecords(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}

		instanceParamRecords = append(instanceParamRecords, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SqlserverService) DescribeSqlserverProjectSecurityGroupsByFilter(ctx context.Context, param map[string]interface{}) (projectSecurityGroups []*sqlserver.SecurityGroup, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeProjectSecurityGroupsRequest()
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
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeProjectSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.SecurityGroupSet) < 1 {
		return
	}

	projectSecurityGroups = append(projectSecurityGroups, response.Response.SecurityGroupSet...)

	return
}

func (me *SqlserverService) DescribeSqlserverDatasourceRegionsByFilter(ctx context.Context) (datasourceRegions []*sqlserver.RegionInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeRegionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeRegions(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || *response.Response.TotalCount == 0 {
		return
	}

	datasourceRegions = response.Response.RegionSet

	return
}

func (me *SqlserverService) DescribeSqlserverRollbackTimeByFilter(ctx context.Context, param map[string]interface{}) (rollbackTime []*sqlserver.DbRollbackTimeInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeRollbackTimeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "DBs" {
			request.DBs = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeRollbackTime(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Details) < 1 {
		return
	}

	rollbackTime = response.Response.Details

	return
}

func (me *SqlserverService) DescribeSqlserverSlowlogsByFilter(ctx context.Context, param map[string]interface{}) (slowlogs []*sqlserver.SlowlogInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeSlowlogsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		limit  int64  = 100
		offset uint64 = 0
	)

	for {
		request.Limit = &limit
		request.Offset = &offset
		response, err := me.client.UseSqlserverClient().DescribeSlowlogs(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || *response.Response.TotalCount == 0 {
			break
		}

		slowlogs = append(slowlogs, response.Response.Slowlogs...)
		if len(response.Response.Slowlogs) < int(limit) {
			break
		}

		offset += uint64(limit)
	}

	return
}

func (me *SqlserverService) DescribeSqlserverUploadBackupInfoByFilter(ctx context.Context, param map[string]interface{}) (uploadBackupInfo *sqlserver.DescribeUploadBackupInfoResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeUploadBackupInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "BackupMigrationId" {
			request.BackupMigrationId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeUploadBackupInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	uploadBackupInfo = response.Response

	return
}

func (me *SqlserverService) DescribeSqlserverUploadIncrementalInfoByFilter(ctx context.Context, param map[string]interface{}) (uploadIncrementalInfo *sqlserver.DescribeUploadIncrementalInfoResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeUploadIncrementalInfoRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "BackupMigrationId" {
			request.BackupMigrationId = v.(*string)
		}
		if k == "IncrementalMigrationId" {
			request.IncrementalMigrationId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeUploadIncrementalInfo(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	uploadIncrementalInfo = response.Response

	return
}

func (me *SqlserverService) DescribeSqlserverConfigDatabaseCDCById(ctx context.Context, instanceId string) (configDatabaseCDC []*sqlserver.DbNormalDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsNormalRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBsNormal(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configDatabaseCDC = response.Response.DBList
	return
}

func (me *SqlserverService) DescribeSqlserverGeneralCloudInstanceById(ctx context.Context, instanceId string) (generalCloudInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	generalCloudInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverConfigDatabaseCTById(ctx context.Context, instanceId string) (configDatabaseCT []*sqlserver.DbNormalDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsNormalRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBsNormal(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configDatabaseCT = response.Response.DBList
	return
}

func (me *SqlserverService) DescribeSqlserverConfigDatabaseMdfById(ctx context.Context, instanceId string) (configDatabaseMdf []*sqlserver.DbNormalDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsNormalRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBsNormal(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configDatabaseMdf = response.Response.DBList
	return
}

func (me *SqlserverService) DescribeSqlserverConfigInstanceNetworkById(ctx context.Context, instanceId string) (configInstanceNetwork *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configInstanceNetwork = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverConfigInstanceParamById(ctx context.Context, instanceId string) (configInstanceParam []*sqlserver.ParameterDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeInstanceParamsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeInstanceParams(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configInstanceParam = response.Response.Items
	return
}

func (me *SqlserverService) DescribeSqlserverConfigInstanceRoGroupById(ctx context.Context, instanceId, readOnlyGroupId string) (configInstanceRoGroup *sqlserver.DescribeReadOnlyGroupDetailsResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeReadOnlyGroupDetailsRequest()
	request.InstanceId = &instanceId
	request.ReadOnlyGroupId = &readOnlyGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeReadOnlyGroupDetails(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	configInstanceRoGroup = response.Response
	return
}

func (me *SqlserverService) DescribeSqlserverConfigInstanceSecurityGroupsById(ctx context.Context, instanceId string) (configInstanceSecurityGroups []*sqlserver.SecurityGroup, errRet error) {
	logId := getLogId(ctx)
	request := sqlserver.NewDescribeDBSecurityGroupsRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBSecurityGroups(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SecurityGroupSet) < 1 {
		return
	}

	configInstanceSecurityGroups = response.Response.SecurityGroupSet
	return
}

func (me *SqlserverService) DescribeSqlserverRenewDBInstanceById(ctx context.Context, instanceId string) (renewDBInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	renewDBInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverRenewPostpaidDBInstanceById(ctx context.Context, instanceId string) (renewPostpaidDBInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	renewPostpaidDBInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverRestartDBInstanceById(ctx context.Context, instanceId string) (restartDBInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	restartDBInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverRestoreInstanceById(ctx context.Context, instanceId string, allNameList []string) (restoreInstance *sqlserver.InstanceDBDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 || response.Response.DBInstances == nil {
		return
	}

	restoreInstance = &sqlserver.InstanceDBDetail{}
	restoreInstance.InstanceId = response.Response.DBInstances[0].InstanceId
	tmpDbDetails := make([]*sqlserver.DBDetail, 0)
	for _, v := range allNameList {
		for _, DbDetail := range response.Response.DBInstances[0].DBDetails {
			if v == *DbDetail.Name {
				tmpDbDetails = append(tmpDbDetails, DbDetail)
				break
			}
		}
	}
	restoreInstance.DBDetails = tmpDbDetails

	return
}

func (me *SqlserverService) DescribeSqlserverRollbackInstanceById(ctx context.Context, instanceId string, allNameList []string) (rollBackInstance *sqlserver.InstanceDBDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	rollBackInstance = &sqlserver.InstanceDBDetail{}
	rollBackInstance.InstanceId = response.Response.DBInstances[0].InstanceId
	tmpDbDetails := make([]*sqlserver.DBDetail, 0)
	for _, v := range allNameList {
		for _, DbDetail := range response.Response.DBInstances[0].DBDetails {
			if v == *DbDetail.Name {
				tmpDbDetails = append(tmpDbDetails, DbDetail)
				break
			}
		}
	}
	rollBackInstance.DBDetails = tmpDbDetails

	return
}

func (me *SqlserverService) DescribeSqlserverConfigTerminateDBInstanceById(ctx context.Context, instanceId string) (configTerminateDBInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	configTerminateDBInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverDBS(ctx context.Context, instanceId, dbName string) (restoreInstance *sqlserver.InstanceDBDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsRequest()
	request.InstanceIdSet = []*string{&instanceId}
	request.Name = &dbName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBs(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	restoreInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverGeneralCloudRoInstanceById(ctx context.Context, instanceId string) (generalCloudRoInstance *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	generalCloudRoInstance = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DeleteSqlserverGeneralCloudRoInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDeleteDBInstanceRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DeleteDBInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *SqlserverService) DescribeSqlserverQueryXeventByFilter(ctx context.Context, param map[string]interface{}) (queryXevent []*sqlserver.Events, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeXEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
		if k == "EventType" {
			request.EventType = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseSqlserverClient().DescribeXEvents(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Events) < 1 {
			break
		}
		queryXevent = append(queryXevent, response.Response.Events...)
		if len(response.Response.Events) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *SqlserverService) DescribeSqlserverInsAttributeByFilter(ctx context.Context, param map[string]interface{}) (datasourceInsAttribute *sqlserver.DescribeDBInstancesAttributeResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeDBInstancesAttributeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	datasourceInsAttribute = response.Response

	return
}

func (me *SqlserverService) DescribeSqlserverInstanceTDEById(ctx context.Context, instanceId string) (instanceTDE *sqlserver.DescribeDBInstancesAttributeResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesAttributeRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceTDE = response.Response
	return
}

func (me *SqlserverService) DescribeSqlserverDatabaseTDEById(ctx context.Context, instanceId string, dbNameList []string) (databaseTDE *sqlserver.InstanceDBDetail, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBsRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBs(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.TotalCount == 0 {
		return
	}

	databaseTDE = &sqlserver.InstanceDBDetail{}
	databaseTDE.InstanceId = response.Response.DBInstances[0].InstanceId
	tmpDbDetails := make([]*sqlserver.DBDetail, 0)
	for _, v := range dbNameList {
		for _, DbDetail := range response.Response.DBInstances[0].DBDetails {
			if v == *DbDetail.Name {
				tmpDbDetails = append(tmpDbDetails, DbDetail)
				break
			}
		}
	}
	databaseTDE.DBDetails = tmpDbDetails

	return
}

func (me *SqlserverService) DescribeSqlserverInstanceHaById(ctx context.Context, instanceId string) (instanceHa *sqlserver.DBInstance, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesRequest()
	request.InstanceIdSet = []*string{&instanceId}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstances(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.DBInstances) < 1 {
		return
	}

	instanceHa = response.Response.DBInstances[0]
	return
}

func (me *SqlserverService) DescribeSqlserverInstanceSslById(ctx context.Context, instanceId string) (instanceSsl *sqlserver.DescribeDBInstancesAttributeResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := sqlserver.NewDescribeDBInstancesAttributeRequest()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeDBInstancesAttribute(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	instanceSsl = response.Response
	return
}

func (me *SqlserverService) DescribeSqlserverDescHaLogByFilter(ctx context.Context, param map[string]interface{}) (descHaLog []*sqlserver.SwitchLog, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = sqlserver.NewDescribeHASwitchLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}

		if k == "StartTime" {
			request.StartTime = v.(*string)
		}

		if k == "EndTime" {
			request.EndTime = v.(*string)
		}

		if k == "SwitchType" {
			request.SwitchType = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseSqlserverClient().DescribeHASwitchLog(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SwitchLog) < 1 {
		return
	}

	descHaLog = response.Response.SwitchLog
	return
}
