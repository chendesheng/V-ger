CC=clang
BIN=bin
P_INC=src/player/gui/include
P_SRC=src/player/gui/cocoa
all: gui
	go install player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/config.json $(BIN)/VgerPlayer.app/Contents/MacOS/config.json
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
gui: $(P_SRC)
	$(CC) $(P_SRC)/*.m -c -I$(P_INC)
	libtool -static -o $(P_SRC)/libcocoa.a *.o
	go install player/gui
exe:
	go install vger
	cp $(BIN)/vger ~/Library/Vger/vger
	pkill vger
website: $(BIN)/main.html $(BIN)/assets/main.js $(BIN)/assets/style.css
	cp $(BIN)/main.html ~/Library/Vger/main.html
	cp $(BIN)/assets/main.js ~/Library/Vger/assets/main.js
	cp $(BIN)/assets/style.css ~/Library/Vger/assets/style.css
vp:
	go install player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
runvger:
	go install vger
	# cp $(BIN)/vger $(BIN)/vgerdebug
	$(BIN)/vger -debug -config=config.json
clean:
	rm *.o