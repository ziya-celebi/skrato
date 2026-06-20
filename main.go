package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SkratoApp struct {
	scanner  *Scanner
	jobs     *JobManager

	statusLabel *widget.Label
	lastLabel   *widget.Label
	detectedList *widget.Label

	status     string
	lastResult *JobResult
	runID      uint64
}

func NewSkratoApp(_ fyne.App) *SkratoApp {
	scanner := NewScanner()
	jobs := NewJobManager()

	a := &SkratoApp{
		scanner:    scanner,
		jobs:       jobs,
		status:     "Ready.",
		lastResult: nil,
		runID:      0,
	}

	a.scanner.ScanActions()
	a.status = a.scanner.Status
	return a
}

func (a *SkratoApp) rescan() {
	if a.jobs.Running() {
		return
	}
	a.scanner.ScanActions()
	a.status = a.scanner.Status
	if a.statusLabel != nil {
		a.statusLabel.SetText(a.status)
	}
	if a.detectedList != nil {
		if len(a.scanner.Detected) > 0 {
			a.detectedList.SetText(joinStrings(a.scanner.Detected, ", "))
		} else {
			a.detectedList.SetText("none")
		}
	}
}

func (a *SkratoApp) startAction(action Action) {
	if a.jobs.Running() {
		return
	}
	a.runID++
	a.jobs.Start(action, a.runID)
	a.status = action.Label + "..."
	if a.statusLabel != nil {
		a.statusLabel.SetText(a.status)
	}
	a.lastResult = nil
	if a.lastLabel != nil {
		a.lastLabel.SetText("Last run: (none)")
	}
}

func (a *SkratoApp) pollResults() {
	for {
		res := a.jobs.TryRecv()
		if res == nil {
			break
		}
		if res.RunID != a.runID {
			continue
		}

		if res.OK {
			a.status = res.Label + " finished successfully."
		} else {
			ec := "nil"
			if res.ExitCode != nil {
				ec = string(rune(*res.ExitCode))
			}
			a.status = res.Label + " failed (exit code: " + ec + ")."
		}

		a.lastResult = res
		a.jobs.SetRunning(false)
	}
}

func (a *SkratoApp) BuildUI() fyne.CanvasObject {
	heading := widget.NewLabel("skrato")
	heading.TextStyle = fyne.TextStyle{Bold: true}
	subtitle := widget.NewLabel("Bootloader and initramfs maintenance")

	a.statusLabel = widget.NewLabel(a.status)

	lastBindingText := "Last run: (none)"
	a.lastLabel = widget.NewLabel(lastBindingText)

	detectedLabel := widget.NewLabel("Detected tools:")
	a.detectedList = widget.NewLabel("")
	if len(a.scanner.Detected) > 0 {
		a.detectedList.SetText(joinStrings(a.scanner.Detected, ", "))
	} else {
		a.detectedList.SetText("none")
	}

	separator := widget.NewSeparator()

	rescanBtn := widget.NewButton("Rescan", func() {
		a.rescan()
	})
	rescanBtn.Resize(fyne.NewSize(180, 38))

	actionContainer := container.NewVBox()
	for _, action := range a.scanner.Actions {
		act := action
		btn := widget.NewButton(act.Label, func() {
			a.startAction(act)
		})
		btn.Resize(fyne.NewSize(180, 38))
		actionContainer.Add(btn)
	}
	actionContainer.Add(rescanBtn)

	content := container.NewVBox(
		heading,
		subtitle,
		separator,
		a.statusLabel,
		a.lastLabel,
		separator,
		detectedLabel,
		a.detectedList,
		separator,
		actionContainer,
	)

	go func() {
		for {
			if a.jobs.Running() {
				a.pollResults()
				if a.lastResult != nil {
					res := a.lastResult
					ec := "nil"
					if res.ExitCode != nil {
						ec = string(rune(*res.ExitCode))
					}
					okText := "success"
					if !res.OK {
						okText = "failed"
					}
					a.lastLabel.SetText("Last run: " + okText + " (exit code: " + ec + ")")
					a.statusLabel.SetText(a.status)
					a.lastResult = nil
				}
			}
		}
	}()

	return content
}

func joinStrings(parts []string, sep string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += sep + parts[i]
	}
	return result
}

func main() {
	myApp := app.New()
	window := myApp.NewWindow("skrato")
	window.Resize(fyne.NewSize(400, 500))

	skratoApp := NewSkratoApp(myApp)
	window.SetContent(skratoApp.BuildUI())

	window.ShowAndRun()
}