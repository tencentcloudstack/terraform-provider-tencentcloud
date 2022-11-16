package tencentcloud

import (
	"context"
	"log"

	tat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tat/v20201028"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TatService struct {
	client *connectivity.TencentCloudClient
}

func (me *TatService) DescribeTatCommand(ctx context.Context, commandId string) (command *tat.Command, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tat.NewDescribeCommandsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&tat.Filter{
			Name:   helper.String("command-id"),
			Values: []*string{&commandId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	commandInfo := make([]*tat.Command, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeCommands(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.CommandSet) < 1 {
			break
		}
		commandInfo = append(commandInfo, response.Response.CommandSet...)
		if len(response.Response.CommandSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(commandInfo) < 1 {
		return
	}
	command = commandInfo[0]

	return

}

func (me *TatService) DeleteTatCommandById(ctx context.Context, commandId string) (errRet error) {
	logId := getLogId(ctx)

	request := tat.NewDeleteCommandRequest()

	request.CommandId = &commandId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTatClient().DeleteCommand(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TatService) DescribeTatInvoker(ctx context.Context, invokerId string) (invoker *tat.Invoker, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tat.NewDescribeInvokersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	request.Filters = append(
		request.Filters,
		&tat.Filter{
			Name:   helper.String("invoker-id"),
			Values: []*string{&invokerId},
		},
	)
	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 100
	invokerSet := make([]*tat.Invoker, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTatClient().DescribeInvokers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.InvokerSet) < 1 {
			break
		}
		invokerSet = append(invokerSet, response.Response.InvokerSet...)
		if len(response.Response.InvokerSet) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(invokerSet) < 1 {
		return
	}
	invoker = invokerSet[0]

	return

}

func (me *TatService) DeleteTatInvokerById(ctx context.Context, invokerId string) (errRet error) {
	logId := getLogId(ctx)

	request := tat.NewDeleteInvokerRequest()

	request.InvokerId = &invokerId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTatClient().DeleteInvoker(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
