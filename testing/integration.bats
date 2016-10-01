#!/usr/bin/env bats

###############################################################################
# Constants
###############################################################################

BASE_PATH=$(cd $BATS_TEST_DIRNAME/..; pwd)
PUBLIC_ICE_REGISTRY="http://ice-registry.cfapps.io"

###############################################################################
# Helpers
###############################################################################

setup() {
  mkdir -p $HOME/.ssh
  cat $BASE_PATH/testing/assets/test_1.pub >> $HOME/.ssh/authorized_keys
  cat $BASE_PATH/testing/assets/test_2.pub >> $HOME/.ssh/authorized_keys

  MY_IP=$(curl -s -XGET $PUBLIC_ICE_REGISTRY/v2/my_ip)
  json=$(curl -s -XPOST $PUBLIC_ICE_REGISTRY/v2/sessions \
    -H 'Content-type: application/json' \
    -d "{\"client_ip_addr\": \"$MY_IP\"}")
  SESSION_ID=$(echo $json | jq '._id' | tr -d '"')
}

teardown() {
  curl -s -XDELETE $PUBLIC_ICE_REGISTRY/v2/sessions/$SESSION_ID

  rm -rf $HOME/.ssh
}

###############################################################################
# Test cases
###############################################################################

@test "it complaints in meaninful ways" {
  rm -rf /home/ice/.ssh

  run $BASE_PATH/ice-agent register-self --session-id foo-bar \
    --api-endpoint $PUBLIC_ICE_REGISTRY
  echo "Exit status is $status"
  echo "Output is '$output'"

  [ $status -ne 0 ]
  [[ "$output" =~ "/home/ice/.ssh/authorized_keys: no such file or directory" ]]
}

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
  [ "$(echo $output | jq '.public_ip_addr' | tr -d '"')" == $MY_IP ]

  curl -XDELETE $PUBLIC_ICE_REGISTRY/v2/instances/$inst_id
}
