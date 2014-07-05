CC=clang
BIN=bin
APP=vgerapp
exe:
	go install vger
	cp $(BIN)/vger ~/Library/Vger/vger
	pkill vger
website: $(APP)/index.html $(APP)/assets/main.js $(APP)/assets/style.css
	cp $(APP)/index.html ~/Library/Vger/index.html
	cp $(APP)/assets/main.js ~/Library/Vger/assets/main.js
	cp $(APP)/assets/style.css ~/Library/Vger/assets/style.css
	macgap build -n "V'ger" vgerapp
vp:
	go install player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
runvger:
	go install vger
	$(BIN)/vger -debug -config=config.json