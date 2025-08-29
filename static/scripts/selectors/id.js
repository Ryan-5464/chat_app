window.APP = window.APP || {};
APP.ID = APP.ID || {}; 

APP.ID.GEN = {
    PFL_USERNAME: 'username',
    EMAIL: 'email',
    ERROR: 'error',
    PASSWORD: 'password',
    STATUS: 'online-status',
}

APP.ID.CHAT = {
    CONTAINER: 'chats-container',
    INPUT: {
        EDIT_NAME: 'chat-name-input',
        NEW_CHAT: 'new-chat-input',
    },
}

APP.ID.CONTACT = {
    CONTAINER: 'contacts-container',
    INPUT: {
        NEW_CONTACT_EMAIL: 'contact-email-input',
    },
}

APP.ID.MEMBER = {
}

APP.ID.MESSAGE = {
    CONTAINER: 'messages-container',
    INPUT: {
        EDIT_MSG: 'edit-message-input',
        NEW_MSG: 'message-input',
    },
    BTN: {
        SEND: 'send-message-button', 
    },
}

APP.ID.MODAL = {
    CHAT: {
        MODAL: 'chatModal',
        BTN : {
            EDIT: 'chat-edit-btn',
            LEAVE: 'chat-leave-btn',
            MEMBER: 'chat-members-btn',
        },
    },
    CONTACT: {
        MODAL: 'contactModal',
        BTN: {
            REMOVE_CONTACT: 'contact-remove-btn',
        },

    },
    MEMBERLIST: {
        MODAL: 'chatMemberModal',
        CONTAINER: 'member-list-container',

        INPUT: {
            ADD_MEMBER: 'add-member-input',
        },
        TITLE: {
            ADD_MEMBER: 'add-member',
        },
        MEMBER: {
            MODAL: 'memberModal',
            BTN : {
                REMOVE_MEMBER: 'member-del-btn',
                ADD_CONTACT: 'add-contact-btn',
            },
        },
    },
    MESSAGE: {
        MODAL: 'messageModal',
        BTN: {
            EDIT_MSG: 'msg-edit-btn',
            DELETE_MSG: 'msg-del-btn',
        }
    }
}

APP.ID.USER = {
    NAME: 'username-input',
}