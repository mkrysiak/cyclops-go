package hash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashWithFrames(t *testing.T) {
	// s, _ := New(jsonWithFrames)
	h, _ := HashForGrouping(jsonWithFrames)
	assert.Equal(t, h, "c8db52a12ce141f21a42903d8bf4c195")
}

func TestHashWithoutFrames(t *testing.T) {
	// s, _ := New(jsonWithoutFrames)
	h, _ := HashForGrouping(jsonWithoutFrames)
	assert.Equal(t, h, "b08c53afbb409576ea7a84de53e79873")
}

func TestHashWithEmptyException(t *testing.T) {
	// s, _ := New(unidentifyableException)
	h, _ := HashForGrouping(unidentifyableException)
	assert.Equal(t, h, "")
}

var jsonWithFrames = []byte(`
{
    "project": "4",
    "logger": "javascript",
    "platform": "javascript",
    "request": {
        "headers": {
            "User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/63.0.3239.84 Chrome/63.0.3239.84 Safari/537.36"
        },
        "url": "http://localhost:8000/app/"
    },
    "exception": {
        "values": [
            {
                "type": "Error",
                "value": "[$injector:modulerr] Failed to instantiate module phonecatApp due to:\nError: [$injector:modulerr] Failed to instantiate module ngRaven due to:\nError: [$injector:nomod] Module 'ngRaven' is not available! You either misspelled the module name or forgot to load it. If registering a module ensure that you specify the dependencies as the second argument.\nhttp://errors.angularjs.org/1.4.14/$injector/nomod?p0=ngRaven\n    at http://localhost:8000/app/bower_components/angular/angular.js:68:12\n    at http://localhost:8000/app/bower_components/angular/angular.js:2022:17\n    at ensure (http://localhost:8000/app/bower_components/angular/angular.js:1946:38)\n    at module (http://localhost:8000/app/bower_components/angular/angular.js:2020:14)\n    at http://localhost:8000/app/bower_components/angular/angular.js:4482:22\n    at forEach (http://localhost:8000/app/bower_components/angular/angular.js:341:20)\n    at loadModules (http://localhost:8000/app/bower_components/angular/angular.js:4466:5)\n    at http://localhost:8000/app/bower_components/angular/angular.js:4483:40\n    at forEach (http://localhost:8000/app/bower_components/angular/angular.js:341:20)\n    at loadModules (http://localhost:8000/app/bower_components/angular/angular.js:4466:5)\n    at createInjector (http://localhost:8000/app/bower_components/angular/angular.js:4391:11)\n    at doBootstrap (http://localhost:8000/app/bower_components/angular/angular.js:1693:20)\n    at bootstrap (http://localhost:8000/app/bower_components/angular/angular.js:1714:12)\n    at angularInit (http://localhost:8000/app/bower_components/angular/angular.js:1608:5)\n    at http://localhost:8000/app/bower_components/angular/angular.js:29847:5\n    at HTMLDocument.trigger (http://localhost:8000/app/bower_components/angular/angular.js:3095:7)\n    at defaultHandlerWrapper (http://localhost:8000/app/bower_components/angular/angular.js:3384:11)\n    at HTMLDocument.eventHandler (http://localhost:8000/app/bower_components/angular/angular.js:3372:9)\nhttp://errors.angularjs.org/1.4.14/$injector/modulerr?p0=ngRaven&p1=Error%3A%20%5B%24injector%3Anomod%5D%20Module%20'ngRaven'%20is%20not%20available!%20You%20either%20misspelled%20the%20module%20name%20or%20forgot%20to%20load%20it.%20If%20registering%20a%20module%20ensure%20that%20you%20specify%20the%20dependencies%20as%20the%20second%20argument.%0Ahttp%3A%2F%2Ferrors.angularjs.org%2F1.4.14%2F%24injector%2Fnomod%3Fp0%3DngRaven%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A68%3A12%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A2022%3A17%0A%20%20%20%20at%20ensure%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1946%3A38)%0A%20%20%20%20at%20module%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A2020%3A14)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4482%3A22%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4483%3A40%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20createInjector%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4391%3A11)%0A%20%20%20%20at%20doBootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1693%3A20)%0A%20%20%20%20at%20bootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1714%3A12)%0A%20%20%20%20at%20angularInit%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1608%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A29847%3A5%0A%20%20%20%20at%20HTMLDocument.trigger%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3095%3A7)%0A%20%20%20%20at%20defaultHandlerWrapper%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3384%3A11)%0A%20%20%20%20at%20HTMLDocument.eventHandler%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3372%3A9)\n    at http://localhost:8000/app/bower_components/angular/angular.js:68:12\n    at http://localhost:8000/app/bower_components/angular/angular.js:4505:15\n    at forEach (http://localhost:8000/app/bower_components/angular/angular.js:341:20)\n    at loadModules (http://localhost:8000/app/bower_components/angular/angular.js:4466:5)\n    at http://localhost:8000/app/bower_components/angular/angular.js:4483:40\n    at forEach (http://localhost:8000/app/bower_components/angular/angular.js:341:20)\n    at loadModules (http://localhost:8000/app/bower_components/angular/angular.js:4466:5)\n    at createInjector (http://localhost:8000/app/bower_components/angular/angular.js:4391:11)\n    at doBootstrap (http://localhost:8000/app/bower_components/angular/angular.js:1693:20)\n    at bootstrap (http://localhost:8000/app/bower_components/angular/angular.js:1714:12)\n    at angularInit (http://localhost:8000/app/bower_components/angular/angular.js:1608:5)\n    at http://localhost:8000/app/bower_components/angular/angular.js:29847:5\n    at HTMLDocument.trigger (http://localhost:8000/app/bower_components/angular/angular.js:3095:7)\n    at defaultHandlerWrapper (http://localhost:8000/app/bower_components/angular/angular.js:3384:11)\n    at HTMLDocument.eventHandler (http://localhost:8000/app/bower_components/angular/angular.js:3372:9)\nhttp://errors.angularjs.org/1.4.14/$injector/modulerr?p0=phonecatApp&p1=Error%3A%20%5B%24injector%3Amodulerr%5D%20Failed%20to%20instantiate%20module%20ngRaven%20due%20to%3A%0AError%3A%20%5B%24injector%3Anomod%5D%20Module%20'ngRaven'%20is%20not%20available!%20You%20either%20misspelled%20the%20module%20name%20or%20forgot%20to%20load%20it.%20If%20registering%20a%20module%20ensure%20that%20you%20specify%20the%20dependencies%20as%20the%20second%20argument.%0Ahttp%3A%2F%2Ferrors.angularjs.org%2F1.4.14%2F%24injector%2Fnomod%3Fp0%3DngRaven%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A68%3A12%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A2022%3A17%0A%20%20%20%20at%20ensure%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1946%3A38)%0A%20%20%20%20at%20module%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A2020%3A14)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4482%3A22%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4483%3A40%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20createInjector%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4391%3A11)%0A%20%20%20%20at%20doBootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1693%3A20)%0A%20%20%20%20at%20bootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1714%3A12)%0A%20%20%20%20at%20angularInit%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1608%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A29847%3A5%0A%20%20%20%20at%20HTMLDocument.trigger%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3095%3A7)%0A%20%20%20%20at%20defaultHandlerWrapper%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3384%3A11)%0A%20%20%20%20at%20HTMLDocument.eventHandler%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3372%3A9)%0Ahttp%3A%2F%2Ferrors.angularjs.org%2F1.4.14%2F%24injector%2Fmodulerr%3Fp0%3DngRaven%26p1%3DError%253A%2520%255B%2524injector%253Anomod%255D%2520Module%2520'ngRaven'%2520is%2520not%2520available!%2520You%2520either%2520misspelled%2520the%2520module%2520name%2520or%2520forgot%2520to%2520load%2520it.%2520If%2520registering%2520a%2520module%2520ensure%2520that%2520you%2520specify%2520the%2520dependencies%2520as%2520the%2520second%2520argument.%250Ahttp%253A%252F%252Ferrors.angularjs.org%252F1.4.14%252F%2524injector%252Fnomod%253Fp0%253DngRaven%250A%2520%2520%2520%2520at%2520http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A68%253A12%250A%2520%2520%2520%2520at%2520http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A2022%253A17%250A%2520%2520%2520%2520at%2520ensure%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A1946%253A38)%250A%2520%2520%2520%2520at%2520module%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A2020%253A14)%250A%2520%2520%2520%2520at%2520http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A4482%253A22%250A%2520%2520%2520%2520at%2520forEach%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A341%253A20)%250A%2520%2520%2520%2520at%2520loadModules%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A4466%253A5)%250A%2520%2520%2520%2520at%2520http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A4483%253A40%250A%2520%2520%2520%2520at%2520forEach%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A341%253A20)%250A%2520%2520%2520%2520at%2520loadModules%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A4466%253A5)%250A%2520%2520%2520%2520at%2520createInjector%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A4391%253A11)%250A%2520%2520%2520%2520at%2520doBootstrap%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A1693%253A20)%250A%2520%2520%2520%2520at%2520bootstrap%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A1714%253A12)%250A%2520%2520%2520%2520at%2520angularInit%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A1608%253A5)%250A%2520%2520%2520%2520at%2520http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A29847%253A5%250A%2520%2520%2520%2520at%2520HTMLDocument.trigger%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A3095%253A7)%250A%2520%2520%2520%2520at%2520defaultHandlerWrapper%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A3384%253A11)%250A%2520%2520%2520%2520at%2520HTMLDocument.eventHandler%2520(http%253A%252F%252Flocalhost%253A8000%252Fapp%252Fbower_components%252Fangular%252Fangular.js%253A3372%253A9)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A68%3A12%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4505%3A15%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4483%3A40%0A%20%20%20%20at%20forEach%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A341%3A20)%0A%20%20%20%20at%20loadModules%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4466%3A5)%0A%20%20%20%20at%20createInjector%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A4391%3A11)%0A%20%20%20%20at%20doBootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1693%3A20)%0A%20%20%20%20at%20bootstrap%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1714%3A12)%0A%20%20%20%20at%20angularInit%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A1608%3A5)%0A%20%20%20%20at%20http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A29847%3A5%0A%20%20%20%20at%20HTMLDocument.trigger%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3095%3A7)%0A%20%20%20%20at%20defaultHandlerWrapper%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3384%3A11)%0A%20%20%20%20at%20HTMLDocument.eventHandler%20(http%3A%2F%2Flocalhost%3A8000%2Fapp%2Fbower_components%2Fangular%2Fangular.js%3A3372%3A9)",
                "stacktrace": {
                    "frames": [
                        {
                            "filename": "http://localhost:8000/app/bower_components/angular/angular.js",
                            "lineno": 3372,
                            "colno": 9,
                            "function": "HTMLDocument.eventHandler",
                            "in_app": true
                        },
                        {
                            "filename": "http://localhost:8000/app/bower_components/angular/angular.js",
                            "lineno": 3384,
                            "colno": 11,
                            "function": "defaultHandlerWrapper",
                            "in_app": true
                        }
                    ]
                }
            }
        ]
    },
    "culprit": "http://errors.angularjs.org/1.4.14/$injector/nomod?p0=ngRaven",
    "extra": {
        "session:duration": 16
    },
    "event_id": "29320acdb67d4f26b3dc78446ff95456"
}`)

var jsonWithoutFrames = []byte(`
	{
		"project": "4",
		"logger": "javascript",
		"platform": "javascript",
		"request": {
			"headers": {
				"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/63.0.3239.84 Chrome/63.0.3239.84 Safari/537.36"
			},
			"url": "http://localhost:8000/app/"
		},
		"culprit": "http://errors.angularjs.org/1.4.14/$injector/nomod?p0=ngRaven",
		"extra": {
			"session:duration": 16
		},
		"event_id": "29320acdb67d4f26b3dc78446ff95456"
	}`)

var unidentifyableException = []byte(`
	{
		"project": "4",
		"logger": "javascript",
		"platform": "javascript",
		"request": {
			"headers": {
				"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/63.0.3239.84 Chrome/63.0.3239.84 Safari/537.36"
			},
			"url": "http://localhost:8000/app/"
		},
		"extra": {
			"session:duration": 16
		},
		"event_id": "29320acdb67d4f26b3dc78446ff95456"
	}`)
