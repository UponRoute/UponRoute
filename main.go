package main

import (
    "flag"
    "fmt"
    "os"
    "time"
)

func main() {
    iterations := flag.Int("iterations", 10, "Number of simulation steps")
    delayMs := flag.Int("delay", 100, "Delay between steps in milliseconds")
    quiet := flag.Bool("quiet", false, "Run with minimal output")

    flag.Parse()

    if !*quiet {
        fmt.Println("SilentSim console app: starting simulation...")
    }

    for i := 1; i <= *iterations; i++ {
        if !*quiet {
            fmt.Printf("Step %d/%d\n", i, *iterations)
        }
        time.Sleep(time.Duration(*delayMs) * time.Millisecond)
    }

    if !*quiet {
        fmt.Println("Simulation complete.")
    }

    os.Exit(0)
}