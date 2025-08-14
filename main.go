package main

import (
    "bufio"
    "flag"
    "fmt"
    "os"
    "runtime"
    "time"
)

func main() {
    iterations := flag.Int("iterations", 10, "Number of simulation steps")
    delayMs := flag.Int("delay", 100, "Delay between steps in milliseconds")
    quiet := flag.Bool("quiet", false, "Run with minimal output")
    pause := flag.Bool("pause", false, "Wait for Enter before exiting (useful when double-clicked on Windows)")

    flag.Parse()

    if !*quiet {
        fmt.Println("SilentSim console app: starting simulation...")
        fmt.Printf("Iterations=%d, Delay=%dms\n", *iterations, *delayMs)
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

    if shouldPause(*pause) {
        fmt.Print("Press Enter to exit...")
        _, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
    }

    os.Exit(0)
}

func shouldPause(userRequested bool) bool {
    if userRequested {
        return true
    }
    // Heuristic: when double-clicked on Windows, there are usually no args and no terminal env.
    if runtime.GOOS == "windows" {
        if len(os.Args) == 1 && os.Getenv("WT_SESSION") == "" && os.Getenv("TERM") == "" {
            return true
        }
    }
    return false
}