var React = require('react');
var Task = require('./Task.js');
var PureRenderMixin = require('react/addons').addons.PureRenderMixin;

var TaskList = React.createClass({
	mixins: [PureRenderMixin],

	render:
	function() {
		var props = this.props;
		var lists = [];
		for (var i = 0; i < 30; i++) {
			lists.push([]);
		}

		this.props.tasks.forEach(function(task) {
			if (task.Season <= 0) {
				lists[0].push(task);
			} else {
				lists[task.Season].push(task);
			}
		});
		var rows = [];
		for (var i = lists.length - 1; i >= 0; i--) {
			rows.push(this.getOneSeason(lists, i));
		}
		return <div className="task-list">{rows}</div>
	},

	getOneSeason: function(lists, i) {
		var l = lists[i].length;
		if (l == 0) return;

		var items = [];
		for (var j = 0; j < l; j++) {
			var t = lists[i][j];
			items.push(<Task data={t} />)
		}

			
		return <div>
				<div className="season-title">{"Season "+i}</div>
				<ul className="task-season noscrollbar noselect">
					{items}
				</ul>
			</div>
	}
});

module.exports = TaskList;
