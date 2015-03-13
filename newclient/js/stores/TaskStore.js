var assign = require('object-assign');
var EventEmitter = require('events').EventEmitter;
var Network = require('../utils/Network.js')
var UPDATE_EVENT = 'update';
var Immutable = require('immutable');

//var _tasks = Immutable.fromJS([]);
var _tasks = [];

var TaskStore = assign({}, EventEmitter.prototype, {
	startMonitor:
	function() {
		var self = this;
		Network.startMonitor('progress', function(data) {
			//_tasks = _tasks.merge(Immutable.fromJS(data));
			_tasks = _tasks.concat(data);
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
		_tasks.sort(function(a, b) {
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
		return _tasks;
	}
});

module.exports = TaskStore;
