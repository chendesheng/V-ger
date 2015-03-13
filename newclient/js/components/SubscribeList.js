var React = require('react');
var Subscribe = require('./Subscribe.js');

var SubscribeList = React.createClass({
        handleSelectSubscribe:
	function(subscribe) {
		this.props.handleSelectSubscribe(subscribe);
	},

        render:
	function() {
		var self  = this;
		var rows = this.props.subscribes.map(function(subscribe) {
			return <Subscribe key={subscribe.Name} onClick={self.handleSelectSubscribe} 
				selected={subscribe.Name==self.props.selectedSubscribe} data={subscribe} />;
		});
		return <div className="subscribes">
			<div className="title-bar">
				<span>TV Shows</span>
				<span className="sep"></span>
				<span className="selected-name">{this.props.selectedSubscribe}</span>
			</div>
			<div className="subscribe-list noscrollbar">{rows}</div>
		</div>
	}
});

module.exports = SubscribeList
