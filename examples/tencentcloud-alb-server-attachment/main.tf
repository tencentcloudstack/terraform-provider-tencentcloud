resource "tencentcloud_alb_server_attachment" "service_layer7" {
  loadbalancer_id = "lb-qk1dqox5"
  listener_id     = "lbl-ghoke4tl"
  location_id     = "loc-i858qv1l"

  backends = [
    {
      instance_id = "ins-4j30i5pe"
      port        = 80
      weight      = 10
    },
    {
      instance_id = "ins-4j30i5pe"
      port        = 83
      weight      = 10
    },
    {
      instance_id = "ins-4j30i5pe"
      port        = 84
      weight      = 80
    },
    {
      instance_id = "ins-4j30i5pe"
      port        = 86
      weight      = 10
    },
  ]
}

resource "tencentcloud_alb_server_attachment" "service_layer4" {
  loadbalancer_id = "lb-qk1dqox5"
  listener_id     = "lbl-d9kf6jvn"

  backends = [
    {
      instance_id = "ins-4j30i5pe"
      port        = 8080
      weight      = 10
    },
    {
      instance_id = "ins-4j30i5pe"
      port        = 8082
      weight      = 10
    },
  ]
}
