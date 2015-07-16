.PHONY: run build

default: run

run: build
	./grimmnight_commander || reset

build: gc.go
	go build

src: dependencies.txt
	bash -c 'while read -r dep || [[ -n $$dep ]]; do go get $$dep; done < dependencies.txt'
