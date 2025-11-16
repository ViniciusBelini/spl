package interpreter

import(
	"fmt"
	"os"
	// "reflect"

	// "SPL/models"
	// "SPL/ast"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
)

func module_GUI_run(params interface{}) (interface{}, error){
	// a := app.New()
	// w := a.NewWindow("Título inicial")
	// b := container.NewVBox()
 //
	// btn := widget.NewButton("Mudar título e conteúdo", func() {
	// 	w.SetTitle("Título modificado dinamicamente")
	// 	w.Resize(fyne.NewSize(400, 200))
	// 	label := widget.NewLabel("Texto inicial")
	// 	b.Add(label)
	// })
 //
	// b.Add(btn)
 //
	// w.SetContent(b)
	// w.ShowAndRun()
 //
	// return nil, nil

	paramsI := params.([]interface{})

	switch paramsI[0].(string){
		case "new_window":
			return module_GUI_new_window(paramsI[1].(string), paramsI[2].(int), paramsI[3].(int))
		case "set_title":
			return module_GUI_set_title(paramsI[1], paramsI[2].(string))
		case "set_size":
			return module_GUI_set_size(paramsI[1], paramsI[2].(int), paramsI[3].(int))
		case "show":
			return module_GUI_show(paramsI[1])
		case "new_container":
			return module_GUI_new_container(paramsI[1].(string)), nil
		case "set_content":
			module_GUI_set_content(paramsI[1], paramsI[2])
			return nil, nil
		case "widget":
			return module_GUI_widget(paramsI)
		case "container_add":
			module_GUI_container_add(paramsI[1], paramsI[2])
	}

	return nil, nil
}

func module_GUI_new_window(title string, width int, height int) (interface{}, error){
	newApp := app.New()
	windowN := newApp.NewWindow(title) // fyne.Window

	module_GUI_set_title(windowN, title)
	module_GUI_set_size(windowN, width, height)

	return windowN, nil
}

func module_GUI_set_title(windowN interface{}, title string) (interface{}, error){
	switch windowN.(type){
		case fyne.Window:
			windowN.(fyne.Window).SetTitle(title)
	}

	return true, nil
}

func module_GUI_set_size(windowN interface{}, width int, height int) (interface{}, error){
	switch windowN.(type){
		case *fyne.Container:
			windowN.(*fyne.Container).Resize(fyne.NewSize(float32(width), float32(height)))
		// case fyne.CanvasObject:
			// windowN.(fyne.CanvasObject).Resize(fyne.NewSize(float32(width), float32(height)))
		case fyne.Widget:
			windowN.(fyne.Widget).Resize(fyne.NewSize(float32(width), float32(height)))
		case fyne.Window:
			windowN.(fyne.Window).Resize(fyne.NewSize(float32(width), float32(height)))
	}

	return true, nil
}

func module_GUI_show(windowN interface{}) (interface{}, error){
	switch windowN.(type){
		case fyne.Window:
			windowN.(fyne.Window).ShowAndRun()
	}

	return true, nil
}

func module_GUI_new_container(typeN string) interface{}{
	if typeN == "vertical"{
		return container.NewVBox() // *fyne.Container
	}else if typeN == "horizontal"{
		return container.NewHBox() // *fyne.Container
	}
	return nil
}
func module_GUI_container_add(containerN interface{}, widgetN interface{}){
	switch widgetN.(type){
		case *fyne.Container:
			switch containerN.(type){
				// case fyne.Window:
					// containerN.(fyne.Window).Add(widgetN.(*fyne.Container))
				case *fyne.Container:
					containerN.(*fyne.Container).Add(widgetN.(*fyne.Container))
			}
		case fyne.Widget:
			switch containerN.(type){
				case *fyne.Container:
					containerN.(*fyne.Container).Add(widgetN.(fyne.Widget))
			}
	}
}

func module_GUI_set_content(windowN interface{}, containerN interface{}){
	switch containerN.(type){
		case *fyne.Container:
			switch windowN.(type){
				case fyne.Window:
					windowN.(fyne.Window).SetContent(container.NewMax(containerN.(*fyne.Container)))
			}
		case fyne.Widget:
			switch windowN.(type){
				case fyne.Window:
					windowN.(fyne.Window).SetContent(container.NewMax(containerN.(fyne.Widget)))
			}
	}
}

// WIDGETS =====================================

func module_GUI_widget(params []interface{}) (interface{}, error){
	switch params[1].(string){
		case "new_button":
			funcP := params[3]
			if fnc, ok := funcP.([2]any);ok{
				funcP = fnc[0]
			}
			return module_GUI_new_button(params[2].(string), funcP.(*Func), 1, 1)
		case "new_label":
			return widget.NewLabel(params[2].(string)), nil
		case "new_input":
			input := widget.NewEntry()
			input.SetPlaceHolder(params[2].(string))

			return input, nil
		case "get_text":
			switch w := params[2].(type){
				case *widget.Label:
					return w.Text, nil
				case *widget.Entry:
					return w.Text, nil
				case *widget.Button:
					return w.Text, nil
			}
		case "set_text":
			switch w := params[2].(type){
				case *widget.Label:
					w.SetText(params[3].(string))
				case *widget.Entry:
					w.SetText(params[3].(string))
				case *widget.Button:
					w.SetText(params[3].(string))
			}
	}

	return nil, nil
}

func module_GUI_new_button(label string, function *Func, line int, pos int) (interface{}, error){
	btn := widget.NewButton(label, func() {
		if function.Outer == nil{
			return
		}

		fileName := function.FileName

		_, err := CallFunc("null", function, nil, function.Outer, fileName, line, pos)
		if err != nil{
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	})

	return btn, nil
}
