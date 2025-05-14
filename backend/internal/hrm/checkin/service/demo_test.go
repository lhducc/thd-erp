package service

import (
  "testing"

  "github.com/stretchr/testify/assert"
)


func TestSum(t *testing.T) {
  assert.Equal(t, 5, Sum(2, 3), "2 + 3 must equal 5")
  assert.NotEqual(t, 6, Sum(2, 3), "2 + 3 not must equal  6")
}
