package tencentcloud

import (
	"context"
	"fmt"
	"log"

	mongodb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20180408"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type MongodbService struct {
	client *connectivity.TencentCloudClient
}

func (me *MongodbService) DescribeInstanceById(ctx context.Context, instanceId string) (instance *mongodb.MongoDBInstanceDetail, errRet error) {
	logId := GetLogId(ctx)
	request := mongodb.NewDescribeDBInstancesRequest()
	request.InstanceIds = []*string{&instanceId}
	response, err := me.client.UseMongodbClient().DescribeDBInstances(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.InstanceDetails) < 1 {
		errRet = fmt.Errorf("mongodb id is not found")
		return
	}
	instance = response.Response.InstanceDetails[0]
	return
}

func (me *MongodbService) ModifyInstanceName(ctx context.Context, instanceId, instanceName string) error {
	logId := GetLogId(ctx)
	request := mongodb.NewRenameInstanceRequest()
	request.InstanceId = &instanceId
	request.NewName = &instanceName
	response, err := me.client.UseMongodbClient().RenameInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) SetInstancePassword(ctx context.Context, instanceId, accountName, password string) error {
	logId := GetLogId(ctx)
	request := mongodb.NewSetPasswordRequest()
	request.InstanceId = &instanceId
	request.UserName = &accountName
	request.Password = &password
	response, err := me.client.UseMongodbClient().SetPassword(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) UpgradeInstance(ctx context.Context, instanceId string, memory, volume int) error {
	logId := GetLogId(ctx)
	request := mongodb.NewUpgradeDBInstanceHourRequest()
	request.InstanceId = &instanceId
	request.Memory = intToPointer(memory)
	request.Volume = intToPointer(volume)
	response, err := me.client.UseMongodbClient().UpgradeDBInstanceHour(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) ModifyProjectId(ctx context.Context, instanceId string, projectId int) error {
	logId := GetLogId(ctx)
	request := mongodb.NewAssignProjectRequest()
	request.InstanceIds = []*string{&instanceId}
	request.ProjectId = intToPointer(projectId)
	response, err := me.client.UseMongodbClient().AssignProject(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) DeleteInstance(ctx context.Context, instanceId string) error {
	logId := GetLogId(ctx)
	request := mongodb.NewTerminateDBInstanceRequest()
	request.InstanceId = &instanceId
	response, err := me.client.UseMongodbClient().TerminateDBInstance(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *MongodbService) DescribeSpecInfo(ctx context.Context, zone string) (infos []*mongodb.SpecificationInfo, errRet error) {
	logId := GetLogId(ctx)
	request := mongodb.NewDescribeSpecInfoRequest()
	if zone != "" {
		request.Zone = &zone
	}
	response, err := me.client.UseMongodbClient().DescribeSpecInfo(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	infos = response.Response.SpecInfoList
	return
}

func (me *MongodbService) DescribeInstancesByFilter(ctx context.Context, instanceId string,
	instanceType int) (mongodbs []*mongodb.MongoDBInstanceDetail, errRet error) {

	logId := GetLogId(ctx)
	request := mongodb.NewDescribeDBInstancesRequest()
	if instanceId != "" {
		request.InstanceIds = []*string{&instanceId}
	}
	if instanceType > 0 {
		temp := int64(instanceType)
		request.InstanceType = &temp
	}

	offset := 0
	pageSize := 100
	mongodbs = make([]*mongodb.MongoDBInstanceDetail, 0)
	for {
		request.Offset = intToPointer(offset)
		request.Limit = intToPointer(pageSize)
		response, err := me.client.UseMongodbClient().DescribeDBInstances(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InstanceDetails) < 1 {
			break
		}

		mongodbs = append(mongodbs, response.Response.InstanceDetails...)

		if len(response.Response.InstanceDetails) < pageSize {
			break
		}
		offset += pageSize
	}
	return
}
