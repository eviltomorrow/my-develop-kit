package xshell

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseStdinRequest(t *testing.T) {
	_assert := assert.New(t)
	_, _, err := parseStdinRequest([]byte("list a"))
	_assert.NotNil(err)

	request, args, err := parseStdinRequest([]byte("list "))
	_assert.Equal(list, request)
	_assert.Nil(err)
	_assert.Nil(args)

	request, args, err = parseStdinRequest([]byte("list|grep a"))
	_assert.NotNil(err)

	request, args, err = parseStdinRequest([]byte("list |grep a"))
	_assert.Equal(list, request)
	_assert.NotNil(args)
	_assert.Nil(err)

	request, args, err = parseStdinRequest([]byte("add  a"))
	_assert.Equal(add, request)
	_assert.NotNil(args)
	_assert.Nil(err)

	request, args, err = parseStdinRequest([]byte("quit  a"))
	_assert.Equal(quit, request)
	_assert.Nil(args)
	_assert.Nil(err)
}
