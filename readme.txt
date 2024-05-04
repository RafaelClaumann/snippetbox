Chapter 2 Foundations
    2.1 Project setup and creating a module
            go mod init snippetbox.claumann.net

            use 'go run main.go' command in your terminal to compile and execute the code
    
    2.2 Web application basics
            You can think of handlers as being a bit like controllers. 
            They’re responsible for executing your application logic and for writing HTTP response headers and bodies.

            The second component is a router (or servemux in Go terminology).
            This stores a mapping between the URL patterns for your application and the corresponding handlers

            http.ResponseWriter parameter provides methods for assembling a HTTP response and sending it to the user.
            *http.Request parameter is a pointer to a struct which holds information about the current request.

            Go’s servemux treats the URL pattern "/" like a catch-all.
            So at the moment all HTTP requests to our server will be handled by the home function, regardless of their URL path.

            During development the go run command is a convenient way to try out your code.
            It’s essentially a shortcut that compiles your code, creates an executable binary in your /tmp directory,
            and then runs this binary in one step.

    2.3 Routing requests

            Go’s servemux supports two different types of URL patterns: fixed paths and subtree paths.
            Fixed paths don’t end with a trailing slash, whereas subtree paths do end with a trailing slash.

            Fixed path patterns are only matched when the request URL path exactly matches the fixed path.
            Subtree path patterns are matched whenever the start of a request URL path matches the subtree path.
            You can think of subtree paths as acting a bit like they have a wildcard at the end, like "/**" or "/static/**".
        
        Restricting the root url pattern

            For instance, in the application we’re building we want the home page to be displayed
            if — and only if —the request URL path exactly matches "/".

            It’s not possible to change the behavior of Go’s servemux to do this, but you can include a simple check in the home handler.

        The DefaultServeMux

            The http.Handle() and http.HandleFunc() allow you to register routes without declaring a servemux.
                
                http.HandleFunc("/", home)
                http.HandleFunc("/snippet/view", snippetView)
                http.HandleFunc("/snippet/create", snippetCreate)

            Behind the scenes, these functions register their routes with something called the DefaultServeMux.

            DefaultServeMux is a global variable, any package can access it and register a route — including any third-party packages
            that your application imports. If one of those third-party packages is compromised, they could use DefaultServeMux to 
            expose a malicious handler to the web.

            It’s generally a good idea to avoid DefaultServeMux and the corresponding helper functions.
            Use your own locally-scoped servemux instead
        
        Additional information

            In Go’s servemux, longer URL patterns always take precedence over shorter ones.
            It will always dispatch the request to the handler corresponding to the longest pattern.
            You can register patterns in any order and it won’t change how the servemux behaves.

            It’s possible to include host names in your URL patterns.
                mux.HandleFunc("foo.example.org/", fooHandler)
                mux.HandleFunc("bar.example.org/", barHandler)
                mux.HandleFunc("/baz", bazHandler)
            Only when there isn’t a host-specific match found will the non-host specific patterns also be checked.

        What about RESTful routing?
            
            Go’s servemux doesn’t support routing based on the request method, URLs with variables in them and regexp-based patterns.

    2.4 Customizing HTTP headers

            Update our application so that the /snippet/create route only responds to HTTP requests which use the POST method.
        
        HTTP status codes

            It’s only possible to call w.WriteHeader() once per response, and after the status code has been written it can’t be changed.
            If you don’t call w.WriteHeader() explicitly, then the first call to w.Write() will automatically send a 200 OK status code to the user.
        
        Customizing headers

            Changing the response header map after a call to w.WriteHeader() or w.Write() will have no effect on the headers that the user receives.

        The http.Error shortcut

            If you want to send a non-200 status code and a plain-text response body then it’s a good opportunity to use the http.Error() shortcut.
            It’s quite rare to use the w.Write() and w.WriteHeader() methods directly.
        
        The net/http constants

            Use constants from the net/http package for HTTP methods and status codes, instead of writing the strings and integers.
            We can use the constant http.MethodPost instead of the string "POST".

            https://pkg.go.dev/net/http#pkg-constants
    
    2.5 URL query strings

        Update the snippetView handler so that it accepts an id query string parameter from the user.
            
            /snippet/view?id=1
        
        The r.URL.Query().Get() method will always return a string value for a parameter,
        or the empty string "" if no matching parameter exists.

        For the purpose of our Snippetbox application, we want to check that it contains a positive integer value.

    2.6 Project structure and organization

        The cmd directory will contain the application-specific code for the executable applications in the project.
        
        The internal directory will contain the ancillary non-application-specific code used in the project.
        We’ll use it to hold potentially reusable code like validation helpers and the SQL database models for the project.

        The ui directory will contain the user-interface assets used by the web application.

        It’s important to point out that the directory name internal carries a special meaning and behavior in Go:
        any packages which live under this directory can only be imported by code inside the parent of the internal directory.
        In our case, this means that any packages which live in internal can only be imported by code inside our snippetbox project directory.

    2.7 HTML templating and inheritance

            The .tmpl extension doesn’t convey any special meaning or behavior here.
            I’ve only chosen this extension because it’s a nice way of making it clear that the file
            contains a Go template when you’re browsing a list of files.

            Use Go’s html/template package, which provides a family of functions for safely parsing and rendering HTML templates.

        Template composition

            To save us typing and prevent duplication, it’s a good idea to create a base (or master) template which contains
            this shared content, which we can then compose with the page-specific markup for the individual pages.
            
            Go ahead and create a new ui/html/base.tmpl

            We use the ExecuteTemplate() method to tell Go that we specifically want to respond using the content of the
            base template (which in turn invokes our title and main templates).
        
        Embedding partials

            For some applications you might want to break out certain bits of HTML into partials that can be reused in different pages or layouts.

            Create ui/html/partials/nav.tmpl containing a named template called "nav".
            Update the base template so that it invokes the navigation partial using the {{template "nav" .}} action.
            Update the home handler to include the new ui/html/partials/nav.tmpl file when parsing the template files.
            The base template should now invoke the nav template.
    
    2.8 Serving static files

            Improve the look and feel of the home page by adding some static CSS and image files.
        
        The http.Fileserver handler
            Go’s net/http package ships with a built-in http.FileServer handler which you can use
            to serve files over HTTP from a specific directory.

            Let’s add a new route to our application so that all requests which begin with "/static/" are handled using this.
            The pattern "/static/" is a subtree path pattern, so it acts a bit like there is a wildcard at the end.

            https://stackoverflow.com/a/27946132/15308818

            This handler remove the leading slash from the URL path and then search the ./ui/static directory.

Chapter 3 Configuration and error handling

    3.1 Managing configuration settings

            Our web application’s main.go file currently contains a couple of hard-coded configuration settings:
                - The network address for the server to listen on (currently ":4000")
                - The file path for the static files directory (currently "./ui/static")
            There’s no separation between our configuration settings and code, and we can’t change the settings at runtime.
        
        Command-line flags
            A common and idiomatic way to manage configuration settings is to use command-line flags when starting an application.

            This defines a new command-line flag with the name addr, a default value of ":4000" and
            some short help text explaining what the flag controls.
                addr := flag.String("addr", ":4000", "HTTP network address")

            go run ./cmd/web -addr=":8080"
        
        Type conversions
            Go also has a range of other functions including flag.Int(), flag.Bool() and flag.Float64().
            These work in exactly the same way as flag.String().

        Automated help
            You can use the -help flag to list all the available command-line flags for an application.
            
            $ go run ./cmd/web -help
                Usage of /tmp/go-build3672328037/b001/exe/web:
                    -addr string
                        HTTP network address (default ":4000")
        
        Pre-existing variables
            It’s possible to parse command-line flag values into the memory addresses of pre-existing variables,
            using the flag.StringVar(), flag.IntVar(), flag.BoolVar() and other functions.

                type config struct {
                    addr      string
                    staticDir string
                }

                var cfg config

                flag.StringVar(&cfg.addr, "addr", ":4000", "HTTP network address")
                flag.StringVar(&cfg.staticDir, "static-dir", "./ui/static", "Path to static assets")

                flag.Parse()
    
    3.2 Leveled logging

            In our application, we can break apart our log messages into two distinct types — or levels.
            The first type is informational messages and the second type is error messages.
            The simple and clear approach is use the log.New() function to create two new custom loggers.

        Decoupled logging

            A big benefit of logging your messages to the standard streams (stdout and stderr) like we are is that your
            application and logging are decoupled.

            During development, it’s easy to view the log output because the standard streams are displayed in the terminal

            In staging or production environments, you can redirect the streams to a final destination for viewing and archival.
            This destination could be on-disk files, or a logging service such as Splunk.

            We could redirect the stdout and stderr streams to on-disk files when starting the application like so:

                $ go run ./cmd/web >>/tmp/info.log 2>>/tmp/error.log
