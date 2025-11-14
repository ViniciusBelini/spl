package modules

import(
	// "fmt"
	// "reflect"

	// "SPL/ast"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	// "fyne.io/fyne/v2/widget"
)

func GUI_run(params interface{}) (interface{}, error){
	paramsI := params.([]interface{})

	switch paramsI[0].(string){
		case "new_window":
			return GUI_new_window(paramsI[1].(string), paramsI[2].(int), paramsI[3].(int))
		case "set_title":
			return GUI_set_title(paramsI[1].(fyne.Window), paramsI[2].(string))
		case "set_size":
			return GUI_set_size(paramsI[1].(fyne.Window), paramsI[2].(int), paramsI[3].(int))
		case "show":
			return GUI_show(paramsI[1].(fyne.Window))
	}

	return nil, nil
}

func GUI_new_window(title string, width int, height int) (interface{}, error){
	newApp := app.New()
	window := newApp.NewWindow(title) // fyne.Window

	GUI_set_title(window, title)
	GUI_set_size(window, width, height)

	return window, nil
}

func GUI_set_title(window fyne.Window, title string) (interface{}, error){
	window.SetTitle(title)

	return true, nil
}

func GUI_set_size(window fyne.Window, width int, height int) (interface{}, error){
	window.Resize(fyne.NewSize(float32(width), float32(height)))

	return true, nil
}

func GUI_show(window fyne.Window) (interface{}, error){
	window.ShowAndRun()

	return true, nil
}
