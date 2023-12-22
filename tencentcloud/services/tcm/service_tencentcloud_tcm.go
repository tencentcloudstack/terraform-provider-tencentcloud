package tcm

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	tcm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcm/v20210413"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewTcmService(client *connectivity.TencentCloudClient) TcmService {
	return TcmService{client: client}
}

type TcmService struct {
	client *connectivity.TencentCloudClient
}

func (me *TcmService) DescribeTcmMesh(ctx context.Context, meshId string) (mesh *tcm.DescribeMeshResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tcm.NewDescribeMeshRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.MeshId = &meshId

	response, err := me.client.UseTcmClient().DescribeMesh(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	mesh = response.Response
	return
}

func (me *TcmService) DeleteTcmMeshById(ctx context.Context, meshId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcm.NewDeleteMeshRequest()

	request.MeshId = &meshId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcmClient().DeleteMesh(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TcmService) DeleteTcmClusterAttachmentById(ctx context.Context, meshId, clusterId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcm.NewUnlinkClusterRequest()

	request.MeshId = &meshId
	request.ClusterId = &clusterId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcmClient().UnlinkCluster(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TcmService) DeleteTcmPrometheusAttachmentById(ctx context.Context, meshID string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tcm.NewUnlinkPrometheusRequest()

	request.MeshID = &meshID

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTcmClient().UnlinkPrometheus(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TcmService) DescribeTcmAccessLogConfig(ctx context.Context, meshName string) (accessLogConfig *tcm.DescribeAccessLogConfigResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tcm.NewDescribeAccessLogConfigRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.MeshId = &meshName

	response, err := me.client.UseTcmClient().DescribeAccessLogConfig(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	accessLogConfig = response.Response
	return
}

func (me *TcmService) DescribeTcmMeshByFilter(ctx context.Context, param map[string][]*string) (meshs []*tcm.Mesh, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tcm.NewDescribeMeshListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		filter := &tcm.Filter{
			Name:   helper.String(k),
			Values: []*string{},
		}
		filter.Values = append(filter.Values, v...)
		request.Filters = append(request.Filters, filter)
	}

	ratelimit.Check(request.GetAction())

	var offset int64 = 0
	var pageSize int64 = 50

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTcmClient().DescribeMeshList(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.MeshList) < 1 {
			break
		}
		meshs = append(meshs, response.Response.MeshList...)
		if *response.Response.Total < pageSize {
			break
		}
		offset += pageSize
	}
	return
}
