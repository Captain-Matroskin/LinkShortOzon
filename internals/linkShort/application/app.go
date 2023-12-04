package application

import "linkShortOzon/internals/linkShort/orm"

type LinkShortAppInterface interface {
	CreateLinkShortApp(linkFull string) (string, error)
	TakeLinkFullApp(linkShort string) (string, error)
}

type LinkShortApp struct {
	Wrapper orm.LinkShortWrapperInterface
}

func (l *LinkShortApp) CreateLinkShortApp(linkFull string) (string, error) {
	return "", nil
}

func (l *LinkShortApp) TakeLinkFullApp(linkShort string) (string, error) {
	return "", nil
}
