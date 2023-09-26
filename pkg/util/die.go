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
