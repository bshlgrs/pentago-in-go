package main

import "math/rand"

func FirstHeuristic(board PentagoBoard) []float {
  winner := board.GetWinner()
  if winner {
    returnValue := make([]float, board.numberOfPlayers)
    returnValue[winner - 1] = 1.0
    return returnValue
  }

  // points for three in a row and twos-in-a-row that can be converted
  // points for
}

func GetAlmostRandomMove(board PentagoBoard) PentagoMove {
  found, move := board.GetUnambiguouslyCorrectMoveIfExists(board.NextToPlay())

  if found {
    return move
  }

  for  {
    x := byte(rand.Int()) % BoardSize
    y := byte(rand.Int()) % BoardSize
    if board.Get(x, y) == 0 {
      return PentagoMove{x, y,
        byte(rand.Int()) % NumberOfMiniSquares, byte(rand.Int()) % NumberOfMiniSquares, rand.Int() % 2 == 0}
    }
  }

  return PentagoMove{}
}

func (board PentagoBoard) GetUnambiguouslyCorrectMoveIfExists (player byte) (bool, PentagoMove) {
  for x := byte(0); x < NumberOfMiniSquares; x++ {
    for y := byte(0); y < NumberOfMiniSquares; y++ {
      board.RotateMut(x, y, true)

      found, moveX, moveY := getWinningPlaceIfExists(player, board)
      if found {
        board.RotateMut(x, y, false)

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
    directions := [][]int8{{0, startInt, 1, 0},{startInt, 0, 0, 1},{startInt, 0, 1, 1},{startInt, 0, -1, 1},{0, startInt, 1, 1},{0, startInt, -1, 1}}

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

