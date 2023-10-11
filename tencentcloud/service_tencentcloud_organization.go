package tencentcloud

import (
	"context"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OrganizationService struct {
	client *connectivity.TencentCloudClient
}

func (me *OrganizationService) DescribeOrganizationOrgNode(ctx context.Context, nodeId string) (orgNode *organization.OrgNode, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	var offset int64 = 0
	var pageSize int64 = 50
	instances := make([]*organization.OrgNode, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().DescribeOrganizationNodes(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		instances = append(instances, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}

	for _, instance := range instances {
		if helper.Int64ToStr(*instance.NodeId) == nodeId {
			orgNode = instance
		}
	}

	return
}

func (me *OrganizationService) DeleteOrganizationOrgNodeById(ctx context.Context, nodeId string) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDeleteOrganizationNodesRequest()

	request.NodeId = []*int64{helper.Int64(helper.StrToInt64(nodeId))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseOrganizationClient().DeleteOrganizationNodes(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgMember(ctx context.Context, uin string) (orgMember *organization.OrgMember, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationMembersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	var offset uint64 = 0
	var pageSize uint64 = 50
	instances := make([]*organization.OrgMember, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().DescribeOrganizationMembers(request)
		if err != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), err.Error())
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		instances = append(instances, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}
		offset += pageSize
	}

	if len(instances) < 1 {
		return
	}

	for _, instance := range instances {
		if helper.Int64ToStr(*instance.MemberUin) == uin {
			orgMember = instance
		}
	}

	return

}

func (me *OrganizationService) DeleteOrganizationOrgMemberById(ctx context.Context, uin string) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDeleteOrganizationMembersRequest()

	request.MemberUin = []*int64{helper.Int64(helper.StrToInt64(uin))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseOrganizationClient().DeleteOrganizationMembers(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationPolicySubAccountAttachment(ctx context.Context, policyId, memberUin string) (policySubAccountAttachment *organization.OrgMemberAuthAccount, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationMemberAuthAccountsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()
	request.PolicyId = helper.StrToInt64Point(policyId)
	request.MemberUin = helper.StrToInt64Point(memberUin)
	request.Limit = helper.IntInt64(50)
	request.Offset = helper.IntInt64(0)

	response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberAuthAccounts(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	if len(response.Response.Items) < 1 {
		return
	}
	policySubAccountAttachment = response.Response.Items[0]
	return
}

func (me *OrganizationService) DeleteOrganizationPolicySubAccountAttachmentById(ctx context.Context, policyId, memberUin, orgSubAccountUin string) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewCancelOrganizationMemberAuthAccountRequest()

	request.PolicyId = helper.StrToInt64Point(policyId)
	request.MemberUin = helper.StrToInt64Point(memberUin)
	request.OrgSubAccountUin = helper.StrToInt64Point(orgSubAccountUin)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseOrganizationClient().CancelOrganizationMemberAuthAccount(request)
	if err != nil {
		errRet = err
		return err
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgMemberAuthIdentityById(ctx context.Context, memberUin int64) (identityIds []int64, errRet error) {
	logId := getLogId(ctx)
	request := organization.NewDescribeOrganizationMemberAuthIdentitiesRequest()
	request.MemberUin = &memberUin

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	var (
		offset int64 = 0
		limit  int64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberAuthIdentities(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if len(response.Response.Items) < 1 {
			return
		}

		for _, v := range response.Response.Items {
			if v.MemberUin != nil && *v.MemberUin == memberUin {
				if *v.IdentityId == 1 {
					continue
				}
				identityIds = append(identityIds, *v.IdentityId)
			}
		}
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DeleteOrganizationOrgMemberAuthIdentityById(ctx context.Context, memberUin string, identityId []string) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDeleteOrganizationMemberAuthIdentityRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	for _, id := range identityId {
		log.Printf("[DEBUG]%s api[%s] delete identity, uin [%s], identityId [%s]\n", logId, request.GetAction(), memberUin, id)

		request.MemberUin = helper.StrToUint64Point(memberUin)
		request.IdentityId = helper.StrToUint64Point(id)

		response, err := me.client.UseOrganizationClient().DeleteOrganizationMemberAuthIdentity(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] delete identity success, request memberUin [%s], response body [%s]\n", logId, request.GetAction(), memberUin, response.ToJsonString())
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrgAuthNodeByFilter(ctx context.Context, param map[string]interface{}) (orgAuthNode []*organization.AuthNode, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationAuthNodeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AuthName" {
			request.AuthName = v.(*string)
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
		response, err := me.client.UseOrganizationClient().DescribeOrganizationAuthNode(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		orgAuthNode = append(orgAuthNode, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrganizationById(ctx context.Context) (result *organization.DescribeOrganizationResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDescribeOrganizationRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DescribeOrganization(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	result = response.Response
	return
}

func (me *OrganizationService) DeleteOrganizationOrganizationById(ctx context.Context) (errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDeleteOrganizationRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeleteOrganization(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgMemberEmailById(ctx context.Context, memberUin int64, bindId uint64) (orgMemberEmail *organization.DescribeOrganizationMemberEmailBindResponseParams, errRet error) {
	logId := getLogId(ctx)

	request := organization.NewDescribeOrganizationMemberEmailBindRequest()
	request.MemberUin = &memberUin

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DescribeOrganizationMemberEmailBind(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}
	if *response.Response.BindId != bindId {
		return
	}
	orgMemberEmail = response.Response
	return
}

func (me *OrganizationService) DescribeOrganizationOrgFinancialByMemberByFilter(ctx context.Context, param map[string]interface{}) (orgFinancialByMember *organization.DescribeOrganizationFinancialByMemberResponseParams, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = organization.NewDescribeOrganizationFinancialByMemberRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Month" {
			request.Month = v.(*string)
		}
		if k == "EndMonth" {
			request.EndMonth = v.(*string)
		}
		if k == "MemberUins" {
			request.MemberUins = v.([]*int64)
		}
		if k == "ProductCodes" {
			request.ProductCodes = v.([]*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset    int64 = 0
		limit     int64 = 20
		items     []*organization.OrgMemberFinancial
		totalCost float64 = 0
		total     int64   = 0
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeOrganizationFinancialByMember(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		items = append(items, response.Response.Items...)
		if response.Response != nil && response.Response.TotalCost != nil && totalCost == 0 && total == 0 {
			totalCost = *response.Response.TotalCost
			total = *response.Response.Total
		}
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}
	orgFinancialByMember = &organization.DescribeOrganizationFinancialByMemberResponseParams{
		TotalCost: &totalCost,
		Items:     items,
		Total:     &total,
	}
	return
}
