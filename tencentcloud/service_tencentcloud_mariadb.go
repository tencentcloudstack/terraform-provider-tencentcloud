package tencentcloud

import (
	"context"
	"log"

	mariadb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mariadb/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type MariadbService struct {
	client *connectivity.TencentCloudClient
}

func (me *MariadbService) DescribeMariadbDbInstance(ctx context.Context, instanceIds string) (dbInstance *mariadb.DescribeDBInstancesResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = mariadb.NewDescribeDBInstancesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.InstanceIds = []*string{&instanceIds}

	response, err := me.client.UseMariadbClient().DescribeDBInstances(request)
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

func (me *MariadbService) DeleteMariadbDbInstanceById(ctx context.Context, instanceIds string) (errRet error) {
	logId := getLogId(ctx)

	request := mariadb.NewDestroyDBInstanceRequest()

	request.InstanceId = &instanceIds

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMariadbClient().DestroyDBInstance(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
