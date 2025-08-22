function addMessageModalListenerToMessageContainer() {
    const modal = ConfigureMessageModal();
    const container = document.getElementById('messages-container');
    const configureEditMsgButton = ConfigureEditMsgButton(modal.__controller);
    const configureDeleteMsgButton = ConfigureDeleteMsgButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();

        console.log("modal", modal)
        
        const editButton = document.getElementById('msg-edit-btn')
        const deleteButton = document.getElementById('msg-del-btn')
        editButton.classList.add('hidden')
        deleteButton.classList.add('hidden')

        const message = e.target.closest('.message')
        if (message.classList.contains('me')) {
            editButton.classList.remove('hidden')
            deleteButton.classList.remove('hidden')
        }
        const messageId = e.target.closest('[data-messageid]')?.getAttribute('data-messageid');
        const userId = e.target.closest('[data-userid]')?.getAttribute('data-userid');
        if (!messageId || !userId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditMsgButton(messageId);
        configureDeleteMsgButton(messageId, userId);
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
    let currentMsgId, currentUserId = null;
    deleteMessageButton.addEventListener('click', () => {
        if (!currentMsgId || !currentUserId) return;
        msgModalController.DeleteMessage(currentMsgId, currentUserId);
    });
    return (messageId, userId) => {currentMsgId = messageId; currentUserId = userId };
};

function ConfigureMessageModal() {
    const modal = document.getElementById('messageModal');
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditMessage: (messageId) => EditMessage(messageId, () => CloseModal(modal)),
        DeleteMessage: (chatId, messageId, userId) => DeleteMessage(chatId, messageId, userId, () => CloseModal(modal)),
    };
    return modal;
};

function EditMessage(messageId, closeModal) {
    const message = document.querySelector(`[data-messageid="${messageId}"]`);
    const msgText = message.querySelector('.message-text');
    const userId = message.getAttribute('data-userid');
    const openInput = document.getElementById('edit-message-input');
    if (openInput) {
        openInput.replaceWith(openInput.__oldtext)
    }
    const input = replaceWithInput(msgText, msgText.innerHTML, 'edit-message-input');
    input.focus();
    closeModal();
    input.addEventListener('keydown', (e) => {
        if (e.key == "Enter") {
            EditMessageHandler(input.value, messageId, userId);
        };
    });
    return;
};

function DeleteMessage(chatId, messageId, userId, closeModal) {
    closeModal();
    DeleteMessageHandler(chatId, messageId, userId);
};