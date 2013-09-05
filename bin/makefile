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
window:
	go install client
	cp client ./GoCocoa.app/Contents/MacOS/GoCocoa