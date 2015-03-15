var React = require('react');

var SubscribeInfo = React.createClass({
	render:
	function() {
		if (this.props.data == null) {
			return <div></div>
		}

		return <div className="subscribe-info">
			<div className="name">{this.props.data.Name}</div>
			<div><a target="_blank" className="source" href={this.props.data.URL}>{this.props.data.URL}</a></div>

			<ul style={{display:"none"}} className="attributes">
				<li>
					<span className="key">Network:</span>
					<span>CW</span>
				</li>
				<li>
					<span className="key">First episode date:</span>
					<span>April 15, 2014</span>
				</li>
				<li>
					<span className="key">Last played:</span>
					<span>The.Vampire.Diaries.S05E03.720p.HDTV.X264-DIMENSION.mkv</span>
				</li>
			</ul>
		</div>
	}
});

module.exports = SubscribeInfo;
