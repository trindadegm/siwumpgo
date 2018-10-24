package ia

import (
  "github.com/trindadegm/wump/def"
  "math"
  "fmt"
  //"container/list"
)

//func (agent *StupidCognitiveAgent) toGetGoldDecision() def.Action {
//  X := agent.model.HunterPos.PosX
//  Y := agent.model.HunterPos.PosY
//
//  // If got gold
//  if agent.model.World[Y][X].Shine {
//    agent.model.HasGold = true
//    agent.objective = EXIT_CAVE
//    return def.PICK
//  }
//
//  if agent.pathToUse.Len() < 1 { // Must calculate a PATH
//    var path *list.List
//    closest, unvisited := agent.model.getClosestExplorationSite();
//    if closest.PosX == -1 {
//      fmt.Println(" -----> LASCOU")
//      return def.IDLE
//    }
//    if (closest != agent.model.HunterPos) {
//      path = AStarPathfinding(agent.model.HunterPos, closest, agent.model.ExploredBoundaryX,
//                               agent.model.ExploredBoundaryY, agent.model.World)
//      fmt.Print(" >> **** TRUE PATH: ")
//      for it := path.Front(); it != nil; it = it.Next() {
//        fmt.Print(it.Value.(def.Point), " ")
//      }
//      fmt.Println()
//
//      if path.Len() > 0 {
//        path.Remove(path.Front())
//      }
//    } else {
//      fmt.Println(" >> %%% EMPTY PATH")
//      path = list.New()
//    }
//
//    path.PushBack(unvisited)
//
//    fmt.Print(" >> PATH: ")
//    for it := path.Front(); it != nil; it = it.Next() {
//      fmt.Print(it.Value.(def.Point), " ")
//    }
//    fmt.Println()
//
//    agent.pathToUse = path
//  }
//
//  direction := getDirection(agent.model.HunterPos, agent.pathToUse.Front().Value.(def.Point))
//  agent.pathToUse.Remove(agent.pathToUse.Front())
//
//  agent.wantToMove = true
//  switch direction {
//  case def.NORTH:
//    return def.FACE_NORTH
//    break
//  case def.EAST:
//    return def.FACE_EAST
//    break
//  case def.SOUTH:
//    return def.FACE_SOUTH
//    break
//  case def.WEST:
//    return def.FACE_WEST
//    break
//  }
//  return def.IDLE
//}

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
