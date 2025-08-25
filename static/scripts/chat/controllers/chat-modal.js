function addChatModalListenerToChatContainer() {
    const modal = ConfigureChatModal();
    const container = document.getElementById('chats-container');
    const configureEditButton = ConfigureEditButton(modal.__controller);
    const configureLeaveButton = ConfigureLeaveButton(modal.__controller);
    const configureMembersButton = ConfigureMembersButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();
        const chatId = e.target.closest('[data-chatid]')?.getAttribute('data-chatid');
        if (!chatId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureEditButton(chatId);
        configureLeaveButton(chatId);
        configureMembersButton(chatId);
    });
};

function ConfigureEditButton(chatModalController) {
    const editNameButton = document.getElementById('chat-edit-btn');
    let currentChatId = null;
    editNameButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.EditChatName(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureLeaveButton(chatModalController) {
    const leaveButton = document.getElementById('chat-leave-btn');
    let currentChatId = null;
    leaveButton.addEventListener('click', () => {
        if (!currentChatId) return;
        chatModalController.LeaveChat(currentChatId);
    });
    return (chatId) => { currentChatId = chatId };
};

function ConfigureMembersButton(chatModalController) {
    const membersButton = document.getElementById('chat-members-btn');
    let currentChatId = null;

    membersButton.addEventListener('click', () => {
        if (!currentChatId) return;
        const addMemberInput = document.getElementById('add-member-input')
        addMemberInput.setAttribute('data-chatid', currentChatId)
        chatModalController.DisplayMemberList(currentChatId);
    });
    return (chatId) => {currentChatId = chatId}
}

function ConfigureChatModal() {
    const modal = document.getElementById('chatModal');
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
    const chat = document.querySelector(`[data-chatid="${chatId}"]`);
    const editChatNameHandler = (inputText) => { EditChatNameHandler(inputText, chatId) }
    closeModal();
    textInputController(chat, editChatNameHandler, 'chat-name-input', 'chat-name')
};

function LeaveChat(chatId, closeModal) {
    const activeChat = document.querySelector('.active');
    const activeChatId = activeChat?.getAttribute("data-chatid");
    const isActive = chatId === activeChatId;
    closeModal();
    console.log("is active", isActive)
    LeaveChatHandler(chatId, isActive);
};

function DisplayMemberList(chatId, closeModal) {
    closeModal(); // closes the chat options modal
    const memberModal = ConfigureMemberListModal();
    memberModal.__controller.OpenAt();

    // Fetch & render members
    DisplayMemberListHandler(chatId);
}