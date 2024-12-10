package helpers

import "errors"

func NumberComparison(givenValue float64, targetValue float64, operator string) (bool, error) {
	switch operator {
	case "<":
		return givenValue < targetValue, nil
	case "<=":
		return givenValue <= targetValue, nil
	case "=":
		return givenValue == targetValue, nil
	case ">=":
		return givenValue >= targetValue, nil
	case ">":
		return givenValue > targetValue, nil
	default:
		return false, errors.New("invalid operator")
	}

}
