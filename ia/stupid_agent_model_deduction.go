package ia

import (
  "github.com/trindadegm/wump/def"
  //"fmt"
  //"container/list"
)

func (model *Model) deduceCaveModel(perception def.Perception) {
  X := model.HunterPos.PosX
  Y := model.HunterPos.PosY

  mustUpdateModel := model.modelPerceptionsChanged(perception)

  model.World[Y][X].Smell = perception.Smell
  model.World[Y][X].Breeze = perception.Breeze
  model.World[Y][X].Shine = perception.Shine

  // Never been here before
  if model.World[Y][X].Visited == false {
    model.VisitedList.PushBack(model.HunterPos)
    model.World[Y][X].Visited = true
    model.World[Y][X].HasPit = NO
    model.World[Y][X].HasWumpus = NO
    model.World[Y][X].IsSafe = YES

    mustUpdateModel = true // Because it shall update
  }

  if mustUpdateModel {
    adjs := getAdjacentPositions(model.HunterPos)
    for i := 0; i < 4; i++ {
      if def.IsInBounds(adjs[i], model.ExploredBoundaryX, model.ExploredBoundaryY) {
        if model.World[adjs[i].PosY][adjs[i].PosX].Visited == false {
          // Then point to adjacent cave that it may be wumpy, only if it may have wumpus
          if perception.Smell && model.World[adjs[i].PosY][adjs[i].PosX].HasWumpus != NO {
            list := model.World[Y][X].WumpusPointer
            if findPointOnList(list, adjs[i]) == nil {
              //fmt.Print("PS ", X, Y, adjs[i])
              list.PushBack(adjs[i])
            }
            //fmt.Println(" -> ", list.Len())
          } else {
            model.World[adjs[i].PosY][adjs[i].PosX].HasWumpus = NO
            //fmt.Println("RPS", X, Y, adjs[i])
            model.noMoreWumpusOn(adjs[i])
          }
          // Then point to adjacent cave that it may be pity, only if it may have a pit
          if perception.Breeze && model.World[adjs[i].PosY][adjs[i].PosX].HasPit != NO {
            list := model.World[Y][X].PitPointer
            if findPointOnList(list, adjs[i]) == nil {
              //fmt.Print("PB ", X, Y, adjs[i])
              list.PushBack(adjs[i])
            }
            //fmt.Println(" -> ", list.Len())
          } else {
            model.World[adjs[i].PosY][adjs[i].PosX].HasPit = NO
            //fmt.Println("RPB", X, Y, adjs[i])
            model.noMorePitOn(adjs[i])
          }

          if !perception.Smell && !perception.Breeze {
            model.World[adjs[i].PosY][adjs[i].PosX].IsSafe = YES
          }
        } else { // The adjacent cave was visited once
          list := model.World[adjs[i].PosY][adjs[i].PosX].WumpusPointer
          element := findPointOnList(list, model.HunterPos)
          if element != nil { // So it pointed to the cave I am
            list.Remove(element) // Because if I am here, here is safe
          }
          list = model.World[adjs[i].PosY][adjs[i].PosX].PitPointer
          element = findPointOnList(list, model.HunterPos)
          if element != nil { // So it pointed to the cave I am
            list.Remove(element) // Because if I am here, here is safe
          }
        }
      }
    }
  }

  it := model.VisitedList.Front()
  for it != nil {
    point := it.Value.(def.Point)

    if model.World[point.PosY][point.PosX].WumpusPointer.Len() == 1 {
      hasWumpusPoint := model.World[point.PosY][point.PosX].WumpusPointer.Front().Value.(def.Point)
      model.World[hasWumpusPoint.PosY][hasWumpusPoint.PosX].HasWumpus = YES
      model.WumpusPos = hasWumpusPoint
    }

    if model.World[point.PosY][point.PosX].PitPointer.Len() == 1 {
      hasPitPoint := model.World[point.PosY][point.PosX].PitPointer.Front().Value.(def.Point)
      model.World[hasPitPoint.PosY][hasPitPoint.PosX].HasPit = YES
    }

    it = it.Next()
  }

  if model.WumpusPos.PosX != -1 { // Found wumpus
    // it is a visited node
    for it := model.VisitedList.Front(); it != nil; it = it.Next() {
      // point is a visited point
      point := it.Value.(def.Point)
      // wPointerList is the list of places suspected to have wumpus
      wPointerList := model.World[point.PosY][point.PosX].WumpusPointer
      // itwp is a suspected node of containing a wumpus
      for itwp := wPointerList.Front(); itwp != nil; {
        // pointed is a suspected point of containing wumpus
        pointed := itwp.Value.(def.Point)
        // But pointed is not with the wumpus, because the wumpus is somewhere else
        if pointed.PosX != model.WumpusPos.PosX || pointed.PosY != model.WumpusPos.PosY {
          next := itwp.Next()
          // So forget about pointing to him
          wPointerList.Remove(itwp)

          // He is not guilty
          model.World[pointed.PosY][pointed.PosX].HasWumpus = NO
          model.World[pointed.PosY][pointed.PosX].IsSafe = YES

          itwp = next
          continue
        }
        itwp = itwp.Next()
      }
    }
  }
}
