#!/bin/bash

set -eu

echo -n -e "\n"
echo "===== Starting gRPC test ====="
echo -n -e "\n"

echo "Test GetUsersForMatching"
grpcurl -v -d @ -plaintext -proto pb/mixlunch.proto localhost:8081 pb.MixLunch.GetUsersForMatching \
  < integration-test/grpc_test_inputs/GetUsersForMatching_001.json

echo "Test CreateParties"
grpcurl -v -d @ -plaintext -proto pb/mixlunch.proto localhost:8081 pb.MixLunch.CreateParties \
  < integration-test/grpc_test_inputs/CreateParties_001.json

echo "Test GetParties"
grpcurl -v -d @ -plaintext -proto pb/mixlunch.proto localhost:8081 pb.MixLunch.GetParties \
  < integration-test/grpc_test_inputs/GetParties_001.json

echo -n -e "\n"
echo "===== gRPC test completed ====="
echo -n -e "\n"
