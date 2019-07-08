package tencentcloud

import (
	"context"
	"fmt"
	dc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dc/v20180410"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"log"
)

type DcService struct {
	client *connectivity.TencentCloudClient
}

/////////common
func (me *DcService) fillFilter(ins []*dc.Filter, key, value string) (outs []*dc.Filter) {
	if ins == nil {
		ins = make([]*dc.Filter, 0, 2)
	}

	var filter = dc.Filter{Name: &key, Values: []*string{&value}}
	ins = append(ins, &filter)
	outs = ins
	return
}

func (me *DcService) strPt2str(pt *string) (ret string) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}
func (me *DcService) intPt2int(pt *int) (ret int) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}

func (me *DcService) int64Pt2int64(pt *int64) (ret int64) {
	if pt == nil {
		return
	} else {
		return *pt
	}
}

func (me *DcService) DescribeDirectConnects(ctx context.Context, dcId,
	name string) (infos []dc.DirectConnect, errRet error) {

	logId := GetLogId(ctx)
	request := dc.NewDescribeDirectConnectsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var limit int64 = 100
	var total int64 = -1
	var has = map[string]bool{}

	var filters []*dc.Filter
	if dcId != "" {
		filters = me.fillFilter(filters, "direct-connect-id", dcId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "direct-connect-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}
	infos = make([]dc.DirectConnect, 0, 10)

getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset

	response, err := me.client.UseDcClient().DescribeDirectConnects(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = *response.Response.TotalCount
	}

	if len(response.Response.DirectConnectSet) > 0 {
		offset += limit
	} else {
		//get empty set,we're done
		return
	}

	for _, item := range response.Response.DirectConnectSet {
		if has[*item.DirectConnectId] {
			errRet = fmt.Errorf("get repeated dc_id[%s] when doing DescribeDirectConnects", *item.DirectConnectId)
			return
		}
		has[*item.DirectConnectId] = true
		infos = append(infos, *item)
	}
	goto getMoreData
}

func (me *DcService) DescribeDirectConnectTunnels(ctx context.Context, dcxId,
	name string) (infos []dc.DirectConnectTunnel, errRet error) {

	logId := GetLogId(ctx)
	request := dc.NewDescribeDirectConnectTunnelsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var limit int64 = 100
	var total int64 = -1
	var has = map[string]bool{}

	var filters []*dc.Filter
	if dcxId != "" {
		filters = me.fillFilter(filters, "direct-connect-tunnel-id", dcxId)
	}
	if name != "" {
		filters = me.fillFilter(filters, "direct-connect-tunnel-name", name)
	}
	if len(filters) > 0 {
		request.Filters = filters
	}
	infos = make([]dc.DirectConnectTunnel, 0, 10)
getMoreData:
	if total >= 0 && offset >= total {
		return
	}
	request.Limit = &limit
	request.Offset = &offset

	response, err := me.client.UseDcClient().DescribeDirectConnectTunnels(request)
	if err != nil {
		errRet = err
		return
	}
	if total < 0 {
		total = *response.Response.TotalCount
	}

	if len(response.Response.DirectConnectTunnelSet) > 0 {
		offset += limit
	} else {
		//get empty set,we're done
		return
	}
	for _, item := range response.Response.DirectConnectTunnelSet {
		if has[*item.DirectConnectTunnelId] {
			errRet = fmt.Errorf("get repeated dcx_id[%s] when doing DescribeDirectConnectTunnels", *item.DirectConnectTunnelId)
			return
		}
		has[*item.DirectConnectTunnelId] = true
		infos = append(infos, *item)
	}
	goto getMoreData
}
