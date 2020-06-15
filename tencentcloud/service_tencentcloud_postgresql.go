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
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type PostgresqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *PostgresqlService) CreatePostgresqlInstance(ctx context.Context, name string, dbVersion string, chargeType string, specCode string, autoRenewFlag int, projectId int, period int, subnetId string, vpcId string, zone string, storage int) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewCreateDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Name = helper.String(name)
	request.DBVersion = helper.String(dbVersion)
	request.InstanceChargeType = helper.String(chargeType)
	request.SpecCode = helper.String(specCode)
	request.AutoRenewFlag = helper.IntInt64(autoRenewFlag)
	request.ProjectId = helper.Int64(int64(projectId))
	request.Period = helper.Int64Uint64(int64(period))
	request.SubnetId = helper.String(subnetId)
	request.VpcId = helper.String(vpcId)
	request.Storage = helper.IntUint64(storage)
	request.Zone = helper.String(zone)
	request.InstanceCount = helper.Int64Uint64(1)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().CreateDBInstances(request)
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

func (me *PostgresqlService) InitPostgresqlInstance(ctx context.Context, instanceId string, password string, charset string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewInitDBInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Charset = helper.String(charset)
	request.AdminName = helper.String("root")
	request.AdminPassword = helper.String(password)
	request.DBInstanceIdSet = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().InitDBInstances(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	if len(response.Response.DBInstanceIdSet) == 0 {
		errRet = errors.New("TencentCloud SDK returns empty postgresql instance ID")
	} else if len(response.Response.DBInstanceIdSet) > 1 {
		errRet = errors.New("TencentCloud SDK returns more than one postgresql instance ID")
	}

	return
}

func (me *PostgresqlService) DescribeSpecCodes(ctx context.Context, zone string) (specCodeList []*postgresql.SpecItemInfo, errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDescribeProductConfigRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Zone = helper.String(zone)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().DescribeProductConfig(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
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
		request.DBInstanceId = helper.String(instanceId)
		ratelimit.Check(request.GetAction())

		response, err := me.client.UsePostgresqlClient().OpenDBExtranetAccess(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}

		/*
			time.Sleep(10 * time.Second)

			//get instance internet status
			instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
			if err != nil {
				return err
			}
			if !has {
				errRet = fmt.Errorf("describe postgresql instance error %s", instanceId)
				return
			}
			for _, v := range instance.DBInstanceNetInfo {
				if *v.NetType == "public" {
					if *v.Status == "opening" { //const

						//sleep and retry
					} else if *v.Status == "opened" {
						return nil
					} else {
						err = fmt.Errorf("open public accress failed ,instance Id %s,  status %s", instanceId, *v.Status)
						return err
					}
				} else {
					continue
				}
			}
			//get no public access
			err = fmt.Errorf("There is no public access with posgresql instance %s", instanceId)
			return err
		*/

	} else {
		request := postgresql.NewCloseDBExtranetAccessRequest()
		request.DBInstanceId = helper.String(instanceId)
		ratelimit.Check(request.GetAction())

		response, err := me.client.UsePostgresqlClient().CloseDBExtranetAccess(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}

		/*
			time.Sleep(10 * time.Second)
			//get instance internet status
			instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
			if err != nil {
				return err
			}
			if !has {
				errRet = fmt.Errorf("describe postgresql instance error %s", instanceId)
				return
			}
			for _, v := range instance.DBInstanceNetInfo {
				if *v.NetType == "public" {
					if *v.Status == "closing" { //const
						//sleep and retry
					} else if *v.Status == "closed" {
						return nil
					} else {
						err = fmt.Errorf("close public accress failed ,instance Id %s,  status %s", instanceId, *v.Status)
						return err
					}
				} else {
					continue
				}
			}
			//get no public access
			err = fmt.Errorf("There is no public access with posgresql instance %s", instanceId)
			return err

		*/
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
	request.DBInstanceId = helper.String(instanceId)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UsePostgresqlClient().DescribeDBInstanceAttribute(request)
	if err != nil {
		ee, ok := err.(*sdkErrors.TencentCloudSDKError)
		if !ok {
			errRet = err
			return
		}
		if ee.Code == "InvalidParameter" {
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
	if instance != nil {
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

	var offset, limit uint64 = 0, 20

	request.Offset = &offset
	request.Limit = &limit
	request.Filters = filter

	for {
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
	request.DBInstanceId = helper.String(instanceId)
	request.InstanceName = helper.String(name)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().ModifyDBInstanceName(request)
	if err != nil {
		errRet = err
	}

	return
}

func (me *PostgresqlService) UpgradePostgresqlInstance(ctx context.Context, instanceId string, memory int, storage int) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewUpgradeDBInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = helper.String(instanceId)
	request.Storage = helper.IntInt64(storage)
	request.Memory = helper.IntInt64(memory)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().UpgradeDBInstance(request)
	if err != nil {
		errRet = err
	}

	return
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
	if err != nil {
		errRet = err
	}
	return
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
	if err != nil {
		errRet = err
	}
	return
}

func (me *PostgresqlService) DeletePostgresqlInstance(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewDestroyDBInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = helper.String(instanceId)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().DestroyDBInstance(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *PostgresqlService) SetPostgresqlInstanceRootPassword(ctx context.Context, instanceId string, password string) (errRet error) {
	logId := getLogId(ctx)
	request := postgresql.NewResetAccountPasswordRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.DBInstanceId = helper.String(instanceId)
	request.UserName = helper.String("root")
	request.Password = helper.String(password)

	ratelimit.Check(request.GetAction())
	_, err := me.client.UsePostgresqlClient().ResetAccountPassword(request)
	if err != nil {
		errRet = err
	}
	return
}

func (me *PostgresqlService) CheckDBInstanceStatus(ctx context.Context, instanceId string) error {
	//check status
	err := resource.Retry(2*readRetryTimeout, func() *resource.RetryError {
		instance, has, err := me.DescribePostgresqlInstanceById(ctx, instanceId)
		if err != nil {
			return retryError(err)
		} else if has && *instance.DBInstanceStatus == POSTGRESQL_STAUTS_RUNNING {
			return nil
		} else if !has {
			return resource.NonRetryableError(fmt.Errorf("check postgresql instance fail"))
		} else {
			return resource.RetryableError(fmt.Errorf("checking postgresql instance %s , status %s ", instanceId, *instance.DBInstanceStatus))
		}
	})

	return err
}
