#!/usr/bin/env bats
load suite

@test "it sets the SSH info when the fingerprint is found" {
  setup_ssh

  run $BASE_PATH/ice-agent register-self --session-id $SESSION_ID \
    --api-endpoint $PUBLIC_ICE_REGISTRY
  echo "Exit status is $status"
  echo "Output is '$output'"

  teardown_ssh

  [ $status -eq 0 ]
  inst_id=$(echo $output | jq '.id' | tr -d '"')
  [ "$inst_id" != "" ]
  [ "$(echo $output | jq '.ssh_authorized_fingerprint' | tr -d '"')" == "30:b6:cb:7e:0b:a3:5a:56:b2:f2:c7:c3:16:1d:2f:db" ]
}
