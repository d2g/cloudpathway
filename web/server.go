package web

import (
	"errors"
	"github.com/d2g/controller"
	"github.com/d2g/sessions"
	"net/http"
	"time"
)

type server struct {
	configuration *Configuration
	controllers   controller.HTTPControllers

	sessions *ActiveSessions
}

var activeServer *server = nil

func NewServer(configuration *Configuration) (*server, error) {
	if activeServer != nil {
		return nil, errors.New("Local server has already been created. (Use Server())")
	} else {
		activeServer = new(server)
		activeServer.configuration = configuration

		sessionstore, err := sessions.NewUnqliteStore("Sessions.unqlite")
		if err != nil {
			return nil, err
		}

		sessioninfo := sessions.SessionInfo{}
		sessioninfo.Timeout = (4 * time.Hour)
		sessioninfo.Store = sessionstore
		sessioninfo.Cookie.Name = "GOSESSION"
		sessioninfo.Cache = sessions.RequestSessions{make([]sessions.RequestSession, 0)}

		activeServer.sessions = &ActiveSessions{sessioninfo}

		controllers := activeServer.Controllers()

		http.Handle("/", activeServer.sessions.GetHandler(controllers.Routes()))

		return activeServer, nil
	}
}

/*
 * It's important that this is a function of the server so that you can pass config options,
 * sessions and that sort of information to the child controller.
 */
func (t *server) Controllers() controller.HTTPControllers {
	/*
	 * If we haven't initialised  the controllers initialise them.
	 */
	if t.controllers == nil {
		/*
		 * Define the controllers here.
		 */
		t.controllers = controller.HTTPControllers{
			&Authentication{
				HTTPController: controller.HTTPController(&UsersCollections{
					Sessions: t.sessions,
				}).SetBase("/useraccess/"),
				Sessions: t.sessions,
			},
			&Authentication{
				HTTPController: controller.HTTPController(&Collection{
					Sessions: t.sessions,
				}).SetBase("/collections/"),
				Sessions: t.sessions,
			},
			&Authentication{
				HTTPController: controller.HTTPController(&User{
					Sessions: t.sessions,
				}).SetBase("/users/"),
				Sessions: t.sessions,
			},
			&Authentication{
				HTTPController: controller.HTTPController(&Device{
					Sessions: t.sessions,
				}).SetBase("/devices/"),
				Sessions: t.sessions,
			},
			&Authentication{
				HTTPController: controller.HTTPController(&Internet{
					Sessions: t.sessions,
				}).SetBase("/internet/"),
				Sessions: t.sessions,
			},
		}

		if t.configuration.DeveloperMode {
			t.controllers = append(t.controllers, controller.HTTPController(&Developer{
				Sessions: t.sessions,
			}).SetBase("/developer/"))
		}

		t.controllers = append(t.controllers, controller.HTTPController(&Root{
			Sessions:        t.sessions,
			NotFoundHandler: http.FileServer(http.Dir(t.configuration.Files)),
		}).SetBase("/"))
	}

	return t.controllers
}

/*
 *
 */
func (t *server) ListenAndServe() error {
	return http.ListenAndServe(":"+t.configuration.Port, nil)
}
