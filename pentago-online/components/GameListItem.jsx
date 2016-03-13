
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
