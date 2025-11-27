package objects

import (
	"os"

	"github.com/senither/zen-lang/objects/process"
)

func globalProcessExit(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("exit", 1, len(args))
	}

	codeObj, ok := args[0].(*Integer)
	if !ok {
		return nil, NewInvalidArgumentTypeError("exit", INTEGER_OBJ, 0, args)
	}

	process.Exit(int(codeObj.Value))

	return NULL, nil
}

func globalProcessArgv(args ...Object) (Object, error) {
	if len(args) != 0 {
		return nil, NewWrongNumberOfArgumentsError("argv", 0, len(args))
	}

	argv := make([]Object, len(os.Args))
	for i, arg := range os.Args {
		argv[i] = &String{Value: arg}
	}

	return &Array{Elements: argv}, nil
}

func globalProcessEnv(args ...Object) (Object, error) {
	if len(args) != 1 {
		return nil, NewWrongNumberOfArgumentsError("env", 1, len(args))
	}

	key, ok := args[0].(*String)
	if !ok {
		return nil, NewInvalidArgumentTypeError("env", STRING_OBJ, 0, args)
	}

	value, exists := process.LookupEnv(key.Value)
	if !exists {
		return NULL, nil
	}

	return &String{Value: value}, nil
}
