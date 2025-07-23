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
};

document.addEventListener("DOMContentLoaded", function () {
  document.querySelectorAll('.chat').forEach(function (el) {
    el.addEventListener('click', function () {
      const chatId = el.getAttribute('data-chatid');
        requestChatMessages(chatId)
    });
  });
});

function requestChatMessages(chatId) {
    payload = {
        Type: "SwitchChat",
        Data: {
            ChatId: chatId,
        }
    }
    console.log("sending payload: ", payload)
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
        const chatElement = document.createElement('div');
        chatElement.className = "chat"
        chatElement.setAttribute("data-chatid", chat.ChatId);
        chatElement.setAttribute("data-adminid", chat.AdminId);
        chatElement.setAttribute("data-createdat", chat.CreatedAt);
        
        chatElement.addEventListener("click", function() {
            requestChatMessages(chat.ChatId)
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
            "data-messageid": message.MessageId,
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

function sendMessage() {
    const messageInput = document.getElementById("input");
    const message = messageInput.value.trim();

    if (message) {
        socket.send(message);
        messageInput.value = '';
        console.log("Message sent")
    }
}
