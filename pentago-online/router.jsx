FlowRouter.route('/games/:gameId', {
  action: function(params, queryParams) {
    ReactDOM.render(
      <App/>,
      document.getElementById("render-target")
    );
  },
  name: 'view-game'
});

FlowRouter.route('/', {
  action: function(params, queryParams) {
    ReactDOM.render(
      <App/>,
      document.getElementById("render-target")
    );
  },
  name: 'root'
});
