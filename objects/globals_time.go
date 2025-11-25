package objects

import (
	"github.com/senither/zen-lang/objects/timer"
)

func globalTimeNow(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgumentsError("now", 0, len(args))
	}

	return &Integer{Value: timer.Now()}, nil
}

func globalTimeSleep(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("sleep", 1, len(args))
	}

	switch v := args[0].(type) {
	case *Integer:
		timer.Sleep(v.Value)
		return NULL, nil
	case *Float:
		timer.Sleep(int64(v.Value))
		return NULL, nil

	default:
		return nil, NewInvalidArgumentTypesError("sleep", []ObjectType{INTEGER_OBJ, FLOAT_OBJ}, 0, args)
	}
}

func globalTimeParse(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("parse", 2, len(args))
	}

	dateStringObj, layoutObj := args[0], args[1]

	dateString, ok := dateStringObj.(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("parse", STRING_OBJ, 0, args)
	}

	layout, ok := layoutObj.(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("parse", STRING_OBJ, 1, args)
	}

	timestamp, err := timer.Parse(dateString.Value, layout.Value)
	if err != nil {
		return nil, NewErrorf("parse", "%s", err.Error())
	}

	return &Integer{Value: timestamp}, nil
}

func globalTimeFormat(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("format", 2, len(args))
	}

	timeObj, formatObj := args[0], args[1]

	var timestamp int64
	switch v := timeObj.(type) {
	case *Integer:
		timestamp = v.Value
	case *Float:
		timestamp = int64(v.Value)

	default:
		return nil, NewInvalidArgumentTypesError("format", []ObjectType{INTEGER_OBJ, FLOAT_OBJ}, 0, args)
	}

	format, ok := formatObj.(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("format", STRING_OBJ, 1, args)
	}

	return &String{Value: timer.Format(timestamp, format.Value)}, nil
}
