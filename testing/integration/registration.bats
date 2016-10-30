#!/usr/bin/env bats
load suite

@test "it conveys useful errors from the iCE registry" {
  run $BASE_PATH/ice-agent register-self --session-id foo-bar \
    --api-endpoint $PUBLIC_ICE_REGISTRY
  echo "Exit status is $status"
  echo "Output is '$output'"

  [ $status -ne 0 ]
  [[ "$output" =~ "value 'foo-bar' cannot be converted to a ObjectId" ]]
}

@test "it outputs the instance JSON on success" {
  run $BASE_PATH/ice-agent register-self --session-id $SESSION_ID \
    --api-endpoint $PUBLIC_ICE_REGISTRY
  echo "Exit status is $status"
  echo "Output is '$output'"

  [ $status -eq 0 ]
  inst_id=$(echo $output | jq '.id' | tr -d '"')
  [ "$inst_id" != "" ]
  [ $(echo $output | jq '.public_ip_addr' | tr -d '"') == $MY_IP ]

  curl -XDELETE $PUBLIC_ICE_REGISTRY/v2/instances/$inst_id
}
