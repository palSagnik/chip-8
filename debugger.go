package chip8

type Debugger struct {
	paused bool
	step bool
	breakpoints map[uint16]bool
}

func NewDebugger() *Debugger {
	return &Debugger{
		paused: false,
		step: false,
		breakpoints: map[uint16]bool{},
	}
}

func (d *Debugger) Pause() {
	d.paused = true
	d.step = false
}

func (d *Debugger) AddBreakpoint(addr uint16) {
	d.breakpoints[addr] = true
}

func (d *Debugger) DeleteBreakpoint(addr uint16) {
	delete(d.breakpoints, addr)
}

func (d *Debugger) ShouldExecute(pc uint16) bool {
	// state => paused
	if d.paused {
		return false
	}
	
	// if pc present in breakpoint => pause state
	if d.breakpoints[pc] {
		d.paused = true
		return false
	}
	return true
}

