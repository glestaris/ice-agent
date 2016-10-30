###############################################################################
# Constants
###############################################################################

BASE_PATH=$(cd $BATS_TEST_DIRNAME/../..; pwd)
PUBLIC_ICE_REGISTRY="http://ice-registry.cfapps.io"

###############################################################################
# Helpers
###############################################################################

setup_ssh() {
  mkdir -p $HOME/.ssh
  cat $BASE_PATH/testing/assets/test_1.pub >> $HOME/.ssh/authorized_keys
  cat $BASE_PATH/testing/assets/test_2.pub >> $HOME/.ssh/authorized_keys
}

teardown_ssh() {
  rm -rf $HOME/.ssh
}

setup() {
  MY_IP=$(curl -s -XGET $PUBLIC_ICE_REGISTRY/v2/my_ip)
  json=$(curl -s -XPOST $PUBLIC_ICE_REGISTRY/v2/sessions \
    -H 'Content-type: application/json' \
    -d "{\"client_ip_addr\": \"$MY_IP\"}")
  SESSION_ID=$(echo $json | jq '._id' | tr -d '"')
}

teardown() {
  curl -s -XDELETE $PUBLIC_ICE_REGISTRY/v2/sessions/$SESSION_ID
}
