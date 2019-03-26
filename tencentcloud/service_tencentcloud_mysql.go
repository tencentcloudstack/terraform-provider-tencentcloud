package tencentcloud

import (
	"context"
	"log"

	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

type MysqlService struct {
	client *connectivity.TencentCloudClient
}

func (me *MysqlService) DescribeBackupsByInstanceId(ctx context.Context,
	instanceId string,
	leftNumber int64) (backupInfos []*cdb.BackupInfo, errRet error) {

	logId := GetLogId(ctx)

	listInitSize := leftNumber
	if listInitSize > 500 {
		listInitSize = 500
	}
	backupInfos = make([]*cdb.BackupInfo, 0, listInitSize)

	request := cdb.NewDescribeBackupsRequest()
	request.InstanceId = &instanceId

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
