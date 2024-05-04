Chapter 2
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
