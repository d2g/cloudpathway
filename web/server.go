package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/d2g/controller"
	"github.com/d2g/sessions"
	"github.com/d2g/sessions/boltsessionstore"
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

		b, err := bolt.Open("userdata/Sessions.boltdb", 0666, nil)
		if err != nil {
			return nil, err
		}
		sessionstore := &boltsessionstore.BoltStore{
			DB: b,
		}

		sessioninfo := sessions.SessionInfo{}
		sessioninfo.Timeout = (time.Duration(configuration.Sessions.Timeout) * time.Second)
		sessioninfo.Store = sessionstore
		sessioninfo.Cookie.Name = "GOSESSION"
		sessioninfo.Cache = sessions.RequestSessions{make([]sessions.RequestSession, 0)}

		activeServer.sessions = &ActiveSessions{sessioninfo}

		//Start Session GC
		go activeServer.sessions.GCSessions()

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
			NotFoundHandler: http.Handler(NotFound{NotFoundHandler: http.FileServer(http.Dir(t.configuration.Files))}),
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
