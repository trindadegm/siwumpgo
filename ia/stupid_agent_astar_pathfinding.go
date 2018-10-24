package ia

import (
  "github.com/trindadegm/wump/def"
  //"fmt"
  "container/list"
  "math"
)

const (
  INT_INFINITY = math.MaxInt32
)

const (
  FLOAT_INFINITY = math.MaxFloat64
)

func AStarPathfinding(start, goal def.Point, dimensionX, dimensionY int, world [][]Cave) *list.List {
  //closedSet := list.New()
  //openSet := list.New()
  //openSet.PushBack(start)
  closedSet := make(map[def.Point]bool)
  openSet := make(map[def.Point]bool)
  openSet[start] = true

  cameFrom := make(map[def.Point]def.Point)

  gScore := make(map[def.Point]int)
  initWithIntInfinity(gScore, dimensionX, dimensionY)

  gScore[start] = 0

  fScore := make(map[def.Point]float64)
  initWithFloatInfinity(fScore, dimensionX, dimensionY)

  fScore[start] = heuristicCostEstimate(start, goal)

  for len(openSet) > 0 {
    current := lowerScoreOnMap(openSet, fScore)

    if current == goal {
      return reconstructPath(cameFrom, current)
    }

    delete(openSet, current)
    closedSet[current] = true

    adjs := getAdjacentPositions(current)
    for i := 0; i < 4; i++ {
      if def.IsInBounds(adjs[i], dimensionX, dimensionY) && world[adjs[i].PosY][adjs[i].PosX].IsSafe == YES {
        if closedSet[adjs[i]] {
          continue
        }

        tentativeGScore := gScore[current] + 1

        _, ok := openSet[adjs[i]]
        if !ok {
          openSet[adjs[i]] = true
        } else if tentativeGScore >= gScore[adjs[i]] {
          continue
        }

        cameFrom[adjs[i]] = current
        gScore[adjs[i]] = tentativeGScore
        fScore[adjs[i]] = float64(gScore[adjs[i]]) + heuristicCostEstimate(adjs[i], goal)
      }
    }
  }

  return nil
}

func initWithIntInfinity(vals map[def.Point]int, dimensionX, dimensionY int) {
  for y := 0; y < dimensionY; y++ {
    for x := 0; x < dimensionX;x++ {
      vals[def.Point {x, y}] = INT_INFINITY
    }
  }
}

func initWithFloatInfinity(vals map[def.Point]float64, dimensionX, dimensionY int) {
  for y := 0; y < dimensionY; y++ {
    for x := 0; x < dimensionX;x++ {
      vals[def.Point {x, y}] = FLOAT_INFINITY
    }
  }
}

func heuristicCostEstimate(start, goal def.Point) float64 {
  return math.Sqrt(math.Pow(float64(goal.PosX - start.PosX), 2.0) + math.Pow(float64(goal.PosY - start.PosY), 2.0))
}

func lowerScoreOnMap(points map[def.Point]bool, score map[def.Point]float64) def.Point {
  var toRet def.Point
  lower := float64(FLOAT_INFINITY)
  for key, _ := range points {
    if score[key] <= lower {
      lower = score[key]
      toRet = key
    }
  }
  return toRet
}

func reconstructPath(cameFrom map[def.Point]def.Point, current def.Point) *list.List {
  totalPath := list.New()
  totalPath.PushFront(current)

  for ok := true; ok; _, ok = cameFrom[current] {
    current = cameFrom[current]
    totalPath.PushFront(current)
  }

  return totalPath
}
