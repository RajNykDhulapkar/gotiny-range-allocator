#!/bin/bash

# test_grpc.sh

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color
YELLOW='\033[1;33m'

# Server details
SERVER="localhost:50051"
SERVICE="rangeallocator.v1.RangeAllocator"

# Function to print test header
print_header() {
  echo -e "\n${YELLOW}=== Testing $1 ===${NC}"
}

# Function to check if grpcurl is installed
check_grpcurl() {
  if ! command -v grpcurl &>/dev/null; then
    echo "grpcurl is not installed. Installing..."
    go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
  fi
}

# Function to test an endpoint
test_endpoint() {
  local method=$1
  local data=$2
  local description=$3

  print_header "$description"

  if [ -z "$data" ]; then
    echo "Testing $method (no data)"
    grpcurl -plaintext "$SERVER" "$SERVICE/$method"
  else
    echo "Testing $method with data: $data"
    grpcurl -plaintext -d "$data" "$SERVER" "$SERVICE/$method"
  fi

  if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ $method test passed${NC}"
  else
    echo -e "${RED}✗ $method test failed${NC}"
  fi
}

# Main test execution
main() {
  # Check if grpcurl is installed
  check_grpcurl

  echo -e "${YELLOW}Starting gRPC service tests on $SERVER${NC}"

  # List all available services
  print_header "Available Services"
  grpcurl -plaintext "$SERVER" list

  # Test GetHealth
  test_endpoint "GetHealth" "" "Health Check"

  # Test AllocateRange
  test_endpoint "AllocateRange" '{"service_id": "test-service-1"}' "Allocate Range"

  # Store range_id for subsequent tests
  RANGE_ID=$(grpcurl -plaintext -d '{"service_id": "test-service-1"}' "$SERVER" "$SERVICE/AllocateRange" | jq -r '.range.range_id')

  # Test GetRange
  test_endpoint "GetRange" "{\"range_id\": \"$RANGE_ID\"}" "Get Range"

  # Test ListRanges
  test_endpoint "ListRanges" '{"service_id": "test-service-1", "page_size": 10}' "List Ranges"

  # Test UpdateRangeStatus
  test_endpoint "UpdateRangeStatus" "{\"range_id\": \"$RANGE_ID\", \"service_id\": \"test-service-1\", \"status\": \"RANGE_STATUS_EXHAUSTED\"}" "Update Range Status"

  # Test invalid inputs
  print_header "Testing Invalid Inputs"

  echo "Testing AllocateRange with empty service_id"
  grpcurl -plaintext -d '{"service_id": ""}' "$SERVER" "$SERVICE/AllocateRange" || echo -e "${GREEN}✓ Validation working${NC}"

  echo "Testing GetRange with empty range_id"
  grpcurl -plaintext -d '{"range_id": ""}' "$SERVER" "$SERVICE/GetRange" || echo -e "${GREEN}✓ Validation working${NC}"
}

# Run main if script is executed directly
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
  main "$@"
fi
