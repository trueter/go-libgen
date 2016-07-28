build :
	go build .

run :
	./go-libgen

#.PHONY: all
all : build run

info :
	echo "no info"
