package def

import (
  "math/rand"
  "math"
  "fmt"
)

// The bigger this is, the most likely it is for pits to spawn
// There is a problem, however, because there are places where
// pits cannot spawn (relative to letting a path to gold) and
// the spawn of gold, hunter and wumpus will remove any pit
// where they are spawned.
const pitFactor = 10

// The bigger this is, the more likely is to gold to spawn
// further away from the hunter
const unluckyFactor = 35

const (
  Pit     = iota
  Normal  = iota
)

type Point struct {
  PosX int
  PosY int
}

type Perception struct {
  Smell bool
  Breeze bool
  Shine bool
  Shock bool
  Scream bool
}

type Square struct {
  terrain int
  hasGold bool
  hasWump bool
  hasHunter bool
}

type World struct {
  squares [][]Square
  sizex int
  sizey int
}

/*
Stringer for Square

One may want to optimize this "thing", so I'll explain why it is this
way (as of NOW):
  1: I could not spare 2 characters (the space must be there so I'll
  not count it) to represent something, so I would only use one
  character.
  2: There is a way for the hunter to be over the gold, and for the
  wumpus to be over the gold.

This made so that I had to draw a different symbol for when the
wumpus or hunter were together with the gold.
*/
func (square Square) String() string {
  if square.terrain == Pit {
    return "O "
  } else if !square.hasGold {
    switch {
    case square.hasWump:
      return "W "
    case square.hasHunter:
      return "i "
    }
  } else {
    switch {
    case square.hasWump:
      return "# "
    case square.hasHunter:
      return "$ "
    default:
      return "* "
    }
  }

  return ". " // No other option suffices, so it is nothing
}

/*
Stringer for World
*/
func (world World) String() string {
  if len(world.squares) == 0{
    panic("World was not created correctly, dimension x is 0")
  } else if len(world.squares[0]) == 0 {
    panic("World was not created correctly, dimension y is 0")
  }

  str := ""
  for y := 0; y < world.sizey; y++ {
    for x := 0; x < world.sizex; x++ {
      str += world.squares[y][x].String();
    }
    str += "\n"
  }

  return str
}

/*
Ok, this function is nuts. The reason for it, is that I wanted to have
a simple problem generator for me to create an intelligent agent, right?
It didn't need to be perfect, for now. It is a small college homework.

So, for every problem be solvable, I needed to list some points that
led in a path to the gold, from the hunter, and I should not put
the bottomless pits in that path. So there is always at least this
(crazy, ill created) path from the hunter to the gold.

This function returns a map with all the points one shouldn't put
a pit into (a map for easy lookup) and then the point at the end,
where the gold shall be put.

I'll divide this function in three parts, inside
*/
func (world *World) guaranteePoints(positionHunter Point) (map[Point]bool, Point) {
  // Length calculated based on map size, it may be minor, but never
  // bigger
  length := int(2 * math.Sqrt(float64(world.sizex * world.sizex) + float64(world.sizey * world.sizey)))

  visited := make(map[Point]bool)
  positionPath := positionHunter

  // The big for
  // Image this algorigthm like this, the gold is now a living thing and
  // it is right where the hunter is! Now it will move away from it.
  // From now on, I'll consider the gold to be "moving".
  for n := 0; n < length; n++ {
    // Decide to which direction to go, it must always move as:
    // {0, 1}, {0, -1}, {1, 0}, {-1, 0} (there is a "thing" after this)
    // But diagonal movement is not movement ok? Never.
    which := rand.Intn(2)
    dx := 2 * rand.Intn(2) - 1
    dy := 2 * rand.Intn(2) - 1

    // The unluckyFactor, ah...
    // Without this code, as the gold is dumb, it is very likely to remain
    // next to the hunter, and it is even common to have the gold right
    // on the hunter Square! This is not good, too easy. So there is a 
    // factor for the hunter missfortune: unluckyFactor. This factor
    // makes likely, that every time the gold tries to move, it will
    // move TWO squares AWAY from the hunter. A little help for the
    // dumb gold right?
    randomPick := rand.Intn(100)
    if which == 0 {
      dx *= 0
      if (randomPick < unluckyFactor) {
        dy = 2 * int(math.Abs(float64(dy)))
      }
    } else {
      dy *= 0
      if (randomPick < unluckyFactor) {
        dx = 2 * int(math.Abs(float64(dx)))
      }
    }

    // Keep the gold inbounds! Very important
    // There is a trick here, when the unlucky factor attacks, the dx or dy
    // (one of them) will be 2, if this put them out of bounds, I'm making
    // them go back to 1, so that they wont turn arround and go -2.
    // This is because the unlucky double movement is already a problem
    // regarding the map of positions.
    if positionPath.PosX + dx >= world.sizex || positionPath.PosX + dx < 0 {
      if dx > 1 {
        dx = 1
      }
      dx *= -1
    }
    if positionPath.PosY + dy >= world.sizey || positionPath.PosY + dy < 0 {
      if dy > 1 {
        dy = 1
      }
      dy *= -1
    }

    nextPositionPath := Point {positionPath.PosX + dx, positionPath.PosY + dy}

    // Ok, the final part, if the place I'm putting the gold was already
    // visited by him before, disregard this for iteration, the dumb brat
    // lost another step. If not...
    if _, isVisited := visited[nextPositionPath]; !isVisited {
      positionPath = nextPositionPath
      visited[positionPath] = true

      // Consider that if he performed a unlucky movement, he should also
      // mark as visited the square he jumped
      if dx > 1 {
        visited[Point {positionPath.PosX-1, positionPath.PosY}] = true
      } else if dy > 1 {
        visited[Point {positionPath.PosX, positionPath.PosY-1}] = true
      }
    }
  }

  return visited, positionPath
}

/*
This function creates a new world, based on the seed given to it
*/
func (world *World) New(sizex, sizey int, seed int64) {
  if sizex == 0 || sizey == 0 {
    panic("Invalid dimensions!");
  }

  rand.Seed(seed)

  world.sizex = sizex
  world.sizey = sizey

  // Guarantees a path from 0, 0, to point lastPoint, where the gold will be
  guaranteedPoints, lastPoint := world.guaranteePoints(Point {0, 0}) // Hunter is always at 0, 0

  // Creates a map made from normal terrain and pit terrain
  // No pits on guaranteed points
  world.squares = make([][]Square, sizey, sizey)
  for y := 0; y < sizey; y++ {
    world.squares[y] = make([]Square, sizex, sizex)
    for x := 0; x < sizex; x++ {
      randomPick := rand.Intn(100)
      if randomPick < pitFactor && !guaranteedPoints[Point {x, y}] {
        world.squares[y][x] = Square{terrain: Pit, hasGold: false, hasWump: false, hasHunter: false}
      } else {
        world.squares[y][x] = Square{terrain: Normal, hasGold: false, hasWump: false, hasHunter: false}
      }
    }
  }

  // Place hunter
  world.squares[0][0].terrain = Normal
  world.squares[0][0].hasHunter = true

  // Place gold
  world.squares[lastPoint.PosY][lastPoint.PosX].hasGold = true

  // Place wump
  wx := rand.Intn(world.sizex)
  wy := rand.Intn(world.sizey)
  if wx == 0 && wy == 0 { // Not over hunter. Would be instant defeat
    wx = 1
  }
  world.squares[wy][wx].terrain = Normal
  world.squares[wy][wx].hasWump = true
}

func (world *World) FromString(worldD string, lenX, lenY int) {
  world.squares = make([][]Square, lenY)
  world.sizey = lenY
  world.sizex = lenX

  //lineSize := 0
  //x, y := 0, 0
  //for i := 0; i < len(worldD); i++ {
  //  char := worldD[i]
  //  //world.squares[y] = append(world.squares[y], Square{Normal, false, false, false})
  //  //lineSize++
  //  switch char {
  //  case '\\':
  //    //world.squares = append(world.squares, make([]Square, 0))
  //    //world.sizex = lineSize
  //    //lineSize = 0
  //    y++
  //    world.squares[y] = make([]Square, lenX)
  //    x = -1
  //    break
  //  case 'W':
  //    world.squares[y][x].hasWump = true
  //    break
  //  case 'i':
  //    world.squares[y][x].hasHunter = true
  //    break
  //  case 'O':
  //    world.squares[y][x].terrain = Pit
  //    break
  //  case '.':
  //    world.squares[y][x].terrain = Normal
  //    break
  //  default:
  //    break
  //  }
  //  x++
  //}
  i := 0
  for y := 0; y < lenY; y++ {
    world.squares[y] = make([]Square, lenX)
    for x := 0; x < lenX; x++ {
      fmt.Println(y, x, len(world.squares[y]))
      switch worldD[i] {
      case 'W':
        world.squares[y][x].hasWump = true
        world.squares[y][x].terrain = Normal
        break
      case 'i':
        world.squares[y][x].hasHunter = true
        world.squares[y][x].terrain = Normal
        break
      case 'O':
        world.squares[y][x].terrain = Pit
        break
      case '.':
        world.squares[y][x].terrain = Normal
        break
      case '*':
        world.squares[y][x].hasGold = true
        world.squares[y][x].terrain = Normal
        break
      case '#':
        world.squares[y][x].hasGold = true
        world.squares[y][x].hasWump = true
        world.squares[y][x].terrain = Normal
      default:
        break
      }
      i++
    }
  }
}
