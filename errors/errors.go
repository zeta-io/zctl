/**
2 * @Author: Nico
3 * @Date: 2021/1/18 0:08
4 */
package errors

import (
	"errors"
	"fmt"
)

var (
	ErrSchemaFormatError  = errors.New("schema format error %s")
	ErrFieldTypeError     = errors.New("field type error %s")
	ErrBuildInputNotExist = errors.New("build input not exist")
	ErrOutputIsNotFile    = errors.New("output is not file")
	ErrTypesFormat    = errors.New("type format is error")
	ErrTypeIsNotPrimitive = errors.New("type is not primitive")
)

func WrapperError(err error, args ...interface{}) error {
	return errors.New(fmt.Sprintf(err.Error(), args...))
}
