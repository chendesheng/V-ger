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
			selectedSubscribe: '',
			subscribes: [],
			tasks: []
		};
	},

	_onTaskUpdate:
	function(tasks) {
		var subscribes = SubscribeStore.getAllSubscribes();
		var selectedSubscribe = this.state.selectedSubscribe || '';

		if (selectedSubscribe == '' && subscribes.length > 0) {
			selectedSubscribe = subscribes[0].Name;
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
		var selectedSubscribe = this.state.selectedSubscribe || '';

		if (selectedSubscribe == '' && subscribes.length > 0) {
			selectedSubscribe = subscribes[0].Name;
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
			selectedSubscribe: subscribe.Name,
			subscribes: SubscribeStore.getAllSubscribes(),
			tasks: TaskStore.getAllTasks()
		});
	},

	render:
	function() {
		var state = this.state;
		var subscribeInfo = null;
		
		var subscribes = this.state.subscribes;
		for (var i = 0; i < subscribes.length; i++) {
			var s = subscribes[i];
			if (s.Name == state.selectedSubscribe) {
				subscribeInfo = <SubscribeInfo data={s} />;
			}
		}

		return <div>
			<SubscribeList selectedSubscribe={this.state.selectedSubscribe} handleSelectSubscribe={this.handleSelectSubscribe} subscribes={this.state.subscribes} />
			{subscribeInfo}
			<TaskList filter={this.state.selectedSubscribe} tasks={this.state.tasks}/>
		</div>
	}
});


module.exports = VgerApp;
