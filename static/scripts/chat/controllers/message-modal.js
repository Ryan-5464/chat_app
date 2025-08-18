function addMessageModalListenerToMessageContainer() {
    const modal = ConfigureMessageModal();
    const container = document.getElementById('messages-container');
    const configureEditMsgButton = ConfigureEditMsgButton(modal.__controller);
    const configureDeleteMsgButton = ConfigureDeleteMsgButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();
        const chatId = e.target.closest('[data-chatid]')?.getAttribute('data-chatid');
        const messageId = e.target.closest('[data-messageid]')?.getAttribute('data-messageid');
        const userId = e.target.closest('[data-userid]')?.getAttribute('data-userid');
        if (!chatId || !messageId || !userId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditMsgButton(messageId);
        configureDeleteMsgButton(chatId, messageId, userId);
    });
};

function ConfigureEditMsgButton(msgModalController) {
    const editMsgButton = document.getElementById('msg-edit-btn');
    let currentMsgId = null;
    editMsgButton.addEventListener('click', () => {
        if (!currentMsgId) return;
        msgModalController.EditMessage(currentMsgId);
    });
    return (msgId) => { currentMsgId = msgId };
};


function ConfigureDeleteMsgButton(msgModalController) {
    const deleteMessageButton = document.getElementById('msg-del-btn');
    let currentChatId, currentMsgId, currentUserId = null;
    deleteMessageButton.addEventListener('click', () => {
        if (!currentChatId || !currentMsgId || !currentUserId) return;
        msgModalController.DeleteMessage(currentChatId, currentMsgId, currentUserId);
    });
    return (chatId, messageId, userId) => {currentChatId = chatId; currentMsgId = messageId; currentUserId = userId };
};

function ConfigureMessageModal() {
    const modal = document.getElementById('messageModal');
    const messageModalController = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditMessage: (messageId) => EditMessage(messageId, () => CloseModal(modal)),
        DeleteMessage: (chatId, messageId, userId) => DeleteMessage(chatId, messageId, userId, () => CloseModal(modal)),
    };
    modal.__controller = messageModalController;
    return modal;
};

function EditMessage(messageId, closeModal) {
    const message = document.querySelector(`[data-messageid="${messageId}"]`);
    const msgText = message.querySelector('.message-text');
    const input = replaceWithInput(msgText, msgText.innerHTML);
    input.focus();
    closeModal();
    input.addEventListener('keydown', (e) => {
        if (e.key == "Enter") {
            EditMessageHandler(input.value, messageId);
        };
    });
    return;
};

function DeleteMessage(chatId, messageId, userId, closeModal) {
    closeModal();
    DeleteMessageHandler(chatId, messageId, userId);
};