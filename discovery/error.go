package discovery

import "errors"

var ErrNotFound = errors.New("no service addresses found")
var ServiceIsNotRegistered = errors.New("service is not registered yet")
var ServiceInstanceIsNotRegistered = errors.New("service instance is not registered yet")
