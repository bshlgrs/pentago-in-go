package main

import "testing"
import "math/rand"
import "fmt"

func TestABunchOfThings(t *testing.T) {
  board := PentagoBoard{}
  board.PlayPieceMut(2, 2 + 1, 1);
  board.PlayPieceMut(3, 3 + 1, 1);
  board.PlayPieceMut(4, 4 + 1, 1);
  board.PlayPieceMut(5, 5 + 1, 1);
  if board.GetWinner() != 0 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it had a winner")
  }

  board.RotateMut(1, 1, true);

  if board.GetWinner() != 0 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it had a winner")
  }

  found, move := board.GetUnambiguouslyCorrectMoveIfExists(1)

  if !found {
    t.Errorf("The board did not find its correct move")
  }
  board.PlayMoveMut(move, 1)

  if board.GetWinner() != 1 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it had a winner")
  }
}


func getRandomBoard() PentagoBoard {
  board := PentagoBoard{}

  for x := byte(0); x < 9; x++ {
    for y := byte(0); y < 9; y++ {
      board.PlayPieceMut(x, y, byte(rand.Int() % 5))
    }
  }
  return board;
}

func BenchmarkGetWinner (b *testing.B) {
  board := getRandomBoard();

  b.ResetTimer();
  x := board.GetWinner();

  if x > 4 {
    fmt.Printf("Some really weird shit is happening")
  }
}

func BenchmarkGetUnambiguouslyCorrectMove (b *testing.B) {
  board := getRandomBoard();

  b.ResetTimer();
  _, move := board.GetUnambiguouslyCorrectMoveIfExists(1);

  if move.x > 9 {
    fmt.Printf("Some really weird shit is happening")
  }
}
