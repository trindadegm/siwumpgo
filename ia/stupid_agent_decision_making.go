package ia

import (
  "github.com/trindadegm/siwumpgo/def"
  "math"
  "fmt"
  //"container/list"
)

// Picks the closest visited square that is adjacent to a unvisited safe square, returns also the
// unvisited safe square
func (model *Model) getClosestExplorationSite() (def.Point, def.Point) {
  unvisited := def.Point {-1, 1}
  closest := def.Point {-1, -1}
  list := model.VisitedList
  srcPos := model.HunterPos

  minDist := math.MaxFloat64 // Infinity
  for it := list.Front(); it != nil; it = it.Next() {
    adjs := getAdjacentPositions(it.Value.(def.Point))
    for i := 0; i < 4; i++ {
      if def.IsInBounds(adjs[i], model.ExploredBoundaryX, model.ExploredBoundaryY) &&
         !model.World[adjs[i].PosY][adjs[i].PosX].Visited && model.World[adjs[i].PosY][adjs[i].PosX].IsSafe == YES {
        dstPos := it.Value.(def.Point)
        dist := dist2D(srcPos, dstPos)
        if dist <= minDist {
          minDist = dist
          closest = dstPos
          unvisited = adjs[i]
        }
      }
    }
  }

  return closest, unvisited
}

func dist2D(a, b def.Point) float64 {
  return math.Sqrt(math.Pow(float64(b.PosX - a.PosX), 2.0) + math.Pow(float64(b.PosY - a.PosY), 2.0))
}

func getDirection(src, dst def.Point) def.Direction {
  if src.PosY > dst.PosY {
    return def.NORTH
  } else if src.PosX < dst.PosX {
    return def.EAST
  } else if src.PosY < dst.PosY {
    return def.SOUTH
  } else if src.PosX > dst.PosX {
    return def.WEST
  } else {
    fmt.Println("ERROR on getDirection: compairing ", src, " with ", dst)
    return def.NORTH
  }
}
