#!/usr/bin/env bats
load suite

@test "it sets the tags" {
  run $BASE_PATH/ice-agent register-self --session-id $SESSION_ID \
    --api-endpoint $PUBLIC_ICE_REGISTRY --tag name=my-name --tag version=1
  echo "Exit status is $status"
  echo "Output is '$output'"

  [ $status -eq 0 ]
  inst_id=$(echo $output | jq '.id' | tr -d '"')
  [ "$inst_id" != "" ]
  [ $(echo $output | jq '.tags.name' | tr -d '"') == "my-name" ]
  [ $(echo $output | jq '.tags.version' | tr -d '"') == "1" ]

  curl -XDELETE $PUBLIC_ICE_REGISTRY/v2/instances/$inst_id
}

@test "it returns a useful error when the tags is misused" {
  run $BASE_PATH/ice-agent register-self --session-id $SESSION_ID \
    --api-endpoint $PUBLIC_ICE_REGISTRY --tag foo
  echo "Exit status is $status"
  echo "Output is '$output'"

  [ $status -ne 0 ]
  [[ "$output" =~ "invalid tag provided `foo`" ]]
}
