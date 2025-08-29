const modal = ConfigureMemberModal()
const configureAddContactButton = ConfigureAddContactButton(modal.__controller)
const configureRemoveMemberButton = ConfigureRemoveMemberButton(modal.__controller)



function AddMemberModalToMemberListModal(memberListModal) {

    memberListModal.addEventListener("contextmenu", (e) => {
        e.preventDefault()
        e.stopPropagation()

        const member = e.target.closest(`.${APP.CLS.MEMBER.TAG}`)
        if (!member) return
        const email  = member.children[1].innerHTML
        const chatId = GetDataAttribute(member, APP.DATA.CHAT.ID) 
        const userId = GetDataAttribute(member, APP.DATA.USER.ID)
        if (!chatId || !userId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY)
        configureAddContactButton(email)
        configureRemoveMemberButton(chatId, userId)
    })
}


function ConfigureAddContactButton(memberModalController) {
    let addContactBtn = document.getElementById(APP.ID.MODAL.MEMBERLIST.MEMBER.BTN.ADD_CONTACT)
    let currentEmail = null
    addContactBtn = RemoveAllListeners(addContactBtn);
    addContactBtn.addEventListener('click', (e) => {
        e.preventDefault();
        if (!currentEmail) return;
        memberModalController.AddContact(currentEmail);
    });
    return (email) => { currentEmail = email }
}

function ConfigureRemoveMemberButton(memberModalController) {
    let removeMemberButton = document.getElementById(APP.ID.MODAL.MEMBERLIST.MEMBER.BTN.REMOVE_MEMBER)
    let currentChatId, currentUserId = null
    removeMemberButton = RemoveAllListeners(removeMemberButton);
    removeMemberButton.addEventListener('click', (e) => {
        if (!currentChatId || !currentUserId) return;
        e.preventDefault()
        memberModalController.RemoveMember(currentChatId, currentUserId)
    })
    return (chatId, userId) => { currentChatId = chatId; currentUserId = userId }
}

function ConfigureMemberModal() {
    const modal = document.getElementById(APP.ID.MODAL.MEMBERLIST.MEMBER.MODAL)
    modal.__controller = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        AddContact: (email) => AddContact(email, () => CloseModal(modal)),
        RemoveMember: (chatId, userId) => RemoveMember(chatId, userId, () => CloseModal(modal)),
    }
    return modal
}

function AddContact(email, closeModal) {
    closeModal()
    AddContactHandler(email)
}

function RemoveMember(chatId, userId, closeModal) {
    closeModal()
    RemoveMemberHandler(chatId, userId)
}