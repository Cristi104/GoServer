package main

import (
	"fmt"
	"net/http"
)

func main()  {
	http.HandleFunc("/csrf", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprintf(w, `<!DOCTYPE html>
		<html>
		<head>
		    <title>GoServer</title>
		</head>
		<body>
		<div>
		<script>
		fetch("http://192.168.0.137:8080/api/auth/signin", {
			method: "POST",
		    headers: {
			    Accept: "application/json",
				"Content-type": "application/json",
		    },
			credentials: true,
			})
			.then((response) => response.json())
			.then((data) => console.log(data))
		</script>
		</div>
		</body>
		</html>`)
	})
	http.ListenAndServe(":3000", nil)
}
