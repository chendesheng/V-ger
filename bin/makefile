exe:
	go install vger
	cp vger ~/Library/Vger/vger
	pkill vger
website: main.html ./assets/main.js ./assets/style.css
	cp main.html ~/Library/Vger/main.html
	cp ./assets/main.js ~/Library/Vger/assets/main.js
	cp ./assets/style.css ~/Library/Vger/assets/style.css
all:
	go install vger
	cp vger ~/Library/Vger/vger
	pkill vger
	cp main.html ~/Library/Vger/main.html
	cp ./assets/main.js ~/Library/Vger/assets/main.js
	cp ./assets/style.css ~/Library/Vger/assets/style.css
vp:
	go install player
	cp ./player ./VgerPlayer.app/Contents/MacOS/VgerPlayer
	cp ./player.plist ./VgerPlayer.app/Contents/Info.plist

player:
	clang *.m -c -I../include
	libtool -static -o libcocoa.a *.o
	go install -a player/gui