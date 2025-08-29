package objects

import "math/big"

func NewInteger(value int64) *Integer {
	return &Integer{Value: big.NewInt(value)}
}

func NewIntegerFromString(value string) *Integer {
	i, ok := new(big.Int).SetString(value, 10)
	if !ok {
		return nil
	}

	return &Integer{Value: i}
}

func NewFloat(value float64) *Float {
	return &Float{
		Value: new(big.Float).
			SetPrec(FLOATING_PRECISION).
			SetFloat64(value),
	}
}

func NewFloatFromString(value string) *Float {
	f, _, err := big.ParseFloat(value, 10, FLOATING_PRECISION, big.ToNearestEven)
	if err != nil {
		return nil
	}

	return &Float{Value: f}
}

func NewFloatFromInt64(value int64) *Float {
	return &Float{
		Value: new(big.Float).
			SetPrec(FLOATING_PRECISION).
			SetInt64(value),
	}
}

func NewFloatFromBigInt(value *big.Int) *Float {
	return &Float{
		Value: new(big.Float).
			SetPrec(FLOATING_PRECISION).
			SetInt(value),
	}
}
