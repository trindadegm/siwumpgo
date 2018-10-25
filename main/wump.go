package main

import (
  "fmt"
  "io"
  "github.com/trindadegm/siwumpgo/def"
  "github.com/trindadegm/siwumpgo/ia"
  "time"
)

func getWord(str string) (word, rest string) {
  divisor := 0
  for ; divisor < len(str); divisor++ {
    if str[divisor] == ' ' {
      break
    }
  }
  if divisor == len(str) { // Only one word
    return str, ""
  }

  return str[:divisor], str[divisor+1:]
}

func processMoveLine(sim *def.Simulation, direction string) {
  switch direction {
  case "north":
    sim.Act(def.FACE_NORTH)
    sim.Act(def.MOVE)
  case "south":
    sim.Act(def.FACE_SOUTH)
    sim.Act(def.MOVE)
  case "east":
    sim.Act(def.FACE_EAST)
    sim.Act(def.MOVE)
  case "west":
    sim.Act(def.FACE_WEST)
    sim.Act(def.MOVE)
  default:
    fmt.Println("Unknown direction")
  }
}

func processShootLine(sim *def.Simulation, direction string) {
  switch direction {
  case "north":
    sim.Act(def.FACE_NORTH)
    sim.Act(def.SHOOT)
  case "south":
    sim.Act(def.FACE_SOUTH)
    sim.Act(def.SHOOT)
  case "east":
    sim.Act(def.FACE_EAST)
    sim.Act(def.SHOOT)
  case "west":
    sim.Act(def.FACE_WEST)
    sim.Act(def.SHOOT)
  default:
    fmt.Println("Unknown direction")
  }
}

func main() {
  fmt.Println("Wump, v0.2")
  var world def.World

  //world.New(7, 7, 99)
  //world.FromString("i.OOO.."+
  //                 "....O.."+
  //                 "..O...."+
  //                 ".....O."+
  //                 "......."+
  //                 "......O"+
  //                 "....W.*", 7, 7)
  //world.New(7, 7, 123451)
  //world.New(7, 7, 90908080)
  //world.New(14, 14, seed)
  //world.NewEx(14, 14, 1540378218, 10, 35)
  //world.New(14, 14, 1540341382)

  var mode int
  var seed int64
  var sizex, sizey int
  fmt.Printf("INPUT MODE: ")
  fmt.Scanf("%d", &mode)
  if mode == 0 {
    var pf, uf int
    fmt.Printf("WD: ")
    fmt.Scanf("%d %d %d %d %d\n", &sizex, &sizey, &seed, &pf, &uf)
    if seed == 0 {
      seed = time.Now().Unix()
    }
    fmt.Printf("%d %d %d %d %d\n", sizex, sizey, seed, uf, pf)
    world.NewEx(sizex, sizey, seed, pf, uf)
  } else {
    var swd, read string
    fmt.Scanf("%d %d\n", &sizex, &sizey)
    for {
      _, err := fmt.Scanln(&read)
      if err == io.EOF {
        break
      }
      swd += read
    }
    fmt.Printf("%d %d %s", sizex, sizey, swd)
    world.FromString(swd, sizex, sizey)
  }

  var sim def.Simulation

  sim.FromWorld(world)

  var agent ia.StupidCognitiveAgent
  agent.New()
  //goto labelEndFor
  for {

    perception := sim.Perceive()
    status := sim.GetStatus()
    facing := sim.Compass()

    fmt.Println("SEED: ", seed)
    fmt.Println(sim, "\n", agent.GetSintesisInfo())
    //fmt.Print(agent)
    fmt.Println(facing)
    fmt.Println("smell", "breeze", "shine", "shock", "scream")
    fmt.Println(perception)

    switch status {
    case def.EATEN:
        fmt.Println("Oh! it hurts! My leg! It's... it's tearing apart! He is munching me!\n" +
                  "Crap! My intestine, it is roled up on his teeth! AAAAAAAAHHHHHHHH!\n\n" +
                  "You watch in despair as the hunter is devoured. He couldn't do it.")
      goto labelEndFor
    case def.BROKEN:
      fmt.Println("Well, let's see what is in this cave-eeeeeeeeeeeeeeee!...\n\n" +
                  "You listen to his voice fading as he falls on a seemingly very, very deep pit.\n\n" +
                  "CRASH CRACK PLAFT...\n" +
                  "He is dead.\n")
      goto labelEndFor
    }

    action := agent.Decide(perception, status, facing)
    fmt.Println(action)
    sim.Act(action)

    time.Sleep(50 * time.Millisecond)
  }
  labelEndFor:

  //path := ia.AStarPathfinding(def.Point {0, 0}, def.Point{4, 6}, 8, 8)
  //for it := path.Front(); it != nil; it = it.Next() {
  //  fmt.Print(it.Value.(def.Point))
  //}
  //fmt.Println()

  fmt.Println("Exiting...")
}
