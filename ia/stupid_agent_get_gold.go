package ia

import (
  "github.com/trindadegm/siwumpgo/def"
  //"math"
  "fmt"
  "container/list"
)

func (agent *StupidCognitiveAgent) toGetGoldDecision() def.Action {
  X := agent.model.HunterPos.PosX
  Y := agent.model.HunterPos.PosY

  // If got gold
  if agent.model.World[Y][X].Shine {
    agent.model.HasGold = true
    agent.objective = EXIT_CAVE
    return def.PICK
  }

  if agent.pathToUse.Len() < 1 { // Must calculate a PATH
    var path *list.List
    closest, unvisited := agent.model.getClosestExplorationSite();
    if closest.PosX == -1 {
      agent.objective = KILL_WUMPUS
      //fmt.Println(" -----> LASCOU")
      return def.IDLE
    }
    if (closest != agent.model.HunterPos) {
      path = AStarPathfinding(agent.model.HunterPos, closest, agent.model.ExploredBoundaryX,
                               agent.model.ExploredBoundaryY, agent.model.World)
      //fmt.Print(" >> **** TRUE PATH: ")
      //for it := path.Front(); it != nil; it = it.Next() {
      //  fmt.Print(it.Value.(def.Point), " ")
      //}
      //fmt.Println()

      if path.Len() > 0 {
        path.Remove(path.Front())
      }
    } else {
      //fmt.Println(" >> %%% EMPTY PATH")
      path = list.New()
    }

    path.PushBack(unvisited)

    //fmt.Print(" >> PATH: ")
    //for it := path.Front(); it != nil; it = it.Next() {
    //  fmt.Print(it.Value.(def.Point), " ")
    //}
    //fmt.Println()

    agent.pathToUse = path
  }

  direction := getDirection(agent.model.HunterPos, agent.pathToUse.Front().Value.(def.Point))
  agent.pathToUse.Remove(agent.pathToUse.Front())

  agent.wantToMove = true
  switch direction {
  case def.NORTH:
    return def.FACE_NORTH
    break
  case def.EAST:
    return def.FACE_EAST
    break
  case def.SOUTH:
    return def.FACE_SOUTH
    break
  case def.WEST:
    return def.FACE_WEST
    break
  }
  fmt.Println("toGetGoldDecision ERROR def.IDLE at END")
  return def.IDLE
}

