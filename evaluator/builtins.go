package evaluator

import (
	"charm/object"
	"fmt"
)

var builtins = map[string]*object.Builtin {
    "len": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }

            switch arg := args[0].(type) {
            case *object.String:
                return &object.Integer{Value: int64(len(arg.Value))}
            case *object.Array:
                return &object.Integer{Value: int64(len(arg.Elements))}
            default:
                return newError("argument to `len` not supported, got %s", args[0].Type())
            }
        },
    },
    "push": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return newError("wrong number of arguments. got=%d, want=2", len(args))
            }
            if args[0].Type() != object.ARRAY_OBJ {
                return newError("argument to `push` must be ARRAY, got %s",
                    args[0].Type())
            }
            arrayObj := args[0].(*object.Array)
            newArrayElements := append(arrayObj.Elements, args[1])
            
            return &object.Array{Elements: newArrayElements}
        },
    },
    "pop": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }
            if args[0].Type() != object.ARRAY_OBJ {
                return newError("argument to `pop` must be ARRAY, got %s",
                    args[0].Type())
            }
            arrayObj := args[0].(*object.Array)
            arrayLen := len(arrayObj.Elements)

            if arrayLen == 0 {
                return NULL
            }

            lastELement := arrayObj.Elements[arrayLen - 1]
            arrayObj.Elements = arrayObj.Elements[0:arrayLen - 1]

            return lastELement
        },
    },
    "keys": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 1 {
                return newError("wrong number of arguments. got=%d, want=1", len(args))
            }
            if args[0].Type() != object.HASHMAP_OBJ {
                return newError("argument to `keys` must be HASHMAP, got %s",
                    args[0].Type())
            }

            hashMapObj := args[0].(*object.HashMap)
            keys := &object.Array{Elements: []object.Object{}}

            for _, pair := range hashMapObj.Map {
                keys.Elements = append(keys.Elements, pair.Key)
            }

            return keys
        },
    },
    "delete": {
        Fn: func(args ...object.Object) object.Object {
            if len(args) != 2 {
                return newError("wrong number of arguments. got=%d, want=2", len(args))
            }

            if args[0].Type() != object.HASHMAP_OBJ {
                return newError("argument to `delete` must be HASHMAP, got %s",
                    args[0].Type())
            }

            hashMapObj := args[0].(*object.HashMap)

            hashable, ok := args[1].(object.Hashable)
            if !ok {
                return newError("unusable as a hashkey: %s", args[1].Type())
            }

            delete(hashMapObj.Map, hashable.HashCode())

            return NULL
        },
    },
    "print": {
        Fn: func(args ...object.Object) object.Object {
            for _, arg := range args {
                fmt.Println(arg.Inspect())
            }

            return NULL
        },
    },
}
