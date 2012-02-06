include $(GOROOT)/src/Make.inc

TARG=github.com/ziutek/de
GOFILES=\
	agent.go\
	minimizer.go\

include $(GOROOT)/src/Make.pkg
