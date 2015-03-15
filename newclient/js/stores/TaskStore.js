var assign = require('object-assign');
var EventEmitter = require('events').EventEmitter;
var Network = require('../utils/Network.js');
var Immutable = require('immutable');
var UPDATE_EVENT = 'update';

var _tasks = Immutable.Map();
var TaskStore = assign({}, EventEmitter.prototype, {
	startMonitor:
	function() {
		var self = this;
		Network.startMonitor('progress', function(data) {
			for (var i = 0; i < data.length; i++) {
				_tasks = _tasks.set(data[i].Name, data[i]);
			}

			_tasks = _tasks.sort(function(a, b) {
				if (a.Season > b.Season) {
					return -1;
				} else if (a.Season < b.Season) {
					return 1;
				} else if (a.Episode < b.Episode) {
					return -1;
				} else if (a.Episode > b.Episode) {
					return 1;
				} else {
					return 0;
				}
			});

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

	getAllTasks:
	function() {
		return _tasks;
	}
});

module.exports = TaskStore;
