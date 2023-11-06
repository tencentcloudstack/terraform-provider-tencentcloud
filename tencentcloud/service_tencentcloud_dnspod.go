package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

// basic information

type DnspodService struct {
	client *connectivity.TencentCloudClient
}

// ////////api
func (me *DnspodService) ModifyDnsPodDomainStatus(ctx context.Context, domain string, status string) (errRet error) {
	logId := getLogId(ctx)
	request := dnspod.NewModifyDomainStatusRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.Domain = helper.String(domain)
	request.Status = &status

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseDnsPodClient().ModifyDomainStatus(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify dnspod domain status failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *DnspodService) ModifyDnsPodDomainRemark(ctx context.Context, domain string, remark string) (errRet error) {
	logId := getLogId(ctx)
	request := dnspod.NewModifyDomainRemarkRequest()
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()
	request.Domain = helper.String(domain)
	request.Remark = &remark

	if err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		_, err := me.client.UseDnsPodClient().ModifyDomainRemark(request)
		if err != nil {
			return retryError(err, InternalError)
		}
		return nil
	}); err != nil {
		log.Printf("[CRITAL]%s modify dnspod domain remark failed, reason: %v", logId, err)
		return err
	}
	return
}

func (me *DnspodService) DescribeDomain(ctx context.Context, domain string) (ret *dnspod.DescribeDomainResponse, errRet error) {

	logId := getLogId(ctx)
	request := dnspod.NewDescribeDomainRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	request.Domain = helper.String(domain)

	response, err := me.client.UseDnsPodClient().DescribeDomain(request)

	if err != nil {
		errRet = err
		return
	}

	if response == nil || response.Response == nil {
		return nil, fmt.Errorf("TencentCloud SDK return nil response, %s", request.GetAction())
	}

	return response, nil
}

func (me *DnspodService) DeleteDomain(ctx context.Context, domain string) (errRet error) {

	logId := getLogId(ctx)
	request := dnspod.NewDeleteDomainRequest()
	ratelimit.Check(request.GetAction())
	request.Domain = helper.String(domain)

	response, err := me.client.UseDnsPodClient().DeleteDomain(request)

	defer func() {
		if errRet != nil {
			responseStr := ""
			if response != nil {
				responseStr = response.ToJsonString()
			}
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s],response body [%s], reason[%s]\n",
				logId,
				request.GetAction(),
				request.ToJsonString(),
				responseStr,
				errRet.Error())
		}
	}()

	errRet = err
	return
}

func (me *DnspodService) DescribeRecordList(ctx context.Context, request *dnspod.DescribeRecordListRequest) (list []*dnspod.RecordListItem, info *dnspod.RecordCountInfo, errRet error) {
	logId := getLogId(ctx)
	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
				logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDnsPodClient().DescribeRecordList(request)

	if err != nil {
		errRet = err
		return
	}

	list = response.Response.RecordList
	info = response.Response.RecordCountInfo

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
		logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DnspodService) DescribeDnspodDomainListByFilter(ctx context.Context, param map[string]interface{}) (domain_list []*dnspod.DomainListItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeDomainFilterListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Type" {
			request.Type = v.(*string)
		}
		if k == "GroupId" {
			request.GroupId = v.([]*int64)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "SortField" {
			request.SortField = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
		if k == "Status" {
			request.Status = v.([]*string)
		}
		if k == "Package" {
			request.Package = v.([]*string)
		}
		if k == "Remark" {
			request.Remark = v.(*string)
		}
		if k == "UpdatedAtBegin" {
			request.UpdatedAtBegin = v.(*string)
		}
		if k == "UpdatedAtEnd" {
			request.UpdatedAtEnd = v.(*string)
		}
		if k == "RecordCountBegin" {
			request.RecordCountBegin = v.(*uint64)
		}
		if k == "RecordCountEnd" {
			request.RecordCountEnd = v.(*uint64)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
		}
		if k == "Tags" {
			request.Tags = v.([]*dnspod.TagItemFilter)
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
		response, err := me.client.UseDnsPodClient().DescribeDomainFilterList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.DomainList) < 1 {
			break
		}
		domain_list = append(domain_list, response.Response.DomainList...)
		if len(response.Response.DomainList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DnspodService) DescribeDnspodDomainAnalyticsByFilter(ctx context.Context, param map[string]interface{}) (alias_data []*dnspod.DomainAliasAnalyticsItem, data []*dnspod.DomainAnalyticsDetail, info *dnspod.DomainAnalyticsInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeDomainAnalyticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "DnsFormat" {
			request.DnsFormat = v.(*string)
		}
		if k == "DomainId" {
			request.DomainId = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDnsPodClient().DescribeDomainAnalytics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	alias_data = response.Response.AliasData
	data = response.Response.Data
	info = response.Response.Info

	return
}

func (me *DnspodService) DescribeDnspodDomainLogListByFilter(ctx context.Context, param map[string]interface{}) (domain_log_list []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeDomainLogListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "DomainId" {
			request.DomainId = v.(*uint64)
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
		response, err := me.client.UseDnsPodClient().DescribeDomainLogList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.LogList) < 1 {
			break
		}
		domain_log_list = append(domain_log_list, response.Response.LogList...)
		if len(response.Response.LogList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DnspodService) DescribeDnspodRecordAnalyticsByFilter(ctx context.Context, param map[string]interface{}) (alias_data []*dnspod.SubdomainAliasAnalyticsItem, data []*dnspod.DomainAnalyticsDetail, info *dnspod.SubdomainAnalyticsInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeSubdomainAnalyticsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "StartDate" {
			request.StartDate = v.(*string)
		}
		if k == "EndDate" {
			request.EndDate = v.(*string)
		}
		if k == "Subdomain" {
			request.Subdomain = v.(*string)
		}
		if k == "DnsFormat" {
			request.DnsFormat = v.(*string)
		}
		if k == "DomainId" {
			request.DomainId = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseDnsPodClient().DescribeSubdomainAnalytics(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	alias_data = response.Response.AliasData
	data = response.Response.Data
	info = response.Response.Info

	return
}

func (me *DnspodService) DescribeDnspodRecordLineListByFilter(ctx context.Context, param map[string]interface{}) (line_list []*dnspod.LineInfo, line_group_list []*dnspod.LineGroupInfo, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeRecordLineListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "DomainGrade" {
			request.DomainGrade = v.(*string)
		}
		if k == "DomainId" {
			request.DomainId = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DescribeRecordLineList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	line_list = response.Response.LineList
	line_group_list = response.Response.LineGroupList

	return
}

func (me *DnspodService) DescribeDnspodRecordListByFilter(ctx context.Context, param map[string]interface{}) (record_list []*dnspod.RecordListItem, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeRecordFilterListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Domain" {
			request.Domain = v.(*string)
		}
		if k == "DomainId" {
			request.DomainId = v.(*uint64)
		}
		if k == "SubDomain" {
			request.SubDomain = v.(*string)
		}
		if k == "RecordType" {
			request.RecordType = v.([]*string)
		}
		if k == "RecordLine" {
			request.RecordLine = v.([]*string)
		}
		if k == "GroupId" {
			request.GroupId = v.([]*uint64)
		}
		if k == "Keyword" {
			request.Keyword = v.(*string)
		}
		if k == "SortField" {
			request.SortField = v.(*string)
		}
		if k == "SortType" {
			request.SortType = v.(*string)
		}
		if k == "RecordValue" {
			request.RecordValue = v.(*string)
		}
		if k == "RecordStatus" {
			request.RecordStatus = v.([]*string)
		}
		if k == "WeightBegin" {
			request.WeightBegin = v.(*uint64)
		}
		if k == "WeightEnd" {
			request.WeightEnd = v.(*uint64)
		}
		if k == "MXBegin" {
			request.MXBegin = v.(*uint64)
		}
		if k == "MXEnd" {
			request.MXEnd = v.(*uint64)
		}
		if k == "TTLBegin" {
			request.TTLBegin = v.(*uint64)
		}
		if k == "TTLEnd" {
			request.TTLEnd = v.(*uint64)
		}
		if k == "UpdatedAtBegin" {
			request.UpdatedAtBegin = v.(*string)
		}
		if k == "UpdatedAtEnd" {
			request.UpdatedAtEnd = v.(*string)
		}
		if k == "Remark" {
			request.Remark = v.(*string)
		}
		if k == "IsExactSubDomain" {
			request.IsExactSubDomain = v.(*bool)
		}
		if k == "ProjectId" {
			request.ProjectId = v.(*int64)
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
		response, err := me.client.UseDnsPodClient().DescribeRecordFilterList(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.RecordList) < 1 {
			break
		}
		record_list = append(record_list, response.Response.RecordList...)
		if len(response.Response.RecordList) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *DnspodService) DescribeDnspodRecordTypeByFilter(ctx context.Context, param map[string]interface{}) (type_list []*string, errRet error) {
	var (
		logId   = getLogId(ctx)
		request = dnspod.NewDescribeRecordTypeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "DomainGrade" {
			request.DomainGrade = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DescribeRecordType(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || len(response.Response.TypeList) < 1 {
		return
	}
	type_list = append(type_list, response.Response.TypeList...)

	return
}
func (me *DnspodService) DescribeDnspodRecordGroupById(ctx context.Context, domain string, groupId uint64) (recordGroup *dnspod.RecordGroupInfo, errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDescribeRecordGroupListRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DescribeRecordGroupList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.GroupList) < 1 {
		return
	}

	for _, item := range response.Response.GroupList {
		if *item.GroupId == groupId {
			recordGroup = item
			return
		}
	}
	return
}

func (me *DnspodService) DeleteDnspodRecordGroupById(ctx context.Context, domain string, groupId uint64) (errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDeleteRecordGroupRequest()
	request.Domain = &domain
	request.GroupId = &groupId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DeleteRecordGroup(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DnspodService) DescribeDnspodDomainAliasById(ctx context.Context, domain string, domainAliasId int64) (domainAliasInfo *dnspod.DomainAliasInfo, errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDescribeDomainAliasListRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DescribeDomainAliasList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if len(response.Response.DomainAliasList) < 1 {
		return
	}

	for _, item := range response.Response.DomainAliasList {
		if *item.Id == domainAliasId {
			domainAliasInfo = item
			return
		}
	}

	return
}

func (me *DnspodService) DeleteDnspodDomainAliasById(ctx context.Context, domain string, domainAliasId int64) (errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDeleteDomainAliasRequest()
	request.Domain = &domain
	request.DomainAliasId = &domainAliasId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DeleteDomainAlias(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *DnspodService) DescribeDnspodCustomLineById(ctx context.Context, domain string, name string) (customLineInfo *dnspod.CustomLineInfo, errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDescribeDomainCustomLineListRequest()
	request.Domain = &domain

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DescribeDomainCustomLineList(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response.Response == nil || response.Response.LineList == nil {
		return
	}

	for _, item := range response.Response.LineList {
		if *item.Name == name {
			customLineInfo = item
			return
		}
	}

	return
}

func (me *DnspodService) DeleteDnspodCustomLineById(ctx context.Context, domain string, name string) (errRet error) {
	logId := getLogId(ctx)

	request := dnspod.NewDeleteDomainCustomLineRequest()
	request.Domain = &domain
	request.Name = &name

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseDnsPodClient().DeleteDomainCustomLine(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
