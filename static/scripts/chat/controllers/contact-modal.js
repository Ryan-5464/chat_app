function addContactModalListenerToContactContainer() {
    const modal = ConfigureContactModal();
    const container = document.getElementById(APP.CLS.CONTACT.CONTAINER);
    const configureRemoveContactButton = ConfigureRemoveContactButton(modal.__controller);
    
    container.addEventListener("contextmenu", (e) => {
        e.preventDefault();
        e.stopPropagation();
        const contact = GetClosestTargetByData(e, APP.DATA.CONTACT.ID);
        const contactId = GetDataAttribute(contact, APP.DATA.CONTACT.ID);
        if (!contactId) return;
        modal.__controller.OpenAt(e.clientX, e.clientY);
        configureRemoveContactButton(contactId);
    });
};

function ConfigureRemoveContactButton(contactModalController) {
    const removeButton = document.getElementById(APP.ID.CONTACT.REMOVE_BTN);
    let currentContactId = null;
    removeButton.addEventListener('click', () => {
        if (!currentContactId) return;
        contactModalController.RemoveContact(currentContactId);
    });
    return (contactId) => { currentContactId = contactId };
};

function ConfigureContactModal() {
    const modal = document.getElementById(APP.ID.MODAL.CONTACT);
    const contactModalController = {
        Close: () => CloseModal(modal),
        OpenAt: (clientX, clientY) => OpenModalAt(modal, clientX, clientY),
        RemoveContact: (contactId) => RemoveContact(contactId, () => CloseModal(modal)),
    };
    modal.__controller = contactModalController;
    return modal;
};

function RemoveContact(contactId, closeModal) {
    const activeContactChat = QSelectByClass(document, APP.CLS.ACTIVE);
    const activeContactId = GetDataAttribute(activeContactChat, APP.DATA.CONTACT.ID);
    const isActive = contactId === activeContactId;
    closeModal();
    console.log("is active", isActive);
    RemoveContactHandler(contactId, isActive);
};