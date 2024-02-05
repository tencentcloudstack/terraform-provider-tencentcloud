package common

import (
	"encoding/json"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"log"
	"strconv"
	"strings"
	"time"

	cam "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

const (
	// DefaultSearchLogStartTimestamp sync logs start time 2023-11-07 16:41:00
	DefaultSearchLogStartTimestamp = 1699346460000

	DefaultTopicId = "aef50d54-b17d-4782-8618-a7873203ec29"

	QueryGrammarRule = " AND "
)

// ResourceAccountInfo 资源账户信息
type ResourceAccountInfo struct {
	ResourceType string // 资源类型
	ResourceName string // 资源名称
	AccountId    string // 主账号ID
	PrincipalId  string // 用户ID
	UserName     string // 用户名
}

// GetResourceCreatorAccountInfo get resource creator user info
func GetResourceCreatorAccountInfo(client *connectivity.TencentCloudClient, resourceCreateAction string, resources []*ResourceInstance) map[string]*ResourceAccountInfo {
	resourceIdToSubAccountInfoMap := make(map[string]*ResourceAccountInfo)
	if resourceCreateAction == "" {
		return resourceIdToSubAccountInfoMap
	}

	request := cls.NewSearchLogRequest()
	request.From = helper.IntInt64(DefaultSearchLogStartTimestamp)
	request.To = helper.Int64(CurrentTimeMillisecond())
	request.TopicId = helper.String(DefaultTopicId)

	for _, r := range resources {
		query := resourceCreateAction + QueryGrammarRule
		if r.Id != "" {
			query = query + r.Id
		} else if r.Name != "" {
			query = query + r.Name
		} else {
			continue
		}
		request.Query = helper.String(query)

		response, err := client.UseClsClient().SearchLog(request)
		if err != nil {
			log.Printf("[CRITAL] search resource[%v] log data error: %v", r.Id, err.Error())
			return resourceIdToSubAccountInfoMap
		}
		if response == nil || response.Response == nil {
			log.Printf("[CRITAL] search resource[%v] log data response is nil", r.Id)
			return resourceIdToSubAccountInfoMap
		}
		if len(response.Response.Results) == 0 {
			log.Printf("[CRITAL] search resource[%v] log data response results is empty", r.Id)
			return resourceIdToSubAccountInfoMap
		}

		result := response.Response.Results[0]
		if result != nil {
			var jsonData string
			if len(*result.LogJson) > 2 {
				jsonData = *result.LogJson
			} else if len(*result.RawLog) > 2 {
				jsonData = *result.RawLog
			} else {
				continue
			}

			resourceAccountInfo := ParseLogJsonData(jsonData)
			if resourceAccountInfo.PrincipalId == resourceAccountInfo.UserName &&
				resourceAccountInfo.PrincipalId != resourceAccountInfo.AccountId {
				userName := GetSubAccountUserName(client, resourceAccountInfo.PrincipalId)
				resourceAccountInfo.UserName = userName
			}
			resourceIdToSubAccountInfoMap[r.Id] = resourceAccountInfo
		}
	}

	return resourceIdToSubAccountInfoMap
}

// GetSubAccountUserName get sub account user name
func GetSubAccountUserName(client *connectivity.TencentCloudClient, uin string) string {
	uinNum, err := strconv.ParseUint(uin, 10, 64)
	if err != nil {
		log.Printf("[CRITAL] parse uin[%v] to uint64 type error: %v", uin, err.Error())
		return ""
	}

	request := cam.NewDescribeSubAccountsRequest()

	uinArray := []*uint64{helper.Uint64(uinNum)}
	request.FilterSubAccountUin = uinArray

	response, err := client.UseCamClient().DescribeSubAccounts(request)
	if err != nil {
		log.Printf("[CRITAL] get sub account[%v] data error: %v", uin, err.Error())
		return ""
	}
	if response == nil || response.Response == nil {
		log.Printf("[CRITAL] get sub account[%v] data response is nil", uin)
		return ""
	}

	name := response.Response.SubAccounts[0].Name
	return *name
}

// CurrentTimeMillisecond get the current millisecond timestamp
func CurrentTimeMillisecond() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func ParseLogJsonData(jsonData string) *ResourceAccountInfo {
	if jsonData == "" {
		return nil
	}

	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		log.Printf("[CRITAL] parse log json data[%v] error: %v", jsonData, err.Error())
		return nil
	}

	resourceType := ""
	if v, ok := data["resourceType"]; ok {
		resourceType = v.(string)
	}
	resourceName := ""
	if v, ok := data["resourceName"]; ok {
		resourceName = v.(string)
		if resourceName != "" {
			resourceName = strings.Split(resourceName, "/")[0]
		}
	}
	accountId, principalId, userName := parseUserIdentityFields(data)

	return &ResourceAccountInfo{
		ResourceType: resourceType,
		ResourceName: resourceName,
		AccountId:    accountId,
		PrincipalId:  principalId,
		UserName:     userName,
	}
}

func parseUserIdentityFields(data map[string]interface{}) (accountId, principalId, userName string) {
	if v, ok := data["userIdentity.accountId"]; ok {
		accountId = v.(string)
	}
	if v, ok := data["userIdentity.principalId"]; ok {
		principalId = v.(string)
	}
	if v, ok := data["userIdentity.userName"]; ok {
		userName = v.(string)
	}
	if v, ok := data["userIdentity"]; ok {
		switch v := v.(type) {
		case string:
			var userIdentity map[string]string
			err := json.Unmarshal([]byte(v), &userIdentity)
			if err == nil {
				accountId = userIdentity["accountId"]
				principalId = userIdentity["principalId"]
				userName = userIdentity["userName"]
			}
		}
	}
	return
}
