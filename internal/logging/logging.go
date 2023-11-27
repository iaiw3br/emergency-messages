package logging

import "github.com/sirupsen/logrus"

type Logger struct {
	*logrus.Logger
}

func New() Logger {
	return Logger{
		logrus.New(),
	}
}
