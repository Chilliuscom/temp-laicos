
## API Reference

### Login form

```http
  POST /api/login
```

#### Requires:
email: email

password:password

### Fetching data on initial webpage load

```http
  GET /api/initialLoad
```
Data received: 
```
all user information
all chatrooms information that user is member of(chatroom id and members)
last message info for each chatroom(sender, content, time)
```

### Fetching home page

```
POST /home 
```
Fetch with id of the bottom post of the page as post_ID and number of posts to be loaded as post_amount

---
returns the next x amount of posts that you are allowed to see

### Register form

```http
  POST /api/register
```

#### Requires:
email: email

password: password

first name: firstname

last name: lastname

date of birth: dateOfBirth

gender: gender

country: country

profile privacy (public or private): privacy

username/nickname (optional): username

about me(optional): aboutMe

profile pic(optional): avatar

### Content submission form

```http
  POST /api/submitContent
```


#### Post/comment creation
Requires multi-part-form

#### data.Type
--- 
post: post
```
title: header
content: content
current user's id: creator_id
image: avatar
privacy(PRIVATE: only followers, PUBLIC: everyone, ALMOSTPRIVATE: select followers that should see the post): privacy
followers that were selected to see the post: privacyUsers
```
--- 
group post: groupPost
```
title: header
content: content
current user's id: creator_id
group id: group_id
image: avatar
```
--- 
comment: comment
```
mater post: masterPost
content: content
current user's id: creator_id
image: avatar
```
--- 
group comment: groupComment
```
mater post: masterPost
content: content
current user's id: creator_id
group id: group_id
image: avatar
```

#### Requests and their responses:
Responses respond with the same request info but with accept/decline data

--- 
follow request: followRequest
```
from/current user: creator_id
to/user that receives request: target_user_id
```
--- 
follow request response: followRequestUpdate
```
to(user that sent the request): creator_id
from(user that got the request): target_user_id
update(ACCEPT/DECLINE): content
```
--- 
group join request: groupJoinRequest
```
from/current user: creator_id
group's id that receives the request: group_id
```
--- 
group join request response: joinUpdate
```
to(the user that requested to join): creator_id
group's id that received the request: group_id
update(ACCEPT/DECLINE): content
```
--- 
group invite : groupInvite
```
from(current user): creator_id
to(the user that was invited to join a group): target_user_id
group's id that sent the request: group_id
```
--- 
group invite response: inviteUpdate
```
to(the user that created the invite): creator_id
from(current user/the user that was invited to join a group): target_user_id
group's id that sent the request: group_id
update(ACCEPT/DECLINE): content
```

### Content submission form

```http
  POST /api/submitContent
```

### WEBSOCKET events

#### WS endpoint
```
  /ws/&email
```
#### "join"

Once someone comes online
```
websocket user id: UserID
list of websocket users' websocket user id's and their emails: Users
```

#### "disconnect"

Once someone goes offline
```
websocket user id: UserID
list of websocket users' websocket user id's and their emails: Users
```

#### "message"

Once someone comes online
```
message content: message
username that sent the message: from
chatroom id that contains the message: room
```

#### "notification"

Notification payload
```
notification event name: name
who(user) is responsible for this notification: from
where the nofication came from(your group, post, comment etc): where
link(href) redirecting to the post/group/comment in question: link
```
---
Notification event names
```
FollowRequest
FollowAccept
FollowDecline
GroupInvite
GroupInviteAccept
GroupInviteDecline
GroupJoinRequest
GroupJoinRequestAccept
GroupJoinRequestDecline
```

### WS instant message sending format

```http
  wsconnection.send
```
```
EventName: 'message',
        EventPayload: {
            roomID: data.selectedUserId, (chatroom id)
            message: message.value (message content)
        } 
```

### Previous message loading in chat

```http
  POST /api/messageRequest
```
```
current chatroom id: chatroomID
oldest message id: messageID
amount of messages to be loaded: amount
```

