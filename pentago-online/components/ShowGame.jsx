
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
