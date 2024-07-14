<!-- > ls -l start.sh
-rw-r--r--  1 satyambaran  staff  110 Jun 27 23:54 start.sh
> chmod +x start.sh
> ls -l start.sh
-rwxr-xr-x  1 satyambaran  staff  110 Jun 27 23:54 start.sh -->


<!-- 
In the Authorization Code Grant flow, the client secret is used in the back-channel communication between the client application and the authorization server, not in the front-channel communication with the user-agent. This means the client secret is never exposed to the user, making it less susceptible to interception and misuse.


Instead of getting db and rdb from config, create a connection while using it only -->