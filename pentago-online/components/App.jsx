App = React.createClass({
  mixins: [ReactMeteorData],

  getInitialState () {
    return {
      state: this.props.initialState,
      currentGameId: this.props.initialGameId
    };
  },

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
    this.setState({ state: "view-game", currentGameId: _id});
  },

  goToMainMenu () {
    FlowRouter.go('/');
    this.setState({ state: "list" });
  },

  handleCreateGame (event) {
    event.preventDefault();

    // Find the text field via the React ref
    var name = ReactDOM.findDOMNode(this.refs.nameInput).value.trim();
    var numberOfPlayers = parseInt(document.querySelector('input[name="playerNumber"]:checked').value);
    var that = this;

    Meteor.call("createGame", name, numberOfPlayers, function (err, gameId) {
      that.setState({ state: "view-game", currentGameId: gameId});
    });
  },

  render() {
    // Get tasks from this.data.tasks

    var inner = <p></p>;

    if (this.state.state == "list") {
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
    } else if (this.state.state == "view-game") {
      inner = <ShowGame userId={this.data.userId} game={Games.findOne({_id: this.state.currentGameId})} />
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

GameListItem = React.createClass({
  handleClick () {
    this.props.handleGameSelect(this.props.game._id);
  },
  render () {
    var game = this.props.game;

    var inner;

    if (this.props.game.state == "getting-players") {
      inner = <span>
        {game.players.length ? <span>
          Players: {game.players.map((x) => x.username).join(",")}.
          </span> :
        <span>{game.numberOfPlayers} player game. </span>}

        &nbsp;Needs {game.numberOfPlayers - game.players.length} more players.
      </span>;
    } else if (game.state == "playing") {
      inner = <span>
        Players: {game.players.map((x) => x.username).join(",")}
        . Current turn: {game.players[game.currentTurn].username}.
      </span>;
    } else if (game.state == "finished") {
      inner = <span>{game.winner.username} won!</span>;
    }

    return <li>
      <a onClick={this.handleClick}>{game.name}</a> {inner}
    </li>;
  }
});

ShowGame = React.createClass({
  handleJoin () {
    Meteor.call("joinGame", this.props.game._id);
  },
  render () {
    var game = this.props.game;

    var inner = <p></p>;
    var that = this;

    if (game.state == "getting-players") {
      inner = <div>
        <p>Currently waiting on more players! {game.numberOfPlayers - game.players.length} needed.</p>

        {game.players.length ? <div>
          <p>Players:</p>
          <ul>
            {game.players.map((x) => {
              return <li key={x._id}>{x.username}</li>;
            })}
          </ul>
        </div> : null}

        {game.players.map((x) => x._id).indexOf(that.props.userId) == -1 ?
          <button className="btn btn-primary" onClick={this.handleJoin}>Join game!</button>
          : null}

      </div>;
    } else if (game.state == "playing") {
      inner = <div>
        <Pentago size={9} miniSquareSize={3} game={game}/>
      </div>;
    } else if (game.state == "finished") {
      inner = <div>
        <Pentago size={9} miniSquareSize={3} game={game}/>
      </div>;
    }

    return (
      <div>
        <h3>{this.props.game.name}</h3>

        {inner}
      </div>
    );
  }
})
