function addMessageModalListenerToMessageContainer() {
    const modal = ConfigureMessageModal();
    const container = document.getElementById(APP.ID.MESSAGE.CONTAINER);
    const configureEditMsgButton = ConfigureEditMsgButton(modal.__controller);
    const configureDeleteMsgButton = ConfigureDeleteMsgButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();

        const editButton = document.getElementById(APP.ID.MODAL.MSG_EDIT_BTN);
        const deleteButton = document.getElementById(APP.ID.MODAL.MSG_DEL_BTN);
        editButton.classList.add(APP.CLS.GEN.HIDDEN);
        deleteButton.classList.add(APP.CLS.GEN.HIDDEN);

        const message = e.target.closest(`.${APP.CLS.MESSAGE.TAG}`);
        if (message.classList.contains(APP.CLS.GEN.ME)) {
            editButton.classList.remove(APP.CLS.GEN.HIDDEN);
            deleteButton.classList.remove(APP.CLS.GEN.HIDDEN);
        }
        const m = GetClosestTargetByData(e, APP.DATA.MESSAGE.ID)
        if (!m) { return; }
        const messageId = GetDataAttribute(m, APP.DATA.MESSAGE.ID);
        const u = GetClosestTargetByData(e, APP.DATA.USER.ID);
        if (!u) { return; }
        const userId = GetDataAttribute(u, APP.DATA.USER.ID);
        if (!messageId || !userId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditMsgButton(messageId);
        configureDeleteMsgButton(messageId, userId);
    });
};

function ConfigureEditMsgButton(msgModalController) {
    const editMsgButton = document.getElementById(APP.ID.MODAL.MSG_EDIT_BTN);
    let currentMsgId = null;
    editMsgButton.addEventListener('click', () => {
        if (!currentMsgId) return;
        msgModalController.EditMessage(currentMsgId);
    });
    return (msgId) => { currentMsgId = msgId };
};


function ConfigureDeleteMsgButton(msgModalController) {
    const deleteMessageButton = document.getElementById(APP.ID.MODAL.MSG_DEL_BTN);
    let currentMsgId, currentUserId = null;
    deleteMessageButton.addEventListener('click', () => {
        if (!currentMsgId || !currentUserId) return;
        msgModalController.DeleteMessage(currentMsgId, currentUserId);
    });
    return (messageId, userId) => {currentMsgId = messageId; currentUserId = userId };
};

function ConfigureMessageModal() {
    const modal = document.getElementById(APP.ID.MODAL.MESSAGE);
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditMessage: (messageId) => EditMessage(messageId, () => CloseModal(modal)),
        DeleteMessage: (chatId, messageId, userId) => DeleteMessage(chatId, messageId, userId, () => CloseModal(modal)),
    };
    return modal;
};

function EditMessage(messageId, closeModal) {
    const message = GetElemByDataTag(APP.DATA.MESSAGE.ID, messageId);
    const userId = GetDataAttribute(message, APP.DATA.USER.ID);
    const editMessageHandler = (inputText) => { EditMessageHandler(inputText, messageId, userId) }
    closeModal();
    textInputController(message, editMessageHandler, APP.ID.MESSAGE.EDIT_INPUT, APP.CLS.MESSAGE.TEXT, true)
};

function DeleteMessage(chatId, messageId, userId, closeModal) {
    closeModal();
    DeleteMessageHandler(chatId, messageId, userId);
};