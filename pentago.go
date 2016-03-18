package main
import "fmt"
import "os"
import "bufio"
import "strings"
import "strconv"

func main() {
  PlayGame()
}


func TestRotation () {
  board := PentagoBoard{}
  board.PlayPieceMut(0, 0, 1 + 0)
  board.PlayPieceMut(1, 0, 1 + 1)
  board.PlayPieceMut(2, 0, 1 + 2)
  board.PlayPieceMut(0, 1, 1 + 3)
  board.PlayPieceMut(1, 1, 1 + 4)
  board.PlayPieceMut(2, 1, 1 + 5)
  board.PlayPieceMut(0, 2, 1 + 6)
  board.PlayPieceMut(1, 2, 1 + 7)
  board.PlayPieceMut(2, 2, 1 + 8)

  board.PrintBoard()

  board.RotateMut(0, 0, true)

  board.PrintBoard()
}

func GetMoveLive() {
  board := PentagoBoard{}
  playerJustMoved, _ := strconv.Atoi(os.Args[1])
  numberOfPlayers, _ := strconv.Atoi(os.Args[2])
  boardString := os.Args[3]

  for i := byte(0); i < 81; i++ {
    player := byte(boardString[i] - '0')
    board.PlayPieceMut(i / BoardSize, i % BoardSize, player)
  }

  board.playerJustMoved = byte(playerJustMoved)
  board.numberOfPlayers = byte(numberOfPlayers)

  move := GetAlmostRandomMove(board)
  move.PrintStandardly()
}

func PlayGame () {
  board := PentagoBoard{}

  board.PrintBoard()

  for (board.GetWinner() == 0) {
    fmt.Printf("There are %d moves\n", len(board.GetMoves()))

    move := GetHumanMove()

    board.DoMoveMut(move)

    board.PrintBoard()

    if (board.GetWinner() != 0) {
      fmt.Printf("you win!")
      return
    }

    move = GetAlmostRandomMove(board)
    fmt.Print("computer move: ")
    move.PrintStandardly()
    board.DoMoveMut(move)

    board.PrintBoard()
  }
  fmt.Printf("you lose!")
}

func GetHumanMove() PentagoMove {
  reader := bufio.NewReader(os.Stdin)
  fmt.Print("Enter move:\n")
  text, _ := reader.ReadString('\n')
  things := strings.SplitN(text, " ", 5)

  xPos, _ := strconv.Atoi(things[0])
  yPos, _ := strconv.Atoi(things[1])
  dir := things[2] == "l"
  xRotPos, _ := strconv.Atoi(things[3])
  yRotPos, _ := strconv.Atoi(strings.Trim(things[4], " \n"))

  return PentagoMove{byte(xPos), byte(yPos), byte(xRotPos), byte(yRotPos), dir}
}

