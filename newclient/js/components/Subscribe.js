var React = require('react');

var Subscribe = React.createClass({
	handleClick:
	function() {
		this.props.onClick(this.props.data);
	},

        render: 
	function() {
		var cls = "";
		if (this.props.selected) {
			cls=" selected"
		}

		return <div onClick={this.handleClick} className={"subscribe"+cls}>
			<div><img className="banner" src={"http://192.168.0.110:9527/subscribe/banner/"+this.props.data.Name} /></div>
		</div>;
	}
});

module.exports = Subscribe;
