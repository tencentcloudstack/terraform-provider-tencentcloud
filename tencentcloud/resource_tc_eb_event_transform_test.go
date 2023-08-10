package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

// go test -i; go test -test.run TestAccTencentCloudEbEventTransformResource_basic -v
func TestAccTencentCloudEbEventTransformResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEbEventTransformDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEbEventTransform,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEbEventTransformExists("tencentcloud_eb_event_transform.eb_transform"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "id"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_transform.eb_transform", "transformations.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_transform.eb_transform", "transformations.0.extraction.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_transform.eb_transform", "transformations.0.extraction.0.extraction_input_path", "$"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_transform.eb_transform", "transformations.0.extraction.0.format", "JSON"),
					resource.TestCheckResourceAttr("tencentcloud_eb_event_transform.eb_transform", "transformations.0.transform.#", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "transformations.0.transform.0.output_structs.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "transformations.0.transform.0.output_structs.0.key"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "transformations.0.transform.0.output_structs.0.value"),
					resource.TestCheckResourceAttrSet("tencentcloud_eb_event_transform.eb_transform", "transformations.0.transform.0.output_structs.0.value_type"),
				),
			},
			{
				ResourceName:      "tencentcloud_eb_event_transform.eb_transform",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckEbEventTransformDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_eb_event_transform" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]
		transformationId := idSplit[2]

		target, err := service.DescribeEbEventTransformById(ctx, eventBusId, ruleId, transformationId)
		if err != nil {
			if err.(*sdkErrors.TencentCloudSDKError).Code == "ResourceNotFound.Transformation" {
				return nil
			}
			return err
		}

		if target != nil {
			return fmt.Errorf("eb transformation %s still exists", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckEbEventTransformExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}
		idSplit := strings.Split(rs.Primary.ID, FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		eventBusId := idSplit[0]
		ruleId := idSplit[1]
		transformationId := idSplit[2]

		service := EbService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		target, err := service.DescribeEbEventTransformById(ctx, eventBusId, ruleId, transformationId)
		if err != nil {
			return err
		}

		if target == nil {
			return fmt.Errorf("eb transformation %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testAccEbEventTransformVar = `
resource "tencentcloud_eb_event_bus" "foo" {
  event_bus_name = "tf-event_bus"
  description    = "event bus desc"
  enable_store   = false
  save_days      = 1
  tags = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_eb_event_rule" "foo" {
  event_bus_id = tencentcloud_eb_event_bus.foo.id
  rule_name    = "tf-event_rule"
  description  = "event rule desc"
  enable       = true
  event_pattern = jsonencode(
    {
      source = "apigw.cloud.tencent"
      type = [
        "connector:apigw",
      ]
    }
  )
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testAccEbEventTransform = testAccEbEventTransformVar + `

resource "tencentcloud_eb_event_transform" "eb_transform" {
    event_bus_id = tencentcloud_eb_event_bus.foo.id
    rule_id      = tencentcloud_eb_event_rule.foo.rule_id

    transformations {

        extraction {
            extraction_input_path = "$"
            format                = "JSON"
        }

        transform {
            output_structs {
                key        = "type"
                value      = "connector:ckafka"
                value_type = "STRING"
            }
            output_structs {
                key        = "source"
                value      = "ckafka.cloud.tencent"
                value_type = "STRING"
            }
            output_structs {
                key        = "region"
                value      = "ap-guangzhou"
                value_type = "STRING"
            }
            output_structs {
                key        = "datacontenttype"
                value      = "application/json;charset=utf-8"
                value_type = "STRING"
            }
            output_structs {
                key        = "status"
                value      = "-"
                value_type = "STRING"
            }
            output_structs {
                key        = "data"
                value      = jsonencode(
                    {
                        Partition = 1
                        msgBody   = "Hello from Ckafka again!"
                        msgKey    = "test"
                        offset    = 37
                        topic     = "test-topic"
                    }
                )
                value_type = "STRING"
            }
        }
    }
}

`
