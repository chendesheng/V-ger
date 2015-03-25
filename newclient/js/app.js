// This file bootstraps the entire application.

var React = require('react');
var VgerApp = require('./components/VgerApp.js');

window.React = React; // export for http://fb.me/react-devtools

React.render(
    <VgerApp />,
    document.getElementById('app')
);
