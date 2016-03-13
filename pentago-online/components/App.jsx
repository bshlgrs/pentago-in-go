App = React.createClass({
  mixins: [ReactMeteorData],

  getMeteorData() {
    return {
      games: Games.find({}).fetch(),
      gamesGettingPlayers: Games.find({ state: "getting-players"}).fetch(),
      gamesPlaying: Games.find({ state: "playing"}).fetch(),
      gamesFinished: Games.find({ state: "finished"}).fetch(),
      user: Meteor.user(),
      userId: Meteor.userId()
    }
  },

  handleGameSelect (_id) {
    FlowRouter.go('/games/' + _id);
  },

  goToMainMenu () {
    FlowRouter.go('/');
  },

  handleCreateGame (event) {
    event.preventDefault();

    // Find the text field via the React ref
    var name = ReactDOM.findDOMNode(this.refs.nameInput).value.trim();
    var numberOfPlayers = parseInt(document.querySelector('input[name="playerNumber"]:checked').value);
    var that = this;

    Meteor.call("createGame", name, numberOfPlayers, function (err, gameId) {
      FlowRouter.go('/games/' + gameId);
    });
  },

  render() {
    // Get tasks from this.data.tasks

    var inner = <p></p>;

    if (FlowRouter.getRouteName() == "root") {
      inner = (
        <div>
          <div className="panel panel-default">
            <div className="panel-body">
              <h3>Join a game!</h3>
              {this.data.games.length ?
                <ul>
                  {this.data.gamesGettingPlayers.map((game) => {
                    return <GameListItem handleGameSelect={this.handleGameSelect} key={game._id} game={game} />;
                  })}
                </ul> :
                <p>No games currently exist, create one?</p>
              }
            </div>
          </div>

          <div className="panel panel-default">
            <div className="panel-body">
              <h3>Currently live matches</h3>
              {this.data.games.length ?
                <ul>
                  {this.data.gamesPlaying.map((game) => {
                    return <GameListItem handleGameSelect={this.handleGameSelect} key={game._id} game={game} />;
                  })}
                </ul> :
                <p>No-one is playing right now :(</p>
              }
            </div>
          </div>

          {this.data.user &&
          <div className="panel panel-default">
            <div className="panel-body">
              <h3>Make new game!</h3>
              <form onSubmit={this.handleCreateGame} >
                <input
                  className="form-control"
                  type="text"
                  ref="nameInput"
                  placeholder="Name of game"
                  defaultValue={this.data.user.username + "'s cool game"}/>
                <div className="radio">
                  {[2, 3, 4].map((n) => {
                    return <label key={n} className="radio-inline">
                      <input type="radio" name="playerNumber" value={n} defaultChecked={2 == n}/> {n} players
                    </label>
                  })}
                </div>

                <button className="btn btn-primary" type="submit">Create game</button>
              </form>
            </div>
          </div>}



          <div className="panel panel-default">
            <div className="panel-body">
              <h3>Old matches</h3>
              {this.data.games.length ?
                <ul>
                  {this.data.gamesFinished.map((game) => {
                    return <GameListItem handleGameSelect={this.handleGameSelect} key={game._id} game={game} />;
                  })}
                </ul> :
                <p>No-one has played yet :(</p>
              }
            </div>
          </div>

        </div>
      );
    } else if (FlowRouter.getRouteName() == "view-game") {
      inner = <ShowGame userId={this.data.userId} game={Games.findOne({_id: FlowRouter.getParam("gameId")})} />
    }

    return <div>
      <div className="pull-right">
        <AccountsUIWrapper/>
      </div>
      <h1 onClick={this.goToMainMenu}><a>pentago online</a></h1>
      {!this.data.userId &&
        <div className="alert alert-info">Sign up to play games!</div>}
      {inner}
    </div>
  },
});

