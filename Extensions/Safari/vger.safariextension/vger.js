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
	} else if (event.name == 'verifycode') {
                console.log('show verify code:'+event.message);

                if (!event.message) {
                        var ele = document.getElementById('vger-verifycodebox');
                        if (ele) document.body.removeChild(ele);
                        return;
                }

                var ele = document.getElementById('vger-verifycodebox');
                if (ele) {
                        var img = document.getElementById('vger-verifycodeimg');
                        if (img) {
                                img.src += 'r';
                                this.style.disabled = '';
                        }
                } else {
                        ele = document.createElement('div');
                        ele.id = 'vger-verifycodebox';
                        ele.style.cssText = 'box-shadow: 0 4px 23px 5px rgba(0, 0, 0, 0.2), 0 2px 6px rgba(0,0,0,0.15);background-color:rgb(240,239,233);padding:10px 0;border-radius:4px;display:block;position:fixed;width:150px;left:50%;top:50%;margin-left:-75px;text-align:center;border:none;z-index:2147483639';
                        ele.innerHTML = '<div style="">Input Verify Code</div><div><img style="cursor:pointer;" title="click to refresh" id="vger-verifycodeimg" src="'+event.message+'?"></div><div><input style="text-align:center" id="vger-verifycode" type="text"/></div><div><input id="vger-verifycodebtn" onclick="vgerVerifyCodeClicked();" type="button" value="OK"/></div>';
                        document.body.insertBefore(ele, document.body.firstChild);
                        document.getElementById('vger-verifycodebtn').onclick = function() {
                                console.log('send verifycode:' + code);

                                var code = document.getElementById('vger-verifycode').value;
                                safari.self.tab.dispatchMessage('verifycode', code);
                                this.style.disabled = 'disabled';
                        };
                        document.getElementById('vger-verifycodeimg').onclick = function() {
                                this.src += 'r';
                        };
                }
                document.getElementById('vger-verifycode').focus();
        }
}

safari.self.addEventListener("message", handleMessage, false);

