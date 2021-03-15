package tencentcloud

import (
	"context"
	"fmt"
	"log"

	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TCRService struct {
	client *connectivity.TencentCloudClient
}

func (me *TCRService) CreateTCRInstance(ctx context.Context, name string, instanceType string, tags map[string]string) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryName = &name
	request.RegistryType = &instanceType

	if len(tags) > 0 {
		tagSpec := tcr.TagSpecification{ResourceType: helper.String("instance"), Tags: make([]*tcr.Tag, 0)}
		for k, v := range tags {
			key, value := k, v
			tag := tcr.Tag{Value: &value, Key: &key}
			tagSpec.Tags = append(tagSpec.Tags, &tag)
		}
		request.TagSpecification = &tagSpec
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateInstance(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.RegistryId == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	instanceId = *response.Response.RegistryId
	return
}

func (me *TCRService) DescribeTCRInstanceById(ctx context.Context, instanceId string) (instance *tcr.Registry, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Registryids = []*string{&instanceId}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeInstances(request)
	if err != nil {
		ee, ok := err.(*sdkErrors.TencentCloudSDKError)
		if !ok {
			errRet = err
			return
		}
		if ee.Code == "ResourceNotFound" {
			errRet = nil
		} else {
			errRet = err
		}
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	if len(response.Response.Registries) == 0 {
		return
	}
	if len(response.Response.Registries) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than one instances, instanceId %s, %s", instanceId, request.GetAction())
	}
	has = true
	instance = response.Response.Registries[0]
	return
}

func (me *TCRService) DescribeTCRInstances(ctx context.Context, instanceId string, filter []*tcr.Filter) (instanceList []*tcr.Registry, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeInstancesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()

	var offset, limit int64 = 0, 200
	request.Filters = filter
	if instanceId != "" {
		request.Registryids = []*string{&instanceId}
	}
	for {
		request.Offset = &offset
		request.Limit = &limit

		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTCRClient().DescribeInstances(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		instanceList = append(instanceList, response.Response.Registries...)
		if len(response.Response.Registries) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TCRService) DeleteTCRInstance(ctx context.Context, instanceId string, deleteBucket bool) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDeleteInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.DeleteBucket = &deleteBucket

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().DeleteInstance(request)
	return err
}

//long term token

//name space
func (me *TCRService) CreateTCRNameSpace(ctx context.Context, instanceId string, name string, isPublic bool) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateNamespaceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.IsPublic = &isPublic
	request.NamespaceName = &name

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateNamespace(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return
}

func (me *TCRService) ModifyTCRNameSpace(ctx context.Context, instanceId string, name string, isPublic bool) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewModifyNamespaceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.NamespaceName = &name
	request.IsPublic = &isPublic

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().ModifyNamespace(request)
	return err
}

func (me *TCRService) DeleteTCRNameSpace(ctx context.Context, instanceId string, name string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDeleteNamespaceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.NamespaceName = &name

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().DeleteNamespace(request)
	return err
}

func (me *TCRService) DescribeTCRNameSpaces(ctx context.Context, instanceId string, name string) (namespaceList []*tcr.TcrNamespaceInfo, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeNamespacesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	if name != "" {
		request.NamespaceName = &name
	}

	var offset, limit int64 = 0, 200
	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTCRClient().DescribeNamespaces(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				errRet = err
				return
			}
			if ee.Code == "ResourceNotFound" {
				errRet = nil
			} else {
				errRet = err
			}
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		namespaceList = append(namespaceList, response.Response.NamespaceList...)
		if len(response.Response.NamespaceList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TCRService) DescribeTCRNameSpaceById(ctx context.Context, instanceId string, name string) (namespace *tcr.TcrNamespaceInfo, has bool, errRet error) {
	namespaces, err := me.DescribeTCRNameSpaces(ctx, instanceId, name)
	if err != nil {
		return nil, false, err
	}

	if len(namespaces) == 0 {
		return nil, has, nil
	} else if len(namespaces) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than 1 namespaces, %s %s", instanceId, name)
		return
	}

	namespace = namespaces[0]
	has = true
	return
}

//repository
func (me *TCRService) CreateTCRRepository(ctx context.Context, instanceId string, namespace string, repositoryName string, briefDesc string, description string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateRepositoryRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.NamespaceName = &namespace
	request.RepositoryName = &repositoryName

	if briefDesc != "" {
		request.BriefDescription = &briefDesc
	}

	if description != "" {
		request.Description = &description
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateRepository(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return
}

func (me *TCRService) ModifyTCRRepository(ctx context.Context, instanceId string, namespace string, repositoryName string, briefDesc string, description string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewModifyRepositoryRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.NamespaceName = &namespace
	request.RepositoryName = &repositoryName
	if briefDesc != "" {
		request.BriefDescription = &briefDesc
	}

	if description != "" {
		request.Description = &description
	}

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().ModifyRepository(request)
	return err
}

func (me *TCRService) DeleteTCRRepository(ctx context.Context, instanceId string, namespace string, repositoryName string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDeleteRepositoryRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.NamespaceName = &namespace
	request.RepositoryName = &repositoryName

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().DeleteRepository(request)
	return err
}

func (me *TCRService) DescribeTCRRepositories(ctx context.Context, instanceId string, namespace string, repositoryName string) (repositoryList []*tcr.TcrRepositoryInfo, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeRepositoriesRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	if namespace != "" {
		request.NamespaceName = &namespace
	}
	if repositoryName != "" {
		request.RepositoryName = &repositoryName
	}
	var offset, limit int64 = 0, 200
	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTCRClient().DescribeRepositories(request)
		if err != nil {
			ee, ok := err.(*sdkErrors.TencentCloudSDKError)
			if !ok {
				errRet = err
				return
			}
			if ee.Code == "ResourceNotFound" {
				errRet = nil
			} else {
				errRet = err
			}
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		repositoryList = append(repositoryList, response.Response.RepositoryList...)
		if len(response.Response.RepositoryList) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TCRService) DescribeTCRRepositoryById(ctx context.Context, instanceId string, namespace string, repositoryName string) (repository *tcr.TcrRepositoryInfo, has bool, errRet error) {
	repositories, err := me.DescribeTCRRepositories(ctx, instanceId, namespace, repositoryName)
	if err != nil {
		return nil, false, err
	}

	if len(repositories) == 0 {
		return nil, has, nil
	} else if len(repositories) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than 1 namespaces, %s %s %s", instanceId, namespace, repositoryName)
		return
	}

	repository = repositories[0]
	has = true
	return
}

//longterm token
func (me *TCRService) CreateTCRLongTermToken(ctx context.Context, instanceId string, description string) (tokenId string, token string, userName string, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateInstanceTokenRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.TokenType = helper.String("longterm")
	request.Desc = &description

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateInstanceToken(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	tokenId = *response.Response.TokenId
	token = *response.Response.Token
	userName = *response.Response.Username
	return
}

func (me *TCRService) ModifyTCRLongTermToken(ctx context.Context, instanceId string, tokenId string, enable bool) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewModifyInstanceTokenRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.TokenId = &tokenId
	request.Enable = &enable

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().ModifyInstanceToken(request)
	return err
}

func (me *TCRService) DeleteTCRLongTermToken(ctx context.Context, instanceId string, tokenId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDeleteInstanceTokenRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.TokenId = &tokenId

	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().DeleteInstanceToken(request)
	return err
}

func (me *TCRService) DescribeTCRTokens(ctx context.Context, instanceId string, tokenId string) (tokenList []*tcr.TcrInstanceToken, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeInstanceTokenRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId

	var offset, limit int64 = 0, 200
	for {
		request.Offset = &offset
		request.Limit = &limit
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseTCRClient().DescribeInstanceToken(request)
		if err != nil {
			errRet = err
			return
		}
		if response == nil || response.Response == nil {
			errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		}
		if tokenId == "" {
			tokenList = append(tokenList, response.Response.Tokens...)
		} else {
			for _, v := range response.Response.Tokens {
				if *v.Id == tokenId {
					tokenList = append(tokenList, v)
				}
			}
		}
		if len(response.Response.Tokens) < int(limit) {
			return
		}
		offset += limit
	}
}

func (me *TCRService) DescribeTCRLongTermTokenById(ctx context.Context, instanceId string, tokenId string) (token *tcr.TcrInstanceToken, has bool, errRet error) {
	tokens, err := me.DescribeTCRTokens(ctx, instanceId, tokenId)
	if err != nil {
		return nil, false, err
	}

	if len(tokens) == 0 {
		return nil, has, nil
	} else if len(tokens) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than 1 namespaces, %s %s", instanceId, tokenId)
		return
	}

	token = tokens[0]
	has = true
	return
}

//VPC attachment
func (me *TCRService) CreateTCRVPCAttachment(ctx context.Context, instanceId string, vpcId string, subnetId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewManageInternalEndpointRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.VpcId = &vpcId
	request.SubnetId = &subnetId
	request.Operation = helper.String("Create")

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().ManageInternalEndpoint(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *TCRService) DeleteTCRVPCAttachment(ctx context.Context, instanceId string, vpcId string, subnetId string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewManageInternalEndpointRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId
	request.VpcId = &vpcId
	request.SubnetId = &subnetId
	request.Operation = helper.String("Delete")

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().ManageInternalEndpoint(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *TCRService) DescribeTCRVPCAttachments(ctx context.Context, instanceId string, vpcId string, subnetId string) (vpcList []*tcr.AccessVpc, errRet error) {
	logId := getLogId(ctx)
	//sdk has internal error as invalid instance id para result
	//to avoid error code check
	//check instance exist first
	_, insHas, err := me.DescribeTCRInstanceById(ctx, instanceId)
	if err != nil {
		errRet = err
		return
	}
	if !insHas {
		return
	}

	request := tcr.NewDescribeInternalEndpointsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeInternalEndpoints(request)
	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	for _, v := range response.Response.AccessVpcSet {
		if vpcId != "" && *v.VpcId != vpcId {
			continue
		}
		if subnetId != "" && *v.SubnetId != subnetId {
			continue
		}
		vpcList = append(vpcList, v)
	}
	return
}

func (me *TCRService) DescribeTCRVPCAttachmentById(ctx context.Context, instanceId string, vpcId string, subnetId string) (vpcAccess *tcr.AccessVpc, has bool, errRet error) {
	vpcAccesses, err := me.DescribeTCRVPCAttachments(ctx, instanceId, vpcId, subnetId)
	if err != nil {
		return nil, false, err
	}

	if len(vpcAccesses) == 0 {
		return nil, has, nil
	} else if len(vpcAccesses) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than 1 namespaces, %s %s %s", instanceId, vpcId, subnetId)
		return
	}

	vpcAccess = vpcAccesses[0]
	has = true
	return
}
