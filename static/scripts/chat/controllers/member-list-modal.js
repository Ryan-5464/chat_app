function ConfigureMemberListModal() {
    const modal = document.getElementById('chatMemberModal');
    
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        AddMemberToChat: (email, chatId) => AddMemberToChat(email, chatId),
    } 
    
    console.log("configure member list => modal:", modal)
    AddMemberModalToMemberListModal(modal) 

    document.addEventListener("click", (e) => {
    const modal = document.getElementById("chatMemberModal");
    if (!modal.classList.contains("open")) return; // only when open

    if (!e.target.closest('.modal-content')) {
        modal.__controller.Close();
    }
    });

    let addMemberInput = document.getElementById('add-member-input')
    addMemberInput = RemoveAllListeners(addMemberInput);
    addMemberInput.addEventListener('keydown', (e) => {
        if (e.key !== 'Enter') return;
        e.preventDefault()
        const chatId = addMemberInput.getAttribute('data-chatid')
        modal.__controller.AddMemberToChat(addMemberInput.value, chatId)
        addMemberInput.value = ''
    })

    return modal;
}

function AddMemberToChat(email, chatId) {
    AddMemberToChatHandler(email, chatId);
}