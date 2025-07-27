/* IMPORTANT 

    INFORMATION IS EXTRACTED FROM ELEMENTS FROM THE FRONTEND FOR DEV. THIS IS A VULNERABILITY. 
    NEEDS TO BE REWRITTEN USING SECURE TOKENS.

*/

const socket = new WebSocket('ws://localhost:8081/ws');

socket.onopen = function () {
    console.log("WebSocket connection established.");
}

socket.onmessage = function (event) {
    console.log("New message received.");
    const payload = JSON.parse(event.data);
    console.log(payload)
    renderChats(payload.Chats, false);
    renderMessages(payload.Messages, true);
}

document.addEventListener("DOMContentLoaded", function() {
    addSwitchChatListenerToChats()
    addNewMsgListenerToMsgInput()
    addNewMsgListenerToSendMsgButton()
    addNewChatEventListener()
})

function addNewChatEventListener() {
    elem = document.getElementById("new-chat-button")
    elem.addEventListener('click', function () {
        const chatNameInput = document.getElementById("chat-name-input")
        const chatName = chatNameInput.value.trim()
        const messages = document.getElementsByClassName("message")
        userId = messages[0].getAttribute("data-userid")
        sendNewChatInfo(userId, chatName)
        chatNameInput.value = '';
    })
}

function addSwitchChatListenerToChats() {
    document.querySelectorAll('.chat').forEach(function (elem) {
        elem.addEventListener('click', function () {
            const chatId = elem.getAttribute('data-chatid')
            requestChatMessages(chatId)
        })
    })
}

function addNewMsgListenerToMsgInput() {
    const input = document.getElementById("input")
    input.addEventListener('keydown', function (event) {
        if (event.key === "Enter") {
            const chatId = getChatIdFromExistingMessage()
            const userId = getUserIdFromExistingChat()
            const replyId = null
            const msgText = input.value.trim()
            sendMessage(userId, msgText, chatId, replyId)
            input.value = '';
        }
    })
}

function addNewMsgListenerToSendMsgButton() {
    const button = document.getElementById("send-message-button")
    button.addEventListener('click', function () {
        const input = document.getElementById("input")
        const userId = getUserIdFromExistingMessage()
        const chatId = getChatIdFromExistingMessage()
        const replyId = null
        const msgText = input.value.trim()
        sendMessage(userId, msgText, chatId, replyId)
    })
}

function getChatIdFromExistingMessage() {
    const message = document.getElementsByClassName("message")
    chatId = message[0].getAttribute("data-chatid")
    console.log(chatId)
    return chatId
}

function getUserIdFromExistingMessage() {
    const chat = document.getElementsByClassName("message")
    userId = chat[0].getAttribute("data-userid")
    console.log(userId)
    return userId
}

function sendNewChatInfo(userId, chatName) {
    console.log("userId: ", userId, "chatName: ", chatName)
    if (chatName) {
        payload = {
            Type: "NewChat",
            Data: {
                UserId: userId,
                ChatName: chatName,
            }
        }
        socket.send(JSON.stringify(payload));
    }
}

function sendMessage(userId, msgText, chatId, replyId) {
    console.log("userId: ", userId, "chatId: ", chatId, "replyId: ", replyId, "msgText: ", msgText)
    if (msgText) {
        payload = {
            Type: "NewMessage",
            Data: {
                UserId: userId,
                ChatId: chatId,
                MsgText: msgText,
                ReplyId: replyId
            }
        }
        socket.send(JSON.stringify(payload));
    }
}

function requestChatMessages(chatId) {
    console.log("chatId: ", chatId)
    payload = {
        Type: "SwitchChat",
        Data: {
            ChatId: chatId,
        }
    }
    socket.send(JSON.stringify(payload))
}

function renderChats(chatData, overwrite) {
    console.log(chatData)
    if (chatData == null) {
        return
    }

    const chatContainer = document.getElementById('chats-container');
    if (overwrite == true) {
        chatContainer.innerHTML = ''; 
    }

    chatData.forEach(chat => {
        console.log("CHAT: ", chat)
        const chatElement = document.createElement('div');
        chatElement.className = "chat"
        chatElement.setAttribute("data-chatid", chat.Id);
        chatElement.setAttribute("data-adminid", chat.AdminId);
        chatElement.setAttribute("data-createdat", chat.CreatedAt);
        
        chatElement.addEventListener("click", function() {
            requestChatMessages(chat.Id)
        })
        
        chatContainer.appendChild(chatElement);

        const data = {
            "chat-name": `${chat.ChatName}`,
            "chat-admin-name": `${chat.AdminName}`,
            "chat-member-count": `${chat.MemberCount}`,
            "chat-unread-message-count": `${chat.UnreadMsgCount}`,
        }
        for (const [key, value] of Object.entries(data)) {
            const element = document.createElement('div')
            element.className = key
            element.textContent = `${value}`
            chatElement.appendChild(element);
        }
    });
}

function renderMessages(messageData, overwrite) {
    const messageContainer = document.getElementById('messages-container');
    if (overwrite == true) {
        messageContainer.innerHTML = ''
    }

    messageData.forEach(message => {
        const messageElement = document.createElement('div');
        messageElement.className = "message";

        const attributes = {
            "data-userid": message.UserId,
            "data-chatid": message.ChatId,
            "data-messageid": message.Id,
            "data-replyid": message.ReplyId,
        };
        for (const [key, value] of Object.entries(attributes)) {
            messageElement.setAttribute(key, value);
        }
        messageContainer.appendChild(messageElement);

        const data = {
            "message-author": `${message.Author}`,
            "message-createdat": `${message.CreatedAt}`,
            "message-lasteditat": `${message.lastEditAt}`,
            "message-text": `${message.Text}`,
        }
        for (const [key, value] of Object.entries(data)) {
            const element = document.createElement('div')
            element.className = key
            element.textContent = `${value}`
            messageElement.appendChild(element);
        }

        messageContainer.scrollTop = messageContainer.scrollHeight;
    });
}


