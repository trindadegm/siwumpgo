package def

import (
  "fmt"
)

type Direction int

// Directions
const (
  NORTH Direction = iota
  EAST
  SOUTH
  WEST
)

type Status int

// Health status
const (
  ALIVE Status = iota
  BROKEN
  EATEN
)

type Action int

// Actions
const (
  IDLE Action = iota
  FACE_NORTH
  FACE_EAST
  FACE_SOUTH
  FACE_WEST
  MOVE
  SHOOT
  PICK
)

// This is used to keep track of the state of the simulation
type Simulation struct {
  world World             // The World, meaning the position of elements
  hunterPos Point         // Keeps track of the hunter easily, It is also on the world, but this way it is easier to find
  iterations int          // (!) (UNUSED) Counts the number of movements performed
  isHunterEaten bool      // When the hunter dies eaten by the wumpus
  isHunterBroken bool     // When the hunter dies falling of a pit
  hasHunterShot bool      // If the hunter spend an arrow
  hunterFacing Direction  // The direction the hunter faces
  hasHunterGold bool
  isHunterKnockedUp bool
  isWumpusScreaming bool
}

// Used to get where the hunter is on the world
func findHunter(world *World) Point {
  for y := 0; y < world.sizey; y++ {
    for x := 0; x < world.sizex; x++ {
      if world.squares[y][x].hasHunter {
        return Point {x, y} // Found
      }
    }
  }
  return Point {-1, -1}     // Not found
}

// Creates a new simulation from a world
func (s *Simulation) FromWorld(world World) {
  s.world = world
  s.hunterPos = findHunter(&world)
  s.iterations = 0
  s.isHunterEaten = false
  s.isHunterBroken = false
  s.hasHunterShot = false
  s.hunterFacing = NORTH
  s.hasHunterGold = false
  s.isWumpusScreaming = false
  s.isHunterKnockedUp = false
}

func IsInBounds(p Point, sizex, sizey int) bool {
  return (p.PosX < sizex && p.PosX > -1 && p.PosY < sizey && p.PosY > -1)
}

func perceiveSquare(perception *Perception, square Square) {
  if square.hasWump {
    perception.Smell = true
  }
  if square.terrain == Pit {
    perception.Breeze = true
  }
}

// Perceives the cave
func (sim *Simulation) Perceive() Perception {
  var perception Perception
  var point Point

  perception.Shock = sim.isHunterKnockedUp
  perception.Scream = sim.isWumpusScreaming

  hunterSquare := sim.world.squares[sim.hunterPos.PosY][sim.hunterPos.PosX]
  if hunterSquare.hasGold {
    perception.Shine = true
  }

  point = sim.hunterPos
  point.PosX++ // Move on X
  if IsInBounds(point, sim.world.sizex, sim.world.sizey) {
    perceiveSquare(&perception, sim.world.squares[point.PosY][point.PosX])
  }

  point = sim.hunterPos
  point.PosX-- // Move on X
  if IsInBounds(point, sim.world.sizex, sim.world.sizey) {
    perceiveSquare(&perception, sim.world.squares[point.PosY][point.PosX])
  }

  point = sim.hunterPos
  point.PosY++ // Move on Y
  if IsInBounds(point, sim.world.sizex, sim.world.sizey) {
    perceiveSquare(&perception, sim.world.squares[point.PosY][point.PosX])
  }

  point = sim.hunterPos
  point.PosY-- // Move on Y
  if IsInBounds(point, sim.world.sizex, sim.world.sizey) {
    perceiveSquare(&perception, sim.world.squares[point.PosY][point.PosX])
  }

  return perception
}

// Perceives status ailments... Yeah...
func (sim *Simulation) GetStatus() Status {
  switch {
  case sim.isHunterBroken:
    return BROKEN
  case sim.isHunterEaten:
    return EATEN
  }
  return ALIVE
}

// Gets the actual facing direction
func (sim *Simulation) Compass() Direction {
  return sim.hunterFacing
}

// Set the hunter where to look at (direction, NORTH, EAST, SOUTH or WEST)
func (sim *Simulation) face(direction Direction) {
  sim.hunterFacing = direction
}

// Sets the hunter position to a point
func (sim *Simulation) setHunterPos(point Point) {
  sim.world.squares[sim.hunterPos.PosY][sim.hunterPos.PosX].hasHunter = false // Hunter moved
  sim.world.squares[point.PosY][point.PosX].hasHunter = true // Hunter moved
  sim.hunterPos = point

  if sim.world.squares[point.PosY][point.PosX].hasWump { // Oh no.
    sim.isHunterEaten = true
  } else if sim.world.squares[point.PosY][point.PosX].terrain == Pit { // Oh no.
    sim.isHunterBroken = true
  }
}

// Moves the hunter to the square being faced. If the square is
// out of bounds, it will set the shock perception
func (sim *Simulation) move() {
  point := sim.hunterPos

  switch (sim.hunterFacing) {
  case NORTH:
    point.PosY--
    break
  case EAST:
    point.PosX++
    break
  case SOUTH:
    point.PosY++
    break
  case WEST:
    point.PosX--
  }

  if IsInBounds(point, sim.world.sizex, sim.world.sizey) {
    sim.setHunterPos(point)
    sim.isHunterKnockedUp = false
  } else {
    sim.isHunterKnockedUp = true
  }
}

// Shoots an arrow to the direction being faced
func (sim *Simulation) shoot() {
  if sim.hasHunterShot {
    return // No more arrows
  }

  diff := Point {0, 0}

  switch (sim.hunterFacing) {
  case NORTH:
    diff.PosY--
    break;
  case EAST:
    diff.PosX++
    break;
  case SOUTH:
    diff.PosY++
    break;
  case WEST:
    diff.PosX--
  }

  point := Point {sim.hunterPos.PosX+diff.PosX, sim.hunterPos.PosY+diff.PosY}

  for ; IsInBounds(point, sim.world.sizex, sim.world.sizey); {
    if sim.world.squares[point.PosY][point.PosX].hasWump {
      sim.isWumpusScreaming = true
      sim.world.squares[point.PosY][point.PosX].hasWump = false // Killed him, YAY!
    }
    point.PosX += diff.PosX
    point.PosY += diff.PosY
  }
}

func (sim* Simulation) pick() {
  point := sim.hunterPos

  if sim.world.squares[point.PosY][point.PosX].hasGold {
    sim.world.squares[point.PosY][point.PosX].hasGold = false
    sim.hasHunterGold = true
  }
}

func (sim *Simulation) Act(a Action) {
  switch a {
  case FACE_NORTH:
    sim.face(NORTH);
    break;
  case FACE_EAST:
    sim.face(EAST);
    break;
  case FACE_SOUTH:
    sim.face(SOUTH);
    break;
  case FACE_WEST:
    sim.face(WEST);
    break;
  case MOVE:
    sim.move();
    break;
  case SHOOT:
    sim.shoot();
    break;
  case PICK:
    sim.pick();
    break;
  case IDLE:
    // Do nothing
  }
  sim.iterations++;
}

func (sim Simulation) String() string {
  return fmt.Sprintf("%v %v\n%s", sim.iterations, sim.hunterPos, sim.world)
}

func (d Direction) String() string {
  var directionName string
  switch d {
  case NORTH:
    directionName = "↑ NORTH"
    break
  case EAST:
    directionName = "→ EAST"
    break
  case SOUTH:
    directionName = "↓ SOUTH"
    break
  case WEST:
    directionName = "← WEST"
    break
  }
  return directionName;
}

func (action Action) String() string {
  switch action {
  case IDLE:
    return "IDLE"
    break
  case FACE_NORTH:
    return "FACE_NORTH"
    break
  case FACE_EAST:
    return "FACE_EAST"
    break
  case FACE_SOUTH:
    return "FACE_SOUTH"
    break
  case FACE_WEST:
    return "FACE_WEST"
    break
  case MOVE:
    return "MOVE"
    break
  case SHOOT:
    return "SHOOT"
    break
  case PICK:
    return "PICK"
    break
  }
  return "?"
}
