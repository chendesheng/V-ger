CC=clang
BIN=bin
P_INC=src/player/gui/include
P_SRC=src/player/gui/cocoa
all:
	$(CC) $(P_SRC)/*.m -c -I$(P_INC)
	libtool -static -o $(P_SRC)/libcocoa.a *.o
	go install -a player/gui
	go install player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/config.json $(BIN)/VgerPlayer.app/Contents/MacOS/config.json
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
	rm *.o
gui:
	$(CC) $(P_SRC)/*.m -c -I$(P_INC)
	libtool -static -o $(P_SRC)/libcocoa.a *.o
	go install -a player/gui
	rm *.o

# exe:
# 	go install vger
# 	cp vger ~/Library/Vger/vger
# 	pkill vger
# website: main.html $(BIN)/assets/main.js $(BIN)/assets/style.css
# 	cp main.html ~/Library/Vger/main.html
# 	cp $(BIN)/assets/main.js ~/Library/Vger/assets/main.js
# 	cp $(BIN)/assets/style.css ~/Library/Vger/assets/style.css
# vger:
# 	go install vger
# 	cp vger ~/Library/Vger/vger
# 	pkill vger
# 	cp main.html ~/Library/Vger/main.html
# 	cp $(BIN)/assets/main.js ~/Library/Vger/assets/main.js
# 	cp $(BIN)/assets/style.css ~/Library/Vger/assets/style.css
# vp:
# 	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
# 	cp $(BIN)/config.json $(BIN)/VgerPlayer.app/Contents/MacOS/config.json
# 	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
