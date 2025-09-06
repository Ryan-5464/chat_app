window.APP = window.APP || {};

APP.URL = {
    BASE: 'http://localhost:8081',
    WEBSOCKET: 'ws://localhost:8081/ws',
}

APP.ENDPOINT = {
    CHAT: '/chat',
    LOGIN: '/api/login',
    EDIT_USERNAME: '/api/profile/name/edit',
    LEAVE_CHAT: '/api/chat/leave',
    DEL_MSG: '/api/message/delete',
    EDIT_CHAT_NAME: '/api/chat/edit',
    CHAT_SWITCH: '/api/chat/switch',
    NEW_CHAT: '/api/chat/new',
    ADD_CONTACT: '/api/chat/contact/add',
    REMOVE_CONTACT: '/api/chat/contact/remove',
    CONTACT_CHAT_SWITCH: '/api/chat/contact/switch',
    EDIT_MESSAGE: '/api/message/edit',
    GET_MEMBERS: '/api/chat/members',
    ADD_MEMBER: '/api/chat/members/add',
    REMOVE_MEMBER: '/api/chat/member/remove',
    ONLINE_STATUS: '/api/status',
    GET_ONLINE_STATUS: '/api/status/get',
}