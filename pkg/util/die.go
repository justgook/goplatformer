package util

import "log/slog"

func OrDie(err error) {
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}
}

func GetOrDie[T any](arg T, err error) T {
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return arg
}
func Get2OrDie[T1 any, T2 any](arg1 T1, arg2 T2, err error) (T1, T2) {
	if err != nil {
		slog.Error(err.Error())
		panic(err)
	}

	return arg1, arg2
}
