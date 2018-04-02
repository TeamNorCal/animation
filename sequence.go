package animation

// Sequencing of animation effects across the entire model, to implement
// higher-level transitions requiring model-level animation effects

import (
	"image/color"
	"time"
)

// A sequencing step. Contains information about the effect(s) to perform
// and the universe[s] being targetted.
// Can be gated on completion of another referenced step, and/or delayed by
// an amount of time. If both gating step and delay are specified, the delay
// will be applied after the gating step completes.
type step struct {
	// The universe to which the step is applied
	universeId uint
	// Step that must complete before this step commences
	onCompletion *step
	// Delay before starting step [after prior step completes if set]
	delay time.Duration
	// The animation effect to play
	effect Animation
}

type sequence struct {
	steps []step
}

type SequenceBuilder interface {
	Then() SequenceBuilder
}

func Do() *sequence {
	s := sequence{nil}
	return s.Then()
}

func (b *sequence) Then() *sequence {
	return b
}

/*
 * SequenceRunner
 */

// Encapsulate a step waiting for another step to complete
type stepAndTime struct {
	runAt time.Time
	toRun *step
}

// Encapsulate a step waiting for another step to complete
type stepAndGatingStep struct {
	waitingOn *step // Step waiting for completion
	toRun     *step // Step to run
}

// SequenceRunner is responsible for executing a given sequence
type SequenceRunner struct {
	awaitingTime     []stepAndTime       // Queue of steps waiting on a particular time
	awaitingStep     []stepAndGatingStep // Queue of steps waiting on another step to complete
	activeByUniverse [][]*step           // Queue of steps that can be run on a particular universe. Only head of queue is procssed
	buffers          [][]color.RGBA      // Buffers to hold universe data
}

func deleteStep(a []*step, i int) []*step {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = nil
	return a[:len(a)-1]
}

func deleteSAGS(a []stepAndGatingStep, i int) []stepAndGatingStep {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = stepAndGatingStep{}
	return a[:len(a)-1]
}

func deleteSAT(a []stepAndTime, i int) []stepAndTime {
	copy(a[i:], a[i+1:])
	a[len(a)-1] = stepAndTime{}
	return a[:len(a)-1]
}

func (sr *SequenceRunner) scheduleAt(s *step, runAt time.Time) {
	sr.awaitingTime = append(sr.awaitingTime, stepAndTime{runAt, s})
}

// Check for steps that are waiting on another step to complete.
// 'now' is the time that should be considered to be the current time
func (sr *SequenceRunner) handleStepComplete(completed *step, now time.Time) {
	uniSteps := sr.activeByUniverse[completed.universeId]
	if len(uniSteps) > 0 && uniSteps[0] == completed {
		sr.activeByUniverse[completed.universeId] = deleteStep(uniSteps, 0)
	}
	for idx := 0; idx < len(sr.awaitingStep); {
		waiting := sr.awaitingStep[idx]
		if waiting.waitingOn == completed {
			s := waiting.toRun
			if s.delay > 0 {
				// Schedule to run after delay
				runAt := now.Add(s.delay)
				sr.scheduleAt(s, runAt)
			} else {
				// Run immediately
				sr.activeByUniverse[s.universeId] = append(sr.activeByUniverse[s.universeId], s)
			}
			// Delete this from the list of waiting steps (and don't increment index)
			sr.awaitingStep = deleteSAGS(sr.awaitingStep, idx)
		} else {
			idx++
		}
	}
}

// Check for any tasks that should run at this point
// 'now' is the time that should be considered to be the current time
func (sr *SequenceRunner) checkScheuledTasks(now time.Time) {
	for idx := 0; idx < len(sr.awaitingTime); {
		waiting := sr.awaitingTime[idx]
		if waiting.runAt.After(now) {
			// Time to run it!
			s := waiting.toRun
			sr.activeByUniverse[s.universeId] = append(sr.activeByUniverse[s.universeId], s)
			// Delete this from the list of waiting steps (and don't increment index)
			sr.awaitingTime = deleteSAT(sr.awaitingTime, idx)
		} else {
			idx++
		}
	}
}

// Generate frame data corresponding to the specified time (which should be
// monotonically increasing with each call)
// Return value indicates whether the sequence is complete.
func (sr *SequenceRunner) ProcessFrame(now time.Time) (done bool) {
	done = true
	for universeId, universe := range sr.activeByUniverse {
		if len(universe) > 0 {
			// We have an active step on this universe
			s := universe[0]
			// ...so we're not done yet
			done = false
			// Process the animation for the universe
			effectDone := s.effect.Frame(sr.buffers[universeId], now)
			if effectDone {
				sr.handleStepComplete(s, now)
			}
		}
	}

	// We are done if we procssed nothing and there are no more queued-up steps
	done = done && len(sr.awaitingStep) == 0 && len(sr.awaitingTime) == 0
	return
}

// Get current data for the specified universe. This data is updated by calling
// ProcessFrame for the universe
func (sr *SequenceRunner) UniverseData(universeId uint) []color.RGBA {
	return sr.buffers[universeId]
}
