package mqtt

import (
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	mqttv20240516 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mqtt/v20240516"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func NewMqttService(client *connectivity.TencentCloudClient) MqttService {
	return MqttService{client: client}
}

type MqttService struct {
	client *connectivity.TencentCloudClient
}

func (me *MqttService) DescribeMqttById(ctx context.Context, instanceId string) (ret *mqttv20240516.DescribeInstanceResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeInstanceRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttInstancePublicEndpointById(ctx context.Context, instanceId string) (ret *mqttv20240516.DescribeInsPublicEndpointsResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeInsPublicEndpointsRequest()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeInsPublicEndpoints(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttTopicById(ctx context.Context, instanceId string, topic string) (ret *mqttv20240516.DescribeTopicResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeTopicRequest()
	request.InstanceId = helper.String(instanceId)
	request.Topic = helper.String(topic)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseMqttV20240516Client().DescribeTopic(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttRegistrationCodeByFilter(ctx context.Context, param map[string]interface{}) (ret *mqttv20240516.ApplyRegistrationCodeResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mqttv20240516.NewApplyRegistrationCodeRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMqttV20240516Client().ApplyRegistrationCode(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *MqttService) GetCertificateSerialNumber(certData string) (string, error) {
	certData = strings.TrimSpace(certData)
	var cert *x509.Certificate
	var err error

	if strings.Contains(certData, "-----BEGIN CERTIFICATE-----") {
		block, _ := pem.Decode([]byte(certData))
		if block == nil {
			return "", fmt.Errorf("failed to parse certificate PEM")
		}

		cert, err = x509.ParseCertificate(block.Bytes)
	} else {
		cert, err = x509.ParseCertificate([]byte(certData))
	}

	if err != nil {
		return "", fmt.Errorf("failed to parse certificate: %v", err)
	}

	serialHex := hex.EncodeToString(cert.SerialNumber.Bytes())
	return serialHex, nil
}

func (me *MqttService) DescribeMqttCaCertificateById(ctx context.Context, instanceId, caSn string) (ret *mqttv20240516.DescribeCaCertificateResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeCaCertificateRequest()
	response := mqttv20240516.NewDescribeCaCertificateResponse()
	request.InstanceId = &instanceId
	request.CaSn = &caSn
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeCaCertificate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttCaCertificatesById(ctx context.Context, instanceId, caSn string) (data *mqttv20240516.CaCertificateItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeCaCertificatesRequest()
	response := mqttv20240516.NewDescribeCaCertificatesResponse()
	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeCaCertificates(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	if len(response.Response.Data) == 0 {
		return
	}

	for _, item := range response.Response.Data {
		if *item.CaSn == caSn {
			data = item
			return
		}
	}

	return
}

func (me *MqttService) DescribeMqttUserById(ctx context.Context, instanceId, userName string) (ret *mqttv20240516.MQTTUserItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeUserListRequest()
	response := mqttv20240516.NewDescribeUserListResponse()

	request.InstanceId = &instanceId
	if userName != "" {
		request.Filters = []*mqttv20240516.Filter{
			{
				Name:   helper.String("Username"),
				Values: helper.Strings([]string{userName}),
			},
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeUserList(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	if len(response.Response.Data) == 0 {
		return
	}

	ret = response.Response.Data[0]
	return
}

func (me *MqttService) DescribeMqttDeviceCertificateById(ctx context.Context, instanceId, deviceCertificateSn string) (ret *mqttv20240516.DescribeDeviceCertificateResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeDeviceCertificateRequest()
	response := mqttv20240516.NewDescribeDeviceCertificateResponse()
	request.InstanceId = &instanceId
	request.DeviceCertificateSn = &deviceCertificateSn

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeDeviceCertificate(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttJwtAuthenticatorById(ctx context.Context, instanceId string) (ret *mqttv20240516.MQTTAuthenticatorItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeAuthenticatorRequest()
	response := mqttv20240516.NewDescribeAuthenticatorResponse()
	request.InstanceId = &instanceId
	request.Type = helper.String("JWT")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeAuthenticator(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	if len(response.Response.Authenticators) == 0 {
		return
	}

	ret = response.Response.Authenticators[0]
	return
}

func (me *MqttService) DescribeMqttJwksAuthenticatorById(ctx context.Context, instanceId string) (ret *mqttv20240516.MQTTAuthenticatorItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeAuthenticatorRequest()
	response := mqttv20240516.NewDescribeAuthenticatorResponse()
	request.InstanceId = &instanceId
	request.Type = helper.String("JWKS")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeAuthenticator(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	if len(response.Response.Authenticators) == 0 {
		return
	}

	ret = response.Response.Authenticators[0]
	return
}

func (me *MqttService) DescribeMqttHttpAuthenticatorById(ctx context.Context, instanceId string) (ret *mqttv20240516.MQTTAuthenticatorItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeAuthenticatorRequest()
	response := mqttv20240516.NewDescribeAuthenticatorResponse()
	request.InstanceId = &instanceId
	request.Type = helper.String("HTTP")

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeAuthenticator(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	if len(response.Response.Authenticators) == 0 {
		return
	}

	ret = response.Response.Authenticators[0]

	return
}

func (me *MqttService) DescribeMqttAuthorizationPolicyById(ctx context.Context, instanceId string) (ret *mqttv20240516.DescribeAuthorizationPoliciesResponseParams, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeAuthorizationPoliciesRequest()
	response := mqttv20240516.NewDescribeAuthorizationPoliciesResponse()
	request.InstanceId = &instanceId

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := me.client.UseMqttV20240516Client().DescribeAuthorizationPolicies(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttInstancesByFilter(ctx context.Context, param map[string]interface{}) (ret []*mqttv20240516.MQTTInstanceItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mqttv20240516.NewDescribeInstanceListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "Filters" {
			request.Filters = v.([]*mqttv20240516.Filter)
		}

		if k == "TagFilters" {
			request.TagFilters = v.([]*mqttv20240516.TagFilter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)

	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMqttV20240516Client().DescribeInstanceList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}

		ret = append(ret, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MqttService) DescribeMqttInstanceDetailByFilter(ctx context.Context, param map[string]interface{}) (ret *mqttv20240516.DescribeInstanceResponseParams, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mqttv20240516.NewDescribeInstanceRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}
	}

	ratelimit.Check(request.GetAction())
	response, err := me.client.UseMqttV20240516Client().DescribeInstance(request)
	if err != nil {
		errRet = err
		return
	}

	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil {
		return
	}

	ret = response.Response
	return
}

func (me *MqttService) DescribeMqttTopicsByFilter(ctx context.Context, param map[string]interface{}) (ret []*mqttv20240516.MQTTTopicItem, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = mqttv20240516.NewDescribeTopicListRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "InstanceId" {
			request.InstanceId = v.(*string)
		}

		if k == "Filters" {
			request.Filters = v.([]*mqttv20240516.Filter)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset int64 = 0
		limit  int64 = 100
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseMqttV20240516Client().DescribeTopicList(request)
		if err != nil {
			errRet = err
			return
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Data) < 1 {
			break
		}

		ret = append(ret, response.Response.Data...)
		if len(response.Response.Data) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *MqttService) DescribeMqttMessageEnrichmentRuleById(ctx context.Context, instanceId, ruleId string) (ret *mqttv20240516.MessageEnrichmentRuleItem, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDescribeMessageEnrichmentRulesRequest()
	response := mqttv20240516.NewDescribeMessageEnrichmentRulesResponse()
	request.InstanceId = helper.String(instanceId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMqttV20240516Client().DescribeMessageEnrichmentRules(request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe message enrichment rules failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	// Find the specific rule by rule ID
	for _, rule := range response.Response.Data {
		if rule != nil && rule.Id != nil && helper.Int64ToStr(*rule.Id) == ruleId {
			ret = rule
			break
		}
	}

	return
}

func (me *MqttService) DeleteMqttMessageEnrichmentRuleById(ctx context.Context, instanceId, ruleId string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := mqttv20240516.NewDeleteMessageEnrichmentRuleRequest()
	request.InstanceId = helper.String(instanceId)
	request.Id = helper.StrToInt64Point(ruleId)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseMqttV20240516Client().DeleteMessageEnrichmentRuleWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		errRet = reqErr
		return
	}

	return
}
