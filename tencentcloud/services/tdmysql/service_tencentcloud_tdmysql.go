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

func (me *TdmysqlService) CreateTdmysqlDBSBackup(ctx context.Context, request *tdmysqlv20211122.CreateDBSBackupRequest) (ret *tdmysqlv20211122.CreateDBSBackupResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *tdmysqlv20211122.CreateDBSBackupResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().CreateDBSBackupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create tdmysql DBS backup failed, Response is nil."))
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

func (me *TdmysqlService) DescribeTdmysqlDBSBackupSets(ctx context.Context, request *tdmysqlv20211122.DescribeDBSBackupSetsRequest) (ret []*tdmysqlv20211122.BackupSetModel, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *tdmysqlv20211122.DescribeDBSBackupSetsResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DescribeDBSBackupSetsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe tdmysql DBS backup sets failed, Response is nil."))
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		return
	}

	if response.Response.Items != nil {
		ret = response.Response.Items
	}
	return
}

func (me *TdmysqlService) ModifyTdmysqlDBSBackupSetComment(ctx context.Context, request *tdmysqlv20211122.ModifyDBSBackupSetCommentRequest) (ret *tdmysqlv20211122.ModifyDBSBackupSetCommentResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *tdmysqlv20211122.ModifyDBSBackupSetCommentResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().ModifyDBSBackupSetCommentWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify tdmysql DBS backup set comment failed, Response is nil."))
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

func (me *TdmysqlService) DeleteTdmysqlDBSBackupSets(ctx context.Context, request *tdmysqlv20211122.DeleteDBSBackupSetsRequest) (ret *tdmysqlv20211122.DeleteDBSBackupSetsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var response *tdmysqlv20211122.DeleteDBSBackupSetsResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTdmysqlV20211122Client().DeleteDBSBackupSetsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete tdmysql DBS backup sets failed, Response is nil."))
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
