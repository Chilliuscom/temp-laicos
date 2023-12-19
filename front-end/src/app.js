//RESTFUL api STUFF

let app = document.getElementById("app")

let notificationButton = document.getElementById("notificationTest")
notificationButton.addEventListener("click", notificationTest)
async function notificationTest() {
    const response = await fetch("/testing")
    //location.href = ("http://localhost:8080/notificationTest")
}

let loginBtn = document.getElementById("HeaderLogin")
loginBtn.addEventListener("click", loginPage)
function loginPage() {
    app.innerHTML = `<form id="loginform">
    <label for="uname"><b>E-mail</b></label>
    <input type="email" placeholder="Enter E-mail address" name="uname" required>
    <br>

    <label for="psw"><b>Password</b></label>
    <input type="password" placeholder="Enter Password" name="psw" required>
    <br>

    <button type="submit">Login</button>
    </form>`

    let loginForm = document.getElementById("loginform")
    loginForm.addEventListener("submit", sendDataLogin)
}
async function sendDataLogin(e) {
    e.preventDefault()
    let loginForm = document.getElementById("loginform")
    let data = {
        email: loginForm.querySelector('input[name="uname"]').value,
        password: loginForm.querySelector('input[name="psw"]').value,
    }
    await fetch("api/login", {
        method: 'POST',
        headers: {
            'Content-type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(response => {
        console.log(response)
        let text = response.json().then(result => {
            console.log(result)
            console.log(result.status)
            //alert(result.message)
            if (result.message == "SUCCESSFUL LOGIN") {
                console.log("LOGINSUCCESS")
                localStorage.loggedIn = true
                localStorage.username = loginForm.querySelector('input[name="uname"]').value
                location.href = "http://localhost:8080"
            }
            if (!result.status) {
                alert(result.message)
            }
        })
    })
}

let logoutBtn = document.getElementById("HeaderLogOut")
logoutBtn.addEventListener("click", logout)
async function logout() {
    const response = await fetch("/api/logout");
    var data = await response.json()
    alert(data.message)
    localStorage.loggedIn = false
    localStorage.username = ""
}

let registerBtn = document.getElementById("HeaderRegister")
registerBtn.addEventListener("click", registerPage)
function registerPage() {
    app.innerHTML = `<form id="registerform" enctype="multipart/form-data">
    <label for="email"><b>E-mail</b></label>
    <input type="email" placeholder="Enter E-mail address" name="email" required>
    <br>

    <label for="psw"><b>Password</b></label>
    <input type="password" placeholder="Enter Password" name="password" required>
    <br>

    <label for="firstName"><b>First name</b></label>
    <input type="text" placeholder="Enter your first name" name="firstName" required>
    <br>
    
    

    <label for="lastName"><b>Last name</b></label>
    <input type="text" placeholder="Enter your last name" name="lastName" required>
    <br>
    
    <label for="username"><b>Username</b></label>
    <input type="text" placeholder="Enter your username" name="username" required>
    <br>
    
    <label for="lastName"><b>Phone</b></label>
    <input type="number" placeholder="Enter your phone" name="phone" required>
    <br>

    <label for="dateOfBirth"><b>Date of Birth</b></label>
    <input type="date" placeholder="Enter your date of birth" name="dateOfBirth" required>
    <br>
    
    <label for="cars">Choose a gender:</label>

    <select name="gender" id="gender">
      <option value="Male">Male</option>
      <option value="Female">Female</option>
    </select>
    <br>
    
    <label for="lastName"><b>Country</b></label>
    <input type="text" placeholder="Enter your country" name="country" required>
    <br>

    <label for="aboutMe"><b>About me (optional)</b></label>
    <input type="text" placeholder="Tell us something about yourself" name="aboutMe">
    <br>

    <label for="avatar">Select image:</label>
    <input type="file" id="avatar" name="avatar" accept="image/jpeg,image/png">
    <br>
    
    <fieldset>
      <legend>Select privacy:</legend>
    
      <div>
        <input type="radio" id="huey" name="privacy" value="1" checked />
        <label for="huey">Public</label>
      </div>
    
      <div>
        <input type="radio" id="dewey" name="privacy" value="0" />
        <label for="dewey">Private</label>
      </div>
    </fieldset>
    <br>

    <button type="submit">Register</button>
    </form>`
    let registerform = document.getElementById("registerform")
    registerform.addEventListener("submit", sendDataRegister)
}
async function sendDataRegister(e) {
    e.preventDefault()
    let loginForm = document.getElementById("registerform")
    let formData = new FormData(loginForm)
    let data = {
        email: loginForm.querySelector('input[name="email"]').value,
        password: loginForm.querySelector('input[name="password"]').value,
        firstName: loginForm.querySelector('input[name="firstName"]').value,
        lastName: loginForm.querySelector('input[name="lastName"]').value,
        username: loginForm.querySelector('input[name="username"]').value,
        dateOfBirth: loginForm.querySelector('input[name="dateOfBirth"]').value,
        phone: loginForm.querySelector('input[name="phone"]').value,
        country: loginForm.querySelector('input[name="country"]').value,
        gender: loginForm.querySelector('select[name="gender"]').value,
        aboutMe: loginForm.querySelector('input[name="aboutMe"]').value,
        avatar: loginForm.querySelector('input[name="avatar"]').value,
        privacy: loginForm.querySelector('input[name="privacy"]').value,
        //some kind of default/empty value for avatar/nickname/aboutMe
    }
    await fetch("api/register", {
        method: 'POST',
        body: formData
    }).then(response => {
        console.log(response)
        let text = response.json().then(result => {
            console.log(result)
            console.log(result.status)
            alert(result.message)
            if (!result.status) {
                alert(result.message)
            }
        })
    })
}

let homeBtn = document.getElementById("homeButton")
homeBtn.addEventListener("click", homePage)
function homePage() {
    app.innerHTML = `
    
    
    
    `
}

var username = ""
let headerUsername = document.getElementById("username")
document.getElementById("changeUsername").addEventListener("click", function () {
    username = ""
    location.reload()
})

var data = { //messaging data
    chatUserList: [],
    message: null,
    selectedUserId: null,
    userID: null,
    messageID: 0
}

let chatBtn = document.getElementById("chatButton")
chatBtn.addEventListener("click", Chat)

function setUsername() {
    username = prompt("please enter your username: ")
    headerUsername.innerHTML = username
}


window.addEventListener("load", (event) => {
    makeChatWebUi()

    if (localStorage.loggedIn != "true") {
        return
    }
    headerUsername.innerHTML = username
    username = localStorage.username
    initWS()
    initPage()
})

async function initPage() {
    const response = await fetch("/api/initialLoad");
    var data = await response.json()
    console.log(data)
}

var webSocketConnection = null

function initWS() {
    if (window["WebSocket"]) {
        const conn = new WebSocket("ws://" + document.location.host + "/ws/" + username)
        webSocketConnection = conn

        window.addEventListener('beforeunload', () => {
            if (webSocketConnection) {
                webSocketConnection.close(1000, "User is leaving the site");
            }
        })

        webSocketConnection.onclose = function (event) {
            console.log("WebSocket closed connection", "Code:", event.code, "Reason:", event.reason);
            if (event.code == 1006) {
                alert("Websocket closed connection")
            }

        };

        conn.onmessage = function (evt) {
            try {
                const socketPayload = JSON.parse(evt.data)
                console.log(socketPayload)
                switch (socketPayload.eventName) {
                    case 'notification':
                        notificationService(socketPayload)
                    case 'join':
                        setTimeout(() => {
                            getChats()
                        }, 50);

                    case 'disconnect':
                        if (!socketPayload.eventPayload) {
                            return
                        }
                        for (let i = 0; i < data.chatUserList.length; i++) {
                            if (data.chatUserList[i].userID == socketPayload.eventPayload.userID) {
                                document.getElementById(data.chatUserList[i].username).remove()
                            }
                        }
                        let selectedUserDisconnected = true
                        for (let i = 0; i < socketPayload.eventPayload.users.length; i++) {
                            if (socketPayload.eventPayload.users[i].username == data.userID) {
                                selectedUserDisconnected = false
                                setTimeout(() => {
                                    document.getElementById(socketPayload.eventPayload.users[i].username).className = "selecteduser"
                                }, 55);
                            }
                        }
                        if (selectedUserDisconnected) {
                            data.userID = null
                            document.getElementById("log").innerHTML = `<div id="anchor"></div>`
                            document.getElementById("chatheadertext").innerHTML = "Please select recipient from the left"
                        }
                        const userInitPayload = socketPayload.eventPayload;

                        data.chatUserList = userInitPayload.users
                        data.userID = data.userID === null ? userInitPayload.userID : data.userID
                        break

                    case 'message response':
                        if (!socketPayload.eventPayload) {
                            return
                        }
                        moveToTop(socketPayload.eventPayload.from)
                        if (socketPayload.eventPayload.from != data.userID) {
                            let newMSG = document.getElementById(socketPayload.eventPayload.from)
                            newMSG.className = "newMessage"

                            return
                        }

                        const messageContent = socketPayload.eventPayload;
                        const sentBy = messageContent.from ? messageContent.from : 'An unnamed fellow'
                        const actualMessage = messageContent.message;

                        data.message = `${actualMessage}`
                        printMessage("x", false, "NOW")
                        break
                    case 'previous message':
                        if (!socketPayload.eventPayload) {
                            return
                        }

                        const messageContent2 = socketPayload.eventPayload;
                        let you = "x"
                        if (messageContent2.from == username) {
                            you = "you"
                        }
                        let actualMessage2 = messageContent2.message
                        data.message = `${actualMessage2}`
                        let messageID = messageContent2.messageID
                        data.messageID = messageID
                        let messageTimestamp = messageContent2.time
                        printMessage(you, true, messageTimestamp)
                    default:
                        break;
                }
            } catch (error) {
                console.log(error)
                console.warn('Something went wrong while decoding the Message Payload')

            }
            return false;
        }
    }
}

function notificationService(socketPayload) {
    console.log(socketPayload)
}

let chatMessages

function Chat() {
    let webChat = document.getElementById("webchat")
    webChat.style.display = ""
    document.getElementById("SENDPM").addEventListener("click", () => {
        if (data.selectedUserId == null) {

            alert("please select user to send message to")
        } else {
            sendMessage()
        }
    })
    document.getElementById("MESSAGEREQUEST").addEventListener("click", () => {
        if (data.selectedUserId == null) {

            alert("please select user to send message to")
        } else {
            loadMessages()
        }
    })
    chatMessages = document.getElementById("log")
    scrollMoreMsg()
}

function makeChatWebUi() {
    const app = document.getElementById("app")
    let chatDiv = document.createElement("div")
    chatDiv.innerHTML = `
        <div class="webchat" id="webchat" style="display:none;">
            <div class="messages" id="messages">
                <div class="chatUser">
                </div>
            </div>
            <div class="chat">
                <div id="chatheader" class="chatheader">
                    <p id="chatheadertext">Please select recipient from the left</p>
                </div>
                <div class="chathistory" id="log">
                    <div id="anchor"></div>
                </div>
                <div id="chatForm">
                    <input type="text" id="chatInput">
                    <input type="button" id="SENDPM" value="SEND">
                    <input type="button" id="MESSAGEREQUEST" value="REQUEST MORE MESSAGES FROM DB">
                </div>
            </div>
        </div>
    `
    app.appendChild(chatDiv)
}

function getChats() {

    let messages = document.getElementById("messages")
    messages.innerHTML = ``
    for (const key in data.chatUserList) {
        let chat = addChat(data.chatUserList[key].username, data.chatUserList[key].userID)
        messages.append(chat)
    }

}

var lastChat = ""

function addChat(uName, uID) { //for selecting recipient
    let chat = document.createElement("div")
    chat.className = "chatUser"
    chat.style.cursor = 'pointer'
    chat.id = uName
    chat.innerHTML = uName
    chat.addEventListener("click", () => {
        if (lastChat != "") {
            lastChat.className = "chatUser"
        }
        chat.className = "selecteduser"
        document.getElementById("log").innerHTML = `<div id="anchor"></div>`
        data.messageID = 0
        data.selectedUserId = uID
        data.userID = uName
        document.getElementById("chatheadertext").innerHTML = uName
        loadMessages(true)
        document.scrollingElement.scroll(0, 1);
        lastChat = chat
    })
    return chat
}
function scrollMoreMsg() {
    document.getElementById("log").addEventListener('scroll', debounce(function () {
        loadMessages()
    }, 100))
}

function printMessage(you, previousMSG, timestamp) {
    let messageBox = document.querySelector("#log")
    let anchor = document.querySelector("#anchor")


    let message = document.createElement("p")

    message.innerHTML = data.message
    message.className = "message"
    if (timestamp == "NOW") {
        var now = new Date()
        timestamp = now.getFullYear() + '-' + (now.getMonth() + 1) + '-' + now.getDate() + ' ' + now.getHours() + ':' + now.getMinutes() + ':' + now.getSeconds()
    }
    let time = document.createElement("p")
    time.innerHTML = timestamp
    time.className = "timestamp"

    if (you == "you") {
        message.className = "yourmessage"
        time.className = "yourtimestamp"
    }
    if (previousMSG == true) {
        message.id = data.messageID

        messageBox.insertBefore(time, messageBox.firstChild)
        messageBox.insertBefore(message, messageBox.firstChild)
    } else {
        messageBox.insertBefore(message, anchor)
        messageBox.insertBefore(time, anchor)
    }
}

function sendMessage() {
    let message = document.getElementById("chatInput")
    webSocketConnection.send(JSON.stringify({
        EventName: 'message',
        EventPayload: {
            //roomID: data.selectedUserId,
            roomID: "3", //replace with room id
            message: message.value
        }
    }))
    data.message = message.value
    message.value = ''
    printMessage("you", false, "NOW")
    moveToTop(data.userID)
}

document.addEventListener("keydown", (f) => {
    if (f.key == 'Enter') {
        if (document.location.href == ("http://" + document.location.host + "/chat")) {
            if (data.selectedUserId == null) {
                alert("please select user to send message to")
            } else {
                sendMessage()
            }
        }
    }
})

function loadMessages(chatOpened) {

    if (chatMessages.scrollTop < chatMessages.scrollHeight * 0.1) {
        let messageID = "0"
        if (!chatOpened) {
            messageID = document.getElementById("log").firstChild.id
        }
        if (messageID == 1) { return }
        let amount = "10"
        requestMessages(data.selectedUserId, messageID, amount)
        /*
        webSocketConnection.send(JSON.stringify({
            EventName: 'messageRequest',
            EventPayload: {
                userID: data.selectedUserId,
                messageID: messageID,
                amount: amount
            }
        }))*/
    }
}

async function requestMessages(chatroomID, messageID, amount) {
    let data = {
        chatroomID: parseInt(chatroomID),
        messageID: parseInt(messageID),
        amount: parseInt(amount),

    }
    const response = await fetch("/api/messageRequest", {
        method: 'POST',
        headers: {
            'Content-type': 'application/json'
        },
        body: JSON.stringify(data)
    }).then(response => {
        console.log(response)
        let data = response.json().then(result => {
            console.log(result)
        })
    })
}

function moveToTop(username) {
    let messages = document.getElementById("messages")
    let user = document.getElementById(username)
    messages.insertBefore(user, messages.firstChild)
}

function throttle(func, wait) {
    let lastexec = 0
    return function (...args) {
        let currentTime = new Date()
        if (currentTime - lastexec >= wait) {
            func(...args)
            lastexec = currentTime
        }
    }
}

function debounce(func, wait, leading) {
    var timeout;
    return function () {
        var context = this, args = arguments;
        clearTimeout(timeout);
        timeout = setTimeout(function () {
            timeout = null;
            if (!leading) func.apply(context, args);
        }, wait);
        if (leading && !timeout) func.apply(context, args);
    };
}