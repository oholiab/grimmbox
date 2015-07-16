.PHONY: run build

default: run

run: build
	./grimmbox || reset

build: gc.go src
	go build

src: dependencies.txt
	bash -c 'while read -r dep || [[ -n $$dep ]]; do go get $$dep; done < dependencies.txt'
