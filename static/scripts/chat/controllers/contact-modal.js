function addContactModalListenerToContactContainer() {
    const modal = ConfigureContactModal();
    const container = document.getElementById('contacts-container');
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
        RemoveContact: (contactId) => RemoveContact(contactId, () => CloseModal(modal)),
    };
    modal.__controller = contactModalController;
    return modal;
};

function RemoveContact(contactId, closeModal) {
    const activeContactChat = document.querySelector('.active');
    const activeContactId = activeContactChat?.dataset.contactid;
    const isActive = contactId === activeContactId;
    closeModal();
    console.log("is active", isActive);
    RemoveContactHandler(contactId, isActive);
};