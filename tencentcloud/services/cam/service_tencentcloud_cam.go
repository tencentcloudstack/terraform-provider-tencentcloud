package cam

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewCamService(client *connectivity.TencentCloudClient) CamService {
	return CamService{client: client}
}

type CamService struct {
	client *connectivity.TencentCloudClient
}

func (me *CamService) DescribeRoleById(ctx context.Context, roleId string) (camInstance *cam.RoleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewDescribeRoleListRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM) //to save in extension
	result := make([]*cam.RoleInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().DescribeRoleList(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, role := range response.Response.List {
			if *role.RoleId == roleId {
				result = append(result, role)
			}
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}

	if len(result) == 0 {
		return
	}
	camInstance = result[0]
	return
}

func (me *CamService) DescribeRolesByFilter(ctx context.Context, params map[string]interface{}) (roles []*cam.RoleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	//need travel
	request := cam.NewDescribeRoleListRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	roles = make([]*cam.RoleInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().DescribeRoleList(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		for _, role := range response.Response.List {
			if params["role_id"] != nil {
				if *role.RoleId != params["role_id"].(string) {
					continue
				}
			}
			if params["name"] != nil {
				if *role.RoleName != params["name"].(string) {
					continue
				}
			}
			if params["description"] != nil {
				if *role.Description != params["description"].(string) {
					continue
				}
			}
			roles = append(roles, role)
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) DeleteRoleById(ctx context.Context, roleId string) error {

	logId := tccommon.GetLogId(ctx)
	request := cam.NewDeleteRoleRequest()
	request.RoleId = &roleId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DeleteRole(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CamService) DeleteRoleByName(ctx context.Context, roleName string) error {

	logId := tccommon.GetLogId(ctx)
	request := cam.NewDeleteRoleRequest()
	request.RoleName = &roleName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DeleteRole(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return nil
}

func (me *CamService) decodeCamPolicyAttachmentId(id string) (instanceId string, policyId64 uint64, errRet error) {
	items := strings.Split(id, "#")
	if len(items) != 2 {
		return instanceId, policyId64, fmt.Errorf(" id is not exist %s", id)
	}
	instanceId = items[0]
	policyId, e := strconv.Atoi(items[1])
	if e != nil {
		errRet = e
		return
	}
	policyId64 = uint64(policyId)
	return
}

func (me *CamService) DescribeRolePolicyAttachmentByName(ctx context.Context, roleName string, params map[string]interface{}) (policyOfRole *cam.AttachedPolicyOfRole, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewListAttachedRolePoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	result := make([]*cam.AttachedPolicyOfRole, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.RoleName = &roleName
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedRolePolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") || errCode == "InvalidParameter.RoleNotExist" {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			policyName, ok := params["policy_name"]
			if ok && *policy.PolicyName == policyName.(string) {
				result = append(result, policy)
			}
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}

	if len(result) == 0 {
		return
	}
	policyOfRole = result[0]
	return
}

func (me *CamService) DescribeRolePolicyAttachmentById(ctx context.Context, rolePolicyAttachmentId string) (policyOfRole *cam.AttachedPolicyOfRole, errRet error) {
	logId := tccommon.GetLogId(ctx)
	roleId, policyId, e := me.decodeCamPolicyAttachmentId(rolePolicyAttachmentId)
	if e != nil {
		return nil, e
	}
	request := cam.NewListAttachedRolePoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	result := make([]*cam.AttachedPolicyOfRole, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.RoleId = &roleId
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedRolePolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") || errCode == "InvalidParameter.RoleNotExist" {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			if *policy.PolicyId == policyId {
				result = append(result, policy)
			}
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}

	if len(result) == 0 {
		return
	}
	policyOfRole = result[0]
	return
}

func (me *CamService) DescribeRolePolicyAttachmentsByFilter(ctx context.Context, params map[string]interface{}) (policyOfRoles []*cam.AttachedPolicyOfRole, errRet error) {
	logId := tccommon.GetLogId(ctx)
	roleId := params["role_id"].(string)
	request := cam.NewListAttachedRolePoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	policyOfRoles = make([]*cam.AttachedPolicyOfRole, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.RoleId = &roleId
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedRolePolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") || errCode == "InvalidParameter.RoleNotExist" {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			if params["policy_id"] != nil {
				if *policy.PolicyId != params["policy_id"].(uint64) {
					continue
				}
			}
			if params["policy_type"] != nil {
				if *policy.PolicyType != params["policy_type"].(string) {
					continue
				}
			}
			if params["create_mode"] != nil {
				if int(*policy.CreateMode) != params["create_mode"].(int) {
					continue
				}
			}
			policyOfRoles = append(policyOfRoles, policy)
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) DeleteRolePolicyAttachmentByName(ctx context.Context, roleName, policyName string) error {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewDetachRolePolicyRequest()
	request.DetachRoleName = &roleName
	request.PolicyName = &policyName
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DetachRolePolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DeleteRolePolicyAttachmentById(ctx context.Context, rolePolicyAttachmentId string) error {
	logId := tccommon.GetLogId(ctx)
	roleId, policyId, e := me.decodeCamPolicyAttachmentId(rolePolicyAttachmentId)
	if e != nil {
		return e
	}
	request := cam.NewDetachRolePolicyRequest()
	request.DetachRoleId = &roleId
	request.PolicyId = &policyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DetachRolePolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DescribeUserPolicyAttachmentById(ctx context.Context, userPolicyAttachmentId string) (policyResults *cam.AttachPolicyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	userId, policyId, e := me.decodeCamPolicyAttachmentId(userPolicyAttachmentId)
	if e != nil {
		return nil, e
	}
	user, err := me.DescribeUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Response == nil || user.Response.Uid == nil {
		return
	}
	uin := user.Response.Uin

	request := cam.NewListAttachedUserPoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	result := make([]*cam.AttachPolicyInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.TargetUin = uin
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedUserPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			if *policy.PolicyId == policyId {
				result = append(result, policy)
			}
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}

	if len(result) == 0 {
		return
	}
	policyResults = result[0]
	return
}

func (me *CamService) DescribeUserPolicyAttachmentsByFilter(ctx context.Context, params map[string]interface{}) (policyResults []*cam.AttachPolicyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	userId := params["user_id"].(string)
	user, err := me.DescribeUserById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if user == nil || user.Response == nil || user.Response.Uid == nil {
		return
	}
	uin := user.Response.Uin
	request := cam.NewListAttachedUserPoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	policyResults = make([]*cam.AttachPolicyInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.TargetUin = uin
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedUserPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			if params["policy_id"] != nil {
				if *policy.PolicyId != params["policy_id"].(uint64) {
					continue
				}
			}
			if params["policy_type"] != nil {
				if *policy.PolicyType != params["policy_type"].(string) {
					continue
				}
			}
			if params["create_mode"] != nil {
				if int(*policy.CreateMode) != params["create_mode"].(int) {
					continue
				}
			}
			policyResults = append(policyResults, policy)
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) AddUserPolicyAttachment(ctx context.Context, userId string, policyId string) error {
	logId := tccommon.GetLogId(ctx)

	user, err := me.DescribeUserById(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil || user.Response == nil || user.Response.Uid == nil {
		return nil
	}
	uin := user.Response.Uin
	policyIdInt, e := strconv.Atoi(policyId)
	if e != nil {
		return e
	}
	policyIdInt64 := uint64(policyIdInt)
	request := cam.NewAttachUserPolicyRequest()
	request.AttachUin = uin
	request.PolicyId = &policyIdInt64
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().AttachUserPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DeleteUserPolicyAttachmentById(ctx context.Context, userPolicyAttachmentId string) error {
	logId := tccommon.GetLogId(ctx)
	userId, policyId, e := me.decodeCamPolicyAttachmentId(userPolicyAttachmentId)
	if e != nil {
		return e
	}
	user, err := me.DescribeUserById(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil || user.Response == nil || user.Response.Uid == nil {
		return nil
	}
	uin := user.Response.Uin

	request := cam.NewDetachUserPolicyRequest()
	request.DetachUin = uin
	request.PolicyId = &policyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DetachUserPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DescribeGroupPolicyAttachmentById(ctx context.Context, groupPolicyAttachmentId string) (policyResults *cam.AttachPolicyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	groupId, policyId, e := me.decodeCamPolicyAttachmentId(groupPolicyAttachmentId)
	if e != nil {
		errRet = e
		return
	}
	groupIdInt, ee := strconv.Atoi(groupId)
	if ee != nil {
		errRet = ee
		return
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewListAttachedGroupPoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	result := make([]*cam.AttachPolicyInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.TargetGroupId = &groupIdInt64
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedGroupPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		for _, policy := range response.Response.List {
			if *policy.PolicyId == policyId {
				result = append(result, policy)
			}
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	if len(result) == 0 {
		return
	}
	policyResults = result[0]
	return
}

func (me *CamService) DescribeGroupPolicyAttachmentsByFilter(ctx context.Context, params map[string]interface{}) (policyResults []*cam.AttachPolicyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	groupId := params["group_id"].(string)
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		errRet = e
		return
	}
	groupIdInt64 := uint64(groupIdInt)
	request := cam.NewListAttachedGroupPoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	policyResults = make([]*cam.AttachPolicyInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		request.TargetGroupId = &groupIdInt64
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedGroupPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}

		for _, policy := range response.Response.List {
			if params["policy_id"] != nil {
				if *policy.PolicyId != params["policy_id"].(uint64) {
					continue
				}
			}
			if params["policy_type"] != nil {
				if *policy.PolicyType != params["policy_type"].(string) {
					continue
				}
			}
			if params["create_mode"] != nil {
				if int(*policy.CreateMode) != params["create_mode"].(int) {
					continue
				}
			}
			policyResults = append(policyResults, policy)
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) AddGroupPolicyAttachment(ctx context.Context, groupId string, policyId string) error {
	logId := tccommon.GetLogId(ctx)

	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		return e
	}
	groupIdInt64 := uint64(groupIdInt)
	policyIdInt, ee := strconv.Atoi(policyId)
	if ee != nil {
		return ee
	}
	policyIdInt64 := uint64(policyIdInt)

	request := cam.NewAttachGroupPolicyRequest()
	request.AttachGroupId = &groupIdInt64
	request.PolicyId = &policyIdInt64
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().AttachGroupPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DeleteGroupPolicyAttachmentById(ctx context.Context, groupPolicyAttachmentId string) error {
	logId := tccommon.GetLogId(ctx)
	groupId, policyId, e := me.decodeCamPolicyAttachmentId(groupPolicyAttachmentId)
	if e != nil {
		return e
	}
	groupIdInt, ee := strconv.Atoi(groupId)
	if ee != nil {
		return ee
	}
	groupIdInt64 := uint64(groupIdInt)

	request := cam.NewDetachGroupPolicyRequest()
	request.DetachGroupId = &groupIdInt64
	request.PolicyId = &policyId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DetachGroupPolicy(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return nil
}

func (me *CamService) DescribePolicyById(ctx context.Context, policyId string) (result *cam.GetPolicyResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewGetPolicyRequest()
	policyIdInt, e := strconv.Atoi(policyId)
	if e != nil {
		errRet = e
		return
	}
	policyIdInt64 := uint64(policyIdInt)
	request.PolicyId = &policyIdInt64
	result, err := me.client.UseCamClient().GetPolicy(request)

	if err != nil {
		log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		return nil, err
	} else {
		if result == nil || result.Response == nil || result.Response.PolicyName == nil {
			return
		}
	}

	return
}

func (me *CamService) DescribePoliciesByFilter(ctx context.Context, params map[string]interface{}) (policies []*cam.StrategyInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	policyId := -1
	policyName := ""
	//notice this policy type is different from the policy attachment, this sdk returns int while the attachments returns string
	policyType := -1
	description := ""
	createMode := -1

	request := cam.NewListPoliciesRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)

	for k, v := range params {
		if k == "policy_id" {
			policyId = v.(int)
		}
		if k == "name" {
			policyName = v.(string)
		}
		if k == "type" {
			policyType = v.(int)
		}
		if k == "description" {
			description = v.(string)
		}
		if k == "create_mode" {
			createMode = v.(int)
		}
	}
	policies = make([]*cam.StrategyInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListPolicies(request)
		if err != nil {
			log.Printf("[CRITAL]%s read CAM policy failed, reason:%s\n", logId, err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		for _, policy := range response.Response.List {
			if policyId != -1 {
				if int(*policy.PolicyId) != policyId {
					continue
				}
			}
			if policyName != "" {
				if *policy.PolicyName != policyName {
					continue
				}
			}
			if policyType != -1 {
				if int(*policy.Type) != policyType {
					continue
				}
			}
			if description != "" {
				if *policy.Description != description {
					continue
				}
			}
			if createMode != -1 {
				if int(*policy.CreateMode) != createMode {
					continue
				}
			}
			policies = append(policies, policy)
		}
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) DescribeUserById(ctx context.Context, userId string) (response *cam.GetUserResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewGetUserRequest()
	request.Name = &userId
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().GetUser(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	return
}

func (me *CamService) DescribeUsersByFilter(ctx context.Context, params map[string]interface{}) (result []*cam.SubAccountInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewListUsersRequest()

	result = make([]*cam.SubAccountInfo, 0)

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().ListUsers(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	for _, user := range response.Response.Data {
		if params["name"] != nil {
			if params["name"].(string) != *user.Name {
				continue
			}
		}
		if params["remark"] != nil {
			if user.Remark != nil {
				if params["remark"].(string) != *user.Remark {
					continue
				}
			} else {
				continue
			}
		}
		if params["phone_num"] != nil {
			if user.PhoneNum != nil {
				if params["phone_num"].(string) != *user.PhoneNum {
					continue
				}
			} else {
				continue
			}
		}
		if params["country_code"] != nil {
			if user.CountryCode != nil {
				if params["country_code"].(string) != *user.CountryCode {
					continue
				}
			} else {
				continue
			}
		}
		if params["email"] != nil {
			if user.Email != nil {
				if params["email"].(string) != *user.Email {
					continue
				}
			} else {
				continue
			}
		}
		if params["uin"] != nil {
			if user.Uin != nil {
				if params["uin"].(int) != int(*user.Uin) {
					continue
				}
			} else {
				continue
			}
		}
		if params["uid"] != nil {
			if user.Uid != nil {
				if params["uid"].(int) != int(*user.Uid) {
					continue
				}
			} else {
				continue
			}
		}
		if params["console_login"] != nil {
			if user.ConsoleLogin != nil {
				if params["console_login"].(int) != int(*user.ConsoleLogin) {
					continue
				}
			} else {
				continue
			}
		}
		result = append(result, user)
	}

	return
}

func (me *CamService) DescribeGroupById(ctx context.Context, groupId string) (camInstance *cam.GetGroupResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewGetGroupRequest()
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		errRet = e
		return
	}
	groupIdInt64 := uint64(groupIdInt)
	request.GroupId = &groupIdInt64
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().GetGroup(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	camInstance = response
	return
}

func (me *CamService) DescribeGroupsByFilter(ctx context.Context, params map[string]interface{}) (groups []*cam.GroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewListGroupsRequest()
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	groups = make([]*cam.GroupInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListGroups(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.GroupInfo) < 1 {
			break
		}
		for _, group := range response.Response.GroupInfo {
			if params["group_id"] != nil {
				if int(*group.GroupId) != params["group_id"].(int) {
					continue
				}
			}
			if params["name"] != nil {
				if *group.GroupName != params["name"].(string) {
					continue
				}
			}
			if params["remark"] != nil {
				if group.Remark == nil || (group.Remark != nil && *group.Remark != params["remark"].(string)) {
					continue
				}
				log.Printf("in")
			}
			groups = append(groups, group)
		}
		if len(response.Response.GroupInfo) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) DescribeGroupMembershipById(ctx context.Context, groupId string) (members []*string, errRet error) {
	logId := tccommon.GetLogId(ctx)
	groupIdInt, e := strconv.Atoi(groupId)
	if e != nil {
		errRet = e
		return
	}
	groupIdInt64 := uint64(groupIdInt)
	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM)
	members = make([]*string, 0)
	request := cam.NewListUsersForGroupRequest()
	request.GroupId = &groupIdInt64
	for {
		request.Page = &pageStart
		request.Rp = &rp
		response, err := me.client.UseCamClient().ListUsersForGroup(request)
		if err != nil {
			log.Printf("[CRITAL]%s read CAM group membership failed, reason:%s\n", logId, err.Error())
			if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
				errCode := ee.GetCode()
				//check if read empty
				if strings.Contains(errCode, "ResourceNotFound") {
					return
				}
			}
			errRet = err
			return
		}

		if response == nil || len(response.Response.UserInfo) < 1 {
			break
		}
		for _, member := range response.Response.UserInfo {

			members = append(members, member.Name)
		}
		if len(response.Response.UserInfo) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	return
}

func (me *CamService) DescribeSAMLProviderById(ctx context.Context, providerName string) (result *cam.GetSAMLProviderResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)
	request := cam.NewGetSAMLProviderRequest()
	request.Name = &providerName
	result, err := me.client.UseCamClient().GetSAMLProvider(request)

	if err != nil {
		log.Printf("[CRITAL]%s read cam SAML provider failed, reason:%s\n", logId, err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		return nil, err
	} else {
		if result == nil || result.Response == nil || result.Response.Name == nil {
			return
		}
	}

	return
}

func (me *CamService) DescribeSAMLProvidersByFilter(ctx context.Context, params map[string]interface{}) (providers []*cam.SAMLProviderInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)
	name := ""
	description := ""

	request := cam.NewListSAMLProvidersRequest()

	for k, v := range params {
		if k == "name" {
			name = v.(string)
		}
		if k == "description" {
			description = v.(string)
		}
	}
	providers = make([]*cam.SAMLProviderInfo, 0)
	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().ListSAMLProviders(request)
	if err != nil {
		log.Printf("[CRITAL]%s read CAM SAML provider failed, reason:%s\n", logId, err.Error())
		if ee, ok := err.(*sdkErrors.TencentCloudSDKError); ok {
			errCode := ee.GetCode()
			//check if read empty
			if strings.Contains(errCode, "ResourceNotFound") {
				return
			}
		}
		errRet = err
		return
	}
	for _, provider := range response.Response.SAMLProviderSet {
		if name != "" {
			if *provider.Name != name {
				continue
			}
		}
		if description != "" {
			if *provider.Description != description {
				continue
			}
		}
		providers = append(providers, provider)
	}

	return
}

func (me *CamService) PolicyDocumentForceCheck(document string) error {

	//Policy syntax allows multi formats, but returns with only one format. In this case, the user's input may be different from the output value. To avoid this, terraform must make sure the syntax of the input policy document consists with the syntax of final returned output
	type Principal struct {
		Qcs []string `json:"qcs"`
	}
	type Statement struct {
		Resource interface{} `json:"resource"`
		//to avoid json unmarshal eats up with '/'
		Action    []json.RawMessage `json:"action"`
		Principal Principal         `json:"principal"`
	}
	type Document struct {
		Version   string      `json:"version"`
		Statement []Statement `json:"statement"`
	}
	var documentJson Document
	err := json.Unmarshal([]byte(document), &documentJson)
	if err != nil {
		return err
	}
	for _, state := range documentJson.Statement {
		//multi value case in elemant `resource`, `action`: input:""/[""], output: [""]
		if state.Resource != nil {
			if reflect.TypeOf(state.Resource) == reflect.TypeOf("string") {
				return fmt.Errorf("The format of `resource` in policy document is invalid, its type must be array")
			}
		}

		if state.Action != nil {
			if reflect.TypeOf(state.Action) == reflect.TypeOf("string") {
				return fmt.Errorf("The format of `resource` in policy document is invalid, its type must be array")
			}

		}
		//multi value case in elemant `principal.qcs`:input :root/[uin of the user], output:[uin of the user]
		for _, qcs := range state.Principal.Qcs {
			if strings.Contains(qcs, "root") {
				return fmt.Errorf("`root` format is not supported, please replace it with uin")
			}
		}
	}
	return nil
}

func (me *CamService) DescribeCamServiceLinkedRole(ctx context.Context, roleId string) (serviceLinkedRole *cam.RoleInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewGetRoleRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.RoleId = &roleId

	response, err := me.client.UseCamClient().GetRole(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		serviceLinkedRole = response.Response.RoleInfo
	}

	return
}

func (me *CamService) DeleteCamServiceLinkedRoleById(ctx context.Context, roleId string) (deletionTaskId string, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDeleteServiceLinkedRoleRequest()

	request.RoleName = &roleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseCamClient().DeleteServiceLinkedRole(request)
	if err != nil {
		errRet = err
		return "", err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response != nil && response.Response != nil {
		deletionTaskId = *response.Response.DeletionTaskId
	}
	return
}

func (me *CamService) DescribeCamServiceLinkedRoleDeleteStatus(ctx context.Context, deletionTaskId string) (response *cam.GetServiceLinkedRoleDeletionStatusResponse, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewGetServiceLinkedRoleDeletionStatusRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.DeletionTaskId = &deletionTaskId

	response, err := me.client.UseCamClient().GetServiceLinkedRoleDeletionStatus(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamUserSamlConfigById(ctx context.Context) (userSamlConfig *cam.DescribeUserSAMLConfigResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDescribeUserSAMLConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DescribeUserSAMLConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	userSamlConfig = response
	return
}

func (me *CamService) DeleteCamUserSamlConfigById(ctx context.Context) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewUpdateUserSAMLConfigRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Operate = helper.String("disable")

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().UpdateUserSAMLConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamMfaFlagById(ctx context.Context, id uint64) (loginFlag *cam.LoginActionFlag, actionFlag *cam.LoginActionFlag, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDescribeSafeAuthFlagCollRequest()
	request.SubUin = &id
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DescribeSafeAuthFlagColl(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.ActionFlag == nil && response.Response.LoginFlag == nil {
		return
	}

	loginFlag = response.Response.LoginFlag
	actionFlag = response.Response.ActionFlag
	return
}

func (me *CamService) DescribeCamAccessKeyById(ctx context.Context, targetUin uint64, accessKey string) (AccessKey *cam.AccessKey, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewListAccessKeysRequest()
	request.TargetUin = &targetUin

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().ListAccessKeys(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.AccessKeys) < 1 {
		return
	}

	for _, v := range response.Response.AccessKeys {
		if *v.AccessKeyId == accessKey {
			AccessKey = v
		}
	}
	return
}

func (me *CamService) DeleteCamAccessKeyById(ctx context.Context, uin, accessKeyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDeleteAccessKeyRequest()
	request.AccessKeyId = &accessKeyId
	request.TargetUin = helper.StrToUint64Point(uin)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DeleteAccessKey(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamUserPermissionBoundaryById(ctx context.Context, targetUin string) (UserPermissionBoundary *cam.GetUserPermissionBoundaryResponse, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewGetUserPermissionBoundaryRequest()
	request.TargetUin = helper.StrToInt64Point(targetUin)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetUserPermissionBoundary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil {
		return
	}

	UserPermissionBoundary = response
	return
}

func (me *CamService) DeleteCamUserPermissionBoundaryById(ctx context.Context, targetUin string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDeleteUserPermissionsBoundaryRequest()
	request.TargetUin = helper.StrToInt64Point(targetUin)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DeleteUserPermissionsBoundary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamPolicyVersionById(ctx context.Context, policyId uint64, versionId uint64) (policyVersion *cam.PolicyVersionDetail, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewGetPolicyVersionRequest()
	request.PolicyId = &policyId
	request.VersionId = &versionId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetPolicyVersion(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.PolicyVersion == nil {
		return
	}

	policyVersion = response.Response.PolicyVersion
	return
}

func (me *CamService) DeleteCamPolicyVersionById(ctx context.Context, policyId uint64, versionId uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDeletePolicyVersionRequest()
	request.PolicyId = &policyId
	request.VersionId = []*uint64{helper.Uint64(versionId)}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DeletePolicyVersion(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamListEntitiesForPolicyByFilter(ctx context.Context, param map[string]interface{}) (ListEntitiesForPolicy []*cam.AttachEntityOfPolicy, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewListEntitiesForPolicyRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "PolicyId" {
			request.PolicyId = v.(*uint64)
		}
		if k == "Rp" {
			request.Rp = v.(*uint64)
		}
		if k == "EntityFilter" {
			request.EntityFilter = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM) //to save in extension
	result := make([]*cam.AttachEntityOfPolicy, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListEntitiesForPolicy(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		result = append(result, response.Response.List...)
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}
	ListEntitiesForPolicy = result
	return
}

func (me *CamService) DescribeCamListAttachedUserPolicyByFilter(ctx context.Context, param map[string]interface{}) (ListAttachedUserPolicy []*cam.AttachedUserPolicy, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewListAttachedUserAllPoliciesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TargetUin" {
			request.TargetUin = v.(*uint64)
		}
		if k == "AttachType" {
			request.AttachType = v.(*uint64)
		}
		if k == "StrategyType" {
			request.StrategyType = v.(*uint64)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
	}

	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM) //to save in extension
	result := make([]*cam.AttachedUserPolicy, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseCamClient().ListAttachedUserAllPolicies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.PolicyList) < 1 {
			break
		}
		result = append(result, response.Response.PolicyList...)
		if len(response.Response.PolicyList) < PAGE_ITEM {
			break
		}

		pageStart += 1
	}
	ListAttachedUserPolicy = result
	return
}

func (me *CamService) DescribeCamTagRoleById(ctx context.Context, roleName, roleId string) (TagRole *cam.RoleInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewGetRoleRequest()
	if roleName == "" {
		request.RoleId = &roleId
	} else {
		request.RoleName = &roleName
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetRole(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || response.Response.RoleInfo == nil {
		return
	}
	TagRole = response.Response.RoleInfo
	return
}

func (me *CamService) DeleteCamTagRoleById(ctx context.Context, roleName, roleId string, keys []*string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewUntagRoleRequest()
	if roleName == "" {
		request.RoleId = &roleId
	} else {
		request.RoleName = &roleName
	}
	request.TagKeys = keys
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().UntagRole(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamRolePermissionBoundaryAttachmentById(ctx context.Context, roleId string, policyId string) (RolePermissionBoundaryAttachment *cam.GetRolePermissionBoundaryResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewGetRolePermissionBoundaryRequest()
	request.RoleId = &roleId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetRolePermissionBoundary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if response == nil || response.Response == nil {
		return
	}
	if *response.Response.PolicyId != helper.StrToInt64(policyId) {
		return
	}
	RolePermissionBoundaryAttachment = response.Response
	return
}

func (me *CamService) DeleteCamRolePermissionBoundaryAttachmentById(ctx context.Context, roleId string, roleName string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewDeleteRolePermissionsBoundaryRequest()
	if roleId == "" {
		request.RoleName = &roleName
	} else {
		request.RoleId = &roleId
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().DeleteRolePermissionsBoundary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *CamService) DescribeCamSecretLastUsedTimeByFilter(ctx context.Context, param map[string]interface{}) (SecretLastUsedTime []*cam.SecretIdLastUsed, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewGetSecurityLastUsedRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SecretIdList" {
			request.SecretIdList = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetSecurityLastUsed(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.SecretIdLastUsedRows) < 1 {
		return
	}
	SecretLastUsedTime = append(SecretLastUsedTime, response.Response.SecretIdLastUsedRows...)
	return
}

func (me *CamService) DescribeCamPolicyGrantingServiceAccessByFilter(ctx context.Context, param map[string]interface{}) (PolicyGrantingServiceAccess []*cam.ListGrantServiceAccessNode, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewListPoliciesGrantingServiceAccessRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "TargetUin" {
			request.TargetUin = v.(*uint64)
		}
		if k == "RoleId" {
			request.RoleId = v.(*uint64)
		}
		if k == "GroupId" {
			request.GroupId = v.(*uint64)
		}
		if k == "ServiceType" {
			request.ServiceType = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().ListPoliciesGrantingServiceAccess(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.List) < 1 {
		return
	}
	PolicyGrantingServiceAccess = response.Response.List
	return
}

func (me *CamService) DescribeCamSetPolicyVersionById(ctx context.Context, policyId, versionId string) (SetPolicyVersion *cam.PolicyVersionItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := cam.NewListPolicyVersionsRequest()
	request.PolicyId = helper.StrToUint64Point(policyId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().ListPolicyVersions(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.Versions) < 1 {
		return
	}
	for _, v := range response.Response.Versions {
		if *v.IsDefaultVersion == int64(1) && *v.VersionId == helper.StrToUInt64(versionId) {
			SetPolicyVersion = v
		}
	}

	return
}

func (me *CamService) DescribeCamAccountSummaryByFilter(ctx context.Context) (AccountSummary *cam.GetAccountSummaryResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewGetAccountSummaryRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseCamClient().GetAccountSummary(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	AccountSummary = response.Response
	return
}

func (me *CamService) DescribeCamGroupUserAccountByFilter(ctx context.Context, param map[string]interface{}) (GroupUserAccount []*cam.GroupInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = cam.NewListGroupsForUserRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Uid" {
			request.Uid = v.(*uint64)
		}
		if k == "SubUin" {
			request.SubUin = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM) //to save in extension
	result := make([]*cam.GroupInfo, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		response, err := me.client.UseCamClient().ListGroupsForUser(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.GroupInfo) < 1 {
			break
		}
		result = append(result, response.Response.GroupInfo...)
		if len(response.Response.GroupInfo) < PAGE_ITEM {
			break
		}

		pageStart += 1
	}
	GroupUserAccount = result
	return
}
