Pentago = React.createClass({
  getInitialState () {
    return {
      currentX: null,
      currentY: null,
      choosingDirection: false
    }
  },
  componentWillMount () {
    this.squareSpacing = 60;
    this.colors = ["gray", "red", "blue", "green", "yellow"];
  },
  currentPlayerId () {
    if (this.props.game.state == "playing") {
      return this.props.game.players[this.props.game.currentTurn]._id;
    }
  },
  isMyTurn () {
    return this.props.game.state == "playing" && Meteor.userId() == this.currentPlayerId();
  },
  componentDidMount () {
    var canvas = this.refs.canvas;
    var that = this;

    var squareSpacing = this.squareSpacing;
    $(canvas).on("mousemove", function (e) {
      var rect = canvas.getBoundingClientRect();
      var x = e.clientX - rect.left;
      var y = e.clientY - rect.top;
      that.setState({
        currentX: (x / squareSpacing | 0),
        currentY: (y / squareSpacing | 0),
      });
    });

    $(canvas).on("click", function (e) {
      that.handleCanvasClick(e);
    })

    this.renderCanvas();
  },
  handleCanvasClick (e) {
    if (this.isMyTurn()) {
      if (this.props.game.placingPiece) {
        if (this.props.game.board[this.state.currentY][this.state.currentX] == 0) {
          Meteor.call("playPiece", this.props.game._id, this.state.currentX, this.state.currentY);
        }
      }
    }
  },
  handleCounterclockwiseRotate () {
    this.handleRotate(false);
  },
  handleClockwiseRotate () {
    this.handleRotate(true);
  },
  handleRotate (clockwise) {
    Meteor.call("playRotation", this.props.game._id, this.miniSquarePosition().x, this.miniSquarePosition().y, clockwise);
  },
  componentDidUpdate () {
    this.renderCanvas();
  },
  miniSquarePosition() {
    var miniSquareSize = this.props.miniSquareSize;
    return {
      x: (this.state.currentX / miniSquareSize) | 0,
      y: (this.state.currentY / miniSquareSize) | 0
    };
  },
  renderCanvas () {
    var canvas = this.refs.canvas;
    var ctx = canvas.getContext("2d");

    var colors = this.colors;

    var game = this.props.game;

    var size = this.props.size;
    var miniSquareSize = this.props.miniSquareSize;
    var squareSpacing = this.squareSpacing;

    ctx.fillStyle = "gray";

    for (var x = 0; x < size / miniSquareSize; x++) {
      for (var y = 0; y < size / miniSquareSize; y++) {
        ctx.beginPath();
        ctx.rect((x + 0.04) * squareSpacing * miniSquareSize,
                 (y + 0.04) * squareSpacing * miniSquareSize,
                 squareSpacing * miniSquareSize * 0.92,
                 squareSpacing * miniSquareSize * 0.92);

        if (x == this.miniSquarePosition().x &&
            y == this.miniSquarePosition().y &&
            this.isMyTurn() &&
            !game.placingPiece) {
          ctx.fillStyle = "darkred";
          ctx.fill();
          ctx.fillStyle = "gray";
        } else {
          ctx.fill();
        }

        ctx.stroke();
      }
    }

    ctx.strokeStyle = "black";

    for (var x = 0; x < size; x++) {
      for (var y = 0; y < size; y++) {
        var pixelX = squareSpacing * (x + 0.5);
        var pixelY = squareSpacing * (y + 0.5);

        ctx.fillStyle = colors[game.board[y][x]];
        ctx.beginPath();
        ctx.arc(pixelX, pixelY, this.squareSpacing * 0.35, 0, 2*Math.PI);
        ctx.fill();
        if (x == this.state.currentX && y == this.state.currentY && this.isMyTurn() && game.placingPiece) {
          ctx.strokeStyle = "red";
          ctx.stroke();
          ctx.strokeStyle = "black";
        } else {
          ctx.stroke();
        }
      };
    };
  },
  render () {
    var game = this.props.game;
    return (
      <div>
        <div className="panel panel-default">
          <div className="panel-body">
            {game.winner ?
            <strong>{game.winner.username} has won the game!</strong> :
            <strong>{game.players[game.currentTurn].username} is about to
              {game.placingPiece ? " place a piece": " do a rotation"}
            </strong>}
          </div>
        </div>
        <div>
          {this.isMyTurn() && !game.placingPiece &&
            <div>
              <button
                className="btn btn-default btn-sml"
                style={{
                  left: this.miniSquarePosition().x * 60 * 3 + 10,
                  marginTop: this.miniSquarePosition().y * 60 * 3,
                  position: "absolute"}}
                onClick={this.handleCounterclockwiseRotate}>
                <i className="fa fa-rotate-left" />
              </button>
              <button
                className="btn btn-default btn-sml"
                style={{
                  left: (this.miniSquarePosition().x + 1) * 60 * 3 - 20,
                  marginTop: this.miniSquarePosition().y * 60 * 3,
                  position: "absolute"}}
                onClick={this.handleClockwiseRotate}>
                <i className="fa fa-rotate-right" />
              </button>
            </div>
          }
          <canvas
            width={this.props.size * this.squareSpacing}
            height={this.props.size * this.squareSpacing}
            ref="canvas" />
        </div>
        <p>Players:
          {game.players.map((x, i) => {
            return <span key={x.username} style={{color: this.colors[i + 1]}}>{x.username}&nbsp;</span>;
          })}
        </p>
      </div>
    );
  }
})
