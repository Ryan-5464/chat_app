/* IMPORTANT 

    INFORMATION IS EXTRACTED FROM ELEMENTS FROM THE FRONTEND FOR DEV. THIS IS A VULNERABILITY. 
    NEEDS TO BE REWRITTEN USING SECURE TOKENS.

*/


const typeNewMessage = "1"



const socket = new WebSocket('ws://localhost:8081/ws');

socket.onopen = function () {
    console.log("WebSocket connection established.");
}

socket.onmessage = function (event) {
    console.log("New message received.");
    const payload = JSON.parse(event.data);
    console.log(payload)
    if (payload.Type == typeNewMessage) {
        renderChats(payload.Chats, false);
        renderMessages(payload.Messages, false);
        return
    }
    renderChats(payload.Chats, false);
    renderMessages(payload.Messages, true);
}

document.addEventListener("DOMContentLoaded", function() {
    addSwitchChatListenerToChats()
    addNewMsgListenerToMsgInput()
    addNewMsgListenerToSendMsgButton()
    addNewChatEventListenerToButton()
    addNewChatEventListenerToInput()
    addChatToggleEventListenerToContainer()
    highlightActiveChat() 
})

function highlightActiveChat() {
    const chats = document.querySelectorAll('.chat')
    const chatId = chats[0].getAttribute('data-chatid')
    switchActiveChat(chatId)
}

function addNewChatEventListenerToButton() {
    elem = document.getElementById("new-chat-button")
    elem.addEventListener('click', function () {
        const chatNameInput = document.getElementById("chat-name-input")
        const chatName = chatNameInput.value.trim()
        console.log("newChatName: ", chatName)
        chatNameInput.value = '';
        newChat(chatName)
    })
}

function addNewChatEventListenerToInput() {
    elem = document.getElementById("chat-name-input")
    elem.addEventListener('keydown', function (event) {
        if (event.key == "Enter") {
            const chatNameInput = document.getElementById("chat-name-input")
            const chatName = chatNameInput.value.trim()
            console.log("newChatName: ", chatName)
            chatNameInput.value = '';
            newChat(chatName)
        }
    })
}

function addSwitchChatListenerToChats() {
    document.querySelectorAll('.chat').forEach(function (elem) {
        elem.addEventListener('click', function () {
            const chatId = elem.getAttribute('data-chatid')
            switchChat(chatId)
        })
    })
}

// need to refactor to handle ids better
function addNewMsgListenerToMsgInput() {
    const input = document.getElementById("input")
    input.addEventListener('keydown', function (event) {
        if (event.key === "Enter") {
            const chatId = getActiveChatId()
            const replyId = null
            const msgText = input.value.trim()
            input.value = '';
            sendMessage(msgText, chatId, replyId)
        }
    })
}

// need to refactor to handle ids better
function addNewMsgListenerToSendMsgButton() {
    const button = document.getElementById("send-message-button")
    button.addEventListener('click', function () {
        const input = document.getElementById("input")
        const chatId = getActiveChatId()
        const replyId = null
        const msgText = input.value.trim()
        sendMessage(msgText, chatId, replyId)
        input.value = '';
    })
}

function getActiveChatId() {
    const chat = document.querySelector(".active")
    chatId = chat.getAttribute("data-chatid")
    console.log("active chat id:", chatId)
    return chatId
}


function sendMessage(msgText, chatId, replyId) {
    console.log("MESSAGE::chatId: ", chatId, "replyId: ", replyId, "msgText: ", msgText)
    if (msgText) {
        payload = {
            Type: typeNewMessage,
            Data: {
                ChatId: chatId,
                MsgText: msgText,
                ReplyId: replyId
            }
        }
        socket.send(JSON.stringify(payload));
    }
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
            "chat-name": `${chat.Name}`,
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

function addChatToggleEventListenerToContainer() {
  const container = document.querySelector('#chats-container');
  if (!container) return;

  container.addEventListener('click', function (e) {
    const clickedChat = e.target.closest('.chat');
    if (!clickedChat || !container.contains(clickedChat)) return;

    // Remove 'active' from all chats
    container.querySelectorAll('.chat').forEach(chat => {
      chat.classList.remove('active');
    });

    // Add 'active' to the clicked chat
    clickedChat.classList.add('active');
  });
}

/* NEW CHAT REQUEST ===================================================== */

function newChat(chatName) {
    fetch(BASEURL + '/api/chat/new', newChatRequestBody(chatName))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[NewChat]Received:', responsePayload);
        renderChats(responsePayload.Chats, false);
        renderMessages(responsePayload.Messages, true);
        const newActiveChatId = responsePayload.NewActiveChatId 
        switchActiveChat(newActiveChatId)
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function newChatRequestBody(chatName) {
    return {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            Name: chatName,
        })
    }
}

/* SWITCH CHAT REQUEST ================================================== */

function switchChat(chatId) {
    fetch(BASEURL + '/api/chat/switch', switchChatRequestBody(chatId))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[SwitchChat]Received:', responsePayload);
        renderMessages(responsePayload.Messages, true);
        const newActiveChatId = responsePayload.NewActiveChatId 
        switchActiveChat(newActiveChatId)
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function switchChatRequestBody(chatId) {
    return {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            ChatId: chatId,
        })
    }
}

function switchActiveChat(newActiveChatId) {
     document.querySelectorAll(`.chat`).forEach(chat => {
        chat.classList.remove(`active`)
    })
    
    const newChat = document.querySelector(`[data-chatid="${newActiveChatId}"]`)
    newChat.classList.add(`active`)
}
