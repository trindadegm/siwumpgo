package ia

import (
  "github.com/trindadegm/siwumpgo/def"
  //"container/list"
  "math"
  "fmt"
)

func (agent *StupidCognitiveAgent) toKillWumpusDecision() def.Action {
  if agent.model.HasKilled {
    agent.model.removeWumpus()
    agent.objective = GET_GOLD
    return def.IDLE
  }

  if agent.model.HunterPos.PosX == agent.model.WumpusPos.PosX ||
      agent.model.HunterPos.PosY == agent.model.WumpusPos.PosY {
    direction := getDirection(agent.model.HunterPos, agent.model.WumpusPos)
    agent.wantToShoot = true
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
  }

  if agent.model.WumpusPos.PosX != -1 { // There is a wumpus
    goal := agent.model.getClosestSquareToShoot()
    world := agent.model.World
    path := AStarPathfinding(agent.model.HunterPos, goal, agent.model.ExploredBoundaryX, agent.model.ExploredBoundaryY, world)

    if path.Len() > 0 {
      path.Remove(path.Front())
    }

    agent.pathToUse = path
  }

  fmt.Println("toKillWumpusDecision ERROR def.IDLE at END")
  return def.IDLE
}

func (model *Model) getClosestSquareToShoot() def.Point {
  minDist := math.MaxFloat64
  closest := def.Point {0, 0}

  for x := 0; x < model.ExploredBoundaryX; x++ {
    if model.World[model.HunterPos.PosY][x].IsSafe == YES {
      dist := heuristicCostEstimate(model.HunterPos, def.Point {x, model.HunterPos.PosY})
      if dist <= minDist {
        minDist = dist
        closest = def.Point {x, model.HunterPos.PosY}
      }
    }
  }

  for y := 0; y < model.ExploredBoundaryY; y++ {
    if model.World[y][model.HunterPos.PosY].IsSafe == YES {
      dist := heuristicCostEstimate(model.HunterPos, def.Point {model.HunterPos.PosX, y})
      if dist <= minDist {
        minDist = dist
        closest = def.Point {model.HunterPos.PosX, y}
      }
    }
  }

  return closest
}
