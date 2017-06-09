## turk

A super simple GraphQL client. 
Retries if your GraphQL server isn't responding with exponential jitter backoff startegy.     
 
 
 Sample Usage:
 
 ```go
    email := "g@mail.com"
	password := "startww3now"
	gql := turk.GraphQL{
		Statement: `query CustomerLoginQuery {
					customerLogin(Email: {Email}, Password: {Password}){
						Avatar
						FirstName
						LastName
						Phone
						id
						Response
					}
				}`,
		Parameters: turk.Props{"Email": email, "Password": password},
	}.Build()

	resp, err := turk.NewTurkClient("http://localhost:8080/gql", &gql).Send()
```