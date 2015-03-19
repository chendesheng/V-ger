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
		var props = this.props;
		var rows = props.subscribes.map(function(subscribe) {
			return <Subscribe key={subscribe.Name} onClick={props.handleSelectSubscribe} 
				selected={subscribe===props.selectedSubscribe} data={subscribe} />
		});
		return <div className="subscribes">
			<div className="subscribe-list noscrollbar">{rows}</div>
			<span className="peek peek-left"></span>
			<span className="peek peek-right"></span>
		</div>
	}
});

module.exports = SubscribeList
