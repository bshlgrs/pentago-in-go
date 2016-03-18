package main

import "testing"
import "math/rand"
import "fmt"

func TestDiagonal(t *testing.T) {
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
  board.DoMoveMut(move)

  if board.GetWinner() != 1 {
    board.PrintBoard();
    t.Errorf("The board did not find a winner")
  }
}

func TestHorizontalBoard(t *testing.T) {
  board := PentagoBoard{}
  board.PlayPieceMut(2, 5, 1);
  board.PlayPieceMut(3, 5, 1);
  board.PlayPieceMut(4, 5, 1);
  board.PlayPieceMut(5, 5, 1);
  if board.GetWinner() != 0 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it had a winner")
  }

  found, move := board.GetUnambiguouslyCorrectMoveIfExists(1)

  if !found {
    t.Errorf("The board did not find its correct move")
  }
  board.DoMoveMut(move)
  board.PlayPieceMut(6, 5, 1);
  if board.GetWinner() != 1 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it didn't have a winner")
  }
}

func TestVerticalBoard(t *testing.T) {
  board := PentagoBoard{}
  board.PlayPieceMut(5, 2, 1);
  board.PlayPieceMut(5, 3, 1);
  board.PlayPieceMut(5, 4, 1);
  board.PlayPieceMut(5, 5, 1);
  if board.GetWinner() != 0 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it had a winner")
  }

  found, move := board.GetUnambiguouslyCorrectMoveIfExists(1)

  if !found {
    t.Errorf("The board did not find its correct move")
  }

  board.DoMoveMut(move)
  board.PlayPieceMut(5, 1, 1);
  if board.GetWinner() != 1 {
    board.PrintBoard();
    t.Errorf("The board incorrectly thought that it didn't have a winner")
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

func BenchmarkPlayGameRandomly (b *testing.B) {
  for i := 0; i < 100; i++ {
    playGameRandomly()
  }
}

func playGameRandomly() byte {
  board := PentagoBoard{}

  for (board.GetWinner() == 0 && board.IsNotDone()) {
    move := GetAlmostRandomMove(board)
    // move.PrintStandardly()
    board.DoMoveMut(move)
    // board.PrintBoard()
  }

  return board.GetWinner();
}
