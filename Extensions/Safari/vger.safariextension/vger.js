document.addEventListener("contextmenu", handleContextMenu, false);

function handleContextMenu(event) {
	// safari.self.tab.setContextMenuEventUserInfo(event, event.target);
	if (event.target.tagName == 'A') {
		safari.self.tab.dispatchMessage("contextmenu", event.target.href);
	} else {
		var currentElement = event.target;
		while (currentElement != null) {
			if (currentElement.tagName == 'A') {
				break;
			}
			currentElement = currentElement.parentNode;
		}
		if (currentElement && currentElement.tagName == 'A') {
			safari.self.tab.dispatchMessage("contextmenu", currentElement.href);
		} else {
			safari.self.tab.dispatchMessage("contextmenu", "");
		}
	}

	// safari.self.tab.dispatchMessage('sendNsaArray');
}

function handleMessage(event) {
	if (event.name == 'alert') {
		alert(event.message);
	}
}

safari.self.addEventListener("message", handleMessage, false);

