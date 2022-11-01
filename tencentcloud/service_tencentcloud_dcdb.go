package tencentcloud

import (
	"context"
	"log"

	dcdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dcdb/v20180411"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
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
	request.InstanceIds[0] = &instanceId

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

func (me *DcdbService) DeleteDcdbDbInstanceById(ctx context.Context, instanceId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDestroyDCDBInstanceRequest()

	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDcdbClient().DestroyDCDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
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
func (me *DcdbService) DescribeDcdbSecurityGroupAttachment(ctx context.Context, instanceId, securityGroupId string) (securityGroupAttachment *dcdb.DescribeDBSecurityGroupsResponseParams, errRet error) {
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

	response, err := me.client.UseDcdbClient().DescribeDBSecurityGroups(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	securityGroupAttachment = response.Response

	return
}

func (me *DcdbService) DeleteDcdbSecurityGroupAttachmentById(ctx context.Context, instanceId, securityGroupId string) (errRet error) {
	logId := getLogId(ctx)

	request := dcdb.NewDisassociateSecurityGroupsRequest()

	request.InstanceIds[0] = &instanceId
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
