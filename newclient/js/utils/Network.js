var $ = require('jquery');

module.exports = {
	startMonitor:
	function (path, ondata, onclose, onerror) {
		var websocket = new WebSocket('ws://192.168.0.111:9527/' + path);

		websocket.onopen = onOpen;
		websocket.onclose = onClose;
		websocket.onmessage = onMessage;
		websocket.onerror = onError;

		function onOpen(evt) {
		}

		function onClose(evt) {
			if (onclose) {
				onclose(evt);
			}
		}

		function onMessage(evt) {
			ondata(JSON.parse(evt.data));
		}

		function onError(evt) {
			if (onerror) {
				onerror(evt.data);
			}
		}

		function doSend(message) {
			websocket.send(message);
		}

		return websocket;
	},

	getJson: 
	function (path, onsuccess, onerror) {
		return $.ajax({
				url: path,
				dataType: 'json',
				success: onsuccess,
				error: onerror
			});
	}
}
