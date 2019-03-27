package tencentcloud

import (
	"context"
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
