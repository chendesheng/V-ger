all: *.m
	clang *.m -c -I../include
	libtool -static -o libcocoa.a *.o
	go install -a player/gui
run:
	../../../../bin/playerguitest