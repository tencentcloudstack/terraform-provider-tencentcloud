package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type DcdbService struct {
	client *connectivity.TencentCloudClient
}

//dc_account
func (me *DcdbService) DescribeDcdbAccount(ctx context.Context, instanceId, userName string) (account *dcdb.DescribeAccountsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId

	response, err := me.client.UseDcdbClient().DescribeAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	account = response.Response
	return
}

func (me *DcdbService) DeleteDcdbAccountById(ctx context.Context, instanceId, userName, host string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDeleteAccountRequest()

	request.InstanceId = &instanceId
	request.UserName = &userName
	request.Host = &host

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DeleteAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

//dc_db_instance
func (me *DcdbService) DescribeDcdbDbInstance(ctx context.Context, instanceId string) (dbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceIds = []*string{&instanceId}

	response, err := me.client.UseDcdbClient().DescribeDCDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	dbInstance = response.Response
	return
}

func (me *DcdbService) InitDcdbDbInstance(ctx context.Context, instanceId string, params []*dcdb.DBParamValue) (initRet bool, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDCDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(15*readRetryTimeout, func() *resource.RetryError {
		dbInstances, errResp := me.DescribeDcdbDbInstance(ctx, instanceId)
		if errResp != nil {
			return retryError(errResp, InternalError)
		}
		if dbInstances.Instances[0] == nil {
			return resource.NonRetryableError(fmt.Errorf("DescribeDcdbDbInstance return result(dcdb instance) is nil!"))
		}

		dbInstance := dbInstances.Instances[0]
		if *dbInstance.Status < 0 {
			return resource.NonRetryableError(fmt.Errorf("dcdb instance init status is %v, operate failed", *dbInstance.Status))
		}
		if *dbInstance.Status == 2 {
			return nil
		}
		if *dbInstance.Status == 3 {
			iniRequest := dcdb.NewInitDCDBInstancesRequest()
			iniRequest.InstanceIds = []*string{&instanceId}
			iniRequest.Params = params
			initErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := me.client.UseDcdbClient().InitDCDBInstances(iniRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
						logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
				}
				return nil
			})
			if initErr != nil {
				return resource.NonRetryableError(fmt.Errorf("dcdb instance init error %v, operate failed", initErr))
			}
			time.Sleep(10 * time.Second)
			return resource.RetryableError(fmt.Errorf("dcdb instance initializing, retry..."))
		}
		return resource.RetryableError(fmt.Errorf("dcdb instance init status is %v, retry...", *dbInstance.Status))
	})
	if err != nil {
		return false, err
	}
	return true, nil
}

//dc_hourdb_instance
func (me *DcdbService) DescribeDcdbHourdbInstance(ctx context.Context, instanceId string) (hourdbInstance *dcdb.DescribeDCDBInstancesResponseParams, errRet error) {
	return me.DescribeDcdbDbInstance(ctx, instanceId)
}

func (me *DcdbService) DeleteDcdbHourdbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyHourDCDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DestroyHourDCDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

//dc_sg
func (me *DcdbService) DescribeDcdbSecurityGroup(ctx context.Context, instanceId string) (securityGroup *dcdb.DescribeDBSecurityGroupsResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dcdb.NewDescribeDBSecurityGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.InstanceId = &instanceId
	request.Product = helper.String("dcdb") // api only use this fixed value

	response, err := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityGroup = response.Response

	return
}

func (me *DcdbService) DeleteDcdbSecurityGroupAttachmentById(ctx context.Context, instanceId, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDisassociateSecurityGroupsRequest()

	request.Product = helper.String("dcdb") // api only use this fixed value
	request.InstanceIds = []*string{&instanceId}
	request.SecurityGroupId = &securityGroupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DisassociateSecurityGroups(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
