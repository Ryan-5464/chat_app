function addContactModalListenerToContactContainer() {
    const modal = ConfigureContactModal();
    const container = document.getElementById('contacts-container');
    // const configureEditButton = ConfigureEditButton(modal.__controller);
    const configureRemoveContactButton = ConfigureRemoveContactButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();
        const contactId = e.target.closest('[data-contactid]')?.dataset.contactid;
        if (!contactId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureRemoveContactButton(contactId);
    });
};

// function ConfigureEditButton(contactModalController) {
//     const editNameButton = document.getElementById('contact-edit-btn');
//     let currentContactId = null;
//     editNameButton.addEventListener('click', () => {
//         if (!currentContactId) return;
//         contactModalController.EditContactName(currentContactId);
//     });
//     return (contactId) => { currentContactId = contactId };
// };

function ConfigureRemoveContactButton(contactModalController) {
    const removeButton = document.getElementById('contact-remove-btn');
    let currentContactId = null;
    removeButton.addEventListener('click', () => {
        if (!currentContactId) return;
        contactModalController.RemoveContact(currentContactId);
    });
    return (contactId) => { currentContactId = contactId };
};

function ConfigureContactModal() {
    const modal = document.getElementById('contactModal');
    const contactModalController = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        // EditContactName: (contactId) => EditContactName(contactId, () => CloseModal(modal)),
        RemoveContact: (contactId) => RemoveContact(contactId, () => CloseModal(modal)),
    };
    modal.__controller = contactModalController;
    return modal;
};

// function EditContactName(contactId, closeModal) {
//     const openInput = document.getElementById('contact-name-input')
//     if (openInput) {
//         const openInputContactName = document.createElement('div')
//         openInputContactName.classList.add('contact-name')
//         openInputContactName.innerHTML = openInput.__oldtext
//         openInput.replaceWith(openInputContactName)
//     }
//     if (openInput) { openInput.remove()}
//     const contact = document.querySelector(`[data-contactid="${contactId}"]`);
//     const contactName = contact.querySelector('.contact-name');
//     const input = replaceWithInput(contactName, "Enter new name");
//     input.focus();
//     closeModal();
//     input.addEventListener('keydown', (e) => {
//         if (e.key === "Enter") {
//             EditContactNameHandler(input.value, contactId);
//         };
//     });
// };

function RemoveContact(contactId, closeModal) {
    const activeContactChat = document.querySelector('.active');
    const activeContactId = activeContactChat?.dataset.contactid;
    const isActive = contactId === activeContactId;
    closeModal();
    console.log("is active", isActive);
    RemoveContactHandler(contactId, isActive);
};