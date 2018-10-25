package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"time"

	"github.com/c9s/goprocinfo/linux"
	flags "github.com/jessevdk/go-flags"
)

// Version by Makefile
var Version string

type cmdOpts struct {
	Pid       uint64 `short:"p" long:"pid" description:"PID" required:"true"`
	KeyPrefix string `long:"key-prefix" description:"Metric key prefix" required:"true"`
	Version   bool   `short:"v" long:"version" description:"Show version"`
}

type processStats struct {
	Now    uint64 `json:"now"`
	CPU    uint64 `json:"cpu"`
	Utime  uint64 `json:"utime"`
	Stime  uint64 `json:"stime"`
	Cutime int64  `json:"cutime"`
	Cstime int64  `json:"cstime"`
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func writeStats(filename string, ps processStats) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	jb, err := json.Marshal(ps)
	if err != nil {
		return err
	}
	_, err = file.Write(jb)
	return err
}

func readStats(filename string) (processStats, error) {
	ps := processStats{}
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		return ps, err
	}
	err = json.Unmarshal(d, &ps)
	if err != nil {
		return ps, err
	}
	return ps, nil
}

func cpuJiffer() uint64 {

	cpu, _ := linux.ReadStat("/proc/stat")

	// cpu's jiffies
	e := cpu.CPUStatAll
	return e.User + e.Nice + e.System + e.Idle
}

func getStats(opts cmdOpts) error {

	c := cpuJiffer()
	now := uint64(time.Now().Unix())
	p, err := linux.ReadProcess(opts.Pid, "/proc")
	if err != nil {
		return fmt.Errorf("failed to fetch stats: %v", err)
	}

	ps := processStats{
		Now:    now,
		CPU:    c,
		Utime:  p.Stat.Utime,
		Cutime: p.Stat.Cutime,
		Stime:  p.Stat.Stime,
		Cstime: p.Stat.Cstime,
	}

	tmpDir := os.TempDir()
	curUser, _ := user.Current()
	prevPath := filepath.Join(tmpDir, fmt.Sprintf("%s-process-status-%d", curUser.Uid, opts.Pid))

	if !fileExists(prevPath) {
		err = writeStats(prevPath, ps)
		if err != nil {
			return fmt.Errorf("failed to save stats: %v", err)
		}
		fmt.Fprintf(os.Stderr, "Notice: first time execution command\n")
		return nil
	}

	prev, err := readStats(prevPath)
	if err != nil {
		return fmt.Errorf("failed to load stats: %v", err)
	}

	user := int64(ps.Utime-prev.Utime) + (ps.Cutime - prev.Cutime)
	system := int64(ps.Stime-prev.Stime) + (ps.Cstime - prev.Cstime)
	up := (float64(user) / float64(ps.CPU-prev.CPU)) * 100
	sp := (float64(system) / float64(ps.CPU-prev.CPU)) * 100
	fmt.Printf("process-status.cpu_%s.user\t%f\t%d\n", opts.KeyPrefix, up, now)
	fmt.Printf("process-status.cpu_%s.system\t%f\t%d\n", opts.KeyPrefix, sp, now)
	err = writeStats(prevPath, ps)
	if err != nil {
		return fmt.Errorf("failed to save stats: %v", err)
	}
	return nil
}

func main() {
	os.Exit(_main())
}

func _main() int {
	opts := cmdOpts{}
	psr := flags.NewParser(&opts, flags.Default)
	_, err := psr.Parse()
	if err != nil {
		return 1
	}
	if opts.Version {
		fmt.Printf(`%s %s
Compiler: %s %s
`,
			os.Args[0],
			Version,
			runtime.Compiler,
			runtime.Version())
		return 0
	}

	err = getStats(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}
	return 0
}
