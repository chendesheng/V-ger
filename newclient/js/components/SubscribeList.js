var React = require('react');
var Subscribe = require('./Subscribe.js');
var PureRenderMixin = require('react/addons').addons.PureRenderMixin;

var SubscribeList = React.createClass({
	mixins: [PureRenderMixin],

        handleSelectSubscribe:
	function(subscribe) {
		this.props.handleSelectSubscribe(subscribe);
	},

        render:
	function() {
		var self  = this;
		var selectedName = null;
		if (this.props.selectedSubscribe != null) {
			selectedName = this.props.selectedSubscribe.Name;
		}

		var rows = this.props.subscribes.map(function(subscribe) {
			return <Subscribe key={subscribe.Name} onClick={self.handleSelectSubscribe} 
				selected={subscribe.Name==selectedName} data={subscribe} />
		});
		return <div className="subscribes">
			<div className="subscribe-list noscrollbar">{rows}</div>
		</div>
	}
});

module.exports = SubscribeList
