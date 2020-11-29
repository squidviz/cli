#!/usr/bin/env bats

declare output
declare status

setup() {
  rm .svfile | true
}

@test "add" {
  run eval 'echo "2" | go run main.go add --label "x"'
  [ "$status" -eq 0 ]
  [ "$output" = "x 2" ]

  run jq -r '.pull_request.data | length' .svfile
  [ "$output" = 1 ]

  run jq -r '.pull_request.data[0].label' .svfile
  [ "$output" = "x" ]

  run jq -r '.pull_request.data[0].value' .svfile
  [ "$output" = 2 ]

  run eval 'echo "3" | go run main.go add --label "y"'
  [ "$status" -eq 0 ]
  [ "$output" = "y 3" ]

  run jq -r '.pull_request.data | length' .svfile
  [ "$output" = 2 ]
}

@test "time" {
  run eval 'echo "1" | go run main.go time --label "x"'
  [ "$status" -eq 0 ]
  [ -r .svfile ]

  run jq -r '.pull_request.data | length' .svfile
  [ "$output" = "1" ]

  run jq -r '.pull_request.data[0].label' .svfile
  [ "$output" = "x" ]

  run jq -r '.pull_request.data[0].value | type == "number"' .svfile
  [ "$output" = "true" ]
}
