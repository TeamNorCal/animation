package animation

import (
	"image/color"
	"testing"
	"time"
)

type testAnimation uint8
type multiRunAnimation int

func (a *testAnimation) Frame(buf []color.RGBA, frameTime time.Time) bool {
	//	fmt.Printf("In Frame for %v buf len %d, time %v\n", *a, len(buf), frameTime.Nanosecond()/1000)
	count := uint8(*a)
	for idx := range buf {
		buf[idx] = color.RGBA{count, count + 1, count + 2, 0x00}
		count += 3
	}
	return true
}

func (a *multiRunAnimation) Frame(buf []color.RGBA, frameTime time.Time) bool {
	*a--
	return *a <= 0
}

func TestDeleteStep(t *testing.T) {
	s1 := Step{UniverseID: 1}
	s2 := Step{UniverseID: 2}
	s3 := Step{UniverseID: 3}
	stepsArr := [3]*Step{&s1, &s2, &s3}
	steps := stepsArr[:]
	no0 := deleteStep(steps, 0)
	if len(no0) != 2 || no0[0] != &s2 || no0[1] != &s3 {
		t.Errorf("Delete 0 not as expected: %v", no0)
	}
	if stepsArr != [3]*Step{&s2, &s3, nil} {
		t.Errorf("Underlying array not as expected after delete 0: %v", no0)
	}
}

func TestSimpleSingleSequence(t *testing.T) {
	t.SkipNow()
	const universeID uint = 3
	ta := testAnimation(1)
	s := Step{UniverseID: universeID, Effect: &ta}
	seq := Sequence{[]Step{s}}
	sr := NewSequenceRunner([]uint{1, 1, 1, 10})
	startTime := time.Unix(0, 0)
	sr.InitSequence(seq, startTime)

	// Check for right-sized empty buffer
	buf := sr.UniverseData(universeID)
	if len(buf) != 10 {
		t.Errorf("Unexpected buffer size %d", len(buf))
	}
	zero := color.RGBA{0, 0, 0, 0}
	for idx, p := range buf {
		if p != zero {
			t.Errorf("Pixel %d not 0 (%x)", idx, p)
		}
	}

	// Call the effect to generate data
	done := sr.ProcessFrame(startTime)
	if done {
		t.Error("Done on first frame")
	}

	// Check for updated data
	count := uint8(1)
	buf = sr.UniverseData(universeID)
	if len(buf) != 10 {
		t.Errorf("Unexpected buffer size %d", len(buf))
	}
	for idx, p := range buf {
		expected := color.RGBA{count, count + 1, count + 2, 0x00}
		if p != expected {
			t.Errorf("Pixel %d value (%v) not expected (%v)", idx, p, expected)
		}
		count += 3
	}

	done = sr.ProcessFrame(startTime.Add(time.Millisecond))
	if !done {
		t.Error("Not done on second call to ProcessFrame")
	}
}

func TestTwoSimpleSequences(t *testing.T) {
	ta1 := testAnimation(1)
	s1 := Step{UniverseID: 3, Effect: &ta1}
	ta2 := testAnimation(31)
	s2 := Step{UniverseID: 1, Effect: &ta2}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 6, 1, 10})
	startTime := time.Unix(0, 0)
	sr.InitSequence(seq, startTime)

	// Check for right-sized empty buffer
	buf := sr.UniverseData(1)
	if len(buf) != 6 {
		t.Errorf("Unexpected buffer size %d", len(buf))
	}
	zero := color.RGBA{0, 0, 0, 0}
	for idx, p := range buf {
		if p != zero {
			t.Errorf("Pixel %d not 0 (%x)", idx, p)
		}
	}
	buf = sr.UniverseData(3)
	if len(buf) != 10 {
		t.Errorf("Unexpected buffer size %d", len(buf))
	}
	for idx, p := range buf {
		if p != zero {
			t.Errorf("Pixel %d not 0 (%x)", idx, p)
		}
	}

	// Call the effect to generate data
	done := sr.ProcessFrame(startTime)
	if done {
		t.Error("Done on first frame")
	}

	// Check for updated data
	buf = sr.UniverseData(1)
	count := uint8(31)
	for idx, p := range buf {
		expected := color.RGBA{count, count + 1, count + 2, 0x00}
		if p != expected {
			t.Errorf("Pixel %d value (%v) not expected (%v)", idx, p, expected)
		}
		count += 3
	}
	buf = sr.UniverseData(3)
	count = uint8(1)
	for idx, p := range buf {
		expected := color.RGBA{count, count + 1, count + 2, 0x00}
		if p != expected {
			t.Errorf("Pixel %d value (%v) not expected (%v)", idx, p, expected)
		}
		count += 3
	}

	done = sr.ProcessFrame(startTime.Add(time.Millisecond))
	if !done {
		t.Error("Not done on second call to ProcessFrame")
	}
}

func TestMultiRun(t *testing.T) {
	ta1 := multiRunAnimation(3)
	s1 := Step{UniverseID: 3, Effect: &ta1}
	ta2 := multiRunAnimation(1)
	s2 := Step{UniverseID: 0, Effect: &ta2}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 1, 1, 1})
	now := time.Unix(0, 0)
	sr.InitSequence(seq, now)

	if sr.ProcessFrame(now) {
		t.Error("Done on call 1")
	}
	if sr.ProcessFrame(now.Add(time.Millisecond)) {
		t.Error("Done on call 2")
	}
	if sr.ProcessFrame(now.Add(2 * time.Millisecond)) {
		t.Error("Done on call 3")
	}
	if !sr.ProcessFrame(now.Add(3 * time.Millisecond)) {
		t.Error("Not done on call 4")
	}
}

func TestDelay(t *testing.T) {
	ta1 := testAnimation(1)
	s1 := Step{UniverseID: 3, Effect: &ta1, Delay: 9 * time.Millisecond}
	ta2 := testAnimation(31)
	s2 := Step{UniverseID: 1, Effect: &ta2}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 6, 1, 10})
	now := time.Unix(0, 0)
	sr.InitSequence(seq, now)

	for tick := 0; tick < 10; tick++ {
		frameTime := now.Add(time.Duration(tick) * time.Millisecond)
		if sr.ProcessFrame(frameTime) {
			t.Errorf("Done on call for tick %d, time %v", tick, frameTime)
		}
		if sr.UniverseData(3)[0].R != 0 {
			t.Error("Delayed effect ran too early")
		}
		if sr.UniverseData(1)[0].R == 0 {
			t.Error("Immediate effect didn't run")
		}
	}

	if sr.ProcessFrame(now.Add(10 * time.Millisecond)) {
		t.Error("Done at time 10")
	}
	if sr.UniverseData(3)[0].R == 0 {
		t.Error("Delayed effect didn't run at right time")
	}

	if !sr.ProcessFrame(now.Add(11 * time.Millisecond)) {
		t.Error("Not done at time 11")
	}
}

func TestScheduleAfter(t *testing.T) {
	ta1 := testAnimation(1)
	ta2 := testAnimation(31)
	s2 := Step{UniverseID: 3, Effect: &ta2, StepID: 2}
	s1 := Step{UniverseID: 1, Effect: &ta1, OnCompletionOf: 2}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 1, 1, 1})
	now := time.Unix(0, 0)
	sr.InitSequence(seq, now)

	for tick := 0; tick < 1; tick++ {
		frameTime := now.Add(time.Duration(tick) * time.Millisecond)
		if sr.ProcessFrame(frameTime) {
			t.Errorf("Done on call for tick %d, time %v", tick, frameTime)
		}
		if sr.UniverseData(1)[0].R != 0 {
			t.Error("Contingent effect ran too early")
		}
		if sr.UniverseData(3)[0].R == 0 {
			t.Error("Immediate effect didn't run")
		}
	}

	if sr.ProcessFrame(now.Add(1 * time.Millisecond)) {
		t.Error("Done at time 1")
	}
	if sr.UniverseData(1)[0].R == 0 {
		t.Error("Contingent effect didn't run at right time")
	}

	if !sr.ProcessFrame(now.Add(2 * time.Millisecond)) {
		t.Error("Not done at time 2")
	}
}

func TestScheduleAfterPlusDelay(t *testing.T) {
	ta1 := testAnimation(1)
	ta2 := testAnimation(31)
	s2 := Step{UniverseID: 3, Effect: &ta2, StepID: 2}
	s1 := Step{UniverseID: 1, Effect: &ta1, OnCompletionOf: 2, Delay: time.Millisecond * 2}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 1, 1, 1})
	now := time.Unix(0, 0)
	sr.InitSequence(seq, now)

	for tick := 0; tick < 3; tick++ {
		frameTime := now.Add(time.Duration(tick) * time.Millisecond)
		if sr.ProcessFrame(frameTime) {
			t.Errorf("Done on call for tick %d, time %v", tick, frameTime)
		}
		if sr.UniverseData(1)[0].R != 0 {
			t.Error("Contingent effect ran too early")
		}
		if sr.UniverseData(3)[0].R == 0 {
			t.Error("Immediate effect didn't run")
		}
	}
	if sr.ProcessFrame(now.Add(3 * time.Millisecond)) {
		t.Error("Done at time 3")
	}
	if sr.UniverseData(1)[0].R == 0 {
		t.Error("Contingent effect didn't run at right time")
	}

	if !sr.ProcessFrame(now.Add(4 * time.Millisecond)) {
		t.Error("Not done at time 4")
	}
}

func TestScheduleAfterBadStepId(t *testing.T) {
	ta1 := testAnimation(1)
	ta2 := testAnimation(31)
	s2 := Step{UniverseID: 3, Effect: &ta2, StepID: 2}
	// Step waiting on invalid step should be ignored with warning
	s1 := Step{UniverseID: 1, Effect: &ta1, OnCompletionOf: 9}
	seq := Sequence{[]Step{s1, s2}}
	sr := NewSequenceRunner([]uint{1, 1, 1, 1})
	now := time.Unix(0, 0)
	sr.InitSequence(seq, now)

	for tick := 0; tick < 1; tick++ {
		frameTime := now.Add(time.Duration(tick) * time.Millisecond)
		if sr.ProcessFrame(frameTime) {
			t.Errorf("Done on call for tick %d, time %v", tick, frameTime)
		}
		if sr.UniverseData(1)[0].R != 0 {
			t.Error("Invalid contingent effect ran")
		}
		if sr.UniverseData(3)[0].R == 0 {
			t.Error("Immediate effect didn't run")
		}
	}

	if !sr.ProcessFrame(now.Add(1 * time.Millisecond)) {
		t.Error("Not done at time 1")
	}

	if sr.UniverseData(1)[0].R != 0 {
		t.Error("Contingent effect ran")
	}
}
