const LEAVE_CHAT_ENDPOINT = '/api/chat/leave'
const DEL_MSG_ENDPOINT = '/api/message/delete'
const EDIT_CHAT_NAME_ENDPOINT = '/api/chat/edit'
const CHAT_SWITCH_ENDPOINT = '/api/chat/switch'
const NEW_CHAT_ENDPOINT = '/api/chat/new'

function GET() {
    return {
        method: 'GET',
    }
}

function POST(json) {
    return {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(json),
    }
}

function DELETE() {
    return {
        method: 'DELETE',
    }
}

async function DeleteMessageRequest(chatId, messageId, userId) {
    
    const params = {
            ChatId: chatId,
            MessageId: messageId,
            Userid: userId,
        }
    const url = BuildURLWithParams(DEL_MSG_ENDPOINT, params)
    
    return fetch(url, DELETE())
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("Delete message response data: ", data)
            return data
        })
        .catch(error => {
            console.log('Error:', error);
        });
}

async function EditChatNameRequest(name, chatId) {
    return fetch(EDIT_CHAT_NAME_ENDPOINT, POST({ Name: name, ChatId: chatId }))
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            return data
        })
        .catch(error => {
            console.log('Error:', error);
        });
}

async function LeaveChatRequest(chatId) {
    const url = BuildURLWithParams(LEAVE_CHAT_ENDPOINT, { ChatId: chatId })
    return fetch(url, DELETE())
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("leave chat response data: ", data)
            return data
        })
        .catch(error => {
            console.log('Error:', error);
        });
}

async function SwitchChatRequest(chatType, chatId) {
    const url = BuildURLWithParams(CHAT_SWITCH_ENDPOINT, { ChatType: chatType, ChatId: chatId })
    return fetch(url, GET())
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("leave chat response data: ", data)
            return data
        })
        .catch(error => {
            console.log('Error:', error);
        });
}

async function NewChatRequest(name) {
    return fetch(NEW_CHAT_ENDPOINT, POST({ Name: name }))
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("new chat response data: ", data)
            return data
        })
        .catch(error => {
            console.log('New chat request failed:', error);
        })
}

async function AddContactRequest(email) {
    return fetch(ADD_CONTACT_ENDPOINT, POST({ Email: email}))
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("add contact response data: ", data)
            return data
        })
        .catch(error => {
            console.log('add contact request failed:', error);
        })
}

async function NewMessageRequest(chatId, replyId, msgText) {
    return fetch(NEW_CHAT_ENDPOINT, POST({ ChatId: chatId, ReplyId: replyId, MsgText: msgText}))
        .then(response => {
            if (!response.ok) throw new Error("Network response was not ok");
            return response.json();
        })
        .then(data => {
            console.log("new message response data: ", data)
            return data
        })
        .catch(error => {
            console.log('New message request failed:', error);
        })
}