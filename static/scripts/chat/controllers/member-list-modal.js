function ConfigureMemberListModal() {
    const modal = document.getElementById(APP.ID.MODAL.CHAT_MEMBER);
    
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        AddMemberToChat: (email, chatId) => AddMemberToChat(email, chatId),
    };
    
    AddMemberModalToMemberListModal(modal);

    document.addEventListener("click", (e) => {
        const modal = document.getElementById(APP.ID.MODAL.CHAT_MEMBER);
        if (!modal.classList.contains(APP.CLS.OPEN)) return; 

        if (!e.target.closest(`.${APP.CLS.MODAL.CONTENT}`)) {
            modal.__controller.Close();
        };
    });

    let addMemberInput = document.getElementById(APP.ID.MEMBER.INPUT);
    addMemberInput = RemoveAllListeners(addMemberInput);
    addMemberInput.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter') return;
        e.preventDefault();
        const chatId = GetDataAttribute(addMemberInput, APP.DATA.CHAT.ID);
        modal.__controller.AddMemberToChat(addMemberInput.value, chatId);
        addMemberInput.value = '';
    });

    return modal;
};

function AddMemberToChat(email, chatId) {
    AddMemberToChatHandler(email, chatId);
};