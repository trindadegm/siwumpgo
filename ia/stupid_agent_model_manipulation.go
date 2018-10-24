package ia

import (
  "github.com/trindadegm/wump/def"
  "container/list"
  //"fmt"
)

func (model *Model) moveAgent() {
  switch model.HunterFacing {
  case def.NORTH:
    model.HunterPos.PosY--
    break
  case def.EAST:
    model.HunterPos.PosX++
    if (model.HunterPos.PosX+1 == model.ExploredBoundaryX) {
      model.increaseModelOnX()
    }
    break
  case def.SOUTH:
    model.HunterPos.PosY++
    if (model.HunterPos.PosY+1 == model.ExploredBoundaryY) {
      model.increaseModelOnY()
    }
    break
  case def.WEST:
    model.HunterPos.PosX--
    break
  }
}

func (model *Model) modelPerceptionsChanged(perceptions def.Perception) bool {
  X := model.HunterPos.PosX
  Y := model.HunterPos.PosY

  toReturn := false

  if (model.World[Y][X].Smell != perceptions.Smell) {
    toReturn = true
  } else if (model.World[Y][X].Breeze != perceptions.Breeze) {
    toReturn = true
  } else if (model.World[Y][X].Shine != perceptions.Shine) {
    toReturn = true
  }

  return toReturn
}

func (model *Model) noMoreWumpusOn(cave def.Point) {
  adjs := getAdjacentPositions(cave)
  //fmt.Print(" NO MORE WUMPUS ON ", cave)
  for i := 0; i < 4; i++ {
    if (def.IsInBounds(adjs[i], model.ExploredBoundaryX, model.ExploredBoundaryY)) {
      pointer := model.World[adjs[i].PosY][adjs[i].PosX].WumpusPointer
      elementToRemove := findPointOnList(pointer, cave)
      if (elementToRemove != nil) {
        pointer.Remove(elementToRemove)
        //removed := pointer.Remove(elementToRemove)
        //fmt.Print(" REMOVED ", removed.(def.Point))
      }
    }
  }
  //fmt.Println()
}

func (model *Model) noMorePitOn(cave def.Point) {
  adjs := getAdjacentPositions(cave)
  //fmt.Print(" NO MORE PIT ON ", cave)
  for i := 0; i < 4; i++ {
    if (def.IsInBounds(adjs[i], model.ExploredBoundaryX, model.ExploredBoundaryY)) {
      //fmt.Println(adjs[i], len(model.World[adjs[i].PosY]), len(model.World))
      pointer := model.World[adjs[i].PosY][adjs[i].PosX].PitPointer
      elementToRemove := findPointOnList(pointer, cave)
      if (elementToRemove != nil) {
        pointer.Remove(elementToRemove)
        //removed := pointer.Remove(elementToRemove)
        //fmt.Print(" REMOVED ", removed.(def.Point))
      }
    }
  }
  //fmt.Println()
}

func findPointOnList(list *list.List, point def.Point) *list.Element {
  for it := list.Front(); it != nil; it = it.Next() {
    if it.Value.(def.Point) == point {
      return it
    }
  }
  return nil
}

func (model *Model) increaseModelOnX() {
  for y := 0; y < len(model.World); y++ {
    model.World[y] = append(model.World[y], NewCave())
  }

  model.ExploredBoundaryX++
}

func (model *Model) increaseModelOnY() {
  caveSlice := make([]Cave, model.ExploredBoundaryX)
  for i := 0; i < len(caveSlice); i++ {
    caveSlice[i] = NewCave()
  }

  model.World = append(model.World, caveSlice)

  model.ExploredBoundaryY++
}

func getAdjacentPositions(point def.Point) [4]def.Point {
  var adjs [4]def.Point
  adjs[0] = def.Point {point.PosX, point.PosY-1}
  adjs[1] = def.Point {point.PosX+1, point.PosY}
  adjs[2] = def.Point {point.PosX, point.PosY+1}
  adjs[3] = def.Point {point.PosX-1, point.PosY}

  return adjs;
}

func (model *Model) getBestBoundaryX() int {
  if model.UpperBoundaryX != UNKNOWN {
    return model.UpperBoundaryX
  } else {
    return model.ExploredBoundaryX
  }
}

func (model *Model) getBestBoundaryY() int {
  if model.UpperBoundaryY != UNKNOWN {
    return model.UpperBoundaryY
  } else {
    return model.ExploredBoundaryY
  }
}

func (model *Model) removeWumpus() {
  if model.WumpusPos.PosX != -1 { // Knows where wumpus is
    adjs := getAdjacentPositions(model.WumpusPos)
    for i := 0; i < 4; i++ {
      if def.IsInBounds(adjs[i], model.ExploredBoundaryX, model.ExploredBoundaryY) {
        model.World[adjs[i].PosY][adjs[i].PosX].WumpusPointer.Init()
      }
    }
    model.World[model.WumpusPos.PosY][model.WumpusPos.PosX].HasWumpus = NO
  }

  // Searches for pointers
  for it := model.VisitedList.Front(); it != nil; it = it.Next() {
    point := it.Value.(def.Point)
    wplist := model.World[point.PosY][point.PosX].WumpusPointer
    for itwp := wplist.Front(); itwp != nil; itwp = itwp.Next() {
      pointed := itwp.Value.(def.Point)
      model.World[pointed.PosY][pointed.PosX].HasWumpus = NO
    }
    wplist.Init()
  }
}
