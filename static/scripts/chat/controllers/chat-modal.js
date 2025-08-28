function addChatModalListenerToChatContainer() {
    const modal = ConfigureChatModal();
    const container = document.getElementById(APP.ID.CHAT.CONTAINER);
    const configureEditButton = ConfigureEditButton(modal.__controller);
    const configureLeaveButton = ConfigureLeaveButton(modal.__controller);
    const configureMembersButton = ConfigureMembersButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();

        const editButton = document.getElementById(APP.ID.CHAT.EDIT_BTN)
        const chat = e.target.closest(".".concat(APP.CLS.CHAT))

        if (chat.classList.contains(APP.CLS.ME)) {
            editButton.classList.remove(APP.CLS.HIDDEN)
        } else {
            editButton.classList.add(APP.CLS.HIDDEN) 
        }

        const c = GetClosestTargetByData(e, APP.DATA.CHAT.ID)
        if (!c) { return; }
        const chatId = GetDataAttribute(chat, APP.DATA.CHAT.ID);
        if (!chatId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditButton(chatId);
        configureLeaveButton(chatId);
        configureMembersButton(chatId, chat);
    });
};

function ConfigureEditButton(chatModalController) {
    const editNameButton = document.getElementById(APP.ID.CHAT.EDIT_BTN);
    let currentChatId = null;
    editNameButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.EditChatName(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureLeaveButton(chatModalController) {
    const leaveButton = document.getElementById(APP.ID.CHAT.LEAVE_BTN);
    let currentChatId = null;
    leaveButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.LeaveChat(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureMembersButton(chatModalController) {
    const membersButton = document.getElementById(APP.ID.CHAT.MEMBER_BTN);
    let currentChatId = null;
    let chatElem = null;

    membersButton.addEventListener('click', () => {
        if (!currentChatId) return;
        const addMember = document.getElementById(APP.ID.CHAT.ADD)
        
        if (chatElem.classList.contains(APP.CLS.ME)) {
            addMember.classList.remove(APP.CLS.HIDDEN)
        } else {
            addMember.classList.add(APP.CLS.HIDDEN) 
        }

        const addMemberInput = document.getElementById(APP.ID.CHAT.INPUT)
        addMemberInput.setAttribute(`data-${APP.DATA.CHAT.ID}`, currentChatId)
        chatModalController.DisplayMemberList(currentChatId);
    });
    return (chatId, chat) => {currentChatId = chatId; chatElem = chat}
}

function ConfigureChatModal() {
    const modal = document.getElementById(APP.ID.MODAL.CHAT);
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        EditChatName: (chatId) => EditChatName(chatId, () => CloseModal(modal)),
        LeaveChat: (chatId) => LeaveChat(chatId, () => CloseModal(modal)),
        DisplayMemberList: (chatId) => DisplayMemberList(chatId, () => CloseModal(modal)),
    };
    return modal;
};

function EditChatName(chatId, closeModal) {
    const chat = GetElemByDataTag(APP.DATA.CHAT.ID, chatId);
    const editChatNameHandler = (inputText) => { EditChatNameHandler(inputText, chatId) }
    closeModal();
    textInputController(chat, editChatNameHandler, APP.ID.CHAT.NAME_INPUT, APP.CLS.CHAT.NAME)
};

function LeaveChat(chatId, closeModal) {
    const activeChat = QSelectByClass(document, APP.CLS.ACTIVE);
    if (!activeChat) { return; }
    const activeChatId = GetDataAttribute(activeChat, APP.DATA.CHAT.ID);
    const isActive = chatId === activeChatId;
    closeModal();
    console.log("is active", isActive)
    LeaveChatHandler(chatId, isActive);
};

function DisplayMemberList(chatId, closeModal) {
    closeModal(); 
    const memberModal = ConfigureMemberListModal();
    memberModal.__controller.OpenAt();
    DisplayMemberListHandler(chatId);
}