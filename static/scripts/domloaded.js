const typeNewMessage = "1"
const typeNewContactMessage = "4"

const socket = new WebSocket('ws://localhost:8081/ws');

socket.onopen = function () {
    console.log("WebSocket connection established.");
}

socket.onmessage = function (event) {
    console.log(":: New message received.");
    const payload = JSON.parse(event.data);
    console.log(":: payload, ", payload);

    let messageChatId;
    if (!payload.Messages) return;
    Object.values(payload.Messages).forEach(message => {
        messageChatId = message.ChatId;
    });
    let activeChatId
    const activeChats = document.querySelectorAll('.active')
    if (!activeChats || activeChats.length === 0) {
        throw new Error("no active chat found")
    }

    const activeChat = activeChats[0]
    console.log("activeChat", activeChat)

    if (activeChat.classList.contains('chat')) {
        activeChatId = activeChat.getAttribute('data-chatid')
    } 

    if (activeChat.classList.contains('contact')) {
        activeChatId = activeChat.getAttribute('data-contactchatid')
    }

    if (!activeChat) {
        throw new Error("failed to get active chat id")
    }

    console.log("activeChatId", activeChatId)
    console.log("messageChatId", messageChatId)


    if (activeChatId == messageChatId) {
        console.log("active chat message received")
        HandleNewMessageResponse(payload);
    } else {
        console.log("other chat message received")
    }

    AutoScrollToBottom()
}

document.addEventListener("DOMContentLoaded", function() {
    SetupListeners()
    
    addNewMsgListenerToMsgInput()
    addNewMsgListenerToSendMsgButton()

    document.querySelectorAll('.message').forEach(msg => {
        formatMessageDates(msg)
    })

    AutoScrollToBottom()

})

function AutoScrollToBottom() {
    const container = document.getElementById('messages-container')
    requestAnimationFrame(() => {
        container.scrollTop = container.scrollHeight;
    });

}


function SetupListeners() {
    addChatModalListenerToChatContainer()
    addMessageModalListenerToMessageContainer()
    addContactModalListenerToContactContainer()
    addSwitchChatControllerToChatsContainer()
    addSwitchContactChatControllerToContactsContainer()
    addNewChatEventListenerToNewChatInput()
    addAddContactEventListenerToAddContactInput()
    ConfigureMemberListModal() 
}





function addNewMsgListenerToMsgInput() {
    const input = document.getElementById("message-input")
    input.addEventListener('keydown', function (event) {
        if (event.key === "Enter") {
            const replyId = null
            const msgText = input.value.trim()
            const chat = document.querySelector(".active")
            chatId = chat.getAttribute("data-chatid")
            if (chatId == null) {
                chatId = chat.getAttribute("data-contactchatid")
                sendContactMessage(msgText, chatId, replyId)
            } else {
                sendMessage(msgText, chatId, replyId)
            }
            input.value = '';
        }
    })
}

function addNewMsgListenerToSendMsgButton() {
    const button = document.getElementById("send-message-button")
    button.addEventListener('click', function () {
        const input = document.getElementById("message-input")
        const replyId = null
        const msgText = input.value.trim()
        const chat = document.querySelector(".active")
        let chatId = chat.getAttribute("data-chatid")
        if (chatId == null) {
            chatId = chat.getAttribute("data-contactchatid")
            sendContactMessage(msgText, chatId, replyId)
        } else {
            sendMessage(msgText, chatId, replyId)
        }
        input.value = '';
    })
}

function sendContactMessage(msgText, chatId, replyId) {
    console.log(":: sending new message")
    if (msgText) {
        payload = {
            Type: typeNewContactMessage,
            Data: {
                ChatId: chatId,
                MsgText: msgText,
                ReplyId: replyId
            }
        }
        console.log(":: request payload, ", payload)
        socket.send(JSON.stringify(payload));
    }
}


function sendMessage(msgText, chatId, replyId) {
    console.log(":: sending new message")
    if (msgText) {
        payload = {
            Type: typeNewMessage,
            Data: {
                ChatId: chatId,
                MsgText: msgText,
                ReplyId: replyId
            }
        }
        console.log(":: request payload, ", payload)
        socket.send(JSON.stringify(payload));
    }
}

function formatMessageDates(messageElem) {
    const createdAtRaw = messageElem.dataset.createdat
    const lastEditAtRaw = messageElem.dataset.lasteditedat

    const createdAt = new Date(createdAtRaw)
    const lastEditAt = new Date(lastEditAtRaw)

    const header = document.createElement('div')
    header.classList.add('message-header')

    if (createdAt < lastEditAt) {
        const lastEditDiv = document.createElement('div')
        lastEditDiv.classList.add('message-lasteditat')
        lastEditDiv.innerText = `Edited: ${lastEditAt.toLocaleString()}`
        header.appendChild(lastEditDiv)
    } else {
        const createdDiv = document.createElement('div')
        createdDiv.classList.add('message-createdat')
        createdDiv.innerText = `Sent: ${createdAt.toLocaleString()}`
        header.appendChild(createdDiv)
    }

    messageElem.insertBefore(header, messageElem.querySelector('.message-text'))
}
