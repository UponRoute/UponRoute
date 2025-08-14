package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

func (v Vector3) Add(o Vector3) Vector3 { return Vector3{v.X + o.X, v.Y + o.Y, v.Z + o.Z} }
func (v Vector3) Sub(o Vector3) Vector3 { return Vector3{v.X - o.X, v.Y - o.Y, v.Z - o.Z} }
func (v Vector3) MulScalar(s float64) Vector3 { return Vector3{v.X * s, v.Y * s, v.Z * s} }
func (v Vector3) Dot(o Vector3) float64 { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }
func (v Vector3) Length() float64 { return math.Sqrt(v.Dot(v)) }
func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if l == 0 {
		return v
	}
	return v.MulScalar(1.0 / l)
}

func parseVector3(s string) (Vector3, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 3 {
		return Vector3{}, fmt.Errorf("expected 3 comma-separated numbers, got %q", s)
	}
	x, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	y, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	z, err3 := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
	if err1 != nil || err2 != nil || err3 != nil {
		return Vector3{}, fmt.Errorf("invalid vector components in %q", s)
	}
	return Vector3{X: x, Y: y, Z: z}, nil
}

func parseTargets(s string) ([]Vector3, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	chunks := strings.Split(s, ";")
	out := make([]Vector3, 0, len(chunks))
	for _, c := range chunks {
		v, err := parseVector3(strings.TrimSpace(c))
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, nil
}

func distancePointToSegment(point, a, b Vector3) float64 {
	ab := b.Sub(a)
	ap := point.Sub(a)
	abLenSq := ab.Dot(ab)
	if abLenSq == 0 {
		return ap.Length()
	}
	t := ap.Dot(ab) / abLenSq
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	closest := a.Add(ab.MulScalar(t))
	return point.Sub(closest).Length()
}

func runNearHitSimulation(origin Vector3, direction Vector3, distance float64, radius float64, targets []Vector3, quiet bool) int {
	if !quiet {
		fmt.Printf("Origin=(%.2f,%.2f,%.2f) Direction=(%.2f,%.2f,%.2f) Distance=%.2f Radius=%.2f\n",
			origin.X, origin.Y, origin.Z, direction.X, direction.Y, direction.Z, distance, radius)
	}
	if len(targets) == 0 {
		// Default demo targets near the shot path
		targets = []Vector3{{X: 8, Y: 0, Z: 0.5}, {X: 12, Y: 1.2, Z: -0.8}, {X: 20, Y: -0.5, Z: 0}}
		if !quiet {
			fmt.Println("No targets specified; using demo targets at:", targets)
		}
	}
	end := origin.Add(direction.Normalize().MulScalar(distance))
	var hits int
	for idx, t := range targets {
		d := distancePointToSegment(t, origin, end)
		if d <= radius {
			hits++
			if !quiet {
				fmt.Printf("Target %d at (%.2f,%.2f,%.2f): HIT (distance %.2f <= %.2f)\n", idx+1, t.X, t.Y, t.Z, d, radius)
			}
		} else if !quiet {
			fmt.Printf("Target %d at (%.2f,%.2f,%.2f): miss (distance %.2f > %.2f)\n", idx+1, t.X, t.Y, t.Z, d, radius)
		}
	}
	if !quiet {
		fmt.Printf("Total hits: %d of %d\n", hits, len(targets))
	}
	return hits
}

func main() {
	// General flags
	quiet := flag.Bool("quiet", false, "Run with minimal output")
	pause := flag.Bool("pause", false, "Wait for Enter before exiting (useful when double-clicked on Windows)")

	// Legacy simple simulation flags (retained)
	iterations := flag.Int("iterations", 10, "Number of simulation steps (legacy demo mode)")
	delayMs := flag.Int("delay", 100, "Delay between steps in milliseconds (legacy demo mode)")

	// Near-hit simulation flags
	mode := flag.String("mode", "bullet", "Mode: 'bullet' for near-hit simulation or 'steps' for legacy demo")
	originStr := flag.String("origin", "0,0,0", "Shot origin as x,y,z")
	directionStr := flag.String("direction", "1,0,0", "Shot direction as x,y,z (will be normalized)")
	distance := flag.Float64("distance", 50, "Shot distance")
	radius := flag.Float64("radius", 2.0, "Near-hit radius")
	targetsStr := flag.String("targets", "", "Semicolon-separated list of target positions as x,y,z;x,y,z;...")

	flag.Parse()

	if *mode == "steps" {
		if !*quiet {
			fmt.Println("SilentSim: legacy step simulation...")
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

	// Near-hit bullet simulation
	origin, err := parseVector3(*originStr)
	if err != nil {
		fmt.Println("Invalid --origin:", err)
		os.Exit(2)
	}
	direction, err := parseVector3(*directionStr)
	if err != nil {
		fmt.Println("Invalid --direction:", err)
		os.Exit(2)
	}
	targets, err := parseTargets(*targetsStr)
	if err != nil {
		fmt.Println("Invalid --targets:", err)
		os.Exit(2)
	}
	_ = runNearHitSimulation(origin, direction, *distance, *radius, targets, *quiet)

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