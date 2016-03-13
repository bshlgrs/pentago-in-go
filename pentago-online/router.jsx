FlowRouter.route('/games/:gameId', {
  action: function(params, queryParams) {
    ReactDOM.render(
      <App initialState="view-game" initialGameId={params.gameId}/>,
      document.getElementById("render-target")
    );
  }
});

FlowRouter.route('/', {
  action: function(params, queryParams) {
    ReactDOM.render(
      <App initialState="list"/>,
      document.getElementById("render-target")
    );
  }
});
