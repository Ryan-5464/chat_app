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
        
        chatContainer.appendChild(chat);
        
        const chatName = document.createElement.div('div')
        chatName.className = "chat-name"
        chatName.textContent = `${chat.ChatName}`
        chatElement.appendChild(chatName)
        
        const adminName = document.createElement.div('div')
        adminName.className = "chat-admin-name"
        adminName.textContent = `${chat.AdminName}`
        chatElement.appendChild(adminName)
        
        const memberCount = document.createElement.div('div')
        memberCount.className = "chat-member-count"
        memberCount.textContent = `${chat.MemberCount}`
        chatElement.appendChild(memberCount)
        
        const UnreadMsgCount = document.createElement.div('div')
        UnreadMsgCount.className = "chat-unread-message-count"
        UnreadMsgCount.textContent = `${chat.UnreadMsgCount}`
        chatElement.appendChild(UnreadMsgCount)
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
        messageElement.setAttribute("data-userid", message.UserId);
        messageElement.setAttribute("data-chatid", message.ChatId);
        messageElement.setAttribute("data-messageid", message.MessageId);
        messageElement.setAttribute("data-replyid", message.ReplyId);
        messageContainer.appendChild(messageElement);

        const author = document.createElement('div')
        author.className = "message-author"
        author.textContent = `${message.Author}`

        const createdAt = document.createElement('div')
        createdAt.className = "message-createdat"
        createdAt.textContent = `${message.CreatedAt}`

        const lastEditAt = document.createElement('div')
        lastEditAt.className = "message-lasteditat"
        lastEditAt.textContent = `${message.lastEditAt}`

        const messageText = document.createElement('div')
        messageText.className = "message-text"
        messageText.textContent = `${message.Text}`

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
