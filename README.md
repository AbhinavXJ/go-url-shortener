# go-url-shortener

![Screenshot from 2025-03-31 22-56-15](https://github.com/user-attachments/assets/5952437b-4cb4-464b-b8fa-0b41d0a76d36)



(Have docker installed in your system)
Just run these commands after cloning,

"docker compose up --build"

To shorten url send POST request to http://localhost:8000/api/v1 in this format:-
eg:-
Body(JSON) of POST request

{
  "url":"https://www.cricbuzz.com/",
  "short":"",
  "expiry":45
}

You will get your shortened url in response-

then enter this in your browser http://localhost:8000/api/v1/(shortened-url-u-got) to redirect to your original url 

and you are good to go
