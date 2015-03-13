var assign = require('object-assign');
var EventEmitter = require('events').EventEmitter;
var Network = require('../utils/network.js')
var UPDATE_EVENT = 'update';

var _subscribes = [];

var SubscribeStore = assign({}, EventEmitter.prototype, {
	fetch:
	function() {
		var self = this;
		Network.getJson('subscribe', function(data) {
			_subscribes = data;
			self.emitChange();
		});
	},

	emitChange:
	function() {
		this.emit(UPDATE_EVENT);
	},
	
	addUpdateListener:
	function(callback) {
		this.on(UPDATE_EVENT, callback);
	},

	getAllSubscribes:
	function() {
		return _subscribes;
	}
});

module.exports = SubscribeStore;
