package main
import "fmt"

const BoardSize byte = 9
const Penta byte = 5
const MiniSquareSize byte = 3
const NumberOfMiniSquares byte = BoardSize / MiniSquareSize


////////////////// PENTAGO BOARD

// I'm trying to make this analogous to mcts.py

// Clone is unnecessary.

type PentagoBoard struct {
  grid [BoardSize][BoardSize]byte
  playerJustMoved byte
  numberOfPlayers byte
  numberOfPlayedPieces byte
}


func (b PentagoBoard) Get (x, y byte) byte {
  return b.grid[y][x];
}

// equivalent to DoMove

func (b *PentagoBoard) DoMoveMut (move PentagoMove) {
  b.PlayPieceMut(move.x, move.y, b.NextToPlay());

  b.RotateMut(move.rotationX, move.rotationY, move.clockwise);

  b.playerJustMoved = b.NextToPlay();
}

func (b *PentagoBoard) PlayPieceMut (x, y, player byte) {
  b.grid[y][x] = player;
  b.numberOfPlayedPieces ++
}

func (b PentagoBoard) NextToPlay () byte {
  return b.playerJustMoved % 4 + 1
}

func (b PentagoBoard) PlayPiece (x, y, player byte) PentagoBoard {
  b.grid[y][x] = player;
  b.numberOfPlayedPieces ++
  return b;
}

func (b *PentagoBoard) RotateMut (x, y byte, clockwise bool) {
  // fail if the size isn't 3
  startX := x * MiniSquareSize;
  startY := y * MiniSquareSize;
  temp := b.grid[startY][startX]

  if clockwise {
    b.grid[startY][startX] = b.grid[startY + 2][startX]
    b.grid[startY + 2][startX] = b.grid[startY + 2][startX + 2]
    b.grid[startY + 2][startX + 2] = b.grid[startY][startX + 2]
    b.grid[startY][startX + 2] = temp

    temp = b.grid[startY + 1][startX]
    b.grid[startY + 1][startX] = b.grid[startY + 2][startX + 1]
    b.grid[startY + 2][startX + 1] = b.grid[startY + 1][startX + 2]
    b.grid[startY + 1][startX + 2] = b.grid[startY][startX + 1]
    b.grid[startY][startX + 1] = temp
  } else {
    b.RotateMut(x, y, true)
    b.RotateMut(x, y, true)
    b.RotateMut(x, y, true)
  }
}

func (board PentagoBoard) GetResult () float64 {
  winner := board.GetWinner()
  if winner == 0 {
    return 0
  } else if winner == board.playerJustMoved  {
    return 1
  } else {
    return -1
  }
}

func (b PentagoBoard) GetMoves() []PentagoMove {
  isAnyMiniSquareRotationallySymmetrical, emptySquareX, emptySquareY :=
    b.isAnyMiniSquareRotationallySymmetrical()

  rotations := []PentagoRotation{}

  if isAnyMiniSquareRotationallySymmetrical {
    rotations = append(rotations, PentagoRotation{emptySquareX, emptySquareY, true})
  }

  for x := byte(0); x < BoardSize / MiniSquareSize; x++ {
    for y := byte(0); x < BoardSize / MiniSquareSize; x++ {
      if isAnyMiniSquareRotationallySymmetrical {
        if ! b.isMiniSquareRotationallySymmetrical(x, y) {
          rotations = append(rotations, PentagoRotation{x, y, true}, PentagoRotation{x, y, false})
        }
      } else {
        rotations = append(rotations, PentagoRotation{x, y, true}, PentagoRotation{x, y, false})
      }
    }
  }

  moves := []PentagoMove{}

  for x := byte(0); x < BoardSize; x++ {
    for y := byte(0); y < BoardSize; y++ {
      // If a minisquare is rotationally symmetrical and the board has another
      // rotationally symmetrical minisquare,
      // there is no reason to move to this minisquare and then rotate it.

      if (b.Get(x, y) == 0) {
        for _, rotation := range(rotations) {
          moves = append(moves, PentagoMove{x, y, rotation.x, rotation.y, rotation.clockwise})
        }
      }
    }
  }

  return moves
}

func (b PentagoBoard) GetApproximateMoves() []PentagoMove {
  isAnyMiniSquareRotationallySymmetrical, emptySquareX, emptySquareY :=
    b.isAnyMiniSquareRotationallySymmetrical()

  rotations := []PentagoRotation{}

  if isAnyMiniSquareRotationallySymmetrical {
    rotations = append(rotations, PentagoRotation{emptySquareX, emptySquareY, true})
  }

  for x := byte(0); x < BoardSize / MiniSquareSize; x++ {
    for y := byte(0); x < BoardSize / MiniSquareSize; x++ {
      if isAnyMiniSquareRotationallySymmetrical {
        if ! b.isMiniSquareRotationallySymmetrical(x, y) {
          rotations = append(rotations, PentagoRotation{x, y, true}, PentagoRotation{x, y, false})
        }
      } else {
        rotations = append(rotations, PentagoRotation{x, y, true}, PentagoRotation{x, y, false})
      }
    }
  }

  moves := []PentagoMove{}

  for x := byte(0); x < BoardSize; x++ {
    for y := byte(0); y < BoardSize; y++ {
      // If a minisquare is rotationally symmetrical and the board has another
      // rotationally symmetrical minisquare,
      // there is no reason to move to this minisquare and then rotate it.

      if (b.Get(x, y) == 0) {
        for _, rotation := range(rotations) {
          moves = append(moves, PentagoMove{x, y, rotation.x, rotation.y, rotation.clockwise})
        }
      }
    }
  }

  return moves
}

func (board PentagoBoard) GetWinner () byte {
  if board.numberOfPlayedPieces < board.numberOfPlayers * (Penta - 1) + 1 {
    return 0
  }

  for y := byte(0); y < BoardSize; y++ {
    winner := board.checkLine(0, y, 1, 0)
    if winner != 0 {
      return winner;
    }
  }

  for x := byte(0); x < BoardSize; x++ {
    winner := board.checkLine(x, 0, 0, 1)
    if winner != 0 {
      return winner;
    }
  }

  for start := byte(0); start < BoardSize; start++ {
    winner := board.checkLine(start, 0, 1, 1)
    if winner != 0 {
      return winner;
    }

    winner = board.checkLine(start, 0, 1, -1)
    if winner != 0 {
      return winner;
    }

    winner = board.checkLine(0, start, 1, 1)
    if winner != 0 {
      return winner;
    }

    winner = board.checkLine(0, start, 1, -1)
    if winner != 0 {
      return winner;
    }
  }

  return 0;
}

// uninteresting

func (board PentagoBoard) PrintBoard () {
  fmt.Println("  012 345 678")
  fmt.Println(" +-----------+")
  for y := byte(0); y < BoardSize; y++ {
    fmt.Printf("%d|", y);
    for x := byte(0); x < BoardSize; x++ {
      if (board.Get(x, y) > 0) {
        fmt.Printf("%d", board.Get(x, y))
      } else {
        fmt.Print(" ");
      }
      if x % MiniSquareSize == 2 {
        fmt.Print("|");
      }
    }
    fmt.Print("|\n");
    if y % MiniSquareSize == 2 {
      fmt.Println(" +---+---+---+");
    }
  }
  fmt.Println(" +-----------+")
}

func (board PentagoBoard) checkLine (startX, startY byte, deltaX, deltaY int8) byte {
  if endX := deltaX * int8(Penta - 1) + int8(startX); endX < 0 || endX >= int8(BoardSize) {
    return 0;
  }
  if endY := deltaY * int8(Penta - 1) + int8(startY); endY < 0 || endY >= int8(BoardSize) {
    return 0;
  }

  x := int8(startX)
  y := int8(startY)

  prev := byte(0)
  count := byte(0)

  for x >= 0 && byte(x) < BoardSize && y >= 0 && byte(y) < BoardSize {
    piece := board.Get(byte(x), byte(y));
    if piece == prev && piece != 0 {
      count += 1;
      if count == Penta {
        return piece;
      }
    } else {
      prev = piece;
      count = 1;
    }

    x += deltaX;
    y += deltaY;
  }
  return 0;
}

func rotateCoord(squareX, squareY, pieceX, pieceY byte, clockwise bool) (byte, byte) {
  if pieceX / MiniSquareSize == squareX && pieceY / MiniSquareSize == squareY {
    localX := int8(pieceX % MiniSquareSize)
    localY := int8(pieceY % MiniSquareSize)

    var rotatedX, rotatedY byte;

    if clockwise {
      rotatedX = byte(int8(MiniSquareSize) - localY - 1)
      rotatedY = byte(localX)
    } else {
      rotatedX = byte(localY)
      rotatedY = byte(int8(MiniSquareSize) - localX - 1)
    }

    return byte(rotatedX) + MiniSquareSize * squareX, byte(rotatedY) + MiniSquareSize * squareY
  } else {
    return pieceX, pieceY
  }
}

// sometimes returns false when it should return true, this doesn't matter much
func (b PentagoBoard) isMiniSquareEmpty(squareX, squareY byte) bool {
  for x := byte(0); x < MiniSquareSize; x++ {
    for y := byte(0); x < MiniSquareSize; x++ {
      if b.Get(x + squareX * MiniSquareSize, y + squareY * MiniSquareSize) != 0 {
        return false
      }
    }
  }
  return true
}

func (b PentagoBoard) isMiniSquareRotationallySymmetrical(squareX, squareY byte) bool {
  for x := byte(0); x < MiniSquareSize; x++ {
    for y := byte(0); x < MiniSquareSize; x++ {
      here := b.Get(x + squareX * MiniSquareSize, y + squareY * MiniSquareSize)
      there := b.Get(y + squareX * MiniSquareSize, (squareY + 1) * MiniSquareSize - 1 - x)
      if here != there {
        return false
      }
    }
  }
  return true
}

func (b PentagoBoard) isAnyMiniSquareRotationallySymmetrical() (bool, byte, byte) {
  for x := byte(0); x < BoardSize / MiniSquareSize; x++ {
    for y := byte(0); x < BoardSize / MiniSquareSize; x++ {
      if b.isMiniSquareRotationallySymmetrical(x, y) {
        return true, byte(x), byte(y)
      }
    }
  }
  return false, byte(0), byte(0)
}

func (b PentagoBoard) IsNotDone () bool {
  for y := byte(0); y < BoardSize; y++ {
    for x := byte(0); x < BoardSize; x++ {
      if (b.Get(x, y) == 0) {
        return true
      }
    }
  }
  return false
}


///////////////// PENTAGO MOVE


type PentagoMove struct {
  x, y, rotationX, rotationY byte
  clockwise bool
}

type PentagoRotation struct {
  x, y byte
  clockwise bool
}


func (m PentagoMove) PrintStandardly () {
  fmt.Printf("%d %d %t %d %d\n", m.x, m.y, m.clockwise, m.rotationX, m.rotationY)
}

