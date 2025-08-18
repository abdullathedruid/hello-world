package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type PageData struct {
	Title string
	Time  string
}

func main() {
	// Configure logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Create router
	r := mux.NewRouter()

	// Routes
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/debug", debugHandler).Methods("GET")
	r.HandleFunc("/time", timeHandler).Methods("GET")
	r.HandleFunc("/click", clickHandler).Methods("POST")

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	logrus.Info("Server starting on http://localhost:8080")
	logrus.Fatal(http.ListenAndServe(":8080", r))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>HTMX + Go</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-2xl">
        <h1 class="text-4xl font-bold text-center text-gray-800 mb-8">HTMX + Go Demo</h1>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Current Time</h2>
            <div id="time-display">
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
                        hx-get="/time" hx-target="#time-display">Get Current Time</button>
            </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Click Counter</h2>
            <div id="click-counter">
                <button class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
                        hx-post="/click" hx-target="#click-counter">Click Me!</button>
            </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">🚀 Farcaster Debug</h2>
            <p class="text-gray-600 mb-4">Debug Farcaster MiniApp SDK values and context</p>
            <a href="/debug" class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded transition duration-200 inline-block">
                Open Debug Page
            </a>
        </div>
    </div>

    <script type="module">
        import { sdk } from 'https://esm.sh/@farcaster/miniapp-sdk'
        
        // After your app is fully loaded and ready to display
        await sdk.actions.ready()
    </script>
</body>
</html>`

	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Time endpoint accessed")
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, `<div class="bg-blue-50 border border-blue-200 rounded-md p-4">
		<p class="text-gray-700">Current time: <strong class="text-blue-600">%s</strong></p>
		<button class="mt-3 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-get="/time" hx-target="#time-display">Refresh Time</button>
	</div>`, currentTime)
}

var clickCount = 0

func clickHandler(w http.ResponseWriter, r *http.Request) {
	clickCount++
	logrus.WithField("count", clickCount).Info("Button clicked")
	fmt.Fprintf(w, `<div class="bg-green-50 border border-green-200 rounded-md p-4">
		<p class="text-gray-700">Button clicked <strong class="text-green-600">%d</strong> times!</p>
		<button class="mt-3 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-post="/click" hx-target="#click-counter">Click Me Again!</button>
	</div>`, clickCount)
}

func debugHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Debug page accessed")
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Farcaster MiniApp Debug</title>
    <meta name="fc:miniapp" content="Farcaster MiniApp Debug Page">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-4xl">
        <h1 class="text-4xl font-bold text-center text-blue-600 mb-8">🚀 Farcaster MiniApp Debug</h1>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Context Information -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">📊 Context Information</h2>
                <div id="context-info" class="space-y-2 text-sm">
                    <p class="text-gray-500">Loading context...</p>
                </div>
            </div>
            
            <!-- User Information -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">👤 User Information</h2>
                <div id="user-info" class="space-y-2 text-sm">
                    <p class="text-gray-500">Loading user info...</p>
                </div>
            </div>
            
            <!-- Location & Platform -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">📍 Location & Platform</h2>
                <div id="location-info" class="space-y-2 text-sm">
                    <p class="text-gray-500">Loading location info...</p>
                </div>
            </div>
            
            <!-- SDK Status -->
            <div class="bg-white rounded-lg shadow-md p-6">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">⚙️ SDK Status</h2>
                <div id="sdk-status" class="space-y-2 text-sm">
                    <p class="text-gray-500">Initializing SDK...</p>
                </div>
            </div>
            
            <!-- Actions & Methods -->
            <div class="bg-white rounded-lg shadow-md p-6 md:col-span-2">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">🎯 Available Actions</h2>
                <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                    <button id="refresh-btn" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200">
                        🔄 Refresh Data
                    </button>
                    <button id="close-btn" class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded transition duration-200">
                        ❌ Close MiniApp
                    </button>
                    <button id="external-link-btn" class="bg-purple-500 hover:bg-purple-700 text-white font-bold py-2 px-4 rounded transition duration-200">
                        🔗 Open External Link
                    </button>
                </div>
                <div id="action-result" class="mt-4"></div>
            </div>
            
            <!-- Raw Data -->
            <div class="bg-white rounded-lg shadow-md p-6 md:col-span-2">
                <h2 class="text-2xl font-semibold text-gray-700 mb-4">🔍 Raw SDK Data</h2>
                <pre id="raw-data" class="bg-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
Loading SDK data...
                </pre>
            </div>
        </div>
        
        <div class="text-center mt-8">
            <a href="/" class="text-blue-500 hover:text-blue-700 font-semibold">← Back to Home</a>
        </div>
    </div>

    <script type="module">
        import { sdk } from 'https://esm.sh/@farcaster/miniapp-sdk'
        
        let sdkReady = false;
        let contextData = null;
        
        async function initializeSDK() {
            try {
                // Mark app as ready
                await sdk.actions.ready();
                sdkReady = true;
                
                // Get context information
                contextData = await sdk.context;
                
                updateDisplay();
                
                document.getElementById('sdk-status').innerHTML = '<div class="text-green-600 font-semibold">✅ SDK Ready</div>';
                
            } catch (error) {
                console.error('SDK initialization error:', error);
                document.getElementById('sdk-status').innerHTML = '<div class="text-red-600 font-semibold">❌ SDK Error: ' + error.message + '</div>';
            }
        }
        
        function updateDisplay() {
            if (!contextData) {
                document.getElementById('raw-data').textContent = 'No context data available';
                return;
            }
            
            // Update context info
            const contextHtml = Object.entries(contextData).map(([key, value]) => {
                if (key === 'user') return ''; // Handle user separately
                return '<div class="flex justify-between border-b pb-1"><span class="font-medium text-gray-600">' + key + ':</span><span class="text-gray-800">' + JSON.stringify(value) + '</span></div>';
            }).filter(html => html).join('');
            document.getElementById('context-info').innerHTML = contextHtml || '<p class="text-gray-500">No context data</p>';
            
            // Update user info
            if (contextData.user) {
                const userHtml = Object.entries(contextData.user).map(([key, value]) => 
                    '<div class="flex justify-between border-b pb-1"><span class="font-medium text-gray-600">' + key + ':</span><span class="text-gray-800">' + JSON.stringify(value) + '</span></div>'
                ).join('');
                document.getElementById('user-info').innerHTML = userHtml;
            } else {
                document.getElementById('user-info').innerHTML = '<p class="text-gray-500">No user data available</p>';
            }
            
            // Update location info
            const locationData = {
                'URL': window.location.href,
                'User Agent': navigator.userAgent,
                'Platform': navigator.platform,
                'Language': navigator.language,
                'Screen': screen.width + 'x' + screen.height,
                'Viewport': window.innerWidth + 'x' + window.innerHeight
            };
            
            const locationHtml = Object.entries(locationData).map(([key, value]) => 
                '<div class="flex justify-between border-b pb-1"><span class="font-medium text-gray-600">' + key + ':</span><span class="text-gray-800 text-xs">' + value + '</span></div>'
            ).join('');
            document.getElementById('location-info').innerHTML = locationHtml;
            
            // Update raw data
            document.getElementById('raw-data').textContent = JSON.stringify(contextData, null, 2);
        }
        
        // Event handlers
        document.getElementById('refresh-btn').addEventListener('click', async () => {
            document.getElementById('action-result').innerHTML = '<div class="bg-blue-50 border border-blue-200 rounded p-2 text-sm">Refreshing data...</div>';
            try {
                contextData = await sdk.context;
                updateDisplay();
                document.getElementById('action-result').innerHTML = '<div class="bg-green-50 border border-green-200 rounded p-2 text-sm text-green-700">✅ Data refreshed successfully</div>';
            } catch (error) {
                document.getElementById('action-result').innerHTML = '<div class="bg-red-50 border border-red-200 rounded p-2 text-sm text-red-700">❌ Refresh failed: ' + error.message + '</div>';
            }
        });
        
        document.getElementById('close-btn').addEventListener('click', async () => {
            try {
                await sdk.actions.close();
                document.getElementById('action-result').innerHTML = '<div class="bg-yellow-50 border border-yellow-200 rounded p-2 text-sm text-yellow-700">⚠️ Close action sent</div>';
            } catch (error) {
                document.getElementById('action-result').innerHTML = '<div class="bg-red-50 border border-red-200 rounded p-2 text-sm text-red-700">❌ Close failed: ' + error.message + '</div>';
            }
        });
        
        document.getElementById('external-link-btn').addEventListener('click', async () => {
            try {
                await sdk.actions.openUrl('https://docs.farcaster.xyz');
                document.getElementById('action-result').innerHTML = '<div class="bg-green-50 border border-green-200 rounded p-2 text-sm text-green-700">✅ External link opened</div>';
            } catch (error) {
                document.getElementById('action-result').innerHTML = '<div class="bg-red-50 border border-red-200 rounded p-2 text-sm text-red-700">❌ Open URL failed: ' + error.message + '</div>';
            }
        });
        
        // Initialize on load
        initializeSDK();
    </script>
</body>
</html>`
	
	t, _ := template.New("debug").Parse(tmpl)
	t.Execute(w, nil)
}
