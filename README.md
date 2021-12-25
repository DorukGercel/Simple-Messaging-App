# Simple-Messaging-App
## Usages:
- Send messages to other online users
- Query your messages

## How to start:
### Server:
go run Server/server.go {PORT_NO}
### Client:
go run Client/client.go {NICKNAME} {IP_ADDR} {PORT_NO}

## Client Usages:
The application runs on terminal. The CLI commands are as follows.

After starting the application:
#### To send message to other user:
{USER_TO_SEND_NICK}[SIGLE-SPACE]{MESSAGE}

Ex: dodo hi bro how are things

These messages are stored under chat_records.txt file under application directory by the server

#### To query your related messages:
All queries must start with letter Q (capital)
##### To Me:
Q[SINGLE-SPACE]T
##### From Me:
Q[SINGLE-SPACE]F

Also in order to limit the number of requests returned:
Q[SINGLE-SPACE][T or F][SINGLE-SPACE]{NUMBER-OF-LIMIT}

## NOTES: 
- These syntax rules are highly enforced therefore error message would be displayed in case of mis-usage.
- Port number values checked accoring to the reserved port numbers.
- IP Addr Checker method not implemented as the net mask is not pre-determined.
- User can't send messahe to him/her-self.
- When server closes, clients are also closed.
