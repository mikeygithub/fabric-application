#!/usr/bin/env bash

function install() {

    for peer in 1 2 3; do
       echo "Install chain-code ${CHAINCODE_NAME} to peer0.org${peer}"
       docker exec cli${peer} peer chaincode install -n ${CHAINCODE_NAME} -v ${VERSION} -p ${CC_SRC_PATH} -l ${LANGUAGE} >&log.txt
       cat log.txt
    done
}

LANGUAGE="golang"
VERSION="1.0"
CHAINCODE_NAME="chaincode-name"

install @?