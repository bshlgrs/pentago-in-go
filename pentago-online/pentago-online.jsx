// Define a collection to hold our tasks
Games = new Mongo.Collection("games");

if (Meteor.isClient) {
  Accounts.ui.config({
    passwordSignupFields: "USERNAME_ONLY"
  });

  Meteor.startup(function () {
    var mainSound = new buzz.sound('/sounds/main.m4a');
    var otherSound = new buzz.sound('/sounds/other.m4a');


    Games.find({}).observe({
      changed: function (newGame, oldGame) {
        if (Meteor.userId()) {
          if (newGame.players.map((x) => {return x._id}).indexOf(Meteor.userId()) != -1) {
            if (oldGame.players[oldGame.currentTurn]._id != newGame.players[newGame.currentTurn]._id) {
              mainSound.play();
            }
          }
        }
      }
    })
  });

  Meteor.subscribe("games");
}

if (Meteor.isServer) {
  Meteor.startup(function () {

  });

  // Meteor.publish("new-games", function () {
  //   return Games.find({state: "getting-players"});
  // });

  // Meteor.publish("playing-games", function () {
  //   return Games.find({state: "playing"});
  // });

  // Meteor.publish("finished-games", function () {
  //   return Games.find({state: "finished"});
  // });

  Meteor.publish("games", function () {
    return Games.find({});
  });
}

Meteor.methods({
  createGame(name, numberOfPlayers) {
    // Make sure the user is logged in before inserting a task
    if (! Meteor.userId()) {
      throw new Meteor.Error("not-authorized");
    }

    var gameId = Games.insert({
      name: name,
      numberOfPlayers: numberOfPlayers,
      state: "getting-players",
      players: [Meteor.user()],
      createdAt: new Date(),
      stateHistory: [],
      moveHistory: [],
      creator: Meteor.user()
    });

    return gameId;
  },

  joinGame(gameId) {
    var game = Games.findOne({_id: gameId});

    if (game.players.length == game.numberOfPlayers - 1) {
      var shuffledPlayers = shuffle(game.players.concat([Meteor.user()]));

      var board = [];

      for (var y = 0; y < 9; y++) {
        var row = [];
        for (var x = 0; x < 9; x++) {
          row.push(0);
        };
        board.push(row);
      };

      Games.update(game._id, {
        $set: {state: "playing", players: shuffledPlayers, board: board, currentTurn: 0, placingPiece: true }
      });
    } else {
      Games.update(gameId, {
        $set: {players: game.players.concat([Meteor.user()])}
      });
    }
  },

  playPiece(gameId, x, y) {
    var game = Games.findOne({_id: gameId});
    var board = game.board;
    var playerNumber = game.players.map((x) => x._id).indexOf(Meteor.userId());

    if (!game.placingPiece || playerNumber == -1 || playerNumber != game.currentTurn) {
      throw new Meteor.Error("not-authorized");
    }

    board[y][x] = playerNumber + 1;

    Games.update(game._id, {
      $set: {
        board: board,
        placingPiece: false
      }
    });
  },

  playPiece(gameId, x, y) {
    var game = Games.findOne({_id: gameId});
    var board = game.board;
    var playerNumber = game.players.map((x) => x._id).indexOf(Meteor.userId());

    if (!game.placingPiece || playerNumber == -1 || playerNumber != game.currentTurn) {
      throw new Meteor.Error("not-authorized");
    }

    if (board[y][x] != 0) {
      throw new Meteor.Error("bad-request");
    }

    board[y][x] = playerNumber + 1;

    if (playerHasWon(board)) {
      winGame(game, board);
    } else {
      Games.update(game._id, {
        $set: {
          board: board,
          placingPiece: false
        }
      });
    }
  },

  playRotation(gameId, rotateX, rotateY, clockwise) {
    var game = Games.findOne({_id: gameId});
    var board = game.board;
    var playerNumber = game.players.map((x) => x._id).indexOf(Meteor.userId());

    // if (game.placingPiece || playerNumber == -1 || playerNumber != game.currentTurn) {
    //   throw new Meteor.Error("not-authorized");
    // }

    // shitty deep copy
    var newBoard = JSON.parse(JSON.stringify(game.board));

    var startX = rotateX * 3;
    var startY = rotateY * 3;

    for (var x = 0; x < 3; x++) {
      for (var y = 0; y < 3; y++) {
        if (clockwise) {
          newBoard[x + startY][3 - y - 1 + startX] = board[startY + y][startX + x];
        } else {
          newBoard[3 - x - 1 + startY][y + startX] = board[startY + y][startX + x];
        }
      }
    }

    if (playerHasWon(newBoard)) {
      winGame(game, newBoard);
    } else {
      Games.update(game._id, {
        $set: {
          board: newBoard,
          placingPiece: true,
          currentTurn: (game.currentTurn + 1) % game.numberOfPlayers
        }
      });
    }
  }
});

function winGame(game, board) {
  Games.update(game._id, {
    $set: {
      board: board,
      state: "finished",
      winner: Meteor.user(),
      currentTurn: -1
    }
  });
}

function playerHasWon(board) {
  var directions = [[0, 1], [1, 0], [1, -1], [1, 1]];
  for (var x = 0; x < 9; x++) {
    for (var y = 0; y < 9; y++) {
      for (var i = 0; i < directions.length; i++) {
        var direction = directions[i];

        var length = 0;
        var player = board[y][x];

        if (player != 0) {
          var currX = x;
          var currY = y;

          while (currX < 9 && currX >= 0 && currY < 9 && currY >= 0 && board[currY][currX] == player) {
            length += 1;
            currX += direction[0];
            currY += direction[1];
            if (length == 5) {
              return true;
            }
          }
        }
      };
    }
  };

  return false;
}

function shuffle(array) {
  var currentIndex = array.length, temporaryValue, randomIndex;

  // While there remain elements to shuffle...
  while (0 !== currentIndex) {

    // Pick a remaining element...
    randomIndex = Math.floor(Math.random() * currentIndex);
    currentIndex -= 1;

    // And swap it with the current element.
    temporaryValue = array[currentIndex];
    array[currentIndex] = array[randomIndex];
    array[randomIndex] = temporaryValue;
  }

  return array;
}
