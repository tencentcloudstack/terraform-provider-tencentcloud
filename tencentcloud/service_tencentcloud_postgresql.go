package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	postgresql "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/postgres/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type PostgresqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *PostgresqlService) CreatePostgresqlInstance(ctx context.Context, name, dbVersion, chargeType, specCode string, autoRenewFlag, projectId, period int, subnetId, vpcId, zone string, securityGroups []string, storage int, username, password, charset string) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewCreateInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Name = &name
	request.DBVersion = &dbVersion
	request.InstanceChargeType = &chargeType
	request.SpecCode = &specCode
	request.AutoRenewFlag = helper.IntInt64(autoRenewFlag)
	request.ProjectId = helper.Int64(int64(projectId))
	request.Period = helper.Int64Uint64(int64(period))
	request.SubnetId = &subnetId
	request.VpcId = &vpcId
	request.Storage = helper.IntUint64(storage)
	request.Zone = &zone
	request.InstanceCount = helper.Int64Uint64(1)
	request.AdminName = &username
	request.AdminPassword = &password
	request.Charset = &charset

	if len(securityGroups) > 0 {
		request.SecurityGroupIds = make([]*string, 0, len(securityGroups))
		for v := range securityGroups {
			request.SecurityGroupIds = append(request.SecurityGroupIds, &securityGroups[v])
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().CreateInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	if len(response.Response.DBInstanceIdSet) == 0 {
		errRet = errors.New("TencentCloud SDK returns empty postgresql ID")
		return
	} else if len(response.Response.DBInstanceIdSet) > 1 {
		errRet = errors.New("TencentCloud SDK returns more than one postgresql ID")
		return
	}
	instanceId = *response.Response.DBInstanceIdSet[0]
	return
}

func (me *PostgresqlService) InitPostgresqlInstance(ctx context.Context, instanceId string, username string, password string, charset string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewInitDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Charset = &charset
	request.AdminName = &username
	request.AdminPassword = &password
	request.DBInstanceIdSet = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().InitDBInstances(request)
	if err != nil {
		return err
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return
}

func (me *PostgresqlService) DescribeSpecinfos(ctx context.Context, zone string) (specCodeList []*postgresql.SpecItemInfo, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDescribeProductConfigRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Zone = &zone

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().DescribeProductConfig(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || len(response.Response.SpecInfoList) == 0 {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}

	specCodeList = response.Response.SpecInfoList[0].SpecItemInfoList
	return
}

func (me *PostgresqlService) ModifyPublicService(ctx context.Context, openInternet bool, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s-open:%t] fail,reason[%s]", logId, "modifyInternetService", openInternet, errRet.Error())
		}
	}()

	if openInternet {
		request := postgresql.NewOpenDBExtranetAccessRequest()
		request.DBInstanceId = &instanceId
		ratelimit.Check(request.GetAction())

		response, err := me.client.UsePostgresqlClient().OpenDBExtranetAccess(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}

		// check open or not
		err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
			instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
			if err != nil {
				return retryError(err)
			} else if has {
				if len(instance.DBInstanceNetInfo) > 0 {
					for _, v := range instance.DBInstanceNetInfo {
						if *v.NetType == "public" {
							if *v.Status == "opened" || *v.Status == "1" {
								return nil
							} else if *v.Status == "opening" || *v.Status == "initing" || *v.Status == "3" || *v.Status == "0" {
								return resource.RetryableError(fmt.Errorf("status %s, postgresql instance %s waiting", *v.Status, instanceId))
							} else {
								return resource.NonRetryableError(fmt.Errorf("status %s, postgresql instance %s open public service fail", *v.Status, instanceId))
							}
						}
					}
					// there is no public service yet
					return resource.RetryableError(fmt.Errorf("cannot find public status, postgresql instance %s watiting", instanceId))
				} else {
					return resource.NonRetryableError(fmt.Errorf("illegal net info of postgresql instance %s", instanceId))
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("check postgresql instance %s fail, instance is not exist", instanceId))
			}
		})
		if err != nil {
			return err
		}

	} else {
		request := postgresql.NewCloseDBExtranetAccessRequest()
		request.DBInstanceId = &instanceId
		ratelimit.Check(request.GetAction())

		response, err := me.client.UsePostgresqlClient().CloseDBExtranetAccess(request)
		if err != nil {
			return err
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		// check close or not
		err = resource.Retry(3*readRetryTimeout, func() *resource.RetryError {
			instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
			if err != nil {
				return retryError(err)
			} else if has {
				if len(instance.DBInstanceNetInfo) > 0 {
					for _, v := range instance.DBInstanceNetInfo {
						if *v.NetType == "public" {
							if *v.Status == "closed" || *v.Status == "2" {
								return nil
							} else if *v.Status == "closing" || *v.Status == "4" {
								return resource.RetryableError(fmt.Errorf("status %s, postgresql instance %s waiting", *v.Status, instanceId))
							} else {
								return resource.NonRetryableError(fmt.Errorf("status %s, postgresql instance %s open public service fail", *v.Status, instanceId))
							}
						}
					}
					// there is no public service
					return nil
				} else {
					return resource.NonRetryableError(fmt.Errorf("illegal net info of postgresql instance %s", instanceId))
				}
			} else {
				return resource.NonRetryableError(fmt.Errorf("check postgresql instance %s fail, instance is not exist", instanceId))
			}
		})
		if err != nil {
			return err
		}
	}
	return
}

func (me *PostgresqlService) DescribePostgresqlInstanceById(ctx context.Context, instanceId string) (instance *postgresql.DBInstance, has bool, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDescribeDBInstanceAttributeRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().DescribeDBInstanceAttribute(request)
	if err != nil {
		ee, ok := err.(*sdkErrors.TencentCloudSDKError)
		if !ok {
			errRet = err
			return
		}
		if ee.Code == "InvalidParameter" || ee.Code == "ResourceNotFound.InstanceNotFoundError" {
			errRet = nil
		} else {
			errRet = err
		}
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	instance = response.Response.DBInstance
	if instance != nil && *instance.DBInstanceStatus != "isolated" && *instance.DBInstanceStatus != "recycled" && *instance.DBInstanceStatus != "offline" {
		has = true
	}
	return
}

func (me *PostgresqlService) DescribePostgresqlInstances(ctx context.Context, filter []*postgresql.Filter) (instanceList []*postgresql.DBInstance, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDescribeDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit uint64 = 0, 10

	for {
		request.Offset = &offset
		request.Limit = &limit
		request.Filters = filter
		ratelimit.Check(request.GetAction())
		response, err := me.client.UsePostgresqlClient().DescribeDBInstances(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		instanceList = append(instanceList, response.Response.DBInstanceSet...)
		if len(response.Response.DBInstanceSet) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *PostgresqlService) ModifyPostgresqlInstanceName(ctx context.Context, instanceId string, name string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewModifyDBInstanceNameRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = &instanceId
	request.InstanceName = &name

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().ModifyDBInstanceName(request)
	return err
}

func (me *PostgresqlService) UpgradePostgresqlInstance(ctx context.Context, instanceId string, memory int, storage int) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewUpgradeDBInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = &instanceId
	request.Storage = helper.IntInt64(storage)
	request.Memory = helper.IntInt64(memory)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().UpgradeDBInstance(request)
	return err
}

func (me *PostgresqlService) ModifyPostgresqlInstanceProjectId(ctx context.Context, instanceId string, projectId int) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewModifyDBInstancesProjectRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceIdSet = []*string{&instanceId}
	request.ProjectId = helper.String(strconv.Itoa(projectId))

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().ModifyDBInstancesProject(request)
	return err
}

func (me *PostgresqlService) SetPostgresqlInstanceAutoRenewFlag(ctx context.Context, instanceId string, autoRenewFlag int) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewSetAutoRenewFlagRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceIdSet = []*string{&instanceId}
	request.AutoRenewFlag = helper.IntInt64(autoRenewFlag)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().SetAutoRenewFlag(request)
	return err
}

func (me *PostgresqlService) IsolatePostgresqlInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := postgresql.NewIsolateDBInstancesRequest()
	request.DBInstanceIdSet = []*string{&instanceId}
	ratelimit.Check(request.GetAction())
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	errRet = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		_, errRet = me.client.UsePostgresqlClient().IsolateDBInstances(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s describe account failed, reason: %v", logId, errRet)
			return retryError(errRet)
		}
		return nil
	})
	return errRet
}

func (me *PostgresqlService) DeletePostgresqlInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDestroyDBInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = &instanceId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().DestroyDBInstance(request)
	return err
}

func (me *PostgresqlService) SetPostgresqlInstanceRootPassword(ctx context.Context, instanceId string, password string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewResetAccountPasswordRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = &instanceId
	request.UserName = helper.String("root")
	request.Password = &password

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().ResetAccountPassword(request)
	return err
}

func (me *PostgresqlService) CheckDBInstanceStatus(ctx context.Context, instanceId string) error {
	// check status
	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		} else if has && *instance.DBInstanceStatus == POSTGRESQL_STAUTS_RUNNING {
			return nil
		} else if !has {
			return resource.NonRetryableError(fmt.Errorf("check postgresql instance %s fail", instanceId))
		} else {
			return resource.RetryableError(fmt.Errorf("checking postgresql instance %s , status %s ", instanceId, *instance.DBInstanceStatus))
		}
	})

	return err
}

func (me *PostgresqlService) DescribeRootUser(ctx context.Context, instanceId string) (accounts []*postgresql.AccountInfo, errRet error) {
	logId := getLogId(ctx)
	orderBy := "createTime"
	orderByType := "asc"

	request := postgresql.NewDescribeAccountsRequest()
	request.DBInstanceId = &instanceId
	request.OrderBy = &orderBy
	request.OrderByType = &orderByType
	var response *postgresql.DescribeAccountsResponse
	errRet = resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		response, errRet = me.client.UsePostgresqlClient().DescribeAccounts(request)
		if errRet != nil {
			log.Printf("[CRITAL]%s describe account failed, reason: %v", logId, errRet)
			return retryError(errRet)
		}
		return nil
	})
	if errRet != nil {
		return nil, errRet
	}
	if response == nil || response.Response == nil || response.Response.Details == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %+v, %s", response, request.GetAction())
	} else {
		accounts = response.Response.Details
	}

	return accounts, errRet
}

func (me *PostgresqlService) ModifyPgParams(ctx context.Context, instanceId string, paramEntrys map[string]string) (err error) {
	logId := getLogId(ctx)
	request := postgresql.NewModifyDBInstanceParametersRequest()
	request.DBInstanceId = &instanceId
	request.ParamList = make([]*postgresql.ParamEntry, 0)

	for key, value := range paramEntrys {
		request.ParamList = append(request.ParamList, &postgresql.ParamEntry{
			Name:          helper.String(key),
			ExpectedValue: helper.String(value),
		})
	}
	_, err = me.client.UsePostgresqlClient().ModifyDBInstanceParameters(request)

	if err != nil {
		log.Printf("[CRITAL]%s modify pgInstance parameter failed, reason: %v", logId, err)
		return err
	}
	return nil
}

func (me *PostgresqlService) DescribePgParams(ctx context.Context, instanceId string) (params map[string]string, err error) {
	logId := getLogId(ctx)
	stateCheckRequest := postgresql.NewDescribeParamsEventRequest()
	stateCheckRequest.DBInstanceId = &instanceId

	// wait for param effective
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		stateCheckResponse, inErr := me.client.UsePostgresqlClient().DescribeParamsEvent(stateCheckRequest)
		if inErr != nil {
			ee, ok := inErr.(*sdkErrors.TencentCloudSDKError)
			if ok && ee.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				return nil
			}
			return retryError(inErr)
		}

		for _, eventItem := range stateCheckResponse.Response.EventItems {
			eventDetail := eventItem.EventDetail
			if eventDetail == nil {
				continue
			}
			for _, eventDetailItem := range eventDetail {
				if *(eventDetailItem.State) != "success" {
					return resource.RetryableError(
						fmt.Errorf("params is updating!"))
				}
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	request := postgresql.NewDescribeDBInstanceParametersRequest()
	request.DBInstanceId = &instanceId

	response, err := me.client.UsePostgresqlClient().DescribeDBInstanceParameters(request)

	if err != nil {
		log.Printf("[CRITAL]%s fetch pg instance parameter failed, reason: %v", logId, err)
		return
	}
	detail := response.Response.Detail
	params = make(map[string]string)
	for _, item := range detail {
		params[*item.Name] = *item.CurrentValue
	}

	return
}
