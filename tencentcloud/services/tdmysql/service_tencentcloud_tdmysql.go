package tdmysql

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTdmysqlService(client *connectivity.TencentCloudClient) TdmysqlService {
	return TdmysqlService{client: client}
}

type TdmysqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *TdmysqlService) CreateDBSBackup(ctx context.Context, request *tdmysqlv20211122.CreateDBSBackupRequest) (backupSetId int64, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.CreateDBSBackupResponse
	)

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().CreateDBSBackup(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("create tdmysql dbs backup failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create tdmysql dbs backup failed, reason:%+v", logId, err)
		return
	}

	if response.Response.BackupSetId == nil {
		err = fmt.Errorf("api %s returns empty backup set id", request.GetAction())
		return
	}

	backupSetId = *response.Response.BackupSetId
	return
}

func (me *TdmysqlService) DescribeDBSBackupSets(ctx context.Context, request *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (items []*tdmysqlv20211122.BackupSetModel, err error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		response *tdmysqlv20211122.DescribeDBSBackupSetsResponse
		offset   = int64(0)
		limit    = int64(200)
	)

	// The named result is initialised here (NOT in the var block above): a
	// named result and a body-level var of the same name do not compile.
	items = make([]*tdmysqlv20211122.BackupSetModel, 0)

	for {
		request.Offset = &offset
		request.Limit = &limit
		errRet := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseTdmysqlV20211122Client().DescribeDBSBackupSets(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("describe tdmysql dbs backup sets failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			log.Printf("[CRITAL]%s describe tdmysql dbs backup sets failed, reason:%+v", logId, errRet)
			return nil, errRet
		}

		if len(response.Response.Items) < 1 {
			break
		}

		items = append(items, response.Response.Items...)

		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TdmysqlService) ModifyDBSBackupSetComment(ctx context.Context, request *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().ModifyDBSBackupSetComment(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("modify tdmysql dbs backup set comment failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s modify tdmysql dbs backup set comment failed, reason:%+v", logId, err)
	}

	return
}

func (me *TdmysqlService) DeleteDBSBackupSets(ctx context.Context, request *tdmysqlv20211122.DeleteDBSBackupSetsRequest) (err error) {
	logId := tccommon.GetLogId(ctx)
	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DeleteDBSBackupSets(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("delete tdmysql dbs backup sets failed, Response is nil."))
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete tdmysql dbs backup sets failed, reason:%+v", logId, err)
	}

	return
}

func (me *TdmysqlService) DescribeTdmysqlDbInstanceById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeDBInstanceDetailResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeDBInstanceDetailRequest()
	response := tdmysqlv20211122.NewDescribeDBInstanceDetailResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBInstanceDetail(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBInstanceSSLStatusById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeInstanceSSLStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeInstanceSSLStatusRequest()
	response := tdmysqlv20211122.NewDescribeInstanceSSLStatusResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeInstanceSSLStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance ssl status failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBInstanceDataReservedSpaceById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeInstanceDataReservedSpaceResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeInstanceDataReservedSpaceRequest()
	response := tdmysqlv20211122.NewDescribeInstanceDataReservedSpaceResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeInstanceDataReservedSpace(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db instance data reserved space failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *TdmysqlService) DescribeTdmysqlDBSecurityGroupsById(ctx context.Context, instanceId string) (ret *tdmysqlv20211122.DescribeDBSecurityGroupsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tdmysqlv20211122.NewDescribeDBSecurityGroupsRequest()
	response := tdmysqlv20211122.NewDescribeDBSecurityGroupsResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBSecurityGroups(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql db security groups failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}
