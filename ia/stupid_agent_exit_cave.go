package ia

import (
  "github.com/trindadegm/siwumpgo/def"
  "fmt"
)

func (agent *StupidCognitiveAgent) toExitCaveDecision() def.Action {
  if agent.model.HunterPos.PosX == 0 && agent.model.HunterPos.PosY == 0 {
    fmt.Println(" *DOING THE VICTORY DANCE* ")
    return def.IDLE
  }

  if agent.pathToUse.Len() < 1 {
    agent.pathToUse = AStarPathfinding(agent.model.HunterPos, def.Point {0, 0}, agent.model.ExploredBoundaryX,
                                       agent.model.ExploredBoundaryY, agent.model.World)

    agent.pathToUse.Remove(agent.pathToUse.Front())
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
  fmt.Println("toExitCaveDecision ERROR def.IDLE at END")
  return def.IDLE
}
