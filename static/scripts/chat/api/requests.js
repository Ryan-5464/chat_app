const AddContactRequest = (email) => safeRequest(() => POST(APP.ENDPOINT.ADD_CONTACT, { Email: email }));

const RemoveContactRequest = (contactId) => safeRequest(() => DELETE(APP.ENDPOINT.REMOVE_CONTACT, { ContactId: contactId }));

const DeleteMessageRequest = (messageId, userId) => safeRequest(() => DELETE(APP.ENDPOINT.DEL_MSG, { MessageId: messageId, UserId: userId }));

const EditChatNameRequest = (newName, chatId) => safeRequest(() => POST(APP.ENDPOINT.EDIT_CHAT_NAME, { Name: newName, ChatId: chatId}));

const LeaveChatRequest = (chatId) => safeRequest(() => DELETE(APP.ENDPOINT.LEAVE_CHAT, { ChatId: chatId}));

const NewChatRequest = (newChatName) => safeRequest(() => POST(APP.ENDPOINT.NEW_CHAT, { Name: newChatName }));

const SwitchChatRequest = (chatId) => safeRequest(() => GET(APP.ENDPOINT.CHAT_SWITCH, { ChatId: chatId}));

const SwitchContactChatRequest = (contactChatId) => safeRequest(() => GET(APP.ENDPOINT.CONTACT_CHAT_SWITCH, { ContactChatId: contactChatId }));

const EditMessageRequest = (messageText, messageId, userId) => safeRequest(() => POST(APP.ENDPOINT.EDIT_MESSAGE, { MsgText: messageText, MessageId: messageId, UserId: userId }));

const GetMemberListRequest = (chatId) => safeRequest(() => GET(APP.ENDPOINT.GET_MEMBERS, { ChatId: chatId }));

const AddMemberToChatRequest = (email, chatId) => safeRequest(() => POST(APP.ENDPOINT.ADD_MEMBER, { Email: email, ChatId: chatId }));

const RemoveMemberRequest = (chatId, userId) => safeRequest(() => DELETE(APP.ENDPOINT.REMOVE_MEMBER, { ChatId: chatId, UserId: userId }));

const ChangeOnlineStatusRequest = (status) => safeRequest(() => GET(APP.ENDPOINT.ONLINE_STATUS, { Status: status }));

const GetOnlineStatusRequest = () => safeRequest(() => GET(APP.ENDPOINT.GET_ONLINE_STATUS, {}));

async function safeRequest(reqFunc) {
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




