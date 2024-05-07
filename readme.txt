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
        
        The http.Server error log

            By default, if Go’s HTTP server encounters an error it will log it using the standard logger.
            For consistency it’d be better to use our new errorLog logger instead.

            We need to initialize a new http.Server struct containing the configuration settings for our server,
            instead of using the http.ListenAndServe() shortcut.
    
    3.3 Dependency injection

            If you open up your handlers.go file you’ll notice that the home handler function is still
            writing error messages using Go’s standard logger, not the errorLog logger that we want to be using.

            How can we make our new errorLog logger available to our home function from main()?

            The simplest way is just put the dependencies in global variables.
            But it is good practice to inject dependencies into your handlers.
            It makes your code more explicit, less error-prone and easier to unit test than if you use global variables.

            For applications where all your handlers are in the same package, like ours,
            a neat way toinject dependencies is to put them into a custom application struct,
            and then define your handler functions as methods against application.

            https://www.geeksforgeeks.org/how-to-add-a-method-to-struct-type-in-golang/
        
        Closures for dependency injection

            The pattern that we’re using to inject dependencies won’t work if your handlers are spread across multiple packages.
            In that case, an alternative approach is to create a config package exporting an Application struct and have your
            handler functions close over this to form a closure.
    
    3.4 Centralized error handling

            Let’s neaten up our application by moving some of the error handling code into helper methods.
            This will help separate our concerns and stop us repeating code as we progress through the build.

            Add a new helpers.go file under the cmd/web directory.

            We use the debug.Stack() function to get a stack trace for the current goroutine and append it to the log message.
            We use the http.StatusText() function to automatically generate a human-friendly text representation of a given HTTP status code.
            For example, http.StatusText(400) will return the string "Bad Request".

            What we want to report is the file name and line number one step back in the stack trace, which would give
            us a clearer idea of where the error actually originated from.

            We can do this by changing the serverError()  helper to use our logger’s Output() function and setting the frame depth to 2.
    
    3.5 Isolating the application routes

            Our main() function is beginning to get a bit crowded, so to keep it clear and focused I’d
            like to move the route declarations for the application into a standalone routes.go file.

Chapter 4 Setting up MySQL

    4.2 Installing a database driver

        To use MySQL from our Go web application we need to install a database driver.
        https://github.com/golang/go/wiki/SQLDrivers

        Go to your project directory and run the go get command:
            go get github.com/go-sql-driver/mysql@v1
        If you want to download a specific version of a package:
            go get github.com/go-sql-driver/mysql@v1.0.3

    4.3 Modules and reproducible builds

            The new lines in go.mod essentially tells the Go command which exact version of github.com/go-sql-driver/mysql 
            should be used when you run a command like go run, go test or go build from your project directory.

            The go.sum file contains the cryptographic checksums representing the content of the required packages.
            The go.sum file isn’t designed to be human-editable and generally you won’t need to open it.
        
        Upgrading packages
            To upgrade to latest available minor or patch release of a package, you can simply run go get with the -u flag:
                go get -u github.com/foo/bar
            
            If you want to upgrade to a specific version then you should run the same command but with the appropriate @version suffix:
                go get -u github.com/foo/bar@v2.0.0
            
        Removing unused packages
            You could either run go get and postfix the package path with @none:
                go get github.com/foo/bar@none
            
            You can run go mod tidy, which will automatically remove any unused packages from your go.mod and go.sum files.
                go mod tidy -v

    4.4 Creating a database connection pool

            We need Go’s sql.Open() function to connect to the database from our web application.

                // The sql.Open() function initializes a new sql.DB object, which is essentially a
                // pool of database connections.
                db, err := sql.Open("mysql", "web:pass@/snippetbox?parseTime=true")
                if err != nil { ... }

            The parseTime=true part of the DSN above is a driver-specific parameter which instructs
            our driverto convert SQL TIME and DATE fields to Go time.Time object.

            The sql.Open() function returns a sql.DB object.
            This isn’t a database connection — it’s a pool of many connections.
            This is an important difference to understand.
            Go manages the connections in this pool as needed, automatically opening 
            and closing connections to the database via the driver.

            The connection pool is intended to be long-lived.
            In a web application it’s normal to initialize the connection pool in 
            your main() function and then pass the pool to your handlers.
        
        Usage in our web application

            Notice how the import path for our driver is prefixed with an underscore?
            This is because our main.go file doesn’t actually use anything in the mysql package.
            So if we try to import it normally the Go compiler will raise an error.
            We need the driver’s init() function to run so that it can register itself with the database/sql package.

            The sql.Open() function doesn’t actually create any connections, all it does is initialize the pool for future use.
            Actual connections to the database are established lazily, as and when needed for the first time.
            To verify that everything is set up correctly we need to use the db.Ping() method to create a connection and check for any errors.

    4.5 Designing a database model

            W’re going to sketch out a database model, you might want to think of it as a service layer or data access layer instead.
            We will encapsulate the code for working with MySQL in a separate package to the rest of our application.

            Create a new internal/models directory containing a snippets.go

            Snippet struct will represent the data for an individual snippet.
            SnippetModel type with methods on it to access and manipulate the snippets in our database.

        Using the SnippetModel

            To use this model in our handlers we need to establish a new SnippetModel struct in our
            main() function and then inject it as a dependency via the application struct.
        
        Benefits of this structure

            Separation of concerns.
            Database logic isn’t tied to our handlers which means that handler responsibilities are limited to HTTP stuff.
            Easier to write tight, focused, unit tests in the future.

            We have total control over which database is used at runtime, just by using the -dsn command-line flag.
    
    4.6 Executing SQL statements

        Let’s update the SnippetModel.Insert() to create a new record in our
        snippets table and then returns the integer id for the new record.

            INSERT INTO snippets (title, content, created, expires)
            VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
        
        Because the data we’ll be using will ultimately be untrusted user input from a form,
        it’s good practice to use placeholder parameters instead of interpolating data in the SQL query.

    Executing the query

        - DB.Query() is used for SELECT queries which return multiple rows.
            https://pkg.go.dev/database/sql#DB.Query

        - DB.QueryRow() is used for SELECT queries which return a single row.
            https://pkg.go.dev/database/sql#DB.QueryRow

        - DB.Exec() is used for statements which don’t return rows (like INSERT and DELETE).
            https://pkg.go.dev/database/sql#DB.Exec
        

        https://pkg.go.dev/database/sql#Result

        LastInsertId()
            Which returns the integer (an int64) generated by the database in response to a command.
            Typically this will be from an “auto increment” column when inserting a new row, which is exactly what’s happening in our case.
        RowsAffected()
            Which returns the number of rows (as an int64) affected by the statement.

        It is perfectly acceptable (and common) to ignore the sql.Result return value if you don’t need it. Like so:
            _, err := m.DB.Exec(stmt, title, content, expires)

    Using the model in our handlers

        Lets demonstrate how to call this new code from our handlers.
        
        docker exec -it snippet-db mysql -u root -p 
            mysql> SHOW databases;
                +--------------------+
                | Database           |
                +--------------------+
                | snippetbox         |
                +--------------------+
                5 rows in set (0.00 sec)

            mysql> USE snippetbox;
                Database changed
            
            mysql> SELECT id, title, expires FROM snippets;
                +----+------------------------+---------------------+
                | id | title                  | expires             |
                +----+------------------------+---------------------+
                |  1 | An old silent pond     | 2025-05-06 01:42:48 |
                |  2 | Over the wintry forest | 2025-05-06 01:42:48 |
                |  3 | First autumn morning   | 2024-05-13 01:42:48 |
                |  4 | O snail                | 2024-05-14 02:03:51 |
                +----+------------------------+---------------------+
                4 rows in set (0.00 sec)
