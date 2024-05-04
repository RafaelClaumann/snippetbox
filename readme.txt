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
