// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Code generated by thriftrw v1.18.0. DO NOT EDIT.
// @generated

package history

import (
	errors "errors"
	fmt "fmt"
	shared "github.com/uber/cadence/.gen/go/shared"
	multierr "go.uber.org/multierr"
	wire "go.uber.org/thriftrw/wire"
	zapcore "go.uber.org/zap/zapcore"
	strings "strings"
)

// HistoryService_RespondActivityTaskFailed_Args represents the arguments for the HistoryService.RespondActivityTaskFailed function.
//
// The arguments for RespondActivityTaskFailed are sent and received over the wire as this struct.
type HistoryService_RespondActivityTaskFailed_Args struct {
	FailRequest *RespondActivityTaskFailedRequest `json:"failRequest,omitempty"`
}

// ToWire translates a HistoryService_RespondActivityTaskFailed_Args struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *HistoryService_RespondActivityTaskFailed_Args) ToWire() (wire.Value, error) {
	var (
		fields [1]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.FailRequest != nil {
		w, err = v.FailRequest.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

func _RespondActivityTaskFailedRequest_1_Read(w wire.Value) (*RespondActivityTaskFailedRequest, error) {
	var v RespondActivityTaskFailedRequest
	err := v.FromWire(w)
	return &v, err
}

// FromWire deserializes a HistoryService_RespondActivityTaskFailed_Args struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a HistoryService_RespondActivityTaskFailed_Args struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v HistoryService_RespondActivityTaskFailed_Args
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *HistoryService_RespondActivityTaskFailed_Args) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.FailRequest, err = _RespondActivityTaskFailedRequest_1_Read(field.Value)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

// String returns a readable string representation of a HistoryService_RespondActivityTaskFailed_Args
// struct.
func (v *HistoryService_RespondActivityTaskFailed_Args) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [1]string
	i := 0
	if v.FailRequest != nil {
		fields[i] = fmt.Sprintf("FailRequest: %v", v.FailRequest)
		i++
	}

	return fmt.Sprintf("HistoryService_RespondActivityTaskFailed_Args{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this HistoryService_RespondActivityTaskFailed_Args match the
// provided HistoryService_RespondActivityTaskFailed_Args.
//
// This function performs a deep comparison.
func (v *HistoryService_RespondActivityTaskFailed_Args) Equals(rhs *HistoryService_RespondActivityTaskFailed_Args) bool {
	if v == nil {
		return rhs == nil
	} else if rhs == nil {
		return false
	}
	if !((v.FailRequest == nil && rhs.FailRequest == nil) || (v.FailRequest != nil && rhs.FailRequest != nil && v.FailRequest.Equals(rhs.FailRequest))) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of HistoryService_RespondActivityTaskFailed_Args.
func (v *HistoryService_RespondActivityTaskFailed_Args) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	if v == nil {
		return nil
	}
	if v.FailRequest != nil {
		err = multierr.Append(err, enc.AddObject("failRequest", v.FailRequest))
	}
	return err
}

// GetFailRequest returns the value of FailRequest if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Args) GetFailRequest() (o *RespondActivityTaskFailedRequest) {
	if v != nil && v.FailRequest != nil {
		return v.FailRequest
	}

	return
}

// IsSetFailRequest returns true if FailRequest is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Args) IsSetFailRequest() bool {
	return v != nil && v.FailRequest != nil
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the arguments.
//
// This will always be "RespondActivityTaskFailed" for this struct.
func (v *HistoryService_RespondActivityTaskFailed_Args) MethodName() string {
	return "RespondActivityTaskFailed"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Call for this struct.
func (v *HistoryService_RespondActivityTaskFailed_Args) EnvelopeType() wire.EnvelopeType {
	return wire.Call
}

// HistoryService_RespondActivityTaskFailed_Helper provides functions that aid in handling the
// parameters and return values of the HistoryService.RespondActivityTaskFailed
// function.
var HistoryService_RespondActivityTaskFailed_Helper = struct {
	// Args accepts the parameters of RespondActivityTaskFailed in-order and returns
	// the arguments struct for the function.
	Args func(
		failRequest *RespondActivityTaskFailedRequest,
	) *HistoryService_RespondActivityTaskFailed_Args

	// IsException returns true if the given error can be thrown
	// by RespondActivityTaskFailed.
	//
	// An error can be thrown by RespondActivityTaskFailed only if the
	// corresponding exception type was mentioned in the 'throws'
	// section for it in the Thrift file.
	IsException func(error) bool

	// WrapResponse returns the result struct for RespondActivityTaskFailed
	// given the error returned by it. The provided error may
	// be nil if RespondActivityTaskFailed did not fail.
	//
	// This allows mapping errors returned by RespondActivityTaskFailed into a
	// serializable result struct. WrapResponse returns a
	// non-nil error if the provided error cannot be thrown by
	// RespondActivityTaskFailed
	//
	//   err := RespondActivityTaskFailed(args)
	//   result, err := HistoryService_RespondActivityTaskFailed_Helper.WrapResponse(err)
	//   if err != nil {
	//     return fmt.Errorf("unexpected error from RespondActivityTaskFailed: %v", err)
	//   }
	//   serialize(result)
	WrapResponse func(error) (*HistoryService_RespondActivityTaskFailed_Result, error)

	// UnwrapResponse takes the result struct for RespondActivityTaskFailed
	// and returns the erorr returned by it (if any).
	//
	// The error is non-nil only if RespondActivityTaskFailed threw an
	// exception.
	//
	//   result := deserialize(bytes)
	//   err := HistoryService_RespondActivityTaskFailed_Helper.UnwrapResponse(result)
	UnwrapResponse func(*HistoryService_RespondActivityTaskFailed_Result) error
}{}

func init() {
	HistoryService_RespondActivityTaskFailed_Helper.Args = func(
		failRequest *RespondActivityTaskFailedRequest,
	) *HistoryService_RespondActivityTaskFailed_Args {
		return &HistoryService_RespondActivityTaskFailed_Args{
			FailRequest: failRequest,
		}
	}

	HistoryService_RespondActivityTaskFailed_Helper.IsException = func(err error) bool {
		switch err.(type) {
		case *shared.BadRequestError:
			return true
		case *shared.InternalServiceError:
			return true
		case *shared.EntityNotExistsError:
			return true
		case *ShardOwnershipLostError:
			return true
		case *shared.DomainNotActiveError:
			return true
		case *shared.LimitExceededError:
			return true
		case *shared.ServiceBusyError:
			return true
		default:
			return false
		}
	}

	HistoryService_RespondActivityTaskFailed_Helper.WrapResponse = func(err error) (*HistoryService_RespondActivityTaskFailed_Result, error) {
		if err == nil {
			return &HistoryService_RespondActivityTaskFailed_Result{}, nil
		}

		switch e := err.(type) {
		case *shared.BadRequestError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.BadRequestError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{BadRequestError: e}, nil
		case *shared.InternalServiceError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.InternalServiceError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{InternalServiceError: e}, nil
		case *shared.EntityNotExistsError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.EntityNotExistError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{EntityNotExistError: e}, nil
		case *ShardOwnershipLostError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.ShardOwnershipLostError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{ShardOwnershipLostError: e}, nil
		case *shared.DomainNotActiveError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.DomainNotActiveError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{DomainNotActiveError: e}, nil
		case *shared.LimitExceededError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.LimitExceededError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{LimitExceededError: e}, nil
		case *shared.ServiceBusyError:
			if e == nil {
				return nil, errors.New("WrapResponse received non-nil error type with nil value for HistoryService_RespondActivityTaskFailed_Result.ServiceBusyError")
			}
			return &HistoryService_RespondActivityTaskFailed_Result{ServiceBusyError: e}, nil
		}

		return nil, err
	}
	HistoryService_RespondActivityTaskFailed_Helper.UnwrapResponse = func(result *HistoryService_RespondActivityTaskFailed_Result) (err error) {
		if result.BadRequestError != nil {
			err = result.BadRequestError
			return
		}
		if result.InternalServiceError != nil {
			err = result.InternalServiceError
			return
		}
		if result.EntityNotExistError != nil {
			err = result.EntityNotExistError
			return
		}
		if result.ShardOwnershipLostError != nil {
			err = result.ShardOwnershipLostError
			return
		}
		if result.DomainNotActiveError != nil {
			err = result.DomainNotActiveError
			return
		}
		if result.LimitExceededError != nil {
			err = result.LimitExceededError
			return
		}
		if result.ServiceBusyError != nil {
			err = result.ServiceBusyError
			return
		}
		return
	}

}

// HistoryService_RespondActivityTaskFailed_Result represents the result of a HistoryService.RespondActivityTaskFailed function call.
//
// The result of a RespondActivityTaskFailed execution is sent and received over the wire as this struct.
type HistoryService_RespondActivityTaskFailed_Result struct {
	BadRequestError         *shared.BadRequestError      `json:"badRequestError,omitempty"`
	InternalServiceError    *shared.InternalServiceError `json:"internalServiceError,omitempty"`
	EntityNotExistError     *shared.EntityNotExistsError `json:"entityNotExistError,omitempty"`
	ShardOwnershipLostError *ShardOwnershipLostError     `json:"shardOwnershipLostError,omitempty"`
	DomainNotActiveError    *shared.DomainNotActiveError `json:"domainNotActiveError,omitempty"`
	LimitExceededError      *shared.LimitExceededError   `json:"limitExceededError,omitempty"`
	ServiceBusyError        *shared.ServiceBusyError     `json:"serviceBusyError,omitempty"`
}

// ToWire translates a HistoryService_RespondActivityTaskFailed_Result struct into a Thrift-level intermediate
// representation. This intermediate representation may be serialized
// into bytes using a ThriftRW protocol implementation.
//
// An error is returned if the struct or any of its fields failed to
// validate.
//
//   x, err := v.ToWire()
//   if err != nil {
//     return err
//   }
//
//   if err := binaryProtocol.Encode(x, writer); err != nil {
//     return err
//   }
func (v *HistoryService_RespondActivityTaskFailed_Result) ToWire() (wire.Value, error) {
	var (
		fields [7]wire.Field
		i      int = 0
		w      wire.Value
		err    error
	)

	if v.BadRequestError != nil {
		w, err = v.BadRequestError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 1, Value: w}
		i++
	}
	if v.InternalServiceError != nil {
		w, err = v.InternalServiceError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 2, Value: w}
		i++
	}
	if v.EntityNotExistError != nil {
		w, err = v.EntityNotExistError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 3, Value: w}
		i++
	}
	if v.ShardOwnershipLostError != nil {
		w, err = v.ShardOwnershipLostError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 4, Value: w}
		i++
	}
	if v.DomainNotActiveError != nil {
		w, err = v.DomainNotActiveError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 5, Value: w}
		i++
	}
	if v.LimitExceededError != nil {
		w, err = v.LimitExceededError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 6, Value: w}
		i++
	}
	if v.ServiceBusyError != nil {
		w, err = v.ServiceBusyError.ToWire()
		if err != nil {
			return w, err
		}
		fields[i] = wire.Field{ID: 7, Value: w}
		i++
	}

	if i > 1 {
		return wire.Value{}, fmt.Errorf("HistoryService_RespondActivityTaskFailed_Result should have at most one field: got %v fields", i)
	}

	return wire.NewValueStruct(wire.Struct{Fields: fields[:i]}), nil
}

// FromWire deserializes a HistoryService_RespondActivityTaskFailed_Result struct from its Thrift-level
// representation. The Thrift-level representation may be obtained
// from a ThriftRW protocol implementation.
//
// An error is returned if we were unable to build a HistoryService_RespondActivityTaskFailed_Result struct
// from the provided intermediate representation.
//
//   x, err := binaryProtocol.Decode(reader, wire.TStruct)
//   if err != nil {
//     return nil, err
//   }
//
//   var v HistoryService_RespondActivityTaskFailed_Result
//   if err := v.FromWire(x); err != nil {
//     return nil, err
//   }
//   return &v, nil
func (v *HistoryService_RespondActivityTaskFailed_Result) FromWire(w wire.Value) error {
	var err error

	for _, field := range w.GetStruct().Fields {
		switch field.ID {
		case 1:
			if field.Value.Type() == wire.TStruct {
				v.BadRequestError, err = _BadRequestError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 2:
			if field.Value.Type() == wire.TStruct {
				v.InternalServiceError, err = _InternalServiceError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 3:
			if field.Value.Type() == wire.TStruct {
				v.EntityNotExistError, err = _EntityNotExistsError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 4:
			if field.Value.Type() == wire.TStruct {
				v.ShardOwnershipLostError, err = _ShardOwnershipLostError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 5:
			if field.Value.Type() == wire.TStruct {
				v.DomainNotActiveError, err = _DomainNotActiveError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 6:
			if field.Value.Type() == wire.TStruct {
				v.LimitExceededError, err = _LimitExceededError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		case 7:
			if field.Value.Type() == wire.TStruct {
				v.ServiceBusyError, err = _ServiceBusyError_Read(field.Value)
				if err != nil {
					return err
				}

			}
		}
	}

	count := 0
	if v.BadRequestError != nil {
		count++
	}
	if v.InternalServiceError != nil {
		count++
	}
	if v.EntityNotExistError != nil {
		count++
	}
	if v.ShardOwnershipLostError != nil {
		count++
	}
	if v.DomainNotActiveError != nil {
		count++
	}
	if v.LimitExceededError != nil {
		count++
	}
	if v.ServiceBusyError != nil {
		count++
	}
	if count > 1 {
		return fmt.Errorf("HistoryService_RespondActivityTaskFailed_Result should have at most one field: got %v fields", count)
	}

	return nil
}

// String returns a readable string representation of a HistoryService_RespondActivityTaskFailed_Result
// struct.
func (v *HistoryService_RespondActivityTaskFailed_Result) String() string {
	if v == nil {
		return "<nil>"
	}

	var fields [7]string
	i := 0
	if v.BadRequestError != nil {
		fields[i] = fmt.Sprintf("BadRequestError: %v", v.BadRequestError)
		i++
	}
	if v.InternalServiceError != nil {
		fields[i] = fmt.Sprintf("InternalServiceError: %v", v.InternalServiceError)
		i++
	}
	if v.EntityNotExistError != nil {
		fields[i] = fmt.Sprintf("EntityNotExistError: %v", v.EntityNotExistError)
		i++
	}
	if v.ShardOwnershipLostError != nil {
		fields[i] = fmt.Sprintf("ShardOwnershipLostError: %v", v.ShardOwnershipLostError)
		i++
	}
	if v.DomainNotActiveError != nil {
		fields[i] = fmt.Sprintf("DomainNotActiveError: %v", v.DomainNotActiveError)
		i++
	}
	if v.LimitExceededError != nil {
		fields[i] = fmt.Sprintf("LimitExceededError: %v", v.LimitExceededError)
		i++
	}
	if v.ServiceBusyError != nil {
		fields[i] = fmt.Sprintf("ServiceBusyError: %v", v.ServiceBusyError)
		i++
	}

	return fmt.Sprintf("HistoryService_RespondActivityTaskFailed_Result{%v}", strings.Join(fields[:i], ", "))
}

// Equals returns true if all the fields of this HistoryService_RespondActivityTaskFailed_Result match the
// provided HistoryService_RespondActivityTaskFailed_Result.
//
// This function performs a deep comparison.
func (v *HistoryService_RespondActivityTaskFailed_Result) Equals(rhs *HistoryService_RespondActivityTaskFailed_Result) bool {
	if v == nil {
		return rhs == nil
	} else if rhs == nil {
		return false
	}
	if !((v.BadRequestError == nil && rhs.BadRequestError == nil) || (v.BadRequestError != nil && rhs.BadRequestError != nil && v.BadRequestError.Equals(rhs.BadRequestError))) {
		return false
	}
	if !((v.InternalServiceError == nil && rhs.InternalServiceError == nil) || (v.InternalServiceError != nil && rhs.InternalServiceError != nil && v.InternalServiceError.Equals(rhs.InternalServiceError))) {
		return false
	}
	if !((v.EntityNotExistError == nil && rhs.EntityNotExistError == nil) || (v.EntityNotExistError != nil && rhs.EntityNotExistError != nil && v.EntityNotExistError.Equals(rhs.EntityNotExistError))) {
		return false
	}
	if !((v.ShardOwnershipLostError == nil && rhs.ShardOwnershipLostError == nil) || (v.ShardOwnershipLostError != nil && rhs.ShardOwnershipLostError != nil && v.ShardOwnershipLostError.Equals(rhs.ShardOwnershipLostError))) {
		return false
	}
	if !((v.DomainNotActiveError == nil && rhs.DomainNotActiveError == nil) || (v.DomainNotActiveError != nil && rhs.DomainNotActiveError != nil && v.DomainNotActiveError.Equals(rhs.DomainNotActiveError))) {
		return false
	}
	if !((v.LimitExceededError == nil && rhs.LimitExceededError == nil) || (v.LimitExceededError != nil && rhs.LimitExceededError != nil && v.LimitExceededError.Equals(rhs.LimitExceededError))) {
		return false
	}
	if !((v.ServiceBusyError == nil && rhs.ServiceBusyError == nil) || (v.ServiceBusyError != nil && rhs.ServiceBusyError != nil && v.ServiceBusyError.Equals(rhs.ServiceBusyError))) {
		return false
	}

	return true
}

// MarshalLogObject implements zapcore.ObjectMarshaler, enabling
// fast logging of HistoryService_RespondActivityTaskFailed_Result.
func (v *HistoryService_RespondActivityTaskFailed_Result) MarshalLogObject(enc zapcore.ObjectEncoder) (err error) {
	if v == nil {
		return nil
	}
	if v.BadRequestError != nil {
		err = multierr.Append(err, enc.AddObject("badRequestError", v.BadRequestError))
	}
	if v.InternalServiceError != nil {
		err = multierr.Append(err, enc.AddObject("internalServiceError", v.InternalServiceError))
	}
	if v.EntityNotExistError != nil {
		err = multierr.Append(err, enc.AddObject("entityNotExistError", v.EntityNotExistError))
	}
	if v.ShardOwnershipLostError != nil {
		err = multierr.Append(err, enc.AddObject("shardOwnershipLostError", v.ShardOwnershipLostError))
	}
	if v.DomainNotActiveError != nil {
		err = multierr.Append(err, enc.AddObject("domainNotActiveError", v.DomainNotActiveError))
	}
	if v.LimitExceededError != nil {
		err = multierr.Append(err, enc.AddObject("limitExceededError", v.LimitExceededError))
	}
	if v.ServiceBusyError != nil {
		err = multierr.Append(err, enc.AddObject("serviceBusyError", v.ServiceBusyError))
	}
	return err
}

// GetBadRequestError returns the value of BadRequestError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetBadRequestError() (o *shared.BadRequestError) {
	if v != nil && v.BadRequestError != nil {
		return v.BadRequestError
	}

	return
}

// IsSetBadRequestError returns true if BadRequestError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetBadRequestError() bool {
	return v != nil && v.BadRequestError != nil
}

// GetInternalServiceError returns the value of InternalServiceError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetInternalServiceError() (o *shared.InternalServiceError) {
	if v != nil && v.InternalServiceError != nil {
		return v.InternalServiceError
	}

	return
}

// IsSetInternalServiceError returns true if InternalServiceError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetInternalServiceError() bool {
	return v != nil && v.InternalServiceError != nil
}

// GetEntityNotExistError returns the value of EntityNotExistError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetEntityNotExistError() (o *shared.EntityNotExistsError) {
	if v != nil && v.EntityNotExistError != nil {
		return v.EntityNotExistError
	}

	return
}

// IsSetEntityNotExistError returns true if EntityNotExistError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetEntityNotExistError() bool {
	return v != nil && v.EntityNotExistError != nil
}

// GetShardOwnershipLostError returns the value of ShardOwnershipLostError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetShardOwnershipLostError() (o *ShardOwnershipLostError) {
	if v != nil && v.ShardOwnershipLostError != nil {
		return v.ShardOwnershipLostError
	}

	return
}

// IsSetShardOwnershipLostError returns true if ShardOwnershipLostError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetShardOwnershipLostError() bool {
	return v != nil && v.ShardOwnershipLostError != nil
}

// GetDomainNotActiveError returns the value of DomainNotActiveError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetDomainNotActiveError() (o *shared.DomainNotActiveError) {
	if v != nil && v.DomainNotActiveError != nil {
		return v.DomainNotActiveError
	}

	return
}

// IsSetDomainNotActiveError returns true if DomainNotActiveError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetDomainNotActiveError() bool {
	return v != nil && v.DomainNotActiveError != nil
}

// GetLimitExceededError returns the value of LimitExceededError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetLimitExceededError() (o *shared.LimitExceededError) {
	if v != nil && v.LimitExceededError != nil {
		return v.LimitExceededError
	}

	return
}

// IsSetLimitExceededError returns true if LimitExceededError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetLimitExceededError() bool {
	return v != nil && v.LimitExceededError != nil
}

// GetServiceBusyError returns the value of ServiceBusyError if it is set or its
// zero value if it is unset.
func (v *HistoryService_RespondActivityTaskFailed_Result) GetServiceBusyError() (o *shared.ServiceBusyError) {
	if v != nil && v.ServiceBusyError != nil {
		return v.ServiceBusyError
	}

	return
}

// IsSetServiceBusyError returns true if ServiceBusyError is not nil.
func (v *HistoryService_RespondActivityTaskFailed_Result) IsSetServiceBusyError() bool {
	return v != nil && v.ServiceBusyError != nil
}

// MethodName returns the name of the Thrift function as specified in
// the IDL, for which this struct represent the result.
//
// This will always be "RespondActivityTaskFailed" for this struct.
func (v *HistoryService_RespondActivityTaskFailed_Result) MethodName() string {
	return "RespondActivityTaskFailed"
}

// EnvelopeType returns the kind of value inside this struct.
//
// This will always be Reply for this struct.
func (v *HistoryService_RespondActivityTaskFailed_Result) EnvelopeType() wire.EnvelopeType {
	return wire.Reply
}