package objects

import (
	"fmt"
	"os"

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

func globalTimeTimezone(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("timezone", 1, len(args))
	}

	timezone, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("timezone", STRING_OBJ, 0, args)
	}

	err := timer.SetTimezone(timezone.Value)
	if err != nil {
		return nil, NewErrorf("timezone", "%s", err.Error())
	}

	return NULL, nil
}

func globalTimeDelayTimer(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("delayTimer", 2, len(args))
	}

	callable, ok := args[0].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypeError("delayTimer", FUNCTION_OBJ, 0, args)
	}

	if callable.ParametersCount() != 0 {
		return nil, NewErrorf("delayTimer", "function passed to `delayTimer` must take zero arguments")
	}

	delayTime, ok := args[1].(*Integer)
	if !ok {
		return nil, NewInvalidArgumentTypeError("delayTimer", INTEGER_OBJ, 1, args)
	}

	if delayTime.Value < 0 {
		return nil, NewErrorf("delayTimer", "delay time must be non-negative")
	}

	time := timer.StartDelayedTimer(func() {
		if rs := callable.Call(); IsError(rs) {
			fmt.Fprintf(os.Stdout, "%s\n", rs.Inspect())
		}
	}, delayTime.Value)

	return BuildImmutableHash(
		HashPair{
			Key: &String{Value: "stop"},
			Value: &Builtin{Fn: func(args ...Object) (Object, error) {
				if timer.StopDelayedTimer(time) {
					return TRUE, nil
				}

				return FALSE, nil
			}},
		},
		HashPair{
			Key:   &String{Value: "timer"},
			Value: &String{Value: fmt.Sprintf("%p", time)},
		},
	), nil
}

func globalTimeScheduleTimer(args ...Object) (Object, error) {
	if len(args) != 2 {
		return nil, NewWrongNumberOfArgumentsError("scheduleTimer", 2, len(args))
	}

	callable, ok := args[0].(Callable)
	if !ok {
		return nil, NewInvalidArgumentTypeError("scheduleTimer", FUNCTION_OBJ, 0, args)
	}

	if callable.ParametersCount() != 0 {
		return nil, NewErrorf("scheduleTimer", "function passed to `scheduleTimer` must take zero arguments")
	}

	intervalTime, ok := args[1].(*Integer)
	if !ok {
		return nil, NewInvalidArgumentTypeError("scheduleTimer", INTEGER_OBJ, 1, args)
	}

	if intervalTime.Value < 0 {
		return nil, NewErrorf("scheduleTimer", "interval time must be non-negative")
	}

	ticker := timer.StartScheduledTimer(func() {
		if rs := callable.Call(); IsError(rs) {
			fmt.Fprintf(os.Stdout, "%s\n", rs.Inspect())
		}
	}, intervalTime.Value)

	return BuildImmutableHash(
		HashPair{
			Key: &String{Value: "stop"},
			Value: &Builtin{Fn: func(args ...Object) (Object, error) {
				timer.StopScheduledTimer(ticker)
				return TRUE, nil
			}},
		},
		HashPair{
			Key:   &String{Value: "timer"},
			Value: &String{Value: fmt.Sprintf("%p", ticker)},
		},
	), nil
}
