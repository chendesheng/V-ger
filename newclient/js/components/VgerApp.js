var React = require('react');
var SubscribeList = require('./SubscribeList.js');
var TaskList = require('./TaskList.js');
var SubscribeInfo = require('./SubscribeInfo.js');
var TaskStore = require('../stores/TaskStore.js');
var SubscribeStore = require('../stores/SubscribeStore.js');

var VgerApp = React.createClass({
	getInitialState:
	function() {
		return {
			selectedSubscribe: null, 
			subscribes: [],
			tasks: []
		};
	},

	_onTaskUpdate:
	function(tasks) {
		var subscribes = SubscribeStore.getAllSubscribes();
		var selectedSubscribe = this.state.selectedSubscribe;

		if (selectedSubscribe == null && subscribes.length > 0) {
			selectedSubscribe = subscribes[0];
		}

		this.setState({
			selectedSubscribe: selectedSubscribe,
			subscribes: subscribes,
			tasks: TaskStore.getAllTasks()
		});
	},

	_onSubscribeUpdate:
	function() {
		var subscribes = SubscribeStore.getAllSubscribes();
		var selectedSubscribe = this.state.selectedSubscribe;

		if (selectedSubscribe == null && subscribes.length > 0) {
			selectedSubscribe = subscribes[0];
		}

		this.setState({
			selectedSubscribe: selectedSubscribe,
			subscribes: subscribes,
			tasks: TaskStore.getAllTasks()
		});
	},

	componentDidMount:
	function() {
		TaskStore.addUpdateListener(this._onTaskUpdate);
		SubscribeStore.addUpdateListener(this._onSubscribeUpdate);

		TaskStore.startMonitor();
		SubscribeStore.fetch();
	},

	componentWillUnmount:
	function() {
		TaskStore.removeUpdateListener(this._onTaskUpdate);
		SubscribeStore.removeUpdateListener(this._onSubscribeUpdate);
	},


	handleSelectSubscribe:
	function(subscribe) {
		this.setState({
			selectedSubscribe: subscribe,
			subscribes: SubscribeStore.getAllSubscribes(),
			tasks: TaskStore.getAllTasks()
		});
	},

	render:
	function() {
		var state = this.state;
		var selectedTasks = state.tasks.filter(function(task){
			if (state.selectedSubscribe) {
				return task.Subscribe==state.selectedSubscribe.Name;
			} else {
				return false;
			}
		});

		return <div>
			<SubscribeList selectedSubscribe={state.selectedSubscribe} handleSelectSubscribe={this.handleSelectSubscribe} subscribes={state.subscribes} />
			<SubscribeInfo data={state.selectedSubscribe} />
			<TaskList tasks={selectedTasks}/>
		</div>
	}
});


module.exports = VgerApp;
