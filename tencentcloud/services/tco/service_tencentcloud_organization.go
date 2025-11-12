package tco

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	organization "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/organization/v20210331"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

type OrganizationService struct {
	client *connectivity.TencentCloudClient
}

func (me *OrganizationService) DescribeOrganizationOrgNode(ctx context.Context, nodeId string) (orgNode *organization.OrgNode, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId    = tccommon.GetLogId(ctx)
		request  = organization.NewDescribeOrganizationMembersRequest()
		response = organization.NewDescribeOrganizationMembersResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "query object", request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset    uint64 = 0
		pageSize  uint64 = 50
		instances        = make([]*organization.OrgMember, 0)
	)

	for {
		request.Offset = &offset
		request.Limit = &pageSize

		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().DescribeOrganizationMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe organization members failed, Response is nil."))
			}

			response = result
			return nil
		})

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
			break
		}
	}

	return
}

func (me *OrganizationService) DeleteOrganizationAccountById(ctx context.Context, uin string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteAccountRequest()
	request.MemberUin = helper.StrToInt64Point(uin)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOrganizationClient().DeleteAccount(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	return errRet
}

func (me *OrganizationService) DeleteOrganizationOrgMemberById(ctx context.Context, uin string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteOrganizationMembersRequest()
	request.MemberUin = []*int64{helper.Int64(helper.StrToInt64(uin))}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, "delete object", request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOrganizationClient().DeleteOrganizationMembers(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	return errRet
}

func (me *OrganizationService) DescribeOrganizationPolicySubAccountAttachment(ctx context.Context, policyId, memberUin string) (policySubAccountAttachment *organization.OrgMemberAuthAccount, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
	logId := tccommon.GetLogId(ctx)

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
		logId   = tccommon.GetLogId(ctx)
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

func (me *OrganizationService) DescribeOrganizationOrgFinancialByMonthByFilter(ctx context.Context, param map[string]interface{}) (orgFinancialByMonth []*organization.OrgFinancialByMonth, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewDescribeOrganizationFinancialByMonthRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
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

	response, err := me.client.UseOrganizationClient().DescribeOrganizationFinancialByMonth(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Items) < 1 {
		return
	}
	orgFinancialByMonth = append(orgFinancialByMonth, response.Response.Items...)
	return
}
func (me *OrganizationService) DescribeOrganizationOrgFinancialByProductByFilter(ctx context.Context, param map[string]interface{}) (orgFinancialByProduct *organization.DescribeOrganizationFinancialByProductResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewDescribeOrganizationFinancialByProductRequest()
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
		items     []*organization.OrgProductFinancial
		totalCost float64 = 0
		total     int64   = 0
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeOrganizationFinancialByProduct(request)
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
	orgFinancialByProduct = &organization.DescribeOrganizationFinancialByProductResponseParams{
		TotalCost: &totalCost,
		Items:     items,
		Total:     &total,
	}
	return
}

func (me *OrganizationService) DescribeOrganizationOrgIdentityById(ctx context.Context, identityId string) (orgIdentity *organization.OrgIdentity, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListOrganizationIdentityRequest()
	request.IdentityId = helper.StrToUint64Point(identityId)

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
	var tmp []*organization.OrgIdentity
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().ListOrganizationIdentity(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		tmp = append(tmp, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}
	}
	for _, item := range tmp {
		if *item.IdentityId == helper.StrToInt64(identityId) {
			orgIdentity = item
		}
	}
	return
}

func (me *OrganizationService) DeleteOrganizationOrgIdentityById(ctx context.Context, identityId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteOrganizationIdentityRequest()
	request.IdentityId = helper.StrToUint64Point(identityId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeleteOrganizationIdentity(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DeleteOrganizationOrgMemberPolicyAttachmentById(ctx context.Context, policyId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteOrganizationMembersPolicyRequest()
	request.PolicyId = helper.StrToUint64Point(policyId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeleteOrganizationMembersPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationMembersByFilter(ctx context.Context, param map[string]interface{}) (members []*organization.OrgMember, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewDescribeOrganizationMembersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Lang" {
			request.Lang = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
		if k == "AuthName" {
			request.AuthName = v.(*string)
		}
		if k == "Product" {
			request.Product = v.(*string)
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
		response, err := me.client.UseOrganizationClient().DescribeOrganizationMembers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		members = append(members, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitById(ctx context.Context, area, unitId string) (orgShareUnit *organization.ManagerShareUnit, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribeShareUnitsRequest()
	request.SearchKey = &unitId
	request.Area = &area
	request.Limit = helper.IntUint64(20)
	request.Offset = helper.IntUint64(0)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DescribeShareUnits(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.Items) < 1 {
		return
	}

	orgShareUnit = response.Response.Items[0]
	return
}

func (me *OrganizationService) DeleteOrganizationOrgShareUnitById(ctx context.Context, unitId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteShareUnitRequest()
	request.UnitId = &unitId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeleteShareUnit(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitMemberById(ctx context.Context, unitId, area, shareMemberUins string) (orgShareUnitMembers []*organization.ShareUnitMember, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribeShareUnitMembersRequest()
	request.UnitId = &unitId
	request.Area = &area

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	ratelimit.Check(request.GetAction())
	var (
		offset uint64 = 0
		limit  uint64 = 10
	)
	var tmp []*organization.ShareUnitMember
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeShareUnitMembers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		tmp = append(tmp, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	memberIdMap := make(map[string]struct{})

	for _, item := range tmp {
		memberIdMap[helper.Int64ToStr(*item.ShareMemberUin)] = struct{}{}
	}

	split := strings.Split(shareMemberUins, tccommon.COMMA_SP)
	if len(split) < 1 {
		errRet = fmt.Errorf("shareMemberUins is broken, %s", shareMemberUins)
		return
	}
	for _, v := range split {
		if _, ok := memberIdMap[v]; !ok {
			return
		}
	}
	orgShareUnitMembers = tmp
	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitMemberV2ById(ctx context.Context, unitId, area string) (orgShareUnitMembers []*organization.ShareUnitMember, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribeShareUnitMembersRequest()
	response := organization.NewDescribeShareUnitMembersResponse()
	request.UnitId = &unitId
	request.Area = &area

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset uint64 = 0
		limit  uint64 = 50
	)

	for {
		request.Offset = &offset
		request.Limit = &limit

		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().DescribeShareUnitMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			response = result
			return nil
		})

		if err != nil {
			errRet = err
			return
		}

		if response == nil || len(response.Response.Items) < 1 {
			break
		}

		orgShareUnitMembers = append(orgShareUnitMembers, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DeleteOrganizationOrgShareUnitMemberById(ctx context.Context, unitId, area, shareMemberUins string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteShareUnitMembersRequest()
	request.UnitId = &unitId
	request.Area = &area
	split := strings.Split(shareMemberUins, tccommon.COMMA_SP)
	if len(split) < 1 {
		errRet = fmt.Errorf("shareMemberUins is broken, %s", shareMemberUins)
		return
	}
	for _, v := range split {
		request.Members = append(request.Members, &organization.ShareMember{ShareMemberUin: helper.StrToInt64Point(v)})
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeleteShareUnitMembers(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DeleteOrganizationOrgShareUnitMemberV2ById(ctx context.Context, unitId, area string, orgShareUnitMembers []*organization.ShareUnitMember) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteShareUnitMembersRequest()
	request.UnitId = &unitId
	request.Area = &area

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for i := 0; i < len(orgShareUnitMembers); i += batchSize {
		end := i + batchSize
		if end > len(orgShareUnitMembers) {
			end = len(orgShareUnitMembers)
		}

		batch := orgShareUnitMembers[i:end]
		// clear Members value
		request.Members = nil
		for _, item := range batch {
			shareMember := organization.ShareMember{}
			shareMember.ShareMemberUin = item.ShareMemberUin
			request.Members = append(request.Members, &shareMember)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().DeleteShareUnitMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			errRet = err
			return
		}
	}

	return
}

func (me *OrganizationService) AddOrganizationOrgShareUnitMemberV2ById(ctx context.Context, unitId, area string, orgShareUnitMembers []*organization.ShareUnitMember) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewAddShareUnitMembersRequest()
	request.UnitId = &unitId
	request.Area = &area

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for i := 0; i < len(orgShareUnitMembers); i += batchSize {
		end := i + batchSize
		if end > len(orgShareUnitMembers) {
			end = len(orgShareUnitMembers)
		}

		batch := orgShareUnitMembers[i:end]
		// clear Members value
		request.Members = nil
		for _, item := range batch {
			shareMember := organization.ShareMember{}
			shareMember.ShareMemberUin = item.ShareMemberUin
			request.Members = append(request.Members, &shareMember)
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().AddShareUnitMembers(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			errRet = err
			return
		}
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareAreaByFilter(ctx context.Context, param map[string]interface{}) (orgShareArea []*organization.ShareArea, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewDescribeShareAreasRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Lang" {
			request.Lang = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DescribeShareAreas(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.Items) < 1 {
		return
	}

	orgShareArea = response.Response.Items
	return
}

func (me *OrganizationService) DescribeOrganizationOrgManagePolicyConfigById(ctx context.Context, organizationId string, policyType string) (OrgManagePolicyConfig *organization.DescribePolicyConfigResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribePolicyConfigRequest()
	request.OrganizationId = helper.StrToUint64Point(organizationId)
	request.Type = helper.IntUint64(ServiceControlPolicyCode)

	if policyType == TagPolicyType {
		request.Type = helper.IntUint64(TagPolicyCode)
	}
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DescribePolicyConfig(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if *response.Response.Status == 1 {
		OrgManagePolicyConfig = response.Response
	}
	return
}

func (me *OrganizationService) DeleteOrganizationOrgManagePolicyConfigById(ctx context.Context, organizationId string, policyType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDisablePolicyTypeRequest()
	request.OrganizationId = helper.StrToUint64Point(organizationId)
	request.PolicyType = &policyType

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DisablePolicyType(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgManagePolicyById(ctx context.Context, policyId, policyType string) (OrgManagePolicy *organization.DescribePolicyResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListPoliciesRequest()
	request.PolicyType = helper.String(policyType)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	pageStart := uint64(1)
	rp := uint64(PAGE_ITEM) //to save in extension
	result := make([]*organization.ListPolicyNode, 0)
	for {
		request.Page = &pageStart
		request.Rp = &rp
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().ListPolicies(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.List) < 1 {
			break
		}
		result = append(result, response.Response.List...)
		if len(response.Response.List) < PAGE_ITEM {
			break
		}
		pageStart += 1
	}

	for _, item := range result {
		if helper.UInt64ToStr(*item.PolicyId) == policyId {
			requestDescribe := organization.NewDescribePolicyRequest()
			requestDescribe.PolicyId = item.PolicyId
			requestDescribe.PolicyType = helper.String(policyType)
			responseDescribe, err := me.client.UseOrganizationClient().DescribePolicy(requestDescribe)
			if err != nil {
				errRet = err
				return
			}
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), responseDescribe.ToJsonString())

			if responseDescribe == nil || responseDescribe.Response == nil {
				break
			}
			OrgManagePolicy = responseDescribe.Response
		}
	}
	return
}

func (me *OrganizationService) DeleteOrganizationOrgManagePolicyById(ctx context.Context, policyId, policyType string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeletePolicyRequest()
	request.PolicyId = helper.StrToUint64Point(policyId)
	request.Type = helper.String(policyType)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DeletePolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationOrgManagePolicyTargetById(ctx context.Context, policyType string, policyId string, targetType string, targetId string) (OrgManagePolicyTarget *organization.ListTargetsForPolicyNode, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListTargetsForPolicyRequest()
	request.PolicyType = &policyType
	request.PolicyId = helper.StrToUint64Point(policyId)
	switch targetType {
	case TargetTypeNode:
		request.TargetType = helper.String(DescribeTargetTypeNode)
	case TargetTypeMember:
		request.TargetType = helper.String(DescribeTargetTypeMember)
	}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().ListTargetsForPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	for _, item := range response.Response.List {
		if item.Uin != nil && helper.UInt64ToStr(*item.Uin) == targetId {
			OrgManagePolicyTarget = item
		}
	}
	return
}

func (me *OrganizationService) DeleteOrganizationOrgManagePolicyTargetById(ctx context.Context, policyType string, policyId string, targetType string, targetId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDetachPolicyRequest()
	request.Type = &policyType
	request.PolicyId = helper.StrToUint64Point(policyId)
	request.TargetType = &targetType
	request.TargetId = helper.StrToUint64Point(targetId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().DetachPolicy(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *OrganizationService) DescribeOrganizationServiceAssignMemberById(ctx context.Context, serviceId string) (items []*organization.OrganizationServiceAssignMember, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListOrgServiceAssignMemberRequest()
	serviceIdInt, _ := strconv.ParseUint(serviceId, 10, 64)
	request.ServiceId = &serviceIdInt

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var offset uint64 = 0
	var pageSize uint64 = 10
	items = make([]*organization.OrganizationServiceAssignMember, 0)

	for {
		request.Offset = &offset
		request.Limit = &pageSize
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().ListOrgServiceAssignMember(request)
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

		items = append(items, response.Response.Items...)
		if len(response.Response.Items) < int(pageSize) {
			break
		}

		offset += pageSize
	}

	return
}

func (me *OrganizationService) DeleteOrganizationServiceAssignMemberById(ctx context.Context, serviceId string, memberUinList []*int64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDeleteOrgServiceAssignRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	serviceIdInt, _ := strconv.ParseUint(serviceId, 10, 64)
	for _, memberUin := range memberUinList {
		ratelimit.Check(request.GetAction())
		request.ServiceId = &serviceIdInt
		request.MemberUin = memberUin
		response, err := me.client.UseOrganizationClient().DeleteOrgServiceAssign(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
	}

	return
}

func (me *OrganizationService) DescribeOrganizationServicesByFilter(ctx context.Context, param map[string]interface{}) (members []*organization.OrganizationServiceAssign, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewListOrganizationServiceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 10
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().ListOrganizationService(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}

		members = append(members, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeIdentityCenterUserById(ctx context.Context, zoneId string, userId string) (ret *organization.UserInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetUserRequest()
	request.UserId = helper.String(userId)
	request.ZoneId = helper.String(zoneId)
	response := organization.NewGetUserResponse()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOrganizationClient().GetUser(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		errRet = err
		log.Printf("[CRITAL]%s update identity center user failed, reason:%+v", logId, err)
		return
	}

	if response.Response == nil {
		return
	}

	ret = response.Response.UserInfo
	return
}

func (me *OrganizationService) DescribeIdentityCenterGroupById(ctx context.Context, zoneId string, groupId string) (ret *organization.GroupInfo, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetGroupRequest()
	request.ZoneId = helper.String(zoneId)
	request.GroupId = helper.String(groupId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.GroupInfo
	return
}

func (me *OrganizationService) DescribeIdentityCenterUserGroupAttachmentById(ctx context.Context, zoneId, groupId, userId string) (joinedGroup *organization.JoinedGroups, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListJoinedGroupsForUserRequest()
	request.ZoneId = helper.String(zoneId)
	request.UserId = helper.String(userId)
	request.MaxResults = helper.Int64(100)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().ListJoinedGroupsForUser(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		for _, v := range response.Response.JoinedGroups {
			if *v.GroupId == groupId {
				joinedGroup = v
				return
			}
		}
		if len(response.Response.JoinedGroups) < int(*request.MaxResults) {
			break
		} else {
			request.NextToken = response.Response.NextToken
		}
	}
	return
}

func (me *OrganizationService) DescribeIdentityCenterExternalSamlIdentityProviderById(ctx context.Context, zoneId string) (ret *organization.SAMLServiceProvider, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetZoneSAMLServiceProviderInfoRequest()
	request.ZoneId = helper.String(zoneId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetZoneSAMLServiceProviderInfo(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.SAMLServiceProvider
	return
}

func (me *OrganizationService) DescribeIdentityCenterExternalSamlIdentityProviderById1(ctx context.Context, zoneId string) (ret *organization.SAMLIdentityProviderConfiguration, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetExternalSAMLIdentityProviderRequest()
	request.ZoneId = helper.String(zoneId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetExternalSAMLIdentityProvider(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.SAMLIdentityProviderConfiguration
	return
}

func (me *OrganizationService) DescribeIdentityCenterRoleConfigurationById(ctx context.Context, zoneId string, roleConfigurationId string) (ret *organization.RoleConfiguration, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetRoleConfigurationRequest()
	request.ZoneId = helper.String(zoneId)
	request.RoleConfigurationId = helper.String(roleConfigurationId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetRoleConfiguration(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.RoleConfigurationInfo
	return
}

func (me *OrganizationService) DescribeIdentityCenterRoleConfigurationPermissionPolicyAttachmentById(ctx context.Context, zoneId, roleConfigurationId, rolePolicyType string) (ret *organization.ListPermissionPoliciesInRoleConfigurationResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListPermissionPoliciesInRoleConfigurationRequest()
	request.ZoneId = helper.String(zoneId)
	request.RoleConfigurationId = helper.String(roleConfigurationId)
	request.RolePolicyType = helper.String(rolePolicyType)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().ListPermissionPoliciesInRoleConfiguration(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *OrganizationService) DescribeIdentityCenterUserSyncProvisioningById(ctx context.Context, zoneId, userProvisioningId string) (ret *organization.UserProvisioning, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetUserSyncProvisioningRequest()
	request.ZoneId = helper.String(zoneId)
	request.UserProvisioningId = helper.String(userProvisioningId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetUserSyncProvisioning(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil {
		return
	}

	ret = response.Response.UserProvisioning
	return
}

func (me *OrganizationService) DescribeIdentityCenterRoleAssignmentById(ctx context.Context, roleAssignmentId string) (ret *organization.ListRoleAssignmentsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	idSplit := strings.Split(roleAssignmentId, tccommon.FILED_SP)
	if len(idSplit) != 6 {
		errRet = fmt.Errorf("roleAssignmentId is broken,%s", roleAssignmentId)
		return
	}

	zoneId := idSplit[0]
	roleConfigurationId := idSplit[1]
	targetType := idSplit[2]
	targetUinString := idSplit[3]
	principalType := idSplit[4]
	principalId := idSplit[5]

	request := organization.NewListRoleAssignmentsRequest()
	request.ZoneId = helper.String(zoneId)
	request.RoleConfigurationId = helper.String(roleConfigurationId)
	request.TargetType = helper.String(targetType)
	targetUin, err := strconv.ParseInt(targetUinString, 10, 64)
	if err != nil {
		errRet = err
		return
	}
	request.TargetUin = helper.Int64(targetUin)
	request.PrincipalType = helper.String(principalType)
	request.PrincipalId = helper.String(principalId)
	request.MaxResults = helper.Int64(10)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().ListRoleAssignments(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *OrganizationService) GetAssignmentTaskStatus(ctx context.Context, zoneId, taskId string) (taskStatus *organization.TaskStatus, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetTaskStatusRequest()
	request.ZoneId = helper.String(zoneId)
	request.TaskId = helper.String(taskId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, e := me.client.UseOrganizationClient().GetTaskStatus(request)
	if e != nil {
		errRet = e
		return
	}

	if response.Response != nil {
		taskStatus = response.Response.TaskStatus
	}
	return

}
func (me *OrganizationService) AssignmentTaskStatusStateRefreshFunc(zoneId, taskId string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		ctx := tccommon.ContextNil

		var object *organization.TaskStatus
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			result, e := me.GetAssignmentTaskStatus(ctx, zoneId, taskId)
			if e != nil {
				return tccommon.RetryError(e)
			}
			object = result
			return nil
		})
		if err != nil {
			return nil, "", err
		}

		return object, *object.Status, nil
	}
}

func (me *OrganizationService) UpdateOrganizationRootNodeName(ctx context.Context, orgId uint64, name string) (errRet error) {
	logId := tccommon.GetLogId(ctx)
	organizationResponseParams, err := me.DescribeOrganizationOrganizationById(ctx)
	if err != nil {
		return err
	}
	if organizationResponseParams == nil {
		return fmt.Errorf("organization is nil")
	}
	rootNodeId := organizationResponseParams.RootNodeId

	request := organization.NewUpdateOrganizationNodeRequest()
	request.NodeId = helper.Int64Uint64(*rootNodeId)
	request.Name = &name

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseOrganizationClient().UpdateOrganizationNode(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s update organization orgNode name failed, reason:%+v", logId, err)
		return err
	}

	return nil
}

func (me *OrganizationService) DescribeIdentityCenterUsersByFilter(ctx context.Context, param map[string]interface{}) (users []*organization.UserInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewListUsersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}
		if k == "UserStatus" {
			request.UserStatus = v.(*string)
		}
		if k == "UserType" {
			request.UserType = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.(*string)
		}
		if k == "FilterGroups" {
			request.FilterGroups = v.([]*string)
		}
		if k == "SortField" {
			request.SortField = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
	}

	users = make([]*organization.UserInfo, 0)
	for {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseOrganizationClient().ListUsers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		users = append(users, response.Response.Users...)

		if response.Response.IsTruncated != nil {
			if *response.Response.IsTruncated {
				request.NextToken = response.Response.NextToken
			} else {
				break
			}
		} else {
			errRet = fmt.Errorf("ListUsers IsTruncated is nil")
			return
		}
	}

	return
}

func (me *OrganizationService) DescribeIdentityCenterGroupsByFilter(ctx context.Context, param map[string]interface{}) (groups []*organization.GroupInfo, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewListGroupsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.(*string)
		}
		if k == "GroupType" {
			request.GroupType = v.(*string)
		}
		if k == "FilterUsers" {
			request.FilterUsers = v.([]*string)
		}
		if k == "SortField" {
			request.SortField = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
	}

	groups = make([]*organization.GroupInfo, 0)
	for {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseOrganizationClient().ListGroups(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		groups = append(groups, response.Response.Groups...)

		if response.Response.IsTruncated != nil {
			if *response.Response.IsTruncated {
				request.NextToken = response.Response.NextToken
			} else {
				break
			}
		} else {
			errRet = fmt.Errorf("ListGroups IsTruncated is nil")
			return
		}
	}

	return
}

func (me *OrganizationService) DescribeIdentityCenterRoleConfigurationsByFilter(ctx context.Context, param map[string]interface{}) (roleConfigurations []*organization.RoleConfiguration, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewListRoleConfigurationsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.(*string)
		}
		if k == "FilterTargets" {
			request.FilterTargets = v.([]*int64)
		}
		if k == "PrincipalId" {
			request.PrincipalId = v.(*string)
		}
	}

	roleConfigurations = make([]*organization.RoleConfiguration, 0)
	for {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseOrganizationClient().ListRoleConfigurations(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		roleConfigurations = append(roleConfigurations, response.Response.RoleConfigurations...)

		if response.Response.IsTruncated != nil {
			if *response.Response.IsTruncated {
				request.NextToken = response.Response.NextToken
			} else {
				break
			}
		} else {
			errRet = fmt.Errorf("ListRoleConfigurations IsTruncated is nil")
			return
		}
	}

	return
}

func (me *OrganizationService) DescribeOrganizationNodesByFilter(ctx context.Context, param map[string]interface{}) (nodes []*organization.OrgNode, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewDescribeOrganizationNodesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Tags" {
			request.Tags = v.([]*organization.Tag)
		}
	}

	var (
		limit  int64 = 50
		offset int64 = 0
	)
	request.Limit = &limit
	request.Offset = &offset
	nodes = make([]*organization.OrgNode, 0)

	for {
		ratelimit.Check(request.GetAction())
		response, err := me.client.UseOrganizationClient().DescribeOrganizationNodes(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		nodes = append(nodes, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}
	}
	return
}

func (me *OrganizationService) DescribeIdentityCenterScimCredentialById(ctx context.Context, zoneId string, credentialId string) (ret *organization.SCIMCredential, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListSCIMCredentialsRequest()
	request.ZoneId = helper.String(zoneId)
	request.CredentialId = helper.String(credentialId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().ListSCIMCredentials(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SCIMCredentials) < 1 {
		return
	}

	ret = response.Response.SCIMCredentials[0]
	return
}

func (me *OrganizationService) DescribeIdentityCenterScimCredentialStatusById(ctx context.Context, zoneId string, credentialId string) (ret *organization.SCIMCredential, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListSCIMCredentialsRequest()
	request.ZoneId = helper.String(zoneId)
	request.CredentialId = helper.String(credentialId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().ListSCIMCredentials(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.SCIMCredentials) < 1 {
		return
	}

	ret = response.Response.SCIMCredentials[0]
	return
}

func (me *OrganizationService) DescribeIdentityCenterScimSynchronizationStatusById(ctx context.Context, zoneId string) (ret *organization.GetSCIMSynchronizationStatusResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetSCIMSynchronizationStatusRequest()
	request.ZoneId = helper.String(zoneId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseOrganizationClient().GetSCIMSynchronizationStatus(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitResourceById(ctx context.Context, unitId string, area string, shareResourceType string, productResourceId string) (ret *organization.ShareUnitResource, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribeShareUnitResourcesRequest()
	request.UnitId = helper.String(unitId)
	request.Area = helper.String(area)
	request.SearchKey = helper.String(productResourceId)
	request.Type = helper.String(shareResourceType)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		limit  uint64 = 50
		offset uint64 = 0
	)
	request.Limit = &limit
	request.Offset = &offset

	for {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseOrganizationClient().DescribeShareUnitResources(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || response.Response == nil {
			return
		}

		for _, item := range response.Response.Items {
			if *item.ProductResourceId == productResourceId {
				ret = item
				return
			}
		}
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit

	}
	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitResourcesByFilter(ctx context.Context, param map[string]interface{}) (ret []*organization.ShareUnitResource, errRet error) {
	var (
		logId   = common.GetLogId(ctx)
		request = organization.NewDescribeShareUnitResourcesRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "UnitId" {
			request.UnitId = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
		if k == "Type" {
			request.Type = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 50
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeShareUnitResources(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitsByFilter(ctx context.Context, param map[string]interface{}) (ret []*organization.ManagerShareUnit, errRet error) {
	var (
		logId   = common.GetLogId(ctx)
		request = organization.NewDescribeShareUnitsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeShareUnits(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationOrgShareUnitMembersByFilter(ctx context.Context, param map[string]interface{}) (ret []*organization.ShareUnitMember, errRet error) {
	var (
		logId   = common.GetLogId(ctx)
		request = organization.NewDescribeShareUnitMembersRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "UnitId" {
			request.UnitId = v.(*string)
		}
		if k == "Area" {
			request.Area = v.(*string)
		}
		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 50
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseOrganizationClient().DescribeShareUnitMembers(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Items) < 1 {
			break
		}
		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeRoleConfigurationProvisioningsByFilter(ctx context.Context, param map[string]interface{}) (roleConfigurationProvisionings []*organization.RoleConfigurationProvisionings, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = organization.NewListRoleConfigurationProvisioningsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}
		if k == "RoleConfigurationId" {
			request.RoleConfigurationId = v.(*string)
		}
		if k == "TargetType" {
			request.TargetType = v.(*string)
		}
		if k == "TargetUin" {
			request.TargetUin = v.(*int64)
		}
		if k == "DeploymentStatus" {
			request.DeploymentStatus = v.(*string)
		}
		if k == "Filter" {
			request.Filter = v.(*string)
		}
	}

	request.MaxResults = helper.IntInt64(100)
	for {
		ratelimit.Check(request.GetAction())

		response, err := me.client.UseOrganizationClient().ListRoleConfigurationProvisionings(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		roleConfigurationProvisionings = append(roleConfigurationProvisionings, response.Response.RoleConfigurationProvisionings...)

		if response.Response.IsTruncated != nil && *response.Response.IsTruncated {
			request.NextToken = response.Response.NextToken
		} else {
			break
		}
	}

	return
}

func (me *OrganizationService) DescribeOrganizationResourceToShareMemberByFilter(ctx context.Context, param map[string]interface{}) (ret []*organization.ShareResourceToMember, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = organization.NewDescribeResourceToShareMemberRequest()
		response = organization.NewDescribeResourceToShareMemberResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Area" {
			request.Area = v.(*string)
		}

		if k == "SearchKey" {
			request.SearchKey = v.(*string)
		}

		if k == "Type" {
			request.Type = v.(*string)
		}

		if k == "ProductResourceIds" {
			request.ProductResourceIds = v.([]*string)
		}
	}

	var (
		offset uint64 = 0
		limit  uint64 = 50
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().DescribeResourceToShareMember(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe resource to share member failed, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			return
		}

		if len(response.Response.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationMembersAuthPolicyAttachmentById(ctx context.Context, policyId, orgSubAccountUin string) (ret []*organization.OrgMembersAuthPolicy, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewDescribeOrganizationMembersAuthPolicyRequest()
	response := organization.NewDescribeOrganizationMembersAuthPolicyResponse()
	request.PolicyId = helper.StrToInt64Point(policyId)
	request.OrgSubAccountUin = helper.StrToInt64Point(orgSubAccountUin)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	var (
		offset int64 = 0
		limit  int64 = 50
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())
			result, e := me.client.UseOrganizationClient().DescribeOrganizationMembersAuthPolicy(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe organization members auth policy, Response is nil."))
			}

			response = result
			return nil
		})

		if errRet != nil {
			return
		}

		if len(response.Response.Items) < 1 {
			break
		}

		ret = append(ret, response.Response.Items...)
		if len(response.Response.Items) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *OrganizationService) DescribeOrganizationExternalSamlIdpCertificateById(ctx context.Context, zoneId, certificateId string) (ret *organization.SAMLIdPCertificate, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewListExternalSAMLIdPCertificatesRequest()
	response := organization.NewListExternalSAMLIdPCertificatesResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOrganizationClient().ListExternalSAMLIdPCertificates(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.SAMLIdPCertificates == nil || len(result.Response.SAMLIdPCertificates) == 0 {
			return resource.NonRetryableError(fmt.Errorf("List external saml idp certificate failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	for _, item := range response.Response.SAMLIdPCertificates {
		if item.CertificateId != nil && *item.CertificateId == certificateId {
			ret = item
			break
		}
	}

	return
}

func (me *OrganizationService) DescribeOrganizationExternalSamlIdentityProviderById(ctx context.Context, zoneId string) (ret *organization.SAMLIdentityProviderConfiguration, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := organization.NewGetExternalSAMLIdentityProviderRequest()
	response := organization.NewGetExternalSAMLIdentityProviderResponse()
	request.ZoneId = &zoneId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	errRet = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseOrganizationClient().GetExternalSAMLIdentityProvider(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Get external saml identity provider failed, Response is nil."))
		}

		response = result
		return nil
	})

	if errRet != nil {
		return
	}

	ret = response.Response.SAMLIdentityProviderConfiguration
	return
}
