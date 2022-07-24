# scores_reflectorTN2
Simple Golang web site that caches or reflects score data from ScoresTN2 application

## BACKGROUND

As I was thinking about how to improve the ScoresTN2 application, I wanted to find a way to solve the problem of people shouting "Hey, what is the score...I still can't read it".

Several larger screen with BT and ardino, or even LED glasses and arduino rolled thru my head, or a service that texted scores etc.  All of these seem overkiil.  I finally settled on a hosted server.

But I wanted to keep it really simple.  I don't want the hastle or possible cost of a hacked or run away service.  I can tolerate someyhing small per month.  I have been using simple VPN that one of the Google/Alphabet companies created, "Outline VPN".  This shows you how to easily set up your own VPN and hosts on Digital Ocean for $5/month.  That is my goal for hosting this reflector.

I debated about a real web server with authentication or a websever with websockets.  First thought about using Dart (as shown in the Boaring show videos) or Golang with Websockets (a cominations used at my day job).  Several links are show in the REFERENCES section below.  One key take away is that I could create a server in any language and wrap it with a Docker container that should be easy to host.

I finally decided to try a very simple server without a real API.


## REALLY SIMPLE DESIGN

- Simple site that tracks the most recent scores
- Use get with params as *input*, not real api

## IDEAS ON HOSTING AND CONTROLLING COST

## ADDRESSING SECURITY


## TODO
- [ ] start go server implementation


## REFERENCES

- https://tutorialedge.net/golang/go-websocket-tutorial/
- https://gowebexamples.com/websockets/
- https://yalantis.com/blog/how-to-build-websockets-in-go/
- https://dev.to/heroku/deploying-your-first-golang-webapp-11b3
- https://medium.com/google-cloud/building-a-go-web-app-from-scratch-to-deploying-on-google-cloud-part-1-building-a-simple-go-aee452a2e654
- Google Flutter Boaring Show parts 1 and 2:
- https://www.youtube.com/watch?v=AaQzV1LTmo0&t=1s
- https://www.youtube.com/watch?v=K85PUBjFhn8&t=1s
- https://getoutline.org/
