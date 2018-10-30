package main

import (
  "fmt"
  "io"
  "github.com/trindadegm/siwumpgo/def"
  "github.com/trindadegm/siwumpgo/ia"
  "time"
)

func main() {
  fmt.Println("Wump, v0.2")
  var world def.World

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
  for {

    perception := sim.Perceive()
    status := sim.GetStatus()
    facing := sim.Compass()

    fmt.Println("SEED: ", seed)
    fmt.Println(sim, "\n", agent.GetSintesisInfo())
    fmt.Print(agent)
    fmt.Println(facing)
    fmt.Println("smell", "breeze", "shine", "shock", "scream")
    fmt.Println(perception)

    if status == def.EATEN {
      fmt.Println("Oh! it hurts! My leg! It's... it's tearing apart! He is munching me!\n" +
                "Crap! My intestine, it is roled up on his teeth! AAAAAAAAHHHHHHHH!\n\n" +
                "You watch in despair as the hunter is devoured. He couldn't do it.")
      break
    } else if status == def.BROKEN {
      fmt.Println("Well, let's see what is in this cave-eeeeeeeeeeeeeeee!...\n\n" +
                  "You listen to his voice fading as he falls on a seemingly very, very deep pit.\n\n" +
                  "CRASH CRACK PLAFT...\n" +
                  "He is dead.\n")
      break
    }

    action := agent.Decide(perception, status, facing)
    fmt.Println(action)
    sim.Act(action)

    time.Sleep(50 * time.Millisecond)
  }

  fmt.Println("Exiting...")
}
