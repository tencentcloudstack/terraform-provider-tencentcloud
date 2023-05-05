package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/pkg/errors"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tcr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tcr/v20190924"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type TCRService struct {
	client *connectivity.TencentCloudClient
}

func (me *TCRService) CreateTCRInstance(ctx context.Context, name string, instanceType string, params map[string]interface{}) (instanceId string, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryName = &name
	request.RegistryType = &instanceType

	var chargeType, period, renewFlag int
	if v, ok := params["registry_charge_type"]; ok {
		chargeType = v.(int)
		request.RegistryChargeType = helper.IntInt64(chargeType - 1)
	}

	if v, ok := params["instance_charge_type_prepaid_period"]; ok {
		period = v.(int)
	}
	if v, ok := params["instance_charge_type_prepaid_renew_flag"]; ok {
		renewFlag = v.(int)
	}
	if chargeType == 2 {
		if period == 0 || renewFlag == 0 {
			errRet = errors.New("Must set instance_charge_type_prepaid_period and instance_charge_type_prepaid_renew_flag when registry_charge_type is postpaid")
			return
		}
		request.RegistryChargePrepaid = &tcr.RegistryChargePrepaid{
			Period:    helper.IntInt64(period),
			RenewFlag: helper.IntInt64(renewFlag - 1),
		}
	}

	var tags map[string]string
	if v, ok := params["tags"]; ok {
		tags = v.(map[string]string)
	}
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

func (me *TCRService) ManageTCRExternalEndpoint(ctx context.Context, instanceId, operation string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewManageExternalEndpointRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.Operation = &operation
	request.RegistryId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().ManageExternalEndpoint(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil || response.Response.RegistryId == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return
}

func (me *TCRService) DescribeExternalEndpointStatus(ctx context.Context, instanceId string) (status string, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeExternalEndpointStatusRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = &instanceId

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeExternalEndpointStatus(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
		return
	}
	if response.Response.Status == nil {
		errRet = fmt.Errorf("TencentCloud SDK return more than one instances, instanceId %s, %s", instanceId, request.GetAction())
		return
	}
	has = true
	status = *response.Response.Status
	return
}

func (me *TCRService) DescribeSecurityPolicies(ctx context.Context, request *tcr.DescribeSecurityPoliciesRequest) (policies []*tcr.SecurityPolicy, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeSecurityPolicies(request)

	if err != nil {
		errRet = err
		return
	}

	if response.Response.SecurityPolicySet != nil {
		policies = response.Response.SecurityPolicySet
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

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
func (me *TCRService) CreateTCRNameSpace(ctx context.Context, instanceId string, name string, isPublic, isAutoScan, isPreventVUL bool, severity string, whitelistItems []interface{}) (errRet error) {
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
	request.IsAutoScan = &isAutoScan
	request.IsPreventVUL = &isPreventVUL
	if severity != "" {
		request.Severity = &severity
	}

	if len(whitelistItems) > 0 {
		for _, item := range whitelistItems {
			whitelistItemMap := item.(map[string]interface{})
			whitelistItemItem := tcr.CVEWhitelistItem{}
			if v, ok := whitelistItemMap["cve_id"]; ok {
				whitelistItemItem.CVEID = helper.String(v.(string))
			}
			request.CVEWhitelistItems = append(request.CVEWhitelistItems, &whitelistItemItem)
		}
	}

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

func (me *TCRService) ModifyInstance(ctx context.Context, registryId, registryType string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewModifyInstanceRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.RegistryId = helper.String(registryId)
	request.RegistryType = helper.String(registryType)
	ratelimit.Check(request.GetAction())
	_, err := me.client.UseTCRClient().ModifyInstance(request)
	return err

}
func (me *TCRService) ModifyTCRNameSpace(ctx context.Context, instanceId string, name string, isPublic, isAutoScan, isPreventVUL bool, severity string, whitelistItems []interface{}) (errRet error) {
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
	request.IsAutoScan = &isAutoScan
	request.IsPreventVUL = &isPreventVUL
	if severity != "" {
		request.Severity = &severity
	}

	if len(whitelistItems) > 0 {
		for _, item := range whitelistItems {
			whitelistItemMap := item.(map[string]interface{})
			whitelistItemItem := tcr.CVEWhitelistItem{}
			if v, ok := whitelistItemMap["cve_id"]; ok {
				whitelistItemItem.CVEID = helper.String(v.(string))
			}
			request.CVEWhitelistItems = append(request.CVEWhitelistItems, &whitelistItemItem)
		}
	}

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
	}

	for i := range namespaces {
		if name == *namespaces[i].Name {
			namespace = namespaces[i]
			has = true
			return
		}
	}
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
	}

	for i := range repositories {
		if *repositories[i].Name == namespace+"/"+repositoryName {
			repository = repositories[i]
			has = true
			return
		}
	}

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
func (me *TCRService) CreateTCRVPCAttachment(ctx context.Context, instanceId string, vpcId string,
	subnetId string, regionId int64, regionName string) (errRet error) {
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
	request.RegionId = helper.Int64Uint64(regionId)
	request.RegionName = &regionName

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

func (me *TCRService) DeleteTCRVPCAttachment(ctx context.Context, instanceId string, vpcId string,
	subnetId string, regionId int, regionName string) (errRet error) {
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
	request.RegionId = helper.IntUint64(regionId)
	request.RegionName = &regionName

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

func (me *TCRService) CreateTcrVpcDns(ctx context.Context, instanceId string, vpcId string, accessIp string, usePublicDomain bool, regionName string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewCreateInternalEndpointDnsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.VpcId = &vpcId
	request.EniLBIp = &accessIp
	request.UsePublicDomain = &usePublicDomain
	request.RegionName = &regionName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateInternalEndpointDns(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *TCRService) DeleteTcrVpcDns(ctx context.Context, instanceId string, vpcId string, accessIp string, usePublicDomain bool, regionName string) (errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDeleteInternalEndpointDnsRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.InstanceId = &instanceId
	request.VpcId = &vpcId
	request.EniLBIp = &accessIp
	request.UsePublicDomain = &usePublicDomain
	request.RegionName = &regionName

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DeleteInternalEndpointDns(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	return
}

func (me *TCRService) DescribeTcrVpcDnsById(ctx context.Context, instanceId string, vpcId string, accessIp string, usePublicDomain bool) (vpcPrivateDomainStatus *tcr.VpcPrivateDomainStatus, has bool, errRet error) {
	logId := getLogId(ctx)
	request := tcr.NewDescribeInternalEndpointDnsStatusRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail,reason[%s]", logId, request.GetAction(), errRet.Error())
		}
	}()
	request.VpcSet = append(request.VpcSet, &tcr.VpcAndDomainInfo{
		InstanceId:      &instanceId,
		VpcId:           &vpcId,
		EniLBIp:         &accessIp,
		UsePublicDomain: &usePublicDomain,
	})

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeInternalEndpointDnsStatus(request)
	if err != nil {
		errRet = err
		return
	}
	if response == nil || response.Response == nil {
		errRet = fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}
	if len(response.Response.VpcSet) == 0 {
		return nil, has, nil
	} else if len(response.Response.VpcSet) > 1 {
		errRet = fmt.Errorf("TencentCloud SDK return more than 1 vpcPrivateDomainStatus, %s %s %s %t", instanceId, vpcId, accessIp, usePublicDomain)
		return
	}

	vpcPrivateDomainStatus = response.Response.VpcSet[0]
	has = true
	return
}

func (me *TCRService) CreateReplicationInstance(ctx context.Context, request *tcr.CreateReplicationInstanceRequest) (id string, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().CreateReplicationInstance(request)

	if err != nil {
		errRet = err
		return
	}

	id = *response.Response.ReplicationRegistryId

	startPolling := false

	err = resource.Retry(readRetryTimeout*5, func() *resource.RetryError {
		req := tcr.NewDescribeReplicationInstancesRequest()
		req.RegistryId = request.RegistryId
		req.Limit = helper.IntInt64(100)
		replicas, err := me.DescribeReplicationInstances(ctx, req)
		if err != nil {
			return retryError(err)
		}
		if len(replicas) == 0 {
			return resource.NonRetryableError(fmt.Errorf("no replica found in registry %s", *request.RegistryId))
		}

		for i := range replicas {
			item := replicas[i]
			if *item.Status == "Running" {
				continue
			}
			startPolling = true
			return resource.RetryableError(fmt.Errorf("replica %s is %s, waiting for task finish", *request.RegistryId, *item.Status))
		}

		if !startPolling {
			return resource.RetryableError(fmt.Errorf("waiting for polling start"))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DeleteReplicationInstance(ctx context.Context, request *tcr.DeleteReplicationInstanceRequest) (errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DeleteReplicationInstance(request)

	if err != nil {
		errRet = err
		return
	}

	startPolling := false

	err = resource.Retry(readRetryTimeout*3, func() *resource.RetryError {
		req := tcr.NewDescribeReplicationInstancesRequest()
		req.RegistryId = request.RegistryId
		req.Limit = helper.IntInt64(100)
		replicas, err := me.DescribeReplicationInstances(ctx, req)
		if err != nil {
			return retryError(err)
		}

		if len(replicas) == 0 {
			return nil
		}

		for i := range replicas {
			item := replicas[i]
			if *item.Status == "Running" {
				continue
			}
			startPolling = true
			return resource.RetryableError(fmt.Errorf("replica %s is %s, waiting for task finish", *request.RegistryId, *item.Status))
		}

		if !startPolling {
			return resource.RetryableError(fmt.Errorf("waiting for polling start"))
		}

		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DescribeReplicationInstances(ctx context.Context, request *tcr.DescribeReplicationInstancesRequest) (list []*tcr.ReplicationRegistry, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseTCRClient().DescribeReplicationInstances(request)

	if err != nil {
		errRet = err
		return
	}

	list = response.Response.ReplicationRegistries

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DescribeTcrTagRetentionRuleById(ctx context.Context, registryId, namespaceName string, retentionId *string) (TagRetentionRule *tcr.RetentionPolicy, errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDescribeTagRetentionRulesRequest()
	request.RegistryId = &registryId
	request.NamespaceName = &namespaceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DescribeTagRetentionRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.RetentionPolicyList) < 1 {
		return
	}

	if retentionId != nil {
		for _, policy := range response.Response.RetentionPolicyList {
			if *policy.RetentionId == helper.StrToInt64(*retentionId) {
				TagRetentionRule = policy
				return
			}
		}
		return nil, fmt.Errorf("[ERROR]%sThe TagRetentionRules[%v] not found in the qurey results. \n", logId, *retentionId)
	}

	TagRetentionRule = response.Response.RetentionPolicyList[0]
	return
}

func (me *TCRService) DeleteTcrTagRetentionRuleById(ctx context.Context, registryId string, retentionId string) (errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDeleteTagRetentionRuleRequest()
	request.RegistryId = &registryId
	request.RetentionId = helper.StrToInt64Point(retentionId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DeleteTagRetentionRule(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DescribeTcrWebhookTriggerById(ctx context.Context, registryId string, triggerId int64, namespaceName string) (WebhookTrigger *tcr.WebhookTrigger, errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDescribeWebhookTriggerRequest()
	request.RegistryId = &registryId
	request.Namespace = &namespaceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DescribeWebhookTrigger(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Triggers) < 1 {
		return
	}

	for _, trigger := range response.Response.Triggers {
		if *trigger.Id == triggerId {
			WebhookTrigger = trigger
			return
		}
	}

	return
}

func (me *TCRService) DeleteTcrWebhookTriggerById(ctx context.Context, registryId string, namespaceName string, triggerId int64) (errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDeleteWebhookTriggerRequest()
	request.RegistryId = &registryId
	request.Id = &triggerId
	request.Namespace = &namespaceName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DeleteWebhookTrigger(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DescribeTcrWebhookTriggerLogByFilter(ctx context.Context, param map[string]interface{}) (DescribeWebhookTriggerLog []*tcr.WebhookTriggerLog, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tcr.NewDescribeWebhookTriggerLogRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "registry_id" {
			request.RegistryId = v.(*string)
		}
		if k == "namespace" {
			request.Namespace = v.(*string)
		}
		if k == "trigger_id" {
			request.Id = v.(*int64)
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
		response, err := me.client.UseTCRClient().DescribeWebhookTriggerLog(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Logs) < 1 {
			break
		}
		DescribeWebhookTriggerLog = append(DescribeWebhookTriggerLog, response.Response.Logs...)
		if len(response.Response.Logs) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TCRService) DescribeTcrCustomizedDomainById(ctx context.Context, registryId string, domainName *string) (CustomizedDomain []*tcr.CustomizedDomainInfo, errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDescribeInstanceCustomizedDomainRequest()
	request.RegistryId = &registryId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DescribeInstanceCustomizedDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainInfoList) < 1 {
		return
	}

	if domainName != nil {
		for _, domain := range response.Response.DomainInfoList {
			if *domain.DomainName == *domainName {
				CustomizedDomain = []*tcr.CustomizedDomainInfo{domain}
				return
			}
		}
	} else {
		CustomizedDomain = response.Response.DomainInfoList
	}

	return
}

func (me *TCRService) DeleteTcrCustomizedDomainById(ctx context.Context, registryId string, domainName string) (errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDeleteInstanceCustomizedDomainRequest()
	request.RegistryId = &registryId
	request.DomainName = &domainName

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DeleteInstanceCustomizedDomain(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) TcrCustomizedDomainStateRefreshFunc(registryId, domainName string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := contextNil

		object, err := me.DescribeTcrCustomizedDomainById(ctx, registryId, &domainName)

		if err != nil {
			return nil, "", err
		}

		return object, helper.PString(object[0].Status), nil
	}
}

func (me *TCRService) DescribeTcrImmutableTagRuleById(ctx context.Context, registryId string, namespaceName, ruleId *string) (ImmutableTagRules []*tcr.ImmutableTagRule, errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDescribeImmutableTagRulesRequest()
	request.RegistryId = &registryId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DescribeImmutableTagRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Rules) < 1 {
		return
	}

	// filter by ns and rule id
	if ruleId != nil && namespaceName != nil {
		targetId := helper.StrToInt64Point(*ruleId)

		for _, rule := range response.Response.Rules {
			if *targetId == *rule.RuleId && *namespaceName == *rule.NsName {
				ImmutableTagRules = []*tcr.ImmutableTagRule{rule}
			}
		}
		return
	}

	// only specify ns
	if namespaceName != nil {
		for _, rule := range response.Response.Rules {
			if *namespaceName == *rule.NsName {
				ImmutableTagRules = append(ImmutableTagRules, rule)
			}
		}
		return
	}

	ImmutableTagRules = response.Response.Rules
	return
}

func (me *TCRService) DeleteTcrImmutableTagRuleById(ctx context.Context, registryId string, namespaceName string, ruleId string) (errRet error) {
	logId := getLogId(ctx)

	request := tcr.NewDeleteImmutableTagRulesRequest()
	request.RegistryId = &registryId
	request.NamespaceName = &namespaceName
	request.RuleId = helper.StrToInt64Point(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTCRClient().DeleteImmutableTagRules(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TCRService) DescribeTcrImagesByFilter(ctx context.Context, param map[string]interface{}) (Images []*tcr.TcrImageInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = tcr.NewDescribeImagesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "registry_id" {
			request.RegistryId = v.(*string)
		}
		if k == "namespace_name" {
			request.NamespaceName = v.(*string)
		}
		if k == "repository_name" {
			request.RepositoryName = v.(*string)
		}
		if k == "image_version" {
			request.ImageVersion = v.(*string)
		}
		if k == "digest" {
			request.Digest = v.(*string)
		}
		if k == "exact_match" {
			request.ExactMatch = v.(*bool)
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
		response, err := me.client.UseTCRClient().DescribeImages(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.ImageInfoList) < 1 {
			break
		}
		Images = append(Images, response.Response.ImageInfoList...)
		if len(response.Response.ImageInfoList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}
