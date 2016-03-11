package main
import "fmt"

const BoardSize byte = 9
const Penta byte = 5
const MiniSquareSize byte = 3
const NumberOfMiniSquares byte = BoardSize / MiniSquareSize

func main() {
  board := PentagoBoard{}
  board.PlayPieceMut(2, 2 + 1, 1);
  board.PlayPieceMut(3, 3 + 1, 1);
  board.PlayPieceMut(4, 4 + 1, 1);
  board.PlayPieceMut(5, 5 + 1, 1);
  board.RotateMut(1, 1, true);
  board.PrintBoard();

  fmt.Printf("and the winner is %d\n", board.GetWinner());

  found, move := board.GetUnambiguouslyCorrectMoveIfExists(1)

  fmt.Printf("get winning place says %b %d %d %d %d %b\n", found, move.x, move.y, move.rotationX, move.rotationY, move.clockwise)
}

type PentagoBoard struct {
  grid [BoardSize][BoardSize]byte
}

type PentagoMove struct {
  x, y, rotationX, rotationY byte
  clockwise bool
}

func (b PentagoBoard) Get (x, y byte) byte {
  return b.grid[y][x];
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
  // something like...
  var protectiveMoveFound = false
  var protectiveMoveX, protectiveMoveY byte

  for start := byte(0); start < BoardSize; start++ {
    startInt := int8(start)
    directions := [][]int8{{0, startInt, 0, 1},{startInt, 0, 1, 0},{startInt, 0, 1, 1},{startInt, 0, -1, 1},{0, startInt, 1, 1},{0, startInt, -1, 1}}

    for _, direction := range directions {
      startX := byte(direction[0])
      startY := byte(direction[1])
      deltaX := direction[2]
      deltaY := direction[3]
      success, x, y, winningPlayer := getWinningLineIfExists(board, startX, startY, deltaX, deltaY)
      if success {
        if winningPlayer == player {
          fmt.Printf("thing2: %d %d %d\n", x, y, winningPlayer)
          return true, x, y
        }
        protectiveMoveFound = true
        protectiveMoveX = x
        protectiveMoveY = y
      }
    }
  }

  if protectiveMoveFound {
    return true, protectiveMoveX, protectiveMoveY
  }

  return false, 0, 0
}

// if BoardSize were greater than Penta * 2 - 1, there could be multiple return values
// returns (isThereASolution?, x, y, player)
func getWinningLineIfExists (board PentagoBoard, startX, startY byte, deltaX, deltaY int8) (bool, byte, byte, byte) {
  if endX := deltaX * int8(Penta) + int8(startX); endX < 0 || endX >= int8(BoardSize) {
    return false, 0, 0, 0;
  }
  if endY := deltaY * int8(Penta) + int8(startY); endY < 0 || endY >= int8(BoardSize) {
    return false, 0, 0, 0;
  }

  x := int8(startX)
  y := int8(startY)

  prev := byte(0)
  count := byte(0)
  prevWasBlank := false

  for x >= 0 && byte(x) < BoardSize && y >= 0 && byte(y) < BoardSize {
    piece := board.Get(byte(x), byte(y));
    if piece == prev && piece != 0 {
      count += 1;
      if count == Penta - 1 {
        if prevWasBlank {
          fmt.Println("here1")
          return true, byte(x - deltaX * int8(Penta - 1)), byte(y - deltaY * int8(Penta - 1)), prev
        }
        if board.Get(byte(x + deltaX), byte(y + deltaY)) == 0 {
          fmt.Println("here2")
          return true, byte(x + deltaX), byte(y + deltaY), prev
        }
      }
    } else {
      prev = piece;
      count = 1;
      prevWasBlank = true;
    }

    x += deltaX;
    y += deltaY;
  }
  return false, 0, 0, 0;
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



// func main() {
//     in := bufio.NewReader(os.Stdin)

//     fmt.Print("Enter string: ")
//     s, err := in.ReadString('\n')
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     s = strings.TrimSpace(s)

//     fmt.Print("Enter 75000: ")
//     s, err = in.ReadString('\n')
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     n, err := strconv.Atoi(strings.TrimSpace(s))
//     if err != nil {
//         fmt.Println(err)
//         return
//     }
//     if n != 75000 {
//         fmt.Println("fail:  not 75000")
//         return
//     }
//     fmt.Println("Good")
// }
