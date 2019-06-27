package tencentcloud

import (
	"context"
	"fmt"
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

//Cnn basic information
type CnnBasicInfo struct {
	cnnId         string
	name          string
	description   string
	state         string
	qos           string
	instanceCount int64
	createTime    string
}

func (me *VpcService) DescribeCcn(ctx context.Context, cnnId string) (info CnnBasicInfo, has int, errRet error) {
	infos, err := me.DescribeCcns(ctx, cnnId, "")
	if err != nil {
		errRet = err
		return
	}
	has = len(infos)
	if has > 0 {
		info = infos[0]
	}
	return
}

func (me *VpcService) DescribeCcns(ctx context.Context, cnnId, name string) (infos []CnnBasicInfo, errRet error) {

	logId := GetLogId(ctx)
	request := vpc.NewDescribeCcnsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	infos = make([]CnnBasicInfo, 0, 100)

	var offset uint64 = 0
	var limit uint64 = 100
	var total = -1
	var has = map[string]bool{}

	var filters []*vpc.Filter
	if cnnId != "" {
		filters = me.fillFilter(filters, "ccn-id", cnnId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "ccn-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}

getMoreData:

	if total >= 0 {
		if int(offset) >= total {
			return
		}
	}
	request.Limit = &limit
	request.Offset = &offset

	response, err := me.client.UseVpcClient().DescribeCcns(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] , request body [%s], response body[%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if total < 0 {
		total = int(*response.Response.TotalCount)
	}

	if len(response.Response.CcnSet) > 0 {
		offset += limit
	} else {
		//get empty ,we're done
		return
	}
	for _, item := range response.Response.CcnSet {
		var basicInfo CnnBasicInfo

		basicInfo.cnnId = *item.CcnId
		basicInfo.name = *item.CcnName
		basicInfo.createTime = *item.CreateTime

		basicInfo.description = *item.CcnDescription
		basicInfo.instanceCount = int64(*item.InstanceCount)
		basicInfo.qos = *item.QosLevel
		basicInfo.state = *item.State

		if has[basicInfo.cnnId] {
			errRet = fmt.Errorf("get repeated cnn_id[%s] when doing DescribeCcns", basicInfo.cnnId)
			return
		}
		has[basicInfo.cnnId] = true
		infos = append(infos, basicInfo)
	}
	goto getMoreData

}
