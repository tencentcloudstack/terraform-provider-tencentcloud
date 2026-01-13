package mqtt_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	svcmqtt "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mqtt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTencentCloudMqttMessageEnrichmentRuleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMqttMessageEnrichmentRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttMessageEnrichmentRule,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMqttMessageEnrichmentRuleExists("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "rule_name", "test-enrichment-rule"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "status", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "rule_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "priority"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "priority", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccMqttMessageEnrichmentRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMqttMessageEnrichmentRuleExists("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "rule_name", "test-enrichment-rule-updated"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "status", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mqtt_message_enrichment_rule.message_enrichment_rule", "remark", "updated remark"),
				),
			},
		},
	})
}

func testAccCheckMqttMessageEnrichmentRuleDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmqtt.NewMqttService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mqtt_message_enrichment_rule" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		ruleId := idSplit[1]

		ruleIdInt64 := helper.StrToInt64(ruleId)

		rule, err := service.DescribeMqttMessageEnrichmentRuleById(ctx, instanceId, ruleIdInt64)
		if err != nil {
			return err
		}

		if rule != nil {
			return fmt.Errorf("mqtt message enrichment rule %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckMqttMessageEnrichmentRuleExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		instanceId := idSplit[0]
		ruleId := idSplit[1]

		ruleIdInt64 := helper.StrToInt64(ruleId)

		service := svcmqtt.NewMqttService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		rule, err := service.DescribeMqttMessageEnrichmentRuleById(ctx, instanceId, ruleIdInt64)
		if err != nil {
			return err
		}

		if rule == nil {
			return fmt.Errorf("mqtt message enrichment rule %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccMqttMessageEnrichmentRule = `
resource "tencentcloud_mqtt_message_enrichment_rule" "message_enrichment_rule" {
  instance_id = "mqtt-zxjwkr98"
  rule_name   = "test-enrichment-rule"
  condition   = "eyJ0eXBlIjoiYW5kIiwicnVsZXMiOlt7ImZpZWxkIjoidG9waWMiLCJvcGVyYXRvciI6Ij09IiwidmFsdWUiOiJ0ZXN0L3RvcGljIn1dfQ=="
  actions     = "W3sidHlwZSI6InNldF9wcm9wZXJ0eSIsImtleSI6InByb3BlcnR5MSIsInZhbHVlIjoidGVzdF92YWx1ZTEifV0="
  status      = 1
  priority    = 1
  remark      = "test remark"
}
`

const testAccMqttMessageEnrichmentRuleUpdate = `
resource "tencentcloud_mqtt_message_enrichment_rule" "message_enrichment_rule" {
  instance_id = "mqtt-zxjwkr98"
  rule_name   = "test-enrichment-rule-updated"
  condition   = "eyJ0eXBlIjoiYW5kIiwicnVsZXMiOlt7ImZpZWxkIjoidG9waWMiLCJvcGVyYXRvciI6Ij09IiwidmFsdWUiOiJ0ZXN0L3RvcGljMiJ9XX0="
  actions     = "W3sidHlwZSI6InNldF9wcm9wZXJ0eSIsImtleSI6InByb3BlcnR5MiIsInZhbHVlIjoidGVzdF92YWx1ZTIifV0="
  status      = 0
  priority    = 2
  remark      = "updated remark"
}
`
