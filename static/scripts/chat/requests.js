const LEAVE_CHAT_ENDPOINT = '/api/chat/leave';
const DEL_MSG_ENDPOINT = '/api/message/delete';
const EDIT_CHAT_NAME_ENDPOINT = '/api/chat/edit';
const CHAT_SWITCH_ENDPOINT = '/api/chat/switch';
const NEW_CHAT_ENDPOINT = '/api/chat/new';
const ADD_CONTACT_ENDPOINT = '/api/chat/contact/add';
const REMOVE_CONTACT_ENDPOINT = 'api/chat/contact/remove';
const CONTACT_CHAT_SWITCH_ENDPOINT = 'api/chat/contact/switch';
const EDIT_MESSAGE_ENDPOINT = 'api/message/edit';

const AddContactRequest = (email) => safeRequest(() => POST(ADD_CONTACT_ENDPOINT, { Email: email }));

const RemoveContactRequest = (contactId) => safeRequest(() => DELETE(REMOVE_CONTACT_ENDPOINT, { ContactId: contactId }));

const DeleteMessageRequest = (chatId, messageId, userId) => safeRequest(() => DELETE(DEL_MSG_ENDPOINT, { ChatId: chatId, MessageId: messageId, UserId: userId }));

const EditChatNameRequest = (newName, chatId) => safeRequest(() => POST(EDIT_CHAT_NAME_ENDPOINT, { Name: newName, ChatId: chatId}));

const LeaveChatRequest = (chatId) => safeRequest(() => DELETE(LEAVE_CHAT_ENDPOINT, { ChatId: chatId}));

const NewChatRequest = (newChatName) => safeRequest(() => POST(NEW_CHAT_ENDPOINT, { Name: newChatName }));

const SwitchChatRequest = (chatId) => safeRequest(() => GET(CHAT_SWITCH_ENDPOINT, { ChatId: chatId}));

const SwitchContactChatRequest = (contactChatId) => safeRequest(() => GET(CONTACT_CHAT_SWITCH_ENDPOINT, { ContactChatId: contactChatId }));

const EditMessageRequest = (messageText, messageId, userId) => safeRequest(() => POST(EDIT_MESSAGE_ENDPOINT, { MsgText: messageText, MessageId: messageId, UserId: userId }));

async function safeRequest(reqFunc) {
    console.log("safe request: ", reqFunc)
    try {
        const responseJSON = await reqFunc();
        if (!responseJSON || Object.keys(responseJSON).length === 0) {
            throw new Error("No response data.");
        };
        return responseJSON;
    } catch (error) {
        console.error(`Request failed => error: `, error);
        return null;
    };
};

async function request(endpoint, options) {
    const response = await fetch(endpoint, options);
    if (!response.ok) {
        const text = await response.text();
        throw new Error(`Network response was not ok: ${response.status} - ${text}`);
    };
    return response.json();
};

async function GET(endpoint, params) {
    const paramsEndpoint = EndpointWithParams(endpoint, params);
    return request(paramsEndpoint, { method: "GET" });
};

async function DELETE(endpoint, params) {
    const paramsEndpoint = EndpointWithParams(endpoint, params);
    return request(paramsEndpoint, { method: 'DELETE' });
};

async function POST(endpoint, payload) {
    console.log("payload:", payload)
    return request(endpoint, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    });
};

function EndpointWithParams(endpoint, params) {
    const query = new URLSearchParams(params).toString();
    return query ? `${endpoint}?${query}` : endpoint;
};





// async function NewMessageRequest(chatId, replyId, msgText) {
//     return fetch(NEW_CHAT_ENDPOINT, POST({ ChatId: chatId, ReplyId: replyId, MsgText: msgText}))
//         .then(response => {
//             if (!response.ok) throw new Error("Network response was not ok");
//             return response.json();
//         })
//         .then(data => {
//             console.log("new message response data: ", data)
//             return data
//         })
//         .catch(error => {
//             console.log('New message request failed:', error);
//         })
// }
