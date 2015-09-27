CC=clang
BIN=../../bin
RESOURCES =../../bin/VgerPlayer.app/Contents/Resources
APP=vgerapp
vp:
	go install vger/player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
	ibtool --compile $(RESOURCES)/MainMenu.nib player/gui/cocoa/MainMenu.xib
	cp $(RESOURCES)/MainMenu.nib $(RESOURCES)/en.lproj/MainMenu.nib
	ibtool --compile $(RESOURCES)/openURL.nib player/gui/cocoa/openURL.xib
	cp $(RESOURCES)/openURL.nib $(RESOURCES)/en.lproj/openURL.nib
exe:
	go install vger
	cp $(BIN)/vger ~/Library/Vger/vger
	pkill vger
website: $(APP)/index.html $(APP)/assets/main.js $(APP)/assets/style.css
	macgap build -n "V'ger" vgerapp
	cp $(APP)/index.html ~/Library/Vger/index.html
	cp $(APP)/assets/main.js ~/Library/Vger/assets/main.js
	cp $(APP)/assets/style.css ~/Library/Vger/assets/style.css
web:
	cp $(APP)/index.html ~/Library/Vger/index.html
	cp $(APP)/assets/main.js ~/Library/Vger/assets/main.js
	cp $(APP)/assets/style.css ~/Library/Vger/assets/style.css
vprace:
	go install -race vger/player
	cp $(BIN)/player $(BIN)/VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp $(BIN)/player.plist $(BIN)/VgerPlayer.app/Contents/Info.plist
runvger:
	go install vger
	$(BIN)/vger -debug -config=config.json
