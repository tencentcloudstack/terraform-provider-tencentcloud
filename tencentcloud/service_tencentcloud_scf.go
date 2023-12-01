package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	scf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/scf/v20180416"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type scfFunctionInfo struct {
	name            string
	handler         *string
	desc            *string
	memSize         *int
	timeout         *int
	environment     map[string]string
	runtime         *string
	vpcId           *string
	subnetId        *string
	role            *string
	clsLogsetId     *string
	clsTopicId      *string
	funcType        *string
	namespace       *string
	layers          []*scf.LayerVersionSimple
	l5Enable        *bool
	publicNetConfig *scf.PublicNetConfigIn

	cosBucketName   *string
	cosObjectName   *string
	cosBucketRegion *string

	zipFile *string

	imageConfig *scf.ImageConfig

	cfsConfig *scf.CfsConfig

	tags map[string]string

	asyncRunEnable *string
	dnsCache       *string
	intranetConfig *scf.IntranetConfigIn
}

type scfTrigger struct {
	name        string
	triggerType string
	triggerDesc string
}

type ScfService struct {
	client *connectivity.TencentCloudClient
}

func (me *ScfService) CreateFunction(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewCreateFunctionRequest()
	request.FunctionName = &info.name
	request.Handler = info.handler
	request.Description = info.desc
	request.MemorySize = helper.IntInt64(*info.memSize)
	request.PublicNetConfig = info.publicNetConfig
	request.Timeout = helper.IntInt64(*info.timeout)
	for k, v := range info.environment {
		if request.Environment == nil {
			request.Environment = new(scf.Environment)
		}
		request.Environment.Variables = append(request.Environment.Variables, &scf.Variable{
			Key:   helper.String(k),
			Value: helper.String(v),
		})
	}
	request.Runtime = info.runtime

	if info.vpcId != nil {
		request.VpcConfig = &scf.VpcConfig{
			VpcId:    info.vpcId,
			SubnetId: info.subnetId,
		}
	}

	request.Namespace = info.namespace
	request.Role = info.role
	request.ClsLogsetId = info.clsLogsetId
	request.ClsTopicId = info.clsTopicId
	request.Layers = info.layers
	request.Type = info.funcType

	request.Code = &scf.Code{
		CosBucketName:   info.cosBucketName,
		CosObjectName:   info.cosObjectName,
		CosBucketRegion: info.cosBucketRegion,
		ZipFile:         info.zipFile,
		ImageConfig:     info.imageConfig,
	}

	if info.cfsConfig != nil && len(info.cfsConfig.CfsInsList) > 0 {
		request.CfsConfig = info.cfsConfig
	}

	if len(info.tags) > 0 {
		for k, v := range info.tags {
			key := k
			value := v
			request.Tags = append(request.Tags, &scf.Tag{
				Key:   &key,
				Value: &value,
			})
		}
	}

	request.AsyncRunEnable = info.asyncRunEnable
	request.DnsCache = info.dnsCache
	request.IntranetConfig = info.intranetConfig

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.CreateFunction(request); err != nil {
			e, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && strings.Contains(e.Code, "ResourceInUse") {
				return resource.NonRetryableError(err)
			}
			return retryError(errors.WithStack(err))
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (me *ScfService) DescribeFunction(ctx context.Context, name, namespace string) (resp *scf.GetFunctionResponse, err error) {
	client := me.client.UseScfClient()

	request := scf.NewGetFunctionRequest()
	request.FunctionName = &name
	request.Namespace = &namespace

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.GetFunction(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}

			return retryError(errors.WithStack(err), InternalError)
		}

		resp = response
		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func (me *ScfService) DescribeFunctions(ctx context.Context, name, namespace, desc *string, tags map[string]string) (functions []*scf.Function, err error) {
	client := me.client.UseScfClient()

	request := scf.NewListFunctionsRequest()
	request.SearchKey = name
	request.Namespace = namespace
	request.Description = desc
	for k, v := range tags {
		request.Filters = append(request.Filters, &scf.Filter{
			Name:   helper.String("tag-" + k),
			Values: []*string{helper.String(v)},
		})
	}
	request.Limit = helper.IntInt64(SCF_FUNCTION_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_FUNCTION_DESCRIBE_LIMIT

	request.Offset = &offset

	// at least run loop once
	for count == SCF_FUNCTION_DESCRIBE_LIMIT {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.ListFunctions(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			functions = append(functions, response.Response.Functions...)
			count = len(response.Response.Functions)

			return nil
		}); err != nil {
			return nil, err
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) ModifyFunctionCode(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewUpdateFunctionCodeRequest()
	request.FunctionName = &info.name
	request.Handler = info.handler
	request.Namespace = info.namespace
	request.Code = &scf.Code{
		CosBucketName:   info.cosBucketName,
		CosObjectName:   info.cosObjectName,
		CosBucketRegion: info.cosBucketRegion,
		ZipFile:         info.zipFile,
		ImageConfig:     info.imageConfig,
	}

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UpdateFunctionCode(request); err != nil {
			return retryError(errors.WithStack(err), InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
}

func (me *ScfService) ModifyFunctionConfig(ctx context.Context, info scfFunctionInfo) error {
	client := me.client.UseScfClient()

	request := scf.NewUpdateFunctionConfigurationRequest()
	request.FunctionName = &info.name
	request.Description = info.desc
	if info.memSize != nil {
		request.MemorySize = helper.IntInt64(*info.memSize)
	}
	if info.timeout != nil {
		request.Timeout = helper.IntInt64(*info.timeout)
	}
	request.Runtime = info.runtime

	request.Environment = new(scf.Environment)
	for k, v := range info.environment {
		request.Environment.Variables = append(request.Environment.Variables, &scf.Variable{
			Key:   helper.String(k),
			Value: helper.String(v),
		})
	}
	// clean all environments
	if len(request.Environment.Variables) == 0 {
		request.Environment.Variables = []*scf.Variable{
			{
				Key:   helper.String(""),
				Value: helper.String(""),
			},
		}
	}

	request.Namespace = info.namespace
	if info.vpcId != nil {
		request.VpcConfig = &scf.VpcConfig{VpcId: info.vpcId}
	}
	if info.subnetId != nil {
		if request.VpcConfig == nil {
			request.VpcConfig = new(scf.VpcConfig)
		}
		request.VpcConfig.SubnetId = info.subnetId
	}
	if info.publicNetConfig != nil {
		request.PublicNetConfig = info.publicNetConfig
	}
	request.Role = info.role
	request.ClsLogsetId = info.clsLogsetId
	request.ClsTopicId = info.clsTopicId
	if info.l5Enable != nil {
		request.L5Enable = helper.String("FALSE")
		if *info.l5Enable {
			request.L5Enable = helper.String("TRUE")
		}
	}

	request.DnsCache = info.dnsCache
	request.IntranetConfig = info.intranetConfig

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UpdateFunctionConfiguration(request); err != nil {
			return retryError(errors.WithStack(err), InternalError)
		}
		return nil
	}); err != nil {
		return err
	}

	return waitScfFunctionReady(ctx, info.name, *info.namespace, client)
}

func (me *ScfService) DeleteFunction(ctx context.Context, name, namespace string) error {
	client := me.client.UseScfClient()

	deleteRequest := scf.NewDeleteFunctionRequest()
	deleteRequest.FunctionName = &name
	deleteRequest.Namespace = &namespace

	if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(deleteRequest.GetAction())

		if _, err := client.DeleteFunction(deleteRequest); err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}
			return retryError(errors.WithStack(err), InternalError)
		}

		return nil
	}); err != nil {
		return err
	}

	descRequest := scf.NewGetFunctionRequest()
	descRequest.FunctionName = &name
	descRequest.Namespace = &namespace

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(descRequest.GetAction())

		if _, err := client.GetFunction(descRequest); err == nil {
			return resource.RetryableError(errors.New("function still exists"))
		} else {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}

			return retryError(errors.WithStack(err), InternalError, "ResourceNotFound.Version")
		}
	})
}

func (me *ScfService) CreateNamespace(ctx context.Context, namespace, desc string) error {
	client := me.client.UseScfClient()

	request := scf.NewCreateNamespaceRequest()
	request.Namespace = &namespace
	request.Description = &desc

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.CreateNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *ScfService) DescribeNamespace(ctx context.Context, namespace string) (ns *scf.Namespace, err error) {
	client := me.client.UseScfClient()

	request := scf.NewListNamespacesRequest()
	request.Limit = helper.IntInt64(SCF_NAMESPACE_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_NAMESPACE_DESCRIBE_LIMIT

	request.Offset = &offset

	// at least run loop once
	for count == SCF_NAMESPACE_DESCRIBE_LIMIT {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.ListNamespaces(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			for _, respNs := range response.Response.Namespaces {
				if *respNs.Name == namespace {
					ns = respNs
					return nil
				}
			}

			count = len(response.Response.Namespaces)
			return nil
		}); err != nil {
			return nil, err
		}

		if ns != nil {
			return
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) DescribeNamespaces(ctx context.Context) (nss []*scf.Namespace, err error) {
	client := me.client.UseScfClient()

	request := scf.NewListNamespacesRequest()
	request.Limit = helper.IntInt64(SCF_NAMESPACE_DESCRIBE_LIMIT)

	var offset int64
	count := SCF_NAMESPACE_DESCRIBE_LIMIT

	request.Offset = &offset

	// at least run loop once
	for count == SCF_NAMESPACE_DESCRIBE_LIMIT {
		if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := client.ListNamespaces(request)
			if err != nil {
				return retryError(errors.WithStack(err))
			}

			count = len(response.Response.Namespaces)
			nss = append(nss, response.Response.Namespaces...)

			return nil
		}); err != nil {
			return nil, err
		}

		offset += int64(count)
	}

	return
}

func (me *ScfService) ModifyNamespace(ctx context.Context, namespace, desc string) error {
	client := me.client.UseScfClient()

	request := scf.NewUpdateNamespaceRequest()
	request.Namespace = &namespace
	request.Description = &desc

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.UpdateNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}
		return nil
	})
}

func (me *ScfService) DeleteNamespace(ctx context.Context, namespace string) error {
	client := me.client.UseScfClient()

	request := scf.NewDeleteNamespaceRequest()
	request.Namespace = &namespace

	return resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := client.DeleteNamespace(request); err != nil {
			return retryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *ScfService) CreateTriggers(ctx context.Context, functionName, namespace string, triggers []scfTrigger) error {
	client := me.client.UseScfClient()

	for _, trigger := range triggers {
		request := scf.NewCreateTriggerRequest()
		request.FunctionName = &functionName
		request.TriggerName = &trigger.name
		request.Type = &trigger.triggerType
		request.TriggerDesc = &trigger.triggerDesc
		request.Namespace = &namespace
		request.Enable = helper.String("OPEN")

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			if _, err := client.CreateTrigger(request); err != nil {
				return retryError(errors.WithStack(err))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (me *ScfService) DeleteTriggers(ctx context.Context, functionName, namespace string, triggers []scfTrigger) error {
	client := me.client.UseScfClient()

	for _, trigger := range triggers {
		request := scf.NewDeleteTriggerRequest()
		request.FunctionName = &functionName
		request.Namespace = &namespace
		request.TriggerName = &trigger.name
		request.Type = &trigger.triggerType
		request.TriggerDesc = &trigger.triggerDesc

		if err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			if _, err := client.DeleteTrigger(request); err != nil {
				return retryError(errors.WithStack(err))
			}
			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}

func (me *ScfService) DescribeLogs(
	ctx context.Context,
	fnName, namespace, order, orderBy string,
	offset, limit int,
	retCode, invokeRequestId, startTime, endTime *string,
) (logs []*scf.FunctionLog, err error) {
	client := me.client.UseScfClient()

	request := scf.NewGetFunctionLogsRequest()
	request.FunctionName = &fnName
	request.Offset = helper.IntInt64(offset)
	request.Limit = helper.IntInt64(limit)
	request.Order = &order
	request.OrderBy = &orderBy
	if retCode != nil {
		request.Filter = &scf.LogFilter{RetCode: retCode}
	}
	request.Namespace = &namespace
	request.FunctionRequestId = invokeRequestId
	request.StartTime = startTime
	request.EndTime = endTime

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.GetFunctionLogs(request)
		if err != nil {
			if sdkError, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				for _, code := range SCF_FUNCTIONS_NOT_FOUND_SET {
					if sdkError.Code == code {
						return nil
					}
				}
			}
			return retryError(errors.WithStack(err))
		}

		logs = response.Response.Data
		return nil
	}); err != nil {
		return nil, err
	}

	return
}

func waitScfFunctionReady(ctx context.Context, name, namespace string, client *scf.Client) error {
	request := scf.NewGetFunctionRequest()
	request.FunctionName = &name
	request.Namespace = &namespace

	return resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		response, err := client.GetFunction(request)
		if err != nil {
			return retryError(errors.WithStack(err), InternalError)
		}

		switch *response.Response.Status {
		case SCF_FUNCTION_STATUS_CREATING, SCF_FUNCTION_STATUS_UPDATING:
			return resource.RetryableError(errors.New("function is not ready"))

		case SCF_FUNCTION_STATUS_ACTIVE:
			return nil

		default:
			return resource.NonRetryableError(errors.Errorf("function status is %s", *response.Response.Status))
		}
	})
}

func (me *ScfService) DescribeScfFunctionAliasById(ctx context.Context, namespace string, functionName string, name string) (functionAlias *scf.GetAliasResponse, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewGetAliasRequest()
	request.Namespace = &namespace
	request.FunctionName = &functionName
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetAlias(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	functionAlias = response
	return
}

func (me *ScfService) DeleteScfFunctionAliasById(ctx context.Context, namespace string, functionName string, name string) (errRet error) {
	logId := getLogId(ctx)

	request := scf.NewDeleteAliasRequest()
	request.Namespace = &namespace
	request.FunctionName = &functionName
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().DeleteAlias(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ScfService) DescribeScfFunctionVersionById(ctx context.Context, functionName string, namespace string, functionVersion string) (FunctionVersion *scf.GetFunctionResponse, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewGetFunctionRequest()
	request.FunctionName = &functionName
	request.Namespace = &namespace
	request.Qualifier = &functionVersion

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetFunction(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	FunctionVersion = response
	return
}

func (me *ScfService) DeleteScfFunctionVersionById(ctx context.Context, functionName string, namespace string, functionVersion string) (errRet error) {
	logId := getLogId(ctx)

	request := scf.NewDeleteFunctionRequest()
	request.FunctionName = &functionName
	request.Namespace = &namespace
	request.Qualifier = &functionVersion

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().DeleteFunction(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ScfService) DescribeScfFunctionEventInvokeConfigById(ctx context.Context, namespace string, functionName string) (FunctionEventInvokeConfig *scf.AsyncTriggerConfig, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewGetFunctionEventInvokeConfigRequest()
	request.Namespace = &namespace
	request.FunctionName = &functionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetFunctionEventInvokeConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	FunctionEventInvokeConfig = response.Response.AsyncTriggerConfig
	return
}

func (me *ScfService) DescribeScfReservedConcurrencyConfigById(ctx context.Context, namespace string, functionName string) (reservedConcurrencyConfig *scf.GetReservedConcurrencyConfigResponse, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewGetReservedConcurrencyConfigRequest()
	request.Namespace = &namespace
	request.FunctionName = &functionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetReservedConcurrencyConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	reservedConcurrencyConfig = response
	return
}

func (me *ScfService) DeleteScfReservedConcurrencyConfigById(ctx context.Context, namespace string, functionName string) (errRet error) {
	logId := getLogId(ctx)

	request := scf.NewDeleteReservedConcurrencyConfigRequest()
	request.Namespace = &namespace
	request.FunctionName = &functionName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().DeleteReservedConcurrencyConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ScfService) DescribeScfProvisionedConcurrencyConfigById(ctx context.Context, functionName string, qualifier string, namespace string) (provisionedConcurrencyConfig *scf.VersionProvisionedConcurrencyInfo, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewGetProvisionedConcurrencyConfigRequest()
	request.FunctionName = &functionName
	request.Qualifier = &qualifier
	request.Namespace = &namespace

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetProvisionedConcurrencyConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	info := response.Response.Allocated

	if len(info) < 1 {
		return
	}

	provisionedConcurrencyConfig = info[0]

	return
}

func (me *ScfService) DeleteScfProvisionedConcurrencyConfigById(ctx context.Context, functionName string, qualifier string, namespace string) (errRet error) {
	logId := getLogId(ctx)

	request := scf.NewDeleteProvisionedConcurrencyConfigRequest()
	request.FunctionName = &functionName
	request.Qualifier = &qualifier
	request.Namespace = &namespace

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().DeleteProvisionedConcurrencyConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *ScfService) DescribeScfAccountInfo(ctx context.Context) (accountInfo *scf.GetAccountResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewGetAccountRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetAccount(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	accountInfo = response.Response

	return
}

func (me *ScfService) DescribeScfAsyncEventManagementByFilter(ctx context.Context, param map[string]interface{}) (AsyncEventManagement []*scf.AsyncEvent, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListAsyncEventsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "Qualifier" {
			request.Qualifier = v.(*string)
		}
		if k == "InvokeType" {
			request.InvokeType = v.([]*string)
		}
		if k == "Status" {
			request.Status = v.([]*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "Orderby" {
			request.Orderby = v.(*string)
		}
		if k == "InvokeRequestId" {
			request.InvokeRequestId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseScfClient().ListAsyncEvents(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.EventList) < 1 {
			break
		}
		AsyncEventManagement = append(AsyncEventManagement, response.Response.EventList...)
		if len(response.Response.EventList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ScfService) DescribeScfTriggersByFilter(ctx context.Context, param map[string]interface{}) (Triggers []*scf.TriggerInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListTriggersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "Filters" {
			request.Filters = v.([]*scf.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseScfClient().ListTriggers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Triggers) < 1 {
			break
		}
		Triggers = append(Triggers, response.Response.Triggers...)
		if len(response.Response.Triggers) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ScfService) DescribeScfAsyncEventStatus(ctx context.Context, param map[string]interface{}) (asyncEventStatus *scf.AsyncEventStatus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewGetAsyncEventStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InvokeRequestId" {
			request.InvokeRequestId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetAsyncEventStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	asyncEventStatus = response.Response.Result

	return
}

func (me *ScfService) DescribeScfFunctionAddress(ctx context.Context, param map[string]interface{}) (functionAddress *scf.GetFunctionAddressResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewGetFunctionAddressRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "Qualifier" {
			request.Qualifier = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetFunctionAddress(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	functionAddress = response.Response

	return
}

func (me *ScfService) DescribeScfRequestStatusByFilter(ctx context.Context, param map[string]interface{}) (requestStatus []*scf.RequestStatus, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewGetRequestStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "FunctionRequestId" {
			request.FunctionRequestId = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "StartTime" {
			request.StartTime = v.(*string)
		}
		if k == "EndTime" {
			request.EndTime = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().GetRequestStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	requestStatus = response.Response.Data

	return
}

func (me *ScfService) DescribeScfFunctionAliasesByFilter(ctx context.Context, param map[string]interface{}) (FunctionAliases []*scf.Alias, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListAliasesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "FunctionVersion" {
			request.FunctionVersion = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = helper.Int64ToStrPoint(offset)
		request.Limit = helper.Int64ToStrPoint(limit)
		response, err := me.client.UseScfClient().ListAliases(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Aliases) < 1 {
			break
		}
		FunctionAliases = append(FunctionAliases, response.Response.Aliases...)
		if len(response.Response.Aliases) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ScfService) DescribeScfLayerVersions(ctx context.Context, param map[string]interface{}) (layerVersions []*scf.LayerVersionInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListLayerVersionsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "LayerName" {
			request.LayerName = v.(*string)
		}
		if k == "CompatibleRuntime" {
			request.CompatibleRuntime = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseScfClient().ListLayerVersions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	layerVersions = response.Response.LayerVersions

	return
}

func (me *ScfService) DescribeScfLayersByFilter(ctx context.Context, param map[string]interface{}) (layers []*scf.LayerVersionInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListLayersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "CompatibleRuntime" {
			request.CompatibleRuntime = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseScfClient().ListLayers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Layers) < 1 {
			break
		}
		layers = append(layers, response.Response.Layers...)
		if len(response.Response.Layers) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ScfService) DescribeScfFunctionVersionsByFilter(ctx context.Context, param map[string]interface{}) (functionVersions []*scf.FunctionVersion, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = scf.NewListVersionByFunctionRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "FunctionName" {
			request.FunctionName = v.(*string)
		}
		if k == "Namespace" {
			request.Namespace = v.(*string)
		}
		if k == "Order" {
			request.Order = v.(*string)
		}
		if k == "OrderBy" {
			request.OrderBy = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseScfClient().ListVersionByFunction(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.FunctionVersion) < 1 {
			break
		}
		functionVersions = append(functionVersions, response.Response.Versions...)
		if len(response.Response.FunctionVersion) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *ScfService) DescribeScfTriggerConfigById(ctx context.Context, functionName string, namespace string, triggerName string) (triggerConfig *scf.TriggerInfo, errRet error) {
	logId := getLogId(ctx)

	request := scf.NewListTriggersRequest()
	request.FunctionName = &functionName
	request.Namespace = &namespace

	filter := scf.Filter{
		Name:   helper.String("TriggerName"),
		Values: []*string{&triggerName},
	}

	request.Filters = append(request.Filters, &filter)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*scf.TriggerInfo, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseScfClient().ListTriggers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Triggers) < 1 {
			break
		}
		instances = append(instances, response.Response.Triggers...)
		if len(response.Response.Triggers) < int(limit) {
			break
		}

		offset += limit
	}

	if len(instances) < 1 {
		return
	}
	triggerConfig = instances[0]
	return
}
