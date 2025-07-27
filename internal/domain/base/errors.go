package base

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrDuplicatedKey = errors.New("duplicate key error")
