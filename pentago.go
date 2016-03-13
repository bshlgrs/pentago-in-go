package main
import "fmt"
import "os"
import "strconv"

func main() {
  board := PentagoBoard{}
  currentPlayer, _ := strconv.Atoi(os.Args[1])
  numberOfPlayers, _ := strconv.Atoi(os.Args[2])
  boardString := os.Args[3]

  for i := byte(0); i < 81; i++ {
    player := byte(boardString[i] - '0')
    board.PlayPieceMut(i / BoardSize, i % BoardSize, player)
  }

  move := GetMove(board, byte(currentPlayer), byte(numberOfPlayers))
  move.PrintStandardly()
  // board.PrintBoard()
  // fmt.Printf("%d %d\n", currentPlayer, numberOfPlayers)
}

func GetMove(board PentagoBoard, currentPlayer, numberOfPlayers byte) PentagoMove {
  found, move := board.GetUnambiguouslyCorrectMoveIfExists(currentPlayer)

  if found {
    return move
  }

  for i := byte(0); i < 81; i++ {
    if board.Get(i % 9, i / 9) == 0 {
      return PentagoMove{i % 9, i / 9, 1, 1, true}
    }
  }

  return PentagoMove{}
}

const BoardSize byte = 9
const Penta byte = 5
const MiniSquareSize byte = 3
const NumberOfMiniSquares byte = BoardSize / MiniSquareSize

type PentagoBoard struct {
  grid [BoardSize][BoardSize]byte
}

type PentagoMove struct {
  x, y, rotationX, rotationY byte
  clockwise bool
}

func (m PentagoMove) PrintStandardly () {
  fmt.Printf("%d %d %d %d %t\n", m.x, m.y, m.rotationX, m.rotationY, m.clockwise)
}

func (b PentagoBoard) Get (x, y byte) byte {
  return b.grid[y][x];
}

func (b *PentagoBoard) PlayMoveMut (move PentagoMove, player byte) {
  b.PlayPieceMut(move.x, move.y, player);
  b.RotateMut(move.rotationX, move.rotationY, move.clockwise);
}

func (b *PentagoBoard) PlayPieceMut (x, y, player byte) {
  b.grid[y][x] = player;
}

func (b PentagoBoard) PlayPiece (x, y, player byte) PentagoBoard {
  b.grid[y][x] = player;
  return b;
}

func (b *PentagoBoard) RotateMut (x, y byte, clockwise bool) {
  startX := x * MiniSquareSize;
  startY := y * MiniSquareSize;
  temp := b.grid[startY][startX]

  if clockwise {
    b.grid[startY][startX] = b.grid[startY + 2][startX]
    b.grid[startY + 2][startX] = b.grid[startY + 2][startX + 2]
    b.grid[startY + 2][startX + 2] = b.grid[startY][startX + 2]
    b.grid[startY][startX + 2] = temp

    temp := b.grid[startY][startX + 1]
    b.grid[startY][startX + 1] = b.grid[startY + 1][startX + 2]
    b.grid[startY + 1][startX + 2] = b.grid[startY + 2][startX + 1]
    b.grid[startY + 2][startX + 1] = b.grid[startY + 1][startX]
    b.grid[startY + 1][startX] = temp
  } else {
    b.RotateMut(x, y, true)
    b.RotateMut(x, y, true)
    b.RotateMut(x, y, true)
  }
}


func (board PentagoBoard) PrintBoard () {
  fmt.Println("+-----------+")
  for y := byte(0); y < BoardSize; y++ {
    fmt.Print("|");
    for x := byte(0); x < BoardSize; x++ {
      if (board.Get(x, y) > 0) {
        fmt.Printf("%d", board.Get(x, y))
      } else {
        fmt.Print(" ");
      }
      if x == 2 || x == 5 {
        fmt.Print("|");
      }
    }
    fmt.Print("|\n");
    if y == 2 || y == 5 {
      fmt.Println("+---+---+---+");
    }
  }
  fmt.Println("+-----------+")
}

func (board PentagoBoard) checkLine (startX, startY byte, deltaX, deltaY int8) byte {
  if endX := deltaX * int8(Penta) + int8(startX); endX < 0 || endX >= int8(BoardSize) {
    return 0;
  }
  if endY := deltaY * int8(Penta) + int8(startY); endY < 0 || endY >= int8(BoardSize) {
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

func (board PentagoBoard) GetWinner () byte {
  for y := byte(0); y < BoardSize; y++ {
    winner := board.checkLine(0, y, 0, 1)
    if winner != 0 {
      return winner;
    }
  }

  for x := byte(0); x < BoardSize; x++ {
    winner := board.checkLine(x, 0, 1, 0)
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

func (board PentagoBoard) GetUnambiguouslyCorrectMoveIfExists (player byte) (bool, PentagoMove) {
  for x := byte(0); x < NumberOfMiniSquares; x++ {
    for y := byte(0); y < NumberOfMiniSquares; y++ {
      board.RotateMut(x, y, true)

      found, moveX, moveY := getWinningPlaceIfExists(player, board)
      if found {
        board.RotateMut(x, y, false)

        fmt.Printf("thing: %d %d %d %d\n", x, y, moveX, moveY)
        realMoveX, realMoveY := rotateCoord(x, y, moveX, moveY, false)

        return true, PentagoMove{realMoveX, realMoveY, x, y, true};
      }

      board.RotateMut(x, y, true)
      board.RotateMut(x, y, true)

      found, moveX, moveY = getWinningPlaceIfExists(player, board)
      if found {
        board.RotateMut(x, y, false)
        realMoveX, realMoveY := rotateCoord(x, y, moveX, moveY, true)
        return true, PentagoMove{realMoveX, realMoveY, x, y, false};
      }

      board.RotateMut(x, y, true)
    }
  }

  return false, PentagoMove{}
}

func getWinningPlaceIfExists (player byte, board PentagoBoard) (bool, byte, byte) {
  for start := byte(0); start < BoardSize; start++ {
    startInt := int8(start)
    directions := [][]int8{{0, startInt, 0, 1},{startInt, 0, 1, 0},{startInt, 0, 1, 1},{startInt, 0, -1, 1},{0, startInt, 1, 1},{0, startInt, -1, 1}}

    for _, direction := range directions {
      startX := byte(direction[0])
      startY := byte(direction[1])
      deltaX := direction[2]
      deltaY := direction[3]
      success, x, y := getWinningLineIfExists(board, startX, startY, deltaX, deltaY, player)
      if success {
        return true, x, y
      }
    }
  }

  return false, 0, 0
}

// if BoardSize were greater than Penta * 2 - 1, there could be multiple return values
// returns (isThereASolution?, x, y)
func getWinningLineIfExists (board PentagoBoard, startX, startY byte, deltaX, deltaY int8, player byte) (bool, byte, byte) {
  if endX := deltaX * int8(Penta) + int8(startX); endX < 0 || endX >= int8(BoardSize) {
    return false, 0, 0;
  }
  if endY := deltaY * int8(Penta) + int8(startY); endY < 0 || endY >= int8(BoardSize) {
    return false, 0, 0;
  }

  x := int8(startX)
  y := int8(startY)

  previousLineLength := byte(0)
  currentLineLength := byte(0)

  emptySpaceFound := false
  resultX := startX
  resultY := startY

  for x >= 0 && byte(x) < BoardSize && y >= 0 && byte(y) < BoardSize {
    piece := board.Get(byte(x), byte(y));
    if piece == player {
      currentLineLength += 1
      if currentLineLength + previousLineLength == Penta - 1 && emptySpaceFound {
        return true, resultX, resultY
      }
    } else if piece == 0 {
      previousLineLength = currentLineLength
      currentLineLength = 0
      resultX = byte(x)
      resultY = byte(y)
      emptySpaceFound = true

      if previousLineLength == Penta - 1 {
        return true, resultX, resultY
      }
    } else {
      currentLineLength = 0
      previousLineLength = 0
      emptySpaceFound = false
    }

    x += deltaX;
    y += deltaY;
  }
  return false, 0, 0;
}

// given a list of numbers like [0, 1, 1, 1, 0, 1, 2, 2]


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

