/* IMPORTANT 

    INFORMATION IS EXTRACTED FROM ELEMENTS FROM THE FRONTEND FOR DEV. THIS IS A VULNERABILITY. 
    NEEDS TO BE REWRITTEN USING SECURE TOKENS.

*/

const typeNewMessage = "1"
const typeNewContactMessage = "4"

const socket = new WebSocket('ws://localhost:8081/ws');

socket.onopen = function () {
    console.log("WebSocket connection established.");
}

socket.onmessage = function (event) {
    console.log(":: New message received.");
    const payload = JSON.parse(event.data);
    console.log(":: payload, ", payload)
    if (payload.Type == typeNewMessage) {
        console.log(":: appending message")
        renderChats(payload.Chats, false);
        renderMessages(payload.Messages, false);
        return
    }
    console.log(":: overwriting messages")
    renderChats(payload.Chats, false);
    renderMessages(payload.Messages, true);
}

document.addEventListener("DOMContentLoaded", function() {
    addSwitchChatListenerToChats()
    addSwitchChatListenerToContacts()
    addNewMsgListenerToMsgInput()
    addNewMsgListenerToSendMsgButton()
    addNewChatEventListenerToButton()
    addNewChatEventListenerToInput()
    addAddContactEventListenerToButton()
    addAddContactEventListenerToInput()
    addChatToggleEventListenerToContainer()
    highlightActiveChat() 
})

function highlightActiveChat() {
    console.log("highlighting active chat")
    const chats = document.querySelectorAll('.chat')
    if (chats.length == 0) {
        return
    }
    const chatId = chats[0].getAttribute('data-chatid')
    switchActiveChat(chatId)
}

function addAddContactEventListenerToButton() {
    elem = document.getElementById("add-contact-button")
    elem.addEventListener('click', function () {
        const contactEmailInput = document.getElementById("contact-email-input")
        const email = contactEmailInput.value.trim()
        console.log("contactEmail: ", email)
        contactEmailInput.value = '';
        addContact(email)
    })
}

function addAddContactEventListenerToInput() {
    elem = document.getElementById("contact-email-input")
    elem.addEventListener('keydown', function (event) {
        if (event.key == "Enter") {
            const contactEmailInput = document.getElementById("contact-email-input")
            const email = contactEmailInput.value.trim()
            console.log("contactEmail: ", email)
            contactEmailInput.value = '';
            addContact(email)
        }
    })
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

function addSwitchChatListenerToContacts() {
    document.querySelectorAll('.contact').forEach(function (elem) {
        elem.addEventListener('click', function () {
            const contactChatId = elem.getAttribute('data-contactchatid')
            switchContactChat(contactChatId)
        })
    })
}

// need to refactor to handle ids better
function addNewMsgListenerToMsgInput() {
    const input = document.getElementById("input")
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

// need to refactor to handle ids better
function addNewMsgListenerToSendMsgButton() {
    const button = document.getElementById("send-message-button")
    button.addEventListener('click', function () {
        const input = document.getElementById("input")
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

function renderChats(chatData, overwrite) {
    console.log(":: rendering chats")
    console.log(":: chat data, ", chatData)
    if (chatData == null) {
        console.log(":: chat data is null")
        return
    }

    const chatContainer = document.getElementById('chats-container');
    if (overwrite == true) {
        console.log(":: overwriting chats")
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
            switchChat(chat.Id)
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

function renderContactList(contactListData, overwrite) {
    console.log(":: rendering contact list")
    console.log(":: contact list data, ", contactListData)
    const contactListContainer = document.getElementById('contact-list-container')
    if (overwrite == true) {
        console.log(":: overwriting contact list")
        contactListContainer.innerHTML = ''; 
    }

    contactListData.forEach(contact => {
        console.log("contact: ", contact)
        const contactElement = document.createElement('div');
        contactElement.className = "contact"
        console.log(contact.Id)
        contactElement.setAttribute("data-contactid", contact.Id)
        contactElement.setAttribute("data-contactchatid", contact.ContactChatId)

        contactElement.addEventListener("click", function() {
            switchChat(contact.ContactChatId)
        })
        
        contactListContainer.appendChild(contactElement);

        const data = {
            "contact-name": `${contact.Name}`,
            "contact-email": `${contact.Email}`,
            "contact-since": `${contact.contactSince}`,
            "contact-status": `${contact.OnlineStatus}`,
        }
        for (const [key, value] of Object.entries(data)) {
            const element = document.createElement('div')
            element.className = key
            element.textContent = `${value}`
            contactElement.appendChild(element);
        }
    });
}

function renderMessages(messageData, overwrite) {
    console.log(":: rendering messages")
    console.log(":: message data, ", messageData)
    const messageContainer = document.getElementById('messages-container');
    if (overwrite == true) {
        console.log(":: overwriting old messages")
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

    container.querySelectorAll('.chat').forEach(chat => {
      chat.classList.remove('active');
    });
     document.querySelectorAll(`.contact`).forEach(chat => {
        chat.classList.remove(`active`)
    })

    clickedChat.classList.add('active');
  });
}

/* NEW Contact REQUEST =================================================== */

function addContact(email) {
    console.log(":: adding new contact")
    fetch(BASEURL + '/api/chat/contact/add', addContactRequestBody(email))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[AddContact]Received:', responsePayload);
        console.log(":: payload, ", responsePayload)
        renderContactList(responsePayload.Contacts, false);
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function addContactRequestBody(email) {
    request = {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            Email: email,
        })
    }
    console.log(":: add contact request details, ", request)
    return request
}

/* NEW CHAT REQUEST ===================================================== */

function newChat(chatName) {
    console.log(":: Creating new chat")
    fetch(BASEURL + '/api/chat/new', newChatRequestBody(chatName))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[NewChat]Received:', responsePayload);
        console.log(":: payload, ", responsePayload)
        renderChats(responsePayload.Chats, false);
        renderMessages(responsePayload.Messages, true);
        console.log(":: new active chat id, ", responsePayload.NewActiveChatId)
        const newActiveChatId = responsePayload.NewActiveChatId 
        switchActiveChat(newActiveChatId)
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function newChatRequestBody(chatName) {
    request = {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            Name: chatName,
        })
    }
    console.log(":: new chat request details, ", request)
    return request
}

/* SWITCH CHAT REQUEST ================================================== */

function switchChat(chatId) {
    console.log(":: Switching chat")
    fetch(BASEURL + '/api/chat/switch', switchChatRequestBody(chatId))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[Switch Chat]Received:', responsePayload);
        console.log(":: payload, ", responsePayload)
        renderMessages(responsePayload.Messages, true);
        console.log(":: new active chat id, ", responsePayload.NewActiveChatId)
        const newActiveChatId = responsePayload.NewActiveChatId 
        switchActiveChat(newActiveChatId)
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function switchChatRequestBody(chatId) {
    request = {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            ChatId: String(chatId),
        })
    }
    console.log(":: switch contact chat request details, ", request)
    return request
}

function switchActiveChat(newActiveChatId) {
    console.log(":: highlighting new active chat")
     document.querySelectorAll(`.chat`).forEach(chat => {
        chat.classList.remove(`active`)
    })
     document.querySelectorAll(`.contact`).forEach(chat => {
        chat.classList.remove(`active`)
    })

    const newChat = document.querySelector(`[data-chatid="${newActiveChatId}"]`)
    newChat.classList.add(`active`)
}

/* SWITCH CONTACT CHAT REQUEST ================================================== */

function switchContactChat(chatId) {
    console.log(":: Switching contact chat")
    fetch(BASEURL + '/api/chat/contact/switch', switchContactChatRequestBody(chatId))
    .then(response => {
        if (!response.ok) {
            throw new Error(`Server error: ${response.status}`);
        }
        return response.json(); 
    })
    .then(responsePayload => {
        console.log('[Switch Contact Chat]Received:', responsePayload);
        console.log(":: payload, ", responsePayload)
        renderMessages(responsePayload.Messages, true);
        console.log(":: new active chat id, ", responsePayload.NewActiveChatId)
        const newActiveContactChatId = responsePayload.NewActiveChatId 
        switchActiveContactChat(newActiveContactChatId)
    })
    .catch(error => {
        console.error('Fetch error:', error);
    });
}

function switchContactChatRequestBody(chatId) {
    request = {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json' 
        },
        body: JSON.stringify({
            ChatId: String(chatId),
        })
    }
    console.log(":: switch contact chat request details, ", request)
    return request
}

function switchActiveContactChat(newActiveContactChatId) {
    console.log(":: highlighting new active contact chat")
     document.querySelectorAll(`.contact`).forEach(chat => {
        chat.classList.remove(`active`)
    })
     document.querySelectorAll(`.chat`).forEach(chat => {
        chat.classList.remove(`active`)
    })
    
    const newChat = document.querySelector(`[data-contactchatid="${newActiveContactChatId}"]`)
    newChat.classList.add(`active`)
}
