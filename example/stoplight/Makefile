topdir=../../
statecraft=$(topdir)/statecraft

all: build dot go python

.PHONY: build dot go python

build:
	cd $(topdir) && go build 

dot:
	$(statecraft) stoplight.statecraft stoplight.dot
	$(statecraft) car.statecraft car.dot

go:
	$(statecraft) stoplight.statecraft go/stoplight/stoplight.go
	$(statecraft) car.statecraft go/car/car.go

python:
	$(statecraft) stoplight.statecraft python/stoplight.py
	$(statecraft) car.statecraft python/car.py

