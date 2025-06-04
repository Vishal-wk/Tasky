package main

import (
	"embed"
	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/w32"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "tasky",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(&GreetService{}),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	app.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "tasky",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		//Frameless: true,
		Windows: application.WindowsWindow{
			ExStyle: w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOPMOST,
		},
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})
	 html := `
    <!DOCTYPE html>
    <html>
    <head>
        <link rel="stylesheet" href="https://unpkg.com/98.css" />
        <style>
            body {
                background-color: #c0c0c0;
                margin: 32px;
            }
            
            .window {
                width: 100%;
                max-width: 700px;
                margin: 0 auto;
            }
            
            .window-body {
                padding: 16px;
            }
            
            .button-row {
                display: flex;
                gap: 6px;
                margin-top: 16px;
            }
            
            .content-area {
                margin: 16px 0;
                background: white;
                padding: 8px;
                border: 2px inset #fff;
            }
        </style>
    </head>
    <body>
        <div class="window">
            <div class="title-bar">
                <div class="title-bar-text">My Windows 98 App</div>
                <div class="title-bar-controls">
                    <button aria-label="Minimize"></button>
                    <button aria-label="Maximize"></button>
                    <button aria-label="Close"></button>
                </div>
            </div>
            
            <div class="window-body">
                <menu role="tablist">
                    <li role="tab"><a href="#file">File</a></li>
                    <li role="tab"><a href="#edit">Edit</a></li>
                    <li role="tab"><a href="#help">Help</a></li>
                </menu>
                
                <div class="content-area">
                    <p>Welcome to my Windows 98 style application!</p>
                    
                    <fieldset>
                        <legend>Input Example</legend>
                        <div class="field-row">
                            <label for="text">Text:</label>
                            <input id="text" type="text" />
                        </div>
                    </fieldset>
                    
                    <fieldset>
                        <legend>Options</legend>
                        <div class="field-row">
                            <input id="check1" type="checkbox" />
                            <label for="check1">Option 1</label>
                        </div>
                        <div class="field-row">
                            <input id="check2" type="checkbox" />
                            <label for="check2">Option 2</label>
                        </div>
                    </fieldset>
                    
                    <div class="button-row">
                        <button>OK</button>
                        <button>Cancel</button>
                        <button>Apply</button>
                    </div>
                </div>
                
                <div class="status-bar">
                    <p class="status-bar-field">Status: Ready</p>
                    <p class="status-bar-field">CPU Usage: 0%</p>
                    <p class="status-bar-field">Memory: 32MB / 64MB</p>
                </div>
            </div>
        </div>

        <script>
            // Add some basic functionality
            document.querySelectorAll('button').forEach(button => {
                button.addEventListener('click', () => {
                    if (button.getAttribute('aria-label') === 'Close') {
                        // You might want to handle this in Go instead
                        window.close();
                    }
                });
            });
            
            // Example of updating status bar
            setInterval(() => {
                const cpu = Math.floor(Math.random() * 100);
                document.querySelector('.status-bar-field:nth-child(2)').textContent = 
                    'CPU Usage: ' + cpu + '%';
            }, 1000);
        </script>
    </body>
    </html>
    `
 
	// Create a goroutine that emits an event containing the current time every second.
	// The frontend can listen to this event and update the UI accordingly.
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
