var React = require('react');
var Network = require('../utils/Network.js');

var Task = React.createClass({
	render:
	function() {
		var status = this.props.data.Status;
		var w = Math.floor((1-this.props.data.LastPlaying)*30) || 0;
		var cls = "";
		var onclick = this.play;

		if (status == 'Finished') {
			onclick = this.open;
		}

		if (status == 'New') {
			cls = 'disabled';
			w = 0;
			onclick = null;
		}

		return <li onClick={onclick} className={cls} title={this.props.data.Name}>
			{this.props.data.Episode<10?"0"+this.props.data.Episode:this.props.data.Episode}
			<span style={{width:w+"px"}} className="progress"></span>
		</li>
	},

	play:
	function() {
		Network.getJson('/play/'+this.props.data.Name);
	},

	open:
	function() {
		Network.getJson('/open/'+this.props.data.Name);
	}
});

module.exports = Task;
