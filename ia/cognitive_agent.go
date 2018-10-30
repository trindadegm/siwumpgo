package ia

import (
  "github.com/trindadegm/siwumpgo/def"
  "container/list"
  "fmt"
)

const (
  UNKNOWN = -1
)

type Objective int

const (
  GET_GOLD Objective = iota
  EXIT_CAVE
  KILL_WUMPUS
)

type DoubtBool int

const (
  MAYBE DoubtBool = iota
  NO
  YES
)

type Cave struct {
  Visited bool
  Smell bool
  Breeze bool
  Shine bool

  HasPit DoubtBool
  HasWumpus DoubtBool
  IsSafe DoubtBool

  PitPointer *list.List
  WumpusPointer *list.List
}

func (cave *Cave) New() {
  cave.Visited = false
  cave.Smell = false
  cave.Breeze = false
  cave.Shine = false
  cave.HasPit = MAYBE
  cave.HasWumpus = MAYBE
  cave.IsSafe = MAYBE
  cave.PitPointer = list.New()
  cave.WumpusPointer = list.New()
}

func NewCave() Cave {
  var cave Cave
  cave.New()
  return cave
}

type Model struct {
  World [][]Cave
  HasShot bool
  HasGold bool
  HasKilled bool

  UpperBoundaryX int
  UpperBoundaryY int

  ExploredBoundaryX int
  ExploredBoundaryY int

  HunterPos def.Point
  HunterFacing def.Direction

  VisitedList *list.List

  WumpusPos def.Point
}

type StupidCognitiveAgent struct {
  model Model
  wantToMove bool
  wantToShoot bool
  pathToUse *list.List
  wumpusAim def.Point
  moved bool
  objective Objective
}

func (agent *StupidCognitiveAgent) New() {
  agent.model.World = [][]Cave { {NewCave(), NewCave()}, {NewCave(), NewCave()} }

  agent.model.HasShot = false
  agent.model.HasGold = false
  agent.model.HasKilled = false

  agent.model.ExploredBoundaryX = 2
  agent.model.ExploredBoundaryY = 2

  agent.model.HunterPos = def.Point {0, 0}
  //agent.model.World[0][0].Visited = true
  //agent.model.World[0][0].IsSafe = true
  agent.model.HunterFacing = def.NORTH

  agent.model.VisitedList = list.New()

  agent.model.WumpusPos = def.Point {-1, -1}

  agent.wantToMove = false
  agent.wantToShoot = false
  agent.objective = GET_GOLD
  agent.pathToUse = list.New()
  agent.wumpusAim = def.Point {-1, -1}
}

func (agent *StupidCognitiveAgent) Decide(perception def.Perception, status def.Status, facing def.Direction) def.Action {
  if status != def.ALIVE { // Is dead. So it will refuse to move, the game should be over by now anyway
    return def.IDLE
  }

  agent.model.HunterFacing = facing

  if perception.Scream {
    agent.model.HasKilled = true
  }

  if agent.wantToShoot {
    agent.wantToShoot = false
    agent.model.updateShoot(facing)
    agent.objective = GET_GOLD
    return def.SHOOT
  } else if agent.wantToMove {
    agent.moved = true
    agent.wantToMove = false
    return def.MOVE
  }

  if perception.Shock && agent.moved {
    switch facing {
    case def.SOUTH:
      agent.model.shockOnBoundaryY()
      break
    case def.EAST:
      agent.model.shockOnBoundaryX()
    }
  } else if agent.moved { // If moved and didn't hit
    agent.model.moveAgent() // Update model positions and things like that
  }
  agent.moved = false

  agent.model.deduceCaveModel(perception)

  var action def.Action
  switch agent.objective {
  case GET_GOLD:
    action = agent.toGetGoldDecision()
  case EXIT_CAVE:
    action = agent.toExitCaveDecision()
  case KILL_WUMPUS:
    action = agent.toKillWumpusDecision()
  }

  if (action == def.FACE_NORTH && facing == def.NORTH) ||
     (action == def.FACE_EAST && facing == def.EAST) ||
     (action == def.FACE_SOUTH && facing == def.SOUTH) ||
     (action == def.FACE_WEST && facing == def.WEST) {
    if agent.wantToShoot {
      agent.wantToShoot = false
      return def.SHOOT
    } else if agent.wantToMove {
      agent.moved = true
      agent.wantToMove = false
      return def.MOVE
    }
  } else {
    return action
  }

  fmt.Println("Decide ERROR def.IDLE at END")
  return def.IDLE
}

// Generates a string representing the agent's knowlege about the world as a map.
// It is a very long and loop based method.
func (agent StupidCognitiveAgent) String() string {
  var world [][]Cave
  var toReturn string
  world = agent.model.World

  // Prints the size of the explored boundary, which is the size of the world on
  // the agent's head.
  lenY := agent.model.ExploredBoundaryY
  lenX := agent.model.ExploredBoundaryX
  toReturn += fmt.Sprintf("%d %d\n    ", lenX, lenY)
  // Prints the column numbers
  for x := 0; x < lenX; x++ {
    toReturn += fmt.Sprintf("%2d  ", x)
  }
  toReturn += "\n"
  // For each line of squares
  for y := 0; y < lenY; y++ {
    // For each square on a line, prints the text representing some part of the state.
    // This is a complicated thing to explain. I am doing three loops for each line
    // here because each square is made out of NINE characters representing the knowledge
    // of the state!
    // The square is like this:
    // ↑↑→    Which I know is crap. A squre is made out of nine characters like that
    // ←i→    So each square is made out of three lines, so each for prints the three
    // ←↓↓    top characters of each square. Deal with it.
    for x := 0; x < lenX; x++ {
      smellPtrUpChar := '.'
      if findPointOnList(world[y][x].WumpusPointer, def.Point {x, y-1}) != nil {
        smellPtrUpChar = '↑'
      }
      pitPtrUpChar := '.'
      if findPointOnList(world[y][x].PitPointer, def.Point {x, y-1}) != nil {
        pitPtrUpChar = '↑'
      }
      smellPtrRightChar := '.'
      if findPointOnList(world[y][x].WumpusPointer, def.Point {x+1, y}) != nil {
        smellPtrRightChar = '→'
      }
      if x == 0 {
        toReturn += "    "
      }
      toReturn += fmt.Sprintf("%c%c%c|", smellPtrUpChar, pitPtrUpChar, smellPtrRightChar)
    }
    toReturn += "\n"
    for x := 0; x < lenX; x++ {
      inPlaceChar := '?'
      if agent.model.HunterPos.PosX == x && agent.model.HunterPos.PosY == y {
        inPlaceChar = 'i'
      } else if world[y][x].HasWumpus == YES {
        inPlaceChar = 'W'
      } else if world[y][x].HasPit == YES {
        inPlaceChar = 'O'
      } else if world[y][x].Visited {
        inPlaceChar = '_'
      } else if world[y][x].IsSafe == YES {
        inPlaceChar = '+'
      }
      pitPtrRightChar := '.'
      if findPointOnList(world[y][x].PitPointer, def.Point {x+1, y}) != nil {
        pitPtrRightChar = '→'
      }
      pitPtrLeftChar := '.'
      if findPointOnList(world[y][x].PitPointer, def.Point {x-1, y}) != nil {
        pitPtrLeftChar = '←'
      }
      if x == 0 {
        toReturn += fmt.Sprintf(" %2d ", y)
      }
      toReturn += fmt.Sprintf("%c%c%c|", pitPtrLeftChar, inPlaceChar, pitPtrRightChar)
    }
    toReturn += "\n"
    for x := 0; x < lenX; x++ {
      smellPtrDownChar := '.'
      if findPointOnList(world[y][x].WumpusPointer, def.Point {x, y+1}) != nil {
        smellPtrDownChar = '↓'
      }
      pitPtrDownChar := '.'
      if findPointOnList(world[y][x].PitPointer, def.Point {x, y+1}) != nil {
        pitPtrDownChar = '↓'
      }
      smellPtrLeftChar := '.'
      if findPointOnList(world[y][x].WumpusPointer, def.Point {x-1, y}) != nil {
        smellPtrLeftChar = '←'
      }
      if x == 0 {
        toReturn += "    "
      }
      toReturn += fmt.Sprintf("%c%c%c|", smellPtrLeftChar, pitPtrDownChar, smellPtrDownChar)
    }
    toReturn += "\n    "
    for x := 0; x < lenX; x++ {
      toReturn += "---+"
    }
    toReturn += "\n"
  }

  return toReturn
}

func (cave Cave) String() string {
  if cave.HasWumpus == YES {
    return "W "
  } else if cave.HasPit == YES{
    return "O "
  } else if cave.Shine {
    return "* "
  } else if cave.Visited {
    return "_ "
  } else if cave.IsSafe == YES{
    return ". "
  } else {
    return "? "
  }
}

// Returns simple information about the agent's state in a string
func (agent *StupidCognitiveAgent) GetSintesisInfo() string {
  objPosX, objPosY := -1, -1
  if agent.pathToUse.Len() > 0 {
    objPosX, objPosY = agent.pathToUse.Back().Value.(def.Point).PosX, agent.pathToUse.Back().Value.(def.Point).PosY
  }

  var objDesc string
  switch agent.objective {
  case GET_GOLD:
    objDesc = "get_gold"
    break
  case EXIT_CAVE:
    objDesc = "exit_cave"
    break
  case KILL_WUMPUS:
    objDesc = "kill_wumpus"
    break
  }

  return fmt.Sprintf("MODEL: {%d, %d} OBJECTIVE: %s ON {%d, %d}", len(agent.model.World[0]),
                      len(agent.model.World), objDesc, objPosX, objPosY)
}
