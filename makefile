CC=clang
BIN=bin
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
	$(BIN)/vger -debug -config=config.json