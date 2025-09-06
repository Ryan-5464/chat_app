function addMessageModalListenerToMessageContainer() {
    const modal = ConfigureMessageModal();
    const container = document.getElementById(APP.ID.MESSAGE.CONTAINER);
    const configureEditMsgButton = ConfigureEditMsgButton(modal.__controller);
    const configureDeleteMsgButton = ConfigureDeleteMsgButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();

        const message = e.target.closest(`.${APP.CLS.MESSAGE.TAG}`);
        if (!message) { return; }
        if (!message.classList.contains(APP.CLS.GEN.ME)) { return; };

        const editButton = document.getElementById(APP.ID.MODAL.MESSAGE.BTN.EDIT_MSG);
        const deleteButton = document.getElementById(APP.ID.MODAL.MESSAGE.BTN.DELETE_MSG);
        editButton.classList.add(APP.CLS.GEN.HIDDEN);
        deleteButton.classList.add(APP.CLS.GEN.HIDDEN);

        editButton.classList.remove(APP.CLS.GEN.HIDDEN);
        deleteButton.classList.remove(APP.CLS.GEN.HIDDEN);

        const m = GetClosestTargetByData(e, APP.DATA.MESSAGE.ID)
        if (!m) { return; }
        const messageId = GetDataAttribute(m, APP.DATA.MESSAGE.ID);

        const u = GetClosestTargetByData(e, APP.DATA.USER.ID);
        if (!u) { return; }
        const userId = GetDataAttribute(u, APP.DATA.USER.ID);

        const c = GetClosestTargetByData(e, APP.DATA.CHAT.ID);
        if (!c) { return; }
        const chatId = GetDataAttribute(c, APP.DATA.CHAT.ID);

        if (!messageId || !userId || !chatId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditMsgButton(messageId);
        configureDeleteMsgButton(messageId, userId, chatId);
    });
};

function ConfigureEditMsgButton(msgModalController) {
    const editMsgButton = document.getElementById(APP.ID.MODAL.MESSAGE.BTN.EDIT_MSG);
    let currentMsgId = null;
    editMsgButton.addEventListener('click', () => {
        if (!currentMsgId) return;
        msgModalController.EditMessage(currentMsgId);
    });
    return (msgId) => { currentMsgId = msgId };
};


function ConfigureDeleteMsgButton(msgModalController) {
    const deleteMessageButton = document.getElementById(APP.ID.MODAL.MESSAGE.BTN.DELETE_MSG);
    let currentMsgId = null;
    let currentUserId = null;
    let currentChatId = null;
    deleteMessageButton.addEventListener('click', () => {
        if (!currentMsgId || !currentUserId || !currentChatId) return;
        msgModalController.DeleteMessage(currentMsgId, currentUserId, currentChatId);
    });
    return (messageId, userId, chatId) => {currentMsgId = messageId; currentUserId = userId; currentChatId = chatId };
};

function ConfigureMessageModal() {
    const modal = document.getElementById(APP.ID.MODAL.MESSAGE.MODAL);
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditMessage: (messageId) => EditMessage(messageId, () => CloseModal(modal)),
        DeleteMessage: (messageId, userId, chatId) => DeleteMessage(messageId, userId, chatId, () => CloseModal(modal)),
    };
    return modal;
};

function EditMessage(messageId, closeModal) {
    const message = GetElemByDataTag(document, APP.DATA.MESSAGE.ID, messageId);
    const userId = GetDataAttribute(message, APP.DATA.USER.ID);
    const editMessageHandler = (inputText) => { EditMessageHandler(inputText, messageId, userId) }
    closeModal();
    textInputController(message, editMessageHandler, APP.ID.MESSAGE.INPUT.EDIT_MSG, APP.CLS.MESSAGE.TEXT, true)
};

function DeleteMessage(messageId, userId, chatId, closeModal) {
    console.log("CHATID, ", chatId)
    closeModal();
    DeleteMessageHandler(messageId, userId, chatId);
};