# simple add, read
curl localhost:80/reset
curl localhost:80/

# assume data=<keeper>,backcolor1,backcolor2,color1,color2,name1,name2,sets1,sets2,score1,score2,posession
# *color*: 6 char hex, rgb
# name*: characters, use %20 for space
# sets*, score*: integer
# posession: 0-none, 1-them, 2-us
# ie. shannon,000000,ffffff,ffffff,000000,Them,Us,0,0,10,8,0
# may have spaces declared as %20, looks like gin converts them to ' '

curl "localhost:3000/add?data=shannon,000000,ffffff,ffffff,000000,Them,Us,0,0,0,0,0"
#sleep 1
curl "localhost:3000/add?data=joe,000000,ffffff,ffffff,000000,Them,Us,0,0,1,0,1"
#sleep 2
curl "localhost:3000/add?data=joe,000000,ffffff,ffffff,000000,Them,Us,0,0,2,0,1"
#sleep 1
curl "localhost:3000/add?data=shannon,000000,ffffff,ffffff,000000,Them,Us,0,0,2,1,2"
#sleep 3
curl "localhost:3000/add?data=shannon,000000,ffffff,ffffff,000000,Them,Us,0,0,2,2,2"
//sleep 30
curl "localhost:3000/add?data=jim,000000,ffffff,ffffff,000000,Them,Us,0,0,1,0,1"
#sleep 2
curl "localhost:3000/add?data=jim,000000,ffffff,ffffff,000000,Them,Us,0,0,2,0,1"
curl "localhost:3000/add?data=shannon,GOROCKY"


echo "\nfrom /"
curl localhost:80/
echo "\nfrom /raw"
curl localhost:80/raw
echo "\nfrom /json"
curl localhost:80/json
echo ""
echo "\nfrom /shannon/raw"
curl localhost:80/shannon/raw
echo "\nfrom /shannon/json"
curl localhost:80/shannon/json
echo ""
