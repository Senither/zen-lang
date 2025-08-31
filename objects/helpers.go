package objects

func BuildImmutableHash(args ...HashPair) *ImmutableHash {
	pairs := make(map[HashKey]HashPair)

	for _, arg := range args {
		hash := arg.Key.(Hashable)
		pairs[hash.HashKey()] = arg
	}

	return &ImmutableHash{Value: Hash{Pairs: pairs}}
}

func WrapBuiltinFunctionInMap(name string, fn func(...Object) Object) HashPair {
	return HashPair{
		Key:   &String{Value: name},
		Value: &Builtin{Fn: fn},
	}
}
