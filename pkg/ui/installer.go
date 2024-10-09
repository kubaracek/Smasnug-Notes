package ui

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

// Installer function type to simulate installation logic
type installerFn func() error

// InstallDialog creates a dialog window for installing Samsung Notes
// Returns (bool, error) where bool indicates success and error indicates any failure during installation
func InstallDialog(installer installerFn) (bool, error) {
	// Create a new window
	w := new(app.Window)
	w.Option(app.Title("Samsung Notes Installer"), app.Size(unit.Dp(600), unit.Dp(400)))

	// Run the installation dialog
	success, err := runInstallDialog(w, installer)
	if err != nil {
		return false, err
	}
	
	return success, nil
}

// UI Context and Dimensions shorthand
type C = layout.Context
type D = layout.Dimensions

// runInstallDialog handles the rendering and event loop for the installation dialog
func runInstallDialog(w *app.Window, installer installerFn) (bool, error) {
	// UI operations
	var ops op.Ops

	// UI elements
	var installButton widget.Clickable

	// UI state flags
	state := "idle"     // "idle", "installing", "done"
	installing := false // whether installation is ongoing

	// Channel to signal when installation is done
	doneChan := make(chan error)

	// Material theme
	th := material.NewTheme()

	for {
		select {
		// Installation completed
		case err := <-doneChan:
			if err != nil {
				// Installation failed
				return false, err
			}
			// Installation successful, show "Done" state
			state = "done"

		// Main UI event loop
		default:
			// Handle window events
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				// If window is closed, return false without error
				return false, nil

			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)

				// Handle button click
				if installButton.Clicked(gtx) {
					if state == "idle" && !installing {
						// Start installation
						state = "installing"
						installing = true

						// Start installation in background
						go func() {
							if err := installer(); err != nil {
								doneChan <- err
								return
							}
							// Notify installation completion
							doneChan <- nil
						}()
					} else if state == "done" {
						// Installation is done, close the dialog
						return true, nil
					}
				}

				// Render UI
				layout.Flex{
					Axis: layout.Vertical,
				}.Layout(gtx,
					// Title
					layout.Rigid(func(gtx C) D {
						title := material.H2(th, "Samsung Notes Installer")
						title.Alignment = text.Middle
						return title.Layout(gtx)
					}),
					// Body text
					layout.Rigid(func(gtx C) D {
						margins := layout.Inset{
							Top:   unit.Dp(25),
							Right: unit.Dp(35),
							Left:  unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx C) D {
							body := material.Body1(th, "Samsung Notes is not installed on your system.\n"+
								"Click 'Install Samsung Notes' to automatically install the app.\n"+
								"This may take a minute.\n\n"+
								"Alternatively, you can manually install it from the Microsoft Store.\n"+
								"Once installed, relaunch this app.")
							body.Alignment = text.Middle
							return body.Layout(gtx)
						})
					}),
					// Install/Done button
					layout.Rigid(func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx, func(gtx C) D {
							var btnText string
							switch state {
							case "idle":
								btnText = "Install Samsung Notes"
							case "installing":
								btnText = "Installing... Please wait"
							case "done":
								btnText = "Samsung Notes Installed, Click to launch!"
							}
							btn := material.Button(th, &installButton, btnText)
							return btn.Layout(gtx)
						})
					}),
				)
				e.Frame(gtx.Ops)
			}
		}
	}
}
