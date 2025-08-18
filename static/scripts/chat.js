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
    SetupListeners()
    
    addNewMsgListenerToMsgInput()
    addNewMsgListenerToSendMsgButton()
    addAddContactEventListenerToInput()
})

function SetupListeners() {
    addChatModalListenerToChatContainer()
    addMessageModalListenerToMessageContainer()
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

function renderContactList(contactListData, overwrite) {
    console.log(":: rendering contact list")
    console.log(":: contact list data, ", contactListData)
    const contactListContainer = document.getElementById('contacts-container')
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

